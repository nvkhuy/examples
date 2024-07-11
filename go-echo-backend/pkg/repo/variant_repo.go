package repo

import (
	"database/sql"
	"strings"

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

type VariantRepo struct {
	db     *db.DB
	logger *logger.Logger
}

func NewVariantRepo(db *db.DB) *VariantRepo {
	return &VariantRepo{
		db:     db,
		logger: logger.New("repo/variant"),
	}
}

type PaginateVariantParams struct {
	models.PaginationParams

	Name    string `json:"name" query:"name" form:"name"`
	ForRole enums.Role
}

func (r *VariantRepo) PaginateVariants(params PaginateVariantParams) *query.Pagination {
	var builder = queryfunc.NewVariantBuilder(queryfunc.VariantBuilderOptions{
		QueryBuilderOptions: queryfunc.QueryBuilderOptions{
			Role: enums.RoleClient,
		},
	})

	if params.Limit == 0 {
		params.Limit = 20
	}

	var result = query.New(r.db, builder).
		WhereFunc(func(builder *query.Builder) {

			if strings.TrimSpace(params.Keyword) != "" {
				var q = "%" + params.Keyword + "%"
				builder.Where("name ILIKE @query OR name ILIKE @query", sql.Named("query", q))
			}

		}).
		Page(params.Page).
		Limit(params.Limit).
		PagingFunc()

	return result
}

func (r *VariantRepo) GetVariantByID(VariantID string, options queryfunc.VariantBuilderOptions) (*models.Variant, error) {
	var builder = queryfunc.NewVariantBuilder(options)
	var Variant models.Variant
	var err = query.New(r.db, builder).
		WhereFunc(func(builder *query.Builder) {
			builder.Where("v.id = ?", VariantID)
		}).
		FirstFunc(&Variant)

	if err != nil {
		if r.db.IsRecordNotFoundError(err) {
			return nil, errs.ErrRecordNotFound
		}
		return nil, err
	}

	return &Variant, nil
}

func (r *VariantRepo) UpdateVariantByID(VariantID string, form models.Variant) (*models.Variant, error) {
	var update models.Variant

	var err = copier.Copy(&update, &form)
	if err != nil {
		return nil, eris.Wrap(err, err.Error())
	}

	err = r.db.Omit(clause.Associations).Model(&models.Variant{}).Where("id = ?", VariantID).Updates(&update).Error
	if err != nil {
		return nil, err
	}

	return r.GetVariantByID(VariantID, queryfunc.VariantBuilderOptions{
		QueryBuilderOptions: queryfunc.QueryBuilderOptions{
			Role: enums.RoleClient,
		},
	})
}

func (r *VariantRepo) DeleteVariantByID(VariantID string) error {
	var err = r.db.Delete(&models.Variant{}, "id = ?", VariantID).Error

	return err
}

func (r *VariantRepo) GetVariantsByProductID(ProductID string, options queryfunc.VariantBuilderOptions) ([]*models.Variant, error) {
	var builder = queryfunc.NewVariantBuilder(options)
	var variants []*models.Variant
	var err = query.New(r.db, builder).
		WhereFunc(func(builder *query.Builder) {
			builder.Where("v.product_id = ?", ProductID)
		}).
		FindFunc(&variants)

	if err != nil {
		if r.db.IsRecordNotFoundError(err) {
			return nil, errs.ErrRecordNotFound
		}
		return nil, err
	}

	return variants, nil
}
