package repo

import (
	"github.com/engineeringinflow/inflow-backend/pkg/db"
	"github.com/engineeringinflow/inflow-backend/pkg/errs"
	"github.com/engineeringinflow/inflow-backend/pkg/logger"
	"github.com/engineeringinflow/inflow-backend/pkg/models"
	"github.com/engineeringinflow/inflow-backend/pkg/models/enums"
	"github.com/engineeringinflow/inflow-backend/pkg/repo/query"
	"github.com/engineeringinflow/inflow-backend/pkg/repo/query/queryfunc"
	"github.com/jinzhu/copier"

	"github.com/rotisserie/eris"

	"gorm.io/gorm/clause"
)

type QuantityPriceTierRepo struct {
	db     *db.DB
	logger *logger.Logger
}

func NewQuantityPriceTierRepo(db *db.DB) *QuantityPriceTierRepo {
	return &QuantityPriceTierRepo{
		db:     db,
		logger: logger.New("repo/QuantityPriceTier"),
	}
}

type PaginateQuantityPriceTiersParams struct {
	models.PaginationParams

	ForRole enums.Role
}

// func (r *QuantityPriceTierRepo) PaginateQuantityPriceTiers(params PaginateQuantityPriceTiersParams) *query.Pagination {
// 	var builder = queryfunc.NewQuantityPriceTierBuilder(queryfunc.QuantityPriceTierBuilderOptions{
// 		ForRole: params.ForRole,
// 	})

// 	var result = query.New(r.adb, builder).
// 		WhereFunc(func(builder *query.Builder) {
// 			if params.Email != "" {
// 				builder.Where("sub.email ILIKE ?", "%"+params.Email)
// 			}
// 		}).
// 		Page(params.Page).
// 		Limit(params.Limit).
// 		PagingFunc()

// 	return result
// }

type SearchQuantityPriceTiersParams struct {
	models.PaginationParams

	Roles []string `json:"roles" query:"roles" form:"roles"`

	AccountStatuses []enums.AccountStatus `json:"account_statuses" query:"account_statuses" form:"account_statuses"`

	ForRole enums.Role
}

func (r *QuantityPriceTierRepo) CreateQuantityPriceTier(form models.SubscribeByEmailForm) (*models.QuantityPriceTier, error) {
	var QuantityPriceTier models.QuantityPriceTier
	err := copier.Copy(&QuantityPriceTier, &form)
	if err != nil {
		return nil, eris.Wrap(err, "")
	}

	err = r.db.Omit(clause.Associations).Create(&QuantityPriceTier).Error
	if err != nil {
		if r.db.IsRecordNotFoundError(err) {
			return nil, errs.ErrRecordNotFound
		}
		return nil, eris.Wrap(err, "")
	}

	return &QuantityPriceTier, nil
}

func (r *QuantityPriceTierRepo) GetQuantityPriceTierByID(QuantityPriceTierID string, options queryfunc.QuantityPriceTierBuilderOptions) (*models.QuantityPriceTier, error) {
	var builder = queryfunc.NewQuantityPriceTierBuilder(options)
	var QuantityPriceTier models.QuantityPriceTier
	var err = query.New(r.db, builder).
		WhereFunc(func(builder *query.Builder) {
			builder.Where("sub.id = ?", QuantityPriceTierID)
		}).
		FirstFunc(&QuantityPriceTier)

	if err != nil {
		if r.db.IsRecordNotFoundError(err) {
			return nil, errs.ErrRecordNotFound
		}
		return nil, err
	}

	return &QuantityPriceTier, nil
}

func (r *QuantityPriceTierRepo) GetQuantityPriceTierByProductID(productID string, options queryfunc.QuantityPriceTierBuilderOptions) ([]*models.QuantityPriceTier, error) {
	var builder = queryfunc.NewQuantityPriceTierBuilder(options)
	var tiers []*models.QuantityPriceTier
	var err = query.New(r.db, builder).
		WhereFunc(func(builder *query.Builder) {
			builder.Where("tier.product_id = ?", productID)
		}).
		FindFunc(&tiers)

	if err != nil {
		if r.db.IsRecordNotFoundError(err) {
			return nil, errs.ErrRecordNotFound
		}
		return nil, err
	}

	return tiers, nil
}
