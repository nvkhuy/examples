package repo

import (
	"github.com/engineeringinflow/inflow-backend/pkg/db"
	"github.com/engineeringinflow/inflow-backend/pkg/logger"
	"github.com/engineeringinflow/inflow-backend/pkg/models"
	"github.com/engineeringinflow/inflow-backend/pkg/models/enums"
	"github.com/engineeringinflow/inflow-backend/pkg/repo/query"
	"github.com/engineeringinflow/inflow-backend/pkg/repo/query/queryfunc"
	"github.com/jinzhu/copier"
	"github.com/rotisserie/eris"
	"golang.org/x/exp/slices"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type BulkPurchaseOrderTrackingRepo struct {
	db     *db.DB
	logger *logger.Logger
}

func NewBulkPurchaseOrderTrackingRepo(db *db.DB) *BulkPurchaseOrderTrackingRepo {
	return &BulkPurchaseOrderTrackingRepo{
		db:     db,
		logger: logger.New("repo/BulkPurchaseOrderTracking"),
	}
}

func (r *BulkPurchaseOrderTrackingRepo) CreateBulkPurchaseOrderTrackingTx(tx *gorm.DB, params models.BulkPurchaseOrderTrackingCreateForm) error {
	var form models.BulkPurchaseOrderTracking
	err := copier.Copy(&form, &params)
	if err != nil {
		return eris.Wrap(err, "")
	}

	oneRowLogActions := []enums.BulkPoTrackingAction{enums.BulkPoTrackingActionUpdatePps, enums.BulkPoTrackingActionUpdateProduction}
	if ok := slices.Contains(oneRowLogActions, params.ActionType); ok {
		err = tx.Model(&models.BulkPurchaseOrderTracking{}).Where("purchase_order_id = ? AND action_type = ?", params.PurchaseOrderID, params.ActionType).Assign(form).FirstOrCreate(&form).Error
		return err
	}

	err = tx.Omit(clause.Associations).Create(&form).Error
	if err != nil {
		return eris.Wrap(err, "")
	}
	return err
}

type PaginateBulkPurchaseOrderTrackingParams struct {
	models.PaginationParams
	models.JwtClaimsInfo

	BulkPurchaseOrderID string                    `json:"bulk_purchase_order_id" query:"bulk_purchase_order_id" form:"bulk_purchase_order_id" param:"bulk_purchase_order_id" validate:"required"`
	UserID              string                    `json:"user_id" query:"user_id" form:"user_id" param:"user_id"`
	ActionTypes         []enums.AuditActionType   `json:"action_types" query:"action_types" form:"action_types" param:"action_types"`
	UserGroup           enums.PoTrackingUserGroup `json:"user_group" query:"user_group" form:"user_group" param:"user_group"`
	IncludeTrackings    bool                      `json:"-"`
}

func (r *BulkPurchaseOrderTrackingRepo) PaginateBulkPurchaseOrderTrackings(params PaginateBulkPurchaseOrderTrackingParams) *query.Pagination {
	var result = query.New(r.db, queryfunc.NewBulkPurchaseOrderTrackingBuilder(queryfunc.BulkPurchaseOrderTrackingBuilderOptions{
		QueryBuilderOptions: queryfunc.QueryBuilderOptions{
			Role: params.GetRole(),
		},
	})).
		WhereFunc(func(builder *query.Builder) {
			builder.Where("pot.purchase_order_id = ?", params.BulkPurchaseOrderID)

			if len(params.ActionTypes) > 0 {
				builder.Where("pot.action_type IN ?", params.ActionTypes)
			}

			if params.UserGroup != "" {
				builder.Where("pot.user_group = ?", params.UserGroup)
			}

			if params.GetRole().IsSeller() {
				builder.Where("(pot.user_id = @user_id OR bposq.user_id = @user_id)")

				builder.Where(map[string]interface{}{
					"user_id": params.GetUserID(),
				})
			}
		}).
		Limit(params.Limit).
		Page(params.Page).
		PagingFunc()
	return result
}
