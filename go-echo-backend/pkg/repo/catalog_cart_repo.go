package repo

import (
	"fmt"
	"strings"
	"time"

	"github.com/engineeringinflow/inflow-backend/pkg/db"
	"github.com/engineeringinflow/inflow-backend/pkg/errs"
	"github.com/engineeringinflow/inflow-backend/pkg/helper"
	"github.com/engineeringinflow/inflow-backend/pkg/logger"
	"github.com/engineeringinflow/inflow-backend/pkg/models"
	"github.com/engineeringinflow/inflow-backend/pkg/models/enums"
	"github.com/engineeringinflow/inflow-backend/pkg/models/price"
	"github.com/engineeringinflow/inflow-backend/pkg/repo/query"
	"github.com/engineeringinflow/inflow-backend/pkg/repo/query/queryfunc"
	"github.com/engineeringinflow/inflow-backend/pkg/stripehelper"
	"github.com/jinzhu/copier"
	"github.com/rotisserie/eris"
	"github.com/samber/lo"
	"github.com/stripe/stripe-go/v74"
	"github.com/thaitanloi365/go-utils/values"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type CatalogCartRepo struct {
	db     *db.DB
	logger *logger.Logger
}

func NewCatalogCartRepo(db *db.DB) *CatalogCartRepo {
	return &CatalogCartRepo{
		db:     db,
		logger: logger.New("repo/catalog_cart"),
	}
}

type PaginateCatalogCartsParams struct {
	models.PaginationParams
	models.JwtClaimsInfo
}

func (r *CatalogCartRepo) PaginateCatalogCarts(params PaginateCatalogCartsParams) *query.Pagination {
	var builder = queryfunc.NewCatalogCartBuilder(queryfunc.CatalogCartBuilderOptions{
		QueryBuilderOptions: queryfunc.QueryBuilderOptions{
			Role: params.GetRole(),
		},
	})

	return query.New(r.db, builder).
		WhereFunc(func(builder *query.Builder) {
			builder.Where("po.is_cart = ?", true)

			builder.Where("po.user_id = ?", params.GetUserID())
		}).
		Limit(params.Limit).
		Page(params.Page).
		PagingFunc()

}

type GetCatalogCartParams struct {
	models.JwtClaimsInfo

	CartID string `json:"cart_id"`
	UserID string `json:"user_id"`
}

func (r *CatalogCartRepo) GetCatalogCart(params GetCatalogCartParams) (*models.PurchaseOrder, error) {
	var record models.PurchaseOrder
	var builder = queryfunc.NewCatalogCartBuilder(queryfunc.CatalogCartBuilderOptions{
		QueryBuilderOptions: queryfunc.QueryBuilderOptions{
			Role: params.GetRole(),
		},
	})

	var err = query.New(r.db, builder).
		WhereFunc(func(builder *query.Builder) {
			if params.UserID != "" {
				builder.Where("po.user_id = ?", params.UserID)
			}

			if params.CartID != "" {
				builder.Where("po.id = ?", params.CartID)
			}
		}).
		Limit(1).
		FirstFunc(&record)

	return &record, err

}

type UpdateCatalogCartsParams struct {
	models.JwtClaimsInfo

	Records []UpdateCatalogCartParams `json:"records"`
}

func (r *CatalogCartRepo) UpdateCatalogCarts(params UpdateCatalogCartsParams) ([]*models.PurchaseOrder, error) {
	var result []*models.PurchaseOrder
	var err = r.db.Transaction(func(tx *gorm.DB) error {
		for _, v := range params.Records {
			po, err := r.UpdateCatalogCart(tx, v)
			if err != nil {
				return err
			}

			result = append(result, po)
		}

		return nil
	})

	return result, err

}

type UpdateCatalogCartParams struct {
	models.JwtClaimsInfo

	CartID string                      `json:"cart_id" param:"cart_id" validate:"required"`
	Items  []*models.PurchaseOrderItem `json:"items"`
}

func (r *CatalogCartRepo) UpdateCatalogCart(tx *gorm.DB, params UpdateCatalogCartParams) (*models.PurchaseOrder, error) {
	var purchaseOrder models.PurchaseOrder
	var err = tx.First(&purchaseOrder, "id = ? AND is_cart = ?", params.CartID, true).Error
	if err != nil {
		if r.db.IsRecordNotFoundError(err) {
			return nil, errs.ErrPONotFound
		}
		return nil, err
	}

	if params.Items != nil && len(params.Items) == 0 {
		err = tx.Unscoped().Delete(&models.PurchaseOrder{}, "id = ?", purchaseOrder.ID).Error
		if err != nil {
			return nil, err
		}

		err = tx.Unscoped().Delete(&models.PurchaseOrderItem{}, "purchase_order_id = ?", purchaseOrder.ID).Error
		if err != nil {
			return nil, err
		}

		return &purchaseOrder, nil
	}

	var itemIDs []string
	var subTotalPrice = price.NewFromFloat(0)

	var items = lo.Map(params.Items, func(item *models.PurchaseOrderItem, index int) *models.PurchaseOrderItem {
		if item.ID == "" {
			item.ID = helper.GenerateXID()
		}

		item.PurchaseOrderID = purchaseOrder.ID
		item.UserID = purchaseOrder.UserID
		item.TotalPrice = item.UnitPrice.MultipleInt(item.Quantity).ToPtr()
		itemIDs = append(itemIDs, item.ID)
		subTotalPrice = subTotalPrice.AddPtr(item.TotalPrice)
		return item
	})

	err = tx.Unscoped().Delete(&models.PurchaseOrderItem{}, "purchase_order_id = ? AND id NOT IN ?", purchaseOrder.ID, itemIDs).Error
	if err != nil {
		return nil, err
	}

	purchaseOrder.SubTotal = subTotalPrice.ToPtr()
	err = purchaseOrder.UpdatePrices()
	if err != nil {
		return nil, err
	}

	err = tx.Omit(clause.Associations).Save(&items).Error
	if err != nil {
		return nil, err
	}

	err = tx.Model(&models.PurchaseOrder{}).Where("id = ?", params.CartID).Updates(&purchaseOrder).Error
	if err != nil {
		return nil, err
	}

	purchaseOrder.Items = items

	return &purchaseOrder, err
}

type CreateCatalogCartForm struct {
	models.JwtClaimsInfo

	ProductID string `param:"product_id" validate:"required"`

	Items []*models.PurchaseOrderItem `json:"items"`
}

func (r *CatalogCartRepo) CreateCatalogCart(params CreateCatalogCartForm) (*models.PurchaseOrder, error) {
	var product models.Product
	var err = r.db.Select("ID", "Currency").First(&product, "id = ?", params.ProductID).Error
	if err != nil {
		return nil, err
	}
	var user models.User
	if err := r.db.Select("ContactOwnerIDs").First(&user, "id = ?", params.GetUserID()).Error; err != nil {
		return nil, err
	}

	var itemIDs []string
	var subTotalPrice = price.NewFromFloat(0)

	var purchaseOrder = models.PurchaseOrder{
		UserID:      params.GetUserID(),
		Currency:    product.Currency,
		IsCart:      values.Bool(true),
		FromCatalog: true,
		AssigneeIDs: user.ContactOwnerIDs,
	}
	purchaseOrder.ID = helper.GenerateXID()

	var items = lo.Map(params.Items, func(item *models.PurchaseOrderItem, index int) *models.PurchaseOrderItem {
		item.ID = helper.GenerateXID()
		item.PurchaseOrderID = purchaseOrder.ID
		item.UserID = purchaseOrder.UserID
		item.ProductID = params.ProductID
		item.TotalPrice = item.UnitPrice.MultipleInt(item.Quantity).ToPtr()
		itemIDs = append(itemIDs, item.ID)
		subTotalPrice = subTotalPrice.AddPtr(item.TotalPrice)
		return item
	})

	err = r.db.Transaction(func(tx *gorm.DB) error {
		err = r.db.Create(&items).Error
		if err != nil {
			return err
		}

		purchaseOrder.SubTotal = subTotalPrice.ToPtr()
		err = purchaseOrder.UpdatePrices()
		if err != nil {
			return err
		}

		err = r.db.Omit(clause.Associations).Clauses(clause.OnConflict{UpdateAll: true}).Create(&purchaseOrder).Error
		return err
	})
	if err != nil {
		return nil, err
	}

	return &purchaseOrder, err
}

type CreateCatalogCartOrderParams struct {
	models.JwtClaimsInfo

	CartID  string   `json:"cart_id" param:"cart_id" validate:"required"`
	ItemIDs []string `json:"item_ids" validate:"required"`
}

func (r *CatalogCartRepo) CreateCatalogCartOrder(tx *gorm.DB, params CreateCatalogCartOrderParams) (*models.PurchaseOrder, error) {
	var purchaseOrder models.PurchaseOrder
	var err = tx.First(&purchaseOrder, "id = ?", params.CartID).Error
	if err != nil {
		return nil, err
	}

	var currentItems []*models.PurchaseOrderItem
	err = tx.Find(&currentItems, "purchase_order_id = ?", purchaseOrder.ID).Error
	if err != nil {
		return nil, err
	}

	var currentItemIDs = lo.Map(currentItems, func(item *models.PurchaseOrderItem, index int) string {
		return item.ID
	})

	_, removedItemIDs := lo.Difference(currentItemIDs, params.ItemIDs)
	if len(removedItemIDs) > 0 {
		var clonePurchaseOrder models.PurchaseOrder
		err = copier.Copy(&clonePurchaseOrder, &purchaseOrder)
		if err != nil {
			return nil, err
		}
		clonePurchaseOrder.ID = helper.GenerateXID()
		clonePurchaseOrder.ReferenceID = helper.GeneratePurchaseOrderReferenceID()
		var items = lo.Filter(currentItems, func(item *models.PurchaseOrderItem, index int) bool {
			return lo.Contains(removedItemIDs, item.ID)
		})

		items = lo.Map(items, func(item *models.PurchaseOrderItem, index int) *models.PurchaseOrderItem {
			item.PurchaseOrderID = clonePurchaseOrder.ID
			return item
		})

		err = tx.Create(&items).Error
		if err != nil {
			return nil, err
		}

		err = tx.Omit(clause.Associations).Create(&clonePurchaseOrder).Error
		if err != nil {
			return nil, err
		}

		return &clonePurchaseOrder, nil
	}

	var updates = models.PurchaseOrder{
		IsCart: values.Bool(false),
	}
	err = tx.Model(&models.PurchaseOrder{}).Where("id = ?", params.CartID).Updates(&updates).Error
	if err != nil {
		return nil, err
	}

	purchaseOrder.IsCart = updates.IsCart
	return &purchaseOrder, err
}

type CreateCatalogCartOrdersParams struct {
	models.JwtClaimsInfo

	Records []CreateCatalogCartOrderParams
}

func (r *CatalogCartRepo) CreateCatalogCartOrders(params CreateCatalogCartOrdersParams) ([]*models.PurchaseOrder, error) {
	var result []*models.PurchaseOrder

	var err = r.db.Transaction(func(tx *gorm.DB) error {
		for _, v := range params.Records {
			po, err := r.CreateCatalogCartOrder(tx, v)
			if err != nil {
				return err
			}

			result = append(result, po)
		}

		return nil
	})

	return result, err
}

type CatalogCartCheckoutInfoParams struct {
	models.JwtClaimsInfo

	CartID      string            `json:"cart_id" validate:"required"`
	PaymentType enums.PaymentType `json:"payment_type" validate:"required"`
}

func (r *CatalogCartRepo) CatalogCartCheckoutInfo(params CatalogCartCheckoutInfoParams) (*models.PurchaseOrder, error) {
	cart, err := r.GetCatalogCart(GetCatalogCartParams{
		JwtClaimsInfo: params.JwtClaimsInfo,
		CartID:        params.CartID,
	})
	if err != nil {
		return nil, err
	}

	cart.PaymentType = params.PaymentType
	err = cart.UpdatePrices()
	if err != nil {
		return nil, err
	}

	err = r.db.Omit(clause.Associations).Save(&cart).Error

	return cart, err
}

type MultiCatalogCartCheckoutInfoParams struct {
	models.JwtClaimsInfo

	Records []CatalogCartCheckoutInfoParams `json:"records" validate:"required"`
}

func (r *CatalogCartRepo) MultiCatalogCartCheckoutInfo(params MultiCatalogCartCheckoutInfoParams) ([]*models.PurchaseOrder, error) {
	var result []*models.PurchaseOrder
	for _, record := range params.Records {
		po, err := r.CatalogCartCheckoutInfo(record)
		if err != nil {
			return nil, err
		}

		result = append(result, po)
	}
	return result, nil
}

type MultiCatalogCartCheckoutParams struct {
	models.JwtClaimsInfo

	CartIDs         []string          `json:"cart_ids" params:"cart_ids" validate:"required"`
	PaymentType     enums.PaymentType `json:"payment_type" validate:"oneof=bank_transfer card"`
	PaymentMethodID string            `json:"payment_method_id" validate:"required_if=PaymentType card"`

	TransactionRefID      string             `json:"transaction_ref_id" validate:"required_if=PaymentType bank_transfer"`
	TransactionAttachment *models.Attachment `json:"transaction_attachment" validate:"required_if=PaymentType bank_transfer"`
}

type MultiCatalogCartCheckoutResponse struct {
	Orders                    []*models.PurchaseOrder         `json:"orders"`
	CheckoutSessionID         string                          `json:"checkout_session_id"`
	PaymentTransaction        *models.PaymentTransaction      `json:"payment_transaction"`
	PaymentIntentNextAction   *stripe.PaymentIntentNextAction `json:"payment_intent_next_action,omitempty"`
	PaymentIntentClientSecret string                          `json:"payment_intent_client_secret,omitempty"`
}

func (r *CatalogCartRepo) MultiCatalogCartCheckout(params MultiCatalogCartCheckoutParams) (*MultiCatalogCartCheckoutResponse, error) {
	var checkoutSessionID = helper.GenerateCheckoutSessionID()

	var resp = MultiCatalogCartCheckoutResponse{
		CheckoutSessionID: checkoutSessionID,
	}
	var orders []*models.PurchaseOrder
	var err = r.db.Find(&orders, "id IN ?", params.CartIDs).Error
	if err != nil {
		return nil, err
	}

	if len(orders) <= 0 {
		return nil, errs.ErrInquiryCartInvalidToCheckout
	}

	if params.PaymentType == enums.PaymentTypeBankTransfer {
		err = r.db.Transaction(func(tx *gorm.DB) error {
			var totalAmount price.Price
			var currency enums.Currency
			var purchaseOrderIDs []string
			var purchaseOrderReferenceIDs []string
			var inquiryIDs []string
			var inquiryReferenceIDs []string
			var transactionRefID = helper.GeneratePaymentTransactionReferenceID()
			for _, purchaseOrder := range orders {
				currency = purchaseOrder.Currency
				totalAmount = totalAmount.AddPtr(purchaseOrder.TotalPrice)

				var updates = models.PurchaseOrder{
					Status:                        enums.PurchaseOrderStatusWaitingConfirm,
					PaymentType:                   params.PaymentType,
					TransactionRefID:              params.TransactionRefID,
					TransactionAttachment:         params.TransactionAttachment,
					TransferedAt:                  values.Int64(time.Now().Unix()),
					Currency:                      purchaseOrder.Currency,
					CheckoutSessionID:             checkoutSessionID,
					PaymentTransactionReferenceID: transactionRefID,
				}

				updates.TaxPercentage = purchaseOrder.TaxPercentage
				if len(purchaseOrder.Quotations) > 0 {
					sampleQuotation, _ := lo.Find(purchaseOrder.Quotations, func(item *models.InquiryQuotationItem) bool {
						return item.Type == enums.InquiryTypeSample
					})
					if sampleQuotation != nil {
						updates.LeadTime = int(values.Int64Value(sampleQuotation.LeadTime))
						updates.StartDate = updates.MarkAsPaidAt
						updates.CompletionDate = values.Int64(time.Unix(*updates.StartDate, 0).AddDate(0, 0, updates.LeadTime).Unix())
					}
				}
				var sqlResult = tx.Model(&models.PurchaseOrder{}).Where("id = ?", purchaseOrder.ID).Updates(&updates)
				if sqlResult.Error != nil {
					return eris.Wrap(sqlResult.Error, sqlResult.Error.Error())
				}

				if sqlResult.RowsAffected == 0 {
					return eris.New("Purchase order not found")
				}

				purchaseOrderIDs = append(purchaseOrderIDs, purchaseOrder.ID)
				purchaseOrderReferenceIDs = append(purchaseOrderReferenceIDs, purchaseOrder.ReferenceID)

				if purchaseOrder.Inquiry != nil {
					inquiryIDs = append(inquiryIDs, purchaseOrder.Inquiry.ID)
					inquiryReferenceIDs = append(inquiryReferenceIDs, purchaseOrder.Inquiry.ReferenceID)
				}

			}

			// create transaction
			var transaction = models.PaymentTransaction{
				ReferenceID:       transactionRefID,
				PaidAmount:        totalAmount.ToPtr(),
				PaymentType:       params.PaymentType,
				Milestone:         enums.PaymentMilestoneFinalPayment,
				UserID:            params.GetUserID(),
				TransactionRefID:  params.TransactionRefID,
				Status:            enums.PaymentStatusWaitingConfirm,
				PaymentPercentage: values.Float64(100),
				TotalAmount:       totalAmount.ToPtr(),
				Currency:          currency,
				CheckoutSessionID: checkoutSessionID,
				PurchaseOrderIDs:  purchaseOrderIDs,
				Metadata: &models.PaymentTransactionMetadata{
					InquiryIDs:                inquiryIDs,
					InquiryReferenceIDs:       inquiryReferenceIDs,
					PurchaseOrderReferenceIDs: purchaseOrderReferenceIDs,
					PurchaseOrderIDs:          purchaseOrderIDs,
				},
			}
			if params.TransactionAttachment != nil {
				transaction.Attachments = &models.Attachments{params.TransactionAttachment}
			}
			err = tx.Create(&transaction).Error
			if err != nil {
				return eris.Wrap(err, "Create transaction error")
			}
			resp.PaymentTransaction = &transaction
			return err

		})

		if err != nil {
			return nil, err
		}
	} else if params.PaymentType == enums.PaymentTypeCard {
		var user models.User
		err = r.db.Select("ID", "StripeCustomerID").First(&user, "id = ?", params.GetUserID()).Error
		if err != nil {
			return nil, err
		}
		var totalAmount price.Price
		var currency enums.Currency
		var finalCardItemIDs []string

		for _, purchaseOrder := range orders {
			currency = purchaseOrder.Currency
			totalAmount = totalAmount.AddPtr(purchaseOrder.TotalPrice)
			finalCardItemIDs = append(finalCardItemIDs, purchaseOrder.CartItemIDs...)
		}

		stripeConfig, err := stripehelper.GetCurrencyConfig(currency)
		if err != nil {
			return nil, err
		}
		var stripeParams = stripehelper.CreatePaymentIntentParams{
			Amount:                  totalAmount.MultipleInt(stripeConfig.SmallestUnitFactor).ToInt64(),
			Currency:                currency,
			PaymentMethodID:         params.PaymentMethodID,
			CustomerID:              user.StripeCustomerID,
			IsCaptureMethodManually: false,
			Description:             fmt.Sprintf("Charges for checkout session %s", checkoutSessionID),
			PaymentMethodTypes:      []string{"card"},
			Metadata: map[string]string{
				"cart_item_ids":       strings.Join(finalCardItemIDs, ","),
				"checkout_session_id": checkoutSessionID,
				"action_source":       string(stripehelper.ActionSourceMultiInquiryPayment),
			},
		}

		pi, err := stripehelper.GetInstance().CreatePaymentIntent(stripeParams)
		if err != nil {
			return nil, err
		}
		if pi.Status != stripe.PaymentIntentStatusSucceeded {
			if pi.NextAction != nil {
				intent, err := stripehelper.GetInstance().ConfirmPaymentIntent(stripehelper.ConfirmPaymentIntentParams{
					PaymentIntentID: pi.ID,
					ReturnURL:       fmt.Sprintf("%s/api/v1/callback/stripe/payment_intents/catalog_carts/%s/confirm?cart_ids=%s", r.db.Configuration.ServerBaseURL, checkoutSessionID, strings.Join(params.CartIDs, ",")),
				})
				if err != nil {
					return nil, err
				}

				if intent.Status == stripe.PaymentIntentStatusSucceeded {
					goto PaymentSuccess
				}

				resp.PaymentIntentNextAction = intent.NextAction
				resp.PaymentIntentClientSecret = intent.ClientSecret
				return &resp, nil

			} else {
				return nil, eris.Errorf("Payment error with status %s", pi.Status)
			}
		}

	PaymentSuccess:
		err = r.db.Transaction(func(tx *gorm.DB) error {
			var purchaseOrderIDs []string
			var purchaseOrderReferenceIDs []string
			var inquiryIDs []string
			var inquiryReferenceIDs []string
			var transactionRefID = helper.GeneratePaymentTransactionReferenceID()

			for _, purchaseOrder := range orders {
				var updates = models.PurchaseOrder{
					PaymentIntentID:               pi.ID,
					Status:                        enums.PurchaseOrderStatusPaid,
					PaymentType:                   params.PaymentType,
					MarkAsPaidAt:                  values.Int64(time.Now().Unix()),
					TransactionRefID:              params.TransactionRefID,
					TransactionAttachment:         params.TransactionAttachment,
					TransferedAt:                  values.Int64(time.Now().Unix()),
					Currency:                      purchaseOrder.Currency,
					CheckoutSessionID:             checkoutSessionID,
					PaymentTransactionReferenceID: transactionRefID,
				}

				if len(purchaseOrder.Quotations) > 0 {
					sampleQuotation, _ := lo.Find(purchaseOrder.Quotations, func(item *models.InquiryQuotationItem) bool {
						return item.Type == enums.InquiryTypeSample
					})
					if sampleQuotation != nil {
						updates.LeadTime = int(values.Int64Value(sampleQuotation.LeadTime))
						updates.StartDate = updates.TransferedAt
						updates.CompletionDate = values.Int64(time.Unix(*updates.StartDate, 0).AddDate(0, 0, updates.LeadTime).Unix())
					}
				}

				updates.TaxPercentage = purchaseOrder.TaxPercentage
				var sqlResult = tx.Model(&models.PurchaseOrder{}).Where("id = ?", purchaseOrder.ID).Updates(&updates)
				if sqlResult.Error != nil {
					return sqlResult.Error
				}

				if sqlResult.RowsAffected == 0 {
					return eris.New("Purchase order is not found")
				}

				var sqlResult2 = tx.Model(&models.Inquiry{}).Where("id = ?", purchaseOrder.InquiryID).UpdateColumn("Status", enums.InquiryStatusFinished)
				if sqlResult2.Error != nil {
					return sqlResult2.Error
				}

				purchaseOrderIDs = append(purchaseOrderIDs, purchaseOrder.ID)
				purchaseOrderReferenceIDs = append(purchaseOrderReferenceIDs, purchaseOrder.ReferenceID)

				if purchaseOrder.Inquiry != nil {
					inquiryIDs = append(inquiryIDs, purchaseOrder.Inquiry.ID)
					inquiryReferenceIDs = append(inquiryReferenceIDs, purchaseOrder.Inquiry.ReferenceID)
				}
			}

			// create transaction
			var transaction = models.PaymentTransaction{
				ReferenceID:       transactionRefID,
				PaidAmount:        totalAmount.ToPtr(),
				PaymentType:       params.PaymentType,
				UserID:            params.GetUserID(),
				TransactionRefID:  params.TransactionRefID,
				PaymentIntentID:   pi.ID,
				Status:            enums.PaymentStatusPaid,
				TotalAmount:       totalAmount.ToPtr(),
				Milestone:         enums.PaymentMilestoneFinalPayment,
				PaymentPercentage: values.Float64(100),
				MarkAsPaidAt:      values.Int64(time.Now().Unix()),
				Currency:          currency,
				CheckoutSessionID: checkoutSessionID,
				PurchaseOrderIDs:  purchaseOrderIDs,
				Metadata: &models.PaymentTransactionMetadata{
					InquiryIDs:                inquiryIDs,
					InquiryReferenceIDs:       inquiryReferenceIDs,
					PurchaseOrderReferenceIDs: purchaseOrderReferenceIDs,
					PurchaseOrderIDs:          purchaseOrderIDs,
				},
			}
			if params.TransactionAttachment != nil {
				transaction.Attachments = &models.Attachments{params.TransactionAttachment}
			}

			err = tx.Create(&transaction).Error
			resp.PaymentTransaction = &transaction
			return err
		})

		if err != nil {
			return nil, eris.Wrap(err, err.Error())
		}

	}

	resp.Orders = orders
	resp.CheckoutSessionID = checkoutSessionID

	return &resp, nil
}
