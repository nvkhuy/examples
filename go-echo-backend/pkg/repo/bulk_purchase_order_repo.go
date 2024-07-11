package repo

import (
	"bytes"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/engineeringinflow/inflow-backend/pkg/excel"
	"github.com/engineeringinflow/inflow-backend/pkg/helper"
	"github.com/engineeringinflow/inflow-backend/pkg/runner"
	"github.com/engineeringinflow/inflow-backend/pkg/s3"
	"github.com/lib/pq"
	"github.com/samber/lo"

	"github.com/engineeringinflow/inflow-backend/pkg/db"
	"github.com/engineeringinflow/inflow-backend/pkg/errs"
	"github.com/engineeringinflow/inflow-backend/pkg/models"
	"github.com/engineeringinflow/inflow-backend/pkg/models/enums"
	"github.com/engineeringinflow/inflow-backend/pkg/models/price"
	"github.com/engineeringinflow/inflow-backend/pkg/repo/query"
	"github.com/engineeringinflow/inflow-backend/pkg/repo/query/queryfunc"
	"github.com/engineeringinflow/inflow-backend/pkg/stripehelper"
	"github.com/jinzhu/copier"
	"github.com/rotisserie/eris"
	"github.com/stripe/stripe-go/v74"
	stripePrice "github.com/stripe/stripe-go/v74/price"
	"github.com/thaitanloi365/go-utils"
	"github.com/thaitanloi365/go-utils/values"
	"golang.org/x/exp/slices"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type BulkPurchaseOrderRepo struct {
	db *db.DB
}

func NewBulkPurchaseOrderRepo(db *db.DB) *BulkPurchaseOrderRepo {
	return &BulkPurchaseOrderRepo{
		db: db,
	}
}

type PaginateBulkPurchaseOrderParams struct {
	models.PaginationParams
	models.JwtClaimsInfo

	TeamID   string `json:"team_id" query:"team_id"`
	UserID   string `json:"user_id" query:"user_id"`
	SellerID string `json:"seller_id" query:"seller_id"`

	Statuses                []enums.BulkPurchaseOrderStatus       `json:"statuses" query:"statuses"`
	TrackingStatuses        []enums.BulkPoTrackingStatus          `json:"tracking_statuses" query:"tracking_statuses"`
	SellerTrackingStatuses  []enums.SellerBulkPoTrackingStatus    `json:"seller_tracking_statuses" query:"seller_tracking_statuses"`
	AssigneeIDs             []string                              `json:"assignee_ids" query:"assignee_ids"`
	AssigneeID              string                                `json:"assignee_id" query:"assignee_id"`
	PostedDateFrom          int64                                 `json:"posted_date_from" query:"posted_date_from"`
	PostedDateTo            int64                                 `json:"posted_date_to" query:"posted_date_to"`
	SellerQuotationStatuses []enums.BulkPurchaseOrderSellerStatus `json:"seller_quotation_statuses" query:"seller_quotation_statuses"`

	IncludeUser       bool `json:"-"`
	IncludeAssignee   bool `json:"-"`
	IncludeTrackings  bool `json:"-"`
	IsQueryAll        bool `json:"-"`
	IncludeCollection bool `json:"-"`
}

func (r *BulkPurchaseOrderRepo) PaginateBulkPurchaseOrder(params PaginateBulkPurchaseOrderParams) *query.Pagination {
	var userID = params.GetUserID()
	if params.TeamID != "" && !params.GetRole().IsAdmin() && !params.IsQueryAll {
		if err := r.db.Select("ID").First(&models.BrandTeam{}, "team_id = ? AND user_id = ?", params.TeamID, userID).Error; err != nil {
			return &query.Pagination{
				Records: []*models.BulkPurchaseOrder{},
			}
		}
		userID = params.TeamID
	}
	var result = query.New(r.db, queryfunc.NewBulkPurchaseOrderBuilder(queryfunc.BulkPurchaseOrderBuilderOptions{
		IncludeUser:       params.IncludeUser,
		IncludeAssignee:   params.IncludeAssignee,
		IncludeTrackings:  params.IncludeTrackings,
		IncludeCollection: params.IncludeCollection,
		QueryBuilderOptions: queryfunc.QueryBuilderOptions{
			Role: params.GetRole(),
		},
	})).
		WhereFunc(func(builder *query.Builder) {
			if !params.IsQueryAll {
				if params.GetRole().IsAdmin() {
					if params.UserID != "" {
						builder.Where("bpo.user_id = ?", params.UserID)
					}
				} else if params.GetRole().IsSeller() {

				} else {
					builder.Where("bpo.user_id = ? AND bpo.deleted_at IS NULL", userID)
				}
			}

			if params.SellerID != "" {
				builder.Where("bpo.seller_id = ?", params.SellerID)
			}

			if len(params.TrackingStatuses) > 0 {
				builder.Where("bpo.tracking_status IN ?", params.TrackingStatuses)
			}
			if len(params.SellerTrackingStatuses) > 0 {
				builder.Where("bpo.seller_tracking_status IN ?", params.SellerTrackingStatuses)
			}

			if len(params.Statuses) > 0 {
				builder.Where("bpo.status IN ?", params.Statuses)
			}

			if params.PostedDateFrom > 0 {
				builder.Where("bpo.created_at >= ?", params.PostedDateFrom)
			}

			if params.PostedDateTo > 0 {
				builder.Where("bpo.created_at <= ?", params.PostedDateTo)
			}

			if params.AssigneeID != "" {
				builder.Where("count_elements(bpo.assignee_ids,?) >= 1", pq.StringArray([]string{params.AssigneeID}))
			}

			if len(params.AssigneeIDs) > 0 {
				builder.Where("count_elements(bpo.assignee_ids,?) >= 1", pq.StringArray(params.AssigneeIDs))
			}

			if keyword := strings.TrimSpace(params.Keyword); keyword != "" {
				var q = "%" + keyword + "%"
				if strings.HasPrefix(keyword, "BPO-") {
					builder.Where("bpo.reference_id ILIKE ?", q)
				} else if strings.HasPrefix(keyword, "GBPO-") {
					builder.Where("bpo.group_id ILIKE ?", q)
				} else if strings.HasPrefix(keyword, "IQ-") {
					builder.Where("iq.reference_id ILIKE ?", q)
				} else {
					builder.Where("(bpo.id ILIKE @keyword OR iq.title ILIKE @keyword OR iq.id ILIKE @keyword)", sql.Named("keyword", q))
				}
			}

			if params.GetRole().IsSeller() {
				builder.Where("(bpo.seller_id = @user_id OR bposq.id IS NOT NULL) AND bpo.deleted_at IS NULL", sql.Named("user_id", userID))
				if len(params.SellerQuotationStatuses) > 0 {
					builder.Where("bposq.status IN @SellerQuotationStatuses", sql.Named("SellerQuotationStatuses", params.SellerQuotationStatuses))
				}
			}
		}).
		OrderBy("bpo.order_group_id ASC, bpo.updated_at DESC").
		Limit(params.Limit).
		Page(params.Page).
		PagingFunc()

	return result
}

type GetBulkPurchaseOrderParams struct {
	models.JwtClaimsInfo

	BulkPurchaseOrderID string `json:"bulk_purchase_order_id" param:"bulk_purchase_order_id" query:"bulk_purchase_order_id" validate:"required"`
	InquiryID           string `json:"inquiry_id"  param:"inquiry_id" query:"inquiry_id"`
	PurchaseOrderID     string `json:"purchase_order_id"  param:"purchase_order_id" query:"purchase_order_id"`
	UserID              string `json:"user_id"  param:"user_id" query:"user_id"`
	TeamID              string `json:"team_id" param:"team_id" query:"team_id"`

	IncludeShippingAddress     bool `json:"-"`
	IncludeUser                bool `json:"-"`
	IncludePaymentTransactions bool `json:"-"`
	IncludeAssignee            bool `json:"-"`
	IncludeInvoice             bool `json:"-"`
	IncludeItems               bool `json:"-"`
	IncludeSellerQuotation     bool `json:"-"`
	IsQueryAll                 bool `json:"-"`
}

func (r *BulkPurchaseOrderRepo) GetBulkPurchaseOrder(params GetBulkPurchaseOrderParams) (*models.BulkPurchaseOrder, error) {
	var userID = params.GetUserID()
	if params.TeamID != "" && !params.GetRole().IsAdmin() {
		if err := r.db.Select("ID").First(&models.BrandTeam{}, "team_id = ? AND user_id = ?", params.TeamID, userID).Error; err != nil {
			return nil, errs.ErrBulkPoNotFound
		}
		userID = params.TeamID
	}

	var bulkPurchaseOrder models.BulkPurchaseOrder
	var err = query.New(r.db, queryfunc.NewBulkPurchaseOrderBuilder(queryfunc.BulkPurchaseOrderBuilderOptions{
		IncludeShippingAddress:     params.IncludeShippingAddress,
		IncludeUser:                params.IncludeUser,
		IncludePaymentTransactions: params.IncludePaymentTransactions,
		IncludeAssignee:            params.IncludeAssignee,
		IncludeInvoice:             params.IncludeInvoice,
		IncludeItems:               params.IncludeItems,
		IncludeSellerQuotation:     params.IncludeSellerQuotation,
		QueryBuilderOptions: queryfunc.QueryBuilderOptions{
			Role: params.GetRole(),
		},
	})).
		WhereFunc(func(builder *query.Builder) {
			if params.BulkPurchaseOrderID != "" {
				if strings.HasPrefix(params.BulkPurchaseOrderID, "BPO-") {
					builder.Where("bpo.reference_id = ?", params.BulkPurchaseOrderID)
				} else {
					builder.Where("bpo.id = ?", params.BulkPurchaseOrderID)
				}
			}

			if params.GetRole().IsAdmin() {
				if params.UserID != "" && params.UserID != params.GetUserID() {
					builder.Where("bpo.user_id = ?", params.UserID)
				}
			} else if params.GetRole().IsSeller() {

			} else {
				builder.Where("bpo.user_id = ? AND bpo.deleted_at IS NULL", userID)
			}

			if params.InquiryID != "" {
				builder.Where("bpo.inquiry_id = ?", params.InquiryID)
			}

			if params.PurchaseOrderID != "" {
				builder.Where("bpo.purchase_order_id = ?", params.PurchaseOrderID)
			}

			if params.GetRole().IsSeller() {
				builder.Where("(bpo.seller_id = @user_id OR bposq.id IS NOT NULL) AND bpo.deleted_at IS NULL", sql.Named("user_id", userID))
			}
		}).
		Limit(1).
		FirstFunc(&bulkPurchaseOrder)

	return &bulkPurchaseOrder, err
}

func (r *BulkPurchaseOrderRepo) UpdateBulkPurchaseOrder(form models.BulkPurchaseOrderUpdateForm) (*models.BulkPurchaseOrder, error) {
	order, err := r.GetBulkPurchaseOrder(GetBulkPurchaseOrderParams{
		JwtClaimsInfo:       form.JwtClaimsInfo,
		BulkPurchaseOrderID: form.BulkPurchaseOrderID,
		UserID:              form.UserID,
	})
	if err != nil {
		return nil, err
	}

	var updates models.BulkPurchaseOrder
	err = copier.Copy(&updates, &form)
	if err != nil {
		return nil, err
	}

	var items []*models.BulkPurchaseOrderItem
	err = copier.Copy(&items, &form.Items)
	if err != nil {
		return nil, err
	}

	err = r.db.Transaction(func(tx *gorm.DB) error {
		if form.ShippingAddress != nil {
			if err := updates.ShippingAddress.CreateOrUpdate(r.db); err == nil {
				updates.ShippingAddressID = updates.ShippingAddress.ID
			}
			if err != nil {
				return err
			}
		}

		if order.TrackingStatus == enums.BulkPoTrackingStatusNew {
			updates.TrackingStatus = enums.BulkPoTrackingStatusWaitingForSubmitOrder
		}

		err = tx.Model(&models.BulkPurchaseOrder{}).Where("id = ?", form.BulkPurchaseOrderID).Updates(&updates).Error
		if err != nil {
			return err
		}

		for _, item := range items {
			item.PurchaseOrderID = order.ID
		}

		var err = tx.Unscoped().Delete(&models.BulkPurchaseOrderItem{}, "purchase_order_id = ?", order.ID).Error
		if err != nil {
			return err
		}

		if len(items) > 0 {
			if err := tx.Create(&items).Error; err != nil {
				return err
			}
		} else {
			order.Items = nil
		}

		// handle order items
		if err := tx.Unscoped().Delete(&models.OrderCartItem{}, "bulk_purchase_order_id = ?", order.ID).Error; err != nil {
			return err
		}
		var orderItems []*models.OrderCartItem
		for _, bpoItem := range items {
			orderItems = append(orderItems, &models.OrderCartItem{
				Model: models.Model{
					ID:        bpoItem.ID,
					CreatedAt: bpoItem.CreatedAt,
					UpdatedAt: bpoItem.UpdatedAt,
					DeletedAt: &bpoItem.DeletedAt,
				},
				BulkPurchaseOrderID: order.ID,
				Size:                bpoItem.Size,
				ColorName:           bpoItem.ColorName,
				Qty:                 bpoItem.Qty,
				UnitPrice:           *bpoItem.UnitPrice,
				TotalPrice:          *bpoItem.TotalPrice,
			})
		}
		if len(orderItems) > 0 {
			if err := tx.Create(&orderItems).Error; err != nil {
				return err
			}
			order.OrderCartItems = orderItems
		} else {
			order.OrderCartItems = nil
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	err = copier.CopyWithOption(order, &updates, copier.Option{IgnoreEmpty: true, DeepCopy: true})

	return order, err
}

func (r *BulkPurchaseOrderRepo) SubmitBulkPurchaseOrder(form models.BulkPurchaseOrderUpdateForm) (*models.BulkPurchaseOrder, error) {
	order, err := r.GetBulkPurchaseOrder(GetBulkPurchaseOrderParams{
		JwtClaimsInfo:       form.JwtClaimsInfo,
		BulkPurchaseOrderID: form.BulkPurchaseOrderID,
		UserID:              form.UserID,
		IncludeUser:         true,
	})
	if err != nil {
		return nil, err
	}

	var updates models.BulkPurchaseOrder

	updates.SubmittedAt = values.Int64(time.Now().Unix())
	updates.TrackingStatus = enums.BulkPoTrackingStatusWaitingForQuotation
	err = r.db.Transaction(func(tx *gorm.DB) error {
		err = NewBulkPurchaseOrderTrackingRepo(r.db).CreateBulkPurchaseOrderTrackingTx(tx, models.BulkPurchaseOrderTrackingCreateForm{
			PurchaseOrderID: order.ID,
			ActionType:      enums.BulkPoTrackingActionSubmitOrder,
			UserID:          order.UserID,
			CreatedByUserID: form.JwtClaimsInfo.GetUserID(),
			Description: func() string {
				if order.User != nil && order.User.Name != "" {
					return fmt.Sprintf("%s submitted order and waiting for quotation", order.User.Name)
				}
				return "User submitted order and waiting for quotation"
			}(),
		})
		if err != nil {
			return err
		}
		return tx.Model(&models.BulkPurchaseOrder{}).Where("id = ?", form.BulkPurchaseOrderID).Updates(&updates).Error
	})
	if err != nil {
		return nil, err
	}

	err = copier.CopyWithOption(order, &updates, copier.Option{IgnoreEmpty: true, DeepCopy: true})

	return order, err
}

func (r *BulkPurchaseOrderRepo) AdminSubmitBulkPurchaseOrder(form models.BulkPurchaseOrderUpdateForm) (*models.BulkPurchaseOrder, error) {
	order, err := r.GetBulkPurchaseOrder(GetBulkPurchaseOrderParams{
		JwtClaimsInfo:       form.JwtClaimsInfo,
		BulkPurchaseOrderID: form.BulkPurchaseOrderID,
		UserID:              form.UserID,
		IncludeUser:         true,
	})
	if err != nil {
		return nil, err
	}

	var updates models.BulkPurchaseOrder

	updates.SubmittedAt = values.Int64(time.Now().Unix())
	updates.TrackingStatus = enums.BulkPoTrackingStatusWaitingForQuotation
	err = r.db.Transaction(func(tx *gorm.DB) error {
		if admin, err := NewUserRepo(r.db).GetShortUserInfo(form.GetUserID()); err == nil {
			err = NewBulkPurchaseOrderTrackingRepo(r.db).CreateBulkPurchaseOrderTrackingTx(tx, models.BulkPurchaseOrderTrackingCreateForm{
				PurchaseOrderID: order.ID,
				ActionType:      enums.BulkPoTrackingActionSubmitOrder,
				UserID:          form.JwtClaimsInfo.GetUserID(),
				CreatedByUserID: form.JwtClaimsInfo.GetUserID(),
				Description:     fmt.Sprintf("%s submitted order and waiting for quotation", admin.Name),
			})
			if err != nil {
				return err
			}
		}

		return tx.Model(&models.BulkPurchaseOrder{}).Where("id = ?", form.BulkPurchaseOrderID).Updates(&updates).Error
	})
	if err != nil {
		return nil, err
	}

	err = copier.CopyWithOption(order, &updates, copier.Option{IgnoreEmpty: true, DeepCopy: true})

	return order, err
}

type BulkPurchaseOrderCheckoutParams struct {
	models.JwtClaimsInfo
	BulkPurchaseOrderID string `param:"bulk_purchase_order_id" validate:"required"`

	PaymentType     enums.PaymentType      `json:"payment_type" validate:"oneof=bank_transfer card"`
	PaymentMethodID string                 `json:"payment_method_id" validate:"required_if=PaymentType card"`
	Milestone       enums.PaymentMilestone `json:"milestone" param:"milestone" query:"milestone"`
	Note            string                 `json:"note"`

	TransactionRefID      string             `json:"transaction_ref_id" validate:"required_if=PaymentType bank_transfer"`
	TransactionAttachment *models.Attachment `json:"transaction_attachment" validate:"required_if=PaymentType bank_transfer"`
}

func (r *BulkPurchaseOrderRepo) BulkPurchaseOrderCheckout(params BulkPurchaseOrderCheckoutParams) (*models.BulkPurchaseOrder, error) {
	order, err := r.BulkPurchaseOrderPreviewCheckout(BulkPurchaseOrderPreviewCheckoutParams{
		JwtClaimsInfo:       params.JwtClaimsInfo,
		BulkPurchaseOrderID: params.BulkPurchaseOrderID,
		PaymentType:         params.PaymentType,
		Milestone:           params.Milestone,
	})
	if err != nil {
		return nil, err
	}

	if params.Milestone == enums.PaymentMilestoneFinalPayment {
		if order.TrackingStatus != enums.BulkPoTrackingStatusFinalPayment {
			return nil, errs.ErrBulkPoInvalidToCheckout
		}
	} else {
		if order.TrackingStatus != enums.BulkPoTrackingStatusFirstPayment {
			return nil, errs.ErrBulkPoInvalidToCheckout
		}
	}

	if params.PaymentType == enums.PaymentTypeBankTransfer {

		// create transaction
		var transaction = models.PaymentTransaction{
			BulkPurchaseOrderID: order.ID,
			Currency:            order.Currency,
			PaidAmount:          order.FirstPaymentTotal,
			PaymentType:         params.PaymentType,
			Milestone:           enums.PaymentMilestoneFirstPayment,
			UserID:              order.UserID,
			TransactionRefID:    params.TransactionRefID,
			Status:              enums.PaymentStatusWaitingConfirm,
			Note:                params.Note,
			BalanceAmount:       order.GetBalanceAmountAfterFirstPayment().ToPtr(),
			PaymentPercentage:   order.FirstPaymentPercentage,
			TotalAmount:         order.TotalPrice,
			Metadata: &models.PaymentTransactionMetadata{
				InquiryID:                    order.Inquiry.ID,
				InquiryReferenceID:           order.Inquiry.ReferenceID,
				BulkPurchaseOrderReferenceID: order.ReferenceID,
				BulkPurchaseOrderID:          order.ID,
			},
		}
		if params.TransactionAttachment != nil {
			transaction.Attachments = &models.Attachments{params.TransactionAttachment}
		}
		if params.Milestone == enums.PaymentMilestoneFinalPayment {
			transaction.PaidAmount = order.GetBalanceAmountAfterFirstPayment().ToPtr()
			transaction.BalanceAmount = price.NewFromFloat(0).ToPtr()
			transaction.Milestone = enums.PaymentMilestoneFinalPayment
		}

		err = r.db.Create(&transaction).Error
		if err != nil {
			return nil, eris.Wrap(err, err.Error())
		}

		var updates models.BulkPurchaseOrder
		updates.TaxPercentage = order.Inquiry.TaxPercentage
		if params.Milestone == enums.PaymentMilestoneFinalPayment {
			updates.TrackingStatus = enums.BulkPoTrackingStatusFinalPaymentConfirm
			updates.FinalPaymentTransferedAt = values.Int64(time.Now().Unix())
			updates.FinalPaymentTransactionRefID = params.TransactionRefID
			updates.FinalPaymentTransactionAttachment = params.TransactionAttachment
			updates.FinalPaymentTransactionReferenceID = transaction.ReferenceID
		} else {
			updates.TrackingStatus = enums.BulkPoTrackingStatusFirstPaymentConfirm
			updates.FirstPaymentTransferedAt = values.Int64(time.Now().Unix())
			updates.FirstPaymentTransactionRefID = params.TransactionRefID
			updates.FirstPaymentTransactionAttachment = params.TransactionAttachment
			updates.FirstPaymentTransactionReferenceID = transaction.ReferenceID
		}

		err = r.db.Transaction(func(tx *gorm.DB) error {
			err = r.db.Model(&models.BulkPurchaseOrder{}).Where("id = ?", order.ID).Updates(&updates).Error
			if err != nil {
				return eris.Wrap(err, err.Error())
			}

			var tracking = models.BulkPurchaseOrderTrackingCreateForm{
				ActionType:      enums.BulkPoTrackingActionMakeFirstPayment,
				UserID:          order.UserID,
				CreatedByUserID: params.JwtClaimsInfo.GetUserID(),
				Description:     "Processed payment",
				PurchaseOrderID: order.ID,
				Metadata: &models.PoTrackingMetadata{
					After: map[string]interface{}{
						"payment_transaction_id": transaction.ID,
					},
				},
			}
			if params.Milestone == enums.PaymentMilestoneFinalPayment {
				tracking.ActionType = enums.BulkPoTrackingActionMakeFinalPayment
			}
			return NewBulkPurchaseOrderTrackingRepo(r.db).CreateBulkPurchaseOrderTrackingTx(tx, tracking)
		})
		if err != nil {
			return nil, eris.Wrap(err, err.Error())
		}

		return order, err
	}

	var user models.User
	err = r.db.Select("ID", "StripeCustomerID").First(&user, "id = ?", params.GetUserID()).Error
	if err != nil {
		return nil, err
	}

	stripeConfig, err := stripehelper.GetCurrencyConfig(order.Currency)
	if err != nil {
		return nil, err
	}

	var stripeParams = stripehelper.CreatePaymentIntentParams{
		Amount:                  order.FirstPaymentTotal.MultipleInt(stripeConfig.SmallestUnitFactor).ToInt64(),
		Currency:                order.Currency,
		PaymentMethodID:         params.PaymentMethodID,
		CustomerID:              user.StripeCustomerID,
		IsCaptureMethodManually: false,
		Description:             fmt.Sprintf("Charges for first payment %s", order.ReferenceID),
		PaymentMethodTypes:      []string{"card"},
		Metadata: map[string]string{
			"milestone":                        string(params.Milestone),
			"inquiry_id":                       order.InquiryID,
			"bulk_purchase_order_id":           order.ID,
			"bulk_purchase_order_reference_id": order.ReferenceID,
			"action_source":                    string(stripehelper.ActionSourceBulkPOFirstPayment),
		},
	}

	if params.Milestone == enums.PaymentMilestoneFinalPayment {
		stripeParams.Amount = order.FinalPaymentTotal.MultipleInt(stripeConfig.SmallestUnitFactor).ToInt64()
		stripeParams.Description = fmt.Sprintf("Charges for final payment %s", order.ReferenceID)
		stripeParams.Metadata["action_source"] = string(stripehelper.ActionSourceBulkPOFinalPayment)
	}

	if order.ShippingAddress != nil {
		stripeParams.Shipping = &stripe.ShippingDetailsParams{
			Name: &order.ShippingAddress.Name,
		}
		if order.ShippingAddress.Coordinate != nil {
			stripeParams.Shipping.Address = &stripe.AddressParams{
				State:      &order.ShippingAddress.Coordinate.Level1,
				City:       &order.ShippingAddress.Coordinate.Level2,
				PostalCode: &order.ShippingAddress.Coordinate.PostalCode,
				Country:    stripe.String(order.ShippingAddress.Coordinate.CountryCode.String()),
				Line1:      &order.ShippingAddress.Coordinate.FormattedAddress,
			}
		}

	}

	pi, err := stripehelper.GetInstance().CreatePaymentIntent(stripeParams)
	if err != nil {
		return nil, err
	}

	if pi.Status != stripe.PaymentIntentStatusSucceeded {
		if pi.NextAction != nil {
			intent, err := stripehelper.GetInstance().ConfirmPaymentIntent(stripehelper.ConfirmPaymentIntentParams{
				PaymentIntentID: pi.ID,
				ReturnURL:       fmt.Sprintf("%s/api/v1/callback/stripe/payment_intents/inquiries/%s/bulk_purchase_orders/%s/confirm", r.db.Configuration.ServerBaseURL, order.InquiryID, order.ID),
			})
			if err != nil {
				return nil, err
			}

			if intent.Status == stripe.PaymentIntentStatusSucceeded {
				goto PaymentSuccess
			}

			order.PaymentIntentNextAction = intent.NextAction
			order.PaymentIntentClientSecret = intent.ClientSecret
			return order, nil
		} else {
			return nil, eris.Errorf("Payment error with status %s", pi.Status)
		}
	}

PaymentSuccess:

	// create transaction
	var transaction = models.PaymentTransaction{
		BulkPurchaseOrderID: order.ID,
		Currency:            order.Currency,
		PaidAmount:          order.FirstPaymentTotal,
		PaymentType:         params.PaymentType,
		Milestone:           enums.PaymentMilestoneFirstPayment,
		UserID:              order.UserID,
		TransactionRefID:    params.TransactionRefID,
		Status:              enums.PaymentStatusPaid,
		Note:                params.Note,
		BalanceAmount:       order.FinalPaymentTotal,
		PaymentPercentage:   order.FirstPaymentPercentage,
		TotalAmount:         order.TotalPrice,
		PaymentIntentID:     pi.ID,
		MarkAsPaidAt:        values.Int64(time.Now().Unix()),
		Metadata: &models.PaymentTransactionMetadata{
			InquiryID:                    order.Inquiry.ID,
			InquiryReferenceID:           order.Inquiry.ReferenceID,
			BulkPurchaseOrderReferenceID: order.ReferenceID,
			BulkPurchaseOrderID:          order.ID,
		},
	}
	if params.TransactionAttachment != nil {
		transaction.Attachments = &models.Attachments{params.TransactionAttachment}
	}
	if params.Milestone == enums.PaymentMilestoneFinalPayment {
		transaction.PaidAmount = order.FinalPaymentTotal
		transaction.BalanceAmount = price.NewFromFloat(0).ToPtr()
		transaction.Milestone = enums.PaymentMilestoneFinalPayment
		transaction.PaymentPercentage = values.Float64(100 - values.Float64Value(order.FirstPaymentPercentage))
	}

	r.db.Transaction(func(tx *gorm.DB) error {
		err = tx.Create(&transaction).Error
		if err != nil {
			return eris.Wrap(err, err.Error())
		}

		var updates = models.BulkPurchaseOrder{
			TrackingStatus:                     enums.BulkPoTrackingStatusFirstPaymentConfirmed,
			FirstPaymentTransferedAt:           values.Int64(time.Now().Unix()),
			FirstPaymentMarkAsPaidAt:           values.Int64(time.Now().Unix()),
			FirstPaymentTransactionReferenceID: transaction.ReferenceID,
		}

		if len(order.AdminQuotations) > 0 {
			var bulkQuotation = order.GetBulkQuotation()
			if bulkQuotation != nil {
				updates.LeadTime = int(values.Int64Value(bulkQuotation.LeadTime))
				updates.StartDate = updates.FirstPaymentMarkAsPaidAt
				updates.CompletionDate = values.Int64(time.Unix(*updates.StartDate, 0).AddDate(0, 0, updates.LeadTime).Unix())
			}
		}

		updates.TaxPercentage = order.Inquiry.TaxPercentage

		if params.Milestone == enums.PaymentMilestoneFinalPayment {
			updates.TrackingStatus = enums.BulkPoTrackingStatusFinalPaymentConfirmed
			updates.FinalPaymentTransferedAt = values.Int64(time.Now().Unix())
			updates.FinalPaymentMarkAsPaidAt = values.Int64(time.Now().Unix())
			updates.FinalPaymentTransactionReferenceID = transaction.ReferenceID
		}

		err = tx.Model(&models.BulkPurchaseOrder{}).Where("id = ?", order.ID).Updates(&updates).Error
		return err
	})
	if err != nil {
		return nil, err
	}

	return order, err

}

type BulkPurchaseOrderPreviewCheckoutParams struct {
	models.JwtClaimsInfo

	BulkPurchaseOrderID string                 `json:"bulk_purchase_order_id" param:"bulk_purchase_order_id" validate:"required"`
	PaymentType         enums.PaymentType      `json:"payment_type" validate:"oneof=bank_transfer card"`
	Milestone           enums.PaymentMilestone `json:"milestone" param:"milestone" query:"milestone"`

	Inquiry            *models.Inquiry           `json:"-"`
	BulkPurchaseOrder  *models.BulkPurchaseOrder `json:"-"`
	UpdatePricing      bool                      `json:"-"`
	IsSkipFirstPayment bool                      `json:"-"`
}

// deprecated
func (r *BulkPurchaseOrderRepo) BulkPurchaseOrderPreviewCheckout(params BulkPurchaseOrderPreviewCheckoutParams) (*models.BulkPurchaseOrder, error) {
	cancel, err := r.db.Locker.AcquireLock(fmt.Sprintf("bulk_purchase_order_%s", params.BulkPurchaseOrderID), time.Second*30)
	if err != nil {
		return nil, err
	}
	defer cancel()

	var bulkPurchaseOrder = params.BulkPurchaseOrder
	if bulkPurchaseOrder == nil || len(bulkPurchaseOrder.Items) == 0 {
		bulkPurchaseOrder, err = r.GetBulkPurchaseOrder(GetBulkPurchaseOrderParams{
			JwtClaimsInfo:       params.JwtClaimsInfo,
			BulkPurchaseOrderID: params.BulkPurchaseOrderID,
			IncludeItems:        true,
		})
		if err != nil {
			return nil, err
		}
	}

	if bulkPurchaseOrder.ShippingAddressID != "" && bulkPurchaseOrder.ShippingAddress == nil {
		bulkPurchaseOrder.ShippingAddress, _ = NewAddressRepo(r.db).GetAddress(GetAddressParams{
			AddressID: bulkPurchaseOrder.ShippingAddressID,
		})
	}

	if err != nil {
		return nil, err
	}

	var subTotalPrice = price.NewFromFloat(0)
	var orderCartItemsToUpdate []*models.OrderCartItem
	for _, item := range bulkPurchaseOrder.OrderCartItems {
		if item.UnitPrice.LessThanOrEqual(0) || item.TotalPrice.LessThanOrEqual(0) {
			item.UnitPrice = *bulkPurchaseOrder.QuotedPrice
			item.TotalPrice = item.UnitPrice.MultipleInt(int64(item.Qty))
			orderCartItemsToUpdate = append(orderCartItemsToUpdate, item)
		}
		subTotalPrice = subTotalPrice.Add(item.TotalPrice)
	}

	for index, additionalItem := range bulkPurchaseOrder.AdditionalItems {
		additionalItem.BulkPurchaseOrderID = bulkPurchaseOrder.ID
		additionalItem.TotalPrice = additionalItem.UnitPrice.MultipleInt(int64(additionalItem.Qty)).ToPtr()
		subTotalPrice = subTotalPrice.AddPtr(additionalItem.TotalPrice)

		bulkPurchaseOrder.AdditionalItems[index] = additionalItem
	}

	bulkPurchaseOrder.SubTotal = subTotalPrice.ToPtr()
	if values.Float64Value(bulkPurchaseOrder.TaxPercentage) == 0 {
		if bulkPurchaseOrder.Inquiry != nil {
			bulkPurchaseOrder.TaxPercentage = bulkPurchaseOrder.Inquiry.TaxPercentage
		}

		if bulkPurchaseOrder.PurchaseOrder != nil {
			bulkPurchaseOrder.TaxPercentage = bulkPurchaseOrder.PurchaseOrder.TaxPercentage
		}
	}

	if bulkPurchaseOrder.CommercialInvoice != nil {
		bulkPurchaseOrder.TaxPercentage = bulkPurchaseOrder.CommercialInvoice.TaxPercentage
	}

	if params.Milestone == enums.PaymentMilestoneFirstPayment {
		bulkPurchaseOrder.FirstPaymentType = params.PaymentType
	} else {
		bulkPurchaseOrder.FinalPaymentType = params.PaymentType
	}

	if params.IsSkipFirstPayment {
		bulkPurchaseOrder.FirstPaymentPercentage = values.Float64(0)
		bulkPurchaseOrder.TrackingStatus = enums.BulkPoTrackingStatusFirstPaymentConfirmed
		bulkPurchaseOrder.FirstPaymentTransferedAt = values.Int64(time.Now().Unix())
		bulkPurchaseOrder.FirstPaymentMarkAsPaidAt = values.Int64(time.Now().Unix())
	}
	err = bulkPurchaseOrder.UpdatePrices()
	if err != nil {
		return nil, err
	}

	if err := r.db.Transaction(func(tx *gorm.DB) error {
		if !params.GetRole().IsAdmin() || params.UpdatePricing {
			if err := tx.Model(&models.BulkPurchaseOrder{}).Where("id = ?", bulkPurchaseOrder.ID).Updates(bulkPurchaseOrder).Error; err != nil {
				return err
			}
		}
		if len(orderCartItemsToUpdate) > 0 {
			if err := tx.Clauses(clause.OnConflict{
				Columns:   []clause.Column{{Name: "id"}},
				DoUpdates: clause.AssignmentColumns([]string{"unit_price", "total_price"}),
			}).Create(&orderCartItemsToUpdate).Error; err != nil {
				return err
			}
		}
		return nil
	}); err != nil {
		return nil, err
	}

	return bulkPurchaseOrder, err
}

type BulkPurchaseOrdersPreviewCheckoutParams struct {
	models.JwtClaimsInfo
	PaymentType        enums.PaymentType `json:"payment_type" validate:"oneof=bank_transfer card"`
	UpdatePricing      bool              `json:"-"`
	IsSkipFirstPayment bool              `json:"-"`

	Params []BulkPurchaseOrderPreviewCheckoutParams `json:"checkout_params"`
}

// deprecated
func (r *BulkPurchaseOrderRepo) BulkPurchaseOrdersPreviewCheckout(params BulkPurchaseOrdersPreviewCheckoutParams) (bulks []*models.BulkPurchaseOrder, err error) {
	for _, param := range params.Params {
		param.JwtClaimsInfo = params.JwtClaimsInfo
		param.PaymentType = params.PaymentType
		param.UpdatePricing = params.UpdatePricing
		param.IsSkipFirstPayment = params.IsSkipFirstPayment

		var bulk *models.BulkPurchaseOrder
		bulk, err = r.BulkPurchaseOrderPreviewCheckout(param)
		if err != nil {
			return
		}
		bulks = append(bulks, bulk)
	}
	return
}

func (r *BulkPurchaseOrderRepo) BulkPurchaseOrderMarkAsPaid(params models.BulkPurchaseOrderMarkAsPaidParams) (*models.BulkPurchaseOrder, error) {
	cancel, err := r.db.Locker.AcquireLock(fmt.Sprintf("bulk_purchase_order_payment_%s", params.BulkPurchaseOrderID), time.Second*20)
	if err != nil {
		return nil, err
	}
	defer cancel()

	var bulkPurchaseOrder models.BulkPurchaseOrder
	err = r.db.Select("ID", "TrackingStatus", "UserID", "PurchaseOrderID", "AdminQuotations").First(&bulkPurchaseOrder, "id = ?", params.BulkPurchaseOrderID).Error
	if err != nil {
		return nil, err
	}

	var updates = models.BulkPurchaseOrder{
		FirstPaymentMarkAsPaidAt: values.Int64(time.Now().Unix()),
		TrackingStatus:           enums.BulkPoTrackingStatusFirstPaymentConfirmed,
	}

	var paymentTransactionUpdates = models.PaymentTransaction{
		MarkAsPaidAt: values.Int64(time.Now().Unix()),
		Status:       enums.PaymentStatusPaid,
	}

	if len(bulkPurchaseOrder.AdminQuotations) > 0 && (params.Milestone == enums.PaymentMilestoneFirstPayment || params.Milestone == enums.PaymentMilestoneDeposit) {
		var purchaseOrder models.PurchaseOrder
		r.db.Select("Quotations").First(&purchaseOrder, "id = ?", bulkPurchaseOrder.PurchaseOrderID)

		sampleQuotation, _ := lo.Find(purchaseOrder.Quotations, func(item *models.InquiryQuotationItem) bool {
			return item.Type == enums.InquiryTypeSample
		})
		if sampleQuotation != nil {
			updates.LeadTime = int(values.Int64Value(sampleQuotation.LeadTime))
			updates.StartDate = paymentTransactionUpdates.MarkAsPaidAt
			updates.CompletionDate = values.Int64(time.Unix(*updates.StartDate, 0).AddDate(0, 0, updates.LeadTime).Unix())
		}
	}

	if params.Milestone == enums.PaymentMilestoneFinalPayment {
		updates.FinalPaymentMarkAsPaidAt = values.Int64(time.Now().Unix())
		updates.TrackingStatus = enums.BulkPoTrackingStatusFinalPaymentConfirmed

		var err = r.db.Transaction(func(tx *gorm.DB) error {
			var err = tx.Model(&models.BulkPurchaseOrder{}).
				Where("id = ? AND final_payment_mark_as_paid_at IS NULL", params.BulkPurchaseOrderID).
				Updates(&updates).Error
			if err != nil {
				return err
			}

			return tx.Model(&models.PaymentTransaction{}).
				Where("bulk_purchase_order_id = ? AND payment_type = ? AND mark_as_paid_at IS NULL AND milestone = ?", params.BulkPurchaseOrderID, enums.PaymentTypeBankTransfer, enums.PaymentMilestoneFinalPayment).
				Updates(&paymentTransactionUpdates).Error
		})

		return &bulkPurchaseOrder, err
	}

	err = r.db.Transaction(func(tx *gorm.DB) error {
		err = tx.Model(&models.BulkPurchaseOrder{}).
			Where("id = ? AND first_payment_mark_as_paid_at IS NULL", params.BulkPurchaseOrderID).
			Updates(&updates).Error
		if err != nil {
			return err
		}

		err = tx.Model(&models.PaymentTransaction{}).
			Where("bulk_purchase_order_id = ? AND payment_type = ? AND mark_as_paid_at IS NULL AND milestone = ?", params.BulkPurchaseOrderID, enums.PaymentTypeBankTransfer, enums.PaymentMilestoneFirstPayment).
			Updates(&paymentTransactionUpdates).Error

		return err
	})

	if err != nil {
		return &bulkPurchaseOrder, err
	}

	return &bulkPurchaseOrder, err
}

func (r *BulkPurchaseOrderRepo) BulkPurchaseOrderMarkAsUnpaid(params models.BulkPurchaseOrderMarkAsUnpaidParams) (*models.BulkPurchaseOrder, error) {
	cancel, err := r.db.Locker.AcquireLock(fmt.Sprintf("bulk_purchase_order_payment_%s", params.BulkPurchaseOrderID), time.Second*20)
	if err != nil {
		return nil, err
	}
	defer cancel()

	var bulkPurchaseOrder models.BulkPurchaseOrder
	err = r.db.Select("ID", "TrackingStatus", "UserID").First(&bulkPurchaseOrder, "id = ?", params.BulkPurchaseOrderID).Error
	if err != nil {
		return nil, err
	}

	var updates = models.BulkPurchaseOrder{
		FirstPaymentMarkAsUnpaidAt: values.Int64(time.Now().Unix()),
	}

	var paymentTransactionUpdates = models.PaymentTransaction{
		MarkAsUnpaidAt: values.Int64(time.Now().Unix()),
		Status:         enums.PaymentStatusUnpaid,
	}

	if params.Milestone == enums.PaymentMilestoneFinalPayment {
		updates.FinalPaymentMarkAsUnpaidAt = values.Int64(time.Now().Unix())

		var err = r.db.Transaction(func(tx *gorm.DB) error {
			var err = tx.Model(&models.BulkPurchaseOrder{}).
				Where("id = ? AND final_payment_mark_as_paid_at IS NULL", params.BulkPurchaseOrderID).
				Updates(&updates).Error
			if err != nil {
				return err
			}

			return tx.Model(&models.PaymentTransaction{}).
				Where("bulk_purchase_order_id = ? AND payment_type = ? AND mark_as_paid_at IS NULL AND milestone = ?", params.BulkPurchaseOrderID, enums.PaymentTypeBankTransfer, enums.PaymentMilestoneFinalPayment).
				Updates(&paymentTransactionUpdates).Error
		})

		return &bulkPurchaseOrder, err
	}

	err = r.db.Model(&models.BulkPurchaseOrder{}).
		Where("id = ? AND first_payment_mark_as_paid_at IS NULL", params.BulkPurchaseOrderID).
		Updates(&updates).Error

	if err != nil {
		return &bulkPurchaseOrder, err
	}

	err = r.db.Model(&models.PaymentTransaction{}).
		Where("bulk_purchase_order_id = ? AND payment_type = ? AND mark_as_paid_at IS NULL AND milestone = ?", params.BulkPurchaseOrderID, enums.PaymentTypeBankTransfer, enums.PaymentMilestoneFirstPayment).
		Updates(&paymentTransactionUpdates).Error

	return &bulkPurchaseOrder, err
}

type UpdateBulkPurchaseOrderTrackingStatusParams struct {
	models.JwtClaimsInfo

	BulkPurchaseOrderID string                           `json:"bulk_purchase_order_id" param:"bulk_purchase_order_id" query:"bulk_purchase_order_id" validate:"required"`
	TrackingStatus      enums.PoTrackingStatus           `json:"tracking_status" param:"tracking_status" query:"tracking_status" validate:"required"`
	ApproveRejectMeta   *models.InquiryApproveRejectMeta `json:"approve_reject_meta" param:"approve_reject_meta" query:"approve_reject_meta"`

	UserID string `json:"-"`
}

func (r *BulkPurchaseOrderRepo) UpdatePurchaseOrderTrackingStatus(params UpdateBulkPurchaseOrderTrackingStatusParams) (*models.BulkPurchaseOrder, error) {
	order, err := r.GetBulkPurchaseOrder(GetBulkPurchaseOrderParams{
		JwtClaimsInfo:       params.JwtClaimsInfo,
		BulkPurchaseOrderID: params.BulkPurchaseOrderID,
		UserID:              params.UserID,
	})

	if err != nil {
		return nil, err
	}

	var updates models.BulkPurchaseOrder
	err = copier.Copy(&updates, &params)
	if err != nil {
		return nil, err
	}

	err = r.db.Model(&models.BulkPurchaseOrder{}).
		Where("id = ?", params.BulkPurchaseOrderID).
		Updates(&updates).Error

	if err != nil {
		return nil, eris.Wrap(err, err.Error())
	}

	if order.TrackingStatus != updates.TrackingStatus {
		err = NewInquiryAuditRepo(r.db).CreateInquiryAudit(models.InquiryAuditCreateForm{
			InquiryID:           order.InquiryID,
			ActionType:          enums.AuditActionTypeInquiryBulkPoCreated,
			UserID:              order.UserID,
			Description:         "Bulk PO status is updated",
			BulkPurchaseOrderID: order.ID,
			Metadata: &models.InquiryAuditMetadata{
				Before: map[string]interface{}{
					"tracking_status": order.TrackingStatus,
				},
				After: map[string]interface{}{
					"tracking_status": updates.TrackingStatus,
				},
			},
		})
	}

	if err != nil {
		return nil, eris.Wrap(err, err.Error())
	}

	order.TrackingStatus = updates.TrackingStatus
	return order, err
}

func (r *BulkPurchaseOrderRepo) SendQuotationToBuyer(form models.SendBulkPurchaseOrderQuotationParams) (*models.BulkPurchaseOrder, error) {
	order, err := r.GetBulkPurchaseOrder(GetBulkPurchaseOrderParams{
		BulkPurchaseOrderID: form.BulkPurchaseOrderID,
		JwtClaimsInfo:       form.JwtClaimsInfo,
		IncludeUser:         true,
	})
	if err != nil {
		return nil, err
	}

	var updates models.BulkPurchaseOrder
	err = copier.Copy(&updates, &form)
	if err != nil {
		return nil, err
	}
	order.FirstPaymentPercentage = &form.FirstPaymentPercentage
	updates.FirstPaymentPercentage = &form.FirstPaymentPercentage

	updates.TrackingStatus = enums.BulkPoTrackingStatusFirstPayment
	updates.QuotationAt = values.Int64(time.Now().Unix())

	err = r.db.Transaction(func(tx *gorm.DB) error {
		err = NewBulkPurchaseOrderTrackingRepo(r.db).CreateBulkPurchaseOrderTrackingTx(tx, models.BulkPurchaseOrderTrackingCreateForm{
			PurchaseOrderID: order.ID,
			ActionType:      enums.BulkPoTrackingActionSubmitQuotation,
			UserID:          order.UserID,
			CreatedByUserID: form.JwtClaimsInfo.GetUserID(),
			Description:     "Admin has sent quotation to buyer",
			Metadata: &models.PoTrackingMetadata{
				After: map[string]interface{}{
					"quotations": form.AdminQuotations,
				},
			},
		})
		if err != nil {
			return eris.Wrap(err, err.Error())
		}

		return tx.Model(&models.BulkPurchaseOrder{}).
			Where("id = ?", form.BulkPurchaseOrderID).
			Updates(&updates).Error

	})
	if err != nil {
		return nil, eris.Wrap(err, err.Error())
	}

	order.TrackingStatus = updates.TrackingStatus
	order.QuotationAt = updates.QuotationAt

	if form.FirstPaymentPercentage == 0 {
		return r.BulkPurchaseOrderPreviewCheckout(BulkPurchaseOrderPreviewCheckoutParams{
			JwtClaimsInfo:       form.JwtClaimsInfo,
			BulkPurchaseOrderID: form.BulkPurchaseOrderID,
			PaymentType:         enums.PaymentTypeBankTransfer,
			Milestone:           enums.PaymentMilestoneFirstPayment,
			UpdatePricing:       true,
			IsSkipFirstPayment:  true,
			BulkPurchaseOrder:   order,
		})
	}

	return order, nil
}

type AdminCreateQcReportParams struct {
	models.JwtClaimsInfo

	BulkPurchaseOrderID string                 `json:"bulk_purchase_order_id" param:"bulk_purchase_order_id" query:"bulk_purchase_order_id" validate:"required"`
	PoQcReports         []*models.PoReportMeta `json:"po_qc_reports" param:"po_qc_reports" query:"po_qc_reports" validate:"required"`
	ApproveQCAt         *int64                 `json:"approve_qc_at" param:"approve_qc_at" query:"approve_qc_at" `
}

func (r *BulkPurchaseOrderRepo) BulkPurchaseOrderCreateQcReport(params AdminCreateQcReportParams) (*models.BulkPurchaseOrder, error) {
	order, err := r.GetBulkPurchaseOrder(GetBulkPurchaseOrderParams{
		BulkPurchaseOrderID: params.BulkPurchaseOrderID,
		JwtClaimsInfo:       params.JwtClaimsInfo,
		IncludeUser:         true,
	})

	if err != nil {
		return nil, eris.Wrap(err, err.Error())
	}

	var updates models.BulkPurchaseOrder
	updates.TrackingStatus = enums.BulkPoTrackingStatusQc

	var reports models.PoReportMetas = params.PoQcReports
	updates.PoQcReports = &reports
	updates.ApproveQCAt = params.ApproveQCAt

	var trackings = lo.Map(params.PoQcReports, func(item *models.PoReportMeta, index int) *models.BulkPurchaseOrderTracking {
		return &models.BulkPurchaseOrderTracking{
			PurchaseOrderID: order.ID,
			ActionType:      enums.BulkPoTrackingActionCreateQcReport,
			UserID:          order.UserID,
			CreatedByUserID: params.JwtClaimsInfo.GetUserID(),
			ReportStatus:    item.Status,
			Attachments:     &item.Attachments,
			Description:     item.Description,
			Metadata: &models.PoTrackingMetadata{
				After: map[string]interface{}{
					"po_report": item,
				},
			},
		}
	})

	err = r.db.Transaction(func(tx *gorm.DB) error {
		err = tx.Create(&trackings).Error
		if err != nil {
			return err
		}
		return tx.Model(&models.BulkPurchaseOrder{}).Where("id = ?", params.BulkPurchaseOrderID).Updates(&updates).Error

	})
	if err != nil {
		return nil, eris.Wrap(err, err.Error())
	}

	order.TrackingStatus = updates.TrackingStatus
	order.PoQcReports = updates.PoQcReports

	return order, err
}

type BulkPurchaseOrderUpdateRawMaterialParams struct {
	models.JwtClaimsInfo

	BulkPurchaseOrderID  string                      `json:"bulk_purchase_order_id" param:"bulk_purchase_order_id" query:"bulk_purchase_order_id" validate:"required"`
	PoRawMaterials       []*models.PoRawMaterialMeta `json:"po_raw_materials" param:"po_raw_materials" query:"po_raw_materials"`
	ApproveRawMaterialAt *int64                      `json:"approve_raw_material_at" param:"approve_raw_material_at" query:"approve_raw_material_at"`
}

func (r *BulkPurchaseOrderRepo) BulkPurchaseOrderUpdateRawMaterial(params BulkPurchaseOrderUpdateRawMaterialParams) (*models.BulkPurchaseOrder, error) {
	order, err := r.GetBulkPurchaseOrder(GetBulkPurchaseOrderParams{
		BulkPurchaseOrderID: params.BulkPurchaseOrderID,
		JwtClaimsInfo:       params.JwtClaimsInfo,
		IncludeUser:         true,
	})
	if err != nil {
		return nil, err
	}

	var updates models.BulkPurchaseOrder
	err = copier.Copy(&updates, &params)
	if err != nil {
		return nil, err
	}

	_ = updates.GenerateRawMaterialRefID(updates.PoRawMaterials)

	updates.TrackingStatus = enums.BulkPoTrackingStatusRawMaterial
	err = r.db.Transaction(func(tx *gorm.DB) error {
		err = NewBulkPurchaseOrderTrackingRepo(r.db).CreateBulkPurchaseOrderTrackingTx(tx, models.BulkPurchaseOrderTrackingCreateForm{
			PurchaseOrderID: order.ID,
			ActionType:      enums.BulkPoTrackingActionUpdateMaterial,
			UserID:          order.UserID,
			CreatedByUserID: params.JwtClaimsInfo.GetUserID(),
			Metadata: &models.PoTrackingMetadata{
				Before: map[string]interface{}{
					"po_raw_materials": order.PoRawMaterials,
				},
				After: map[string]interface{}{
					"po_raw_materials": params.PoRawMaterials,
				},
			},
		})
		if err != nil {
			return err
		}
		return tx.Model(&models.BulkPurchaseOrder{}).
			Where("id = ?", params.BulkPurchaseOrderID).
			Updates(&updates).Error
	})
	if err != nil {
		return nil, eris.Wrap(err, err.Error())
	}

	order.TrackingStatus = updates.TrackingStatus

	return order, err
}

type BulkPurchaseOrderUpdateTrackingStatusParams struct {
	models.JwtClaimsInfo

	BulkPurchaseOrderID string                     `json:"bulk_purchase_order_id" param:"bulk_purchase_order_id" query:"bulk_purchase_order_id" validate:"required"`
	TrackingStatus      enums.BulkPoTrackingStatus `json:"tracking_status" param:"tracking_status" query:"tracking_status"`
	TrackingAction      enums.BulkPoTrackingAction `json:"tracking_action" param:"tracking_action" query:"tracking_action"`

	UserID string `json:"-"`
}

func (r *BulkPurchaseOrderRepo) BulkPurchaseOrderUpdateTrackingStatus(params BulkPurchaseOrderUpdateTrackingStatusParams) (*models.BulkPurchaseOrder, error) {
	order, err := r.GetBulkPurchaseOrder(GetBulkPurchaseOrderParams{
		BulkPurchaseOrderID: params.BulkPurchaseOrderID,
		JwtClaimsInfo:       params.JwtClaimsInfo,
		IncludeUser:         true,
	})

	if err != nil {
		return nil, err
	}

	var updates = models.BulkPurchaseOrder{
		TrackingStatus: params.TrackingStatus,
	}

	err = r.db.Transaction(func(tx *gorm.DB) error {
		if order.TrackingStatus != updates.TrackingStatus {
			err = NewBulkPurchaseOrderTrackingRepo(r.db).CreateBulkPurchaseOrderTrackingTx(tx, models.BulkPurchaseOrderTrackingCreateForm{
				PurchaseOrderID: params.BulkPurchaseOrderID,
				ActionType:      params.TrackingAction,
				UserID:          order.UserID,
				CreatedByUserID: params.JwtClaimsInfo.GetUserID(),
				Metadata: &models.PoTrackingMetadata{
					Before: map[string]interface{}{
						"tracking_status": order.TrackingStatus,
					},
					After: map[string]interface{}{
						"tracking_status": updates.TrackingStatus,
					},
				},
			})

			if err != nil {
				return eris.Wrap(err, err.Error())
			}

		}

		return tx.Model(&models.BulkPurchaseOrder{}).
			Where("id = ?", params.BulkPurchaseOrderID).
			Updates(&updates).Error
	})

	order.TrackingStatus = updates.TrackingStatus
	return order, err
}

type BulkPurchaseBuyerApproveRawMaterialParams struct {
	models.JwtClaimsInfo

	BulkPurchaseOrderID string `json:"bulk_purchase_order_id" param:"bulk_purchase_order_id" query:"bulk_purchase_order_id" validate:"required"`

	ItemIDs []string `json:"item_ids" param:"item_ids" query:"item_ids" validate:"required"`
}

func (r *BulkPurchaseOrderRepo) BulkPurchaseOrderBuyerApproveRawMaterial(params BulkPurchaseBuyerApproveRawMaterialParams) (*models.BulkPurchaseOrder, error) {
	order, err := r.GetBulkPurchaseOrder(GetBulkPurchaseOrderParams{
		BulkPurchaseOrderID: params.BulkPurchaseOrderID,
		JwtClaimsInfo:       params.JwtClaimsInfo,
	})

	if err != nil {
		return nil, err
	}

	if *order.PoRawMaterials == nil {
		return nil, errs.ErrBulkPoInvalidToApproveRawMaterial
	}

	var updates models.BulkPurchaseOrder
	updates.PoRawMaterials = order.PoRawMaterials
	var shouldTracking = false
	for _, existItem := range *updates.PoRawMaterials {
		// empty itemIds mean auto approve from system -> approve all item
		if (slices.Contains(params.ItemIDs, existItem.ReferenceID) && !values.BoolValue(existItem.BuyerApproved)) || len(params.ItemIDs) == 0 {
			existItem.BuyerApproved = values.Bool(true)
			existItem.Status = enums.PoRawMaterialStatusApproved
			shouldTracking = true
		}
	}

	err = r.db.Transaction(func(tx *gorm.DB) error {
		if shouldTracking {
			err = NewBulkPurchaseOrderTrackingRepo(r.db).CreateBulkPurchaseOrderTrackingTx(tx, models.BulkPurchaseOrderTrackingCreateForm{
				PurchaseOrderID: params.BulkPurchaseOrderID,
				ActionType:      enums.BulkPoTrackingActionBuyerApproveRawMaterial,
				UserID:          order.UserID,
				CreatedByUserID: params.JwtClaimsInfo.GetUserID(),
				Metadata: &models.PoTrackingMetadata{
					Before: map[string]interface{}{
						"po_raw_materials": order.PoRawMaterials,
					},
					After: map[string]interface{}{
						"po_raw_materials": updates.PoRawMaterials,
					},
				},
			})

			if err != nil {
				return eris.Wrap(err, err.Error())
			}

		}

		return tx.Model(&models.BulkPurchaseOrder{}).
			Where("id = ?", params.BulkPurchaseOrderID).
			Updates(&updates).Error
	})

	if err != nil {
		return nil, eris.Wrap(err, err.Error())
	}
	order.PoRawMaterials = updates.PoRawMaterials
	return order, err
}

type BulkOrderConfirmDeliveredParams struct {
	models.JwtClaimsInfo

	BulkPurchaseOrderID string `json:"bulk_purchase_order_id" param:"bulk_purchase_order_id" query:"bulk_purchase_order_id" validate:"required"`
}

func (r *BulkPurchaseOrderRepo) BuyerConfirmDelivered(params BulkOrderConfirmDeliveredParams) (*models.BulkPurchaseOrder, error) {
	order, err := r.GetBulkPurchaseOrder(GetBulkPurchaseOrderParams{
		BulkPurchaseOrderID: params.BulkPurchaseOrderID,
		JwtClaimsInfo:       params.JwtClaimsInfo,
	})
	if err != nil {
		return nil, err
	}

	validStatus := []enums.BulkPoTrackingStatus{enums.BulkPoTrackingStatusDelivering}
	if ok := slices.Contains(validStatus, order.TrackingStatus); !ok {
		return nil, errs.ErrPoInvalidToConfirmDelivered
	}

	var updates = models.BulkPurchaseOrder{
		TrackingStatus:      enums.BulkPoTrackingStatusDeliveryConfirmed,
		ReceiverConfirmedAt: values.Int64(time.Now().Unix()),
	}

	err = r.db.Transaction(func(tx *gorm.DB) error {
		err = NewBulkPurchaseOrderTrackingRepo(r.db).CreateBulkPurchaseOrderTrackingTx(tx, models.BulkPurchaseOrderTrackingCreateForm{
			PurchaseOrderID: params.BulkPurchaseOrderID,
			ActionType:      enums.BulkPoTrackingActionConfirmDelivered,
			UserID:          order.UserID,
			CreatedByUserID: params.JwtClaimsInfo.GetUserID(),
			Metadata: &models.PoTrackingMetadata{
				Before: map[string]interface{}{
					"tracking_status": order.TrackingStatus,
				},
				After: map[string]interface{}{
					"tracking_status": updates.TrackingStatus,
				},
			},
		})
		if err != nil {
			return err
		}
		return tx.Model(&models.BulkPurchaseOrder{}).
			Where("id = ?", params.BulkPurchaseOrderID).
			Updates(&updates).Error
	})
	if err != nil {
		return nil, eris.Wrap(err, err.Error())
	}

	return order, err
}

type BulkPurchaseOrderMarkDeliveringParams struct {
	models.JwtClaimsInfo

	BulkPurchaseOrderID string                 `json:"bulk_purchase_order_id" param:"bulk_purchase_order_id" query:"bulk_purchase_order_id" validate:"required"`
	LogisticInfo        *models.PoLogisticMeta `json:"logistic_info" param:"logistic_info" query:"logistic_info" form:"logistic_info" validate:"required"`
}

func (r *BulkPurchaseOrderRepo) MarkDelivering(params BulkPurchaseOrderMarkDeliveringParams) (*models.BulkPurchaseOrder, error) {
	order, err := r.GetBulkPurchaseOrder(GetBulkPurchaseOrderParams{
		BulkPurchaseOrderID: params.BulkPurchaseOrderID,
		JwtClaimsInfo:       params.JwtClaimsInfo,
		IncludeUser:         true,
	})

	validStatus := []enums.BulkPoTrackingStatus{enums.BulkPoTrackingStatusFinalPaymentConfirm, enums.BulkPoTrackingStatusFinalPaymentConfirmed}
	if ok := slices.Contains(validStatus, order.TrackingStatus); !ok {
		return nil, errs.ErrBulkPoInvalidToChangeTrackingStatus
	}

	if err != nil {
		return nil, err
	}

	var updates = models.BulkPurchaseOrder{
		DeliveryStartedAt: values.Int64(time.Now().Unix()),
		TrackingStatus:    enums.BulkPoTrackingStatusDelivering,
	}

	err = r.db.Transaction(func(tx *gorm.DB) error {
		if order.TrackingStatus != updates.TrackingStatus {
			err = NewBulkPurchaseOrderTrackingRepo(r.db).CreateBulkPurchaseOrderTrackingTx(tx, models.BulkPurchaseOrderTrackingCreateForm{
				PurchaseOrderID: params.BulkPurchaseOrderID,
				ActionType:      enums.BulkPoTrackingActionMarkDelivering,
				UserID:          order.UserID,
				CreatedByUserID: params.JwtClaimsInfo.GetUserID(),
				Metadata: &models.PoTrackingMetadata{
					Before: map[string]interface{}{
						"tracking_status": order.TrackingStatus,
					},
					After: map[string]interface{}{
						"tracking_status": updates.TrackingStatus,
						"logistic_info":   params.LogisticInfo,
					},
				},
			})

			if err != nil {
				return eris.Wrap(err, err.Error())
			}

			if err != nil {
				return eris.Wrap(err, err.Error())
			}
		}

		return tx.Model(&models.BulkPurchaseOrder{}).
			Where("id = ?", params.BulkPurchaseOrderID).
			Updates(&updates).Error
	})
	if err != nil {
		return nil, err
	}

	order.TrackingStatus = updates.TrackingStatus
	order.LogisticInfo = updates.LogisticInfo

	return order, err
}

type BulkPurchaseOrderMarkFirstPaymentParams struct {
	JwtClaimsInfo models.JwtClaimsInfo

	BulkPurchaseOrderID string `json:"bulk_purchase_order_id" param:"bulk_purchase_order_id" query:"bulk_purchase_order_id" validate:"required"`
}

func (r *BulkPurchaseOrderRepo) BulkPurchaseOrderMarkFirstPayment(params BulkPurchaseOrderMarkFirstPaymentParams) (*models.BulkPurchaseOrder, error) {
	order, err := r.GetBulkPurchaseOrder(GetBulkPurchaseOrderParams{
		BulkPurchaseOrderID: params.BulkPurchaseOrderID,
		JwtClaimsInfo:       params.JwtClaimsInfo,
		IncludeUser:         true,
	})

	if err != nil {
		return nil, err
	}

	var updates models.BulkPurchaseOrder
	err = copier.Copy(&updates, &params)
	if err != nil {
		return nil, err
	}
	updates.FirstPaymentReceivedAt = values.Int64(time.Now().Unix())
	updates.FirstPaymentMarkAsPaidAt = updates.FirstPaymentReceivedAt
	updates.TrackingStatus = enums.BulkPoTrackingStatusFirstPaymentConfirmed

	err = r.db.Transaction(func(tx *gorm.DB) error {
		if order.TrackingStatus != updates.TrackingStatus {
			err = NewBulkPurchaseOrderTrackingRepo(r.db).CreateBulkPurchaseOrderTrackingTx(tx, models.BulkPurchaseOrderTrackingCreateForm{
				PurchaseOrderID: params.BulkPurchaseOrderID,
				ActionType:      enums.BulkPoTrackingActionFirstPaymentConfirmed,
				UserID:          order.UserID,
				CreatedByUserID: params.JwtClaimsInfo.GetUserID(),
				Metadata: &models.PoTrackingMetadata{
					Before: map[string]interface{}{
						"tracking_status": order.TrackingStatus,
					},
					After: map[string]interface{}{
						"tracking_status": updates.TrackingStatus,
					},
				},
			})

			if err != nil {
				return eris.Wrap(err, err.Error())
			}
		}
		return tx.Model(&models.BulkPurchaseOrder{}).
			Where("id = ?", params.BulkPurchaseOrderID).
			Updates(&updates).Error
	})

	order.FirstPaymentReceivedAt = updates.FirstPaymentReceivedAt
	order.FirstPaymentMarkAsPaidAt = updates.FirstPaymentMarkAsPaidAt
	order.TrackingStatus = updates.TrackingStatus
	return order, err
}

type BulkPurchaseOrderMarkFinalPaymentParams struct {
	JwtClaimsInfo models.JwtClaimsInfo

	BulkPurchaseOrderID string `json:"bulk_purchase_order_id" param:"bulk_purchase_order_id" query:"bulk_purchase_order_id" validate:"required"`
}

func (r *BulkPurchaseOrderRepo) BulkPurchaseOrderMarkFinalPayment(params BulkPurchaseOrderMarkFinalPaymentParams) (*models.BulkPurchaseOrder, error) {
	order, err := r.GetBulkPurchaseOrder(GetBulkPurchaseOrderParams{
		BulkPurchaseOrderID: params.BulkPurchaseOrderID,
		JwtClaimsInfo:       params.JwtClaimsInfo,
		IncludeUser:         true,
	})

	if err != nil {
		return nil, err
	}

	var updates models.BulkPurchaseOrder
	err = copier.Copy(&updates, &params)
	if err != nil {
		return nil, err
	}
	updates.FinalPaymentReceivedAt = values.Int64(time.Now().Unix())
	updates.FinalPaymentMarkAsPaidAt = updates.FinalPaymentReceivedAt
	updates.TrackingStatus = enums.BulkPoTrackingStatusFinalPaymentConfirmed

	err = r.db.Transaction(func(tx *gorm.DB) error {
		if order.TrackingStatus != updates.TrackingStatus {
			err = NewBulkPurchaseOrderTrackingRepo(r.db).CreateBulkPurchaseOrderTrackingTx(tx, models.BulkPurchaseOrderTrackingCreateForm{
				PurchaseOrderID: params.BulkPurchaseOrderID,
				ActionType:      enums.BulkPoTrackingActionFinalPaymentConfirmed,
				UserID:          order.UserID,
				CreatedByUserID: params.JwtClaimsInfo.GetUserID(),
				Metadata: &models.PoTrackingMetadata{
					Before: map[string]interface{}{
						"tracking_status": order.TrackingStatus,
					},
					After: map[string]interface{}{
						"tracking_status": updates.TrackingStatus,
					},
				},
			})
			if err != nil {
				return eris.Wrap(err, err.Error())
			}

		}

		return r.db.Model(&models.BulkPurchaseOrder{}).
			Where("id = ?", params.BulkPurchaseOrderID).
			Updates(&updates).Error
	})

	order.FinalPaymentReceivedAt = updates.FinalPaymentReceivedAt
	order.FinalPaymentMarkAsPaidAt = updates.FinalPaymentMarkAsPaidAt
	order.TrackingStatus = updates.TrackingStatus
	return order, err
}

type BulkPurchaseOrderConfirmQcReportParams struct {
	models.JwtClaimsInfo

	BulkPurchaseOrderID string `json:"bulk_purchase_order_id" param:"bulk_purchase_order_id" query:"bulk_purchase_order_id" validate:"required"`

	ShippingFee       price.Price                                `json:"shipping_fee" validate:"required"`
	ShippingAddress   *models.AddressForm                        `json:"shipping_address,omitempty" validate:"required"`
	CommercialInvoice *models.BulkPurchaseOrderCommercialInvoice `json:"commercial_invoice" validate:"required"`
	AdditionalItems   models.BulkPurchaseOrderAdditionalItems    `json:"additional_items"`

	DeductSampleAmount bool `json:"deduct_sample_amount"`
}

func (r *BulkPurchaseOrderRepo) BulkPurchaseOrderConfirmQCReport(params BulkPurchaseOrderConfirmQcReportParams) (*models.BulkPurchaseOrder, error) {
	order, err := r.GetBulkPurchaseOrder(GetBulkPurchaseOrderParams{
		BulkPurchaseOrderID: params.BulkPurchaseOrderID,
		JwtClaimsInfo:       params.JwtClaimsInfo,
		IncludeUser:         true,
	})

	if err != nil {
		return nil, err
	}

	var updates models.BulkPurchaseOrder
	err = copier.Copy(&updates, &params)
	if err != nil {
		return nil, err
	}

	if updates.ShippingAddress != nil {
		updates.ShippingAddress.UserID = order.UserID
		if err = updates.ShippingAddress.CreateOrUpdate(r.db); err == nil {
			updates.ShippingAddressID = updates.ShippingAddress.ID
		}
	}

	updates.TrackingStatus = enums.BulkPoTrackingStatusFinalPayment
	updates.CommercialInvoice = params.CommercialInvoice
	updates.AdditionalItems = params.AdditionalItems

	if updates.CommercialInvoice != nil {
		updates.CommercialInvoice.Status = enums.InvoiceStatusUnpaid
	}

	if order.PurchaseOrder != nil && params.DeductSampleAmount {
		order.SampleDeductionAmount = order.PurchaseOrder.SubTotal.SubPtr(order.PurchaseOrder.ShippingFee).ToPtr()
		updates.SampleDeductionAmount = order.SampleDeductionAmount
	}

	err = r.db.Transaction(func(tx *gorm.DB) error {
		if order.TrackingStatus != updates.TrackingStatus {
			err = NewBulkPurchaseOrderTrackingRepo(r.db).CreateBulkPurchaseOrderTrackingTx(tx, models.BulkPurchaseOrderTrackingCreateForm{
				PurchaseOrderID: params.BulkPurchaseOrderID,
				ActionType:      enums.BulkPoTrackingActionMarkFinalPayment,
				UserID:          order.UserID,
				CreatedByUserID: params.JwtClaimsInfo.GetUserID(),
				Metadata: &models.PoTrackingMetadata{
					Before: map[string]interface{}{
						"tracking_status": order.TrackingStatus,
					},
					After: map[string]interface{}{
						"tracking_status": updates.TrackingStatus,
					},
				},
			})

			if err != nil {
				return eris.Wrap(err, err.Error())
			}

		}

		order.TrackingStatus = updates.TrackingStatus
		order.CommercialInvoice = updates.CommercialInvoice
		err = order.UpdatePrices()
		if err != nil {
			return eris.Wrap(err, err.Error())
		}

		return r.db.Model(&models.BulkPurchaseOrder{}).
			Where("id = ?", params.BulkPurchaseOrderID).
			Updates(&updates).Error

	})
	if err != nil {
		return nil, eris.Wrap(err, err.Error())
	}

	return order, err
}

type BulkPurchaseOrderMarkDeliveredParams struct {
	models.JwtClaimsInfo

	BulkPurchaseOrderID string `json:"bulk_purchase_order_id" param:"bulk_purchase_order_id" query:"bulk_purchase_order_id" validate:"required"`
}

func (r *BulkPurchaseOrderRepo) AdminMarkDelivered(params BulkPurchaseOrderMarkDeliveredParams) (*models.BulkPurchaseOrder, error) {
	order, err := r.GetBulkPurchaseOrder(GetBulkPurchaseOrderParams{
		BulkPurchaseOrderID: params.BulkPurchaseOrderID,
		JwtClaimsInfo:       params.JwtClaimsInfo,
		IncludeUser:         true,
	})
	if err != nil {
		return nil, err
	}

	validStatus := []enums.BulkPoTrackingStatus{enums.BulkPoTrackingStatusDelivering}
	if ok := slices.Contains(validStatus, order.TrackingStatus); !ok {
		return nil, errs.ErrPoInvalidToConfirmDelivered
	}

	var updates models.BulkPurchaseOrder
	err = copier.Copy(&updates, &params)
	if err != nil {
		return nil, err
	}

	updates.TrackingStatus = enums.BulkPoTrackingStatusDelivered
	updates.DeliveredAt = values.Int64(time.Now().Unix())

	err = r.db.Transaction(func(tx *gorm.DB) error {
		err = NewBulkPurchaseOrderTrackingRepo(r.db).CreateBulkPurchaseOrderTrackingTx(tx, models.BulkPurchaseOrderTrackingCreateForm{
			PurchaseOrderID: params.BulkPurchaseOrderID,
			ActionType:      enums.BulkPoTrackingActionDelivered,
			UserID:          order.UserID,
			CreatedByUserID: params.JwtClaimsInfo.GetUserID(),
			Metadata: &models.PoTrackingMetadata{
				Before: map[string]interface{}{
					"tracking_status": order.TrackingStatus,
				},
				After: map[string]interface{}{
					"tracking_status": updates.TrackingStatus,
				},
			},
		})
		if err != nil {
			return err
		}
		return tx.Model(&models.BulkPurchaseOrder{}).
			Where("id = ?", params.BulkPurchaseOrderID).
			Updates(&updates).Error
	})

	if err != nil {
		return nil, eris.Wrap(err, err.Error())
	}

	return order, err
}

type BulkPurchaseOrderUpdatePpsParams struct {
	models.JwtClaimsInfo

	BulkPurchaseOrderID string            `json:"bulk_purchase_order_id" param:"bulk_purchase_order_id" query:"bulk_purchase_order_id" validate:"required"`
	PpsInfo             *models.PoPpsMeta `json:"pps_info" param:"pps_info" query:"pps_info"`
}

func (r *BulkPurchaseOrderRepo) BulkPurchaseOrderUpdatePps(params BulkPurchaseOrderUpdatePpsParams) (*models.BulkPurchaseOrder, error) {
	order, err := r.GetBulkPurchaseOrder(GetBulkPurchaseOrderParams{
		BulkPurchaseOrderID: params.BulkPurchaseOrderID,
		JwtClaimsInfo:       params.JwtClaimsInfo,
		IncludeUser:         true,
	})
	if err != nil {
		return nil, err
	}
	var ppsInfoArr models.PoPpsMetas
	if order.PpsInfo != nil {
		ppsInfoArr = append(ppsInfoArr, *order.PpsInfo...)
	}
	if params.PpsInfo.ID != "" {
		for idx, p := range ppsInfoArr {
			if p.ID == params.PpsInfo.ID {
				ppsInfoArr[idx] = params.PpsInfo
			}
		}
	} else {
		newPps := params.PpsInfo
		newPps.ID = helper.GenerateXID()
		newPps.Status = enums.PpsStatusNone
		ppsInfoArr = append(ppsInfoArr, newPps)
	}

	var updates = models.BulkPurchaseOrder{
		PpsInfo:        &ppsInfoArr,
		TrackingStatus: enums.BulkPoTrackingStatusPps,
	}

	err = r.db.Transaction(func(tx *gorm.DB) error {
		err = NewBulkPurchaseOrderTrackingRepo(r.db).CreateBulkPurchaseOrderTrackingTx(tx, models.BulkPurchaseOrderTrackingCreateForm{
			PurchaseOrderID: order.ID,
			ActionType:      enums.BulkPoTrackingActionUpdatePps,
			UserID:          order.UserID,
			CreatedByUserID: params.JwtClaimsInfo.GetUserID(),
			Metadata: &models.PoTrackingMetadata{
				Before: map[string]interface{}{
					"pps_info": order.PpsInfo,
				},
				After: map[string]interface{}{
					"pps_info": updates.PpsInfo,
				},
			},
		})
		return tx.Model(&models.BulkPurchaseOrder{}).
			Where("id = ?", params.BulkPurchaseOrderID).
			Updates(&updates).Error
	})
	if err != nil {
		return nil, eris.Wrap(err, err.Error())
	}

	order.TrackingStatus = updates.TrackingStatus
	order.PpsInfo = updates.PpsInfo
	return order, err
}

type BulkPurchaseOrderUpdateProductionParams struct {
	models.JwtClaimsInfo

	BulkPurchaseOrderID string                   `json:"bulk_purchase_order_id" param:"bulk_purchase_order_id" query:"bulk_purchase_order_id" validate:"required"`
	ProductionInfo      *models.PoProductionMeta `json:"production_info" param:"production_info" query:"production_info"`
}

func (r *BulkPurchaseOrderRepo) BulkPurchaseOrderUpdateProduction(params BulkPurchaseOrderUpdateProductionParams) (*models.BulkPurchaseOrder, error) {
	order, err := r.GetBulkPurchaseOrder(GetBulkPurchaseOrderParams{
		BulkPurchaseOrderID: params.BulkPurchaseOrderID,
		JwtClaimsInfo:       params.JwtClaimsInfo,
		IncludeUser:         true,
	})
	if err != nil {
		return nil, err
	}

	var updates models.BulkPurchaseOrder
	err = copier.Copy(&updates, &params)
	if err != nil {
		return nil, err
	}

	if order.TrackingStatus == enums.BulkPoTrackingStatusRawMaterial {
		updates.TrackingStatus = enums.BulkPoTrackingStatusProduction
	}

	err = r.db.Transaction(func(tx *gorm.DB) error {
		err = NewBulkPurchaseOrderTrackingRepo(r.db).CreateBulkPurchaseOrderTrackingTx(tx, models.BulkPurchaseOrderTrackingCreateForm{
			PurchaseOrderID: order.ID,
			ActionType:      enums.BulkPoTrackingActionUpdateProduction,
			UserID:          order.UserID,
			CreatedByUserID: params.JwtClaimsInfo.GetUserID(),
			Metadata: &models.PoTrackingMetadata{
				Before: map[string]interface{}{
					"production_info": order.ProductionInfo,
				},
				After: map[string]interface{}{
					"production_info": params.ProductionInfo,
				},
			},
		})
		if err != nil {
			return err
		}
		return tx.Model(&models.BulkPurchaseOrder{}).
			Where("id = ?", params.BulkPurchaseOrderID).
			Updates(&updates).Error
	})

	if err != nil {
		return nil, eris.Wrap(err, err.Error())
	}

	order.TrackingStatus = updates.TrackingStatus

	return order, err
}

func (r *BulkPurchaseOrderRepo) BulkPurchaseOrderAssignPIC(params models.BulkPurchaseOrderAssignPICParam) (updates models.BulkPurchaseOrder, err error) {
	var user models.User
	err = r.db.Model(&models.User{}).Where("id IN ? AND role IN ?", params.AssigneeIDs, []enums.Role{
		enums.RoleLeader,
		enums.RoleStaff,
	}).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = eris.Wrap(errs.ErrUserNotFound, "cannot get assignee")
		return
	}
	if err != nil {
		return
	}

	var bulkPurchaseOrder models.BulkPurchaseOrder
	err = r.db.Select("ID", "AssigneeIDs", "UserID").First(&bulkPurchaseOrder, "id = ?", params.BulkPurchaseOrderID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = errs.ErrPONotFound
		}
		return
	}

	updates.AssigneeIDs = params.AssigneeIDs

	err = r.db.Transaction(func(tx *gorm.DB) error {
		err = tx.Model(&updates).Clauses(clause.Returning{}).
			Where("id = ?", params.BulkPurchaseOrderID).Updates(&updates).Error

		var chatRoom models.ChatRoom
		err = tx.Select("ID").Where(map[string]interface{}{
			"inquiry_id":             updates.InquiryID,
			"purchase_order_id":      updates.PurchaseOrderID,
			"bulk_purchase_order_id": updates.ID,
			"buyer_id":               updates.UserID,
		}).First(&chatRoom).Error
		if err != nil && !r.db.IsRecordNotFoundError(err) {
			return err
		}

		if r.db.IsRecordNotFoundError(err) {
			chatRoom.InquiryID = updates.InquiryID
			chatRoom.PurchaseOrderID = updates.PurchaseOrderID
			chatRoom.BulkPurchaseOrderID = updates.ID
			chatRoom.HostID = params.GetUserID()
			if err := tx.Create(&chatRoom).Error; err != nil {
				return err
			}
		}
		var chatRoomUsers = []*models.ChatRoomUser{{RoomID: chatRoom.ID, UserID: bulkPurchaseOrder.UserID}}
		for _, userId := range params.AssigneeIDs {
			chatRoomUsers = append(chatRoomUsers, &models.ChatRoomUser{RoomID: chatRoom.ID, UserID: userId})
		}
		if err := tx.Clauses(clause.OnConflict{DoNothing: true}).Create(&chatRoomUsers).Error; err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return
	}
	return
}

type BulkPurchaseOrderGetSamplePOParams struct {
	models.JwtClaimsInfo

	BulkPurchaseOrderID string `json:"bulk_purchase_order_id" param:"bulk_purchase_order_id" query:"bulk_purchase_order_id" validate:"required"`
}

func (r *BulkPurchaseOrderRepo) BulkPurchaseOrderGetSamplePO(params BulkPurchaseOrderGetSamplePOParams) (*models.PurchaseOrder, error) {
	var bulkOrder models.BulkPurchaseOrder
	var err = r.db.Select("ID", "UserID", "PurchaseOrderID").First(&bulkOrder, "id = ?", params.BulkPurchaseOrderID).Error
	if err != nil {
		return nil, err
	}
	order, err := NewPurchaseOrderRepo(r.db).GetPurchaseOrder(GetPurchaseOrderParams{
		PurchaseOrderID: bulkOrder.PurchaseOrderID,
		JwtClaimsInfo:   params.JwtClaimsInfo,
		UserID:          bulkOrder.UserID,
	})
	if err != nil {
		return nil, err
	}
	return order, nil
}

type CreateBulkPurchaseInvoiceParams struct {
	models.JwtClaimsInfo
}

func (r *BulkPurchaseOrderRepo) CreateBulkPurchaseInvoice(params CreateBulkPurchaseInvoiceParams) (interface{}, error) {
	return nil, nil
}

func (r *BulkPurchaseOrderRepo) ExportExcel(params PaginateBulkPurchaseOrderParams) (*models.Attachment, error) {
	params.IncludeUser = true
	params.IncludeAssignee = true
	params.IsQueryAll = true
	params.WithoutCount = true
	var result = r.PaginateBulkPurchaseOrder(params)
	if result == nil || result.Records == nil {
		return nil, errors.New("empty response")
	}

	trans, ok := result.Records.([]*models.BulkPurchaseOrder)
	if !ok {
		return nil, errors.New("empty response")
	}

	fileContent, err := models.BulkPurchaseOrders(trans).ToExcel()
	if err != nil {
		return nil, err
	}
	var contentType = models.ContentTypeXLSX
	url := fmt.Sprintf("uploads/bulk_purchase_orders/export/export_bulk_purchase_order_user_%s%s", params.GetUserID(), contentType.GetExtension())
	_, err = s3.New(r.db.Configuration).UploadFile(s3.UploadFileParams{
		Data:        bytes.NewReader(fileContent),
		Bucket:      r.db.Configuration.AWSS3StorageBucket,
		ContentType: string(contentType),
		ACL:         "private",
		Key:         url,
	})
	if err != nil {
		return nil, err
	}
	var resp = models.Attachment{
		FileKey:     url,
		ContentType: string(contentType),
	}
	return &resp, err
}

type ResetBulkPurchaseOrderParams struct {
	models.JwtClaimsInfo

	BulkPurchaseOrderID string `json:"bulk_purchase_order_id" param:"bulk_purchase_order_id" query:"bulk_purchase_order_id" validate:"required"`
}

func (r *BulkPurchaseOrderRepo) ResetBulkPurchaseOrder(params ResetBulkPurchaseOrderParams) (*models.BulkPurchaseOrder, error) {
	order, err := r.GetBulkPurchaseOrder(GetBulkPurchaseOrderParams{
		BulkPurchaseOrderID: params.BulkPurchaseOrderID,
		JwtClaimsInfo:       params.JwtClaimsInfo,
		IncludeUser:         true,
	})
	if err != nil {
		return nil, err
	}

	if order.FirstPaymentIntentID != "" || values.Int64Value(order.FirstPaymentMarkAsPaidAt) > 0 || values.Int64Value(order.FirstPaymentReceivedAt) > 0 {
		return nil, errs.ErrBulkPoFirstPaymentAlreaydPaid
	}

	var updates = models.BulkPurchaseOrder{
		SubmittedAt:    nil,
		TrackingStatus: enums.BulkPoTrackingStatusNew,
		Status:         enums.BulkPurchaseOrderStatusNew,
		QuotationAt:    nil,
		Attachments:    nil,
	}

	err = r.db.Transaction(func(tx *gorm.DB) error {
		err = tx.Select("SubmittedAt", "QuotationAt", "TrackingStatus", "Status", "Attachments").Model(&models.BulkPurchaseOrder{}).Where("id = ?", order.ID).Updates(&updates).Error
		if err != nil {
			return err
		}

		return tx.Unscoped().Delete(&models.BulkPurchaseOrderItem{}, "purchase_order_id = ?", order.ID).Error

	})
	if err != nil {
		return nil, err
	}

	order.SubmittedAt = updates.SubmittedAt
	order.TrackingStatus = updates.TrackingStatus
	order.Status = updates.Status
	order.QuotationAt = updates.QuotationAt
	return order, err
}

type CreateDepositParams struct {
	models.JwtClaimsInfo

	BulkPurchaseOrderID string `json:"bulk_purchase_order_id" param:"bulk_purchase_order_id" query:"bulk_purchase_order_id" validate:"required"`

	DepositAmount price.Price `json:"deposit_amount" validate:"required"`
	DepositNote   string      `json:"deposit_note" validate:"required"`

	IsPaid                bool               `json:"is_paid"`
	TransactionRefID      string             `json:"transaction_ref_id" validate:"required_if=IsPaid true"`
	TransactionAttachment *models.Attachment `json:"transaction_attachment" validate:"required_if=IsPaid true"`
}

func (r *BulkPurchaseOrderRepo) CreateDeposit(params CreateDepositParams) (*models.BulkPurchaseOrder, error) {
	order, err := r.GetBulkPurchaseOrder(GetBulkPurchaseOrderParams{
		BulkPurchaseOrderID: params.BulkPurchaseOrderID,
		JwtClaimsInfo:       params.JwtClaimsInfo,
		IncludeUser:         true,
	})
	if err != nil {
		return nil, err
	}

	if order.DepositPaidAmount.GreaterThan(0) {
		return nil, errs.ErrBulkPoDepositAlreaydPaid
	}

	if params.DepositAmount.LessThanOrEqual(0) {
		return nil, errs.ErrBulkPoDepositAmountInvalid
	}

	stripeConfig, err := stripehelper.GetCurrencyConfig(order.Currency)
	if err != nil {
		return nil, err
	}

	var updates = models.BulkPurchaseOrder{
		DepositNote:   params.DepositNote,
		DepositAmount: &params.DepositAmount,
	}

	if params.IsPaid {
		if params.TransactionRefID != "" && params.TransactionAttachment != nil {
			updates.DepositPaidAmount = &params.DepositAmount
			updates.DepositTransactionAttachment = params.TransactionAttachment
			updates.DepositTransactionRefID = params.TransactionRefID

			err = r.db.Model(&models.BulkPurchaseOrder{}).Where("id = ?", order.ID).Updates(&updates).Error
			if err != nil {
				return nil, err
			}

			order.DepositPaidAmount = updates.DepositPaidAmount
			order.DepositTransactionAttachment = updates.DepositTransactionAttachment
			order.DepositTransactionRefID = updates.DepositTransactionRefID
			return order, err

		} else {
			return nil, errs.ErrBulkPoDepositTransactionReferenceInvalid
		}
	}

	priceItem, err := stripePrice.New(&stripe.PriceParams{
		Currency:   stripe.String(string(order.Currency)),
		UnitAmount: stripe.Int64(params.DepositAmount.MultipleInt(stripeConfig.SmallestUnitFactor).ToInt64()),
		ProductData: &stripe.PriceProductDataParams{
			Name: stripe.String(fmt.Sprintf("Deposit for %s", order.ReferenceID)),
		},
	})
	if err != nil {
		return nil, err
	}

	var transactionFee = params.DepositAmount.
		MultipleFloat64(stripeConfig.TransactionFee).
		Add(price.NewFromFloat(stripeConfig.AdditionalFee))

	feeItem, err := stripePrice.New(&stripe.PriceParams{
		Currency:   stripe.String(string(order.Currency)),
		UnitAmount: stripe.Int64(transactionFee.MultipleInt(stripeConfig.SmallestUnitFactor).ToInt64()),
		ProductData: &stripe.PriceProductDataParams{
			Name: stripe.String("Transaction Fee"),
		},
	})
	if err != nil {
		return nil, err
	}

	link, err := stripehelper.GetInstance().CreatePaymentLink(stripehelper.CreatePaymentLinkParams{
		Currency: order.Currency,
		Metadata: map[string]string{
			"bulk_purchase_order_id":           order.ID,
			"bulk_purchase_order_reference_id": order.ReferenceID,
			"action_source":                    string(stripehelper.ActionSourceBulkPODepositPayment),
			"milestone":                        string(enums.PaymentMilestoneDeposit),
		},
		RedirectURL: fmt.Sprintf("%s/bulks/%s", r.db.Configuration.BrandPortalBaseURL, order.ID),
		LineItems: []*stripe.PaymentLinkLineItemParams{
			{
				Price:    &priceItem.ID,
				Quantity: values.Int64(1),
			},
			{
				Price:    &feeItem.ID,
				Quantity: values.Int64(1),
			},
		},
	})
	if err != nil {
		return nil, err
	}

	updates.DepositPaymentLink = helper.AddURLQuery(link.URL,
		map[string]string{
			"client_reference_id": order.ReferenceID,
			"prefilled_email":     order.User.Email,
		},
	)

	updates.DepositPaymentLinkID = link.ID

	err = r.db.Model(&models.BulkPurchaseOrder{}).Where("id = ?", order.ID).Updates(&updates).Error
	if err != nil {
		return nil, err
	}

	order.DepositPaymentLinkID = updates.DepositPaymentLinkID
	order.DepositPaymentLink = updates.DepositPaymentLink
	order.DepositAmount = updates.DepositAmount
	order.DepositTransactionFee = transactionFee.ToPtr()

	return order, err
}

type CreateBulkPoPaymentLinkParams struct {
	models.JwtClaimsInfo

	BulkPurchaseOrderID string                 `json:"bulk_purchase_order_id" param:"bulk_purchase_order_id" query:"bulk_purchase_order_id" validate:"required"`
	Milestone           enums.PaymentMilestone `json:"milestone" validate:"required"`
}

func (r *BulkPurchaseOrderRepo) CreateBulkPoPaymentLink(params CreateBulkPoPaymentLinkParams) (*models.BulkPurchaseOrder, error) {
	bulkPO, err := r.BulkPurchaseOrderPreviewCheckout(BulkPurchaseOrderPreviewCheckoutParams{
		BulkPurchaseOrderID: params.BulkPurchaseOrderID,
		JwtClaimsInfo:       params.JwtClaimsInfo,
		UpdatePricing:       true,
	})
	if err != nil {
		return nil, err
	}

	stripeConfig, err := stripehelper.GetCurrencyConfig(bulkPO.Currency)
	if err != nil {
		return nil, err
	}

	var priceParams = stripe.PriceParams{
		Currency: stripe.String(string(bulkPO.Currency)),
		ProductData: &stripe.PriceProductDataParams{
			Name: stripe.String(fmt.Sprintf("Bulk PO %s", bulkPO.ReferenceID)),
		},
	}

	var linkParams = stripehelper.CreatePaymentLinkParams{
		Metadata: map[string]string{
			"bulk_purchase_order_id":           bulkPO.ID,
			"bulk_purchase_order_reference_id": bulkPO.ReferenceID,
			"milestone":                        string(enums.PaymentMilestoneDeposit),
		},
		RedirectURL: fmt.Sprintf("%s/bulks/%s", r.db.Configuration.BrandPortalBaseURL, bulkPO.ID),
		LineItems: []*stripe.PaymentLinkLineItemParams{
			{
				Quantity: values.Int64(1),
			},
		},
	}

	switch params.Milestone {
	case enums.PaymentMilestoneFirstPayment:
		if bulkPO.IsFirstPaymentPaid() {
			return nil, errs.ErrBulkPoFirstPaymentAlreaydPaid
		}
		priceParams.UnitAmount = stripe.Int64(bulkPO.FirstPaymentTotal.MultipleInt(stripeConfig.SmallestUnitFactor).ToInt64())
		priceParams.ProductData.Name = stripe.String(fmt.Sprintf("Bulk PO %s - First payment", bulkPO.ReferenceID))
		linkParams.Metadata["action_source"] = string(stripehelper.ActionSourceBulkPOFirstPayment)

	case enums.PaymentMilestoneSecondPayment:
		if bulkPO.IsSecondPaymentPaid() {
			return nil, errs.ErrBulkPoSecondPaymentAlreaydPaid
		}
		priceParams.UnitAmount = stripe.Int64(bulkPO.SecondPaymentTotal.MultipleInt(stripeConfig.SmallestUnitFactor).ToInt64())
		priceParams.ProductData.Name = stripe.String(fmt.Sprintf("Bulk PO %s - Second payment", bulkPO.ReferenceID))
		linkParams.Metadata["action_source"] = string(stripehelper.ActionSourceBulkPOSecondPayment)

	case enums.PaymentMilestoneFinalPayment:
		if bulkPO.IsFinalPaymentPaid() {
			return nil, errs.ErrBulkPoFinalPaymentAlreaydPaid
		}
		priceParams.UnitAmount = stripe.Int64(bulkPO.FinalPaymentTotal.MultipleInt(stripeConfig.SmallestUnitFactor).ToInt64())
		priceParams.ProductData.Name = stripe.String(fmt.Sprintf("Bulk PO %s - Final payment", bulkPO.ReferenceID))
		linkParams.Metadata["action_source"] = string(stripehelper.ActionSourceBulkPOFinalPayment)
	default:
		return nil, errs.ErrBulkPoInvalidToCheckout
	}

	priceItem, err := stripePrice.New(&priceParams)
	if err != nil {
		return nil, err
	}

	linkParams.LineItems[0].Price = &priceItem.ID
	link, err := stripehelper.GetInstance().CreatePaymentLink(linkParams)
	if err != nil {
		return nil, err
	}

	var updates models.BulkPurchaseOrder

	switch params.Milestone {
	case enums.PaymentMilestoneFirstPayment:
		updates.FirstPaymentLink = helper.AddURLQuery(link.URL,
			map[string]string{
				"client_reference_id": bulkPO.ReferenceID,
				"prefilled_email":     bulkPO.User.Email,
			},
		)
		updates.FirstPaymentLinkID = link.ID
		bulkPO.FirstPaymentLink = updates.FirstPaymentLink
		bulkPO.FirstPaymentLinkID = updates.FirstPaymentLinkID
	case enums.PaymentMilestoneSecondPayment:
		updates.SecondPaymentLink = helper.AddURLQuery(link.URL,
			map[string]string{
				"client_reference_id": bulkPO.ReferenceID,
				"prefilled_email":     bulkPO.User.Email,
			},
		)
		updates.SecondPaymentLinkID = link.ID
		bulkPO.SecondPaymentLink = updates.SecondPaymentLink
		bulkPO.SecondPaymentLinkID = updates.SecondPaymentLinkID
	case enums.PaymentMilestoneFinalPayment:
		updates.FinalPaymentLink = helper.AddURLQuery(link.URL,
			map[string]string{
				"client_reference_id": bulkPO.ReferenceID,
				"prefilled_email":     bulkPO.User.Email,
			},
		)
		updates.FinalPaymentLinkID = link.ID
		bulkPO.FinalPaymentLink = updates.FinalPaymentLink
		bulkPO.FinalPaymentLinkID = updates.FinalPaymentLinkID
	default:
		return nil, errs.ErrBulkPoInvalidToCheckout
	}

	err = r.db.Model(&models.BulkPurchaseOrder{}).Where("id = ?", bulkPO.ID).Updates(&updates).Error
	if err != nil {
		return nil, err
	}

	return bulkPO, err
}

type BulkPurchaseOrderStageCommentsParams struct {
	models.JwtClaimsInfo

	BulkPurchaseOrderID string              `json:"bulk_purchase_order_id" param:"bulk_purchase_order_id" query:"bulk_purchase_order_id" validate:"required"`
	Comment             string              `json:"comment,omitempty" validate:"required"`
	Attachments         *models.Attachments `json:"attachments,omitempty"`
}

func (r *BulkPurchaseOrderRepo) StageCommentsCreate(params BulkPurchaseOrderStageCommentsParams) (*models.BulkPurchaseOrder, error) {
	order, err := r.GetBulkPurchaseOrder(GetBulkPurchaseOrderParams{
		BulkPurchaseOrderID: params.BulkPurchaseOrderID,
		JwtClaimsInfo:       params.JwtClaimsInfo,
	})
	if err != nil {
		return nil, err
	}

	err = NewBulkPurchaseOrderTrackingRepo(r.db).CreateBulkPurchaseOrderTrackingTx(r.db.DB, models.BulkPurchaseOrderTrackingCreateForm{
		PurchaseOrderID: params.BulkPurchaseOrderID,
		ActionType:      enums.BulkPoTrackingActionStageComment,
		UserID:          order.UserID,
		CreatedByUserID: params.JwtClaimsInfo.GetUserID(),
		FromStatus:      order.TrackingStatus,
		ToStatus:        order.TrackingStatus,
		Metadata: &models.PoTrackingMetadata{
			After: map[string]interface{}{
				"tracking_status": order.TrackingStatus,
				"comment":         params.Comment,
				"attachments":     params.Attachments,
			},
		},
	})

	if err != nil {
		return nil, eris.Wrap(err, err.Error())
	}

	return order, err
}

type UploadExcelParams struct {
	models.JwtClaimsInfo

	FileKey string `json:"file_key" validate:"required"`
	UserID  string `param:"user_id" validate:"required"`
}

func (r *BulkPurchaseOrderRepo) UploadExcel(params UploadExcelParams) ([]*models.BulkPurchaseOrder, error) {
	data, err := s3.New(r.db.Configuration).GetObject(&s3.GetObjectParams{
		Bucket: r.db.Configuration.AWSS3StorageBucket,
		Key:    params.FileKey,
	})
	if err != nil {
		return nil, err
	}

	result, err := excel.ParseBulkPurchaseOrders(bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}

	var bulks []*models.BulkPurchaseOrder
	var groupID = helper.GenerateBulkPurchaseOrderGroupID()

	err = r.db.Transaction(func(tx *gorm.DB) error {
		for _, record := range result.Records {
			var purchaseOrder = models.PurchaseOrder{
				UserID:            params.UserID,
				AssigneeIDs:       pq.StringArray{params.GetUserID()},
				ClientReferenceID: &record.ReferenceID,
				Currency:          record.Currency,
				MarkAsPaidAt:      values.Int64(time.Now().Unix()),
				PaymentType:       enums.PaymentTypeBankTransfer,
				TrackingStatus:    enums.PoTrackingStatusNew,
				Status:            enums.PurchaseOrderStatusPaid,
				TransactionRefID:  groupID,
			}
			purchaseOrder.ID = helper.GenerateXID()

			var bulkPurchaseOrder = models.BulkPurchaseOrder{
				UserID:          params.UserID,
				AssigneeIDs:     pq.StringArray{params.GetUserID()},
				PurchaseOrderID: purchaseOrder.ID,
				Currency:        record.Currency,
			}
			bulkPurchaseOrder.ID = helper.GenerateXID()
			var subTotalPrice = price.NewFromFloat(0)

			for _, item := range record.Items {
				if item.Qty == 0 {
					continue
				}

				purchaseOrder.Items = append(purchaseOrder.Items, &models.PurchaseOrderItem{
					Sku:        item.SKU,
					Style:      item.Style,
					Size:       item.Size,
					Quantity:   item.Qty,
					NumSamples: item.NumSamples,
					UnitPrice:  item.UnitPrice,
					TotalPrice: item.TotalPrice,
				})

				var bulkPurchaseOrderItem = &models.BulkPurchaseOrderItem{
					Qty:             item.Qty,
					Size:            item.Size,
					Sku:             item.SKU,
					Style:           item.Style,
					UnitPrice:       item.UnitPrice,
					TotalPrice:      item.TotalPrice,
					PurchaseOrderID: bulkPurchaseOrder.ID,
					ColorName:       item.SKU,
				}
				subTotalPrice = subTotalPrice.AddPtr(bulkPurchaseOrderItem.TotalPrice)
				bulkPurchaseOrder.Items = append(bulkPurchaseOrder.Items, bulkPurchaseOrderItem)
				err = tx.Save(&bulkPurchaseOrderItem).Error
				if err != nil {
					return err
				}
			}

			bulkPurchaseOrder.SubTotal = subTotalPrice.ToPtr()
			bulkPurchaseOrder.TaxPercentage = &record.TaxPercentage
			bulkPurchaseOrder.FirstPaymentType = enums.PaymentTypeBankTransfer
			if record.FirstPaymentPercentage == 0 {
				bulkPurchaseOrder.FirstPaymentPercentage = values.Float64(0)
				bulkPurchaseOrder.TrackingStatus = enums.BulkPoTrackingStatusFirstPaymentConfirmed
				bulkPurchaseOrder.FirstPaymentTransferedAt = values.Int64(time.Now().Unix())
				bulkPurchaseOrder.FirstPaymentMarkAsPaidAt = values.Int64(time.Now().Unix())
			}
			err = bulkPurchaseOrder.UpdatePrices()
			if err != nil {
				return err
			}

			err = tx.Create(&purchaseOrder).Error
			if err != nil {
				return err
			}

			err = tx.Create(&bulkPurchaseOrder).Error
			if err != nil {
				return err
			}

			bulkPurchaseOrder.PurchaseOrder = &purchaseOrder

			bulks = append(bulks, &bulkPurchaseOrder)

		}

		return nil
	})
	if err != nil {
		if duplicated, pgErr := r.db.IsDuplicateConstraint(err); duplicated {
			if pgErr.ConstraintName == "purchase_orders_client_reference_id_key" {
				return nil, errs.ErrPOClientReferenceIDDuplicated
			}
		}
		return nil, err
	}

	return bulks, err
}

type UploadBOMParams struct {
	models.JwtClaimsInfo

	BulkPurchaseOrderID string `param:"bulk_purchase_order_id" validate:"required"`
	FileKey             string `json:"file_key" validate:"required"`
}

func (r *BulkPurchaseOrderRepo) UploadBOM(params UploadBOMParams) ([]*models.Bom, error) {
	data, err := s3.New(r.db.Configuration).GetObject(&s3.GetObjectParams{
		Bucket: r.db.Configuration.AWSS3StorageBucket,
		Key:    params.FileKey,
	})
	if err != nil {
		return nil, err
	}

	var bulkPurchaseOrder models.BulkPurchaseOrder
	err = r.db.Select("ID", "ReferenceID", "UserID").First(&bulkPurchaseOrder, "id = ?", params.BulkPurchaseOrderID).Error
	if err != nil {
		return nil, err
	}

	boms, err := excel.ParseBOM(bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}

	var s3Client = s3.New(r.db.Configuration)
	var runnerMain = runner.New(10)
	defer runnerMain.Release()

	var items = make([]*models.Bom, len(boms))

	for i := 0; i < len(boms); i++ {
		i := i

		runnerMain.Submit(func() {
			items[i] = &models.Bom{
				UserID:                bulkPurchaseOrder.UserID,
				BulkPurchaseOrderID:   bulkPurchaseOrder.ID,
				MS:                    boms[i].MS,
				HenryID:               boms[i].HenryID,
				FitSizing:             boms[i].FitSizing,
				ConfirmedColor:        boms[i].ConfirmedColor,
				Cat1:                  boms[i].Cat1,
				Cat2:                  boms[i].Cat2,
				MainFabric:            boms[i].MainFabric,
				MainFabricComposition: boms[i].MainFabricComposition,
				Rivet:                 boms[i].Rivet,
				ZipperTapeColor:       boms[i].ZipperTapeColor,
				ZipperTeeth:           boms[i].ZipperTeeth,
				Thread:                boms[i].Thread,
				OtherDetail:           boms[i].OtherDetail,
				Artwork:               boms[i].Artwork,
				Lining:                boms[i].Lining,
				SustainableHangtag:    boms[i].SustainableHangtag,
			}

			var runnerDownload = runner.New(10)
			defer runnerDownload.Release()
			runnerDownload.Submit(func() {
				for index, image := range boms[i].Image {
					var uploadParams = s3.UploadFileParams{
						Data:        bytes.NewBuffer(image.Data),
						Bucket:      r.db.Configuration.AWSS3StorageBucket,
						Key:         fmt.Sprintf("boms/%s_image%d%s", bulkPurchaseOrder.ReferenceID, index, image.Ext),
						ContentType: string(models.GetContentTypeFromExt(image.Ext)),
						ACL:         "private",
					}

					_, err := s3Client.UploadFile(uploadParams)
					if err != nil {
						r.db.CustomLogger.Errorf("Upload BOM image %s err=%+v", uploadParams.Key, err)
						continue
					}

					r.db.CustomLogger.Debugf("Upload BOM image %s success", uploadParams.Key)
					items[i].Image = append(items[i].Image, &models.Attachment{
						FileKey:     uploadParams.Key,
						ContentType: uploadParams.ContentType,
					})
				}

				if boms[i].ImageLink != "" {
					data, err := helper.DownloadImageFromURL(boms[i].ImageLink)
					if err != nil {
						r.db.CustomLogger.Errorf("Download BOM image %s err=%+v", boms[i].ImageLink, err)
						return
					}

					var imageData = bytes.NewBuffer(data)
					var ext = helper.GetExtensionFromURL(boms[i].ImageLink)
					var uploadParams = s3.UploadFileParams{
						Bucket:      r.db.Configuration.AWSS3StorageBucket,
						Data:        imageData,
						ContentType: string(models.GetContentTypeFromExt(ext)),
						ACL:         "private",
						Key:         fmt.Sprintf("boms/%s_image%s", bulkPurchaseOrder.ReferenceID, ext),
					}
					_, err = s3Client.UploadFile(uploadParams)
					if err != nil {
						r.db.CustomLogger.Errorf("Upload BOM image %s err=%+v", uploadParams.Key, err)
						return
					}

					r.db.CustomLogger.Debugf("Upload BOM image %s success", uploadParams.Key)
					items[i].Image = append(items[i].Image, &models.Attachment{
						FileKey:     uploadParams.Key,
						ContentType: uploadParams.ContentType,
					})
				}

			})

			runnerDownload.Submit(func() {
				for index, image := range boms[i].Buttons {
					var uploadParams = s3.UploadFileParams{
						Data:        bytes.NewBuffer(image.Data),
						Bucket:      r.db.Configuration.AWSS3StorageBucket,
						Key:         fmt.Sprintf("boms/%s_button_%d%s", bulkPurchaseOrder.ReferenceID, index, image.Ext),
						ACL:         "private",
						ContentType: string(models.GetContentTypeFromExt(image.Ext)),
					}

					_, err := s3Client.UploadFile(uploadParams)
					if err != nil {
						r.db.CustomLogger.Errorf("Upload BOM %s buttons err=%+v", uploadParams.Key, err)
						continue
					}

					r.db.CustomLogger.Debugf("Upload BOM %s buttons success", uploadParams.Key)
					items[i].Buttons = append(items[i].Buttons, &models.Attachment{
						FileKey:     uploadParams.Key,
						ContentType: uploadParams.ContentType,
					})
				}

			})

			runnerDownload.Submit(func() {
				for index, image := range boms[i].MainLabel {
					var uploadParams = s3.UploadFileParams{
						Data:        bytes.NewBuffer(image.Data),
						Bucket:      r.db.Configuration.AWSS3StorageBucket,
						Key:         fmt.Sprintf("boms/%s_main_label_%d%s", bulkPurchaseOrder.ReferenceID, index, image.Ext),
						ACL:         "private",
						ContentType: string(models.GetContentTypeFromExt(image.Ext)),
					}

					_, err := s3Client.UploadFile(uploadParams)
					if err != nil {
						r.db.CustomLogger.Errorf("Upload BOM %s main label err=%+v", uploadParams.Key, err)
						continue
					}

					r.db.CustomLogger.Debugf("Upload BOM %s main label success", uploadParams.Key)
					items[i].MainLabel = append(items[i].MainLabel, &models.Attachment{
						FileKey:     uploadParams.Key,
						ContentType: uploadParams.ContentType,
					})
				}

			})

			runnerDownload.Submit(func() {
				for index, image := range boms[i].SizeLabel {
					var uploadParams = s3.UploadFileParams{
						Data:        bytes.NewBuffer(image.Data),
						Bucket:      r.db.Configuration.AWSS3StorageBucket,
						Key:         fmt.Sprintf("boms/%s_size_label_%d%s", bulkPurchaseOrder.ReferenceID, index, image.Ext),
						ACL:         "private",
						ContentType: string(models.GetContentTypeFromExt(image.Ext)),
					}

					_, err := s3Client.UploadFile(uploadParams)
					if err != nil {
						r.db.CustomLogger.Errorf("Upload BOM %s size label err=%+v", uploadParams.Key, err)
						continue
					}

					r.db.CustomLogger.Debugf("Upload BOM %s size label success", uploadParams.Key)
					items[i].SizeLabel = append(items[i].SizeLabel, &models.Attachment{
						FileKey:     uploadParams.Key,
						ContentType: uploadParams.ContentType,
					})
				}

			})

			runnerDownload.Submit(func() {
				for index, image := range boms[i].MainHangtagString {
					var uploadParams = s3.UploadFileParams{
						Data:        bytes.NewBuffer(image.Data),
						Bucket:      r.db.Configuration.AWSS3StorageBucket,
						Key:         fmt.Sprintf("boms/%s_main_hangtag_string_%d%s", bulkPurchaseOrder.ReferenceID, index, image.Ext),
						ACL:         "private",
						ContentType: string(models.GetContentTypeFromExt(image.Ext)),
					}

					_, err := s3Client.UploadFile(uploadParams)
					if err != nil {
						r.db.CustomLogger.Errorf("Upload BOM %s main hangtag err=%+v", uploadParams.Key, err)
						continue
					}

					r.db.CustomLogger.Debugf("Upload BOM %s main hangtag success", uploadParams.Key)
					items[i].MainHangtagString = append(items[i].MainHangtagString, &models.Attachment{
						FileKey:     uploadParams.Key,
						ContentType: uploadParams.ContentType,
					})
				}

			})

			runnerDownload.Submit(func() {
				for index, image := range boms[i].CareLabel {
					var uploadParams = s3.UploadFileParams{
						Data:        bytes.NewBuffer(image.Data),
						Bucket:      r.db.Configuration.AWSS3StorageBucket,
						Key:         fmt.Sprintf("boms/%s_care_label_%d%s", bulkPurchaseOrder.ReferenceID, index, image.Ext),
						ACL:         "private",
						ContentType: string(models.GetContentTypeFromExt(image.Ext)),
					}

					_, err := s3Client.UploadFile(uploadParams)
					if err != nil {
						r.db.CustomLogger.Errorf("Upload BOM %s care label err=%+v", uploadParams.Key, err)
						continue
					}

					r.db.CustomLogger.Debugf("Upload BOM %s care label success", uploadParams.Key)
					items[i].CareLabel = append(items[i].CareLabel, &models.Attachment{
						FileKey:     uploadParams.Key,
						ContentType: uploadParams.ContentType,
					})
				}

			})

			runnerDownload.Submit(func() {
				for index, image := range boms[i].BarcodeSticker {
					var uploadParams = s3.UploadFileParams{
						Data:        bytes.NewBuffer(image.Data),
						Bucket:      r.db.Configuration.AWSS3StorageBucket,
						Key:         fmt.Sprintf("boms/%s_barcode_sticker_%d%s", bulkPurchaseOrder.ReferenceID, index, image.Ext),
						ACL:         "private",
						ContentType: string(models.GetContentTypeFromExt(image.Ext)),
					}

					_, err := s3Client.UploadFile(uploadParams)
					if err != nil {
						r.db.CustomLogger.Errorf("Upload BOM %s barcode sticker err=%+v", uploadParams.Key, err)
						continue
					}

					r.db.CustomLogger.Debugf("Upload BOM %s barcode sticker success", uploadParams.Key)
					items[i].BarcodeSticker = append(items[i].BarcodeSticker, &models.Attachment{
						FileKey:     uploadParams.Key,
						ContentType: uploadParams.ContentType,
					})
				}

			})
			runnerDownload.Wait()
		})
	}

	runnerMain.Wait()

	err = r.db.Clauses(clause.OnConflict{
		Columns: []clause.Column{
			{Name: "bulk_purchase_order_id"},
			{Name: "ms"},
			{Name: "henry_id"},
		},
		UpdateAll: true,
	}).Create(&items).Error

	return items, err
}

type PaginateBOMParams struct {
	models.PaginationParams

	BulkPurchaseOrderID string `param:"bulk_purchase_order_id" validate:"required"`
}

func (r *BulkPurchaseOrderRepo) PaginateBOM(params PaginateBOMParams) *query.Pagination {
	return query.New(r.db, queryfunc.NewBulkPurchaseOrderBomBuilder(queryfunc.BulkPurchaseOrderBomBuilderOptions{})).
		WhereFunc(func(builder *query.Builder) {
			if params.BulkPurchaseOrderID != "" {
				builder.Where("b.bulk_purchase_order_id = ?", params.BulkPurchaseOrderID)
			}
		}).
		Limit(params.Limit).
		Page(params.Page).
		PagingFunc()
}

type GetCheckoutLinkParams struct {
	models.PaginationParams

	BulkPurchaseOrderID string `param:"bulk_purchase_order_id" validate:"required"`
}

type BulkPurchaseOrderUpdateDesignParams struct {
	models.JwtClaimsInfo

	BulkPurchaseOrderID string              `json:"bulk_purchase_order_id" param:"bulk_purchase_order_id" query:"bulk_purchase_order_id" validate:"required"`
	TechpackAttachments *models.Attachments `json:"techpack_attachments" param:"techpack_attachments" query:"techpack_attachments"`
}

func (r *BulkPurchaseOrderRepo) UpdateDesign(params BulkPurchaseOrderUpdateDesignParams) (*models.BulkPurchaseOrder, error) {
	order, err := r.GetBulkPurchaseOrder(GetBulkPurchaseOrderParams{
		BulkPurchaseOrderID: params.BulkPurchaseOrderID,
		UserID:              params.GetUserID(),
		JwtClaimsInfo:       params.JwtClaimsInfo,
	})
	if err != nil {
		return nil, err
	}

	validStatus := []enums.BulkPoTrackingStatus{enums.BulkPoTrackingStatusNew}
	if ok := slices.Contains(validStatus, order.TrackingStatus); !ok {
		return nil, errs.ErrPoInvalidToUploadDesign
	}

	var updates models.BulkPurchaseOrder
	err = copier.Copy(&updates, &params)
	if err != nil {
		return nil, err
	}

	err = r.db.Model(&models.BulkPurchaseOrder{}).
		Where("id = ?", params.BulkPurchaseOrderID).
		Updates(&updates).Error

	if err != nil {
		return nil, eris.Wrap(err, err.Error())
	}

	return order, err
}

type UpdateBulkPurchaseOrderLogsParams struct {
	models.JwtClaimsInfo

	BulkPurchaseOrderID string `json:"bulk_purchase_order_id" query:"bulk_purchase_order_id" form:"bulk_purchase_order_id" param:"bulk_purchase_order_id" validate:"required"`
	LogID               string `json:"log_id" param:"log_id" validate:"required"`

	Notes       string             `json:"notes" validate:"required"`
	Attachments models.Attachments `json:"attachments" validate:"required"`
}

func (r *BulkPurchaseOrderRepo) UpdateBulkPurchaseOrderLogs(params UpdateBulkPurchaseOrderLogsParams) (*models.BulkPurchaseOrderTracking, error) {
	var log = models.BulkPurchaseOrderTracking{
		Notes:       params.Notes,
		Attachments: &params.Attachments,
	}

	var err = r.db.Model(&models.BulkPurchaseOrderTracking{}).
		Where("purchase_order_id = ? AND id = ?", params.BulkPurchaseOrderID, params.LogID).
		Updates(&log).Error

	return &log, err
}

type DeleteBulkPurchaseOrderLogsParams struct {
	models.JwtClaimsInfo

	BulkPurchaseOrderID string `json:"bulk_purchase_order_id" query:"bulk_purchase_order_id" form:"bulk_purchase_order_id" param:"bulk_purchase_order_id" validate:"required"`
	LogID               string `json:"log_id" param:"log_id" validate:"required"`
}

func (r *BulkPurchaseOrderRepo) DeleteBulkPurchaseOrderLogs(params DeleteBulkPurchaseOrderLogsParams) (*models.BulkPurchaseOrderTracking, error) {
	var inquiryAudit = models.BulkPurchaseOrderTracking{
		Notes:       "",
		Attachments: nil,
	}

	var err = r.db.Model(&models.BulkPurchaseOrderTracking{}).
		Select("Notes", "Attachments").
		Where("purchase_order_id = ? AND id = ?", params.BulkPurchaseOrderID, params.LogID).
		Updates(&inquiryAudit).Error

	return &inquiryAudit, err
}

type GetBulkPurchaseOrderInvoiceParams struct {
	models.JwtClaimsInfo

	BulkPurchaseOrderID string `param:"bulk_purchase_order_id" json:"bulk_purchase_order_id" validate:"required"`
	ReCreate            bool   `param:"re_create"`
}

func (r *BulkPurchaseOrderRepo) GetBulkPurchaseOrderInvoice(params GetBulkPurchaseOrderInvoiceParams) (*models.BulkPurchaseOrder, error) {
	var bulkPurchaseOrder models.BulkPurchaseOrder
	var err = r.db.Select("ID", "CommercialInvoiceAttachment", "CommercialInvoice").First(&bulkPurchaseOrder, "id = ?", params.BulkPurchaseOrderID).Error
	if err != nil {
		return nil, err
	}

	if bulkPurchaseOrder.CommercialInvoiceAttachment != nil && !params.ReCreate {
		return &bulkPurchaseOrder, err
	}

	result, err := NewInvoiceRepo(r.db).CreateBulkPurchaseOrderInvoice(CreateBulkPurchaseOrderInvoiceParams{
		BulkPurchaseOrderID: bulkPurchaseOrder.ID,
		ReCreate:            params.ReCreate,
	})
	if err != nil {
		if eris.Is(err, errs.ErrPOInvoiceAlreadyGenerated) && result.CommercialInvoiceAttachment != nil {
			goto Susccess
		}
		return nil, err
	}

	if result.CommercialInvoiceAttachment == nil {
		return nil, eris.New("Invoice is not able to generate")
	}

Susccess:
	return result, nil
}

type BulkPurchaseOrderFeedbackParams struct {
	BulkPurchaseOrderID string `json:"bulk_purchase_order_id" param:"bulk_purchase_order_id"`
	models.BulkPurchaseOrderFeedback
}

func (r *BulkPurchaseOrderRepo) BuyerGiveBulkPurchaseOrderFeedback(params BulkPurchaseOrderFeedbackParams) error {
	var updates = &models.BulkPurchaseOrder{
		Feedback: &params.BulkPurchaseOrderFeedback,
	}
	return r.db.Model(&models.BulkPurchaseOrder{}).Where("id = ?", params.BulkPurchaseOrderID).Updates(updates).Error
}

type AdminAssignBulkPOMakerParams struct {
	models.JwtClaimsInfo

	BulkPurchaseOrderID string `json:"bulk_purchase_order_id" param:"bulk_purchase_order_id" query:"bulk_purchase_order_id" validate:"required"`

	MakerID string `json:"maker_id"`
}

func (r *BulkPurchaseOrderRepo) AdminAssignBulkPOMaker(params AdminAssignBulkPOMakerParams) (*models.BulkPurchaseOrder, error) {
	bulkPO, err := r.GetBulkPurchaseOrder(GetBulkPurchaseOrderParams{
		JwtClaimsInfo:       params.JwtClaimsInfo,
		BulkPurchaseOrderID: params.BulkPurchaseOrderID,
	})
	if err != nil {
		return nil, err
	}

	var sellerQuotation models.BulkPurchaseOrderSellerQuotation

	if params.MakerID != "inflow" {
		var seller models.User
		err = r.db.Select("ID", "Email", "CompanyName").First(&seller, "id = ?", params.MakerID).Error
		if err != nil {
			return nil, err
		}

		err = r.db.Select("ID", "UserID", "BulkPurchaseOrderID").First(&sellerQuotation, "user_id = ? AND bulk_purchase_order_id = ?", seller.ID, params.BulkPurchaseOrderID).Error
		if err != nil {
			return nil, err
		}

		bulkPO.Seller = &seller
	}

	err = r.db.Transaction(func(tx *gorm.DB) error {
		err = tx.Model(&models.BulkPurchaseOrder{}).
			Where("id = ?", params.BulkPurchaseOrderID).
			Update("SellerID", params.MakerID).Error
		if err != nil {
			return err
		}

		if params.MakerID != "inflow" && sellerQuotation.ID != "" {
			var updates = models.BulkPurchaseOrderSellerQuotation{
				BulkPurchaseOrderID: params.BulkPurchaseOrderID,
				Status:              enums.BulkPurchaseOrderSellerStatusApproved,
			}

			err = tx.Model(&models.BulkPurchaseOrderSellerQuotation{}).Where("id = ?", sellerQuotation.ID).
				Updates(&updates).Error

		}

		return err
	})
	if err != nil {
		return nil, err
	}

	return bulkPO, err
}

type PaginateBulkPurchaseOrderAllocationParams struct {
	models.PaginationParams
	models.JwtClaimsInfo

	BulkPurchaseOrderID string `param:"bulk_purchase_order_id" validate:"required"`

	Statuses                []enums.InquirySellerStatus           `query:"statuses,omitempty"`
	SellerQuotationStatuses []enums.BulkPurchaseOrderSellerStatus `query:"seller_quotation_statuses,omitempty"`
}

func (r *BulkPurchaseOrderRepo) PaginateBulkPurchaseOrderAllocation(params PaginateBulkPurchaseOrderAllocationParams) *query.Pagination {
	var builder = queryfunc.NewBulkPurchaseOrderAllocationBuilder(queryfunc.BulkPurchaseOrderAllocationBuilderOptions{
		QueryBuilderOptions: queryfunc.QueryBuilderOptions{
			Role: params.GetRole(),
		},
	})

	var result = query.New(r.db, builder).
		WhereFunc(func(builder *query.Builder) {
			builder.Where("u.role = ?", enums.RoleSeller)

			builder.Where("bposq.bulk_purchase_order_id = ?", params.BulkPurchaseOrderID)

			if len(params.Statuses) > 0 {
				builder.Where("bposq.status IN ?", params.Statuses)
			}

			if len(params.SellerQuotationStatuses) > 0 {
				builder.Where("bposq.status IN ?", params.SellerQuotationStatuses)
			}

			if keyword := strings.TrimSpace(params.Keyword); keyword != "" {
				var q = "%" + keyword + "%"
				builder.Where("(u.name ILIKE @keyword OR u.company_name ILIKE @keyword OR u.email ILIKE @keyword)", sql.Named("keyword", q))
			}
		}).
		Page(params.Page).
		Limit(params.Limit).
		PagingFunc()

	return result

}

func (r *BulkPurchaseOrderRepo) CreateMultipleBulkPurchaseOrders(req *models.CreateMultipleBulkPurchaseOrdersRequest) ([]*models.BulkPurchaseOrder, error) {
	var user models.User
	if err := r.db.First(&user, "id = ?", req.UserID).Error; err != nil {
		if r.db.IsRecordNotFoundError(err) {
			return nil, errs.ErrUserNotFound
		}
		return nil, err
	}

	var orderGroupIDs = make([]string, 0, len(req.Bulks))
	for _, bulk := range req.Bulks {
		if bulk.OrderGroupID != "" {
			orderGroupIDs = append(orderGroupIDs, bulk.OrderGroupID)
		}
	}
	var orderGroups = make(models.OrderGroups, 0, len(orderGroupIDs))
	if err := r.db.Select("ID").Find(&orderGroups, "id IN ?", orderGroupIDs).Error; err != nil {
		return nil, err
	}
	var dbOrderGroupIDs = orderGroups.IDs()
	for _, id := range orderGroupIDs {
		if !helper.StringContains(dbOrderGroupIDs, id) {
			return nil, eris.Wrapf(errs.ErrOrderGroupNotFound, "order_group_id:%s", id)
		}
	}

	var mapClientRefID = make(map[string]struct{}, len(req.Bulks))
	for _, bulk := range req.Bulks {
		if bulk.PurchaseOrderClientReferenceID != "" {
			if _, ok := mapClientRefID[bulk.PurchaseOrderClientReferenceID]; ok {
				return nil, errs.ErrPORefIDDuplicate
			}
			mapClientRefID[bulk.PurchaseOrderClientReferenceID] = struct{}{}
		}
	}

	var coordinateToCreate *models.Coordinate
	var shippingAddressToCreate *models.Address
	if req.ShippingMethod != enums.ShippingMethodEXW {
		if req.AddressCoordinate == nil {
			return nil, errs.ErrAddressRequired
		}
		coordinateToCreate = req.AddressCoordinate
		coordinateToCreate.ID = utils.MD5(coordinateToCreate.ToJsonString())
		shippingAddressToCreate = &models.Address{
			CoordinateID: coordinateToCreate.ID,
			UserID:       user.ID,
			Name:         user.Name,
			PhoneNumber:  user.PhoneNumber,
		}
		shippingAddressToCreate.ID = shippingAddressToCreate.GenerateID()
	}

	var bulksToCreate = make(models.BulkPurchaseOrders, 0, len(req.Bulks))
	var orderCartItemsToCreate models.OrderCartItems
	var purchaseOrdersToCreate = make(models.PurchaseOrders, 0, len(req.Bulks))

	for idx, bulk := range req.Bulks {
		var shippingAddressID = ""
		if shippingAddressToCreate != nil {
			shippingAddressID = shippingAddressToCreate.ID
		}
		// purchase order to create
		var purchaseOrder models.PurchaseOrder
		if len(bulk.PurchaseOrderItems) > 0 {
			purchaseOrder = models.PurchaseOrder{
				ID:                helper.GenerateXID(),
				UserID:            user.ID,
				ProductName:       bulk.ProductName,
				Attachments:       bulk.Attachments,
				Currency:          req.Currency,
				Status:            enums.PurchaseOrderStatusPending,
				TrackingStatus:    enums.PoTrackingStatusNew,
				AssigneeIDs:       user.ContactOwnerIDs,
				OrderGroupID:      bulk.OrderGroupID,
				ShippingAddressID: shippingAddressID,
				Quotations:        bulk.Quotations,
				PaymentType:       enums.PaymentTypeBankTransfer,
				Pricing: models.Pricing{
					TaxPercentage: &bulk.PurchaseOrderTaxPercentage,
					ShippingFee:   &bulk.PurchaseOrderShippingFee,
				},
			}
			if bulk.PurchaseOrderClientReferenceID != "" {
				purchaseOrder.ClientReferenceID = &bulk.PurchaseOrderClientReferenceID
			}
			// update purchase order prices
			sampleQuotation, ok := lo.Find(bulk.Quotations, func(item *models.InquiryQuotationItem) bool {
				return item.Type == enums.InquiryTypeSample
			})
			if ok && sampleQuotation.Price.GreaterThan(0) {
				var subTotalPrice = price.NewFromFloat(0)

				for _, item := range bulk.PurchaseOrderItems {
					var unitPrice = sampleQuotation.Price
					var totalPrice = unitPrice.MultipleInt(int64(item.Qty))
					subTotalPrice = subTotalPrice.Add(totalPrice)

					orderCartItemsToCreate = append(orderCartItemsToCreate, &models.OrderCartItem{
						PurchaseOrderID: purchaseOrder.ID,
						Size:            item.Size,
						ColorName:       item.ColorName,
						Style:           item.Style,
						Qty:             item.Qty,
						UnitPrice:       unitPrice,
						TotalPrice:      totalPrice,
						NoteToSupplier:  item.NoteToSupplier,
					})
				}
				purchaseOrder.SubTotal = &subTotalPrice
				if err := purchaseOrder.UpdatePrices(); err != nil {
					return nil, err
				}
			} else {
				purchaseOrder.Status = enums.PurchaseOrderStatusPaid
			}
			purchaseOrdersToCreate = append(purchaseOrdersToCreate, &purchaseOrder)

		}
		// bulk to create
		var bulkCreate = models.BulkPurchaseOrder{
			Model: models.Model{
				ID: helper.GenerateXID(),
			},
			UserID:                 user.ID,
			PurchaseOrderID:        purchaseOrder.ID,
			ProductName:            bulk.ProductName,
			Note:                   bulk.Note,
			Attachments:            bulk.Attachments,
			TechpackAttachments:    bulk.TechpackAttachments,
			SizeAttachments:        bulk.SizeAttachments,
			OrderGroupID:           bulk.OrderGroupID,
			AdminQuotations:        bulk.Quotations,
			Status:                 enums.BulkPurchaseOrderStatusNew,
			TrackingStatus:         enums.BulkPoTrackingStatusFirstPayment,
			Currency:               req.Currency,
			ShippingMethod:         req.ShippingMethod,
			ShippingAttachments:    &req.ShippingAttachments,
			ShippingNote:           &req.ShippingNote,
			PackingAttachments:     &req.PackingAttachments,
			PackingNote:            &req.PackingNote,
			ShippingAddressID:      shippingAddressID,
			AssigneeIDs:            user.ContactOwnerIDs,
			FirstPaymentPercentage: &req.FirstPaymentPercentage,
			FirstPaymentType:       enums.PaymentTypeBankTransfer,
			Pricing: models.Pricing{
				TaxPercentage: &req.TaxPercentage,
				ShippingFee:   &req.ShippingFee,
			},
		}
		if req.FirstPaymentPercentage == 0 {
			bulkCreate.TrackingStatus = enums.BulkPoTrackingStatusFirstPaymentConfirmed
		}
		// update bulk prices
		var subTotalPrice = price.NewFromFloat(0)
		bulkQuotation, ok := lo.Find(bulk.Quotations, func(item *models.InquiryQuotationItem) bool {
			return item.Type == enums.InquiryTypeBulk
		})
		if !ok {
			return nil, eris.Wrapf(errs.ErrBulkQuotationEmpty, "bulk index:%d", idx)
		}
		for _, item := range bulk.Items {
			var unitPrice = bulkQuotation.Price
			var totalPrice = unitPrice.MultipleInt(int64(item.Qty))
			subTotalPrice = subTotalPrice.Add(totalPrice)

			orderCartItemsToCreate = append(orderCartItemsToCreate, &models.OrderCartItem{
				BulkPurchaseOrderID: bulkCreate.ID,
				Size:                item.Size,
				ColorName:           item.ColorName,
				Style:               item.Style,
				Qty:                 item.Qty,
				NoteToSupplier:      item.NoteToSupplier,
				UnitPrice:           unitPrice,
				TotalPrice:          totalPrice,
			})
		}
		bulkCreate.SubTotal = &subTotalPrice
		if err := bulkCreate.UpdatePrices(); err != nil {
			return nil, err
		}
		bulkCreate.PurchaseOrder = &purchaseOrder
		bulksToCreate = append(bulksToCreate, &bulkCreate)
	}

	if err := r.db.Transaction(func(tx *gorm.DB) error {
		if coordinateToCreate != nil {
			if err := tx.Clauses(clause.OnConflict{
				Columns:   []clause.Column{{Name: "id"}},
				UpdateAll: true,
			}).Create(coordinateToCreate).Error; err != nil {
				return err
			}
		}

		if shippingAddressToCreate != nil {
			if err := tx.Clauses(clause.OnConflict{
				Columns:   []clause.Column{{Name: "id"}},
				UpdateAll: true,
			}).Create(shippingAddressToCreate).Error; err != nil {
				return err
			}
		}
		if err := tx.Create(&orderCartItemsToCreate).Error; err != nil {
			return err
		}
		if err := tx.Create(&bulksToCreate).Error; err != nil {
			return err
		}
		if len(purchaseOrdersToCreate) > 0 {
			if err := tx.Create(&purchaseOrdersToCreate).Error; err != nil {
				return err
			}
		}
		return nil
	}); err != nil {
		return nil, err
	}
	return bulksToCreate, nil
}

type AdminSellerBulkPurchaseOrderUploadPoParams struct {
	models.JwtClaimsInfo

	BulkPurchaseOrderID string               `json:"bulk_purchase_order_id" param:"bulk_purchase_order_id" query:"bulk_purchase_order_id" validate:"required"`
	SellerPoAttachments models.PoAttachments `json:"seller_po_attachments,omitempty" validate:"required"`
}

func (r *BulkPurchaseOrderRepo) AdminSellerBulkPurchaseOrderUploadPo(params AdminSellerBulkPurchaseOrderUploadPoParams) (*models.BulkPurchaseOrder, error) {
	order, err := r.GetBulkPurchaseOrder(GetBulkPurchaseOrderParams{
		BulkPurchaseOrderID: params.BulkPurchaseOrderID,
		JwtClaimsInfo:       params.JwtClaimsInfo,
	})
	if err != nil {
		return nil, err
	}

	var checkStatuses = []enums.SellerBulkPoTrackingStatus{
		enums.SellerBulkPoTrackingStatusPORejected,
		enums.SellerBulkPoTrackingStatusPO,
	}

	if !lo.Contains(checkStatuses, order.SellerTrackingStatus) {
		return nil, errs.ErrBulkPoNotAbleToUpdatePO
	}
	var updates = models.BulkPurchaseOrder{
		SellerPoAttachments:  &params.SellerPoAttachments,
		SellerTrackingStatus: enums.SellerBulkPoTrackingStatusPO,
	}

	err = r.db.Model(&models.BulkPurchaseOrder{}).
		Where("id = ?", params.BulkPurchaseOrderID).
		Updates(&updates).Error

	if err != nil {
		return nil, eris.Wrap(err, err.Error())
	}

	order.SellerPoAttachments = updates.SellerPoAttachments
	order.SellerTrackingStatus = updates.SellerTrackingStatus

	return order, err
}

func (r *BulkPurchaseOrderRepo) SubmitMultipleBulkQuotations(req *models.SubmitMultipleBulkQuotationsRequest) (models.BulkPurchaseOrders, error) {
	var bulkIDs = make([]string, 0, len(req.Quotations))
	for _, quotation := range req.Quotations {
		bulkIDs = append(bulkIDs, quotation.BulkPurchaseOrderID)
	}
	var bulks = make(models.BulkPurchaseOrders, 0, len(bulkIDs))
	if err := r.db.Find(&bulks, "id IN ?", bulkIDs).Error; err != nil {
		return nil, err
	}
	var dbBulkIDs = bulks.IDs()
	for _, id := range bulkIDs {
		if !helper.StringContains(dbBulkIDs, id) {
			return nil, eris.Wrapf(errs.ErrBulkPoNotFound, "bulk_id:%s", id)
		}
	}
	for _, bulk := range bulks {
		if bulk.TrackingStatus != enums.BulkPoTrackingStatusWaitingForQuotation {
			return nil, eris.Wrapf(errs.ErrBulkPoInvalidToSendQuotationToBuyer, "bulk_id:%s", bulk.ID)
		}
	}

	var orderItems models.OrderCartItems
	if err := r.db.Find(&orderItems, "bulk_purchase_order_id IN ?", bulkIDs).Error; err != nil {
		return nil, err
	}
	var mapBulkIdToOrderItems = make(map[string]models.OrderCartItems, len(orderItems))
	for _, item := range orderItems {
		_, ok := mapBulkIdToOrderItems[item.BulkPurchaseOrderID]
		if ok {
			mapBulkIdToOrderItems[item.BulkPurchaseOrderID] = append(mapBulkIdToOrderItems[item.BulkPurchaseOrderID], item)
		} else {
			mapBulkIdToOrderItems[item.BulkPurchaseOrderID] = []*models.OrderCartItem{item}
		}
	}

	inquiryIDs := bulks.InquiryIDs()
	inquiries := make(models.Inquiries, 0, len(inquiryIDs))
	if err := r.db.Select("ID", "TaxPercentage").Find(&inquiries, "id IN ?", inquiryIDs).Error; err != nil {
		return nil, err
	}
	dbInquiryIDs := inquiries.IDs()
	for _, id := range inquiryIDs {
		if !helper.StringContains(dbInquiryIDs, id) {
			return nil, eris.Wrapf(errs.ErrInquiryNotFound, "inquiry_id:%s", id)
		}
	}
	var mapInquiryIdToInquiry = make(map[string]*models.Inquiry, len(inquiries))
	for _, iq := range inquiries {
		mapInquiryIdToInquiry[iq.ID] = iq
	}

	purchaseOrderIDs := bulks.PurchaseOrderIDs()
	purchaseOrders := make(models.PurchaseOrders, 0, len(purchaseOrderIDs))
	if err := r.db.Find(&purchaseOrders, "id IN ? AND status = ?", purchaseOrderIDs, enums.PurchaseOrderStatusPending).Error; err != nil {
		return nil, err
	}
	var mapPoIdToPurchaseOrder = make(map[string]*models.PurchaseOrder, len(purchaseOrders))
	for _, po := range purchaseOrders {
		mapPoIdToPurchaseOrder[po.ID] = po
	}

	var poOrderItems models.OrderCartItems
	if err := r.db.Find(&poOrderItems, "purchase_order_id IN ?", purchaseOrders.IDs()).Error; err != nil {
		return nil, err
	}
	var mapPoIdToOrderItems = make(map[string]models.OrderCartItems, len(orderItems))
	for _, item := range poOrderItems {
		_, ok := mapPoIdToOrderItems[item.PurchaseOrderID]
		if ok {
			mapPoIdToOrderItems[item.PurchaseOrderID] = append(mapPoIdToOrderItems[item.PurchaseOrderID], item)
		} else {
			mapPoIdToOrderItems[item.PurchaseOrderID] = []*models.OrderCartItem{item}
		}
	}

	var bulksToUpdate = make(models.BulkPurchaseOrders, 0, len(req.Quotations))
	var purchaseOrdersToUpdate = make(models.PurchaseOrders, 0, len(req.Quotations))
	var orderCartItemsToUpdate = make(models.OrderCartItems, 0, len(orderItems))

	var bulkTrackings []*models.BulkPurchaseOrderTracking

	for _, quotation := range req.Quotations {
		bulk, ok := lo.Find(bulks, func(item *models.BulkPurchaseOrder) bool {
			return item.ID == quotation.BulkPurchaseOrderID
		})
		if !ok {
			return nil, eris.Wrapf(errs.ErrBulkPoNotFound, "bulk_id:%s", quotation.BulkPurchaseOrderID)
		}
		orderItems, ok := mapBulkIdToOrderItems[bulk.ID]
		if !ok {
			return nil, eris.Wrapf(errs.ErrOrderItemNotFound, "bulk_purchase_order_id:%s", bulk.ID)
		}
		inquiry, ok := mapInquiryIdToInquiry[bulk.InquiryID]
		if !ok && bulk.InquiryID != "" {
			return nil, eris.Wrapf(errs.ErrInquiryNotFound, "bulk_purchase_order_id:%s", bulk.ID)
		}

		bulk.AdminQuotations = quotation.AdminQuotations
		bulk.FirstPaymentPercentage = &quotation.FirstPaymentPercentage
		bulk.QuotationNote = quotation.QuotationNote
		bulk.QuotationNoteAttachments = quotation.QuotationNoteAttachments
		bulk.QuotationAt = values.Int64(time.Now().Unix())
		bulk.TrackingStatus = enums.BulkPoTrackingStatusFirstPayment
		if quotation.FirstPaymentPercentage == 0 {
			bulk.TrackingStatus = enums.BulkPoTrackingStatusFirstPaymentConfirmed
		}

		//update bulk checkout prices
		var subTotalPrice = price.NewFromFloat(0)
		bulkQuotation, ok := lo.Find(quotation.AdminQuotations, func(item *models.InquiryQuotationItem) bool {
			return item.Type == enums.InquiryTypeBulk
		})
		if !ok {
			return nil, eris.Wrapf(errs.ErrBulkQuotationEmpty, "bulk_purchase_order_id:%s", bulk.ID)
		}
		for _, item := range orderItems {
			item.UnitPrice = bulkQuotation.Price
			item.TotalPrice = item.UnitPrice.MultipleInt(int64(item.Qty))
			subTotalPrice = subTotalPrice.Add(item.TotalPrice)

			orderCartItemsToUpdate = append(orderCartItemsToUpdate, item)
		}
		bulk.SubTotal = subTotalPrice.ToPtr()
		bulk.FirstPaymentType = enums.PaymentTypeBankTransfer
		if inquiry != nil {
			bulk.TaxPercentage = inquiry.TaxPercentage
		}
		if bulk.CommercialInvoice != nil {
			bulk.TaxPercentage = bulk.CommercialInvoice.TaxPercentage
		}
		if err := bulk.UpdatePrices(); err != nil {
			return nil, err
		}

		bulksToUpdate = append(bulksToUpdate, bulk)

		// update purchase checkout prices
		purchaseOrder, ok := mapPoIdToPurchaseOrder[bulk.PurchaseOrderID]
		if ok {
			var poSubTotalPrice = price.NewFromFloat(0)
			sampleQuotation, ok := lo.Find(quotation.AdminQuotations, func(item *models.InquiryQuotationItem) bool {
				return item.Type == enums.InquiryTypeSample
			})
			if !ok {
				return nil, eris.Wrapf(errs.ErrBulkQuotationEmpty, "bulk_purchase_order_id:%s", bulk.ID)
			}
			if sampleQuotation.Price.GreaterThan(0) {
				orderItems, ok := mapPoIdToOrderItems[purchaseOrder.ID]
				if !ok {
					return nil, eris.Wrapf(errs.ErrOrderItemNotFound, "purchase_order_id:%s", purchaseOrder.ID)
				}
				for _, item := range orderItems {
					item.UnitPrice = sampleQuotation.Price
					item.TotalPrice = item.UnitPrice.MultipleInt(int64(item.Qty))
					subTotalPrice = subTotalPrice.Add(item.TotalPrice)

					orderCartItemsToUpdate = append(orderCartItemsToUpdate, item)
				}
				purchaseOrder.SubTotal = poSubTotalPrice.ToPtr()
				purchaseOrder.TaxPercentage = bulk.TaxPercentage
				purchaseOrder.ShippingFee = bulk.ShippingFee
				purchaseOrder.PaymentType = enums.PaymentTypeBankTransfer
				if err := purchaseOrder.UpdatePrices(); err != nil {
					return nil, err
				}
			} else {
				purchaseOrder.Status = enums.PurchaseOrderStatusPaid
			}
			purchaseOrdersToUpdate = append(purchaseOrdersToUpdate, purchaseOrder)
		}

		bulkTrackings = append(bulkTrackings, &models.BulkPurchaseOrderTracking{
			PurchaseOrderID: bulk.ID,
			ActionType:      enums.BulkPoTrackingActionSubmitQuotation,
			UserID:          bulk.UserID,
			CreatedByUserID: req.GetUserID(),
			Description:     "Admin has sent quotation to buyer",
			Metadata: &models.PoTrackingMetadata{
				After: map[string]interface{}{
					"quotations": bulk.AdminQuotations,
				},
			},
		})
	}
	if err := r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "id"}},
			UpdateAll: true,
		}).Create(&bulksToUpdate).Error; err != nil {
			return err
		}

		if len(purchaseOrdersToUpdate) > 0 {
			if err := tx.Clauses(clause.OnConflict{
				Columns:   []clause.Column{{Name: "id"}},
				UpdateAll: true,
			}).Create(&purchaseOrdersToUpdate).Error; err != nil {
				return err
			}
		}

		if err := tx.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "id"}},
			DoUpdates: clause.AssignmentColumns([]string{"unit_price", "total_price"}),
		}).Create(&orderCartItemsToUpdate).Error; err != nil {
			return err
		}

		return tx.Create(&bulkTrackings).Error
	}); err != nil {
		return nil, err
	}
	return bulksToUpdate, nil
}

func (r *BulkPurchaseOrderRepo) UploadBulks(req *models.UploadBulksRequest) (*models.UploadBulksResponse, error) {
	data, err := s3.New(r.db.Configuration).GetObject(&s3.GetObjectParams{
		Bucket: r.db.Configuration.AWSS3StorageBucket,
		Key:    req.FileKey,
	})
	if err != nil {
		return nil, err
	}

	result, err := excel.ParseBulkPurchaseOrdersV2(bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}
	return &models.UploadBulksResponse{
		CreateMultipleBulkPurchaseOrdersRequest: *result,
	}, nil

}
