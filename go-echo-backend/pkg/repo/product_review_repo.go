package repo

import (
	"github.com/engineeringinflow/inflow-backend/pkg/db"
	"github.com/engineeringinflow/inflow-backend/pkg/errs"
	"github.com/engineeringinflow/inflow-backend/pkg/logger"
	"github.com/engineeringinflow/inflow-backend/pkg/models"
	"github.com/engineeringinflow/inflow-backend/pkg/repo/query"
	"github.com/engineeringinflow/inflow-backend/pkg/repo/query/queryfunc"
	"github.com/jinzhu/copier"

	"github.com/rotisserie/eris"

	"gorm.io/gorm/clause"
)

type ProductReviewRepo struct {
	db     *db.DB
	logger *logger.Logger
}

func NewProductReviewRepo(db *db.DB) *ProductReviewRepo {
	return &ProductReviewRepo{
		db:     db,
		logger: logger.New("repo/product_review"),
	}
}

type PaginateProductReviewParams struct {
	models.PaginationParams
	models.JwtClaimsInfo

	ProductID string `json:"product_id" query:"product_id" form:"product_id"`
	ShopID    string `json:"shop_id" query:"shop_id" form:"shop_id"`
	Name      string `json:"name" query:"name" form:"name"`
}

func (r *ProductReviewRepo) PaginateProductReviews(params PaginateProductReviewParams) *query.Pagination {
	var builder = queryfunc.NewProductReviewBuilder(queryfunc.ProductReviewBuilderOptions{
		QueryBuilderOptions: queryfunc.QueryBuilderOptions{
			Role: params.GetRole(),
		},
	})

	if params.Limit == 0 {
		params.Limit = 20
	}

	var result = query.New(r.db, builder).
		WhereFunc(func(builder *query.Builder) {
			if params.ProductID != "" {
				builder.Where("pr.product_id = ?", params.ProductID)
			}
			if params.ShopID != "" {
				builder.Where("pr.shop_id = ?", params.ShopID)
			}
		}).
		Page(params.Page).
		Limit(params.Limit).
		PagingFunc()

	return result
}

func (r *ProductReviewRepo) CreateReview(form models.ProductReviewCreateForm) (*models.ProductReview, error) {
	var productReview models.ProductReview
	err := copier.Copy(&productReview, &form)
	if err != nil {
		return nil, eris.Wrap(err, "")
	}

	err = r.db.Omit(clause.Associations).Create(&productReview).Error
	if err != nil {
		if r.db.IsRecordNotFoundError(err) {
			return nil, errs.ErrRecordNotFound
		}
		return nil, eris.Wrap(err, "")
	}

	return &productReview, nil
}

func (r *ProductReviewRepo) GetReviewByID(reviewID string, options queryfunc.ProductReviewBuilderOptions) (*models.ProductReview, error) {
	var builder = queryfunc.NewProductReviewBuilder(options)
	var review models.ProductReview
	var err = query.New(r.db, builder).
		WhereFunc(func(builder *query.Builder) {
			builder.Where("pr.id = ?", reviewID)
		}).
		FirstFunc(&review)

	if err != nil {
		if r.db.IsRecordNotFoundError(err) {
			return nil, errs.ErrRecordNotFound
		}
		return nil, err
	}

	return &review, nil
}
