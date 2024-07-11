package repo

import (
	"strings"

	"github.com/engineeringinflow/inflow-backend/pkg/db"
	"github.com/engineeringinflow/inflow-backend/pkg/models"
	"github.com/engineeringinflow/inflow-backend/pkg/repo/query"
	"github.com/engineeringinflow/inflow-backend/pkg/repo/query/queryfunc"
	"github.com/jinzhu/copier"
)

type SettingSEORepo struct {
	db *db.DB
}

func NewSettingSEORepo(db *db.DB) *SettingSEORepo {
	return &SettingSEORepo{
		db: db,
	}
}

type PaginateSettingSEOParams struct {
	models.PaginationParams
	models.JwtClaimsInfo
}

func (r *SettingSEORepo) PaginateSettingSEO(params PaginateSettingSEOParams) (qp *query.Pagination) {
	var result = query.New(r.db, queryfunc.NewSettingSEOBuilder(queryfunc.SettingSEOBuilderOptions{
		QueryBuilderOptions: queryfunc.QueryBuilderOptions{
			Role: params.GetRole(),
		},
	})).
		Limit(params.Limit).
		Page(params.Page).
		PagingFunc()
	return result
}

func (r *SettingSEORepo) PaginateSettingSEOLanguageGroup(params PaginateSettingSEOParams) (qp *query.Pagination) {
	var result = query.New(r.db, queryfunc.NewSettingSEOLanguageGroupBuilder(queryfunc.SettingSEOBuilderOptions{
		QueryBuilderOptions: queryfunc.QueryBuilderOptions{
			Role: params.GetRole(),
		},
	})).
		Limit(params.Limit).
		Page(params.Page).
		PagingFunc()
	return result
}

type GetSettingSEOByRouteNameParams struct {
	models.JwtClaimsInfo

	RouteName string `json:"route_name" query:"route_name" validate:"required"`
}

func (r *SettingSEORepo) GetSettingSEOByRouteName(params GetSettingSEOByRouteNameParams) (*models.SettingSEOLanguageGroup, error) {
	if strings.HasPrefix(params.RouteName, "/product/") {
		var slug = strings.ReplaceAll(params.RouteName, "/product/", "")
		var product models.Product
		var err = r.db.First(&product, "slug = ?", slug).Error
		if err != nil {
			return nil, err
		}

		return &models.SettingSEOLanguageGroup{
			Route: params.RouteName,
			EN: models.SettingSEO{
				Route:       params.RouteName,
				Title:       product.Name,
				Description: product.ShortDescription,
			},
			VI: models.SettingSEO{
				Route:       params.RouteName,
				Title:       product.Name,
				Description: product.ShortDescription,
			},
		}, nil
	}

	var result models.SettingSEOLanguageGroup

	var err = query.New(r.db, queryfunc.NewSettingSEOLanguageGroupBuilder(queryfunc.SettingSEOBuilderOptions{
		QueryBuilderOptions: queryfunc.QueryBuilderOptions{
			Role: params.GetRole(),
		},
	})).WhereFunc(func(builder *query.Builder) {
		builder.Where("ss.route = ?", params.RouteName)
	}).
		FirstFunc(&result)

	return &result, err
}

func (r *SettingSEORepo) CreateSettingSEO(form models.CreateSettingSEOForm) (settingSEO models.SettingSEO, err error) {
	err = copier.Copy(&settingSEO, &form)
	if err != nil {
		return models.SettingSEO{}, err
	}
	err = r.db.Create(&settingSEO).Error
	return
}

func (r *SettingSEORepo) UpdateSettingSEO(form models.UpdateSettingSEOForm) (settingSEO models.SettingSEO, err error) {
	err = copier.Copy(&settingSEO, &form)
	if err != nil {
		return models.SettingSEO{}, err
	}
	err = r.db.Updates(&settingSEO).Error
	return
}

func (r *SettingSEORepo) DeleteSettingSEOs(form models.DeleteSettingSEOForm) (err error) {
	err = r.db.Unscoped().Delete(&models.SettingSEO{}, "id = ?", form.ID).Error
	return
}
