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

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type BlogCategoryRepo struct {
	db     *db.DB
	logger *logger.Logger
}

func NewBlogCategoryRepo(db *db.DB) *BlogCategoryRepo {
	return &BlogCategoryRepo{
		db:     db,
		logger: logger.New("repo/BlogCategory"),
	}
}

type PaginateBlogCategoryParams struct {
	models.PaginationParams
	models.JwtClaimsInfo

	Name      string `json:"name" query:"name" form:"name"`
	TotalPost *int   `json:"total_post"`
}

func (r *BlogCategoryRepo) PaginateBlogCategory(params PaginateBlogCategoryParams) []*models.BlogCategory {
	var builder = queryfunc.NewBlogCategoryBuilder(queryfunc.BlogCategoryBuilderOptions{
		QueryBuilderOptions: queryfunc.QueryBuilderOptions{
			Role: params.GetRole(),
		},
	})
	if params.Limit == 0 {
		params.Limit = 1000
	}
	var result []*models.BlogCategory
	query.New(r.db, builder).
		WhereFunc(func(builder *query.Builder) {
			if params.TotalPost != nil {
				builder.Where("(select count(1) from posts where category_id = cate.id) > ?", params.TotalPost)
			}
		}).
		Limit(params.Limit).
		FindFunc(&result)

	return result
}

func (r *BlogCategoryRepo) CreateBlogCategory(form models.BlogCategoryCreateForm) (*models.BlogCategory, error) {
	var bc models.BlogCategory
	err := copier.Copy(&bc, &form)
	if err != nil {
		return nil, eris.Wrap(err, "")
	}

	err = r.db.Omit(clause.Associations).Create(&bc).Error
	if err != nil {
		if r.db.IsRecordNotFoundError(err) {
			return nil, errs.ErrCategoryExisted
		}
		return nil, eris.Wrap(err, "")
	}

	return &bc, nil
}

type GetBlogCategoryParams struct {
	BlogCategoryID string `param:"blog_category_id" query:"blog_category_id" form:"blog_category_id" validate:"required"`

	models.JwtClaimsInfo
}

func (r *BlogCategoryRepo) GetBlogCategory(params GetBlogCategoryParams) (*models.BlogCategory, error) {
	var builder = queryfunc.NewBlogCategoryBuilder(queryfunc.BlogCategoryBuilderOptions{
		QueryBuilderOptions: queryfunc.QueryBuilderOptions{
			Role: params.GetRole(),
		},
	})
	var BlogCategory models.BlogCategory
	var err = query.New(r.db, builder).
		WhereFunc(func(builder *query.Builder) {
			builder.Where("cate.id = ?", params.BlogCategoryID)
		}).
		FirstFunc(&BlogCategory)

	if err != nil {
		if r.db.IsRecordNotFoundError(err) {
			return nil, errs.ErrCategoryNotFound
		}
		return nil, err
	}

	return &BlogCategory, nil
}

func (r *BlogCategoryRepo) UpdateBlogCategory(form models.BlogCategoryUpdateForm) (*models.BlogCategory, error) {
	var update models.BlogCategory

	var err = copier.Copy(&update, &form)
	if err != nil {
		return nil, eris.Wrap(err, err.Error())
	}

	err = r.db.Omit(clause.Associations).Select("*").Model(&update).Where("id = ?", form.BlogCategoryID).Updates(&update).Error
	if err != nil {
		return nil, err
	}

	return r.GetBlogCategory(GetBlogCategoryParams{
		BlogCategoryID: form.BlogCategoryID,
	})
}

type DeleteBlogCategoryParams struct {
	BlogCategoryID string `param:"blog_category_id" validate:"required"`

	models.JwtClaimsInfo
}

func (r *BlogCategoryRepo) DeleteBlogCategory(params DeleteBlogCategoryParams) error {
	var cate models.BlogCategory
	var err = r.db.First(&cate, "id = ?", params.BlogCategoryID).Error
	if cate.ID == "" {
		return errs.ErrCategoryNotFound
	}

	err = r.db.Transaction(func(tx *gorm.DB) error {
		err = tx.Model(&models.Post{}).Where("category_id = ?", params.BlogCategoryID).UpdateColumn("category_id", "").Error
		if err != nil {
			return eris.Wrap(err, "")
		}

		err = tx.Unscoped().Delete(&models.BlogCategory{}, "id = ?", params.BlogCategoryID).Error

		return err
	})

	return err
}
