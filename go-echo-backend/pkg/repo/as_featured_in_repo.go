package repo

import (
	"time"

	"github.com/engineeringinflow/inflow-backend/pkg/db"
	"github.com/engineeringinflow/inflow-backend/pkg/errs"
	"github.com/engineeringinflow/inflow-backend/pkg/logger"
	"github.com/engineeringinflow/inflow-backend/pkg/models"
	"github.com/engineeringinflow/inflow-backend/pkg/models/enums"
	"github.com/engineeringinflow/inflow-backend/pkg/repo/query"
	"github.com/engineeringinflow/inflow-backend/pkg/repo/query/queryfunc"
	"github.com/jinzhu/copier"
	"github.com/thaitanloi365/go-utils/values"

	"github.com/rotisserie/eris"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type AsFeaturedInRepo struct {
	db     *db.DB
	logger *logger.Logger
}

func NewAsFeaturedInRepo(db *db.DB) *AsFeaturedInRepo {
	return &AsFeaturedInRepo{
		db:     db,
		logger: logger.New("repo/as_featured_in"),
	}
}

type PaginateAsFeaturedInParams struct {
	models.PaginationParams
	models.JwtClaimsInfo

	Statuses []enums.PostStatus `json:"statuses" query:"statuses" param:"statuses"`
}

func (r *AsFeaturedInRepo) PaginateAsFeaturedIn(params PaginateAsFeaturedInParams) *query.Pagination {
	var builder = queryfunc.NewAsFeaturedInBuilder(queryfunc.AsFeaturedInBuilderOptions{
		QueryBuilderOptions: queryfunc.QueryBuilderOptions{
			Role: params.GetRole(),
		},
	})
	if params.Limit == 0 {
		params.Limit = 10
	}

	var result = query.New(r.db, builder).
		WhereFunc(func(builder *query.Builder) {
			if len(params.Statuses) > 0 {
				builder.Where("afi.status IN ?", params.Statuses)
			}
		}).
		Page(params.Page).
		Limit(params.Limit).
		PagingFunc()

	return result
}

func (r *AsFeaturedInRepo) CreateAsFeaturedIn(form models.AsFeaturedInCreateForm) (*models.AsFeaturedIn, error) {
	var bc models.AsFeaturedIn
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

type GetAsFeaturedInParams struct {
	AsFeaturedInID string `param:"ads_video_id" query:"ads_video_id" form:"ads_video_id" validate:"required"`

	models.JwtClaimsInfo
}

func (r *AsFeaturedInRepo) GetAsFeaturedIn(params GetAsFeaturedInParams) (*models.AsFeaturedIn, error) {
	var builder = queryfunc.NewAsFeaturedInBuilder(queryfunc.AsFeaturedInBuilderOptions{
		QueryBuilderOptions: queryfunc.QueryBuilderOptions{
			Role: params.GetRole(),
		},
	})
	var AsFeaturedIn models.AsFeaturedIn
	var err = query.New(r.db, builder).
		WhereFunc(func(builder *query.Builder) {
			builder.Where("afi.id = ?", params.AsFeaturedInID)
		}).
		FirstFunc(&AsFeaturedIn)

	if err != nil {
		if r.db.IsRecordNotFoundError(err) {
			return nil, errs.ErrRecordNotFound
		}
		return nil, err
	}

	return &AsFeaturedIn, nil
}

func (r *AsFeaturedInRepo) UpdateAsFeaturedIn(form models.AsFeaturedInUpdateForm) (*models.AsFeaturedIn, error) {
	var update models.AsFeaturedIn

	var err = copier.Copy(&update, &form)
	if err != nil {
		return nil, eris.Wrap(err, err.Error())
	}

	if form.Status == enums.PostStatusPublished {
		update.PublishedAt = values.Int64(time.Now().Unix())
	}

	if form.VI != nil && form.VI.Status == enums.PostStatusPublished {
		update.VI.PublishedAt = values.Int64(time.Now().Unix())
	}

	err = r.db.Omit(clause.Associations).Model(&update).Where("id = ?", form.AsFeaturedInID).Updates(&update).Error
	if err != nil {
		return nil, err
	}

	return r.GetAsFeaturedIn(GetAsFeaturedInParams{
		AsFeaturedInID: form.AsFeaturedInID,
	})
}

type DeleteAsFeaturedInParams struct {
	AsFeaturedInID string `param:"ads_video_id" validate:"required"`

	models.JwtClaimsInfo
}

func (r *AsFeaturedInRepo) DeleteAsFeaturedIn(params DeleteAsFeaturedInParams) error {
	var cate models.AsFeaturedIn
	var err = r.db.First(&cate, "id = ?", params.AsFeaturedInID).Error
	if cate.ID == "" {
		return errs.ErrCategoryNotFound
	}

	err = r.db.Transaction(func(tx *gorm.DB) error {
		err = tx.Unscoped().Delete(&models.AsFeaturedIn{}, "id = ?", params.AsFeaturedInID).Error

		return err
	})

	return err
}
