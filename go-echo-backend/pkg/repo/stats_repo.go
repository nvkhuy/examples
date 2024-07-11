package repo

import (
	"github.com/engineeringinflow/inflow-backend/pkg/db"
	"github.com/engineeringinflow/inflow-backend/pkg/logger"
	"github.com/engineeringinflow/inflow-backend/pkg/models"
	"github.com/engineeringinflow/inflow-backend/pkg/models/enums"
	"github.com/engineeringinflow/inflow-backend/pkg/repo/query"
	"github.com/engineeringinflow/inflow-backend/pkg/repo/query/queryfunc"
)

type StatsRepo struct {
	db     *db.DB
	logger *logger.Logger
}

func NewStatsRepo(db *db.DB) *StatsRepo {
	return &StatsRepo{
		db:     db,
		logger: logger.New("repo/stats"),
	}
}

type StatsSuppliersParams struct {
	models.JwtClaimsInfo

	FromTime int `json:"from_time" query:"from_time" form:"from_time"`
	ToTime   int `json:"to_time" query:"to_time" form:"to_time" validate:"omitempty,gtfield=FromTime"`
}

func (r *StatsRepo) StatsSuppliers(params StatsSuppliersParams) models.StatsSuppliers {
	var builder = queryfunc.NewStatsSuppliersBuilder(queryfunc.StatsSuppliersBuilderOptions{
		QueryBuilderOptions: queryfunc.QueryBuilderOptions{
			Role: params.GetRole(),
		},
	})
	var resp models.StatsSuppliers

	query.New(r.db, builder).
		WhereFunc(func(builder *query.Builder) {
			if params.FromTime > 0 {
				builder.Where("sub.created_at <= ", params.FromTime)
			}

			if params.ToTime > 0 {
				builder.Where("sub.created_at >= ", params.ToTime)
			}
		}).
		Scan(&resp)

	return resp
}

type StatsBuyersParams struct {
	FromTime int `json:"from_time" query:"from_time" form:"from_time"`
	ToTime   int `json:"to_time" query:"to_time" form:"to_time" validate:"omitempty,gtfield=FromTime"`

	ForRole enums.Role
}

func (r *StatsRepo) StatsBuyers(params StatsBuyersParams) models.StatsBuyers {
	var builder = queryfunc.NewStatsBuyersBuilder(queryfunc.StatsBuyersBuilderOptions{
		ForRole: params.ForRole,
	})
	var resp models.StatsBuyers

	query.New(r.db, builder).
		WhereFunc(func(builder *query.Builder) {
			if params.FromTime > 0 {
				builder.Where("u.created_at <= ", params.FromTime)
			}

			if params.ToTime > 0 {
				builder.Where("u.created_at >= ", params.ToTime)
			}
		}).
		Scan(&resp)

	return resp
}

type StatsProductsParams struct {
	FromTime int `json:"from_time" query:"from_time" form:"from_time"`
	ToTime   int `json:"to_time" query:"to_time" form:"to_time" validate:"omitempty,gtfield=FromTime"`

	ForRole enums.Role
}

func (r *StatsRepo) StatsProducts(params StatsProductsParams) models.StatsProducts {
	var builder = queryfunc.NewStatsProductsBuilder(queryfunc.StatsProductsBuilderOptions{
		ForRole: params.ForRole,
	})
	var resp models.StatsProducts

	query.New(r.db, builder).
		WhereFunc(func(builder *query.Builder) {
			if params.FromTime > 0 {
				builder.Where("p.created_at <= ", params.FromTime)
			}

			if params.ToTime > 0 {
				builder.Where("p.created_at >= ", params.ToTime)
			}
		}).
		Scan(&resp)

	return resp
}

type StatsCategoriesParams struct {
	FromTime int `json:"from_time" query:"from_time" form:"from_time"`
	ToTime   int `json:"to_time" query:"to_time" form:"to_time" validate:"omitempty,gtfield=FromTime"`

	ForRole enums.Role
}

func (r *StatsRepo) StatsCategories(params StatsProductsParams) models.StatsCategories {
	var builder = queryfunc.NewStatsCategoriesBuilder(queryfunc.StatsCategoriesBuilderOptions{
		ForRole: params.ForRole,
	})
	var resp models.StatsCategories

	query.New(r.db, builder).
		WhereFunc(func(builder *query.Builder) {
			if params.FromTime > 0 {
				builder.Where("p.created_at <= ", params.FromTime)
			}

			if params.ToTime > 0 {
				builder.Where("p.created_at >= ", params.ToTime)
			}
		}).
		Scan(&resp)

	return resp
}
