package repo

import (
	"database/sql"
	"strings"
	"time"

	"github.com/engineeringinflow/inflow-backend/pkg/db"
	"github.com/engineeringinflow/inflow-backend/pkg/errs"
	"github.com/engineeringinflow/inflow-backend/pkg/helper"
	"github.com/engineeringinflow/inflow-backend/pkg/models"
	"github.com/engineeringinflow/inflow-backend/pkg/models/enums"
	"github.com/engineeringinflow/inflow-backend/pkg/models/price"
	"github.com/engineeringinflow/inflow-backend/pkg/repo/query"
	"github.com/engineeringinflow/inflow-backend/pkg/repo/query/queryfunc"
	"github.com/jinzhu/copier"
	"github.com/lib/pq"
	"github.com/rotisserie/eris"
	"github.com/samber/lo"
	"github.com/thaitanloi365/go-utils/values"
	"golang.org/x/exp/slices"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type SellerBulkPurchaseOrderRepo struct {
	db *db.DB
}

func NewSellerBulkPurchaseOrderRepo(db *db.DB) *SellerBulkPurchaseOrderRepo {
	return &SellerBulkPurchaseOrderRepo{
		db: db,
	}

}

type GetBulkPurchaseOrderSellerQuotationParams struct {
	models.JwtClaimsInfo

	SellerQuotationID string `json:"seller_quotation_id" param:"seller_quotation_id" validate:"required"`
}

func (r *SellerBulkPurchaseOrderRepo) GetBulkPurchaseOrderSellerQuotation(params GetBulkPurchaseOrderSellerQuotationParams) (*models.BulkPurchaseOrderSellerQuotation, error) {
	var builder = queryfunc.NewBulkPurchaseOrderSellerQuotationBuilder(queryfunc.BulkPurchaseOrderSellerQuotationBuilderOptions{
		IncludeBulk: true,
	})
	var sellerQuotation models.BulkPurchaseOrderSellerQuotation
	var err = query.New(r.db, builder).
		WhereFunc(func(builder *query.Builder) {
			builder.Where("rq.id = ?", params.SellerQuotationID)
		}).
		FirstFunc(&sellerQuotation)

	if err != nil {
		if r.db.IsRecordNotFoundError(err) {
			return nil, errs.ErrRecordNotFound
		}
		return nil, err
	}

	return &sellerQuotation, nil
}

type SellerSubmitBulkQuotationParams struct {
	models.JwtClaimsInfo

	SellerQuotationID string `json:"seller_quotation_id" param:"seller_quotation_id" validate:"required"`

	FabricCost     *price.Price                   `json:"fabric_cost,omitempty"`
	DecorationCost *price.Price                   `json:"decoration_cost,omitempty"`
	MakingCost     *price.Price                   `json:"making_cost,omitempty"` // sewing, cut, making, finishing
	OtherCost      *price.Price                   `json:"other_cost,omitempty"`
	SellerRemark   string                         `json:"seller_remark,omitempty"`
	BulkQuotations models.SellerBulkQuotationMOQs `json:"bulk_quotations,omitempty"`
}

func (r *SellerBulkPurchaseOrderRepo) SellerSubmitBulkQuotation(form SellerSubmitBulkQuotationParams) (*models.BulkPurchaseOrderSellerQuotation, error) {
	record, err := r.GetBulkPurchaseOrderSellerQuotation(GetBulkPurchaseOrderSellerQuotationParams{
		JwtClaimsInfo:     form.JwtClaimsInfo,
		SellerQuotationID: form.SellerQuotationID,
	})
	if err != nil {
		return nil, err
	}

	if record.Status == enums.BulkPurchaseOrderSellerStatusApproved {
		return nil, errs.ErrBulkPurchaseOrderSellerQuotationInvalidToSubmit
	}
	var updates models.BulkPurchaseOrderSellerQuotation
	err = copier.Copy(&updates, &form)
	if err != nil {
		return nil, err
	}

	updates.QuotationAt = values.Int64(time.Now().Unix())
	updates.Status = enums.BulkPurchaseOrderSellerStatusWaitingForApproval
	updates.QuotedPrice = updates.FabricCost.AddPtr(updates.MakingCost).AddPtr(updates.DecorationCost).AddPtr(updates.OtherCost).ToPtr()

	if updates.BulkQuotations != nil {
		updates.BulkQuotations = lo.Map(updates.BulkQuotations, func(item *models.SellerBulkQuotationMOQ, index int) *models.SellerBulkQuotationMOQ {
			item.UnitPrice = updates.QuotedPrice.Add(item.UpCharge)
			item.TotalPrice = item.UnitPrice.MultipleInt(values.Int64Value(item.Quantity))
			updates.QuotedPrice = &item.UnitPrice
			return item
		})

	}

	err = r.db.Transaction(func(tx *gorm.DB) error {
		var sqlResult = tx.Clauses(clause.OnConflict{DoNothing: true}).
			Where("id = ? AND user_id = ?", form.SellerQuotationID, form.GetUserID()).
			Updates(updates)
		if sqlResult.Error != nil {
			return sqlResult.Error
		}

		if sqlResult.RowsAffected == 0 {
			return errs.ErrBulkPurchaseOrderSellerQuotationNotFound
		}

		var log = models.BulkPurchaseOrderTracking{
			PurchaseOrderID: record.BulkPurchaseOrderID,
			ActionType:      enums.BulkPoTrackingActionSellerSubmitQuotation,
			UserGroup:       enums.PoTrackingUserGroupSeller,
			UserID:          record.UserID,
			CreatedByUserID: form.JwtClaimsInfo.GetUserID(),
			Metadata: &models.PoTrackingMetadata{
				Before: map[string]interface{}{
					"seller_quotation_status": record.Status,
				},
				After: map[string]interface{}{
					"seller_quotation_status": updates.Status,
				},
			},
		}
		return tx.Create(&log).Error

	})
	if err != nil {
		return nil, err
	}

	record.QuotationAt = updates.QuotationAt
	record.Status = updates.Status

	return record, err
}

func (r *SellerBulkPurchaseOrderRepo) SellerReSubmitBulkQuotation(form SellerSubmitBulkQuotationParams) (*models.BulkPurchaseOrderSellerQuotation, error) {
	record, err := r.GetBulkPurchaseOrderSellerQuotation(GetBulkPurchaseOrderSellerQuotationParams{
		JwtClaimsInfo:     form.JwtClaimsInfo,
		SellerQuotationID: form.SellerQuotationID,
	})
	if err != nil {
		return nil, err
	}

	if record.Status == enums.BulkPurchaseOrderSellerStatusApproved {
		return nil, errs.ErrBulkPurchaseOrderSellerQuotationInvalidToSubmit
	}

	var updates models.BulkPurchaseOrderSellerQuotation
	err = copier.Copy(&updates, &form)
	if err != nil {
		return nil, err
	}

	updates.QuotationAt = values.Int64(time.Now().Unix())
	updates.Status = enums.BulkPurchaseOrderSellerStatusWaitingForApproval
	updates.QuotedPrice = updates.FabricCost.AddPtr(updates.MakingCost).AddPtr(updates.DecorationCost).AddPtr(updates.OtherCost).ToPtr()

	if updates.BulkQuotations != nil {
		updates.BulkQuotations = lo.Map(updates.BulkQuotations, func(item *models.SellerBulkQuotationMOQ, index int) *models.SellerBulkQuotationMOQ {
			item.UnitPrice = updates.QuotedPrice.Add(item.UpCharge)
			item.TotalPrice = item.UnitPrice.MultipleInt(values.Int64Value(item.Quantity))
			updates.QuotedPrice = &item.UnitPrice
			return item
		})

	}

	err = r.db.Transaction(func(tx *gorm.DB) error {
		var sqlResult = tx.Clauses(clause.OnConflict{DoNothing: true}).
			Where("id = ? AND user_id = ?", form.SellerQuotationID, form.GetUserID()).
			Updates(updates)
		if sqlResult.Error != nil {
			return sqlResult.Error
		}

		if sqlResult.RowsAffected == 0 {
			return errs.ErrBulkPurchaseOrderSellerQuotationNotFound
		}

		var log = models.BulkPurchaseOrderTracking{
			PurchaseOrderID: record.BulkPurchaseOrderID,
			ActionType:      enums.BulkPoTrackingActionSellerReSubmitQuotation,
			UserID:          record.UserID,
			UserGroup:       enums.PoTrackingUserGroupSeller,
			CreatedByUserID: form.JwtClaimsInfo.GetUserID(),
			Metadata: &models.PoTrackingMetadata{
				Before: map[string]interface{}{
					"seller_quotation_status": record.Status,
				},
				After: map[string]interface{}{
					"seller_quotation_status": updates.Status,
				},
			},
		}

		return tx.Create(&log).Error
	})

	record.QuotationAt = updates.QuotationAt
	record.Status = updates.Status
	return record, err
}

type SellerSubmitMultipleBulkQuotationsParams struct {
	models.JwtClaimsInfo
	Quotations []*SellerSubmitBulkQuotationParams `json:"quotations" validate:"required"`
}

func (r *SellerBulkPurchaseOrderRepo) SubmitMultipleBulkQuotations(req *SellerSubmitMultipleBulkQuotationsParams) (models.BulkPurchaseOrderSellerQuotations, error) {
	var bulkSellerQtIDs = make([]string, 0, len(req.Quotations))
	for _, quotation := range req.Quotations {
		bulkSellerQtIDs = append(bulkSellerQtIDs, quotation.SellerQuotationID)
	}
	var bulkSellerQts models.BulkPurchaseOrderSellerQuotations
	if err := r.db.Find(&bulkSellerQts, "id IN ? AND user_id = ?", bulkSellerQtIDs, req.GetUserID()).Error; err != nil {
		return nil, err
	}
	var dbBulkSellerQtIDs = bulkSellerQts.IDs()
	for _, id := range bulkSellerQtIDs {
		if !helper.StringContains(dbBulkSellerQtIDs, id) {
			return nil, eris.Wrapf(errs.ErrBulkPurchaseOrderSellerQuotationNotFound, "bulk_seller_quotation_id:%s", id)
		}
	}

	for _, bsq := range bulkSellerQts {
		if bsq.Status == enums.BulkPurchaseOrderSellerStatusApproved {
			return nil, eris.Wrapf(errs.ErrBulkPurchaseOrderSellerQuotationInvalidToSubmit, "bulk_seller_quotation_id:%s", bsq.ID)
		}
	}

	var bulkSellerQtsToUpdate = make(models.BulkPurchaseOrderSellerQuotations, 0, len(req.Quotations))
	var trackingLogsToCreate = make([]models.BulkPurchaseOrderTracking, 0, len(req.Quotations))

	for _, quotation := range req.Quotations {
		bulkSellerQt, _ := lo.Find(bulkSellerQts, func(item *models.BulkPurchaseOrderSellerQuotation) bool {
			return item.ID == quotation.SellerQuotationID
		})
		var prevSellerQuotationStatus = bulkSellerQt.Status

		bulkSellerQt.FabricCost = quotation.FabricCost
		bulkSellerQt.DecorationCost = quotation.DecorationCost
		bulkSellerQt.MakingCost = quotation.MakingCost
		bulkSellerQt.OtherCost = quotation.OtherCost
		bulkSellerQt.SellerRemark = quotation.SellerRemark
		bulkSellerQt.BulkQuotations = quotation.BulkQuotations
		bulkSellerQt.QuotationAt = values.Int64(time.Now().Unix())
		bulkSellerQt.Status = enums.BulkPurchaseOrderSellerStatusWaitingForApproval

		bulkSellerQt.QuotedPrice = bulkSellerQt.FabricCost.AddPtr(bulkSellerQt.MakingCost).AddPtr(bulkSellerQt.DecorationCost).AddPtr(bulkSellerQt.OtherCost).ToPtr()

		if bulkSellerQt.BulkQuotations != nil {
			bulkSellerQt.BulkQuotations = lo.Map(bulkSellerQt.BulkQuotations, func(item *models.SellerBulkQuotationMOQ, index int) *models.SellerBulkQuotationMOQ {
				item.UnitPrice = bulkSellerQt.QuotedPrice.Add(item.UpCharge)
				item.TotalPrice = item.UnitPrice.MultipleInt(values.Int64Value(item.Quantity))
				bulkSellerQt.QuotedPrice = &item.UnitPrice
				return item
			})

		}

		bulkSellerQtsToUpdate = append(bulkSellerQtsToUpdate, bulkSellerQt)

		trackingLogsToCreate = append(trackingLogsToCreate, models.BulkPurchaseOrderTracking{
			PurchaseOrderID: bulkSellerQt.BulkPurchaseOrderID,
			ActionType:      enums.BulkPoTrackingActionSellerSubmitQuotation,
			UserGroup:       enums.PoTrackingUserGroupSeller,
			UserID:          bulkSellerQt.UserID,
			CreatedByUserID: req.GetUserID(),
			Metadata: &models.PoTrackingMetadata{
				Before: map[string]interface{}{
					"seller_quotation_status": prevSellerQuotationStatus,
				},
				After: map[string]interface{}{
					"seller_quotation_status": bulkSellerQt.Status,
				},
			},
		})
	}

	if err := r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Clauses(clause.OnConflict{Columns: []clause.Column{{Name: "id"}}, UpdateAll: true}).
			Create(bulkSellerQtsToUpdate).Error; err != nil {
			return eris.Wrap(err, err.Error())
		}

		return tx.Create(trackingLogsToCreate).Error
	}); err != nil {
		return nil, err
	}

	return bulkSellerQtsToUpdate, nil
}

type AdminBulkPurchaseOrderFirstPayoutParams struct {
	models.JwtClaimsInfo
	BulkPurchaseOrderID string `json:"bulk_purchase_order_id" param:"bulk_purchase_order_id" validate:"required"`

	TransactionRefID      string             `json:"transaction_ref_id" validate:"required_unless=PayoutPercentage 0"`
	TransactionAttachment *models.Attachment `json:"transaction_attachment" validate:"required_unless=PayoutPercentage 0"`
	PayoutPercentage      float64            `json:"payout_percentage,omitempty" validate:"gte=0,lte=100"`
}

func (r *SellerBulkPurchaseOrderRepo) AdminBulkPurchaseOrderFirstPayout(params AdminBulkPurchaseOrderFirstPayoutParams) (*models.BulkPurchaseOrder, error) {
	bulkPO, err := NewBulkPurchaseOrderRepo(r.db).GetBulkPurchaseOrder(GetBulkPurchaseOrderParams{
		JwtClaimsInfo:          params.JwtClaimsInfo,
		BulkPurchaseOrderID:    params.BulkPurchaseOrderID,
		IncludeItems:           true,
		IncludeSellerQuotation: true,
	})
	if err != nil {
		return nil, err
	}

	if bulkPO.SellerID == "" || bulkPO.SellerID == "inflow" {
		return nil, errs.ErrSellerInvalidToPayout
	}

	if bulkPO.SellerBulkQuotation == nil || len(bulkPO.SellerBulkQuotation.BulkQuotations) == 0 {
		return nil, errs.ErrSellerInvalidToPayout.WithDetail("No quotation")
	}

	var updates = models.BulkPurchaseOrder{
		SellerFirstPayoutTransactionRefID:      params.TransactionRefID,
		SellerFirstPayoutTransactionAttachment: params.TransactionAttachment,
		SellerFirstPayoutTransferedAt:          values.Int64(time.Now().Unix()),
		SellerFirstPayoutMarkAsPaidAt:          values.Int64(time.Now().Unix()),
		SellerFirstPayoutPercentage:            &params.PayoutPercentage,
		SellerTrackingStatus:                   enums.SellerBulkPoTrackingStatusFirstPaymentConfirm,
		SellerPayoutTotalAmount:                bulkPO.SellerBulkQuotation.BulkQuotations[0].TotalPrice.ToPtr(),
	}
	updates.SellerFirstPayoutTotalAmount = updates.SellerPayoutTotalAmount.MultipleFloat64(params.PayoutPercentage).DivInt(100).ToPtr()
	updates.SellerFinalPayoutTotalAmount = updates.SellerPayoutTotalAmount.SubPtr(updates.SellerFirstPayoutTotalAmount).ToPtr()

	err = r.db.Transaction(func(tx *gorm.DB) error {
		if params.PayoutPercentage > 0 {
			var transaction = models.PaymentTransaction{
				BulkPurchaseOrderID: bulkPO.ID,
				PaymentType:         enums.PaymentTypeBankTransfer,
				Milestone:           enums.PaymentMilestoneFirstPayment,
				UserID:              bulkPO.SellerID,
				TransactionRefID:    params.TransactionRefID,
				Status:              enums.PaymentStatusPaid,
				PaymentPercentage:   &params.PayoutPercentage,
				PaidAmount:          updates.SellerFirstPayoutTotalAmount,
				TotalAmount:         updates.SellerPayoutTotalAmount,
				Currency:            bulkPO.Currency,
				TransactionType:     enums.TransactionTypeDebit,
				Attachments: &models.Attachments{
					params.TransactionAttachment,
				},
				Metadata: &models.PaymentTransactionMetadata{
					BulkPurchaseOrderID: bulkPO.ID,
				},
			}
			err = tx.Create(&transaction).Error
			if err != nil {
				return err
			}

			updates.SellerFirstPayoutTransactionReferenceID = transaction.ReferenceID

		} else {
			updates.SellerTrackingStatus = enums.SellerBulkPoTrackingStatusFirstPaymentSkipped
		}

		var sqlResult = tx.Model(&models.BulkPurchaseOrder{}).Where("id = ?", bulkPO.ID).Updates(&updates)
		if sqlResult.Error != nil {
			return sqlResult.Error
		}

		if sqlResult.RowsAffected == 0 {
			return errs.ErrBulkPoNotFound
		}
		return sqlResult.Error

	})
	if err != nil {
		return nil, err
	}

	return bulkPO, err
}

type AdminBulkPurchaseOrderFinalPayoutParams struct {
	models.JwtClaimsInfo
	BulkPurchaseOrderID string `json:"bulk_purchase_order_id" param:"bulk_purchase_order_id" validate:"required"`

	TransactionRefID      string             `json:"transaction_ref_id" validate:"required"`
	TransactionAttachment *models.Attachment `json:"transaction_attachment" validate:"required"`
}

func (r *SellerBulkPurchaseOrderRepo) AdminBulkPurchaseOrderFinalPayout(params AdminBulkPurchaseOrderFinalPayoutParams) (*models.BulkPurchaseOrder, error) {
	bulkPO, err := NewBulkPurchaseOrderRepo(r.db).GetBulkPurchaseOrder(GetBulkPurchaseOrderParams{
		JwtClaimsInfo:          params.JwtClaimsInfo,
		BulkPurchaseOrderID:    params.BulkPurchaseOrderID,
		IncludeItems:           true,
		IncludeSellerQuotation: true,
	})
	if err != nil {
		return nil, err
	}

	if bulkPO.SellerID == "" || bulkPO.SellerID == "inflow" {
		return nil, errs.ErrSellerInvalidToPayout
	}

	if bulkPO.SellerBulkQuotation == nil || len(bulkPO.SellerBulkQuotation.BulkQuotations) == 0 {
		return nil, errs.ErrSellerInvalidToPayout.WithDetail("No quotation")
	}

	var updates = models.BulkPurchaseOrder{
		SellerFinalPayoutTransactionRefID:      params.TransactionRefID,
		SellerFinalPayoutTransactionAttachment: params.TransactionAttachment,
		SellerFinalPayoutTransferedAt:          values.Int64(time.Now().Unix()),
		SellerFinalPayoutMarkAsPaidAt:          values.Int64(time.Now().Unix()),
		SellerTrackingStatus:                   enums.SellerBulkPoTrackingStatusFinalPaymentConfirm,
	}

	if bulkPO.SellerPayoutTotalAmount != nil {
		updates.SellerFinalPayoutTotalAmount = bulkPO.SellerPayoutTotalAmount.SubPtr(bulkPO.SellerFirstPayoutTotalAmount).ToPtr()
	} else {
		updates.SellerPayoutTotalAmount = bulkPO.SellerBulkQuotation.BulkQuotations[0].TotalPrice.ToPtr()
		updates.SellerFinalPayoutTotalAmount = updates.SellerPayoutTotalAmount.SubPtr(bulkPO.SellerFirstPayoutTotalAmount).ToPtr()

	}

	err = r.db.Transaction(func(tx *gorm.DB) error {
		var finalPayoutPercentage = 100 - values.Float32Value(bulkPO.SellerFirstPayoutPercentage)
		if finalPayoutPercentage > 0 {
			var transaction = models.PaymentTransaction{
				BulkPurchaseOrderID: bulkPO.ID,
				PaymentType:         enums.PaymentTypeBankTransfer,
				Milestone:           enums.PaymentMilestoneFinalPayment,
				UserID:              bulkPO.SellerID,
				TransactionRefID:    params.TransactionRefID,
				Status:              enums.PaymentStatusPaid,
				PaymentPercentage:   values.Float64(finalPayoutPercentage),
				PaidAmount:          updates.SellerFinalPayoutTotalAmount,
				TotalAmount:         updates.SellerPayoutTotalAmount,
				Currency:            bulkPO.Currency,
				TransactionType:     enums.TransactionTypeDebit,
				Attachments: &models.Attachments{
					params.TransactionAttachment,
				},
				Metadata: &models.PaymentTransactionMetadata{
					BulkPurchaseOrderID: bulkPO.ID,
				},
			}
			err = tx.Create(&transaction).Error
			if err != nil {
				return err
			}

			updates.SellerFinalPayoutTransactionReferenceID = transaction.ReferenceID

		}

		var sqlResult = tx.Model(&models.BulkPurchaseOrder{}).Where("id = ?", bulkPO.ID).Updates(&updates)
		if sqlResult.Error != nil {
			return sqlResult.Error
		}

		if sqlResult.RowsAffected == 0 {
			return errs.ErrBulkPoNotFound
		}
		return sqlResult.Error

	})
	if err != nil {
		return nil, err
	}

	return bulkPO, err
}

type SellerBulkPurchaseOrderUpdateRawMaterialParams struct {
	models.JwtClaimsInfo

	BulkPurchaseOrderID  string                     `json:"bulk_purchase_order_id" param:"bulk_purchase_order_id" query:"bulk_purchase_order_id" validate:"required"`
	PoRawMaterials       *models.PoRawMaterialMetas `json:"po_raw_materials" param:"po_raw_materials" query:"po_raw_materials"`
	ApproveRawMaterialAt *int64                     `json:"approve_raw_material_at" param:"approve_raw_material_at" query:"approve_raw_material_at"`
}

func (r *SellerBulkPurchaseOrderRepo) SellerBulkPurchaseOrderUpdateRawMaterial(params SellerBulkPurchaseOrderUpdateRawMaterialParams) (*models.BulkPurchaseOrder, error) {
	order, err := NewBulkPurchaseOrderRepo(r.db).GetBulkPurchaseOrder(GetBulkPurchaseOrderParams{
		JwtClaimsInfo:       params.JwtClaimsInfo,
		BulkPurchaseOrderID: params.BulkPurchaseOrderID,
	})
	if err != nil {
		return nil, err
	}

	var updates = models.BulkPurchaseOrder{
		SellerPoRawMaterials: params.PoRawMaterials,
		ApproveRawMaterialAt: params.ApproveRawMaterialAt,
	}

	_ = updates.GenerateRawMaterialRefID(updates.SellerPoRawMaterials)

	updates.SellerTrackingStatus = enums.SellerBulkPoTrackingStatusRawMaterial
	err = r.db.Transaction(func(tx *gorm.DB) error {
		err = NewBulkPurchaseOrderTrackingRepo(r.db).CreateBulkPurchaseOrderTrackingTx(tx, models.BulkPurchaseOrderTrackingCreateForm{
			PurchaseOrderID: order.ID,
			ActionType:      enums.BulkPoTrackingActionUpdateMaterial,
			UserID:          params.GetUserID(),
			UserGroup:       enums.PoTrackingUserGroupSeller,
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

	order.SellerTrackingStatus = updates.SellerTrackingStatus

	return order, err
}

type SellerBulkPurchaseOrderUpdatePpsParams struct {
	models.JwtClaimsInfo

	BulkPurchaseOrderID string            `json:"bulk_purchase_order_id" param:"bulk_purchase_order_id" query:"bulk_purchase_order_id" validate:"required"`
	PpsInfo             *models.PoPpsMeta `json:"pps_info" validate:"required"`
}

func (r *SellerBulkPurchaseOrderRepo) SellerBulkPurchaseOrderUpdatePps(params SellerBulkPurchaseOrderUpdatePpsParams) (*models.BulkPurchaseOrder, error) {
	order, err := NewBulkPurchaseOrderRepo(r.db).GetBulkPurchaseOrder(GetBulkPurchaseOrderParams{
		BulkPurchaseOrderID: params.BulkPurchaseOrderID,
		JwtClaimsInfo:       params.JwtClaimsInfo,
	})
	if err != nil {
		return nil, err
	}
	var ppsInfoArr models.PoPpsMetas
	if order.SellerPpsInfo != nil {
		ppsInfoArr = append(ppsInfoArr, *order.SellerPpsInfo...)
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
		SellerPpsInfo:        &ppsInfoArr,
		SellerTrackingStatus: enums.SellerBulkPoTrackingStatusPps,
	}

	err = r.db.Transaction(func(tx *gorm.DB) error {
		NewBulkPurchaseOrderTrackingRepo(r.db).CreateBulkPurchaseOrderTrackingTx(tx, models.BulkPurchaseOrderTrackingCreateForm{
			PurchaseOrderID: order.ID,
			ActionType:      enums.BulkPoTrackingActionSellerUpdatePps,
			UserGroup:       enums.PoTrackingUserGroupSeller,
			UserID:          params.GetUserID(),
			CreatedByUserID: params.JwtClaimsInfo.GetUserID(),
			Metadata: &models.PoTrackingMetadata{
				Before: map[string]interface{}{
					"seller_pps_info": order.PpsInfo,
				},
				After: map[string]interface{}{
					"seller_pps_info": updates.PpsInfo,
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

	order.SellerTrackingStatus = updates.SellerTrackingStatus
	order.PpsInfo = updates.PpsInfo

	return order, err
}

type SellerBulkPurchaseOrderUpdateProductionParams struct {
	models.JwtClaimsInfo

	BulkPurchaseOrderID string                   `json:"bulk_purchase_order_id" param:"bulk_purchase_order_id" query:"bulk_purchase_order_id" validate:"required"`
	ProductionInfo      *models.PoProductionMeta `json:"production_info" param:"production_info" query:"production_info"`
}

func (r *SellerBulkPurchaseOrderRepo) SellerBulkPurchaseOrderUpdateProduction(params SellerBulkPurchaseOrderUpdateProductionParams) (*models.BulkPurchaseOrder, error) {
	order, err := NewBulkPurchaseOrderRepo(r.db).GetBulkPurchaseOrder(GetBulkPurchaseOrderParams{
		JwtClaimsInfo:       params.JwtClaimsInfo,
		BulkPurchaseOrderID: params.BulkPurchaseOrderID,
		IncludeUser:         true,
	})
	if err != nil {
		return nil, err
	}

	var updates = models.BulkPurchaseOrder{
		SellerProductionInfo: params.ProductionInfo,
	}

	if order.SellerTrackingStatus == enums.SellerBulkPoTrackingStatusPps {
		updates.SellerTrackingStatus = enums.SellerBulkPoTrackingStatusProduction
	}

	err = r.db.Transaction(func(tx *gorm.DB) error {
		err = NewBulkPurchaseOrderTrackingRepo(r.db).CreateBulkPurchaseOrderTrackingTx(tx, models.BulkPurchaseOrderTrackingCreateForm{
			PurchaseOrderID: order.ID,
			ActionType:      enums.BulkPoTrackingActionUpdateProduction,
			UserGroup:       enums.PoTrackingUserGroupSeller,
			UserID:          params.GetUserID(),
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
		return tx.Model(&models.BulkPurchaseOrder{}).
			Where("id = ?", params.BulkPurchaseOrderID).
			Updates(&updates).Error
	})
	if err != nil {
		return nil, eris.Wrap(err, err.Error())
	}

	order.SellerTrackingStatus = updates.SellerTrackingStatus

	return order, err
}

type SellerCreateQcReportParams struct {
	models.JwtClaimsInfo

	BulkPurchaseOrderID string               `json:"bulk_purchase_order_id" param:"bulk_purchase_order_id" query:"bulk_purchase_order_id" validate:"required"`
	PoQcReports         models.PoReportMetas `json:"po_qc_reports" param:"po_qc_reports" query:"po_qc_reports" validate:"required"`
	ApproveQCAt         *int64               `json:"approve_qc_at" param:"approve_qc_at" query:"approve_qc_at" `
}

func (r *SellerBulkPurchaseOrderRepo) SellerBulkPurchaseOrderCreateQcReport(params SellerCreateQcReportParams) (*models.BulkPurchaseOrder, error) {
	order, err := NewBulkPurchaseOrderRepo(r.db).GetBulkPurchaseOrder(GetBulkPurchaseOrderParams{
		BulkPurchaseOrderID: params.BulkPurchaseOrderID,
		JwtClaimsInfo:       params.JwtClaimsInfo,
		IncludeUser:         true,
	})

	if err != nil {
		return nil, eris.Wrap(err, err.Error())
	}

	var updates = models.BulkPurchaseOrder{
		SellerTrackingStatus: enums.SellerBulkPoTrackingStatusQc,
		SellerPoQcReports:    &params.PoQcReports,
	}
	var trackings = lo.Map(params.PoQcReports, func(item *models.PoReportMeta, index int) *models.BulkPurchaseOrderTracking {
		var tracking = models.BulkPurchaseOrderTracking{
			PurchaseOrderID: order.ID,
			ActionType:      enums.BulkPoTrackingActionCreateQcReport,
			UserID:          params.GetUserID(),
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

		return &tracking
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

	order.SellerTrackingStatus = updates.SellerTrackingStatus
	order.SellerPoQcReports = updates.SellerPoQcReports

	return order, err
}

type PaginateBulkPurchaseOrderMatchingSellersParams struct {
	models.PaginationParams
	models.JwtClaimsInfo

	SellerID      string
	ProductGroups []string `json:"product_groups" query:"product_groups" form:"product_groups"`
	ProductTypes  []string `json:"product_types" query:"product_types" form:"product_types"`
	FabricTypes   []string `json:"fabric_types" query:"fabric_types" form:"fabric_types"`
}

func (r *SellerBulkPurchaseOrderRepo) PaginateBulkPurchaseOrderMatchingSellers(params PaginateBulkPurchaseOrderMatchingSellersParams) *query.Pagination {
	var builder = queryfunc.NewInquirySellerMatchingBuilder(queryfunc.InquirySellerMatchingBuilderOptions{
		QueryBuilderOptions: queryfunc.QueryBuilderOptions{
			Role: enums.RoleSeller,
		},
	})

	if params.Limit == 0 {
		params.Limit = 20
	}

	var result = query.New(r.db, builder).
		WhereFunc(func(builder *query.Builder) {
			builder.Where("role = ?", enums.RoleSeller)
			builder.Where("u.supplier_type = ?", enums.SupplierTypeManufacturer)

			if len(params.FabricTypes) > 0 {
				builder.Where("count_elements(bu.excepted_fabric_types,?) = 0", pq.StringArray(params.FabricTypes))
			}
			if len(params.ProductTypes) > 0 {
				builder.Where("count_elements(bu.product_types,?) > 0", pq.StringArray(params.ProductTypes))
			}
			if len(params.ProductGroups) > 0 {
				builder.Where("count_elements(bu.product_groups,?) > 0", pq.StringArray(params.ProductGroups))
			}
		}).
		Page(params.Page).
		Limit(params.Limit).
		PagingFunc()

	return result
}

type SendBulkPurchaseOrderToSellerParams struct {
	models.JwtClaimsInfo

	BulkPurchaseOrderID string `json:"bulk_purchase_order_id" param:"bulk_purchase_order_id" query:"bulk_purchase_order_id" validate:"required"`

	Sellers []*models.SellerRequestQuotationInfo `json:"sellers" validate:"required"`
}

func (r *SellerBulkPurchaseOrderRepo) SendBulkPurchaseOrderToSeller(form SendBulkPurchaseOrderToSellerParams) ([]*models.BulkPurchaseOrderSellerQuotation, error) {
	var bulkPO models.BulkPurchaseOrder
	var err = r.db.Select("ID", "Currency").First(&bulkPO, "id = ?", form.BulkPurchaseOrderID).Error
	if err != nil {
		return nil, err
	}

	var records = lo.Map(form.Sellers, func(item *models.SellerRequestQuotationInfo, index int) *models.BulkPurchaseOrderSellerQuotation {
		var record = &models.BulkPurchaseOrderSellerQuotation{
			BulkPurchaseOrderID:         bulkPO.ID,
			Currency:                    bulkPO.Currency,
			UserID:                      item.SellerID,
			Status:                      enums.BulkPurchaseOrderSellerStatusWaitingForQuotation,
			VarianceAmount:              item.VarianceAmount,
			VariancePercentage:          item.VariancePercentage,
			AdminSentAt:                 values.Int64(time.Now().Unix()),
			OfferPrice:                  item.OfferPrice,
			OfferRemark:                 item.OfferRemark,
			OrderType:                   item.OrderType,
			ExpectedStartProductionDate: item.ExpectedStartProductionDate,
		}

		return record
	})

	err = r.db.Clauses(clause.OnConflict{
		DoNothing: true,
		Columns: []clause.Column{
			{Name: "user_id"},
			{Name: "bulk_purchase_order_id"},
		},
	}).Create(&records).Error
	if err != nil {
		return nil, err
	}

	return records, nil
}

type PaginateBulkPurchaseOrderSellerQuotationsParams struct {
	models.PaginationParams
	models.JwtClaimsInfo

	BulkPurchaseOrderID string `json:"bulk_purchase_order_id" query:"bulk_purchase_order_id" form:"bulk_purchase_order_id" param:"bulk_purchase_order_id"`

	UserID string `json:"user_id" query:"user_id" form:"user_id"`

	DateFrom         int64  `json:"date_from" query:"date_from" form:"date_from"`
	DateTo           int64  `json:"date_to" query:"date_to" form:"date_to"`
	OrderReferenceID string `json:"order_reference_id" query:"order_reference_id" form:"order_reference_id"`

	Statuses []enums.BulkPurchaseOrderSellerStatus `json:"statuses" query:"statuses" form:"statuses"`
}

func (r *SellerBulkPurchaseOrderRepo) PaginateBulkPurchaseOrderSellerQuotations(params PaginateBulkPurchaseOrderSellerQuotationsParams) *query.Pagination {
	var builder = queryfunc.NewBulkPurchaseOrderSellerQuotationBuilder(queryfunc.BulkPurchaseOrderSellerQuotationBuilderOptions{
		QueryBuilderOptions: queryfunc.QueryBuilderOptions{
			Role: params.GetRole(),
		},
		CurrentUserID: params.GetUserID(),
	})

	var result = query.New(r.db, builder).
		WhereFunc(func(builder *query.Builder) {
			if params.BulkPurchaseOrderID != "" {
				builder.Where("rq.bulk_purchase_order_id = ?", params.BulkPurchaseOrderID)
			}

			if params.UserID != "" {
				builder.Where("rq.user_id = ?", params.UserID)
			}
			if params.DateFrom > 0 {
				builder.Where("rq.created_at >= ?", params.DateFrom)
			}
			if params.DateTo > 0 {
				builder.Where("rq.created_at <= ?", params.DateTo)
			}

			if params.Statuses != nil {
				builder.Where("rq.status IN ?", params.Statuses)
			}

			if strings.TrimSpace(params.OrderReferenceID) != "" {
				var q = "%" + params.OrderReferenceID + "%"
				builder.Where("rq.order_reference_id ILIKE @query_po", sql.Named("query_po", q))
			}
		}).
		Page(params.Page).
		Limit(params.Limit).
		PagingFunc()

	return result
}

type UpdateBulkPurchaseOrderProductPhotoParams struct {
	models.JwtClaimsInfo

	BulkPurchaseOrderID string             `json:"bulk_purchase_order_id" param:"bulk_purchase_order_id" validate:"required"`
	ProductAttachments  models.Attachments `json:"product_attachments"`
}

func (r *SellerBulkPurchaseOrderRepo) UpdateBulkPurchaseOrderProductPhoto(params UpdateBulkPurchaseOrderProductPhotoParams) (*models.BulkPurchaseOrder, error) {
	var updates = models.BulkPurchaseOrder{
		SellerProductAttachments: &params.ProductAttachments,
	}

	var err = r.db.Model(&models.BulkPurchaseOrder{}).Where("id = ?", params.BulkPurchaseOrderID).Updates(&updates).Error
	if err != nil {
		return nil, err
	}
	return &updates, err
}

type UpdateBulkPurchaseOrderTechpackParams struct {
	models.JwtClaimsInfo

	BulkPurchaseOrderID string             `json:"bulk_purchase_order_id" param:"bulk_purchase_order_id" validate:"required"`
	TechpackAttachments models.Attachments `json:"techpack_attachments"`
}

func (r *SellerBulkPurchaseOrderRepo) UpdateBulkPurchaseOrderTechpack(params UpdateBulkPurchaseOrderTechpackParams) (*models.BulkPurchaseOrder, error) {
	var updates = models.BulkPurchaseOrder{
		SellerTechpackAttachments: &params.TechpackAttachments,
	}

	var err = r.db.Model(&models.BulkPurchaseOrder{}).Where("id = ?", params.BulkPurchaseOrderID).Updates(&updates).Error
	if err != nil {
		return nil, err
	}
	return &updates, err
}

type UpdateBulkPurchaseOrderSizeChartParams struct {
	models.JwtClaimsInfo

	BulkPurchaseOrderID  string             `json:"bulk_purchase_order_id" param:"bulk_purchase_order_id" validate:"required"`
	SizeChartAttachments models.Attachments `json:"size_chart_attachments"`
}

func (r *SellerBulkPurchaseOrderRepo) UpdateBulkPurchaseOrderSizeChart(params UpdateBulkPurchaseOrderSizeChartParams) (*models.BulkPurchaseOrder, error) {
	var updates = models.BulkPurchaseOrder{
		SellerSizeChartAttachments: &params.SizeChartAttachments,
	}

	var err = r.db.Model(&models.BulkPurchaseOrder{}).Where("id = ?", params.BulkPurchaseOrderID).Updates(&updates).Error
	if err != nil {
		return nil, err
	}
	return &updates, err
}

type UpdateBulkPurchaseOrderSizeSpecParams struct {
	models.JwtClaimsInfo

	BulkPurchaseOrderID string             `json:"bulk_purchase_order_id" param:"bulk_purchase_order_id" validate:"required"`
	SizeSpecAttachments models.Attachments `json:"size_spec_attachments"`
}

func (r *SellerBulkPurchaseOrderRepo) UpdateBulkPurchaseOrderSizeSpec(params UpdateBulkPurchaseOrderSizeSpecParams) (*models.BulkPurchaseOrder, error) {
	var updates = models.BulkPurchaseOrder{
		SellerSizeSpecAttachments: &params.SizeSpecAttachments,
	}

	var err = r.db.Model(&models.BulkPurchaseOrder{}).Where("id = ?", params.BulkPurchaseOrderID).Updates(&updates).Error
	if err != nil {
		return nil, err
	}
	return &updates, err
}

type UpdateBulkPurchaseOrderSizeGradingParams struct {
	models.JwtClaimsInfo

	BulkPurchaseOrderID    string             `json:"bulk_purchase_order_id" param:"bulk_purchase_order_id" validate:"required"`
	SizeGradingAttachments models.Attachments `json:"size_grading_attachments"`
}

func (r *SellerBulkPurchaseOrderRepo) UpdateBulkPurchaseOrderSizeGrading(params UpdateBulkPurchaseOrderSizeGradingParams) (*models.BulkPurchaseOrder, error) {
	var updates = models.BulkPurchaseOrder{
		SellerSizeGradingAttachments: &params.SizeGradingAttachments,
	}

	var err = r.db.Model(&models.BulkPurchaseOrder{}).Where("id = ?", params.BulkPurchaseOrderID).Updates(&updates).Error
	if err != nil {
		return nil, err
	}
	return &updates, err
}

type UpdateBulkPurchaseOrderBillOfMaterialParams struct {
	models.JwtClaimsInfo

	BulkPurchaseOrderID       string             `json:"bulk_purchase_order_id" param:"bulk_purchase_order_id" validate:"required"`
	BillOfMaterialAttachments models.Attachments `json:"bill_of_material_attachments"`
}

func (r *SellerBulkPurchaseOrderRepo) UpdateBulkPurchaseOrderBillOfMaterial(params UpdateBulkPurchaseOrderBillOfMaterialParams) (*models.BulkPurchaseOrder, error) {
	var updates = models.BulkPurchaseOrder{
		SellerBillOfMaterialAttachments: &params.BillOfMaterialAttachments,
	}

	var err = r.db.Model(&models.BulkPurchaseOrder{}).Where("id = ?", params.BulkPurchaseOrderID).Updates(&updates).Error
	if err != nil {
		return nil, err
	}
	return &updates, err
}

type AdminSellerBulkPurchaseOrderUpdateInspectionProcedureParams struct {
	models.JwtClaimsInfo

	BulkPurchaseOrderID            string             `json:"bulk_purchase_order_id" param:"bulk_purchase_order_id" validate:"required"`
	InspectionProcedureAttachments models.Attachments `json:"inspection_procedure_attachments"`
	InspectionProcedureNote        string             `json:"inspection_procedure_note"`
}

func (r *SellerBulkPurchaseOrderRepo) AdminSellerBulkPurchaseOrderUpdateInspectionProcedure(params AdminSellerBulkPurchaseOrderUpdateInspectionProcedureParams) (*models.BulkPurchaseOrder, error) {
	var updates = models.BulkPurchaseOrder{
		SellerInspectionProcedureAttachments: &params.InspectionProcedureAttachments,
		SellerInspectionProcedureNote:        params.InspectionProcedureNote,
	}

	var err = r.db.Model(&models.BulkPurchaseOrder{}).Where("id = ?", params.BulkPurchaseOrderID).Updates(&updates).Error
	if err != nil {
		return nil, err
	}
	return &updates, err
}

type UpdateBulkPurchaseOrderTestingRequirementsParams struct {
	models.JwtClaimsInfo

	BulkPurchaseOrderID            string             `json:"bulk_purchase_order_id" param:"bulk_purchase_order_id" validate:"required"`
	TestingRequirementsAttachments models.Attachments `json:"testing_requirements_attachments"`
	TestingRequirementsNote        string             `json:"testing_requirements_note"`
}

func (r *SellerBulkPurchaseOrderRepo) UpdateBulkPurchaseOrderTestingRequirements(params UpdateBulkPurchaseOrderTestingRequirementsParams) (*models.BulkPurchaseOrder, error) {
	var updates = models.BulkPurchaseOrder{
		SellerInspectionTestingRequirementsAttachments: &params.TestingRequirementsAttachments,
		SellerInspectionTestingRequirementsNote:        params.TestingRequirementsNote,
	}

	var err = r.db.Model(&models.BulkPurchaseOrder{}).Where("id = ?", params.BulkPurchaseOrderID).Updates(&updates).Error
	if err != nil {
		return nil, err
	}
	return &updates, err
}

type UpdateBulkPurchaseOrderPackingParams struct {
	models.JwtClaimsInfo

	BulkPurchaseOrderID string             `json:"bulk_purchase_order_id" param:"bulk_purchase_order_id" validate:"required"`
	PackingAttachments  models.Attachments `json:"packing_attachments"`
	PackingNote         string             `json:"packing_note"`
}

func (r *SellerBulkPurchaseOrderRepo) UpdateBulkPurchaseOrderPacking(params UpdateBulkPurchaseOrderPackingParams) (*models.BulkPurchaseOrder, error) {
	var updates = models.BulkPurchaseOrder{
		SellerPackingAttachments: &params.PackingAttachments,
		SellerPackingNote:        params.PackingNote,
	}

	var err = r.db.Model(&models.BulkPurchaseOrder{}).Where("id = ?", params.BulkPurchaseOrderID).Updates(&updates).Error
	if err != nil {
		return nil, err
	}
	return &updates, err
}

type AdminApproveSellerBulkPurchaseOrderQuotationParams struct {
	models.JwtClaimsInfo

	SellerQuotationID string `json:"seller_quotation_id" param:"seller_quotation_id" validate:"required"`
}

func (r *SellerBulkPurchaseOrderRepo) AdminApproveSellerBulkPurchaseOrderQuotation(params AdminApproveSellerBulkPurchaseOrderQuotationParams) error {
	var admin models.User
	var err = r.db.Select("ID", "Name", "Email").First(&admin, "id = ?", params.GetUserID()).Error
	if err != nil {
		return err
	}

	var sellerQuotation models.BulkPurchaseOrderSellerQuotation
	err = r.db.Select("ID", "Status", "BulkPurchaseOrderID", "UserID").First(&sellerQuotation, "id = ?", params.SellerQuotationID).Error
	if err != nil {
		return err
	}

	if sellerQuotation.Status != enums.BulkPurchaseOrderSellerStatusWaitingForApproval {
		return errs.ErrInquiryInvalidToSendQuotationToBuyer
	}

	err = r.db.Transaction(func(tx *gorm.DB) error {
		var quotationUpdates = models.BulkPurchaseOrderSellerQuotation{
			Status: enums.BulkPurchaseOrderSellerStatusApproved,
		}
		err = tx.Model(&models.BulkPurchaseOrderSellerQuotation{}).Where("id = ?", params.SellerQuotationID).Updates(&quotationUpdates).Error
		if err != nil {
			return err
		}

		var otherUpdates = models.BulkPurchaseOrderSellerQuotation{
			Status: enums.BulkPurchaseOrderSellerStatusRejected,
		}
		err = tx.Model(&models.BulkPurchaseOrderSellerQuotation{}).
			Where("id <> ? AND bulk_purchase_order_id = ?", params.SellerQuotationID, sellerQuotation.BulkPurchaseOrderID).
			Updates(&otherUpdates).Error
		if err != nil {
			return err
		}

		var bulkUpdates = models.BulkPurchaseOrder{
			SellerTrackingStatus: enums.SellerBulkPoTrackingStatusPO,
			SellerID:             sellerQuotation.UserID,
		}

		var log = models.BulkPurchaseOrderTracking{
			PurchaseOrderID: sellerQuotation.BulkPurchaseOrderID,
			ActionType:      enums.BulkPoTrackingActionAdminApproveSellerQuotation,
			UserID:          sellerQuotation.UserID,
			UserGroup:       enums.PoTrackingUserGroupSeller,
			CreatedByUserID: params.JwtClaimsInfo.GetUserID(),
			Metadata: &models.PoTrackingMetadata{
				Before: map[string]interface{}{
					"seller_quotation_status": sellerQuotation.Status,
				},
				After: map[string]interface{}{
					"seller_quotation_status": quotationUpdates.Status,
				},
			},
		}
		err = tx.Create(&log).Error
		if err != nil {
			return err
		}

		return tx.Model(&models.BulkPurchaseOrder{}).Where("id = ?", sellerQuotation.BulkPurchaseOrderID).Updates(&bulkUpdates).Error
	})

	return err

}

type AdminRejectSellerBulkPurchaseOrderQuotationParams struct {
	models.JwtClaimsInfo

	SellerQuotationID string `json:"seller_quotation_id" param:"seller_quotation_id" validate:"required"`
	RejectReason      string `json:"reject_reason"`
}

func (r *SellerBulkPurchaseOrderRepo) AdminRejectSellerBulkPurchaseOrderQuotation(params AdminRejectSellerBulkPurchaseOrderQuotationParams) error {
	var admin models.User
	var err = r.db.Select("ID", "Name", "Email").First(&admin, "id = ?", params.GetUserID()).Error
	if err != nil {
		return err
	}

	var sellerQuotation models.BulkPurchaseOrderSellerQuotation
	err = r.db.Select("ID", "Status", "BulkPurchaseOrderID", "UserID").First(&sellerQuotation, "id = ?", params.SellerQuotationID).Error
	if err != nil {
		return err
	}

	if sellerQuotation.Status != enums.BulkPurchaseOrderSellerStatusWaitingForApproval {
		return errs.ErrInquiryInvalidToSendQuotationToBuyer
	}

	var quotationUpdates = models.BulkPurchaseOrderSellerQuotation{
		Status:       enums.BulkPurchaseOrderSellerStatusRejected,
		RejectReason: params.RejectReason,
	}
	err = r.db.Transaction(func(tx *gorm.DB) error {
		err = tx.Model(&models.BulkPurchaseOrderSellerQuotation{}).Where("id = ?", params.SellerQuotationID).Updates(&quotationUpdates).Error
		if err != nil {
			return err
		}
		var log = models.BulkPurchaseOrderTracking{
			PurchaseOrderID: sellerQuotation.BulkPurchaseOrderID,
			ActionType:      enums.BulkPoTrackingActionAdminRejectSellerQuotation,
			UserID:          sellerQuotation.UserID,
			UserGroup:       enums.PoTrackingUserGroupSeller,
			CreatedByUserID: params.JwtClaimsInfo.GetUserID(),
			Metadata: &models.PoTrackingMetadata{
				Before: map[string]interface{}{
					"seller_quotation_status": sellerQuotation.Status,
				},
				After: map[string]interface{}{
					"seller_quotation_status": quotationUpdates.Status,
				},
			},
		}
		return tx.Create(&log).Error
	})
	if err != nil {
		return err
	}

	return err

}

type UpdateBulkPurchaseOrderLabelGuideAttachmentsParams struct {
	models.JwtClaimsInfo

	BulkPurchaseOrderID   string             `json:"bulk_purchase_order_id" param:"bulk_purchase_order_id" validate:"required"`
	LabelGuideAttachments models.Attachments `json:"label_guide_attachments"`
}

func (r *SellerBulkPurchaseOrderRepo) UpdateBulkPurchaseOrderLabelGuideAttachments(params UpdateBulkPurchaseOrderLabelGuideAttachmentsParams) (*models.BulkPurchaseOrder, error) {
	var updates = models.BulkPurchaseOrder{
		SellerLabelGuideAttachments: &params.LabelGuideAttachments,
	}

	var err = r.db.Model(&models.BulkPurchaseOrder{}).Where("id = ?", params.BulkPurchaseOrderID).Updates(&updates).Error
	if err != nil {
		return nil, err
	}
	return &updates, err
}

type UpdateBulkPurchaseOrderPointOfMeasurementAttachmentsParams struct {
	models.JwtClaimsInfo

	BulkPurchaseOrderID           string             `json:"bulk_purchase_order_id" param:"bulk_purchase_order_id" validate:"required"`
	PointOfMeasurementAttachments models.Attachments `json:"point_of_measurement_attachments"`
}

func (r *SellerBulkPurchaseOrderRepo) UpdateBulkPurchaseOrderPointOfMeasurementAttachments(params UpdateBulkPurchaseOrderPointOfMeasurementAttachmentsParams) (*models.BulkPurchaseOrder, error) {
	var updates = models.BulkPurchaseOrder{
		SellerPointOfMeasurementAttachments: &params.PointOfMeasurementAttachments,
	}

	var err = r.db.Model(&models.BulkPurchaseOrder{}).Where("id = ?", params.BulkPurchaseOrderID).Updates(&updates).Error
	if err != nil {
		return nil, err
	}
	return &updates, err
}

type BulkPurchaseOrderApprovePOParams struct {
	models.JwtClaimsInfo

	BulkPurchaseOrderID string `json:"bulk_purchase_order_id" param:"bulk_purchase_order_id" validate:"required"`
}

func (r *SellerBulkPurchaseOrderRepo) BulkPurchaseOrderApprovePO(params BulkPurchaseOrderApprovePOParams) (*models.BulkPurchaseOrder, error) {
	bulkPO, err := NewBulkPurchaseOrderRepo(r.db).GetBulkPurchaseOrder(GetBulkPurchaseOrderParams{
		JwtClaimsInfo:       params.JwtClaimsInfo,
		BulkPurchaseOrderID: params.BulkPurchaseOrderID,
	})
	if err != nil {
		return nil, err
	}

	if bulkPO.SellerPoAttachments != nil {
		var pendingItems = lo.Filter(*bulkPO.SellerPoAttachments, func(item *models.PoAttachment, index int) bool {
			return item.Status != enums.PoAttachmentStatusRejected
		})

		var items models.PoAttachments = lo.Map(pendingItems, func(item *models.PoAttachment, index int) *models.PoAttachment {
			item.Status = enums.PoAttachmentStatusApproved
			return item
		})

		var err = r.db.Transaction(func(tx *gorm.DB) error {
			var updates = models.BulkPurchaseOrder{
				SellerPoAttachments:  &items,
				SellerTrackingStatus: enums.SellerBulkPoTrackingStatusWaitingFirstPayment,
			}
			err = tx.Model(&models.BulkPurchaseOrder{}).Where("id = ?", params.BulkPurchaseOrderID).Updates(&updates).Error
			if err != nil {
				return err
			}

			var log = models.BulkPurchaseOrderTracking{
				PurchaseOrderID: bulkPO.ID,
				ActionType:      enums.BulkPoTrackingActionSellerApprovePO,
				UserGroup:       enums.PoTrackingUserGroupSeller,
				UserID:          bulkPO.SellerID,
				CreatedByUserID: params.JwtClaimsInfo.GetUserID(),
				Metadata: &models.PoTrackingMetadata{
					Before: map[string]interface{}{
						"seller_po_attachments":  bulkPO.SellerPoAttachments,
						"seller_tracking_status": bulkPO.SellerTrackingStatus,
					},
					After: map[string]interface{}{
						"seller_tracking_status": updates.SellerTrackingStatus,
					},
				},
			}

			return tx.Create(&log).Error
		})

		if err != nil {
			return nil, err
		}
	}

	return bulkPO, err
}

type BulkPurchaseOrderRejectPOParams struct {
	models.JwtClaimsInfo

	BulkPurchaseOrderID string `json:"bulk_purchase_order_id" param:"bulk_purchase_order_id" validate:"required"`
	RejectReason        string `json:"reject_reason"`
}

func (r *SellerBulkPurchaseOrderRepo) BulkPurchaseOrderRejectPO(params BulkPurchaseOrderRejectPOParams) (*models.BulkPurchaseOrder, error) {
	bulkPO, err := NewBulkPurchaseOrderRepo(r.db).GetBulkPurchaseOrder(GetBulkPurchaseOrderParams{
		JwtClaimsInfo:       params.JwtClaimsInfo,
		BulkPurchaseOrderID: params.BulkPurchaseOrderID,
	})
	if err != nil {
		return nil, err
	}

	if bulkPO.SellerPoAttachments != nil {
		var pendingItems = lo.Filter(*bulkPO.SellerPoAttachments, func(item *models.PoAttachment, index int) bool {
			return item.Status != enums.PoAttachmentStatusApproved
		})

		var items models.PoAttachments = lo.Map(pendingItems, func(item *models.PoAttachment, index int) *models.PoAttachment {
			item.Status = enums.PoAttachmentStatusRejected
			item.RejectReason = params.RejectReason
			return item
		})

		err = r.db.Transaction(func(tx *gorm.DB) error {
			var updates = models.BulkPurchaseOrder{
				SellerPoAttachments:  &items,
				SellerTrackingStatus: enums.SellerBulkPoTrackingStatusPORejected,
			}
			err = tx.Model(&models.BulkPurchaseOrder{}).Where("id = ?", params.BulkPurchaseOrderID).Updates(&updates).Error
			if err != nil {
				return err
			}

			var log = models.BulkPurchaseOrderTracking{
				PurchaseOrderID: bulkPO.ID,
				ActionType:      enums.BulkPoTrackingActionSellerRejectPO,
				UserID:          bulkPO.SellerID,
				UserGroup:       enums.PoTrackingUserGroupSeller,
				CreatedByUserID: params.JwtClaimsInfo.GetUserID(),
				Metadata: &models.PoTrackingMetadata{
					Before: map[string]interface{}{
						"seller_po_attachments":  bulkPO.SellerPoAttachments,
						"seller_tracking_status": bulkPO.SellerTrackingStatus,
					},
					After: map[string]interface{}{
						"seller_tracking_status": updates.SellerTrackingStatus,
					},
				},
			}

			return tx.Create(&log).Error
		})
		if err != nil {
			return nil, err
		}
	}

	return bulkPO, err
}

type BulkPurchaseOrderStartWithoutFirstPaymentParams struct {
	models.JwtClaimsInfo

	BulkPurchaseOrderID string `json:"bulk_purchase_order_id" param:"bulk_purchase_order_id" validate:"required"`
}

func (r *SellerBulkPurchaseOrderRepo) BulkPurchaseOrderStartWithoutFirstPayment(params BulkPurchaseOrderStartWithoutFirstPaymentParams) (*models.BulkPurchaseOrder, error) {
	bulkPO, err := NewBulkPurchaseOrderRepo(r.db).GetBulkPurchaseOrder(GetBulkPurchaseOrderParams{
		JwtClaimsInfo:       params.JwtClaimsInfo,
		BulkPurchaseOrderID: params.BulkPurchaseOrderID,
	})
	if err != nil {
		return nil, err
	}

	if bulkPO.SellerTrackingStatus != enums.SellerBulkPoTrackingStatusFirstPaymentSkipped {
		return nil, errs.ErrBulkPoNotAbleToStart
	}

	var updates = models.BulkPurchaseOrder{
		SellerTrackingStatus: enums.SellerBulkPoTrackingStatusFirstPaymentConfirmed,
	}

	err = r.db.Transaction(func(tx *gorm.DB) error {
		err = tx.Model(&models.BulkPurchaseOrder{}).Where("id = ?", params.BulkPurchaseOrderID).Updates(&updates).Error
		if err != nil {
			return err
		}

		var log = models.BulkPurchaseOrderTracking{
			PurchaseOrderID: bulkPO.ID,
			ActionType:      enums.BulkPoTrackingActionAdminRejectSellerQuotation,
			UserGroup:       enums.PoTrackingUserGroupSeller,
			UserID:          bulkPO.SellerID,
			CreatedByUserID: params.JwtClaimsInfo.GetUserID(),
			Metadata: &models.PoTrackingMetadata{
				Before: map[string]interface{}{
					"seller_tracking_status": bulkPO.SellerTrackingStatus,
				},
				After: map[string]interface{}{
					"seller_tracking_status": updates.SellerTrackingStatus,
				},
			},
		}

		return tx.Create(&log).Error
	})
	if err != nil {
		return nil, err
	}

	bulkPO.SellerTrackingStatus = updates.SellerTrackingStatus
	return bulkPO, err
}

type BulkPurchaseOrderConfirmReceiveFirstPaymentParams struct {
	models.JwtClaimsInfo

	BulkPurchaseOrderID string `json:"bulk_purchase_order_id" param:"bulk_purchase_order_id" validate:"required"`
}

func (r *SellerBulkPurchaseOrderRepo) BulkPurchaseOrderConfirmReceiveFirstPayment(params BulkPurchaseOrderConfirmReceiveFirstPaymentParams) (*models.BulkPurchaseOrder, error) {
	bulkPO, err := NewBulkPurchaseOrderRepo(r.db).GetBulkPurchaseOrder(GetBulkPurchaseOrderParams{
		JwtClaimsInfo:       params.JwtClaimsInfo,
		BulkPurchaseOrderID: params.BulkPurchaseOrderID,
	})
	if err != nil {
		return nil, err
	}

	if bulkPO.SellerTrackingStatus != enums.SellerBulkPoTrackingStatusFirstPaymentConfirm {
		return nil, errs.ErrBulkPoNotAbleToConfirm
	}

	var updates = models.BulkPurchaseOrder{
		SellerTrackingStatus: enums.SellerBulkPoTrackingStatusFirstPaymentConfirmed,
	}

	err = r.db.Transaction(func(tx *gorm.DB) error {
		err = tx.Model(&models.BulkPurchaseOrder{}).Where("id = ?", params.BulkPurchaseOrderID).Updates(&updates).Error
		if err != nil {
			return err
		}

		var log = models.BulkPurchaseOrderTracking{
			PurchaseOrderID: bulkPO.ID,
			ActionType:      enums.BulkPoTrackingActionFirstPaymentConfirmed,
			UserGroup:       enums.PoTrackingUserGroupSeller,
			UserID:          bulkPO.SellerID,
			CreatedByUserID: params.JwtClaimsInfo.GetUserID(),
			Metadata: &models.PoTrackingMetadata{
				Before: map[string]interface{}{
					"seller_tracking_status": bulkPO.SellerTrackingStatus,
				},
				After: map[string]interface{}{
					"seller_tracking_status": updates.SellerTrackingStatus,
				},
			},
		}

		return tx.Create(&log).Error
	})
	if err != nil {
		return nil, err
	}

	bulkPO.SellerTrackingStatus = updates.SellerTrackingStatus
	return bulkPO, err
}

type BulkPurchaseOrderConfirmReceiveFinalPaymentParams struct {
	models.JwtClaimsInfo

	BulkPurchaseOrderID string `json:"bulk_purchase_order_id" param:"bulk_purchase_order_id" validate:"required"`
}

func (r *SellerBulkPurchaseOrderRepo) BulkPurchaseOrderConfirmReceiveFinalPayment(params BulkPurchaseOrderConfirmReceiveFinalPaymentParams) (*models.BulkPurchaseOrder, error) {
	bulkPO, err := NewBulkPurchaseOrderRepo(r.db).GetBulkPurchaseOrder(GetBulkPurchaseOrderParams{
		JwtClaimsInfo:       params.JwtClaimsInfo,
		BulkPurchaseOrderID: params.BulkPurchaseOrderID,
	})
	if err != nil {
		return nil, err
	}

	if bulkPO.SellerTrackingStatus != enums.SellerBulkPoTrackingStatusFinalPaymentConfirm {
		return nil, errs.ErrBulkPoNotAbleToConfirm
	}

	var updates = models.BulkPurchaseOrder{
		SellerTrackingStatus: enums.SellerBulkPoTrackingStatusFinalPaymentConfirmed,
	}

	err = r.db.Transaction(func(tx *gorm.DB) error {
		err = tx.Model(&models.BulkPurchaseOrder{}).Where("id = ?", params.BulkPurchaseOrderID).Updates(&updates).Error
		if err != nil {
			return err
		}

		var log = models.BulkPurchaseOrderTracking{
			PurchaseOrderID: bulkPO.ID,
			ActionType:      enums.BulkPoTrackingActionFinalPaymentConfirmed,
			UserGroup:       enums.PoTrackingUserGroupSeller,
			UserID:          bulkPO.SellerID,
			CreatedByUserID: params.JwtClaimsInfo.GetUserID(),
			Metadata: &models.PoTrackingMetadata{
				Before: map[string]interface{}{
					"seller_tracking_status": bulkPO.SellerTrackingStatus,
				},
				After: map[string]interface{}{
					"seller_tracking_status": updates.SellerTrackingStatus,
				},
			},
		}

		return tx.Create(&log).Error
	})
	if err != nil {
		return nil, err
	}

	bulkPO.SellerTrackingStatus = updates.SellerTrackingStatus
	return bulkPO, err
}

type AdminSellerBulkPoConfirmDeliveredParams struct {
	models.JwtClaimsInfo

	BulkPurchaseOrderID string `json:"bulk_purchase_order_id" param:"bulk_purchase_order_id" query:"bulk_purchase_order_id" validate:"required"`
}

func (r *SellerBulkPurchaseOrderRepo) AdminSellerBulkPoConfirmDelivered(params AdminSellerBulkPoConfirmDeliveredParams) (*models.BulkPurchaseOrder, error) {
	order, err := NewBulkPurchaseOrderRepo(r.db).GetBulkPurchaseOrder(GetBulkPurchaseOrderParams{
		BulkPurchaseOrderID: params.BulkPurchaseOrderID,
		JwtClaimsInfo:       params.JwtClaimsInfo,
	})
	if err != nil {
		return nil, err
	}
	validStatus := []enums.SellerBulkPoTrackingStatus{enums.SellerBulkPoTrackingStatusDelivering}
	if ok := slices.Contains(validStatus, order.SellerTrackingStatus); !ok {
		return nil, errs.ErrPoInvalidToConfirmDelivered
	}

	var updates = models.BulkPurchaseOrder{
		SellerTrackingStatus: enums.SellerBulkPoTrackingStatusDeliveryConfirmed,
	}

	err = r.db.Transaction(func(tx *gorm.DB) error {
		if order.SellerTrackingStatus != updates.SellerTrackingStatus {
			err = NewBulkPurchaseOrderTrackingRepo(r.db).CreateBulkPurchaseOrderTrackingTx(tx, models.BulkPurchaseOrderTrackingCreateForm{
				PurchaseOrderID: params.BulkPurchaseOrderID,
				ActionType:      enums.BulkPoTrackingActionConfirmDelivered,
				UserGroup:       enums.PoTrackingUserGroupSeller,
				UserID:          params.GetUserID(),
				CreatedByUserID: params.JwtClaimsInfo.GetUserID(),
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

		return tx.Model(&models.BulkPurchaseOrder{}).
			Where("id = ?", params.BulkPurchaseOrderID).
			Updates(&updates).Error
	})

	order.SellerTrackingStatus = updates.SellerTrackingStatus

	return order, err
}

type SellerBulkPurchaseOrderMarkDeliveringParams struct {
	models.JwtClaimsInfo

	BulkPurchaseOrderID string                 `param:"bulk_purchase_order_id" validate:"required"`
	LogisticInfo        *models.PoLogisticMeta `json:"logistic_info" param:"logistic_info" query:"logistic_info" form:"logistic_info" validate:"required"`
}

func (r *SellerBulkPurchaseOrderRepo) SellerBulkPurchaseOrderMarkDelivering(params SellerBulkPurchaseOrderMarkDeliveringParams) (*models.BulkPurchaseOrder, error) {
	order, err := NewBulkPurchaseOrderRepo(r.db).GetBulkPurchaseOrder(GetBulkPurchaseOrderParams{
		BulkPurchaseOrderID: params.BulkPurchaseOrderID,
		JwtClaimsInfo:       params.JwtClaimsInfo,
	})

	if err != nil {
		return nil, err
	}

	var updates = models.BulkPurchaseOrder{
		SellerDeliveryStartedAt: values.Int64(time.Now().Unix()),
		SellerTrackingStatus:    enums.SellerBulkPoTrackingStatusDelivering,
		SellerLogisticInfo:      params.LogisticInfo,
	}

	err = r.db.Transaction(func(tx *gorm.DB) error {
		if order.TrackingStatus != updates.TrackingStatus {
			err = NewBulkPurchaseOrderTrackingRepo(r.db).CreateBulkPurchaseOrderTrackingTx(tx, models.BulkPurchaseOrderTrackingCreateForm{
				PurchaseOrderID: params.BulkPurchaseOrderID,
				ActionType:      enums.BulkPoTrackingActionSellerDelivering,
				UserID:          order.UserID,
				CreatedByUserID: params.JwtClaimsInfo.GetUserID(),
				UserGroup:       enums.PoTrackingUserGroupSeller,
				Metadata: &models.PoTrackingMetadata{
					Before: map[string]interface{}{
						"seller_tracking_status": order.SellerTrackingStatus,
					},
					After: map[string]interface{}{
						"seller_tracking_status": updates.SellerTrackingStatus,
						"logistic_info":          params.LogisticInfo,
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
	order.DeliveryStartedAt = updates.DeliveryStartedAt
	order.TrackingStatus = updates.TrackingStatus
	order.SellerLogisticInfo = updates.SellerLogisticInfo

	return order, err
}

type SellerBulkPurchaseOrderFeedbackParams struct {
	models.JwtClaimsInfo

	BulkPurchaseOrderID string `json:"bulk_purchase_order_id" param:"bulk_purchase_order_id" validate:"required"`

	Feedback *models.BulkPurchaseOrderFeedback `json:"feedback" validate:"required"`
}

func (r *SellerBulkPurchaseOrderRepo) SellerBulkPurchaseOrderFeedback(params SellerBulkPurchaseOrderFeedbackParams) (*models.BulkPurchaseOrder, error) {
	bulkPO, err := NewBulkPurchaseOrderRepo(r.db).GetBulkPurchaseOrder(GetBulkPurchaseOrderParams{
		BulkPurchaseOrderID: params.BulkPurchaseOrderID,
	})
	if err != nil {
		return nil, err
	}
	var updates = models.BulkPurchaseOrder{
		SellerFeedback: params.Feedback,
	}

	err = r.db.Model(&models.BulkPurchaseOrder{}).Where("id = ?", params.BulkPurchaseOrderID).Updates(&updates).Error
	if err != nil {
		return nil, err
	}

	bulkPO.SellerFeedback = updates.SellerFeedback

	return bulkPO, err
}

type SellerBulkPurchaseOrderFirstPaymentInvoiceParams struct {
	models.JwtClaimsInfo

	BulkPurchaseOrderID string `json:"bulk_purchase_order_id" param:"bulk_purchase_order_id" validate:"required"`
}

func (r *SellerBulkPurchaseOrderRepo) SellerBulkPurchaseOrderFirstPaymentInvoice(params SellerBulkPurchaseOrderFirstPaymentInvoiceParams) (*models.Attachment, error) {
	return nil, eris.New("Not implemented yet")
}

type SellerBulkPurchaseOrderFinalPaymentInvoiceParams struct {
	models.JwtClaimsInfo

	BulkPurchaseOrderID string `json:"bulk_purchase_order_id" param:"bulk_purchase_order_id" validate:"required"`
}

func (r *SellerBulkPurchaseOrderRepo) SellerBulkPurchaseOrderFinalPaymentInvoice(params SellerBulkPurchaseOrderFinalPaymentInvoiceParams) (*models.Attachment, error) {
	return nil, eris.New("Not implemented yet")
}

type SellerBulkPurchaseOrderUpdateTrackingStatusParams struct {
	models.JwtClaimsInfo

	BulkPurchaseOrderID  string                           `json:"bulk_purchase_order_id" param:"bulk_purchase_order_id" query:"bulk_purchase_order_id" validate:"required"`
	SellerTrackingStatus enums.SellerBulkPoTrackingStatus `json:"sellertracking_status" param:"sellertracking_status" query:"sellertracking_status"`
	TrackingAction       enums.BulkPoTrackingAction       `json:"tracking_action" param:"tracking_action" query:"tracking_action"`
}

func (r *SellerBulkPurchaseOrderRepo) SellerBulkPurchaseOrderUpdateTrackingStatus(params SellerBulkPurchaseOrderUpdateTrackingStatusParams) (*models.BulkPurchaseOrder, error) {
	order, err := NewBulkPurchaseOrderRepo(r.db).GetBulkPurchaseOrder(GetBulkPurchaseOrderParams{
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

	err = r.db.Transaction(func(tx *gorm.DB) error {
		if order.SellerTrackingStatus != updates.SellerTrackingStatus {
			err = NewBulkPurchaseOrderTrackingRepo(r.db).CreateBulkPurchaseOrderTrackingTx(tx, models.BulkPurchaseOrderTrackingCreateForm{
				PurchaseOrderID: params.BulkPurchaseOrderID,
				ActionType:      params.TrackingAction,
				UserID:          order.SellerID,
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
			})

			if err != nil {
				return eris.Wrap(err, err.Error())
			}

		}
		return r.db.Model(&models.BulkPurchaseOrder{}).
			Where("id = ?", params.BulkPurchaseOrderID).
			Updates(&updates).Error

	})
	if err != nil {
		return nil, err
	}

	order.SellerTrackingStatus = updates.SellerTrackingStatus
	return order, err
}

type SellerBulkPurchaseOrderMarkProductionParams struct {
	models.JwtClaimsInfo

	BulkPurchaseOrderID string                   `json:"bulk_purchase_order_id" param:"bulk_purchase_order_id" query:"bulk_purchase_order_id" validate:"required"`
	ProductionInfo      *models.PoProductionMeta `json:"production_info" param:"production_info" query:"production_info"`
}

func (r *SellerBulkPurchaseOrderRepo) SellerBulkPurchaseOrderMarkProduction(params SellerBulkPurchaseOrderMarkProductionParams) (*models.BulkPurchaseOrder, error) {
	order, err := NewBulkPurchaseOrderRepo(r.db).GetBulkPurchaseOrder(GetBulkPurchaseOrderParams{
		JwtClaimsInfo:       params.JwtClaimsInfo,
		BulkPurchaseOrderID: params.BulkPurchaseOrderID,
		IncludeUser:         true,
	})
	if err != nil {
		return nil, err
	}

	var updates = models.BulkPurchaseOrder{
		SellerProductionInfo: params.ProductionInfo,
		SellerTrackingStatus: enums.SellerBulkPoTrackingStatusProduction,
	}

	err = r.db.Transaction(func(tx *gorm.DB) error {
		if order.SellerTrackingStatus != updates.SellerTrackingStatus {
			var record = models.BulkPurchaseOrderTracking{
				PurchaseOrderID: params.BulkPurchaseOrderID,
				ActionType:      enums.BulkPoTrackingActionSellerMarkProduction,
				UserGroup:       enums.PoTrackingUserGroupSeller,
				UserID:          order.SellerID,
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
			err = tx.Create(&record).Error
			if err != nil {
				return err
			}
		}

		err = tx.Model(&models.BulkPurchaseOrder{}).
			Where("id = ?", params.BulkPurchaseOrderID).
			Updates(&updates).Error

		return err
	})

	if err != nil {
		return nil, eris.Wrap(err, err.Error())
	}

	order.SellerTrackingStatus = updates.SellerTrackingStatus

	return order, err
}

type SellerBulkPurchaseOrderMarkInspectionParams struct {
	models.JwtClaimsInfo

	BulkPurchaseOrderID            string             `json:"bulk_purchase_order_id" param:"bulk_purchase_order_id" query:"bulk_purchase_order_id" validate:"required"`
	InspectionProcedureAttachments models.Attachments `json:"inspection_procedure_attachments"`
	InspectionProcedureNote        string             `json:"inspection_procedure_note"`
}

func (r *SellerBulkPurchaseOrderRepo) SellerBulkPurchaseOrderMarkInspection(params SellerBulkPurchaseOrderMarkInspectionParams) (*models.BulkPurchaseOrder, error) {
	order, err := NewBulkPurchaseOrderRepo(r.db).GetBulkPurchaseOrder(GetBulkPurchaseOrderParams{
		JwtClaimsInfo:       params.JwtClaimsInfo,
		BulkPurchaseOrderID: params.BulkPurchaseOrderID,
		IncludeUser:         true,
	})
	if err != nil {
		return nil, err
	}

	var updates = models.BulkPurchaseOrder{
		SellerInspectionProcedureAttachments: &params.InspectionProcedureAttachments,
		SellerInspectionProcedureNote:        params.InspectionProcedureNote,
		SellerTrackingStatus:                 enums.SellerBulkPoTrackingStatusInspection,
	}

	err = r.db.Transaction(func(tx *gorm.DB) error {
		if order.SellerTrackingStatus != updates.SellerTrackingStatus {
			var record = models.BulkPurchaseOrderTracking{
				PurchaseOrderID: params.BulkPurchaseOrderID,
				ActionType:      enums.BulkPoTrackingActionSellerMarkInspection,
				UserID:          order.SellerID,
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
			err = tx.Create(&record).Error
			if err != nil {
				return err
			}
		}

		err = tx.Model(&models.BulkPurchaseOrder{}).
			Where("id = ?", params.BulkPurchaseOrderID).
			Updates(&updates).Error

		return err
	})

	if err != nil {
		return nil, eris.Wrap(err, err.Error())
	}

	order.SellerTrackingStatus = updates.SellerTrackingStatus

	return order, err
}

type AdminSellerBulkPurchaseOrderUpdatePpsParams struct {
	models.JwtClaimsInfo

	BulkPurchaseOrderID string            `json:"bulk_purchase_order_id" param:"bulk_purchase_order_id" query:"bulk_purchase_order_id" validate:"required"`
	PpsInfo             *models.PoPpsMeta `json:"pps_info" param:"pps_info" query:"pps_info"`
}

func (r *SellerBulkPurchaseOrderRepo) AdminSellerBulkPurchaseOrderUpdatePps(params AdminSellerBulkPurchaseOrderUpdatePpsParams) (*models.BulkPurchaseOrder, error) {
	order, err := NewBulkPurchaseOrderRepo(r.db).GetBulkPurchaseOrder(GetBulkPurchaseOrderParams{
		BulkPurchaseOrderID: params.BulkPurchaseOrderID,
		JwtClaimsInfo:       params.JwtClaimsInfo,
		IncludeUser:         true,
	})
	if err != nil {
		return nil, err
	}
	var ppsInfoArr models.PoPpsMetas
	if order.SellerPpsInfo != nil {
		ppsInfoArr = append(ppsInfoArr, *order.SellerPpsInfo...)
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
		SellerPpsInfo:        &ppsInfoArr,
		SellerTrackingStatus: enums.SellerBulkPoTrackingStatusPps,
	}

	err = r.db.Transaction(func(tx *gorm.DB) error {
		err = NewBulkPurchaseOrderTrackingRepo(r.db).CreateBulkPurchaseOrderTrackingTx(tx, models.BulkPurchaseOrderTrackingCreateForm{
			PurchaseOrderID: order.ID,
			ActionType:      enums.BulkPoTrackingActionUpdatePps,
			UserID:          order.UserID,
			UserGroup:       enums.PoTrackingUserGroupSeller,
			CreatedByUserID: params.JwtClaimsInfo.GetUserID(),
			Metadata: &models.PoTrackingMetadata{
				Before: map[string]interface{}{
					"seller_pps_info": order.PpsInfo,
				},
				After: map[string]interface{}{
					"seller_pps_info": updates.PpsInfo,
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
