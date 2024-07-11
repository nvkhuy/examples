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
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type PurchaseOrderTrackingRepo struct {
	db     *db.DB
	logger *logger.Logger
}

func NewPurchaseOrderTrackingRepo(db *db.DB) *PurchaseOrderTrackingRepo {
	return &PurchaseOrderTrackingRepo{
		db:     db,
		logger: logger.New("repo/PurchaseOrderTracking"),
	}
}

func (r *PurchaseOrderTrackingRepo) CreatePurchaseOrderTrackingTx(tx *gorm.DB, params models.PurchaseOrderTrackingCreateForm) error {
	var form models.PurchaseOrderTracking
	err := copier.Copy(&form, &params)
	if err != nil {
		return eris.Wrap(err, "")
	}

	err = tx.Omit(clause.Associations).Create(&form).Error
	if err != nil {
		return eris.Wrap(err, "")
	}
	return err
}

type PaginatePurchaseOrderTrackingParams struct {
	models.PaginationParams
	models.JwtClaimsInfo

	PurchaseOrderID string                    `json:"purchase_order_id" query:"purchase_order_id" form:"purchase_order_id" param:"purchase_order_id" validate:"required"`
	UserID          string                    `json:"user_id" query:"user_id" form:"user_id" param:"user_id"`
	ActionType      enums.AuditActionType     `json:"action_type" query:"action_type" form:"action_type" param:"action_type"`
	UserGroup       enums.PoTrackingUserGroup `json:"user_group" param:"user_group"`
}

func (r *PurchaseOrderTrackingRepo) PaginatePurchaseOrderTrackings(params PaginatePurchaseOrderTrackingParams) *query.Pagination {
	var result = query.New(r.db, queryfunc.NewPurchaseOrderTrackingBuilder(queryfunc.PurchaseOrderTrackingBuilderOptions{
		QueryBuilderOptions: queryfunc.QueryBuilderOptions{
			Role: params.GetRole(),
		},
	})).
		WhereFunc(func(builder *query.Builder) {
			builder.Where("pot.purchase_order_id = ?", params.PurchaseOrderID)

			if params.UserID != "" {
				builder.Where("pot.user_id = ?", params.UserID)
			}

			if params.UserGroup != "" {
				builder.Where("pot.user_group = ?", params.UserGroup)
			}
		}).
		Limit(params.Limit).
		Page(params.Page).
		PagingFunc()
	return result
}
