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
	"github.com/lib/pq"

	"github.com/rotisserie/eris"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type AdsVideoRepo struct {
	db     *db.DB
	logger *logger.Logger
}

func NewAdsVideoRepo(db *db.DB) *AdsVideoRepo {
	return &AdsVideoRepo{
		db:     db,
		logger: logger.New("repo/AdsVideo"),
	}
}

type PaginateAdsVideoParams struct {
	models.PaginationParams
	models.JwtClaimsInfo

	Sections []enums.AdsVideoSection `json:"sections" query:"sections" form:"sections"`
}

func (r *AdsVideoRepo) PaginateAdsVideo(params PaginateAdsVideoParams) *query.Pagination {
	var builder = queryfunc.NewAdsVideoBuilder(queryfunc.AdsVideoBuilderOptions{
		QueryBuilderOptions: queryfunc.QueryBuilderOptions{
			Role: params.GetRole(),
		},
	})
	if params.Limit == 0 {
		params.Limit = 25
	}

	var result = query.New(r.db, builder).
		WhereFunc(func(builder *query.Builder) {
			if len(params.Sections) > 0 {
				builder.Where("ads.sections && ?", pq.Array(params.Sections))
			}
		}).
		Page(params.Page).
		Limit(params.Limit).
		PagingFunc()

	return result
}

func (r *AdsVideoRepo) CreateAdsVideo(form models.AdsVideoCreateForm) (*models.AdsVideo, error) {
	var bc models.AdsVideo
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

type GetAdsVideoParams struct {
	AdsVideoID string `param:"ads_video_id" query:"ads_video_id" form:"ads_video_id" validate:"required"`

	models.JwtClaimsInfo
}

func (r *AdsVideoRepo) GetAdsVideo(params GetAdsVideoParams) (*models.AdsVideo, error) {
	var builder = queryfunc.NewAdsVideoBuilder(queryfunc.AdsVideoBuilderOptions{
		QueryBuilderOptions: queryfunc.QueryBuilderOptions{
			Role: params.GetRole(),
		},
	})
	var AdsVideo models.AdsVideo
	var err = query.New(r.db, builder).
		WhereFunc(func(builder *query.Builder) {
			builder.Where("ads.id = ?", params.AdsVideoID)
		}).
		FirstFunc(&AdsVideo)

	if err != nil {
		if r.db.IsRecordNotFoundError(err) {
			return nil, errs.ErrRecordNotFound
		}
		return nil, err
	}

	return &AdsVideo, nil
}

func (r *AdsVideoRepo) UpdateAdsVideo(form models.AdsVideoUpdateForm) (*models.AdsVideo, error) {
	var update models.AdsVideo

	var err = copier.Copy(&update, &form)
	if err != nil {
		return nil, eris.Wrap(err, err.Error())
	}

	err = r.db.Omit(clause.Associations).Model(&update).Where("id = ?", form.AdsVideoID).Updates(&update).Error
	if err != nil {
		return nil, err
	}

	return r.GetAdsVideo(GetAdsVideoParams{
		AdsVideoID: form.AdsVideoID,
	})
}

type DeleteAdsVideoParams struct {
	AdsVideoID string `param:"ads_video_id" validate:"required"`

	models.JwtClaimsInfo
}

func (r *AdsVideoRepo) DeleteAdsVideo(params DeleteAdsVideoParams) error {
	var cate models.AdsVideo
	var err = r.db.First(&cate, "id = ?", params.AdsVideoID).Error
	if cate.ID == "" {
		return errs.ErrCategoryNotFound
	}

	err = r.db.Transaction(func(tx *gorm.DB) error {
		err = tx.Unscoped().Delete(&models.AdsVideo{}, "id = ?", params.AdsVideoID).Error

		return err
	})

	return err
}
