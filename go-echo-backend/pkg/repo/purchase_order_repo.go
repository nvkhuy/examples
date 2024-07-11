package repo

import (
	"bytes"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"

	"github.com/engineeringinflow/inflow-backend/pkg/helper"
	"github.com/engineeringinflow/inflow-backend/pkg/s3"
	"github.com/engineeringinflow/inflow-backend/pkg/stripehelper"
	"github.com/lib/pq"
	"github.com/samber/lo"
	"github.com/stripe/stripe-go/v74"
	stripePrice "github.com/stripe/stripe-go/v74/price"

	"github.com/engineeringinflow/inflow-backend/pkg/db"
	"github.com/engineeringinflow/inflow-backend/pkg/errs"
	"github.com/engineeringinflow/inflow-backend/pkg/models"
	"github.com/engineeringinflow/inflow-backend/pkg/models/enums"
	"github.com/engineeringinflow/inflow-backend/pkg/models/price"
	"github.com/engineeringinflow/inflow-backend/pkg/repo/query"
	"github.com/engineeringinflow/inflow-backend/pkg/repo/query/queryfunc"
	"github.com/jinzhu/copier"
	"github.com/rotisserie/eris"
	"github.com/thaitanloi365/go-utils/values"
	"golang.org/x/exp/slices"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type PurchaseOrderRepo struct {
	db *db.DB
}

func NewPurchaseOrderRepo(db *db.DB) *PurchaseOrderRepo {
	return &PurchaseOrderRepo{
		db: db,
	}
}

type PaginatePurchaseOrdersParams struct {
	models.PaginationParams
	models.JwtClaimsInfo

	UserID string `json:"user_id" query:"user_id"`
	TeamID string `json:"team_id" query:"team_id"`

	Statuses []enums.PurchaseOrderStatus `json:"statuses" query:"statuses"`

	ExcludedInquiryStatuses []enums.InquiryStatus `json:"excluded_inquiry_statuses" query:"excluded_inquiry_statuses"`

	InquiryStatuses []enums.InquiryBuyerStatus `json:"inquiry_statuses" query:"inquiry_statuses"`

	TrackingStatuses       []enums.PoTrackingStatus       `json:"tracking_statuses" query:"tracking_statuses"`
	SellerTrackingStatuses []enums.SellerPoTrackingStatus `json:"seller_tracking_statuses" query:"seller_tracking_statuses"`

	AssigneeIDs []string `json:"assignee_ids" query:"assignee_ids"`

	AssigneeID    string `json:"assignee_id" query:"assignee_id"`
	SampleMakerID string `json:"sample_maker_id" param:"sample_maker_id" query:"sample_maker_id"`

	RoundID       string              `json:"round_id" query:"round_id"`
	RoundStatuses []enums.RoundStatus `json:"round_statuses" query:"round_statuses"`

	PostedDateFrom int64 `json:"posted_date_from" query:"posted_date_from"`
	PostedDateTo   int64 `json:"posted_date_to" query:"posted_date_to"`

	CatalogTrackingStatuses []enums.PoCatalogTrackingStatus `json:"catalog_tracking_statuses" query:"catalog_tracking_statuses"`

	CatalogSamples              *bool `json:"catalog_samples" query:"catalog_samples"`
	IncludeAssignee             bool  `json:"-"`
	IncludeSampleMaker          bool  `json:"-"`
	IncludeTrackings            bool  `json:"-"`
	IncludeUnpaidWithoutInquiry bool  `json:"-"`
	IncludeItems                bool  `json:"-"`
	IsQueryAll                  bool  `json:"-"`
	IncludeCollection           bool  `json:"-"`
}

func (r *PurchaseOrderRepo) PaginatePurchaseOrders(params PaginatePurchaseOrdersParams) *query.Pagination {
	var userID = params.GetUserID()
	if params.TeamID != "" && !params.GetRole().IsAdmin() && !params.IsQueryAll {
		if err := r.db.Select("ID").First(&models.BrandTeam{}, "team_id = ? AND user_id = ?", params.TeamID, userID).Error; err != nil {
			return &query.Pagination{
				Records: []*models.PurchaseOrder{},
			}
		}
		userID = params.TeamID
	}

	var result = query.New(r.db, queryfunc.NewPurchaseOrderBuilder(queryfunc.PurchaseOrderBuilderOptions{
		QueryBuilderOptions: queryfunc.QueryBuilderOptions{
			Role: params.GetRole(),
		},
		IncludeItems:       params.IncludeItems,
		IncludeUsers:       params.GetRole().IsAdmin(),
		IncludeAssignee:    params.IncludeAssignee,
		IncludeSampleMaker: params.IncludeSampleMaker,
		IncludeTrackings:   params.IncludeTrackings,
		IncludeCollection:  params.IncludeCollection,
	})).
		WhereFunc(func(builder *query.Builder) {
			if params.IncludeUnpaidWithoutInquiry {
				builder.Where("(po.status = ? OR (po.status = ? AND COALESCE(po.inquiry_id,'') = ?))", enums.PurchaseOrderStatusPaid, enums.PurchaseOrderStatusPending, "")
			} else {
				params.Statuses = append(params.Statuses, enums.PurchaseOrderStatusPaid)
			}

			builder.Where("po.is_cart = ?", false)

			if !params.IsQueryAll {
				if params.GetRole().IsAdmin() {
					if params.UserID != "" {
						builder.Where("po.user_id = ?", params.UserID)
					}
				} else if params.GetRole().IsSeller() {
					// do nothing
				} else {
					builder.Where("po.user_id = ? AND po.deleted_at IS NULL", userID)
				}
			}

			if params.SampleMakerID != "" {
				builder.Where("po.sample_maker_id = ?", params.SampleMakerID)
			}

			if params.RoundID != "" {
				builder.Where("po.round_id = ?", params.RoundID)
			}

			if len(params.RoundStatuses) > 0 {
				builder.Where("po.round_status IN ?", params.RoundStatuses)
			}

			if len(params.Statuses) > 0 {
				builder.Where("po.status IN ?", params.Statuses)
			}

			if len(params.TrackingStatuses) > 0 {
				builder.Where("po.tracking_status IN ?", params.TrackingStatuses)
			}

			if len(params.SellerTrackingStatuses) > 0 {
				builder.Where("po.seller_tracking_status IN ?", params.SellerTrackingStatuses)
			}

			if params.PostedDateFrom > 0 {
				builder.Where("po.created_at >= ?", params.PostedDateFrom)
			}

			if params.PostedDateTo > 0 {
				builder.Where("po.created_at <= ?", params.PostedDateTo)
			}

			if len(params.InquiryStatuses) > 0 {
				builder.Where("iq.status IN ?", params.InquiryStatuses)
			}

			if params.AssigneeID != "" {
				builder.Where("count_elements(po.assignee_ids,?) >= 1", pq.StringArray([]string{params.AssigneeID}))
			}

			if len(params.AssigneeIDs) > 0 {
				builder.Where("count_elements(po.assignee_ids,?) >= 1", pq.StringArray(params.AssigneeIDs))
			}

			if params.CatalogSamples != nil {
				builder.Where("po.from_catalog = ?", *params.CatalogSamples)
			}

			if keyword := strings.TrimSpace(params.Keyword); keyword != "" {
				var q = "%" + keyword + "%"
				if strings.HasPrefix(keyword, "PO-") {
					builder.Where("po.reference_id ILIKE ?", q)
				} else if strings.HasPrefix(keyword, "IQ-") {
					builder.Where("iq.reference_id ILIKE ?", q)
				} else {
					builder.Where("(po.id ILIKE @keyword OR po.client_reference_id ILIKE @keyword OR iq.id ILIKE @keyword)", sql.Named("keyword", q))
				}

			}
		}).
		OrderBy("po.order_group_id ASC, po.updated_at DESC").
		Limit(params.Limit).
		Page(params.Page).
		PagingFunc()
	return result
}

type GetPurchaseOrderParams struct {
	models.JwtClaimsInfo

	PurchaseOrderID           string `json:"purchase_order_id" param:"purchase_order_id" query:"purchase_order_id" validate:"required"`
	TeamID                    string `json:"team_id" param:"team_id" query:"team_id"`
	SampleMakerID             string `json:"sample_maker_id" param:"sample_maker_id" query:"sample_maker_id"`
	UserID                    string `json:"-"`
	IncludeAssignee           bool   `json:"-"`
	IncludeSampleMaker        bool   `json:"-"`
	IncludeInquirySeller      bool   `json:"-"`
	IncludeInvoice            bool   `json:"-"`
	IncludeUsers              bool   `json:"-"`
	IncludePaymentTransaction bool   `json:"-"`
	IsQueryAll                bool   `json:"-"`
}

func (r *PurchaseOrderRepo) GetPurchaseOrder(params GetPurchaseOrderParams) (*models.PurchaseOrder, error) {
	if params.TeamID != "" && !params.GetRole().IsAdmin() && !params.IsQueryAll {
		if err := r.db.Select("ID").First(&models.BrandTeam{}, "team_id = ? AND user_id = ?", params.TeamID, params.GetUserID()); err != nil {
			return nil, errs.ErrPONotFound
		}
		params.UserID = params.GetUserID()
	}

	var purchaseOrder models.PurchaseOrder
	var err = query.New(r.db, queryfunc.NewPurchaseOrderBuilder(queryfunc.PurchaseOrderBuilderOptions{
		QueryBuilderOptions: queryfunc.QueryBuilderOptions{
			Role: params.GetRole(),
		},
		IncludeCartItems:          true,
		IncludeItems:              true,
		IncludeUsers:              params.IncludeUsers,
		IncludeAssignee:           params.IncludeAssignee,
		IncludeSampleMaker:        params.IncludeSampleMaker,
		IncludeInquirySeller:      params.IncludeInquirySeller,
		IncludeInvoice:            params.IncludeInvoice,
		IncludePaymentTransaction: params.IncludePaymentTransaction,
	})).
		WhereFunc(func(builder *query.Builder) {
			if strings.HasPrefix(params.PurchaseOrderID, "PO-") {
				builder.Where("po.reference_id = ?", params.PurchaseOrderID)
			} else {
				builder.Where("po.id = ?", params.PurchaseOrderID)
			}

			if params.UserID != "" {
				builder.Where("po.user_id = ?", params.UserID)
			}
			if params.SampleMakerID != "" {
				builder.Where("po.sample_maker_id = ?", params.SampleMakerID)
			}

			if params.GetRole().IsSeller() && params.UserID != "" {
				builder.Where("po.sample_maker_id = ?", params.UserID)
			} else if params.GetRole().IsBuyer() && params.UserID != "" {
				builder.Where("po.user_id = ?", params.UserID)
			}
		}).
		Limit(1).
		FirstFunc(&purchaseOrder)

	return &purchaseOrder, err
}

func (r *PurchaseOrderRepo) GetPurchaseOrderShortInfo(purchaseOrderID string) (*models.PurchaseOrder, error) {
	var purchaseOrder models.PurchaseOrder
	var err = r.db.First(&purchaseOrder, "id = ?", purchaseOrderID).Error

	return &purchaseOrder, err
}

func (r *PurchaseOrderRepo) GetPurchaseOrderWithSeller(purchaseOrderID string) (*models.PurchaseOrder, error) {
	var purchaseOrder models.PurchaseOrder
	var err = r.db.First(&purchaseOrder, "id = ?", purchaseOrderID).Error
	if err != nil {
		return nil, err
	}

	if purchaseOrder.SampleMakerID != "" {
		var seller models.User
		err = r.db.First(&seller, "id = ?", purchaseOrder.SampleMakerID).Error
		if err != nil {
			return nil, err
		}

		purchaseOrder.SampleMaker = &seller
	}
	return &purchaseOrder, err
}

type UpdatePurchaseOrderTrackingStatusParams struct {
	models.JwtClaimsInfo

	PurchaseOrderID string                 `json:"purchase_order_id" param:"purchase_order_id" query:"purchase_order_id" validate:"required"`
	TrackingStatus  enums.PoTrackingStatus `json:"tracking_status" param:"tracking_status" query:"tracking_status"`
	TrackingAction  enums.PoTrackingAction `json:"tracking_action" param:"tracking_action" query:"tracking_action"`

	UserID string `json:"-"`
}

type AdminUpdatePurchaseOrderParams struct {
	models.JwtClaimsInfo

	PurchaseOrderID     string              `json:"purchase_order_id" param:"purchase_order_id" query:"purchase_order_id" validate:"required"`
	TechpackAttachments *models.Attachments `json:"techpack_attachments" param:"techpack_attachments" query:"techpack_attachments"`
	SampleAttachments   *models.Attachments `json:"sample_attachments" param:"sample_attachments" query:"sample_attachments"`

	UserID string `json:"-"`
}

func (r *PurchaseOrderRepo) UpdatePurchaseOrderTrackingStatus(params UpdatePurchaseOrderTrackingStatusParams) (*models.PurchaseOrder, error) {
	order, err := r.GetPurchaseOrder(GetPurchaseOrderParams{
		PurchaseOrderID: params.PurchaseOrderID,
		JwtClaimsInfo:   params.JwtClaimsInfo,
	})

	if err != nil {
		return nil, err
	}

	var updates models.PurchaseOrder
	err = copier.Copy(&updates, &params)
	if err != nil {
		return nil, err
	}

	if order.TrackingStatus != updates.TrackingStatus {
		var trackingParams = models.PurchaseOrderTrackingCreateForm{
			PurchaseOrderID: params.PurchaseOrderID,
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
		}

		if updates.TrackingStatus == enums.PoTrackingStatusMaking {
			trackingParams.Metadata.After = map[string]interface{}{
				"tracking_status": updates.TrackingStatus,
				"making_info":     updates.MakingInfo,
			}
		}

		if updates.TrackingStatus == enums.PoTrackingStatusSubmit {
			trackingParams.Metadata.After = map[string]interface{}{
				"tracking_status": updates.TrackingStatus,
				"submit_info":     updates.SubmitInfo,
			}
		}

		err = r.db.Transaction(func(tx *gorm.DB) error {
			err = NewPurchaseOrderTrackingRepo(r.db).CreatePurchaseOrderTrackingTx(tx, trackingParams)
			if err != nil {
				return eris.Wrap(err, err.Error())
			}

			return tx.Model(&models.PurchaseOrder{}).
				Where("id = ?", params.PurchaseOrderID).
				Updates(&updates).Error
		})
		if err != nil {
			return nil, eris.Wrap(err, err.Error())
		}
	}

	order.TrackingStatus = updates.TrackingStatus

	return order, err
}

type AdminPurchaseOrderMarkSubmitParams struct {
	models.JwtClaimsInfo

	PurchaseOrderID string                      `json:"purchase_order_id" param:"purchase_order_id" query:"purchase_order_id" validate:"required"`
	SubmitInfo      *models.PoMarkingStatusMeta `json:"submit_info" param:"submit_info" query:"submit_info"`
}

func (r *PurchaseOrderRepo) AdminPurchaseOrderMarkSubmit(params AdminPurchaseOrderMarkSubmitParams) (*models.PurchaseOrder, error) {
	order, err := r.GetPurchaseOrder(GetPurchaseOrderParams{
		PurchaseOrderID: params.PurchaseOrderID,
		JwtClaimsInfo:   params.JwtClaimsInfo,
	})

	if err != nil {
		return nil, err
	}

	validStatus := []enums.PoTrackingStatus{enums.PoTrackingStatusMaking}
	if ok := slices.Contains(validStatus, order.TrackingStatus); !ok {
		return nil, errs.ErrPoInvalidToMarkSubmit
	}

	var updates models.PurchaseOrder
	err = copier.Copy(&updates, &params)
	if err != nil {
		return nil, err
	}
	updates.TrackingStatus = enums.PoTrackingStatusSubmit

	if order.TrackingStatus != updates.TrackingStatus {
		err = r.db.Transaction(func(tx *gorm.DB) error {
			err = NewPurchaseOrderTrackingRepo(r.db).CreatePurchaseOrderTrackingTx(tx, models.PurchaseOrderTrackingCreateForm{
				PurchaseOrderID: params.PurchaseOrderID,
				ActionType:      enums.PoTrackingActionMarkSubmit,
				UserID:          order.UserID,
				CreatedByUserID: params.JwtClaimsInfo.GetUserID(),
				Metadata: &models.PoTrackingMetadata{
					Before: map[string]interface{}{
						"tracking_status": order.TrackingStatus,
					},
					After: map[string]interface{}{
						"tracking_status": updates.TrackingStatus,
						"submit_info":     updates.SubmitInfo,
					},
				},
			})
			if err != nil {
				return eris.Wrap(err, err.Error())
			}

			return tx.Model(&models.PurchaseOrder{}).
				Where("id = ?", params.PurchaseOrderID).
				Updates(&updates).Error
		})
		if err != nil {
			return nil, eris.Wrap(err, err.Error())
		}
	}

	order.SubmitInfo = updates.SubmitInfo
	order.TrackingStatus = updates.TrackingStatus

	return order, err
}

type AdminPurchaseOrderMarkMakingParams struct {
	models.JwtClaimsInfo

	PurchaseOrderID string                      `json:"purchase_order_id" param:"purchase_order_id" query:"purchase_order_id" validate:"required"`
	MakingInfo      *models.PoMarkingStatusMeta `json:"making_info" param:"making_info" query:"making_info"`
}

func (r *PurchaseOrderRepo) AdminPurchaseOrderMarkMaking(params AdminPurchaseOrderMarkMakingParams) (*models.PurchaseOrder, error) {
	order, err := r.GetPurchaseOrder(GetPurchaseOrderParams{
		PurchaseOrderID: params.PurchaseOrderID,
		JwtClaimsInfo:   params.JwtClaimsInfo,
	})

	if err != nil {
		return nil, err
	}

	validStatus := []enums.PoTrackingStatus{enums.PoTrackingStatusRawMaterial}
	if ok := slices.Contains(validStatus, order.TrackingStatus); !ok {
		return nil, errs.ErrPoInvalidToMarkMaking
	}

	var updates models.PurchaseOrder
	err = copier.Copy(&updates, &params)
	if err != nil {
		return nil, err
	}
	updates.TrackingStatus = enums.PoTrackingStatusMaking

	if order.TrackingStatus != updates.TrackingStatus {
		err = r.db.Transaction(func(tx *gorm.DB) error {
			err = NewPurchaseOrderTrackingRepo(r.db).CreatePurchaseOrderTrackingTx(tx, models.PurchaseOrderTrackingCreateForm{
				PurchaseOrderID: params.PurchaseOrderID,
				ActionType:      enums.PoTrackingActionMarkMaking,
				UserID:          order.UserID,
				CreatedByUserID: params.JwtClaimsInfo.GetUserID(),
				Metadata: &models.PoTrackingMetadata{
					Before: map[string]interface{}{
						"tracking_status": order.TrackingStatus,
					},
					After: map[string]interface{}{
						"tracking_status": updates.TrackingStatus,
						"making_info":     updates.MakingInfo,
					},
				},
			})
			if err != nil {
				return eris.Wrap(err, err.Error())
			}

			return tx.Model(&models.PurchaseOrder{}).
				Where("id = ?", params.PurchaseOrderID).
				Updates(&updates).Error
		})
		if err != nil {
			return nil, eris.Wrap(err, err.Error())
		}
	}

	order.MakingInfo = updates.MakingInfo
	order.TrackingStatus = updates.TrackingStatus

	return order, err
}

type AdminPurchaseOrderMarkMakingWithoutRawMaterialParams struct {
	models.JwtClaimsInfo

	PurchaseOrderID string                      `json:"purchase_order_id" param:"purchase_order_id" query:"purchase_order_id" validate:"required"`
	MakingInfo      *models.PoMarkingStatusMeta `json:"making_info" param:"making_info" query:"making_info"`
}

func (r *PurchaseOrderRepo) MarkMakingWithoutRawMaterial(params AdminPurchaseOrderMarkMakingWithoutRawMaterialParams) (*models.PurchaseOrder, error) {
	order, err := r.GetPurchaseOrder(GetPurchaseOrderParams{
		PurchaseOrderID: params.PurchaseOrderID,
		JwtClaimsInfo:   params.JwtClaimsInfo,
	})

	if err != nil {
		return nil, err
	}

	if order.PoRawMaterials != nil && len(*order.PoRawMaterials) > 0 {
		return nil, errs.ErrPoInvalidToSkipMaterial
	}

	var updates models.PurchaseOrder
	err = copier.Copy(&updates, &params)
	if err != nil {
		return nil, err
	}
	updates.TrackingStatus = enums.PoTrackingStatusMaking

	if order.TrackingStatus != updates.TrackingStatus {
		err = r.db.Transaction(func(tx *gorm.DB) error {
			err = NewPurchaseOrderTrackingRepo(r.db).CreatePurchaseOrderTrackingTx(tx, models.PurchaseOrderTrackingCreateForm{
				PurchaseOrderID: params.PurchaseOrderID,
				ActionType:      enums.PoTrackingActionMarkMaking,
				UserID:          order.UserID,
				CreatedByUserID: params.JwtClaimsInfo.GetUserID(),
				Metadata: &models.PoTrackingMetadata{
					Before: map[string]interface{}{
						"tracking_status": order.TrackingStatus,
					},
					After: map[string]interface{}{
						"tracking_status": updates.TrackingStatus,
						"making_info":     updates.MakingInfo,
					},
				},
			})
			if err != nil {
				return eris.Wrap(err, err.Error())
			}

			return tx.Model(&models.PurchaseOrder{}).
				Where("id = ?", params.PurchaseOrderID).
				Updates(&updates).Error
		})

		if err != nil {
			return nil, eris.Wrap(err, err.Error())
		}
	}

	order.MakingInfo = updates.MakingInfo
	order.TrackingStatus = updates.TrackingStatus

	return order, err
}

type AdminPurchaseOrderMarkDeliveringParams struct {
	models.JwtClaimsInfo

	PurchaseOrderID string                 `json:"purchase_order_id" param:"purchase_order_id" query:"purchase_order_id" validate:"required"`
	LogisticInfo    *models.PoLogisticMeta `json:"logistic_info" param:"logistic_info" query:"logistic_info" form:"logistic_info" validate:"required"`
}

func (r *PurchaseOrderRepo) AdminPurchaseOrderMarkDelivering(params AdminPurchaseOrderMarkDeliveringParams) (*models.PurchaseOrder, error) {
	order, err := r.GetPurchaseOrder(GetPurchaseOrderParams{
		PurchaseOrderID: params.PurchaseOrderID,
		JwtClaimsInfo:   params.JwtClaimsInfo,
	})

	if err != nil {
		return nil, err
	}

	var validStatus = []enums.PoTrackingStatus{enums.PoTrackingStatusSubmit}
	var canDelivering = slices.Contains(validStatus, order.TrackingStatus) || (order.Status == enums.PurchaseOrderStatusPaid && order.InquiryID == "")
	if !canDelivering {
		return nil, errs.ErrPoInvalidToMarkDelivering
	}

	var updates models.PurchaseOrder
	err = copier.Copy(&updates, &params)
	if err != nil {
		return nil, err
	}
	updates.DeliveryStartedAt = values.Int64(time.Now().Unix())
	updates.TrackingStatus = enums.PoTrackingStatusDelivering

	if order.TrackingStatus != updates.TrackingStatus {
		r.db.Transaction(func(tx *gorm.DB) error {
			err = NewPurchaseOrderTrackingRepo(r.db).CreatePurchaseOrderTrackingTx(tx, models.PurchaseOrderTrackingCreateForm{
				PurchaseOrderID: params.PurchaseOrderID,
				ActionType:      enums.PoTrackingActionMarkDelivering,
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

			return tx.Model(&models.PurchaseOrder{}).
				Where("id = ?", params.PurchaseOrderID).
				Updates(&updates).Error
		})

		if err != nil {
			return nil, eris.Wrap(err, err.Error())
		}

	}

	order.DeliveryStartedAt = updates.DeliveryStartedAt
	order.TrackingStatus = updates.TrackingStatus
	order.LogisticInfo = updates.LogisticInfo

	return order, err
}

func (r *PurchaseOrderRepo) AdminPurchaseOrderConfirmDelivered(params UpdatePurchaseOrderTrackingStatusParams) (*models.PurchaseOrder, error) {
	order, err := r.GetPurchaseOrder(GetPurchaseOrderParams{
		PurchaseOrderID: params.PurchaseOrderID,
		JwtClaimsInfo:   params.JwtClaimsInfo,
	})
	if err != nil {
		return nil, err
	}
	validStatus := []enums.PoTrackingStatus{enums.PoTrackingStatusDelivering}
	if ok := slices.Contains(validStatus, order.TrackingStatus); !ok {
		return nil, errs.ErrPoInvalidToConfirmDelivered
	}

	var updates models.PurchaseOrder
	err = copier.Copy(&updates, &params)
	if err != nil {
		return nil, err
	}
	updates.ReceiverConfirmedAt = values.Int64(time.Now().Unix())
	updates.TrackingStatus = enums.PoTrackingStatusDeliveryConfirmed

	if order.TrackingStatus != updates.TrackingStatus {
		err = r.db.Transaction(func(tx *gorm.DB) error {
			err = NewPurchaseOrderTrackingRepo(r.db).CreatePurchaseOrderTrackingTx(tx, models.PurchaseOrderTrackingCreateForm{
				PurchaseOrderID: params.PurchaseOrderID,
				ActionType:      enums.PoTrackingActionConfirmDelivered,
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

			return tx.Model(&models.PurchaseOrder{}).
				Where("id = ?", params.PurchaseOrderID).
				Updates(&updates).Error
		})
		if err != nil {
			return nil, eris.Wrap(err, err.Error())
		}
	}

	order.ReceiverConfirmedAt = updates.ReceiverConfirmedAt
	order.TrackingStatus = updates.TrackingStatus

	return order, err
}

type AdminUpdatePurchaseOrderRawMaterialParams struct {
	models.JwtClaimsInfo

	PurchaseOrderID      string                    `json:"purchase_order_id" param:"purchase_order_id" query:"purchase_order_id" validate:"required"`
	PoRawMaterials       models.PoRawMaterialMetas `json:"po_raw_materials" param:"po_raw_materials" query:"po_raw_materials"`
	ApproveRawMaterialAt *int64                    `json:"approve_raw_material_at" param:"approve_raw_material_at" query:"approve_raw_material_at"`
}

func (r *PurchaseOrderRepo) AdminUpdatePurchaseOrderRawMaterial(params AdminUpdatePurchaseOrderRawMaterialParams) (*models.PurchaseOrder, error) {
	order, err := r.GetPurchaseOrder(GetPurchaseOrderParams{
		PurchaseOrderID: params.PurchaseOrderID,
		JwtClaimsInfo:   params.JwtClaimsInfo,
	})
	if err != nil {
		return nil, err
	}
	var updates models.PurchaseOrder
	err = copier.Copy(&updates, &params)
	if err != nil {
		return nil, err
	}

	_ = updates.GenerateRawMaterialRefID(updates.PoRawMaterials)

	updates.TrackingStatus = enums.PoTrackingStatusRawMaterial

	err = r.db.Transaction(func(tx *gorm.DB) error {
		err = NewPurchaseOrderTrackingRepo(r.db).CreatePurchaseOrderTrackingTx(tx, models.PurchaseOrderTrackingCreateForm{
			PurchaseOrderID: params.PurchaseOrderID,
			ActionType:      enums.PoTrackingActionUpdateMaterial,
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
			return eris.Wrap(err, err.Error())
		}

		return tx.Model(&models.PurchaseOrder{}).
			Where("id = ?", params.PurchaseOrderID).
			Updates(&updates).Error
	})
	if err != nil {
		return nil, eris.Wrap(err, err.Error())
	}

	order.TrackingStatus = updates.TrackingStatus
	return order, err
}

type AdminUpdatePurchaseOrderDesignParams struct {
	models.JwtClaimsInfo

	PurchaseOrderID     string              `json:"purchase_order_id" param:"purchase_order_id" query:"purchase_order_id" validate:"required"`
	TechpackAttachments *models.Attachments `json:"techpack_attachments" param:"techpack_attachments" query:"techpack_attachments"`
	ApproveDesignAt     *int64              `json:"approve_design_at" param:"approve_design_at" query:"approve_design_at"`
}

func (r *PurchaseOrderRepo) AdminUpdatePurchaseOrderDesign(params AdminUpdatePurchaseOrderDesignParams) (*models.PurchaseOrder, error) {
	order, err := r.GetPurchaseOrder(GetPurchaseOrderParams{
		PurchaseOrderID: params.PurchaseOrderID,
		JwtClaimsInfo:   params.JwtClaimsInfo,
	})
	if err != nil {
		return nil, err
	}

	validStatus := []enums.PoTrackingStatus{enums.PoTrackingStatusNew, enums.PoTrackingStatusDesignRejected, enums.PoTrackingStatusWaitingForApproved}
	if ok := slices.Contains(validStatus, order.TrackingStatus); !ok {
		return nil, errs.ErrPoInvalidToUploadDesign
	}

	var updates models.PurchaseOrder
	err = copier.Copy(&updates, &params)
	if err != nil {
		return nil, err
	}

	updates.TrackingStatus = enums.PoTrackingStatusWaitingForApproved

	err = r.db.Transaction(func(tx *gorm.DB) error {
		err = NewPurchaseOrderTrackingRepo(r.db).CreatePurchaseOrderTrackingTx(tx, models.PurchaseOrderTrackingCreateForm{
			PurchaseOrderID: params.PurchaseOrderID,
			ActionType:      enums.PoTrackingActionUpdateDesign,
			UserID:          order.UserID,
			CreatedByUserID: params.JwtClaimsInfo.GetUserID(),
			Metadata: &models.PoTrackingMetadata{
				Before: map[string]interface{}{
					"techpack_attachments": order.TechpackAttachments,
				},
				After: map[string]interface{}{
					"techpack_attachments": params.TechpackAttachments,
				},
			},
		})

		if err != nil {
			return eris.Wrap(err, err.Error())
		}

		return tx.Model(&models.PurchaseOrder{}).
			Where("id = ?", params.PurchaseOrderID).
			Updates(&updates).Error
	})

	if err != nil {
		return nil, eris.Wrap(err, err.Error())
	}

	order.TrackingStatus = enums.PoTrackingStatusWaitingForApproved

	return order, err
}

type BuyerUpdatePurchaseOrderDesignParams struct {
	models.JwtClaimsInfo

	PurchaseOrderID     string              `json:"purchase_order_id" param:"purchase_order_id" query:"purchase_order_id" validate:"required"`
	Attachments         *models.Attachments `json:"attachments" param:"attachments" query:"attachments"`
	TechpackAttachments *models.Attachments `json:"techpack_attachments" param:"techpack_attachments" query:"techpack_attachments"`
}

func (r *PurchaseOrderRepo) BuyerUpdatePurchaseOrderDesign(params BuyerUpdatePurchaseOrderDesignParams) (*models.PurchaseOrder, error) {
	order, err := r.GetPurchaseOrder(GetPurchaseOrderParams{
		PurchaseOrderID: params.PurchaseOrderID,
		JwtClaimsInfo:   params.JwtClaimsInfo,
	})
	if err != nil {
		return nil, err
	}

	validStatus := []enums.PoTrackingStatus{enums.PoTrackingStatusNew, enums.PoTrackingStatusDesignRejected, enums.PoTrackingStatusWaitingForApproved}
	if ok := slices.Contains(validStatus, order.TrackingStatus); !ok {
		return nil, errs.ErrPoInvalidToUploadDesign
	}

	var updates models.PurchaseOrder
	err = copier.Copy(&updates, &params)
	if err != nil {
		return nil, err
	}

	updates.TrackingStatus = enums.PoTrackingStatusWaitingForApproved

	r.db.Transaction(func(tx *gorm.DB) error {
		NewPurchaseOrderTrackingRepo(r.db).CreatePurchaseOrderTrackingTx(tx, models.PurchaseOrderTrackingCreateForm{
			PurchaseOrderID: params.PurchaseOrderID,
			ActionType:      enums.PoTrackingActionUpdateDesign,
			UserID:          order.UserID,
			CreatedByUserID: params.JwtClaimsInfo.GetUserID(),
			Metadata: &models.PoTrackingMetadata{
				Before: map[string]interface{}{
					"techpack_attachments": order.TechpackAttachments,
					"attachments":          order.Attachments,
				},
				After: map[string]interface{}{
					"techpack_attachments": params.TechpackAttachments,
					"attachments":          params.Attachments,
				},
			},
		})

		return tx.Model(&models.PurchaseOrder{}).
			Where("id = ?", params.PurchaseOrderID).
			Updates(&updates).Error
	})
	if err != nil {
		return nil, eris.Wrap(err, err.Error())
	}

	order.TrackingStatus = enums.PoTrackingStatusWaitingForApproved

	return order, err
}

type BuyerRejectDesignParams struct {
	models.JwtClaimsInfo

	PurchaseOrderID   string                           `json:"purchase_order_id" param:"purchase_order_id" query:"purchase_order_id" validate:"required"`
	ApproveRejectMeta *models.InquiryApproveRejectMeta `json:"approve_reject_meta" param:"approve_reject_meta" query:"approve_reject_meta"`
}

func (r *PurchaseOrderRepo) BuyerRejectDesign(params BuyerRejectDesignParams) (*models.PurchaseOrder, error) {
	order, err := r.GetPurchaseOrder(GetPurchaseOrderParams{
		PurchaseOrderID: params.PurchaseOrderID,
		JwtClaimsInfo:   params.JwtClaimsInfo,
	})
	if err != nil {
		return nil, err
	}

	validStatus := []enums.PoTrackingStatus{enums.PoTrackingStatusWaitingForApproved}
	if ok := slices.Contains(validStatus, order.TrackingStatus); !ok {
		return nil, errs.ErrPoInvalidToRejectDesign
	}

	var updates = models.PurchaseOrder{
		TrackingStatus:      enums.PoTrackingStatusDesignRejected,
		ApproveRejectMeta:   params.ApproveRejectMeta,
		ApproveDesignAt:     nil,
		TechpackAttachments: nil,
	}

	err = r.db.Transaction(func(tx *gorm.DB) error {
		err = NewPurchaseOrderTrackingRepo(r.db).CreatePurchaseOrderTrackingTx(tx, models.PurchaseOrderTrackingCreateForm{
			PurchaseOrderID: params.PurchaseOrderID,
			ActionType:      enums.PoTrackingActionRejectDesign,
			UserID:          params.JwtClaimsInfo.GetUserID(),
			CreatedByUserID: params.JwtClaimsInfo.GetUserID(),
			Metadata: &models.PoTrackingMetadata{
				After: map[string]interface{}{
					"approve_reject_meta":  params.ApproveRejectMeta,
					"techpack_attachments": order.TechpackAttachments,
				},
			},
		})

		return tx.Select("TrackingStatus", "TechpackAttachments", "ApproveRejectMeta", "ApproveDesignAt").Model(&models.PurchaseOrder{}).
			Where("id = ?", params.PurchaseOrderID).
			Updates(&updates).Error
	})

	if err != nil {
		return nil, eris.Wrap(err, err.Error())
	}

	order.TrackingStatus = updates.TrackingStatus
	order.TechpackAttachments = updates.TechpackAttachments
	order.ApproveRejectMeta = updates.ApproveRejectMeta
	return order, err
}

type BuyerApproveDesignParams struct {
	models.JwtClaimsInfo

	PurchaseOrderID   string                           `json:"purchase_order_id" param:"purchase_order_id" query:"purchase_order_id" validate:"required"`
	ApproveRejectMeta *models.InquiryApproveRejectMeta `json:"approve_reject_meta" param:"approve_reject_meta" query:"approve_reject_meta"`
}

func (r *PurchaseOrderRepo) BuyerApproveDesign(params BuyerApproveDesignParams) (*models.PurchaseOrder, error) {
	order, err := r.GetPurchaseOrder(GetPurchaseOrderParams{
		PurchaseOrderID: params.PurchaseOrderID,
		JwtClaimsInfo:   params.JwtClaimsInfo,
	})
	if err != nil {
		return nil, err
	}

	validStatus := []enums.PoTrackingStatus{enums.PoTrackingStatusWaitingForApproved}
	if ok := slices.Contains(validStatus, order.TrackingStatus); !ok {
		return nil, errs.ErrPoInvalidToApproveDesign
	}

	var updates models.PurchaseOrder
	err = copier.Copy(&updates, &params)
	if err != nil {
		return nil, err
	}

	updates.TrackingStatus = enums.PoTrackingStatusDesignApproved
	updates.SellerTrackingStatus = enums.SellerPoTrackingStatusDesignApprovedByBuyer
	updates.SellerTechpackAttachments = order.TechpackAttachments
	updates.ApproveDesignAt = aws.Int64(time.Now().Unix())

	err = r.db.Transaction(func(tx *gorm.DB) error {
		err = NewPurchaseOrderTrackingRepo(r.db).CreatePurchaseOrderTrackingTx(tx, models.PurchaseOrderTrackingCreateForm{
			PurchaseOrderID: params.PurchaseOrderID,
			ActionType:      enums.PoTrackingActionApproveDesign,
			UserID:          order.UserID,
			CreatedByUserID: params.JwtClaimsInfo.GetUserID(),
			Metadata: &models.PoTrackingMetadata{
				After: map[string]interface{}{
					"techpack_attachments": order.TechpackAttachments,
				},
			},
		})
		if err != nil {
			return eris.Wrap(err, err.Error())
		}

		return tx.Model(&models.PurchaseOrder{}).
			Where("id = ?", params.PurchaseOrderID).
			Updates(&updates).Error
	})

	if err != nil {
		return nil, eris.Wrap(err, err.Error())
	}

	order.TrackingStatus = enums.PoTrackingStatusDesignApproved
	return order, err
}

type AdminSkipDesignParams struct {
	models.JwtClaimsInfo

	PurchaseOrderID string `json:"purchase_order_id" param:"purchase_order_id" query:"purchase_order_id" validate:"required"`
}

func (r *PurchaseOrderRepo) AdminSkipDesign(params AdminSkipDesignParams) (*models.PurchaseOrder, error) {
	order, err := r.GetPurchaseOrder(GetPurchaseOrderParams{
		PurchaseOrderID: params.PurchaseOrderID,
		JwtClaimsInfo:   params.JwtClaimsInfo,
	})
	if err != nil {
		return nil, err
	}

	validStatus := []enums.PoTrackingStatus{enums.PoTrackingStatusNew}
	if ok := slices.Contains(validStatus, order.TrackingStatus); !ok {
		return nil, errs.ErrPoInvalidToApproveDesign
	}

	var updates models.PurchaseOrder
	err = copier.Copy(&updates, &params)
	if err != nil {
		return nil, err
	}

	updates.TrackingStatus = enums.PoTrackingStatusDesignApproved
	updates.SellerTrackingStatus = enums.SellerPoTrackingStatusDesignApprovedByBuyer
	updates.SellerTechpackAttachments = order.TechpackAttachments
	updates.ApproveDesignAt = aws.Int64(time.Now().Unix())

	err = r.db.Model(&models.PurchaseOrder{}).
		Where("id = ?", params.PurchaseOrderID).
		Updates(&updates).Error

	if err != nil {
		return nil, eris.Wrap(err, err.Error())
	}

	order.TrackingStatus = enums.PoTrackingStatusDesignApproved
	return order, err
}

type BuyerConfirmDeliveredParams struct {
	models.JwtClaimsInfo

	PurchaseOrderID string `json:"purchase_order_id" param:"purchase_order_id" query:"purchase_order_id" validate:"required"`
}

func (r *PurchaseOrderRepo) BuyerConfirmDelivered(params BuyerConfirmDeliveredParams) (*models.PurchaseOrder, error) {
	order, err := r.GetPurchaseOrder(GetPurchaseOrderParams{
		PurchaseOrderID: params.PurchaseOrderID,
		JwtClaimsInfo:   params.JwtClaimsInfo,
	})
	if err != nil {
		return nil, err
	}

	var updates models.PurchaseOrder
	err = copier.Copy(&updates, &params)
	if err != nil {
		return nil, err
	}

	updates.TrackingStatus = enums.PoTrackingStatusDeliveryConfirmed
	updates.ReceiverConfirmedAt = values.Int64(time.Now().Unix())

	err = r.db.Transaction(func(tx *gorm.DB) error {
		err = NewPurchaseOrderTrackingRepo(r.db).CreatePurchaseOrderTrackingTx(tx, models.PurchaseOrderTrackingCreateForm{
			PurchaseOrderID: params.PurchaseOrderID,
			ActionType:      enums.PoTrackingActionConfirmDelivered,
			UserID:          params.JwtClaimsInfo.GetUserID(),
			CreatedByUserID: params.JwtClaimsInfo.GetUserID(),
		})
		if err != nil {
			return err
		}

		return tx.Model(&models.PurchaseOrder{}).
			Where("id = ?", params.PurchaseOrderID).
			Updates(&updates).Error
	})
	if err != nil {
		return nil, eris.Wrap(err, err.Error())
	}
	order.TrackingStatus = updates.TrackingStatus
	order.ReceiverConfirmedAt = updates.ReceiverConfirmedAt
	return order, err
}

func (r *PurchaseOrderRepo) PurchaseOrderAssignPIC(params models.PurchaseOrderAssignPICParam) (updates models.PurchaseOrder, err error) {
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
	var purchaseOrder models.PurchaseOrder
	err = r.db.Select("ID", "AssigneeIDs", "UserID").First(&purchaseOrder, "id = ?", params.PurchaseOrderID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = errs.ErrPONotFound
		}
		return
	}
	updates.AssigneeIDs = params.AssigneeIDs

	err = r.db.Transaction(func(tx *gorm.DB) error {
		err = tx.Model(&updates).Clauses(clause.Returning{}).
			Where("id = ?", params.PurchaseOrderID).Updates(&updates).Error

		var chatRoom models.ChatRoom
		err = tx.Select("ID").Where(map[string]interface{}{
			"inquiry_id":             updates.InquiryID,
			"purchase_order_id":      updates.ID,
			"bulk_purchase_order_id": "",
			"buyer_id":               updates.UserID,
		}).First(&chatRoom).Error
		if err != nil && !r.db.IsRecordNotFoundError(err) {
			return err
		}

		if r.db.IsRecordNotFoundError(err) {
			chatRoom.PurchaseOrderID = updates.ID
			chatRoom.InquiryID = updates.InquiryID
			chatRoom.HostID = params.GetUserID()
			if err := tx.Create(&chatRoom).Error; err != nil {
				return err
			}
		}
		var chatRoomUsers = []*models.ChatRoomUser{{RoomID: chatRoom.ID, UserID: purchaseOrder.UserID}}
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

type AdminUnarchivePurchaseOrderParams struct {
	models.JwtClaimsInfo

	PurchaseOrderID string `json:"purchase_order_id" query:"purchase_order_id" param:"purchase_order_id" validate:"required"`
}

func (r *PurchaseOrderRepo) AdminUnarchivePurchaseOrder(params AdminUnarchivePurchaseOrderParams) error {
	return r.db.Unscoped().Model(&models.PurchaseOrder{}).Where("id = ?", params.PurchaseOrderID).UpdateColumn("DeletedAt", nil).Error
}

type AdminArchivePurchaseOrderParams struct {
	models.JwtClaimsInfo

	PurchaseOrderID string `json:"purchase_order_id" query:"purchase_order_id" param:"purchase_order_id" validate:"required"`
}

func (r *PurchaseOrderRepo) AdminArchivePurchaseOrder(params AdminArchivePurchaseOrderParams) error {
	return r.db.Unscoped().
		Model(&models.PurchaseOrder{}).
		Where("id = ? AND tracking_status = ?", params.PurchaseOrderID, enums.PoTrackingStatusNew).
		UpdateColumn("DeletedAt", time.Now().Unix()).Error
}

type AdminDeletePurchaseOrderParams struct {
	models.JwtClaimsInfo

	PurchaseOrderID string `json:"purchase_order_id" query:"purchase_order_id" param:"purchase_order_id" validate:"required"`
}

func (r *PurchaseOrderRepo) AdminDeletePurchaseOrder(params AdminDeletePurchaseOrderParams) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		var sqlResult = tx.Unscoped().
			Model(&models.PurchaseOrder{}).
			Where("id = ? AND deleted_at is not null", params.PurchaseOrderID).
			Delete("DeletedAt", time.Now().Unix())
		if sqlResult.Error != nil {
			return sqlResult.Error
		}

		if sqlResult.RowsAffected == 0 {
			return errs.ErrInquiryNotFound
		}

		var err = tx.Unscoped().Delete(&models.PaymentTransaction{}, "purchase_order_id = ?", params.PurchaseOrderID).Error
		if err != nil {
			return err
		}

		return nil
	})
}

type PurchaseOrderFeedbackParams struct {
	PurchaseOrderID string `json:"purchase_order_id" param:"purchase_order_id"`
	models.PurchaseOrderFeedback
}

func (r *PurchaseOrderRepo) BuyerGivePurchaseOrderFeedback(params PurchaseOrderFeedbackParams) (err error) {
	var updates = &models.PurchaseOrder{
		Feedback: &params.PurchaseOrderFeedback,
	}
	err = r.db.Model(&models.PurchaseOrder{}).Where("id = ?", params.PurchaseOrderID).Updates(updates).Error
	return
}

type PurchaseApproveRawMaterialParams struct {
	models.JwtClaimsInfo

	PurchaseOrderID string `json:"purchase_order_id" param:"purchase_order_id" query:"purchase_order_id" validate:"required"`

	ItemIDs []string `json:"item_ids" param:"item_ids" query:"item_ids" validate:"required"`
}

func (r *PurchaseOrderRepo) PurchaseOrderApproveRawMaterial(params PurchaseApproveRawMaterialParams) (*models.PurchaseOrder, error) {
	order, err := r.GetPurchaseOrder(GetPurchaseOrderParams{
		PurchaseOrderID: params.PurchaseOrderID,
		JwtClaimsInfo:   params.JwtClaimsInfo,
	})

	if err != nil {
		return nil, err
	}

	if *order.PoRawMaterials == nil {
		return nil, errs.ErrPoInvalidToApproveRawMaterial
	}

	var updates models.PurchaseOrder
	updates.PoRawMaterials = order.PoRawMaterials
	var shouldTracking = false
	for _, existItem := range *updates.PoRawMaterials {
		// if empty itemsIDs means it approve by system
		if (slices.Contains(params.ItemIDs, existItem.ReferenceID) && !values.BoolValue(existItem.BuyerApproved)) || len(params.ItemIDs) == 0 {
			existItem.BuyerApproved = values.Bool(true)
			existItem.Status = enums.PoRawMaterialStatusApproved
			shouldTracking = true
		}
	}

	if shouldTracking {
		err = r.db.Transaction(func(tx *gorm.DB) error {
			err = NewPurchaseOrderTrackingRepo(r.db).CreatePurchaseOrderTrackingTx(tx, models.PurchaseOrderTrackingCreateForm{
				PurchaseOrderID: params.PurchaseOrderID,
				ActionType:      enums.PoTrackingActionApproveRawMaterial,
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

			return tx.Model(&models.PurchaseOrder{}).
				Where("id = ?", params.PurchaseOrderID).
				Updates(&updates).Error
		})
		if err != nil {
			return nil, eris.Wrap(err, err.Error())
		}
	}

	order.PoRawMaterials = updates.PoRawMaterials
	return order, err
}

type PoCommentMarkSeenParams struct {
	models.JwtClaimsInfo

	PurchaseOrderID string `json:"purchase_order_id" param:"purchase_order_id"`
	FileKey         string `json:"file_key" param:"file_key"`
}

func (r *PurchaseOrderRepo) DesignCommentMarkSeen(params PoCommentMarkSeenParams) error {
	var err = r.db.Model(&models.Comment{}).
		Where("seen_at IS NULL").
		Where("user_id != ?", params.GetUserID()).
		Where("target_type = ? AND target_id = ? AND file_key = ?", enums.CommentTargetTypePurchaseOrderDesign, params.PurchaseOrderID, params.FileKey).
		Update("seen_at", time.Now().Unix()).Error

	return err
}

type PoCommentStatusCountParams struct {
	models.JwtClaimsInfo
	PurchaseOrderID string `json:"purchase_order_id" param:"purchase_order_id"`
}

type PoDesignCommentStatusCountItem struct {
	FileKey     string `json:"file_key"`
	UnseenCount int64  `json:"unseen_count"`
}

func (r *PurchaseOrderRepo) DesignCommentStatusCount(params PoCommentStatusCountParams) ([]*PoDesignCommentStatusCountItem, error) {
	order, err := r.GetPurchaseOrder(GetPurchaseOrderParams{
		PurchaseOrderID: params.PurchaseOrderID,
		JwtClaimsInfo:   params.JwtClaimsInfo,
	})
	if err != nil {
		return nil, eris.Wrap(err, err.Error())
	}

	var items []*PoDesignCommentStatusCountItem

	if order.TechpackAttachments != nil && len(*order.TechpackAttachments) > 0 {
		for _, fileItem := range *order.TechpackAttachments {
			count := int64(0)
			if fileItem.FileKey != "" {
				r.db.Model(&models.Comment{}).Where("user_id != ? AND target_type = ? AND target_id = ? AND file_key = ? AND seen_at IS NULL", params.GetUserID(), enums.CommentTargetTypePurchaseOrderDesign, params.PurchaseOrderID, fileItem.FileKey).Count(&count)
				items = append(items, &PoDesignCommentStatusCountItem{
					FileKey:     fileItem.FileKey,
					UnseenCount: count,
				})
			}

		}
	}

	return items, nil
}

type AdminAssignMakerParams struct {
	models.JwtClaimsInfo

	PurchaseOrderID string `json:"purchase_order_id" param:"purchase_order_id" query:"purchase_order_id" validate:"required"`

	MakerID string `json:"maker_id"`
}

func (r *PurchaseOrderRepo) AdminAssignMaker(params AdminAssignMakerParams) (*models.PurchaseOrder, error) {
	purchaseOrder, err := r.GetPurchaseOrder(GetPurchaseOrderParams{
		PurchaseOrderID: params.PurchaseOrderID,
	})
	if err != nil {
		return nil, err
	}

	var inquirySeller models.InquirySeller

	if params.MakerID != "inflow" {
		var seller models.User
		err = r.db.Select("ID", "Email", "CompanyName").First(&seller, "id = ?", params.MakerID).Error
		if err != nil {
			return nil, err
		}

		err = r.db.Select("ID", "UserID", "InquiryID").First(&inquirySeller, "user_id = ? AND inquiry_id = ?", seller.ID, purchaseOrder.InquiryID).Error
		if err != nil {
			return nil, err
		}

		purchaseOrder.SampleMaker = &seller
	}

	err = r.db.Model(&models.PurchaseOrder{}).Where("id = ?", params.PurchaseOrderID).Update("SampleMakerID", params.MakerID).Error
	if err != nil {
		return nil, err
	}

	if params.MakerID != "inflow" && inquirySeller.ID != "" {
		err = r.db.Model(&models.InquirySeller{}).Where("id = ?", inquirySeller.ID).Update("PurchaseOrderID", purchaseOrder.ID).Error
		if err != nil {
			return nil, err
		}
	}

	return purchaseOrder, err
}

func (r *PurchaseOrderRepo) ExportExcel(params PaginatePurchaseOrdersParams) (*models.Attachment, error) {
	var s3Client = s3.New(r.db.Configuration)
	params.IncludeAssignee = true
	params.IsQueryAll = true
	var result = r.PaginatePurchaseOrders(params)
	if result == nil || result.Records == nil {
		return nil, errors.New("empty response")
	}

	trans, ok := result.Records.([]*models.PurchaseOrder)
	if !ok {
		return nil, errors.New("empty response")
	}

	fileContent, err := models.PurchaseOrders(trans).ToExcel()
	if err != nil {
		return nil, err
	}

	var contentType = models.ContentTypeXLSX
	url := fmt.Sprintf("uploads/purchase_orders/export/export_purchase_order_user_%s%s", params.GetUserID(), contentType.GetExtension())
	_, err = s3Client.UploadFile(s3.UploadFileParams{
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

type RefundPurchaseOrderParams struct {
	models.JwtClaimsInfo
	PurchaseOrderID string `json:"purchase_order_id" param:"purchase_order_id" validate:"required"`
	Reason          string `json:"reason" validate:"required"`
}

func (r *PurchaseOrderRepo) RefundPurchaseOrder(params RefundPurchaseOrderParams) (*models.PurchaseOrder, error) {
	order, err := r.GetPurchaseOrder(GetPurchaseOrderParams{
		PurchaseOrderID: params.PurchaseOrderID,
		JwtClaimsInfo:   params.JwtClaimsInfo,
	})

	if err != nil {
		return nil, err
	}

	var canCancel = order.Status == enums.PurchaseOrderStatusPaid && order.TrackingStatus == enums.PoTrackingStatusNew
	if !canCancel {
		return nil, eris.Wrapf(errs.ErrPOIsNotAbleToCancel, "Order %s with status %s-%s is not able to cancel", order.ReferenceID, order.Status, order.TrackingStatus)
	}

	if order.PaymentIntentID == "" {
		return nil, eris.Wrapf(errs.ErrPOIsNotAbleToCancel, "Order %s has empty payment intent", order.ReferenceID)
	}

	var admin models.User
	err = r.db.Select("ID", "Name", "Email").First(&admin, "id = ?", params.GetUserID()).Error
	if err != nil {
		return nil, err
	}

	var updatePaymentTransaction = models.PaymentTransaction{
		Status:       enums.PaymentStatusRefunded,
		RefundReason: params.Reason,
	}

	var updatePurchaseOrder = models.PurchaseOrder{
		Status:       enums.PurchaseOrderStatusCanceled,
		RefundReason: params.Reason,
	}

	r.db.Transaction(func(tx *gorm.DB) error {
		var sqlResult = tx.Model(&models.PaymentTransaction{}).Where("payment_intent_id = ?", order.PaymentIntentID).Updates(&updatePaymentTransaction)
		if sqlResult.Error != nil {
			return sqlResult.Error
		}

		if sqlResult.RowsAffected == 0 {
			return errs.ErrPaymentTransactionIsInvalid
		}

		sqlResult = tx.Delete(&models.Invoice{}, "invoice_number = ?", order.InvoiceNumber)
		if sqlResult.Error != nil {
			return sqlResult.Error
		}

		if sqlResult.RowsAffected == 0 {
			return errs.ErrPaymentTransactionIsInvalid
		}

		sqlResult = tx.Model(&models.PurchaseOrder{}).Where("id = ?", order.ID).Updates(&updatePurchaseOrder)
		if sqlResult.Error != nil {
			return sqlResult.Error
		}

		if sqlResult.RowsAffected == 0 {
			return errs.ErrPaymentTransactionIsInvalid
		}

		var audit = models.InquiryAudit{
			InquiryID:       order.InquiryID,
			PurchaseOrderID: order.ID,
			UserID:          order.UserID,
			ActionType:      enums.AuditActionTypeInquiryAdminRefund,
			Description:     fmt.Sprintf("%s has refunded with reason %s", admin.Name, params.Reason),
			Metadata: &models.InquiryAuditMetadata{
				After: map[string]interface{}{
					"refunded_by": map[string]interface{}{
						"id":    admin.ID,
						"name":  admin.Name,
						"email": admin.Email,
					},
				},
			},
		}

		err = tx.Create(&audit).Error
		if err != nil {
			return err
		}

		_, err = stripehelper.GetInstance().RefundPaymentIntent(stripehelper.RefundPaymentIntentParams{
			PaymentIntentID: order.PaymentIntentID,
			Metadata: map[string]string{
				"purchase_order_id":           order.ID,
				"purchase_order_reference_id": order.ReferenceID,
				"refunded_by_name":            admin.Name,
				"refunded_by_email":           admin.Email,
			},
		})

		return err
	})
	if err != nil {
		return nil, err
	}

	order.RefundReason = updatePurchaseOrder.RefundReason
	order.Status = updatePurchaseOrder.Status

	return order, err
}

type PurchaseOrderStageCommentsParams struct {
	models.JwtClaimsInfo

	PurchaseOrderID string              `json:"purchase_order_id" param:"purchase_order_id" query:"purchase_order_id" validate:"required"`
	Comment         string              `json:"comment,omitempty" validate:"required"`
	Attachments     *models.Attachments `json:"attachments,omitempty"`
}

func (r *PurchaseOrderRepo) StageCommentsCreate(params PurchaseOrderStageCommentsParams) (*models.PurchaseOrder, error) {
	order, err := r.GetPurchaseOrder(GetPurchaseOrderParams{
		PurchaseOrderID: params.PurchaseOrderID,
		JwtClaimsInfo:   params.JwtClaimsInfo,
	})

	if err != nil {
		return nil, err
	}

	err = NewPurchaseOrderTrackingRepo(r.db).CreatePurchaseOrderTrackingTx(r.db.DB, models.PurchaseOrderTrackingCreateForm{
		PurchaseOrderID: params.PurchaseOrderID,
		ActionType:      enums.PoTrackingActionStageComment,
		UserID:          order.UserID,
		CreatedByUserID: params.JwtClaimsInfo.GetUserID(),
		FromStatus:      order.TrackingStatus,
		ToStatus:        order.TrackingStatus,
		UserGroup:       enums.PoTrackingUserGroupBuyer,
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

type UpdatePurchaseOrderParams struct {
	models.JwtClaimsInfo
	PurchaseOrderID string `param:"purchase_order_id" validate:"required"`

	UserID string `json:"-"`

	Attachments         *models.Attachments `json:"attachments,omitempty"`
	Document            *models.Attachments `json:"document,omitempty"`
	Design              *models.Attachments `json:"design,omitempty"`
	FabricAttachments   *models.Attachments `json:"fabric_attachments,omitempty"`
	TechpackAttachments *models.Attachments `json:"techpack_attachments,omitempty"`
	DesignNote          string              `json:"design_note,omitempty"`
}

func (r *PurchaseOrderRepo) UpdatePurchaseOrder(params UpdatePurchaseOrderParams) (*models.PurchaseOrder, error) {
	order, err := r.GetPurchaseOrder(GetPurchaseOrderParams{
		PurchaseOrderID: params.PurchaseOrderID,
		JwtClaimsInfo:   params.JwtClaimsInfo,
		UserID:          params.UserID,
		IncludeUsers:    true,
	})
	if err != nil {
		return nil, err
	}

	var updates = models.PurchaseOrder{
		Attachments:         params.Attachments,
		Document:            params.Document,
		Design:              params.Design,
		FabricAttachments:   params.FabricAttachments,
		TechpackAttachments: params.TechpackAttachments,
		DesignNote:          params.DesignNote,
	}
	err = r.db.Model(&models.PurchaseOrder{}).
		Where("id = ?", params.PurchaseOrderID).
		Updates(&updates).Error
	if err != nil {
		return nil, err
	}

	err = copier.CopyWithOption(order, &updates, copier.Option{IgnoreEmpty: true, DeepCopy: true})
	if err != nil {
		return nil, err
	}

	return order, err

}

type CreateBulkPurchaseOrderParams struct {
	models.JwtClaimsInfo

	PurchaseOrderID string `json:"purchase_order_id" query:"purchase_order_id" param:"purchase_order_id" validate:"required"`
}

func (r *PurchaseOrderRepo) CreateBulkPurchaseOrder(params CreateBulkPurchaseOrderParams) (*models.BulkPurchaseOrder, error) {
	cancel, err := r.db.Locker.AcquireLock(fmt.Sprintf("purchase_order_%s", params.PurchaseOrderID), time.Minute)
	if err != nil {
		return nil, err
	}
	defer cancel()

	var purchaseOrder models.PurchaseOrder
	if params.GetRole().IsAdmin() {
		err = r.db.Select("ID", "Status", "TrackingStatus", "Currency", "UserID", "InquiryID", "AssigneeIDs", "OrderGroupID", "SampleMakerID").First(&purchaseOrder, "id = ?", params.PurchaseOrderID).Error
	} else {
		err = r.db.Select("ID", "Status", "TrackingStatus", "Currency", "UserID", "InquiryID", "AssigneeIDs", "OrderGroupID", "SampleMakerID").First(&purchaseOrder, "id = ? AND user_id = ?", params.PurchaseOrderID, params.GetUserID()).Error
	}
	if err != nil {
		return nil, err
	}

	if purchaseOrder.TrackingStatus != enums.PoTrackingStatusDeliveryConfirmed {
		return nil, errs.ErrInquirySampleOrderIsNotPaid
	}

	var inquiry models.Inquiry
	if purchaseOrder.InquiryID != "" {
		err = r.db.Select("ID", "Title", "SkuNote", "UserID", "Currency", "ShippingAddressID").First(&inquiry, "id = ?", purchaseOrder.InquiryID).Error
		if err != nil {
			return nil, err
		}
	}

	record, err := NewBulkPurchaseOrderRepo(r.db).GetBulkPurchaseOrder(GetBulkPurchaseOrderParams{
		JwtClaimsInfo:   params.JwtClaimsInfo,
		PurchaseOrderID: params.PurchaseOrderID,
	})
	if err != nil {
		if r.db.IsRecordNotFoundError(err) {
			var bulkPurchaseOrder = models.BulkPurchaseOrder{
				UserID:            purchaseOrder.UserID,
				Currency:          purchaseOrder.Currency,
				InquiryID:         purchaseOrder.InquiryID,
				PurchaseOrderID:   purchaseOrder.ID,
				Status:            enums.BulkPurchaseOrderStatusNew,
				TrackingStatus:    enums.BulkPoTrackingStatusNew,
				ShippingAddressID: inquiry.ShippingAddressID,
				ShippingAddress:   inquiry.ShippingAddress,
				OrderGroupID:      purchaseOrder.OrderGroupID,
				// SellerID:          purchaseOrder.SampleMakerID,
				ProductName: inquiry.Title,
				Note:        inquiry.SkuNote,
			}

			if params.GetRole().IsAdmin() {
				bulkPurchaseOrder.AssigneeIDs = append(bulkPurchaseOrder.AssigneeIDs, params.GetUserID())
			} else {
				bulkPurchaseOrder.AssigneeIDs = purchaseOrder.AssigneeIDs
			}
			err = r.db.Create(&bulkPurchaseOrder).Error
			if err != nil {
				return nil, err
			}
			bulkPurchaseOrder.Inquiry = &inquiry
			bulkPurchaseOrder.PurchaseOrder = &purchaseOrder
			return &bulkPurchaseOrder, err
		} else {
			return nil, err
		}
	}

	return record, err

}

type CreatePurchaseOrderParams struct {
	models.JwtClaimsInfo

	UserID          string                       `json:"user_id" validate:"required"`
	Currency        enums.Currency               `json:"currency" validate:"required"`
	Items           []*models.PurchaseOrderItem  `json:"items" validate:"required"`
	ProductWeight   *float64                     `json:"product_weight" validate:"required"`
	ShippingFee     price.Price                  `json:"shipping_fee"`
	SizeChart       enums.InquirySizeChart       `json:"size_chart"`
	ShippingAddress *models.Address              `json:"shipping_address"`
	TaxPercentage   *float64                     `json:"tax_percentage"`
	Quotations      models.InquiryQuotationItems `json:"quotations,omitempty"`

	IsPaid                bool               `json:"is_paid"`
	TransactionRefID      string             `json:"transaction_ref_id,omitempty" validate:"required_if=IsPaid true"`
	TransactionAttachment *models.Attachment `json:"transaction_attachment,omitempty" validate:"required_if=IsPaid true"`

	InquiryID string `json:"inquiry_id"`
}

func (r *PurchaseOrderRepo) CreatePurchaseOrder(params CreatePurchaseOrderParams) (*models.PurchaseOrder, error) {
	var purchaseOrder = models.PurchaseOrder{
		Currency:   params.Currency,
		Quotations: params.Quotations,
		UserID:     params.UserID,
		Items:      params.Items,
	}
	purchaseOrder.ID = helper.GenerateXID()
	purchaseOrder.AssigneeIDs = append(purchaseOrder.AssigneeIDs, params.GetUserID())
	purchaseOrder.TaxPercentage = params.TaxPercentage
	purchaseOrder.ShippingFee = &params.ShippingFee
	purchaseOrder.ProductWeight = params.ProductWeight
	purchaseOrder.Quotations = params.Quotations
	purchaseOrder.InquiryID = params.InquiryID

	if params.ShippingAddress != nil {
		if err := params.ShippingAddress.CreateOrUpdate(r.db); err == nil {
			purchaseOrder.ShippingAddressID = params.ShippingAddress.ID
			purchaseOrder.ShippingAddress = params.ShippingAddress
		}
	}

	var subTotalPrice = price.NewFromFloat(0)
	var items = lo.Map(purchaseOrder.Items, func(item *models.PurchaseOrderItem, index int) *models.PurchaseOrderItem {
		item.PurchaseOrderID = purchaseOrder.ID
		item.TotalPrice = item.UnitPrice.MultipleInt(item.Quantity).ToPtr()
		subTotalPrice = subTotalPrice.AddPtr(item.TotalPrice)
		return item
	})

	purchaseOrder.Items = items

	if params.IsPaid {
		purchaseOrder.MarkAsPaidAt = values.Int64(time.Now().Unix())
		purchaseOrder.TransactionRefID = params.TransactionRefID
		purchaseOrder.Attachments = &models.Attachments{
			params.TransactionAttachment,
		}

		var paymentTrans = models.PaymentTransaction{
			PurchaseOrderID:  purchaseOrder.ID,
			MarkAsPaidAt:     purchaseOrder.MarkAsPaidAt,
			TransactionRefID: params.TransactionRefID,
			Attachments:      purchaseOrder.Attachments,
		}
		var err = r.db.Transaction(func(tx *gorm.DB) error {
			var err = tx.Create(&items).Error
			if err != nil {
				return err
			}

			err = tx.Create(&purchaseOrder).Error
			if err != nil {
				return err
			}

			return tx.Create(&paymentTrans).Error
		})
		if err != nil {
			return nil, err
		}

		return &purchaseOrder, nil
	}

	var err = r.db.Transaction(func(tx *gorm.DB) error {
		if len(items) > 0 {
			var err = tx.Create(&items).Error
			if err != nil {
				return err
			}

		}

		return tx.Create(&purchaseOrder).Error
	})
	if err != nil {
		return nil, err
	}

	return &purchaseOrder, err
}

type CreatePurchaseOrderPaymentLinkParams struct {
	models.JwtClaimsInfo

	PurchaseOrderID string                      `param:"purchase_order_id" validate:"required"`
	Items           []*models.PurchaseOrderItem `json:"items" validate:"required"`
	ProductWeight   *float64                    `json:"product_weight" param:"product_weight" query:"product_weight"`
	ShippingFee     price.Price                 `json:"shipping_fee"`
	TaxPercentage   *float64                    `json:"tax_percentage" validate:"min=0,max=100"`
}

func (r *PurchaseOrderRepo) CreatePurchaseOrderPaymentLinkParams(params CreatePurchaseOrderPaymentLinkParams) (*models.PurchaseOrder, error) {
	cancel, err := r.db.Locker.AcquireLock(fmt.Sprintf("purchase_order_payment_link_%s", params.PurchaseOrderID), time.Second*30)
	if err != nil {
		return nil, err
	}
	defer cancel()

	var purchaseOrder models.PurchaseOrder
	err = r.db.First(&purchaseOrder, "id = ?", params.PurchaseOrderID).Error
	if err != nil {
		return nil, eris.Wrap(err, err.Error())
	}

	var buyer models.User
	err = r.db.Select("Name", "Email").First(&buyer, "id = ?", purchaseOrder.UserID).Error
	if err != nil {
		return nil, eris.Wrap(err, err.Error())
	}

	var items []*models.PurchaseOrderItem
	err = copier.Copy(&items, &params.Items)
	if err != nil {
		return nil, err
	}

	var subTotalPrice = price.NewFromFloat(0)
	items = lo.Map(items, func(item *models.PurchaseOrderItem, index int) *models.PurchaseOrderItem {
		item.PurchaseOrderID = purchaseOrder.ID
		item.TotalPrice = item.UnitPrice.MultipleInt(item.Quantity).ToPtr()
		subTotalPrice = subTotalPrice.AddPtr(item.TotalPrice)
		return item
	})

	err = r.db.Transaction(func(tx *gorm.DB) error {
		var err = tx.Unscoped().Delete(&models.PurchaseOrder{}, "purchase_order_id = ?", purchaseOrder.ID).Error
		if err != nil {
			return err
		}

		err = tx.Create(&items).Error

		if err != nil {
			return err
		}
		return err
	})
	if err != nil {
		return nil, err
	}

	purchaseOrder.ShippingFee = &params.ShippingFee
	purchaseOrder.TaxPercentage = params.TaxPercentage
	purchaseOrder.ProductWeight = params.ProductWeight
	purchaseOrder.Items = items
	purchaseOrder.SubTotal = &subTotalPrice
	err = purchaseOrder.UpdatePrices()
	if err != nil {
		return nil, err
	}

	stripeConfig, err := stripehelper.GetCurrencyConfig(purchaseOrder.Inquiry.Currency)
	if err != nil {
		return nil, err
	}

	var paymentLineItems []*stripe.PaymentLinkLineItemParams
	for _, cartItem := range items {
		var itemNames []string
		if cartItem.FabricName != "" {
			itemNames = append(itemNames, cartItem.FabricName)
		} else {
			if cartItem.Fabric != nil && cartItem.Fabric.FabricType != "" {
				itemNames = append(itemNames, cartItem.Fabric.FabricType)
			}
		}

		if cartItem.Size != "" {
			itemNames = append(itemNames, cartItem.Size)
		}

		if cartItem.Color != "" {
			itemNames = append(itemNames, cartItem.Color)
		}

		priceItem, err := stripePrice.New(&stripe.PriceParams{
			Currency:   stripe.String(string(purchaseOrder.Currency)),
			UnitAmount: stripe.Int64(cartItem.UnitPrice.MultipleInt(stripeConfig.SmallestUnitFactor).ToInt64()),
			ProductData: &stripe.PriceProductDataParams{
				Name: stripe.String(strings.Join(itemNames, "/")),
			},
		})
		if err != nil {
			return nil, err
		}

		paymentLineItems = append(paymentLineItems, &stripe.PaymentLinkLineItemParams{
			Quantity: stripe.Int64(int64(cartItem.Quantity)),
			Price:    &priceItem.ID,
		})
	}

	// Append for tax and shipping
	if purchaseOrder.Tax.GreaterThan(0) {
		priceItemTax, err := stripePrice.New(&stripe.PriceParams{
			Currency:   stripe.String(string(purchaseOrder.Currency)),
			UnitAmount: stripe.Int64(purchaseOrder.Tax.MultipleInt(stripeConfig.SmallestUnitFactor).ToInt64()),
			ProductData: &stripe.PriceProductDataParams{
				Name: stripe.String("Tax"),
			},
		})
		if err != nil {
			return nil, err
		}
		paymentLineItems = append(paymentLineItems, &stripe.PaymentLinkLineItemParams{
			Quantity: stripe.Int64(1),
			Price:    &priceItemTax.ID,
		})
	}

	if purchaseOrder.ShippingFee.GreaterThan(0) {
		priceItemShippingFee, err := stripePrice.New(&stripe.PriceParams{
			Currency:   stripe.String(string(purchaseOrder.Currency)),
			UnitAmount: stripe.Int64(purchaseOrder.ShippingFee.MultipleInt(stripeConfig.SmallestUnitFactor).ToInt64()),
			ProductData: &stripe.PriceProductDataParams{
				Name: stripe.String("Shipping Fee"),
			},
		})
		if err != nil {
			return nil, err
		}
		paymentLineItems = append(paymentLineItems, &stripe.PaymentLinkLineItemParams{
			Quantity: stripe.Int64(1),
			Price:    &priceItemShippingFee.ID,
		})
	}

	if purchaseOrder.TransactionFee.GreaterThan(0) {
		priceItemTransactionFee, err := stripePrice.New(&stripe.PriceParams{
			Currency:   stripe.String(string(purchaseOrder.Currency)),
			UnitAmount: stripe.Int64(purchaseOrder.TransactionFee.MultipleInt(stripeConfig.SmallestUnitFactor).ToInt64()),
			ProductData: &stripe.PriceProductDataParams{
				Name: stripe.String("Transaction Fee"),
			},
		})
		if err != nil {
			return nil, err
		}
		paymentLineItems = append(paymentLineItems, &stripe.PaymentLinkLineItemParams{
			Quantity: stripe.Int64(1),
			Price:    &priceItemTransactionFee.ID,
		})
	}

	var stripeParams = stripehelper.CreatePaymentLinkParams{
		Currency: purchaseOrder.Inquiry.Currency,
		Metadata: map[string]string{
			"inquiry_id":                  purchaseOrder.InquiryID,
			"inquiry_reference_id":        purchaseOrder.Inquiry.ReferenceID,
			"purchase_order_id":           purchaseOrder.ID,
			"purchase_order_reference_id": purchaseOrder.ReferenceID,
			"cart_item_ids":               strings.Join(purchaseOrder.CartItemIDs, ","),
			"action_source":               string(stripehelper.ActionSourceInquiryPayment),
		},
		RedirectURL: fmt.Sprintf("%s/samples/%s", r.db.Configuration.BrandPortalBaseURL, purchaseOrder.ID),
		LineItems:   paymentLineItems,
	}

	pl, err := stripehelper.GetInstance().CreatePaymentLink(stripeParams)
	if err != nil {
		return nil, err
	}
	purchaseOrder.PaymentLink = helper.AddURLQuery(pl.URL,
		map[string]string{
			"client_reference_id": purchaseOrder.ReferenceID,
			"prefilled_email":     buyer.Email,
		},
	)
	purchaseOrder.PaymentLinkID = pl.ID

	err = r.db.Model(&models.PurchaseOrder{}).Where("id = ?", purchaseOrder.ID).Updates(&purchaseOrder).Error
	if err != nil {
		return nil, eris.Wrap(err, err.Error())
	}

	return &purchaseOrder, nil

}

type CreateMultiPurchaseOrderPaymentLinkParams struct {
	models.JwtClaimsInfo

	Data []CreateMultiPurchaseOrderPaymentLinkItem `param:"data" validate:"required"`
}

type CreateMultiPurchaseOrderPaymentLinkItem struct {
	PurchaseOrderID string                      `json:"purchase_order_id" param:"purchase_order_id" validate:"required"`
	Items           []*models.PurchaseOrderItem `json:"items" param:"items"`
	ProductWeight   *float64                    `json:"product_weight" param:"product_weight" query:"product_weight"`
	ShippingFee     price.Price                 `json:"shipping_fee"`
	TaxPercentage   *float64                    `json:"tax_percentage" validate:"min=0,max=100"`
}

type CreateMultiPurchaseOrderPaymentLinkResponse struct {
	PurchaseOrders []*models.PurchaseOrder `json:"purchase_orders"`
	PaymentLink    string                  `json:"payment_link"`
}

func (r *PurchaseOrderRepo) MultiPurchaseOrderCreatePaymentLink(params CreateMultiPurchaseOrderPaymentLinkParams) (*CreateMultiPurchaseOrderPaymentLinkResponse, error) {
	var currency enums.Currency
	var paymentLineItems []*stripe.PaymentLinkLineItemParams
	var checkoutSessionID = helper.GenerateCheckoutSessionID()
	var purchaseOrderIDs []string
	var purchaseOrders []*models.PurchaseOrder
	var stripeConfig *stripehelper.Config

	var totalTax = price.NewFromFloat(0)
	var totalTransactionFee = price.NewFromFloat(0)
	var totalShippingFee = price.NewFromFloat(0)
	var buyer models.User
	for _, item := range params.Data {
		cancel, err := r.db.Locker.AcquireLock(fmt.Sprintf("purchase_order_payment_link_%s", item.PurchaseOrderID), time.Second*30)
		if err != nil {
			return nil, err
		}
		defer cancel()

		var purchaseOrder models.PurchaseOrder
		err = r.db.First(&purchaseOrder, "id = ?", item.PurchaseOrderID).Error
		if err != nil {
			return nil, eris.Wrap(err, err.Error())
		}
		currency = purchaseOrder.Currency

		if buyer.ID == "" {
			err = r.db.Select("ID", "Name", "Email").First(&buyer, "id = ?", purchaseOrder.UserID).Error
			if err != nil {
				return nil, eris.Wrap(err, err.Error())
			}
		}

		var items []*models.PurchaseOrderItem
		if len(purchaseOrder.Items) == 0 {
			var err = r.db.Unscoped().Delete(&models.PurchaseOrderItem{}, "purchase_order_id = ?", purchaseOrder.ID).Error
			if err != nil {
				return nil, err
			}
		} else {
			err = copier.Copy(&items, &purchaseOrder.Items)
			if err != nil {
				return nil, err
			}

			var subTotalPrice = price.NewFromFloat(0)
			items = lo.Map(items, func(item *models.PurchaseOrderItem, index int) *models.PurchaseOrderItem {
				item.PurchaseOrderID = purchaseOrder.ID
				item.CheckoutSessionID = checkoutSessionID
				item.TotalPrice = item.UnitPrice.MultipleInt(item.Quantity).ToPtr()
				subTotalPrice = subTotalPrice.AddPtr(item.TotalPrice)
				return item
			})

			err = r.db.Transaction(func(tx *gorm.DB) error {
				var err = tx.Unscoped().Delete(&models.PurchaseOrderItem{}, "purchase_order_id = ?", purchaseOrder.ID).Error
				if err != nil {
					return err
				}

				err = tx.Create(&items).Error

				if err != nil {
					return err
				}

				purchaseOrder.Items = items

				err = purchaseOrder.UpdatePrices()
				if err != nil {
					return err
				}

				err = tx.Omit(clause.Associations).Model(&models.PurchaseOrder{}).Where("id = ?", purchaseOrder.ID).Updates(&purchaseOrder).Error
				if err != nil {
					return eris.Wrap(err, err.Error())
				}

				return err
			})
			if err != nil {
				return nil, err
			}

		}

		stripeConfig, err = stripehelper.GetCurrencyConfig(currency)
		if err != nil {
			return nil, err
		}

		for _, cartItem := range items {
			var itemNames []string
			if cartItem.FabricName != "" {
				itemNames = append(itemNames, cartItem.FabricName)
			} else {
				if cartItem.Fabric != nil && cartItem.Fabric.FabricType != "" {
					itemNames = append(itemNames, cartItem.Fabric.FabricType)
				}
			}

			if cartItem.Size != "" {
				itemNames = append(itemNames, cartItem.Size)
			}

			if cartItem.Color != "" {
				itemNames = append(itemNames, cartItem.Color)
			}

			priceItem, err := stripePrice.New(&stripe.PriceParams{
				Currency:   stripe.String(string(purchaseOrder.Currency)),
				UnitAmount: stripe.Int64(cartItem.UnitPrice.MultipleInt(stripeConfig.SmallestUnitFactor).ToInt64()),
				ProductData: &stripe.PriceProductDataParams{
					Name: stripe.String(strings.Join(itemNames, "/")),
				},
			})
			if err != nil {
				return nil, err
			}

			paymentLineItems = append(paymentLineItems, &stripe.PaymentLinkLineItemParams{
				Quantity: stripe.Int64(cartItem.Quantity),
				Price:    &priceItem.ID,
			})
		}

		purchaseOrders = append(purchaseOrders, &purchaseOrder)
		purchaseOrderIDs = append(purchaseOrderIDs, purchaseOrder.ID)
	}

	// Append for tax and shipping
	if totalTax.GreaterThan(0) {
		priceItemTax, err := stripePrice.New(&stripe.PriceParams{
			Currency:   stripe.String(string(currency)),
			UnitAmount: stripe.Int64(totalTax.MultipleInt(stripeConfig.SmallestUnitFactor).ToInt64()),
			ProductData: &stripe.PriceProductDataParams{
				Name: stripe.String("Tax"),
			},
		})
		if err != nil {
			return nil, err
		}
		paymentLineItems = append(paymentLineItems, &stripe.PaymentLinkLineItemParams{
			Quantity: stripe.Int64(1),
			Price:    &priceItemTax.ID,
		})
	}

	if totalShippingFee.GreaterThan(0) {
		priceItemShippingFee, err := stripePrice.New(&stripe.PriceParams{
			Currency:   stripe.String(string(currency)),
			UnitAmount: stripe.Int64(totalShippingFee.MultipleInt(stripeConfig.SmallestUnitFactor).ToInt64()),
			ProductData: &stripe.PriceProductDataParams{
				Name: stripe.String("Shipping Fee"),
			},
		})
		if err != nil {
			return nil, err
		}
		paymentLineItems = append(paymentLineItems, &stripe.PaymentLinkLineItemParams{
			Quantity: stripe.Int64(1),
			Price:    &priceItemShippingFee.ID,
		})
	}

	if totalTransactionFee.GreaterThan(0) {
		priceItemTransactionFee, err := stripePrice.New(&stripe.PriceParams{
			Currency:   stripe.String(string(currency)),
			UnitAmount: stripe.Int64(totalTransactionFee.MultipleInt(stripeConfig.SmallestUnitFactor).ToInt64()),
			ProductData: &stripe.PriceProductDataParams{
				Name: stripe.String("Transaction Fee"),
			},
		})
		if err != nil {
			return nil, err
		}
		paymentLineItems = append(paymentLineItems, &stripe.PaymentLinkLineItemParams{
			Quantity: stripe.Int64(1),
			Price:    &priceItemTransactionFee.ID,
		})
	}

	var stripeParams = stripehelper.CreatePaymentLinkParams{
		Currency: currency,
		Metadata: map[string]string{
			"checkout_session_id": checkoutSessionID,
			"action_source":       string(stripehelper.ActionSourceMultiPOPayment),
		},
		RedirectURL: fmt.Sprintf("%s/purchase-order-checkout?checkout_session_id=%s", r.db.Configuration.BrandPortalBaseURL, checkoutSessionID),
		LineItems:   paymentLineItems,
	}

	pl, err := stripehelper.GetInstance().CreatePaymentLink(stripeParams)
	if err != nil {
		return nil, err
	}

	var samplePaymentLink = helper.AddURLQuery(pl.URL,
		map[string]string{
			"client_reference_id": checkoutSessionID,
			"prefilled_email":     buyer.Email,
			"checkout_session_id": checkoutSessionID,
		},
	)
	var updates = models.PurchaseOrder{
		CheckoutSessionID: checkoutSessionID,
		PaymentLink:       samplePaymentLink,
		PaymentLinkID:     pl.ID,
	}
	err = r.db.Model(&models.PurchaseOrder{}).Where("id IN ?", purchaseOrderIDs).Updates(&updates).Error
	if err != nil {
		return nil, eris.Wrap(err, err.Error())
	}

	return &CreateMultiPurchaseOrderPaymentLinkResponse{PurchaseOrders: purchaseOrders, PaymentLink: samplePaymentLink}, nil
}

type AdminConfirmPOParams struct {
	models.JwtClaimsInfo

	PurchaseOrderID string      `json:"purchase_order_id" param:"purchase_order_id" validate:"required"`
	ShippingFee     price.Price `json:"shipping_fee" validate:"required"`
	TaxPercentage   float64     `json:"tax_percentage" validate:"gte=0"`
}

func (r *PurchaseOrderRepo) AdminConfirmPO(params AdminConfirmPOParams) (*models.PurchaseOrder, error) {
	var po models.PurchaseOrder
	var err = r.db.First(&po, "id = ?", params.PurchaseOrderID).Error
	if err != nil {
		return nil, err
	}

	var updates = models.PurchaseOrder{
		ConfirmedAt: values.Int64(time.Now().Unix()),
	}
	updates.ShippingFee = &params.ShippingFee
	updates.TaxPercentage = &params.TaxPercentage

	var tracking = models.PurchaseOrderTracking{
		PurchaseOrderID: params.PurchaseOrderID,
		ActionType:      enums.PoTrackingActionAdminConfirm,
		UserID:          po.UserID,
		CreatedByUserID: params.JwtClaimsInfo.GetUserID(),
		UserGroup:       enums.PoTrackingUserGroupBuyer,
		Metadata: &models.PoTrackingMetadata{
			Before: map[string]interface{}{
				"shipping_fee":   po.ShippingFee,
				"tax_percentage": po.TaxPercentage,
			},
			After: map[string]interface{}{
				"shipping_fee":   updates.ShippingFee,
				"tax_percentage": updates.TaxPercentage,
			},
		},
	}

	err = r.db.Transaction(func(tx *gorm.DB) error {
		err = tx.Model(&models.PurchaseOrder{}).Where("id = ?", params.PurchaseOrderID).Updates(&updates).Error
		if err != nil {
			return err
		}

		return tx.Create(&tracking).Error
	})
	if err != nil {
		return nil, err
	}

	po.ConfirmedAt = updates.ConfirmedAt
	po.ShippingFee = updates.ShippingFee
	po.TaxPercentage = updates.TaxPercentage

	return &po, err
}

type AdminCancelPOParams struct {
	models.JwtClaimsInfo

	PurchaseOrderID string `json:"purchase_order_id" param:"purchase_order_id" validate:"required"`
	Reason          string `json:"reason" validate:"required"`
}

func (r *PurchaseOrderRepo) AdminCancelPO(params AdminCancelPOParams) (*models.PurchaseOrder, error) {
	var po models.PurchaseOrder
	var err = r.db.First(&po, "id = ?", params.PurchaseOrderID).Error
	if err != nil {
		return nil, err
	}

	var canCancel = po.Status == enums.PurchaseOrderStatusPending || po.Status == enums.PurchaseOrderStatusWaitingConfirm
	if !canCancel {
		return nil, errs.ErrPOIsNotAbleToCancel.WithMessagef("PO status %s is not able to cancel", po.Status)
	}

	var tracking = models.PurchaseOrderTracking{
		ActionType:      enums.PoTrackingActionAdminCancel,
		FromStatus:      po.TrackingStatus,
		ToStatus:        enums.PoTrackingStatusCanceled,
		UserID:          po.UserID,
		PurchaseOrderID: po.ID,
		CreatedByUserID: params.GetUserID(),
		UserGroup:       enums.PoTrackingUserGroupBuyer,
		Description:     params.Reason,
	}

	var updates = models.PurchaseOrder{
		Status: enums.PurchaseOrderStatusCanceled,
	}
	err = r.db.Transaction(func(tx *gorm.DB) error {
		err = tx.Create(&tracking).Error
		if err != nil {
			return err
		}
		err = tx.Model(&models.PurchaseOrder{}).
			Where("id = ?", params.PurchaseOrderID).
			Updates(&updates).Error

		return err
	})

	return &po, err
}

type MultiPurchaseOrderParams struct {
	models.JwtClaimsInfo

	CheckoutSessionID string `param:"checkout_session_id" validate:"required"`
	Note              string `json:"note"`
}

func (r *PurchaseOrderRepo) MultiPurchaseOrderMarkAsPaid(params MultiPurchaseOrderParams) ([]*models.PurchaseOrder, error) {
	var purchaseOrders []*models.PurchaseOrder

	if params.CheckoutSessionID != "" {
		var err = r.db.Select("ID", "InquiryID", "Status", "UserID", "AssigneeIDs", "CheckoutSessionID").Find(&purchaseOrders, "checkout_session_id = ?", params.CheckoutSessionID).Error
		if err != nil {
			return nil, err
		}
	}

	if len(purchaseOrders) == 0 {
		return nil, errs.ErrInquiryNotFound
	}

	var paymentTransactionUpdates = models.PaymentTransaction{
		MarkAsPaidAt: values.Int64(time.Now().Unix()),
		Status:       enums.PaymentStatusPaid,
	}

	var err = r.db.Transaction(func(tx *gorm.DB) error {
		for _, purchaseOrder := range purchaseOrders {
			if purchaseOrder.Status == enums.PurchaseOrderStatusPaid {
				return errs.ErrPoIsAlreadyPaid
			}

			var updates = models.PurchaseOrder{
				MarkAsPaidAt: paymentTransactionUpdates.MarkAsPaidAt,
				Status:       enums.PurchaseOrderStatusPaid,
			}

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

			var tracking = models.PurchaseOrderTracking{
				PurchaseOrderID: purchaseOrder.ID,
				ActionType:      enums.PoTrackingActionSellerPaymentReceived,
				UserID:          purchaseOrder.UserID,
				CreatedByUserID: params.JwtClaimsInfo.GetUserID(),
				UserGroup:       enums.PoTrackingUserGroupBuyer,
				Metadata: &models.PoTrackingMetadata{
					After: map[string]interface{}{
						"note":                   params.Note,
						"transaction_attachment": purchaseOrder.TransactionAttachment,
						"transaction_ref_id":     purchaseOrder.TransactionRefID,
					},
				},
			}

			var err = tx.Create(&tracking).Error
			if err != nil {
				return err
			}

			err = tx.Model(&updates).Clauses(clause.Returning{Columns: []clause.Column{{Name: "id"}}}).
				Where("checkout_session_id = ? AND payment_type = ?", params.CheckoutSessionID, enums.PaymentTypeBankTransfer).
				Updates(&updates).Error
			if err != nil {
				return err
			}
		}

		return tx.Model(&models.PaymentTransaction{}).
			Where("checkout_session_id = ?", params.CheckoutSessionID).
			Updates(&paymentTransactionUpdates).Error
	})

	return purchaseOrders, err
}

func (r *PurchaseOrderRepo) MultiPurchaseOrderMarkAsUnpaid(params MultiPurchaseOrderParams) ([]*models.PurchaseOrder, error) {
	var purchaseOrders []*models.PurchaseOrder

	if params.CheckoutSessionID != "" {
		var err = r.db.Select("ID", "InquiryID", "Status", "UserID", "AssigneeIDs", "CheckoutSessionID").Find(&purchaseOrders, "checkout_session_id = ?", params.CheckoutSessionID).Error
		if err != nil {
			return nil, err
		}
	}

	if len(purchaseOrders) == 0 {
		return nil, errs.ErrInquiryNotFound
	}

	for _, purchaseOrder := range purchaseOrders {
		if purchaseOrder.Status == enums.PurchaseOrderStatusPaid {
			return nil, errs.ErrPoIsAlreadyPaid
		}
	}

	var err = r.db.Transaction(func(tx *gorm.DB) (e error) {
		var updates = models.PurchaseOrder{
			MarkAsUnpaidAt: values.Int64(time.Now().Unix()),
			Status:         enums.PurchaseOrderStatusUnpaid,
		}
		e = tx.Model(&updates).Clauses(clause.Returning{Columns: []clause.Column{{Name: "id"}}}).
			Where("checkout_session_id = ? AND payment_type = ?", params.CheckoutSessionID, enums.PaymentTypeBankTransfer).
			Updates(&updates).Error
		if e != nil {
			return
		}

		var paymentTrxUpdates = models.PaymentTransaction{
			Status:         enums.PaymentStatusUnpaid,
			Note:           params.Note,
			MarkAsUnpaidAt: updates.MarkAsUnpaidAt,
		}
		e = tx.Model(&paymentTrxUpdates).
			Where("checkout_session_id = ?", updates.ID).
			Updates(&paymentTrxUpdates).Error
		return
	})

	return purchaseOrders, err
}

func (r *PurchaseOrderRepo) PurchaseOrderMarkAsPaid(params models.PurchaseOrderIDParam) (*models.PurchaseOrder, error) {
	cancel, err := r.db.Locker.AcquireLock(fmt.Sprintf("purchase_order_%s", params.PurchaseOrderID), time.Second*20)
	if err != nil {
		return nil, err
	}
	defer cancel()

	var purchaseOrder = params.PurchaseOrder

	if params.PurchaseOrderID != "" && purchaseOrder == nil {
		purchaseOrder, err = r.GetPurchaseOrderShortInfo(params.PurchaseOrderID)
		if err != nil {
			return nil, err
		}
	}

	if purchaseOrder.ID == "" {
		return nil, errs.ErrInquiryNotFound
	}

	if purchaseOrder.Status == enums.PurchaseOrderStatusPaid {
		return nil, errs.ErrPoIsAlreadyPaid
	}

	var updates = models.PurchaseOrder{
		MarkAsPaidAt: values.Int64(time.Now().Unix()),
		Status:       enums.PurchaseOrderStatusPaid,
	}

	var paymentTransactionUpdates = models.PaymentTransaction{
		MarkAsPaidAt: updates.MarkAsPaidAt,
		Status:       enums.PaymentStatusPaid,
	}

	var tracking = models.PurchaseOrderTracking{
		PurchaseOrderID: params.PurchaseOrderID,
		ActionType:      enums.PoTrackingActionPaymentReceived,
		UserID:          purchaseOrder.UserID,
		CreatedByUserID: params.JwtClaimsInfo.GetUserID(),
		UserGroup:       enums.PoTrackingUserGroupBuyer,
		Metadata: &models.PoTrackingMetadata{
			After: map[string]interface{}{
				"note":                   params.Note,
				"transaction_attachment": purchaseOrder.TransactionAttachment,
				"transaction_ref_id":     purchaseOrder.TransactionRefID,
			},
		},
	}

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

	err = r.db.Transaction(func(tx *gorm.DB) error {
		var err = tx.Model(&updates).Clauses(clause.Returning{Columns: []clause.Column{{Name: "id"}}}).
			Where("purchase_order_id = ? AND payment_type = ?", params.PurchaseOrderID, enums.PaymentTypeBankTransfer).
			Updates(&updates).Error
		if err != nil {
			return err
		}

		err = tx.Create(&tracking).Error
		if err != nil {
			return err
		}

		return tx.Model(&models.PaymentTransaction{}).
			Where("purchase_order_id = ?", updates.ID).
			Updates(&paymentTransactionUpdates).Error
	})
	if err != nil {
		return nil, err
	}

	purchaseOrder.Status = updates.Status
	purchaseOrder.MarkAsPaidAt = updates.MarkAsPaidAt

	return purchaseOrder, err
}

func (r *PurchaseOrderRepo) PurchaseOrderMarkAsUnpaid(params models.PurchaseOrderIDParam) (*models.PurchaseOrder, error) {
	cancel, err := r.db.Locker.AcquireLock(fmt.Sprintf("purchase_order_%s", params.PurchaseOrderID), time.Second*20)
	if err != nil {
		return nil, err
	}
	defer cancel()

	var purchaseOrder = params.PurchaseOrder

	if params.PurchaseOrderID != "" && purchaseOrder == nil {
		purchaseOrder, err = r.GetPurchaseOrderShortInfo(params.PurchaseOrderID)
		if err != nil {
			return nil, err
		}
	}

	if purchaseOrder.ID == "" {
		return nil, errs.ErrInquiryNotFound
	}

	if purchaseOrder.Status == enums.PurchaseOrderStatusPaid {
		return nil, errs.ErrPoIsAlreadyPaid
	}

	var updates = models.PurchaseOrder{
		MarkAsUnpaidAt: values.Int64(time.Now().Unix()),
		Status:         enums.PurchaseOrderStatusUnpaid,
	}

	err = r.db.Transaction(func(tx *gorm.DB) (e error) {
		e = tx.Model(&updates).Clauses(clause.Returning{Columns: []clause.Column{{Name: "id"}}}).
			Where("purchase_order_id = ? AND payment_type = ?", params.PurchaseOrderID, enums.PaymentTypeBankTransfer).
			Updates(&updates).Error
		if e != nil {
			return
		}

		var paymentTrxUpdates = models.PaymentTransaction{
			Status:         enums.PaymentStatusUnpaid,
			Note:           params.Note,
			MarkAsUnpaidAt: updates.MarkAsUnpaidAt,
		}
		e = tx.Model(&paymentTrxUpdates).
			Where("purchase_order_id = ?", updates.ID).
			Updates(&paymentTrxUpdates).Error
		return
	})
	if err != nil {
		return nil, err
	}

	purchaseOrder.Status = updates.Status
	purchaseOrder.MarkAsUnpaidAt = updates.MarkAsUnpaidAt
	return purchaseOrder, err
}

type GetPurchaseOrderInvoiceParams struct {
	models.JwtClaimsInfo

	PurchaseOrderID string `param:"purchase_order_id" json:"purchase_order_id" validate:"required"`
	ReCreate        bool   `param:"re_create"`
}

func (r *PurchaseOrderRepo) GetPurchaseOrderInvoice(params GetPurchaseOrderInvoiceParams) (*models.Attachment, error) {
	var purchaseOrder models.PurchaseOrder
	var err = r.db.Select("ID", "InvoiceNumber").First(&purchaseOrder, "id = ?", params.PurchaseOrderID).Error
	if err != nil {
		return nil, err
	}

	if purchaseOrder.InvoiceNumber > 0 {
		var invoice models.Invoice
		if err := r.db.First(&invoice, "invoice_number = ?", purchaseOrder.InvoiceNumber).Error; err == nil && invoice.Document != nil {
			return invoice.Document, err
		}
	}

	result, err := NewInvoiceRepo(r.db).CreatePurchaseOrderInvoice(CreatePurchaseOrderInvoiceParams{
		PurchaseOrderID: purchaseOrder.ID,
		ReCreate:        params.ReCreate,
	})
	if err != nil {
		if eris.Is(err, errs.ErrPOInvoiceAlreadyGenerated) && result.Invoice.Document != nil {
			goto Susccess
		}
		return nil, err
	}

	if result.Invoice == nil || result.Invoice.Document == nil {
		return nil, eris.New("Invoice is not able to generate")
	}

Susccess:
	return result.Invoice.Document, nil
}

type UpdatePurchaseOrderLogsParams struct {
	models.JwtClaimsInfo

	PurchaseOrderID string `json:"purchase_order_id" param:"purchase_order_id" validate:"required"`
	LogID           string `json:"log_id" param:"log_id" validate:"required"`

	Notes       string             `json:"notes" validate:"required"`
	Attachments models.Attachments `json:"attachments" validate:"required"`
}

func (r *PurchaseOrderRepo) UpdatePurchaseOrderLogs(params UpdatePurchaseOrderLogsParams) (*models.PurchaseOrderTracking, error) {
	var log = models.PurchaseOrderTracking{
		Notes:       params.Notes,
		Attachments: &params.Attachments,
	}

	var err = r.db.Model(&models.PurchaseOrderTracking{}).
		Where("purchase_order_id = ? AND id = ?", params.PurchaseOrderID, params.LogID).
		Updates(&log).Error

	return &log, err
}

type DeletePurchaseOrderLogsParams struct {
	models.JwtClaimsInfo

	PurchaseOrderID string `json:"purchase_order_id" param:"purchase_order_id" validate:"required"`
	LogID           string `json:"log_id" param:"log_id" validate:"required"`
}

func (r *PurchaseOrderRepo) DeletePurchaseOrderLogs(params DeletePurchaseOrderLogsParams) (*models.PurchaseOrderTracking, error) {
	var log = models.PurchaseOrderTracking{
		Notes:       "",
		Attachments: nil,
	}

	var err = r.db.Model(&models.PurchaseOrderTracking{}).
		Select("Notes", "Attachments").
		Where("purchase_order_id = ? AND id = ?", params.PurchaseOrderID, params.LogID).
		Updates(&log).Error

	return &log, err
}

type AdminSellerApproveRawMaterialsParmas struct {
	models.JwtClaimsInfo

	PurchaseOrderID string `json:"purchase_order_id" param:"purchase_order_id" query:"purchase_order_id" validate:"required"`
}

func (r *PurchaseOrderRepo) AdminSellerApproveRawMaterials(params AdminSellerApproveRawMaterialsParmas) (*models.PurchaseOrder, error) {
	order, err := r.GetPurchaseOrder(GetPurchaseOrderParams{
		PurchaseOrderID: params.PurchaseOrderID,
		JwtClaimsInfo:   params.JwtClaimsInfo,
	})

	if err != nil {
		return nil, err
	}

	validStatus := []enums.SellerPoTrackingStatus{enums.SellerPoTrackingStatusRawMaterial}
	if ok := slices.Contains(validStatus, order.SellerTrackingStatus); !ok {
		return nil, errs.ErrPoInvalidToMarkMaking
	}

	var updates models.PurchaseOrder
	err = copier.Copy(&updates, &params)
	if err != nil {
		return nil, err
	}

	if order.SellerPoRawMaterials != nil {
		var items models.PoRawMaterialMetas = lo.Map(*order.SellerPoRawMaterials, func(item *models.PoRawMaterialMeta, index int) *models.PoRawMaterialMeta {
			item.Status = enums.PoRawMaterialStatusApproved
			return item
		})
		updates.SellerPoRawMaterials = &items
		updates.PoRawMaterials = &items
	}

	if order.TrackingStatus != updates.TrackingStatus {
		err = r.db.Transaction(func(tx *gorm.DB) error {
			err = NewPurchaseOrderTrackingRepo(r.db).CreatePurchaseOrderTrackingTx(tx, models.PurchaseOrderTrackingCreateForm{
				PurchaseOrderID: params.PurchaseOrderID,
				ActionType:      enums.PoTrackingActionMarkMaking,
				UserID:          order.UserID,
				CreatedByUserID: params.JwtClaimsInfo.GetUserID(),
				Metadata: &models.PoTrackingMetadata{
					Before: map[string]interface{}{
						"tracking_status": order.TrackingStatus,
					},
					After: map[string]interface{}{
						"tracking_status":  updates.TrackingStatus,
						"po_raw_materials": updates.PoRawMaterials,
					},
				},
			})

			if err != nil {
				return eris.Wrap(err, err.Error())
			}

			return tx.Model(&models.PurchaseOrder{}).
				Where("id = ?", params.PurchaseOrderID).
				Updates(&updates).Error
		})
		if err != nil {
			return nil, eris.Wrap(err, err.Error())
		}
	}

	order.MakingInfo = updates.MakingInfo
	order.TrackingStatus = updates.TrackingStatus
	order.SellerTrackingStatus = updates.SellerTrackingStatus

	return order, err
}

type BuyerMakeAnotherSampleParams struct {
	models.JwtClaimsInfo

	PurchaseOrderID string `json:"purchase_order_id" param:"purchase_order_id" validate:"required"`

	Requirement string `json:"requirement"`
	TrimButton  string `json:"trim_button"`
	TrimThread  string `json:"trim_thread"`
	TrimZipper  string `json:"trim_zipper"`
	TrimLabel   string `json:"trim_label"`

	FabricName   string   `json:"fabric_name"`
	FabricWeight *float64 `json:"fabric_weight"`

	Attachments         *models.Attachments `json:"attachments,omitempty"`
	FabricAttachments   *models.Attachments `json:"fabric_attachments,omitempty"`
	TechpackAttachments *models.Attachments `json:"techpack_attachments,omitempty"`

	ShippingAddress *models.AddressForm `gorm:"-" json:"shipping_address,omitempty"`

	Items []*models.InquiryCartItemCreateForm `json:"items"`
}

func (r *PurchaseOrderRepo) BuyerMakeAnotherSample(params BuyerMakeAnotherSampleParams) (*models.PurchaseOrder, error) {
	var purchaseOrder models.PurchaseOrder
	var err = r.db.First(&purchaseOrder, "id = ? AND user_id = ?", params.PurchaseOrderID, params.GetUserID()).Error
	if err != nil {
		if r.db.IsRecordNotFoundError(err) {
			return nil, errs.ErrPONotFound
		}
		return nil, err
	}

	if *purchaseOrder.RoundID == "" {
		*purchaseOrder.RoundID = helper.GenerateSampleRoundID()
	}

	var inquiry models.Inquiry
	err = r.db.First(&inquiry, "id = ? AND user_id = ?", purchaseOrder.InquiryID, params.GetUserID()).Error
	if err != nil {
		if r.db.IsRecordNotFoundError(err) {
			return nil, errs.ErrPONotFound
		}
		return nil, err
	}

	var itemIDs []string
	var items = lo.Map(params.Items, func(item *models.InquiryCartItemCreateForm, index int) *models.InquiryCartItem {
		var cartItem models.InquiryCartItem
		var err = copier.Copy(&cartItem, item)
		if err != nil {
			return nil
		}

		cartItem.ID = helper.GenerateXID()
		cartItem.InquiryID = purchaseOrder.InquiryID
		cartItem.PurchaseOrderID = purchaseOrder.ID
		cartItem.UnitPrice = inquiry.GetSampleUnitPrice()
		cartItem.TotalPrice = item.UnitPrice.MultipleInt(int64(item.Qty))

		itemIDs = append(itemIDs, cartItem.ID)
		return &cartItem
	})

	var newInquiry = models.Inquiry{
		AdminQuotations:     inquiry.AdminQuotations,
		UserID:              inquiry.UserID,
		AssigneeIDs:         inquiry.AssigneeIDs,
		Title:               inquiry.Title,
		SkuNote:             inquiry.SkuNote,
		Currency:            inquiry.Currency,
		Quantity:            inquiry.Quantity,
		ExpiredDate:         inquiry.ExpiredDate,
		TaxPercentage:       inquiry.TaxPercentage,
		ShippingAddressID:   inquiry.ShippingAddressID,
		ShippingAddress:     inquiry.ShippingAddress,
		SizeList:            inquiry.SizeList,
		SizeChart:           inquiry.SizeChart,
		Composition:         inquiry.Composition,
		ColorList:           inquiry.ColorList,
		StyleNo:             inquiry.StyleNo,
		FabricWeightUnit:    inquiry.FabricWeightUnit,
		CollectionID:        inquiry.CollectionID,
		OrderGroupID:        inquiry.OrderGroupID,
		Requirement:         params.Requirement,
		Attachments:         params.Attachments,
		FabricAttachments:   params.FabricAttachments,
		TechpackAttachments: params.TechpackAttachments,
		FabricName:          params.FabricName,
		FabricWeight:        params.FabricWeight,
	}
	newInquiry.ID = helper.GenerateXID()

	var newPurchaseOrder = models.PurchaseOrder{
		RoundID:             purchaseOrder.RoundID,
		Currency:            purchaseOrder.Currency,
		ShippingAddressID:   purchaseOrder.ShippingAddressID,
		ShippingAddress:     purchaseOrder.ShippingAddress,
		UserID:              purchaseOrder.UserID,
		AssigneeIDs:         purchaseOrder.AssigneeIDs,
		CartItemIDs:         itemIDs,
		Attachments:         params.Attachments,
		FabricAttachments:   params.FabricAttachments,
		TechpackAttachments: params.TechpackAttachments,
		InquiryID:           newInquiry.ID,
	}

	if params.ShippingAddress != nil {
		var address models.Address
		err = copier.Copy(&address, params.ShippingAddress)
		if err != nil {
			return nil, err
		}
		address.UserID = inquiry.UserID

		err = address.CreateOrUpdate(r.db)
		if err != nil {
			return nil, err
		}

		newInquiry.ShippingAddress = &address
		newInquiry.ShippingAddressID = address.ID

		newPurchaseOrder.ShippingAddress = &address
		newPurchaseOrder.ShippingAddressID = address.ID
	}

	err = r.db.Transaction(func(tx *gorm.DB) error {
		err = tx.Create(&items).Error
		if err != nil {
			return err
		}

		err = tx.Create(&newPurchaseOrder).Error
		if err != nil {
			return err
		}

		err = tx.Create(&newInquiry).Error
		if err != nil {
			return err
		}

		return tx.Model(&models.PurchaseOrder{}).Where("id = ?", purchaseOrder.ID).UpdateColumn("RoundID", purchaseOrder.RoundID).Error
	})
	if err != nil {
		return nil, err
	}

	return NewInquiryRepo(r.db).InquiryPreviewCheckout(InquiryPreviewCheckoutParams{
		InquiryID:     newInquiry.ID,
		UserID:        newInquiry.UserID,
		Inquiry:       &newInquiry,
		CartItems:     items,
		UpdatePricing: true,
		JwtClaimsInfo: params.JwtClaimsInfo,
	})

}

type AdminPurchaseOrderApproveRoundParams struct {
	models.JwtClaimsInfo
	PurchaseOrderID string `json:"purchase_order_id" param:"purchase_order_id" validate:"required"`
	RoundID         string `json:"round_id" param:"round_id" validate:"required"`
}

func (r *PurchaseOrderRepo) AdminPurchaseOrderApproveRound(params AdminPurchaseOrderApproveRoundParams) error {
	var sqlResult = r.db.Model(&models.PurchaseOrder{}).
		Where("id = ? AND round_id = ?", params.PurchaseOrderID, params.RoundID).UpdateColumn("RoundStatus", enums.RoundStatusApproved)

	if sqlResult.Error != nil {
		return sqlResult.Error
	}

	if sqlResult.RowsAffected == 0 {
		return errs.ErrPONotFound
	}
	return nil
}

type AdminPurchaseOrderRejectRoundParams struct {
	models.JwtClaimsInfo
	PurchaseOrderID string `json:"purchase_order_id" param:"purchase_order_id" validate:"required"`
	RoundID         string `json:"round_id" param:"round_id" validate:"required"`
}

func (r *PurchaseOrderRepo) AdminPurchaseOrderRejectRound(params AdminPurchaseOrderRejectRoundParams) error {
	var sqlResult = r.db.Model(&models.PurchaseOrder{}).
		Where("id = ? AND round_id = ?", params.PurchaseOrderID, params.RoundID).UpdateColumn("RoundStatus", enums.RoundStatusRejected)

	if sqlResult.Error != nil {
		return sqlResult.Error
	}

	if sqlResult.RowsAffected == 0 {
		return errs.ErrPONotFound
	}
	return nil
}
