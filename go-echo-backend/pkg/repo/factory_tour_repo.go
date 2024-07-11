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

type FactoryTourRepo struct {
	db     *db.DB
	logger *logger.Logger
}

func NewFactoryTourRepo(db *db.DB) *FactoryTourRepo {
	return &FactoryTourRepo{
		db:     db,
		logger: logger.New("repo/FactoryTour"),
	}
}

type PaginateFactoryToursParams struct {
	models.PaginationParams

	models.JwtClaimsInfo
}

func (r *FactoryTourRepo) PaginateFactoryTours(params PaginateFactoryToursParams) *query.Pagination {
	var builder = queryfunc.NewFactoryTourBuilder(queryfunc.FactoryTourBuilderOptions{
		QueryBuilderOptions: queryfunc.QueryBuilderOptions{
			Role: params.GetRole(),
		},
	})

	var result = query.New(r.db, builder).
		WhereFunc(func(builder *query.Builder) {
		}).
		Page(params.Page).
		Limit(params.Limit).
		PagingFunc()

	return result
}

type SearchFactoryToursParams struct {
	models.PaginationParams
	models.JwtClaimsInfo
}

func (r *FactoryTourRepo) CreateFactoryTour(form models.FactoryTourUpdateForm) (*models.FactoryTour, error) {
	var FactoryTour models.FactoryTour
	err := copier.Copy(&FactoryTour, &form)
	if err != nil {
		return nil, eris.Wrap(err, "")
	}

	err = r.db.Omit(clause.Associations).Create(&FactoryTour).Error
	if err != nil {
		if r.db.IsRecordNotFoundError(err) {
			return nil, errs.ErrRecordNotFound
		}
		return nil, eris.Wrap(err, "")
	}

	return &FactoryTour, nil
}
