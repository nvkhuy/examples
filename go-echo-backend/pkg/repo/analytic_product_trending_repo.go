package repo

import (
	"database/sql"
	"errors"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/engineeringinflow/inflow-backend/pkg/db"
	"github.com/engineeringinflow/inflow-backend/pkg/logger"
	"github.com/engineeringinflow/inflow-backend/pkg/models"
	"github.com/engineeringinflow/inflow-backend/pkg/models/enums"
	"github.com/engineeringinflow/inflow-backend/pkg/models/price"
	"github.com/engineeringinflow/inflow-backend/pkg/repo/query"
	"github.com/engineeringinflow/inflow-backend/pkg/repo/query/queryfunc"
	"github.com/samber/lo"
	"hash/fnv"
	"strings"
)

type ProductTrending struct {
	db     *db.DB
	adb    *db.DB
	logger *logger.Logger
}

func NewProductTrending(adb *db.DB) *ProductTrending {
	return &ProductTrending{
		adb:    adb,
		logger: logger.New("repo/product_trending"),
	}
}

func (r *ProductTrending) WithDB(db *db.DB) *ProductTrending {
	r.db = db
	return r
}

type PaginateProductTrendingParams struct {
	models.PaginationParams
	models.JwtClaimsInfo
	ProductIDs []string `json:"product_ids" query:"product_ids"`
}

func (r *ProductTrending) PaginateProductTrendings(params PaginateProductTrendingParams) *query.Pagination {
	var builder = queryfunc.NewProductTrendingBuilder(queryfunc.ProductTrendingBuilderOptions{
		QueryBuilderOptions: queryfunc.QueryBuilderOptions{
			Role: params.GetRole(),
		},
		DB: r.db,
	})
	params.ProductIDs = lo.Union(params.ProductIDs)

	if params.Limit == 0 {
		params.Limit = 12
	}
	if params.Page == 0 {
		params.Page = 1
	}

	var result = query.New(r.adb, builder).
		WhereFunc(func(builder *query.Builder) {
			if len(params.ProductIDs) > 0 {
				builder.Where("id IN ?", params.ProductIDs)
			}
			if keyword := strings.TrimSpace(params.Keyword); keyword != "" {
				var q = "%" + keyword + "%"
				builder.Where("(pt.name ILIKE @keyword OR pt.category ILIKE @keyword OR pt.sub_category ILIKE @keyword)", sql.Named("keyword", q))
			}
			if r.db.Configuration.IsProd() {
				builder.Where("pt.is_publish = ?", true)
			}
		}).
		Page(params.Page).
		Limit(params.Limit).
		OrderBy("pt.created_at DESC").
		PagingFunc()

	return result
}

type PaginateProductTrendingGroupParams struct {
	models.PaginationParams
	models.JwtClaimsInfo
	Domains       []enums.Domain `json:"domains,omitempty" query:"domains" params:"domains" form:"domains"`
	Categories    []string       `json:"categories,omitempty" query:"categories" params:"categories" form:"categories"`
	SubCategories []string       `json:"sub_categories" query:"sub_categories" params:"sub_categories" form:"sub_categories"`
	StartDate     int64          `json:"start_date,omitempty" query:"start_date" params:"start_date" form:"start_date"`
	EndDate       int64          `json:"end_date,omitempty" query:"end_date" params:"end_date" form:"end_date"`
}

func (r *ProductTrending) PaginateProductTrendingGroup(params PaginateProductTrendingGroupParams) *query.Pagination {
	var builder = queryfunc.NewProductTrendingGroupBuilder(queryfunc.ProductTrendingGroupBuilderOptions{
		QueryBuilderOptions: queryfunc.QueryBuilderOptions{
			Role: params.GetRole(),
		},
		DB: r.db,
	})

	if params.Limit == 0 {
		params.Limit = 12
	}
	if params.Page == 0 {
		params.Page = 1
	}

	var result = query.New(r.adb, builder).
		WhereFunc(func(builder *query.Builder) {
			if params.StartDate > 0 && params.EndDate > 0 {
				builder.Where("created_at >= ? AND created_at <= ?", params.StartDate, params.EndDate)
			}
			if len(params.Domains) > 0 {
				builder.Where("ptg.domain IN ?", params.Domains)
			}
			if len(params.Categories) > 0 {
				builder.Where("ptg.category IN ?", params.Categories)
			}
			if len(params.SubCategories) > 0 {
				builder.Where("ptg.sub_category IN ?", params.SubCategories)
			}
			if keyword := strings.TrimSpace(params.Keyword); keyword != "" {
				var q = "%" + keyword + "%"
				builder.Where("(ptg.name ILIKE @keyword OR ptg.category ILIKE @keyword OR ptg.sub_category ILIKE @keyword)", sql.Named("keyword", q))
			}
			if r.db.Configuration.IsProd() {
				builder.Where("ptg.is_publish = ?", true)
			}
		}).
		Page(params.Page).
		Limit(params.Limit).
		OrderBy("ptg.created_at DESC").
		PagingFunc()

	return result
}

type ListProductTrendingParams struct {
	models.PaginationParams
	models.JwtClaimsInfo
	ProductIDs []string `json:"product_ids" query:"product_ids"`
}

func (r *ProductTrending) ListProductTrendings(params ListProductTrendingParams) (products []models.AnalyticProductTrending, err error) {
	var builder = queryfunc.NewProductTrendingBuilder(queryfunc.ProductTrendingBuilderOptions{
		QueryBuilderOptions: queryfunc.QueryBuilderOptions{
			Role: params.GetRole(),
		},
	})

	err = query.New(r.adb, builder).
		WhereFunc(func(builder *query.Builder) {
			if len(params.ProductIDs) > 0 {
				builder.Where("id IN ?", params.ProductIDs)
			}
		}).
		Page(params.Page).Limit(params.Limit).OrderBy("p.created_at DESC").
		Find(&products)
	return
}

type ListProductTrendingDomainParams struct {
	models.PaginationParams
	models.JwtClaimsInfo
}

func (r *ProductTrending) ListProductTrendingDomain(params ListProductTrendingDomainParams) (domainStats []models.AnalyticProductTrendingStats, err error) {
	err = r.adb.Model(&models.AnalyticProductTrendingStats{}).Select("domain,count(1)").Group("domain").Having("count(1) > 0").Order("count(1) DESC").Find(&domainStats).Error
	return
}

type ListProductTrendingCategoryParams struct {
	models.PaginationParams
	models.JwtClaimsInfo
}

func (r *ProductTrending) ListProductTrendingCategory(params ListProductTrendingCategoryParams) (categoryStats []models.AnalyticProductTrendingStats, err error) {
	err = r.adb.Model(&models.AnalyticProductTrendingStats{}).Select("category,count(1)").Group("category").Having("count(1) > 0").Order("count(1) DESC").Find(&categoryStats).Error
	return
}

func (r *ProductTrending) ListProductTrendingSubCategory(params ListProductTrendingCategoryParams) (categoryStats []models.AnalyticProductTrendingStats, err error) {
	err = r.adb.Model(&models.AnalyticProductTrendingStats{}).Select("sub_category,count(1)").Group("sub_category").Having("count(1) > 0").Order("count(1) DESC").Find(&categoryStats).Error
	return
}

type PaginateProductTrendingChartParams struct {
	models.JwtClaimsInfo
	models.PaginationParams
	ProductTrendingID string `json:"product_trending_id" param:"product_trending_id" query:"product_trending_id" form:"product_trending_id" validate:"required"`
	ProductID         string
	DateFrom          int64                `json:"date_from" query:"date_from" form:"date_from" param:"date_from"`
	DateTo            int64                `json:"date_to" query:"date_to" form:"date_to" param:"date_to"`
	PredictNext       int                  `json:"predict_next" query:"predict_next" form:"predict_next" param:"predict_next"`
	PredictOn         enums.PredictChartOn `json:"predict_on" query:"predict_on" form:"predict_on" param:"predict_on" validate:"required"`
}

func (r *ProductTrending) Chart(params PaginateProductTrendingChartParams) (pag *query.Pagination) {
	if params.Limit == 0 {
		params.Limit = 20
	}
	params.ProductID = r.toProductId(params.ProductTrendingID)

	var builder = queryfunc.NewChartDAProductBuilder(queryfunc.ChartDAProductBuilderOptions{
		QueryBuilderOptions: queryfunc.QueryBuilderOptions{
			Role: params.GetRole(),
		},
		PredictNext: params.PredictNext,
		PredictOn:   params.PredictOn,
	})

	var queryBuilder = query.New(r.adb, builder).
		WhereFunc(func(builder *query.Builder) {
			if params.ProductID != "" {
				var product models.AnalyticProduct
				r.adb.Model(&models.AnalyticProduct{}).Select("id", "url").Where("id = ?", params.ProductID).First(&product)
				if product.URL != "" {
					builder.Where("pc.url = ?", product.URL)
				}
			}
			if params.DateFrom > 0 && params.DateTo > 0 {
				builder.Where("pc.created_at >= ?", params.DateFrom).Where("pc.created_at <= ?", params.DateTo)
			}
		})
	if params.PredictNext > 0 {
		queryBuilder.GroupBy("pc.url,pc.scrape_date").OrderBy("scrape_date DESC")
	}

	return queryBuilder.Page(params.Page).Limit(params.Limit).PagingFunc()
}

type GetProductTrendingParam struct {
	models.JwtClaimsInfo
	ProductTrendingID string `json:"product_trending_id" param:"product_trending_id" query:"product_trending_id" form:"product_trending_id" validate:"required"`
	ProductID         string
}

func (r *ProductTrending) Get(params GetProductTrendingParam) (result models.AnalyticProductTrending, err error) {
	params.ProductID = r.toProductId(params.ProductTrendingID)
	analyticProduct, err := NewDataAnalyticRepo(r.adb).GetProduct(DataAnalyticGetProductParam{
		ProductID: params.ProductID,
	})
	if err != nil {
		return
	}

	var builder = queryfunc.NewProductTrendingBuilder(queryfunc.ProductTrendingBuilderOptions{
		QueryBuilderOptions: queryfunc.QueryBuilderOptions{
			Role: params.GetRole(),
		},
		DB: r.db,
	})

	err = query.New(r.adb, builder).
		WhereFunc(func(builder *query.Builder) {
			if params.ProductTrendingID != "" {
				builder.Where("pt.id = ?", params.ProductTrendingID)
			}
		}).FirstFunc(&result)
	result.Price = price.NewFromFloat(analyticProduct.Price)
	result.OverallGrowthRate = analyticProduct.OverallGrowthRate
	return
}

func (r *ProductTrending) toProductId(key string) string {
	var setupIds = []string{"cmeh5hbb2hjc375m021g", "cnjco9rb2hjdqovfqjng", "cnprm4jb2hjdqovhgtb0", "cm4ikiqlk3mikinbs4v0", "cm4ik4qlk3mikinbrqk0"}
	hash := fnv.New32a()
	hash.Write([]byte(key))
	idx := int(hash.Sum32()) % len(setupIds)
	return setupIds[idx]
}

// CUD

type CreateProductTrendingParams struct {
	models.JwtClaimsInfo
	Products []*models.AnalyticProductTrending `json:"products,omitempty"`
}

func (r *ProductTrending) CreateProductTrending(params CreateProductTrendingParams) (product []*models.AnalyticProductTrending, err error) {
	for _, p := range params.Products {
		p.Domain = enums.Inflow
		p.IsPublish = aws.Bool(false)
	}

	err = r.adb.Model(&models.AnalyticProductTrending{}).Create(&params.Products).Error
	if err != nil {
		return
	}
	product = params.Products
	return
}

type UpdateProductTrendingParams struct {
	models.JwtClaimsInfo
	Products []*models.AnalyticProductTrending `json:"products,omitempty"`
}

func (r *ProductTrending) UpdateProductTrending(params UpdateProductTrendingParams) (product []*models.AnalyticProductTrending, err error) {
	for _, p := range params.Products {
		if p.ID == "" {
			err = errors.New("empty product id")
			return
		}
	}
	for _, p := range params.Products {
		err = r.adb.Model(&models.AnalyticProductTrending{}).Where("id = ?", p.ID).Updates(&p).Error
		if err != nil {
			return
		}
	}
	product = params.Products
	return
}

type DeleteProductTrendingParams struct {
	models.JwtClaimsInfo
	ProductIDs []string `json:"product_ids" query:"product_ids" params:"product_ids" validate:"required"`
}

func (r *ProductTrending) DeleteProductTrending(params DeleteProductTrendingParams) (err error) {
	if len(params.ProductIDs) == 0 {
		return
	}
	err = r.adb.Unscoped().Delete(&models.AnalyticProductTrending{}, "id IN ?", params.ProductIDs).Error
	return
}
