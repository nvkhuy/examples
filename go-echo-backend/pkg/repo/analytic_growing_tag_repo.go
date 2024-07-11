package repo

import (
	"github.com/engineeringinflow/inflow-backend/pkg/db"
	"github.com/engineeringinflow/inflow-backend/pkg/logger"
	"github.com/engineeringinflow/inflow-backend/pkg/models"
)

type GrowingTagRepo struct {
	db     *db.DB
	adb    *db.DB
	logger *logger.Logger
}

func NewAnalyticGrowingTagRepo(adb *db.DB) *GrowingTagRepo {
	return &GrowingTagRepo{
		adb:    adb,
		logger: logger.New("repo/analytic_growing_tags"),
	}
}

func (r *GrowingTagRepo) WithDB(db *db.DB) *GrowingTagRepo {
	r.db = db
	return r
}

type PaginateAnalyticGrowingTagsParams struct {
	models.JwtClaimsInfo
	models.PaginationParams
}

func (r *GrowingTagRepo) Paginate(params PaginateAnalyticGrowingTagsParams) (results *models.AnalyticGrowingTags, err error) {
	if params.Limit == 0 {
		params.Limit = 4
	}
	err = r.adb.Model(&models.AnalyticGrowingTag{}).Limit(params.Limit).Order("created_at DESC").Find(&results).Error
	return
}
