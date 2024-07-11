package repo

import (
	"time"

	"github.com/engineeringinflow/inflow-backend/pkg/db"
	"github.com/engineeringinflow/inflow-backend/pkg/errs"
	"github.com/engineeringinflow/inflow-backend/pkg/models"
	"github.com/engineeringinflow/inflow-backend/pkg/models/enums"
	"github.com/engineeringinflow/inflow-backend/pkg/models/price"
	"github.com/jinzhu/copier"
	"github.com/rotisserie/eris"
	"github.com/samber/lo"
	"github.com/thaitanloi365/go-utils/values"
	"golang.org/x/exp/slices"
	"gorm.io/gorm"
)

type SellerPurchaseOrderRepo struct {
	db *db.DB
}

func NewSellerPurchaseOrderRepo(db *db.DB) *SellerPurchaseOrderRepo {
	return &SellerPurchaseOrderRepo{
		db: db,
	}
}

type AdminSellerPORawMaterialSendToBuyerParams struct {
	models.JwtClaimsInfo

	PurchaseOrderID string `json:"purchase_order_id" param:"purchase_order_id" query:"purchase_order_id" validate:"required"`
}

func (r *SellerPurchaseOrderRepo) AdminSellerPORawMaterialSendToBuyer(params AdminSellerPORawMaterialSendToBuyerParams) (*models.PurchaseOrder, error) {
	order, err := NewPurchaseOrderRepo(r.db).GetPurchaseOrder(GetPurchaseOrderParams{
		PurchaseOrderID: params.PurchaseOrderID,
		JwtClaimsInfo:   params.JwtClaimsInfo,
	})

	if err != nil {
		return nil, err
	}

	var updates models.PurchaseOrder
	updates.PoRawMaterials = order.PoRawMaterials

	validStatus := []enums.SellerPoTrackingStatus{enums.SellerPoTrackingStatusRawMaterial}
	if ok := slices.Contains(validStatus, order.SellerTrackingStatus); !ok {
		return order, errs.ErrPOInvalidToSendRawMaterial
	}

	for _, rawItem := range *updates.PoRawMaterials {
		rawItem.WaitingForSendToBuyer = values.Bool(false)
	}

	err = r.db.Model(&models.PurchaseOrder{}).
		Where("id = ?", params.PurchaseOrderID).
		Updates(&updates).Error

	if err != nil {
		return nil, eris.Wrap(err, err.Error())
	}

	order.PoRawMaterials = updates.PoRawMaterials
	return order, err
}

func (r *SellerPurchaseOrderRepo) AdminSellerPoRawMaterialCommentMarkSeen(params PoCommentMarkSeenParams) error {
	var err = r.db.Model(&models.Comment{}).
		Where("seen_at IS NULL").
		Where("user_id != ?", params.GetUserID()).
		Where("target_type = ? AND target_id = ? AND file_key = ?", enums.CommentTargetTypeSellerPoRawMaterial, params.PurchaseOrderID, params.FileKey).
		Update("seen_at", time.Now().Unix()).Error

	return err
}

func (r *SellerPurchaseOrderRepo) AdminSellerPoRawMaterialCommentStatusCount(params PoCommentStatusCountParams) ([]*models.CommentStatusCountItem, error) {
	order, err := NewPurchaseOrderRepo(r.db).GetPurchaseOrder(GetPurchaseOrderParams{
		PurchaseOrderID: params.PurchaseOrderID,
		JwtClaimsInfo:   params.JwtClaimsInfo,
	})
	if err != nil {
		return nil, eris.Wrap(err, err.Error())
	}

	var items []*models.CommentStatusCountItem

	if order.PoRawMaterials != nil && len(*order.PoRawMaterials) > 0 {
		for _, fileItem := range *order.PoRawMaterials {
			count := int64(0)
			if fileItem.ReferenceID != "" {
				r.db.Model(&models.Comment{}).Where("user_id != ? AND target_type = ? AND target_id = ? AND file_key = ? AND seen_at IS NULL", params.GetUserID(), enums.CommentTargetTypeSellerPoRawMaterial, params.PurchaseOrderID, fileItem.ReferenceID).Count(&count)
				items = append(items, &models.CommentStatusCountItem{
					FileKey:     fileItem.ReferenceID,
					UnseenCount: count,
				})
			}

		}
	}

	return items, nil
}

func (r *SellerPurchaseOrderRepo) SellerPoRawMaterialCommentStatusCount(params PoCommentStatusCountParams) ([]*models.CommentStatusCountItem, error) {
	order, err := NewPurchaseOrderRepo(r.db).GetPurchaseOrder(GetPurchaseOrderParams{
		PurchaseOrderID: params.PurchaseOrderID,
		JwtClaimsInfo:   params.JwtClaimsInfo,
	})
	if err != nil {
		return nil, eris.Wrap(err, err.Error())
	}

	var items []*models.CommentStatusCountItem

	if order.PoRawMaterials != nil && len(*order.PoRawMaterials) > 0 {
		for _, fileItem := range *order.PoRawMaterials {
			count := int64(0)
			if fileItem.ReferenceID != "" {
				r.db.Model(&models.Comment{}).Where("user_id != ? AND target_type = ? AND target_id = ? AND file_key = ? AND seen_at IS NULL", params.GetUserID(), enums.CommentTargetTypeSellerPoRawMaterial, params.PurchaseOrderID, fileItem.ReferenceID).Count(&count)
				items = append(items, &models.CommentStatusCountItem{
					FileKey:     fileItem.ReferenceID,
					UnseenCount: count,
				})
			}

		}
	}

	return items, nil
}

type SellerPurchaseOrderMarkMakingParams struct {
	models.JwtClaimsInfo

	PurchaseOrderID     string             `json:"purchase_order_id" param:"purchase_order_id"`
	SellerEstMakingAt   *int64             `json:"seller_est_making_at,omitempty"`
	SellerEstDeliveryAt *int64             `json:"seller_est_delivery_at,omitempty"`
	Note                string             `json:"note"`
	Attachments         models.Attachments `json:"attachments"`
}

func (r *SellerPurchaseOrderRepo) SellerPurchaseOrderMarkMaking(params SellerPurchaseOrderMarkMakingParams) error {
	order, err := NewPurchaseOrderRepo(r.db).GetPurchaseOrder(GetPurchaseOrderParams{
		PurchaseOrderID: params.PurchaseOrderID,
		JwtClaimsInfo:   params.JwtClaimsInfo,
	})
	if err != nil {
		return err
	}

	var validStatus = []enums.SellerPoTrackingStatus{
		enums.SellerPoTrackingStatusRawMaterial,
		enums.SellerPoTrackingStatusRawMaterialSkipped,
		enums.SellerPoTrackingStatusMaking,
	}
	if !slices.Contains(validStatus, order.SellerTrackingStatus) {
		return errs.ErrPoInvalidToMarkMaking
	}

	var updates = models.PurchaseOrder{
		SellerEstMakingAt:    params.SellerEstMakingAt,
		SellerEstDeliveryAt:  params.SellerEstDeliveryAt,
		SellerTrackingStatus: enums.SellerPoTrackingStatusMaking,
	}

	if order.SellerTrackingStatus != updates.SellerTrackingStatus {
		err = r.db.Transaction(func(tx *gorm.DB) error {
			err = NewPurchaseOrderTrackingRepo(r.db).CreatePurchaseOrderTrackingTx(tx, models.PurchaseOrderTrackingCreateForm{
				PurchaseOrderID: params.PurchaseOrderID,
				ActionType:      enums.PoTrackingActionSellerMarkMaking,
				UserID:          params.GetUserID(),
				CreatedByUserID: params.GetUserID(),
				UserGroup:       enums.PoTrackingUserGroupSeller,
				Metadata: &models.PoTrackingMetadata{
					Before: map[string]interface{}{
						"seller_tracking_status": order.SellerTrackingStatus,
					},
					After: map[string]interface{}{
						"seller_tracking_status": updates.SellerTrackingStatus,
						"seller_est_making_at":   updates.SellerEstMakingAt,
						"seller_est_delivery_at": updates.SellerEstDeliveryAt,
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
			return eris.Wrap(err, err.Error())
		}
	}

	return err
}

type SellerPurchaseOrderMarkSubmitParams struct {
	models.JwtClaimsInfo

	PurchaseOrderID  string                      `json:"purchase_order_id" param:"purchase_order_id" query:"purchase_order_id" validate:"required"`
	SellerSubmitInfo *models.PoMarkingStatusMeta `json:"seller_submit_info" param:"seller_submit_info" query:"seller_submit_info"`
}

func (r *SellerPurchaseOrderRepo) SellerMarkSubmit(params SellerPurchaseOrderMarkSubmitParams) (*models.PurchaseOrder, error) {
	order, err := NewPurchaseOrderRepo(r.db).GetPurchaseOrder(GetPurchaseOrderParams{
		PurchaseOrderID: params.PurchaseOrderID,
		JwtClaimsInfo:   params.JwtClaimsInfo,
	})

	if err != nil {
		return nil, err
	}

	validStatus := []enums.SellerPoTrackingStatus{enums.SellerPoTrackingStatusMaking, enums.SellerPoTrackingStatusSubmit}
	if ok := slices.Contains(validStatus, order.SellerTrackingStatus); !ok {
		return nil, errs.ErrPoInvalidToMarkSubmit
	}

	var updates models.PurchaseOrder
	err = copier.Copy(&updates, &params)
	if err != nil {
		return nil, err
	}
	updates.SellerTrackingStatus = enums.SellerPoTrackingStatusSubmit
	updates.TrackingStatus = enums.PoTrackingStatusSubmit

	if order.SellerTrackingStatus != updates.SellerTrackingStatus {
		err = r.db.Transaction(func(tx *gorm.DB) error {
			err = NewPurchaseOrderTrackingRepo(r.db).CreatePurchaseOrderTrackingTx(tx, models.PurchaseOrderTrackingCreateForm{
				PurchaseOrderID: params.PurchaseOrderID,
				ActionType:      enums.PoTrackingActionSellerMarkSubmit,
				UserID:          params.GetUserID(),
				CreatedByUserID: params.GetUserID(),
				UserGroup:       enums.PoTrackingUserGroupSeller,
				Metadata: &models.PoTrackingMetadata{
					Before: map[string]interface{}{
						"seller_tracking_status": order.SellerTrackingStatus,
					},
					After: map[string]interface{}{
						"seller_tracking_status": updates.SellerTrackingStatus,
						"seller_submit_info":     updates.SellerSubmitInfo,
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

	order.SellerSubmitInfo = updates.SellerSubmitInfo
	order.SellerTrackingStatus = updates.SellerTrackingStatus

	return order, err
}

type SellerPurchaseOrderMarkDeliveringParams struct {
	models.JwtClaimsInfo

	PurchaseOrderID    string                 `json:"purchase_order_id" param:"purchase_order_id" query:"purchase_order_id" validate:"required"`
	SellerLogisticInfo *models.PoLogisticMeta `json:"seller_logistic_info" param:"seller_logistic_info" query:"seller_logistic_info" form:"seller_logistic_info" validate:"required"`
}

func (r *SellerPurchaseOrderRepo) SellerMarkDelivering(params SellerPurchaseOrderMarkDeliveringParams) (*models.PurchaseOrder, error) {
	order, err := NewPurchaseOrderRepo(r.db).GetPurchaseOrder(GetPurchaseOrderParams{
		PurchaseOrderID: params.PurchaseOrderID,
		JwtClaimsInfo:   params.JwtClaimsInfo,
	})

	if err != nil {
		return nil, err
	}

	validStatus := []enums.SellerPoTrackingStatus{enums.SellerPoTrackingStatusSubmit}
	if ok := slices.Contains(validStatus, order.SellerTrackingStatus); !ok {
		return nil, errs.ErrPoInvalidToMarkDelivering
	}

	var updates models.PurchaseOrder
	err = copier.Copy(&updates, &params)
	if err != nil {
		return nil, err
	}
	updates.SellerDeliveryStartedAt = values.Int64(time.Now().Unix())
	updates.SellerTrackingStatus = enums.SellerPoTrackingStatusDelivering

	if order.SellerTrackingStatus != updates.SellerTrackingStatus {
		err = r.db.Transaction(func(tx *gorm.DB) error {
			err = NewPurchaseOrderTrackingRepo(r.db).CreatePurchaseOrderTrackingTx(tx, models.PurchaseOrderTrackingCreateForm{
				PurchaseOrderID: params.PurchaseOrderID,
				ActionType:      enums.PoTrackingActionSellerMarkDelivering,
				UserID:          params.JwtClaimsInfo.GetUserID(),
				CreatedByUserID: params.JwtClaimsInfo.GetUserID(),
				UserGroup:       enums.PoTrackingUserGroupSeller,
				Metadata: &models.PoTrackingMetadata{
					Before: map[string]interface{}{
						"seller_tracking_status": order.SellerTrackingStatus,
					},
					After: map[string]interface{}{
						"seller_tracking_status": updates.SellerTrackingStatus,
						"seller_logistic_info":   params.SellerLogisticInfo,
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

	order.SellerDeliveryStartedAt = updates.SellerDeliveryStartedAt
	order.SellerTrackingStatus = updates.SellerTrackingStatus
	order.SellerLogisticInfo = updates.LogisticInfo

	return order, err
}

type SellerApproveDesignParams struct {
	models.JwtClaimsInfo

	PurchaseOrderID   string                           `json:"purchase_order_id" param:"purchase_order_id" query:"purchase_order_id" validate:"required"`
	ApproveRejectMeta *models.InquiryApproveRejectMeta `json:"approve_reject_meta" param:"approve_reject_meta" query:"approve_reject_meta"`
}

func (r *SellerPurchaseOrderRepo) SellerApproveDesign(params SellerApproveDesignParams) (*models.PurchaseOrder, error) {
	order, err := NewPurchaseOrderRepo(r.db).GetPurchaseOrder(GetPurchaseOrderParams{
		PurchaseOrderID: params.PurchaseOrderID,
		JwtClaimsInfo:   params.JwtClaimsInfo,
	})
	if err != nil {
		return nil, err
	}

	validStatus := []enums.SellerPoTrackingStatus{enums.SellerPoTrackingStatusDesignApproval}
	if ok := slices.Contains(validStatus, order.SellerTrackingStatus); !ok {
		return nil, errs.ErrPoInvalidToApproveDesign
	}

	var updates models.PurchaseOrder
	err = copier.Copy(&updates, &params)
	if err != nil {
		return nil, err
	}

	updates.SellerTrackingStatus = enums.SellerPoTrackingStatusDesignApprovedBySeller
	updates.SellerTechpackAttachments = order.TechpackAttachments
	updates.SellerDesignApprovedAt = values.Int64(time.Now().Unix())

	err = r.db.Transaction(func(tx *gorm.DB) error {
		var tracking = models.PurchaseOrderTracking{
			PurchaseOrderID: params.PurchaseOrderID,
			ActionType:      enums.PoTrackingActionSellerApprovedDesign,
			UserID:          order.SampleMakerID,
			UserGroup:       enums.PoTrackingUserGroupSeller,
			CreatedByUserID: params.JwtClaimsInfo.GetUserID(),
			Metadata: &models.PoTrackingMetadata{
				Before: map[string]interface{}{
					"seller_tracking_status": order.SellerTrackingStatus,
				},
				After: map[string]interface{}{
					"seller_tracking_status": updates.SellerTrackingStatus,
				},
			},
		}

		err = tx.Create(&tracking).Error
		if err != nil {
			return err
		}

		var sqlResult = tx.Model(&models.PurchaseOrder{}).
			Where("id = ?", params.PurchaseOrderID).
			Updates(&updates)
		if sqlResult.Error != nil {
			return sqlResult.Error
		}

		if sqlResult.RowsAffected == 0 {
			return errs.ErrPONotFound
		}

		return sqlResult.Error
	})

	if err != nil {
		return nil, eris.Wrap(err, err.Error())
	}

	order.SellerTrackingStatus = enums.SellerPoTrackingStatusDesignApprovedBySeller
	return order, err
}

type AdminSellerPurchaseOrderUploadCommentMarkSeenParams struct {
	models.JwtClaimsInfo

	PurchaseOrderID string `json:"purchase_order_id" param:"purchase_order_id"`
}

func (r *SellerPurchaseOrderRepo) AdminSellerPurchaseOrderUploadCommentMarkSeen(params AdminSellerPurchaseOrderUploadCommentMarkSeenParams) error {
	var err = r.db.Model(&models.Comment{}).
		Where("seen_at IS NULL").
		Where("user_id != ?", params.GetUserID()).
		Where("target_type = ? AND target_id = ?", enums.CommentTargetTypeSellerPoUpload, params.PurchaseOrderID).
		Update("seen_at", time.Now().Unix()).Error

	return err
}

type SellerPurchaseOrderUploadCommentMarkSeenParams struct {
	models.JwtClaimsInfo

	PurchaseOrderID string `json:"purchase_order_id" param:"purchase_order_id"`
}

func (r *SellerPurchaseOrderRepo) SellerPurchaseOrderUploadCommentMarkSeen(params SellerPurchaseOrderUploadCommentMarkSeenParams) error {
	var err = r.db.Model(&models.Comment{}).
		Where("seen_at IS NULL").
		Where("user_id != ?", params.GetUserID()).
		Where("target_type = ? AND target_id = ?", enums.CommentTargetTypeSellerPoUpload, params.PurchaseOrderID).
		Update("seen_at", time.Now().Unix()).Error

	return err
}

type SellerApprovePoParams struct {
	models.JwtClaimsInfo

	PurchaseOrderID string `json:"purchase_order_id" param:"purchase_order_id"`
}

func (r *SellerPurchaseOrderRepo) SellerApprovePo(params SellerApprovePoParams) error {
	var purchaseOrder models.PurchaseOrder
	var err = r.db.Select("ID", "SellerPoAttachments", "TrackingStatus", "SellerTrackingStatus").
		First(&purchaseOrder, "id = ? AND sample_maker_id = ?", params.PurchaseOrderID, params.GetUserID()).Error
	if err != nil {
		return err
	}
	var updates = models.PurchaseOrder{
		SellerTrackingStatus: enums.SellerPoTrackingStatusWaitingForPayment,
	}

	if purchaseOrder.SellerPoAttachments != nil {
		var filterItems = lo.Filter(*purchaseOrder.SellerPoAttachments, func(item *models.PoAttachment, index int) bool {
			return item.Status != enums.PoAttachmentStatusApproved
		})

		var items models.PoAttachments = lo.Map(filterItems, func(item *models.PoAttachment, index int) *models.PoAttachment {
			item.Status = enums.PoAttachmentStatusApproved
			return item
		})

		updates.SellerPoAttachments = &items
	}

	err = r.db.Transaction(func(tx *gorm.DB) error {
		var sqlResult = tx.Model(&models.PurchaseOrder{}).Where("id = ? AND sample_maker_id = ?", params.PurchaseOrderID, params.GetUserID()).
			Updates(&updates)
		if sqlResult.RowsAffected == 0 {
			return errs.ErrInquirySellerInvalidToApprove
		}
		var tracking = models.PurchaseOrderTracking{
			PurchaseOrderID: params.PurchaseOrderID,
			ActionType:      enums.PoTrackingActionSellerApprovedPO,
			UserID:          params.GetUserID(),
			CreatedByUserID: params.JwtClaimsInfo.GetUserID(),
			UserGroup:       enums.PoTrackingUserGroupSeller,
			Metadata: &models.PoTrackingMetadata{
				Before: map[string]interface{}{
					"seller_tracking_status": purchaseOrder.SellerTrackingStatus,
				},
				After: map[string]interface{}{
					"seller_tracking_status": updates.SellerTrackingStatus,
					"seller_po_attachments":  purchaseOrder.SellerPoAttachments,
				},
			},
		}

		err = tx.Create(&tracking).Error
		return err
	})
	if err != nil {
		return err
	}

	return err

}

type SellerRejectPoParams struct {
	models.JwtClaimsInfo

	PurchaseOrderID string `json:"purchase_order_id" param:"purchase_order_id"`
	Reason          string `json:"reason"`
}

func (r *SellerPurchaseOrderRepo) SellerRejectPo(params SellerRejectPoParams) error {
	var purchaseOrder models.PurchaseOrder
	var err = r.db.Select("ID", "SellerPoAttachments", "TrackingStatus", "SellerTrackingStatus").
		First(&purchaseOrder, "id = ? AND sample_maker_id = ?", params.PurchaseOrderID, params.GetUserID()).Error
	if err != nil {
		return err
	}

	var updates = models.PurchaseOrder{
		SellerTrackingStatus: enums.SellerPoTrackingStatusRejectPO,
		SellerPORejectReason: params.Reason,
	}

	if purchaseOrder.SellerPoAttachments != nil {
		var filterItems = lo.Filter(*purchaseOrder.SellerPoAttachments, func(item *models.PoAttachment, index int) bool {
			return item.Status != enums.PoAttachmentStatusApproved
		})

		var items models.PoAttachments = lo.Map(filterItems, func(item *models.PoAttachment, index int) *models.PoAttachment {
			item.Status = enums.PoAttachmentStatusRejected
			return item
		})

		updates.SellerPoAttachments = &items
	}

	err = r.db.Transaction(func(tx *gorm.DB) error {
		var sqlResult = tx.Model(&models.PurchaseOrder{}).Where("id = ? AND sample_maker_id = ?", params.PurchaseOrderID, params.GetUserID()).
			Updates(&updates)
		if sqlResult.RowsAffected == 0 {
			return errs.ErrInquirySellerInvalidToApprove
		}
		var tracking = models.PurchaseOrderTracking{
			PurchaseOrderID: params.PurchaseOrderID,
			ActionType:      enums.PoTrackingActionSellerRejectedPO,
			UserID:          params.GetUserID(),
			CreatedByUserID: params.JwtClaimsInfo.GetUserID(),
			UserGroup:       enums.PoTrackingUserGroupSeller,
			Metadata: &models.PoTrackingMetadata{
				Before: map[string]interface{}{
					"seller_tracking_status": purchaseOrder.SellerTrackingStatus,
				},
				After: map[string]interface{}{
					"seller_tracking_status": updates.SellerTrackingStatus,
					"seller_po_attachment":   purchaseOrder.SellerPoAttachments,
				},
			},
		}

		err = tx.Create(&tracking).Error
		return err
	})
	if err != nil {
		return err
	}

	return err
}

func (r *SellerPurchaseOrderRepo) AdminSellerPurchaseOrderUploadCommentStatusCount(params PoCommentStatusCountParams) *models.CommentStatusCountItem {
	var resp models.CommentStatusCountItem

	r.db.Model(&models.Comment{}).Where("user_id != ? AND target_type = ? AND target_id = ? AND seen_at IS NULL", params.GetUserID(), enums.CommentTargetTypeSellerPoUpload, params.PurchaseOrderID).Count(&resp.UnseenCount)

	return &resp
}

type AdminSellerPurchaseOrderMarkDesignApprovalParams struct {
	models.JwtClaimsInfo

	PurchaseOrderID string `json:"purchase_order_id" param:"purchase_order_id" query:"purchase_order_id" validate:"required"`
}

func (r *SellerPurchaseOrderRepo) AdminSellerPurchaseOrderMarkDesignApproval(params AdminSellerPurchaseOrderMarkDesignApprovalParams) (*models.PurchaseOrder, error) {
	order, err := NewPurchaseOrderRepo(r.db).GetPurchaseOrder(GetPurchaseOrderParams{
		PurchaseOrderID: params.PurchaseOrderID,
		JwtClaimsInfo:   params.JwtClaimsInfo,
	})

	if err != nil {
		return nil, err
	}

	validStatus := []enums.SellerPoTrackingStatus{enums.SellerPoTrackingStatusPaymentConfirmed}
	if ok := slices.Contains(validStatus, order.SellerTrackingStatus); !ok {
		return nil, errs.ErrPoInvalidToMarkMaking
	}

	var updates models.PurchaseOrder
	updates.SellerTrackingStatus = enums.SellerPoTrackingStatusDesignApproval

	if order.SellerTrackingStatus != updates.SellerTrackingStatus {
		err = r.db.Transaction(func(tx *gorm.DB) error {
			err = NewPurchaseOrderTrackingRepo(r.db).CreatePurchaseOrderTrackingTx(tx, models.PurchaseOrderTrackingCreateForm{
				PurchaseOrderID: params.PurchaseOrderID,
				ActionType:      enums.PoTrackingActionAdminApprovedDesign,
				CreatedByUserID: params.JwtClaimsInfo.GetUserID(),
				UserGroup:       enums.PoTrackingUserGroupSeller,
				UserID:          order.SampleMakerID,
				Metadata: &models.PoTrackingMetadata{
					Before: map[string]interface{}{
						"seller_tracking_status": order.TrackingStatus,
					},
					After: map[string]interface{}{
						"seller_tracking_status": updates.TrackingStatus,
						"seller_design":          order.SellerDesign,
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

	order.SellerTrackingStatus = updates.SellerTrackingStatus

	return order, err
}

type AdminSellerPurchaseOrderUpdateDesignParams struct {
	models.JwtClaimsInfo

	PurchaseOrderID string                     `json:"purchase_order_id" param:"purchase_order_id" query:"purchase_order_id" validate:"required"`
	SellerDesign    *models.SellerPoDesignMeta `json:"seller_design" param:"seller_design"`
}

func (r *SellerPurchaseOrderRepo) AdminSellerPurchaseOrderUpdateDesign(params AdminSellerPurchaseOrderUpdateDesignParams) (*models.PurchaseOrder, error) {
	order, err := NewPurchaseOrderRepo(r.db).GetPurchaseOrder(GetPurchaseOrderParams{
		PurchaseOrderID: params.PurchaseOrderID,
		JwtClaimsInfo:   params.JwtClaimsInfo,
	})
	if err != nil {
		return nil, err
	}

	var updates = models.PurchaseOrder{
		SellerDesign: params.SellerDesign,
	}

	err = r.db.Transaction(func(tx *gorm.DB) error {
		err = NewPurchaseOrderTrackingRepo(r.db).CreatePurchaseOrderTrackingTx(tx, models.PurchaseOrderTrackingCreateForm{
			PurchaseOrderID: params.PurchaseOrderID,
			ActionType:      enums.PoTrackingActionAdminUpdatedDesign,
			CreatedByUserID: params.JwtClaimsInfo.GetUserID(),
			UserGroup:       enums.PoTrackingUserGroupSeller,
			UserID:          order.SampleMakerID,
			Metadata: &models.PoTrackingMetadata{
				Before: map[string]interface{}{
					"seller_tracking_status": order.TrackingStatus,
				},
				After: map[string]interface{}{
					"seller_tracking_status": updates.TrackingStatus,
					"seller_design":          updates.SellerDesign,
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

	order.SellerDesign = updates.SellerDesign

	return order, err
}

func (r *SellerPurchaseOrderRepo) AdminSellerPoDesignCommentMarkSeen(params PoCommentMarkSeenParams) error {
	var err = r.db.Model(&models.Comment{}).
		Where("seen_at IS NULL").
		Where("user_id != ?", params.GetUserID()).
		Where("target_type = ? AND target_id = ?", enums.CommentTargetTypeSellerPoDesign, params.PurchaseOrderID).
		Update("seen_at", time.Now().Unix()).Error

	return err
}

func (r *SellerPurchaseOrderRepo) SellerPoDesignCommentMarkSeen(params PoCommentMarkSeenParams) error {
	var err = r.db.Model(&models.Comment{}).
		Where("seen_at IS NULL").
		Where("user_id != ?", params.GetUserID()).
		Where("target_type = ? AND target_id = ?", enums.CommentTargetTypeSellerPoDesign, params.PurchaseOrderID).
		Update("seen_at", time.Now().Unix()).Error

	return err
}

func (r *SellerPurchaseOrderRepo) AdminSellerPoDesignCommentStatusCount(params PoCommentStatusCountParams) (*models.CommentStatusCountItem, error) {
	order, err := NewPurchaseOrderRepo(r.db).GetPurchaseOrder(GetPurchaseOrderParams{
		PurchaseOrderID: params.PurchaseOrderID,
		JwtClaimsInfo:   params.JwtClaimsInfo,
	})
	if err != nil {
		return nil, eris.Wrap(err, err.Error())
	}

	var resp models.CommentStatusCountItem

	r.db.Model(&models.Comment{}).Where("user_id != ? AND target_type = ? AND target_id = ? AND seen_at IS NULL", params.GetUserID(), enums.CommentTargetTypeSellerPoDesign, order.ID).Count(&resp.UnseenCount)

	return &resp, nil
}

func (r *SellerPurchaseOrderRepo) SellerPoDesignCommentStatusCount(params PoCommentStatusCountParams) (*models.CommentStatusCountItem, error) {
	order, err := NewPurchaseOrderRepo(r.db).GetPurchaseOrder(GetPurchaseOrderParams{
		PurchaseOrderID: params.PurchaseOrderID,
		JwtClaimsInfo:   params.JwtClaimsInfo,
	})
	if err != nil {
		return nil, eris.Wrap(err, err.Error())
	}

	var resp models.CommentStatusCountItem

	r.db.Model(&models.Comment{}).Where("user_id != ? AND target_type = ? AND target_id = ? AND seen_at IS NULL", params.GetUserID(), enums.CommentTargetTypeSellerPoDesign, order.ID).Count(&resp.UnseenCount)

	return &resp, nil
}

type AdminSellerPurchaseOrderUpdateFinalDesignParams struct {
	models.JwtClaimsInfo

	PurchaseOrderID           string              `json:"purchase_order_id" param:"purchase_order_id" query:"purchase_order_id" validate:"required"`
	SellerTechpackAttachments *models.Attachments `json:"seller_techpack_attachments" param:"seller_techpack_attachments"`
}

func (r *SellerPurchaseOrderRepo) AdminSellerPurchaseOrderUpdateFinalDesign(params AdminSellerPurchaseOrderUpdateFinalDesignParams) (*models.PurchaseOrder, error) {
	order, err := NewPurchaseOrderRepo(r.db).GetPurchaseOrder(GetPurchaseOrderParams{
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

	err = r.db.Model(&models.PurchaseOrder{}).
		Where("id = ?", params.PurchaseOrderID).
		Updates(&updates).Error

	if err != nil {
		return nil, eris.Wrap(err, err.Error())
	}

	return order, err
}

type AdminSellerPurchaseOrderApproveFinalDesignParams struct {
	models.JwtClaimsInfo

	PurchaseOrderID string `json:"purchase_order_id" param:"purchase_order_id" query:"purchase_order_id" validate:"required"`
}

func (r *SellerPurchaseOrderRepo) AdminSellerPurchaseOrderApproveFinalDesign(params AdminSellerPurchaseOrderApproveFinalDesignParams) (*models.PurchaseOrder, error) {
	order, err := NewPurchaseOrderRepo(r.db).GetPurchaseOrder(GetPurchaseOrderParams{
		PurchaseOrderID: params.PurchaseOrderID,
		JwtClaimsInfo:   params.JwtClaimsInfo,
	})
	if err != nil {
		return nil, err
	}

	var updates models.PurchaseOrder
	updates.SellerTrackingStatus = enums.SellerPoTrackingStatusDesignApprovedByAdmin

	err = r.db.Model(&models.PurchaseOrder{}).
		Where("id = ?", params.PurchaseOrderID).
		Updates(&updates).Error

	if err != nil {
		return nil, eris.Wrap(err, err.Error())
	}

	return order, err
}

type SellerMarkRawMaterialParams struct {
	models.JwtClaimsInfo

	PurchaseOrderID string `json:"purchase_order_id" param:"purchase_order_id"`
}

func (r *SellerPurchaseOrderRepo) SellerMarkRawMaterial(params SellerMarkRawMaterialParams) error {
	var updates models.PurchaseOrder
	updates.SellerTrackingStatus = enums.SellerPoTrackingStatusRawMaterial

	var err = r.db.Model(&models.PurchaseOrder{}).Where("id = ?", params.PurchaseOrderID).Updates(&updates).Error
	if err != nil {
		return err
	}

	return err
}

func (r *SellerPurchaseOrderRepo) AdminSellerFinalDesignCommentMarkSeen(params PoCommentMarkSeenParams) error {
	var err = r.db.Model(&models.Comment{}).
		Where("seen_at IS NULL").
		Where("user_id != ?", params.GetUserID()).
		Where("target_type = ? AND target_id = ? AND file_key = ?", enums.CommentTargetTypeSellerPoFinalDesign, params.PurchaseOrderID, params.FileKey).
		Update("seen_at", time.Now().Unix()).Error

	return err
}

func (r *SellerPurchaseOrderRepo) AdminSellerFinalDesignCommentStatusCount(params PoCommentStatusCountParams) ([]*models.CommentStatusCountItem, error) {
	order, err := NewPurchaseOrderRepo(r.db).GetPurchaseOrder(GetPurchaseOrderParams{
		PurchaseOrderID: params.PurchaseOrderID,
		JwtClaimsInfo:   params.JwtClaimsInfo,
	})
	if err != nil {
		return nil, eris.Wrap(err, err.Error())
	}

	var items []*models.CommentStatusCountItem

	if order.SellerTechpackAttachments != nil && len(*order.SellerTechpackAttachments) > 0 {
		for _, fileItem := range *order.SellerTechpackAttachments {
			count := int64(0)
			if fileItem.FileKey != "" {
				r.db.Model(&models.Comment{}).Where("user_id != ? AND target_type = ? AND target_id = ? AND file_key = ? AND seen_at IS NULL", params.GetUserID(), enums.CommentTargetTypeSellerPoFinalDesign, params.PurchaseOrderID, fileItem.FileKey).Count(&count)
				items = append(items, &models.CommentStatusCountItem{
					FileKey:     fileItem.FileKey,
					UnseenCount: count,
				})
			}

		}
	}

	return items, nil
}

func (r *SellerPurchaseOrderRepo) SellerFinalDesignCommentStatusCount(params PoCommentStatusCountParams) ([]*models.CommentStatusCountItem, error) {
	order, err := NewPurchaseOrderRepo(r.db).GetPurchaseOrder(GetPurchaseOrderParams{
		PurchaseOrderID: params.PurchaseOrderID,
		JwtClaimsInfo:   params.JwtClaimsInfo,
	})
	if err != nil {
		return nil, eris.Wrap(err, err.Error())
	}

	var items []*models.CommentStatusCountItem

	if order.SellerTechpackAttachments != nil && len(*order.SellerTechpackAttachments) > 0 {
		for _, fileItem := range *order.SellerTechpackAttachments {
			count := int64(0)
			if fileItem.FileKey != "" {
				r.db.Model(&models.Comment{}).Where("user_id != ? AND target_type = ? AND target_id = ? AND file_key = ? AND seen_at IS NULL", params.GetUserID(), enums.CommentTargetTypeSellerPoFinalDesign, params.PurchaseOrderID, fileItem.FileKey).Count(&count)
				items = append(items, &models.CommentStatusCountItem{
					FileKey:     fileItem.FileKey,
					UnseenCount: count,
				})
			}

		}
	}

	return items, nil
}

type SellerUpdatePurchaseOrderRawMaterialParams struct {
	models.JwtClaimsInfo

	PurchaseOrderID string                    `json:"purchase_order_id" param:"purchase_order_id" query:"purchase_order_id" validate:"required"`
	PoRawMaterials  models.PoRawMaterialMetas `json:"po_raw_materials" param:"po_raw_materials" query:"po_raw_materials"`
}

func (r *SellerPurchaseOrderRepo) SellerUpdatePurchaseOrderRawMaterial(params SellerUpdatePurchaseOrderRawMaterialParams) (*models.PurchaseOrder, error) {
	order, err := NewPurchaseOrderRepo(r.db).GetPurchaseOrder(GetPurchaseOrderParams{
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

	updates.GenerateRawMaterialRefID(updates.PoRawMaterials)

	err = r.db.Model(&models.PurchaseOrder{}).
		Where("id = ?", params.PurchaseOrderID).
		Updates(&updates).Error

	if err != nil {
		return nil, eris.Wrap(err, err.Error())
	}

	return order, err
}

type SellerPurchaseOrderUpdateRawMaterialParams struct {
	models.JwtClaimsInfo

	PurchaseOrderID      string                     `json:"purchase_order_id" param:"purchase_order_id" query:"purchase_order_id" validate:"required"`
	PoRawMaterials       *models.PoRawMaterialMetas `json:"po_raw_materials" param:"po_raw_materials" query:"po_raw_materials"`
	ApproveRawMaterialAt *int64                     `json:"approve_raw_material_at" param:"approve_raw_material_at" query:"approve_raw_material_at"`
}

func (r *PurchaseOrderRepo) SellerPurchaseOrderUpdateRawMaterial(params SellerPurchaseOrderUpdateRawMaterialParams) (*models.PurchaseOrder, error) {
	order, err := r.GetPurchaseOrder(GetPurchaseOrderParams{
		PurchaseOrderID: params.PurchaseOrderID,
		JwtClaimsInfo:   params.JwtClaimsInfo,
	})
	if err != nil {
		return nil, err
	}

	var updates = models.PurchaseOrder{
		SellerPoRawMaterials: params.PoRawMaterials,
		SellerTrackingStatus: enums.SellerPoTrackingStatusRawMaterial,
	}
	updates.GenerateRawMaterialRefID(updates.SellerPoRawMaterials)

	err = r.db.Transaction(func(tx *gorm.DB) error {
		NewPurchaseOrderTrackingRepo(r.db).CreatePurchaseOrderTrackingTx(tx, models.PurchaseOrderTrackingCreateForm{
			PurchaseOrderID: order.ID,
			ActionType:      enums.PoTrackingActionUpdateMaterial,
			UserID:          order.UserID,
			CreatedByUserID: params.JwtClaimsInfo.GetUserID(),
			Metadata: &models.PoTrackingMetadata{
				Before: map[string]interface{}{
					"seller_po_raw_materials": order.SellerPoRawMaterials,
				},
				After: map[string]interface{}{
					"seller_po_raw_materials": updates.SellerPoRawMaterials,
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

	order.SellerTrackingStatus = updates.SellerTrackingStatus

	return order, err
}

type SellerPurchaseOrderReceivePaymentParams struct {
	models.JwtClaimsInfo

	PurchaseOrderID string `json:"purchase_order_id" param:"purchase_order_id" query:"purchase_order_id" validate:"required"`
}

func (r *SellerPurchaseOrderRepo) SellerPurchaseOrderReceivePayment(params SellerPurchaseOrderReceivePaymentParams) (*models.PurchaseOrder, error) {
	order, err := NewPurchaseOrderRepo(r.db).GetPurchaseOrder(GetPurchaseOrderParams{
		PurchaseOrderID: params.PurchaseOrderID,
		JwtClaimsInfo:   params.JwtClaimsInfo,
	})
	if err != nil {
		return nil, err
	}

	var updates = models.PurchaseOrder{
		PayoutMarkAsReceivedAt: values.Int64(r.db.NowFunc().Unix()),
		SellerTrackingStatus:   enums.SellerPoTrackingStatusDesignApproval,
	}

	err = r.db.Transaction(func(tx *gorm.DB) error {
		if order.SellerTrackingStatus != updates.SellerTrackingStatus {
			err = NewPurchaseOrderTrackingRepo(r.db).CreatePurchaseOrderTrackingTx(tx, models.PurchaseOrderTrackingCreateForm{
				PurchaseOrderID: params.PurchaseOrderID,
				ActionType:      enums.PoTrackingActionSellerPaymentReceived,
				UserID:          params.GetUserID(),
				CreatedByUserID: params.JwtClaimsInfo.GetUserID(),
				UserGroup:       enums.PoTrackingUserGroupSeller,
				Metadata: &models.PoTrackingMetadata{
					Before: map[string]interface{}{
						"seller_tracking_status": order.TrackingStatus,
					},
					After: map[string]interface{}{
						"seller_tracking_status": updates.TrackingStatus,
					},
				},
			})
			if err != nil {
				return eris.Wrap(err, err.Error())
			}

		}
		return tx.Model(&models.PurchaseOrder{}).
			Where("id = ?", params.PurchaseOrderID).
			Updates(&updates).Error

	})

	return order, err

}

type AdminSellerPurchaseOrderUploadPoParams struct {
	models.JwtClaimsInfo

	PurchaseOrderID     string                `json:"purchase_order_id" param:"purchase_order_id" query:"purchase_order_id" validate:"required"`
	SellerPoAttachments *models.PoAttachments `json:"seller_po_attachments,omitempty" validate:"required"`
}

func (r *SellerPurchaseOrderRepo) AdminSellerPurchaseOrderUploadPo(params AdminSellerPurchaseOrderUploadPoParams) (*models.PurchaseOrder, error) {
	order, err := NewPurchaseOrderRepo(r.db).GetPurchaseOrder(GetPurchaseOrderParams{
		PurchaseOrderID: params.PurchaseOrderID,
		JwtClaimsInfo:   params.JwtClaimsInfo,
	})
	if err != nil {
		return nil, err
	}

	var checkStatuses = []enums.SellerPoTrackingStatus{
		enums.SellerPoTrackingStatusRejectPO,
		enums.SellerPoTrackingStatusNew,
	}
	if !lo.Contains(checkStatuses, order.SellerTrackingStatus) {
		return nil, errs.ErrPOIsNotAbleToUploadPO.WithDetailMessagef("SellerTrackingStatus=%s", order.SellerTrackingStatus)
	}
	var updates = models.PurchaseOrder{
		SellerPoAttachments:  params.SellerPoAttachments,
		SellerTrackingStatus: enums.SellerPoTrackingStatusNew,
		SellerPORejectReason: "",
	}

	err = r.db.Transaction(func(tx *gorm.DB) error {
		err = NewPurchaseOrderTrackingRepo(r.db).CreatePurchaseOrderTrackingTx(tx, models.PurchaseOrderTrackingCreateForm{
			PurchaseOrderID: params.PurchaseOrderID,
			ActionType:      enums.PoTrackingActionAdminUploadedPO,
			UserID:          order.SampleMakerID,
			UserGroup:       enums.PoTrackingUserGroupSeller,
			CreatedByUserID: params.GetUserID(),
			Metadata: &models.PoTrackingMetadata{
				Before: map[string]interface{}{
					"seller_tracking_status": order.SellerTrackingStatus,
					"seller_po_attachments":  order.SellerPoAttachments,
				},
				After: map[string]interface{}{
					"seller_tracking_status": updates.SellerTrackingStatus,
					"seller_po_attachments":  params.SellerPoAttachments,
				},
			},
		})
		if err != nil {
			return eris.Wrap(err, err.Error())
		}

		return tx.Model(&models.PurchaseOrder{}).
			Select("SellerPoAttachments", "SellerTrackingStatus", "SellerPORejectReason").
			Where("id = ?", params.PurchaseOrderID).
			Updates(&updates).Error
	})
	if err != nil {
		return nil, eris.Wrap(err, err.Error())
	}

	order.SellerPoAttachments = updates.SellerPoAttachments
	order.SellerTrackingStatus = updates.SellerTrackingStatus

	return order, err
}

type AdminSellerPoConfirmDeliveredParams struct {
	models.JwtClaimsInfo

	PurchaseOrderID string `json:"purchase_order_id" param:"purchase_order_id" query:"purchase_order_id" validate:"required"`
}

func (r *SellerPurchaseOrderRepo) AdminSellerPoConfirmDelivered(params AdminSellerPoConfirmDeliveredParams) (*models.PurchaseOrder, error) {
	order, err := NewPurchaseOrderRepo(r.db).GetPurchaseOrder(GetPurchaseOrderParams{
		PurchaseOrderID: params.PurchaseOrderID,
		JwtClaimsInfo:   params.JwtClaimsInfo,
	})
	if err != nil {
		return nil, err
	}
	validStatus := []enums.SellerPoTrackingStatus{enums.SellerPoTrackingStatusDelivering}
	if ok := slices.Contains(validStatus, order.SellerTrackingStatus); !ok {
		return nil, errs.ErrPoInvalidToConfirmDelivered
	}

	var updates models.PurchaseOrder
	err = copier.Copy(&updates, &params)
	if err != nil {
		return nil, err
	}
	updates.SellerDeliveryConfirmedAt = values.Int64(time.Now().Unix())
	updates.SellerTrackingStatus = enums.SellerPoTrackingStatusDeliveryConfirmed

	if order.SellerTrackingStatus != updates.SellerTrackingStatus {
		err = r.db.Transaction(func(tx *gorm.DB) error {
			err = NewPurchaseOrderTrackingRepo(r.db).CreatePurchaseOrderTrackingTx(tx, models.PurchaseOrderTrackingCreateForm{
				PurchaseOrderID: params.PurchaseOrderID,
				ActionType:      enums.PoTrackingActionSellerConfirmDelivered,
				UserID:          order.SampleMakerID,
				UserGroup:       enums.PoTrackingUserGroupSeller,
				CreatedByUserID: params.JwtClaimsInfo.GetUserID(),
				Metadata: &models.PoTrackingMetadata{
					Before: map[string]interface{}{
						"tracking_status": order.SellerTrackingStatus,
					},
					After: map[string]interface{}{
						"tracking_status": updates.SellerTrackingStatus,
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

	order.SellerDeliveryConfirmedAt = updates.SellerDeliveryConfirmedAt
	order.SellerTrackingStatus = updates.SellerTrackingStatus

	return order, err
}

type AdminSellerPurchaseOrderDeliveryFeedbackParams struct {
	models.JwtClaimsInfo
	PurchaseOrderID        string `json:"purchase_order_id" param:"purchase_order_id"`
	SellerDeliveryFeedback string `json:"seller_delivery_feedback" param:"seller_delivery_feedback"`
}

func (r *SellerPurchaseOrderRepo) AdminSellerPurchaseOrderDeliveryFeedback(params AdminSellerPurchaseOrderDeliveryFeedbackParams) (err error) {
	var updates = &models.PurchaseOrder{
		SellerDeliveryFeedback: params.SellerDeliveryFeedback,
	}
	err = r.db.Model(&models.PurchaseOrder{}).Where("id = ?", params.PurchaseOrderID).Updates(updates).Error
	return
}

func (r *SellerPurchaseOrderRepo) SellerFinalDesignCommentMarkSeen(params PoCommentMarkSeenParams) error {
	var err = r.db.Model(&models.Comment{}).
		Where("seen_at IS NULL").
		Where("user_id != ?", params.GetUserID()).
		Where("target_type = ? AND target_id = ? AND file_key = ?", enums.CommentTargetTypeSellerPoFinalDesign, params.PurchaseOrderID, params.FileKey).
		Update("seen_at", time.Now().Unix()).Error

	return err
}

func (r *SellerPurchaseOrderRepo) SellerPurchaseOrderUploadCommentStatusCount(params PoCommentStatusCountParams) *models.CommentStatusCountItem {
	var resp models.CommentStatusCountItem

	r.db.Model(&models.Comment{}).Where("user_id != ? AND target_type = ? AND target_id = ? AND seen_at IS NULL", params.GetUserID(), enums.CommentTargetTypeSellerPoUpload, params.PurchaseOrderID).Count(&resp.UnseenCount)

	return &resp
}

func (r *SellerPurchaseOrderRepo) SellerPoRawMaterialCommentMarkSeen(params PoCommentMarkSeenParams) error {
	var err = r.db.Model(&models.Comment{}).
		Where("seen_at IS NULL").
		Where("user_id != ?", params.GetUserID()).
		Where("target_type = ? AND target_id = ? AND file_key = ?", enums.CommentTargetTypeSellerPoRawMaterial, params.PurchaseOrderID, params.FileKey).
		Update("seen_at", time.Now().Unix()).Error

	return err
}

type SellerPurchaseOrderSkipRawMaterialParams struct {
	models.JwtClaimsInfo
	PurchaseOrderID string `json:"purchase_order_id" param:"purchase_order_id" validate:"required"`
}

func (r *SellerPurchaseOrderRepo) SellerPurchaseOrderSkipRawMaterial(params SellerPurchaseOrderSkipRawMaterialParams) (*models.PurchaseOrder, error) {
	order, err := NewPurchaseOrderRepo(r.db).GetPurchaseOrder(GetPurchaseOrderParams{
		JwtClaimsInfo:   params.JwtClaimsInfo,
		PurchaseOrderID: params.PurchaseOrderID,
	})
	if err != nil {
		return nil, err
	}

	var updates = models.PurchaseOrder{
		SellerTrackingStatus: enums.SellerPoTrackingStatusRawMaterialSkipped,
	}

	r.db.Transaction(func(tx *gorm.DB) error {
		if order.SellerTrackingStatus != updates.SellerTrackingStatus {
			err = NewPurchaseOrderTrackingRepo(r.db).CreatePurchaseOrderTrackingTx(tx, models.PurchaseOrderTrackingCreateForm{
				PurchaseOrderID: params.PurchaseOrderID,
				ActionType:      enums.PoTrackingActionSellerSkipRawMaterial,
				UserID:          params.GetUserID(),
				CreatedByUserID: params.JwtClaimsInfo.GetUserID(),
				UserGroup:       enums.PoTrackingUserGroupSeller,
				Metadata: &models.PoTrackingMetadata{
					Before: map[string]interface{}{
						"seller_tracking_status": order.SellerTrackingStatus,
					},
					After: map[string]interface{}{
						"seller_tracking_status": updates.SellerTrackingStatus,
					},
				},
			})
			if err != nil {
				return eris.Wrap(err, err.Error())
			}

		}

		var sqlResult = tx.Model(&models.PurchaseOrder{}).
			Where("id = ? AND sample_maker_id = ?", params.PurchaseOrderID, params.GetUserID()).
			Updates(&updates)
		if sqlResult.Error != nil {
			return sqlResult.Error
		}

		if sqlResult.RowsAffected == 0 {
			return errs.ErrBulkPoNotFound
		}

		return nil
	})

	order.SellerTrackingStatus = updates.SellerTrackingStatus
	return order, err
}

type AdminSellerPurchaseOrderPayoutParams struct {
	models.JwtClaimsInfo
	PurchaseOrderID string `json:"purchase_order_id" param:"purchase_order_id" validate:"required"`

	TransactionRefID      string             `json:"transaction_ref_id" validate:"required"`
	TransactionAttachment *models.Attachment `json:"transaction_attachment" validate:"required"`
}

func (r *SellerPurchaseOrderRepo) AdminSellerPurchaseOrderPayout(params AdminSellerPurchaseOrderPayoutParams) (*models.PurchaseOrder, error) {
	purchaseOrder, err := NewPurchaseOrderRepo(r.db).GetPurchaseOrder(GetPurchaseOrderParams{
		PurchaseOrderID:      params.PurchaseOrderID,
		IncludeInquirySeller: true,
	})
	if err != nil {
		return nil, err
	}

	if purchaseOrder.SampleMakerID == "" || purchaseOrder.SampleMakerID == "inflow" {
		return nil, errs.ErrSellerInvalidToPayout
	}

	if purchaseOrder.InquirySeller == nil {
		return nil, errs.ErrSellerInvalidToPayout
	}

	var totalPrice price.Price
	if purchaseOrder.InquirySeller.SampleUnitPrice == nil {
		purchaseOrder.InquirySeller.SampleUnitPrice = purchaseOrder.InquirySeller.GetSampleUnitPrice().ToPtr()
		r.db.Model(&models.InquirySeller{}).Where("id = ?", purchaseOrder.InquirySeller.ID).UpdateColumn("SampleUnitPrice", purchaseOrder.InquirySeller.SampleUnitPrice)
	}

	if len(purchaseOrder.Items) > 0 {
		for _, cartItem := range purchaseOrder.Items {
			totalPrice = totalPrice.Add(purchaseOrder.InquirySeller.SampleUnitPrice.MultipleInt(cartItem.Quantity))
		}
	} else if len(purchaseOrder.OrderCartItems) > 0 {
		for _, cartItem := range purchaseOrder.OrderCartItems {
			totalPrice = totalPrice.Add(purchaseOrder.InquirySeller.SampleUnitPrice.MultipleInt(cartItem.Qty))
		}
	} else if len(purchaseOrder.CartItems) > 0 {
		for _, cartItem := range purchaseOrder.CartItems {
			totalPrice = totalPrice.Add(purchaseOrder.InquirySeller.SampleUnitPrice.MultipleInt(cartItem.Qty))
		}
	}

	var updates = models.PurchaseOrder{
		PayoutTransactionRefID:      params.TransactionRefID,
		PayoutTransactionAttachment: params.TransactionAttachment,
		SellerTrackingStatus:        enums.SellerPoTrackingStatusPaymentConfirmed,
		PayoutTransferedAt:          values.Int64(time.Now().Unix()),
		PayoutMarkAsPaidAt:          values.Int64(time.Now().Unix()),
		SellerPricing: models.SellerPricing{
			SellerTotalPrice: totalPrice.ToPtr(),
		},
	}

	var transaction = models.PaymentTransaction{
		PurchaseOrderID:   purchaseOrder.ID,
		PaidAmount:        totalPrice.ToPtr(),
		PaymentType:       enums.PaymentTypeBankTransfer,
		Milestone:         enums.PaymentMilestoneFinalPayment,
		UserID:            purchaseOrder.UserID,
		TransactionRefID:  params.TransactionRefID,
		Status:            enums.PaymentStatusPaid,
		PaymentPercentage: values.Float64(100),
		TotalAmount:       totalPrice.ToPtr(),
		Currency:          purchaseOrder.Inquiry.Currency,
		TransactionType:   enums.TransactionTypeDebit,
		Attachments: &models.Attachments{
			params.TransactionAttachment,
		},
		Metadata: &models.PaymentTransactionMetadata{
			InquiryID:                purchaseOrder.Inquiry.ID,
			InquiryReferenceID:       purchaseOrder.Inquiry.ReferenceID,
			PurchaseOrderReferenceID: purchaseOrder.ReferenceID,
			PurchaseOrderID:          purchaseOrder.ID,
		},
	}

	err = r.db.Transaction(func(tx *gorm.DB) error {
		err = tx.Create(&transaction).Error
		if err != nil {
			return err
		}

		updates.PaymentTransactionReferenceID = transaction.ReferenceID
		err = tx.Model(&models.PurchaseOrder{}).Where("id = ?", params.PurchaseOrderID).Updates(&updates).Error

		return err

	})
	if err != nil {
		return nil, err
	}

	purchaseOrder.PayoutTransactionRefID = updates.PayoutTransactionRefID
	purchaseOrder.PayoutTransactionAttachment = updates.PayoutTransactionAttachment
	purchaseOrder.SellerTrackingStatus = updates.SellerTrackingStatus
	purchaseOrder.SellerTotalPrice = totalPrice.ToPtr()
	return purchaseOrder, err
}

type AdminSellerPurchaseOrderPreviewPayoutParams struct {
	models.JwtClaimsInfo
	PurchaseOrderID string `json:"purchase_order_id" param:"purchase_order_id" validate:"required"`
}

func (r *SellerPurchaseOrderRepo) AdminSellerPurchaseOrderPreviewPayout(params AdminSellerPurchaseOrderPreviewPayoutParams) (*models.PurchaseOrder, error) {
	purchaseOrder, err := NewPurchaseOrderRepo(r.db).GetPurchaseOrder(GetPurchaseOrderParams{
		PurchaseOrderID:      params.PurchaseOrderID,
		IncludeInquirySeller: true,
	})
	if err != nil {
		return nil, err
	}

	if purchaseOrder.SampleMakerID == "" || purchaseOrder.SampleMakerID == "inflow" {
		return nil, errs.ErrSellerInvalidToPayout
	}

	if purchaseOrder.InquirySeller == nil {
		return nil, errs.ErrSellerInvalidToPayout
	}

	var totalPrice price.Price
	if purchaseOrder.InquirySeller.SampleUnitPrice == nil {
		purchaseOrder.InquirySeller.SampleUnitPrice = purchaseOrder.InquirySeller.GetSampleUnitPrice().ToPtr()
		r.db.Model(&models.InquirySeller{}).Where("id = ?", purchaseOrder.InquirySeller.ID).UpdateColumn("SampleUnitPrice", purchaseOrder.InquirySeller.SampleUnitPrice)
	}

	if len(purchaseOrder.Items) > 0 {
		for _, cartItem := range purchaseOrder.Items {
			totalPrice = totalPrice.Add(purchaseOrder.InquirySeller.SampleUnitPrice.MultipleInt(cartItem.Quantity))
		}
	} else if len(purchaseOrder.OrderCartItems) > 0 {
		for _, cartItem := range purchaseOrder.OrderCartItems {
			totalPrice = totalPrice.Add(purchaseOrder.InquirySeller.SampleUnitPrice.MultipleInt(cartItem.Qty))
		}
	} else if len(purchaseOrder.CartItems) > 0 {
		for _, cartItem := range purchaseOrder.CartItems {
			totalPrice = totalPrice.Add(purchaseOrder.InquirySeller.SampleUnitPrice.MultipleInt(cartItem.Qty))
		}
	}

	var updates models.PurchaseOrder
	updates.SellerTotalPrice = totalPrice.ToPtr()

	err = r.db.Model(&models.PurchaseOrder{}).Where("id = ?", params.PurchaseOrderID).Updates(&updates).Error
	if err != nil {
		return nil, err
	}

	return purchaseOrder, err
}
