package repo

import (
	"database/sql"
	"strings"

	"github.com/engineeringinflow/inflow-backend/pkg/db"
	"github.com/engineeringinflow/inflow-backend/pkg/errs"
	"github.com/engineeringinflow/inflow-backend/pkg/logger"
	"github.com/engineeringinflow/inflow-backend/pkg/models"
	"github.com/engineeringinflow/inflow-backend/pkg/repo/query"
	"github.com/engineeringinflow/inflow-backend/pkg/repo/query/queryfunc"
)

type ProductAttributeRepo struct {
	db     *db.DB
	logger *logger.Logger
}

func NewProductAttributeRepo(db *db.DB) *ProductAttributeRepo {
	return &ProductAttributeRepo{
		db:     db,
		logger: logger.New("repo/product_attribute"),
	}
}

type PaginateProductAttributeParams struct {
	models.PaginationParams
	models.JwtClaimsInfo

	Name string `json:"name" query:"name" form:"name"`
}

func (r *ProductAttributeRepo) PaginateProductAttributes(params PaginateProductAttributeParams) *query.Pagination {
	var builder = queryfunc.NewProductAttributeBuilder(queryfunc.ProductAttributeBuilderOptions{
		QueryBuilderOptions: queryfunc.QueryBuilderOptions{
			Role: params.GetRole(),
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

func (r *ProductAttributeRepo) GetProductAttributesByProductID(ProductID string, options queryfunc.ProductAttributeBuilderOptions) ([]*models.ProductAttribute, error) {
	var builder = queryfunc.NewProductAttributeBuilder(options)
	var ProductAttributes []*models.ProductAttribute
	var err = query.New(r.db, builder).
		WhereFunc(func(builder *query.Builder) {
			builder.Where("att.product_id = ?", ProductID)
		}).
		FindFunc(&ProductAttributes)

	if err != nil {
		if r.db.IsRecordNotFoundError(err) {
			return nil, errs.ErrRecordNotFound
		}
		return nil, err
	}

	return ProductAttributes, nil
}
