package repo

import (
	"database/sql"
	"fmt"
	"log"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/engineeringinflow/inflow-backend/pkg/db"
	"github.com/engineeringinflow/inflow-backend/pkg/logger"
	"github.com/engineeringinflow/inflow-backend/pkg/models"
	"github.com/engineeringinflow/inflow-backend/pkg/models/enums"
	"github.com/engineeringinflow/inflow-backend/pkg/repo/query"
	"github.com/engineeringinflow/inflow-backend/pkg/repo/query/queryfunc"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type DataAnalyticRepo struct {
	db     *db.DB
	adb    *db.DB
	logger *logger.Logger
}

func NewDataAnalyticRepo(adb *db.DB) *DataAnalyticRepo {
	return &DataAnalyticRepo{
		adb:    adb,
		logger: logger.New("repo/data_analytic"),
	}
}

func (r *DataAnalyticRepo) WithDB(db *db.DB) *DataAnalyticRepo {
	r.db = db
	return r
}

type DataAnalyticPlatformOverviewParam struct {
	models.PaginationParams
	models.JwtClaimsInfo
}

func (r *DataAnalyticRepo) Overview(params DataAnalyticPlatformOverviewParam) (pag *query.Pagination) {
	var overviews models.DataAnalyticPlatformOverviews
	err := r.adb.Model(&models.AnalyticProduct{}).
		Select("domain, country_code, count(1) as total").
		Group("domain, country_code").
		Order("count(1) desc").
		Find(&overviews).Error
	if err != nil {
		log.Println(err)
	}
	pag = &query.Pagination{
		Records:            overviews,
		TotalRecord:        len(overviews),
		TotalPage:          1,
		TotalCurrentRecord: len(overviews),
	}
	return
}

type DataAnalyticSearchProductsParam struct {
	models.PaginationParams
	models.JwtClaimsInfo
	CountryCodes     []string `json:"country_codes" param:"country_codes" query:"country_codes" form:"country_codes"`
	Domains          []string `json:"domains" param:"domains" query:"domains" form:"domains"`
	SortBy           string   `json:"sort_by" param:"sort_by" query:"sort_by" form:"sort_by"`
	IsSortDescending bool     `json:"is_sort_descending" param:"is_sort_descending" query:"is_sort_descending" form:"is_sort_descending"`
	GrowthRateFrom   *float64 `json:"growth_rate_from" param:"growth_rate_from" query:"growth_rate_from" form:"growth_rate_from"`
	GrowthRateTo     *float64 `json:"growth_rate_to" param:"growth_rate_to" query:"growth_rate_to" form:"growth_rate_to"`
	ProductClasses   []string `json:"product_classes" param:"product_classes" query:"product_classes" form:"product_classes"`
	SubCategories    []string `json:"sub_categories" param:"sub_categories" query:"sub_categories" form:"sub_categories"`
}

func (r *DataAnalyticRepo) SearchProducts(params DataAnalyticSearchProductsParam) (pag *query.Pagination) {
	var builder = queryfunc.NewAnalyticProductBuilder(queryfunc.AnalyticProductBuilderOptions{
		QueryBuilderOptions: queryfunc.QueryBuilderOptions{
			Role: params.GetRole(),
		},
		IncludeProductClass: true,
	})

	if params.Limit == 0 {
		params.Limit = 20
	}

	var querying = query.New(r.adb, builder).
		WhereFunc(func(builder *query.Builder) {
			if len(params.ProductClasses) > 0 {
				builder.Where("pc.class IN ?", params.ProductClasses)
			}
			if len(params.SubCategories) > 0 {
				builder.Where("p.sub_category IN ?", params.SubCategories)
			}
			if len(params.Domains) > 0 {
				builder.Where("p.domain IN ?", params.Domains)
			}
			if len(params.CountryCodes) > 0 {
				builder.Where("p.country_code IN ?", params.CountryCodes)
			}
			if params.GrowthRateFrom != nil && params.GrowthRateTo != nil {
				builder.Where("p.overall_growth_rate > ? AND overall_growth_rate < ?", *params.GrowthRateFrom, *params.GrowthRateTo)
			}
			if keyword := strings.TrimSpace(params.Keyword); keyword != "" {
				var q = "%" + keyword + "%"
				builder.Where("(p.name ILIKE @keyword OR p.category ILIKE @keyword OR p.sub_category ILIKE @keyword)", sql.Named("keyword", q))
			}
		})

	if params.SortBy == "created_at" {
		if params.IsSortDescending {
			querying.OrderBy(" p.created_at DESC")
		} else {
			querying.OrderBy(" p.created_at")
		}
	} else if params.SortBy == "updated_at" {
		if params.IsSortDescending {
			querying.OrderBy(" p.updated_at DESC")
		} else {
			querying.OrderBy(" p.updated_at")
		}
	} else if params.SortBy == "growth_rate" {
		if params.IsSortDescending {
			querying.OrderBy("p.overall_growth_rate DESC")
		} else {
			querying.OrderBy("p.overall_growth_rate")
		}
	} else {
		querying.OrderBy("p.overall_growth_rate DESC")
	}

	return querying.Page(params.Page).Limit(params.Limit).PagingFunc()
}

func (r *DataAnalyticRepo) ProductGroupURL(params DataAnalyticSearchProductsParam) (pag *query.Pagination) {
	var builder = queryfunc.NewAnalyticProductGroupBuilder(queryfunc.AnalyticProductBuilderOptions{
		QueryBuilderOptions: queryfunc.QueryBuilderOptions{
			Role: params.GetRole(),
		},
		IncludeProductClass: true,
	})

	if params.Limit == 0 {
		params.Limit = 20
	}

	var querying = query.New(r.adb, builder).
		WhereFunc(func(builder *query.Builder) {
			if len(params.ProductClasses) > 0 {
				builder.Where("pc.class IN ?", params.ProductClasses)
			}
			if len(params.SubCategories) > 0 {
				builder.Where("p.sub_category IN ?", params.SubCategories)
			}
			if len(params.Domains) > 0 {
				builder.Where("p.domain IN ?", params.Domains)
			}
			if len(params.CountryCodes) > 0 {
				builder.Where("p.country_code IN ?", params.CountryCodes)
			}
			if params.GrowthRateFrom != nil && params.GrowthRateTo != nil {
				builder.Where("p.overall_growth_rate > ? AND overall_growth_rate < ?", *params.GrowthRateFrom, *params.GrowthRateTo)
			}

			builder.Where("p.images is not null")

			if r.adb.Configuration.IsProd() {
				builder.Where("p.is_publish = ?", true)
			}

			if keyword := strings.TrimSpace(params.Keyword); keyword != "" {
				var q = "%" + keyword + "%"
				builder.Where("(p.name ILIKE @keyword OR p.category ILIKE @keyword OR p.sub_category ILIKE @keyword)", sql.Named("keyword", q))
			}
		})

	if params.SortBy == "created_at" {
		if params.IsSortDescending {
			querying.OrderBy(" p.created_at DESC")
		} else {
			querying.OrderBy(" p.created_at")
		}
	} else if params.SortBy == "updated_at" {
		if params.IsSortDescending {
			querying.OrderBy(" p.updated_at DESC")
		} else {
			querying.OrderBy(" p.updated_at")
		}
	} else if params.SortBy == "growth_rate" {
		if params.IsSortDescending {
			querying.OrderBy("p.overall_growth_rate DESC")
		} else {
			querying.OrderBy("p.overall_growth_rate")
		}
	} else {
		querying.OrderBy("p.overall_growth_rate DESC")
	}

	return querying.Page(params.Page).Limit(params.Limit).WithoutCount(true).PagingFunc()
}

type DataAnalyticRecommendProductsParam struct {
	models.JwtClaimsInfo
	models.PaginationParams
	ProductClass       string  `json:"product_class" param:"product_class" query:"product_class" form:"product_class"`
	ProductClasses     string  `json:"product_classes" param:"product_classes" query:"product_classes" form:"product_classes"`
	ProductID          string  `json:"product_id" param:"product_id" query:"product_id" form:"product_id"` //inflow product
	RecommendProductID string  `json:"recommend_product_id" param:"recommend_product_id" query:"recommend_product_id" form:"recommend_product_id"`
	ConfThreshold      float64 `json:"conf_threshold" param:"conf_threshold" query:"conf_threshold" form:"conf_threshold"`
}

func (r *DataAnalyticRepo) RecommendProducts(params DataAnalyticRecommendProductsParam) (pag *query.Pagination) {
	var builder = queryfunc.NewAnalyticProductBuilder(queryfunc.AnalyticProductBuilderOptions{
		IncludeProductClass: true,
	})

	if params.Limit == 0 {
		params.Limit = 20
	}
	if params.ConfThreshold == 0 {
		params.ConfThreshold = 0.7
	}

	var result = query.New(r.adb, builder).
		WhereFunc(func(builder *query.Builder) {
			if params.ProductID != "" {
				var class models.ProductClass
				r.db.Model(&models.ProductClass{}).Where("product_id = ?", params.ProductID).Order("conf DESC").First(&class)
				if class.Class != "" {
					builder.Where("pc.class = ?", class.Class).Where("pc.conf >= ?", params.ConfThreshold)
				}
			} else if params.RecommendProductID != "" {
				var class models.ProductClass
				r.adb.Model(&models.ProductClass{}).Where("product_id = ?", params.RecommendProductID).Order("conf DESC").First(&class)
				if class.Class != "" {
					builder.Where("pc.class = ?", class.Class).Where("pc.conf >= ?", params.ConfThreshold)
				}
			} else if params.ProductClass != "" {
				builder.Where("pc.class = ?", params.ProductClass)
			} else if len(params.ProductClasses) > 0 {
				builder.Where("pc.class IN ?", params.ProductClasses)
			}
			builder.Where("p.domain != ?", enums.Amazon).
				Where("p.images is not null").
				Where("COALESCE(p.category,'') != ''").
				Where("p.overall_growth_rate != 0 AND p.overall_growth_rate < 500").
				Where("pc.conf = (?)", gorm.Expr("SELECT MAX(conf) FROM product_classes WHERE product_id = p.id"))
		}).
		Page(params.Page).
		OrderBy("p.overall_growth_rate DESC").
		Limit(params.Limit).
		PagingFunc()

	return result
}

type DataAnalyticTopProductsParam struct {
	models.PaginationParams
	models.JwtClaimsInfo
	OrderBy      string   `json:"order_by" query:"order_by" form:"order_by" param:"order_by" validate:"required,oneof=sold"`
	DateFrom     int64    `json:"date_from" query:"date_from" form:"date_from" param:"date_from"`
	DateTo       int64    `json:"date_to" query:"date_to" form:"date_to" param:"date_to"`
	Category     string   `json:"category" param:"category" query:"category" form:"category"`
	SubCategory  string   `json:"sub_category" param:"sub_category" query:"sub_category" form:"sub_category"`
	Domain       string   `json:"domain" param:"domain" query:"domain" form:"domain"`
	CountryCodes []string `json:"country_codes" param:"country_codes" query:"country_codes" form:"country_codes"`
}

func (r *DataAnalyticRepo) TopProducts(params DataAnalyticTopProductsParam) (pag *query.Pagination) {
	var builder = queryfunc.NewRankDAProductBuilder(queryfunc.RankDAProductBuilderOptions{
		QueryBuilderOptions: queryfunc.QueryBuilderOptions{
			Role: params.GetRole(),
		},
		OrderBy: params.OrderBy,
	})

	if params.Limit == 0 {
		params.Limit = 20
	}

	var result = query.New(r.adb, builder).
		WhereFunc(func(builder *query.Builder) {
			builder.Where("p.created_at >= ? and p.created_at <= ?", params.DateFrom, params.DateTo)
			if params.Domain != "" {
				params.Domain = strings.TrimSpace(strings.ToLower(params.Domain))
				builder.Where("p.domain = ?", params.Domain)
			}
			if len(params.CountryCodes) > 0 {
				builder.Where("p.country_code IN ?", params.CountryCodes)
			}
			if category := strings.TrimSpace(params.Category); category != "" {
				builder.Where("p.category = ?", category)
			}
			if subcategory := strings.TrimSpace(params.SubCategory); subcategory != "" {
				builder.Where("p.sub_category = ?", subcategory)
			}
		}).
		Page(params.Page).
		Limit(params.Limit).
		PagingFunc()

	return result
}

type DataAnalyticTopCategoriesParam struct {
	models.PaginationParams
	models.JwtClaimsInfo
	Select      []string `json:"select,omitempty" query:"select" form:"select" param:"select"`
	SubCategory string   `json:"sub_category" query:"sub_category" form:"sub_category" param:"sub_category"`
	OrderBy     string   `json:"order_by" query:"order_by" form:"order_by" param:"order_by" validate:"required,oneof=sold rev"`
	DateFrom    int64    `json:"date_from" query:"date_from" form:"date_from" param:"date_from"`
	DateTo      int64    `json:"date_to" query:"date_to" form:"date_to" param:"date_to" `
}

func (r *DataAnalyticRepo) TopCategories(params DataAnalyticTopCategoriesParam) (pag *query.Pagination) {
	var builder = queryfunc.NewRankDACategoryBuilder(queryfunc.RankDACategoryBuilderOptions{
		QueryBuilderOptions: queryfunc.QueryBuilderOptions{
			Role: params.GetRole(),
		},
		PaginationParams: params.PaginationParams,
		Select:           params.Select,
		OrderBy:          params.OrderBy,
		DateFrom:         params.DateFrom,
		DateTo:           params.DateTo,
	})

	if params.Page == 0 {
		params.Page = 1
	}

	if params.Limit == 0 {
		params.Limit = 20
	}

	var result = query.New(r.adb, builder).
		WhereFunc(func(builder *query.Builder) {
			builder.Where("p.created_at >= ?", params.DateFrom).Where("p.created_at <= ?", params.DateTo)
			if len(params.Select) > 0 {
				builder.Where(fmt.Sprintf("%s != ''", strings.Join(params.Select, ",")))
			}
			if keyword := strings.TrimSpace(params.Keyword); keyword != "" {
				var q = "%" + keyword + "%"
				sel := strings.Join(params.Select, ",")
				builder.Where(fmt.Sprintf("%s ILIKE @keyword", sel), sql.Named("keyword", q))
			}
		}).
		Page(params.Page).
		Limit(params.Limit).
		PagingFunc()

	return result
}

type DataAnalyticTopMoverProductsParam struct {
	models.PaginationParams
	models.JwtClaimsInfo
	Month        int64    `json:"month"`
	Year         int64    `json:"year"`
	Category     string   `json:"category" param:"category" query:"category" form:"category"`
	SubCategory  string   `json:"sub_category" param:"sub_category" query:"sub_category" form:"sub_category"`
	Domain       string   `json:"domain" param:"domain" query:"domain" form:"domain"`
	DateFrom     int64    `json:"date_from" query:"date_from" form:"date_from" param:"date_from"`
	DateTo       int64    `json:"date_to" query:"date_to" form:"date_to" param:"date_to" `
	CountryCodes []string `json:"country_codes" param:"country_codes" query:"country_codes" form:"country_codes"`
}

func (r *DataAnalyticRepo) TopMoverProducts(params DataAnalyticTopMoverProductsParam) (pag *query.Pagination) {
	var builder = queryfunc.NewTopMoverProductBuilder(queryfunc.TopTopMoverProductBuilderOptions{
		QueryBuilderOptions: queryfunc.QueryBuilderOptions{
			Role: params.GetRole(),
		},
	})

	if params.Limit == 0 {
		params.Limit = 20
	}
	if params.DateFrom > 0 {
		timeObject := time.Unix(params.DateFrom, 0)

		// Extract month and year
		params.Month = int64(timeObject.Month())
		params.Year = int64(timeObject.Year())
	}

	var result = query.New(r.adb, builder).
		WhereFunc(func(builder *query.Builder) {
			builder.Where(map[string]interface{}{
				"month": params.Month,
				"year":  params.Year,
			})
			if params.Domain != "" {
				params.Domain = strings.TrimSpace(strings.ToLower(params.Domain))
				builder.Where("p.domain = ?", params.Domain)
			}
			if category := strings.TrimSpace(params.Category); category != "" {
				builder.Where("p.category = ?", category)
			}
			if subcategory := strings.TrimSpace(params.SubCategory); subcategory != "" {
				builder.Where("p.sub_category = ?", subcategory)
			}
			if len(params.CountryCodes) > 0 {
				builder.Where("p.country_code IN ?", params.CountryCodes)
			}
		}).
		Page(params.Page).
		Limit(params.Limit).
		PagingFunc()

	return result
}

type DataAnalyticGetProductParam struct {
	models.JwtClaimsInfo
	ProductID string `json:"product_id" param:"product_id" query:"product_id" form:"product_id"`
}

func (r *DataAnalyticRepo) GetProduct(params DataAnalyticGetProductParam) (result models.AnalyticProduct, err error) {
	var builder = queryfunc.NewAnalyticProductBuilder(queryfunc.AnalyticProductBuilderOptions{
		IncludeProductClass: true,
	})

	err = query.New(r.adb, builder).
		WhereFunc(func(builder *query.Builder) {
			if params.ProductID != "" {
				builder.Where("p.id = ?", params.ProductID)
			}
		}).
		FirstFunc(&result)
	return
}

type DataAnalyticGetDAProductChartParam struct {
	models.JwtClaimsInfo
	models.PaginationParams
	ProductID   string               `json:"product_id" param:"product_id" query:"product_id" form:"product_id" validate:"required"`
	DateFrom    int64                `json:"date_from" query:"date_from" form:"date_from" param:"date_from"`
	DateTo      int64                `json:"date_to" query:"date_to" form:"date_to" param:"date_to"`
	PredictNext int                  `json:"predict_next" query:"predict_next" form:"predict_next" param:"predict_next"`
	PredictOn   enums.PredictChartOn `json:"predict_on" query:"predict_on" form:"predict_on" param:"predict_on" validate:"required"`
}

func (r *DataAnalyticRepo) GetProductChart(params DataAnalyticGetDAProductChartParam) (pag *query.Pagination) {
	var builder = queryfunc.NewChartDAProductBuilder(queryfunc.ChartDAProductBuilderOptions{
		QueryBuilderOptions: queryfunc.QueryBuilderOptions{
			Role: params.GetRole(),
		},
		PredictNext: params.PredictNext,
		PredictOn:   params.PredictOn,
	})

	if params.Limit == 0 {
		params.Limit = 1000
	}

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

	return queryBuilder.Page(params.Page).Limit(params.Limit).WithoutCount(true).PagingFunc()
}

type GetAnalyticProductClassGroupParam struct {
	models.JwtClaimsInfo
	models.PaginationParams
}

func (r *DataAnalyticRepo) GetProductClassGroup(params GetAnalyticProductClassGroupParam) (pag *query.Pagination) {
	var builder = queryfunc.NewAnalyticProductClassBuilder(queryfunc.AnalyticProductClassBuilderOptions{
		QueryBuilderOptions: queryfunc.QueryBuilderOptions{
			Role: params.GetRole(),
		},
		IsGroupByClass: true,
	})

	if params.Limit == 0 {
		params.Limit = 20
	}

	var queryBuilder = query.New(r.adb, builder).
		WhereFunc(func(builder *query.Builder) {
			builder.GroupBy("pc.class").OrderBy("count(1) DESC")
		})

	return queryBuilder.Page(params.Page).Limit(params.Limit).PagingFunc()
}

type GetAnalyticProductTrendingGroupParam struct {
	models.JwtClaimsInfo
	models.PaginationParams

	SubCategories []string `json:"query" query:"sub_categories" form:"sub_categories" param:"sub_categories"`
}

func (r *DataAnalyticRepo) GetProductTrendingGroup(params GetAnalyticProductTrendingGroupParam) (pag *query.Pagination) {
	var builder = queryfunc.NewAnalyticProductTrendingBuilder(queryfunc.AnalyticProductTrendingBuilderOptions{
		QueryBuilderOptions: queryfunc.QueryBuilderOptions{
			Role: params.GetRole(),
		},
	})

	if params.Limit == 0 {
		params.Limit = 20
	}

	return query.New(r.adb, builder).
		Page(params.Page).
		Limit(params.Limit).
		WhereFunc(func(builder *query.Builder) {
			if len(params.SubCategories) > 0 {
				builder.Where("count_elements(tp.top_sub_categories,?) > 0", pq.StringArray(params.SubCategories))
			}

			if keyword := strings.TrimSpace(params.Keyword); keyword != "" {
				var q = "%" + keyword + "%"
				builder.Where("(tp.trending ILIKE @keyword)", sql.Named("keyword", q))
			}

		}).
		PagingFunc()

}

type DataAnalyticNewUsersParam struct {
	models.PaginationParams
	models.JwtClaimsInfo
	DateFrom int64 `json:"date_from" query:"date_from" form:"date_from" param:"date_from"`
	DateTo   int64 `json:"date_to" query:"date_to" form:"date_to" param:"date_to" `
}

func (r *DataAnalyticRepo) PaginateNewUsers(params DataAnalyticNewUsersParam) (result *models.DataAnalyticNewUser, err error) {
	type NewUserInDate struct {
		Date  string `json:"date"`
		Count int64  `json:"count"`
	}
	var data []NewUserInDate
	err = r.db.Model(&models.User{}).
		Select("TO_CHAR(TO_TIMESTAMP(created_at), 'YYYY/MM/DD') AS date, COUNT(1) AS count").
		Where("created_at >= ? AND created_at <= ?", params.DateFrom, params.DateTo).
		Group("date").
		Order("date DESC").
		Find(&data).Error
	if err != nil {
		return
	}
	result = &models.DataAnalyticNewUser{
		Count: int64(len(data)),
	}
	for _, v := range data {
		result.Charts = append(result.Charts, models.DataAnalyticChart{
			Date:  v.Date,
			Value: float64(v.Count),
		})
	}
	return
}

type DataAnalyticNewCatalogProductParam struct {
	models.PaginationParams
	models.JwtClaimsInfo
	DateFrom int64 `json:"date_from" query:"date_from" form:"date_from" param:"date_from"`
	DateTo   int64 `json:"date_to" query:"date_to" form:"date_to" param:"date_to" `
}

func (r *DataAnalyticRepo) PaginateNewCatalogProduct(params DataAnalyticNewCatalogProductParam) (result *models.DataAnalyticNewCatalogProduct, err error) {
	type NewProductInDate struct {
		Date  string `json:"date"`
		Count int64  `json:"count"`
	}
	var data []NewProductInDate
	err = r.db.Model(&models.Product{}).
		Select("TO_CHAR(TO_TIMESTAMP(created_at), 'YYYY/MM/DD') AS date, COUNT(1) AS count").
		Where("created_at >= ? AND created_at <= ?", params.DateFrom, params.DateTo).
		Group("date").
		Order("date DESC").
		Find(&data).Error
	if err != nil {
		return
	}
	result = &models.DataAnalyticNewCatalogProduct{
		Count: int64(len(data)),
	}
	for _, v := range data {
		result.Charts = append(result.Charts, models.DataAnalyticChart{
			Date:  v.Date,
			Value: float64(v.Count),
		})
	}
	return
}

type DataAnalyticInquiriesParam struct {
	models.PaginationParams
	models.JwtClaimsInfo
	DateFrom int64 `json:"date_from" query:"date_from" form:"date_from" param:"date_from"`
	DateTo   int64 `json:"date_to" query:"date_to" form:"date_to" param:"date_to" `
}

func (r *DataAnalyticRepo) PaginateInquiries(params DataAnalyticInquiriesParam) (result *models.DataAnalyticInquiry, err error) {
	result = &models.DataAnalyticInquiry{}
	err = r.db.Transaction(func(tx *gorm.DB) (err error) {
		qInquiry := r.db.Model(&models.Inquiry{}).
			Select("count(1)").
			Where("created_at >= ? AND created_at <= ?", params.DateFrom, params.DateTo).
			Session(&gorm.Session{})

		qInquiryAudit := r.db.Model(&models.InquiryAudit{}).
			Select("count(1)").
			Where("created_at >= ? AND created_at <= ?", params.DateFrom, params.DateTo).
			Session(&gorm.Session{})

		err = qInquiry.Find(&result.Num).Error
		err = qInquiry.Where("assignee_ids is not null").Find(&result.AssignedNum).Error
		err = qInquiry.Where("assignee_ids is null").Find(&result.UnAssignedNum).Error
		// 172800 -> 48H
		err = qInquiry.Where("status = ? AND (? - created_at > 172800)", enums.InquiryStatusNew, time.Now().Unix()).Find(&result.QuoteOverdue).Error
		err = qInquiry.Where("status = ? AND (? - created_at <= 172800)", enums.InquiryStatusNew, time.Now().Unix()).Find(&result.QuotePotentialLate).Error
		err = qInquiryAudit.Where("action_type = ?", enums.AuditActionTypeInquiryCreated).Find(&result.RFQSubmitted).Error
		err = qInquiryAudit.Where("action_type = ?", enums.AuditActionTypeInquiryAdminSendBuyerQuotation).Find(&result.SendQuotation).Error
		err = qInquiryAudit.Where("action_type = ?", enums.AuditActionTypeInquiryBuyerApproveQuotation).Find(&result.QuotationApproved).Error
		err = qInquiryAudit.Where("action_type = ?", enums.AuditActionTypeInquirySamplePoCreated).Find(&result.WaitingForPayment).Error
		err = qInquiryAudit.Where("action_type = ?", enums.AuditActionTypeInquiryAdminMarkAsPaid).Find(&result.PaymentConfirmed).Error
		return
	})

	return
}

type DataAnalyticPurchaseOrdersParam struct {
	models.PaginationParams
	models.JwtClaimsInfo
	DateFrom int64 `json:"date_from" query:"date_from" form:"date_from" param:"date_from"`
	DateTo   int64 `json:"date_to" query:"date_to" form:"date_to" param:"date_to" `
}

func (r *DataAnalyticRepo) PaginatePurchaseOrders(params DataAnalyticPurchaseOrdersParam) (result *models.DataAnalyticPO, err error) {
	result = &models.DataAnalyticPO{}
	err = r.db.Transaction(func(tx *gorm.DB) (err error) {
		qPO := r.db.Model(&models.PurchaseOrder{}).
			Select("count(1)").
			Where("created_at >= ? AND created_at <= ?", params.DateFrom, params.DateTo).
			Session(&gorm.Session{})

		err = qPO.Find(&result.Num).Error
		err = qPO.Where("assignee_ids is not null").Find(&result.AssignedNum).Error
		err = qPO.Where("assignee_ids is null").Find(&result.UnAssignedNum).Error
		err = qPO.Where("tracking_status = ? AND (? - created_at > lead_time*86400)", enums.PoTrackingStatusNew, time.Now().Unix()).Find(&result.QuoteOverdue).Error
		err = qPO.Where("tracking_status = ? AND (? - created_at <= lead_time*86400)", enums.PoTrackingStatusNew, time.Now().Unix()).Find(&result.QuotePotentialLate).Error
		err = qPO.Where("tracking_status = ?", enums.PoTrackingStatusNew).Find(&result.SampleOrder).Error
		err = qPO.Where("tracking_status = ? OR tracking_status = ?", enums.PoTrackingStatusDesignRejected, enums.PoTrackingStatusWaitingForApproved).Find(&result.Design).Error
		err = qPO.Where("tracking_status = ? OR tracking_status = ?", enums.PoTrackingStatusDesignApproved, enums.PoTrackingStatusRawMaterial).Find(&result.RawMaterial).Error
		err = qPO.Where("tracking_status = ?", enums.PoTrackingStatusMaking).Find(&result.Making).Error
		err = qPO.Where("tracking_status = ?", enums.PoTrackingActionMarkSubmit).Find(&result.Submit).Error
		err = qPO.Where("tracking_status = ?", enums.PoTrackingActionMarkDelivering).Find(&result.Delivery).Error
		err = qPO.Where("tracking_status = ?", enums.PoTrackingActionConfirmDelivered).Find(&result.Approval).Error
		return
	})
	return
}

type DataAnalyticBulkPurchaseOrdersParam struct {
	models.PaginationParams
	models.JwtClaimsInfo
	DateFrom int64 `json:"date_from" query:"date_from" form:"date_from" param:"date_from"`
	DateTo   int64 `json:"date_to" query:"date_to" form:"date_to" param:"date_to" `
}

func (r *DataAnalyticRepo) PaginateBulkPurchaseOrders(params DataAnalyticBulkPurchaseOrdersParam) (result *models.DataAnalyticBulkPO, err error) {
	result = &models.DataAnalyticBulkPO{}
	err = r.db.Transaction(func(tx *gorm.DB) (err error) {
		qBPO := r.db.Model(&models.BulkPurchaseOrder{}).
			Select("count(1)").
			Where("created_at >= ? AND created_at <= ?", params.DateFrom, params.DateTo).
			Session(&gorm.Session{})

		err = qBPO.Find(&result.Num).Error
		err = qBPO.Where("assignee_ids is not null").Find(&result.AssignedNum).Error
		err = qBPO.Where("assignee_ids is null").Find(&result.UnAssignedNum).Error
		err = qBPO.Where("tracking_status = ? AND (? - created_at > lead_time*86400)", enums.BulkPoTrackingStatusNew, time.Now().Unix()).Find(&result.QuoteOverdue).Error
		err = qBPO.Where("tracking_status = ? AND (? - created_at > lead_time*86400)", enums.BulkPoTrackingStatusNew, time.Now().Unix()).Find(&result.QuotePotentialLate).Error
		err = qBPO.Where("tracking_status = ?", enums.BulkPoTrackingStatusNew).Find(&result.New).Error
		err = qBPO.Where("tracking_status IN ?", []enums.BulkPoTrackingStatus{
			enums.BulkPoTrackingStatusWaitingForSubmitOrder,
			enums.BulkPoTrackingStatusWaitingForQuotation,
		}).Find(&result.Review).Error
		err = qBPO.Where("tracking_status IN ?", []enums.BulkPoTrackingStatus{
			enums.BulkPoTrackingStatusFirstPayment,
			enums.BulkPoTrackingStatusFirstPaymentConfirm,
		}).Find(&result.FirstPayment).Error
		err = qBPO.Where("tracking_status IN ?", []enums.BulkPoTrackingStatus{
			enums.BulkPoTrackingStatusFirstPaymentConfirmed,
			enums.BulkPoTrackingStatusRawMaterial,
			enums.BulkPoTrackingStatusPps,
			enums.BulkPoTrackingStatusProduction,
			enums.BulkPoTrackingStatusQc,
		}).Find(&result.Making).Error
		err = qBPO.Where("tracking_status = ?", enums.BulkPoTrackingStatusSubmit).Find(&result.Submit).Error
		err = qBPO.Where("tracking_status IN ?", []enums.BulkPoTrackingStatus{
			enums.BulkPoTrackingStatusFinalPayment,
			enums.BulkPoTrackingStatusFinalPaymentConfirm,
			enums.BulkPoTrackingStatusFinalPaymentConfirmed,
		}).Find(&result.FinalPayment).Error
		err = qBPO.Where("tracking_status = ?", enums.BulkPoTrackingStatusDelivering).Find(&result.Delivery).Error
		return
	})
	return
}

type DataAnalyticOpsBizPerformance struct {
	models.PaginationParams
	models.JwtClaimsInfo
	DateFrom int64 `json:"date_from" query:"date_from" form:"date_from" param:"date_from"`
	DateTo   int64 `json:"date_to" query:"date_to" form:"date_to" param:"date_to" `
}

func (r *DataAnalyticRepo) GetOpsBizPerformance(params DataAnalyticOpsBizPerformance) (result *models.DataAnalyticPerformance, err error) {
	result = &models.DataAnalyticPerformance{}
	spacing := params.DateTo - params.DateFrom
	var wg sync.WaitGroup
	wg.Add(6)
	go func() {
		defer wg.Done()
		// Inquiry
		q := r.db.Model(&models.Inquiry{}).Select("count(1)").
			Session(&gorm.Session{})
		qCreatedAt := r.db.Model(&models.Inquiry{}).Select("count(1)").
			Where("created_at >= ? AND created_at <= ?", params.DateFrom, params.DateTo).
			Session(&gorm.Session{})

		err = qCreatedAt.Where("quotation_at - created_at < 172800").Find(&result.InquiryQuoteInTime).Error
		var iqInTimeBefore, iqAll int64
		err = qCreatedAt.Find(&iqAll).Error

		err = q.Where("created_at >= ? AND created_at <= ?", params.DateFrom-spacing, params.DateTo-spacing).
			Where("quotation_at - created_at < 172800").
			Find(&iqInTimeBefore).Error
		if iqInTimeBefore != 0 {
			result.InquiryQuoteInTimeDiffPercentage = float64(result.InquiryQuoteInTime) / float64(iqInTimeBefore)
		}
		if iqAll != 0 {
			result.InquiryQuoteInTimeTotalPercentage = float64(result.InquiryQuoteInTime) * 100 / float64(iqAll)
		}
	}()
	go func() {
		defer wg.Done()
		// PO
		qPO := r.db.Model(&models.PurchaseOrder{}).Select("count(1)").
			Session(&gorm.Session{})

		qPOCreatedAt := r.db.Model(&models.PurchaseOrder{}).Select("count(1)").
			Where("created_at >= ? AND created_at <= ?", params.DateFrom, params.DateTo).
			Session(&gorm.Session{})

		err = qPOCreatedAt.Where("delivery_started_at - created_at < lead_time*86400").Find(&result.POInLeadTime).Error
		var poInTimeBefore, poAll int64
		err = qPOCreatedAt.Find(&poAll).Error
		err = qPO.Where("created_at >= ? AND created_at <= ?", params.DateFrom-spacing, params.DateTo-spacing).
			Where("delivery_started_at - created_at < lead_time*86400").
			Find(&poInTimeBefore).Error
		if poInTimeBefore != 0 {
			result.POInLeadTimeDiffPercentage = float64(result.POInLeadTime) / float64(poInTimeBefore)
		}
		if poAll != 0 {
			result.POInLeadTimeTotalPercentage = float64(result.POInLeadTime) * 100 / float64(poAll)
		}
	}()
	go func() {
		defer wg.Done()
		// Bulk PO
		qBPO := r.db.Model(&models.BulkPurchaseOrder{}).Select("count(1)").
			Session(&gorm.Session{})

		qBPOCreatedAt := r.db.Model(&models.BulkPurchaseOrder{}).Select("count(1)").
			Where("created_at >= ? AND created_at <= ?", params.DateFrom, params.DateTo).
			Session(&gorm.Session{})

		err = qBPOCreatedAt.
			Where("delivery_started_at - created_at < (CAST(admin_quotations -> 0 -> 'lead_time' AS INTEGER)) * 86400").
			Find(&result.BulkPOInLeadTime).Error
		var bulkPoInTimeBefore, bulkPoAll int64
		err = qBPOCreatedAt.Find(&bulkPoAll).Error
		err = qBPO.Where("created_at >= ? AND created_at <= ?", params.DateFrom-spacing, params.DateTo-spacing).
			Where("delivery_started_at - created_at < (CAST(admin_quotations -> 0 -> 'lead_time' AS INTEGER)) * 86400").
			Find(&bulkPoInTimeBefore).Error
		if bulkPoInTimeBefore != 0 {
			result.BulkPOInLeadTimeDiffPercentage = float64(result.BulkPOInLeadTime) / float64(bulkPoInTimeBefore)
		}
		if bulkPoAll != 0 {
			result.BulkPOInLeadTimeTotalPercentage = float64(result.BulkPOInLeadTime) * 100 / float64(bulkPoAll)
		}
	}()
	go func() {
		defer wg.Done()
		// Inquiry Approved
		var iqApproveBefore, iqApproveAll int64
		err = r.db.Model(&models.Inquiry{}).
			Joins("JOIN (SELECT inquiry_id, created_at, row_number() OVER (PARTITION BY inquiry_id ORDER BY created_at) as rn FROM inquiry_audits WHERE action_type = ?) AS tmp ON id = tmp.inquiry_id AND rn = 1", enums.AuditActionTypeInquiryBuyerApproveQuotation).
			Where("inquiries.created_at > ? AND inquiries.created_at < ?", params.DateFrom, params.DateTo).
			Count(&result.InquiryApproved).Error
		err = r.db.Model(&models.Inquiry{}).
			Where("inquiries.created_at >= ? AND inquiries.created_at <= ?", params.DateFrom, params.DateTo).
			Count(&iqApproveAll).Error
		err = r.db.Model(&models.Inquiry{}).
			Joins("JOIN (SELECT inquiry_id, created_at, row_number() OVER (PARTITION BY inquiry_id ORDER BY created_at) as rn FROM inquiry_audits WHERE action_type = ?) AS tmp ON id = tmp.inquiry_id AND rn = 1", enums.AuditActionTypeInquiryBuyerApproveQuotation).
			Where("inquiries.created_at > ? AND inquiries.created_at < ?", params.DateFrom-spacing, params.DateTo-spacing).
			Count(&iqApproveBefore).Error

		if iqApproveBefore != 0 {
			result.InquiryApprovedDiffPercentage = float64(result.InquiryApproved) / float64(iqApproveBefore)
		}
		if iqApproveAll != 0 {
			result.InquiryApprovedTotalPercentage = float64(result.InquiryApproved) * 100 / float64(iqApproveAll)
		}
	}()
	go func() {
		defer wg.Done()
		// Sample Payment
		qPOPayment := r.db.Model(&models.PurchaseOrder{}).Select("count(1)").
			Session(&gorm.Session{})

		qPOPaymentCreatedAt := r.db.Model(&models.PurchaseOrder{}).Select("count(1)").
			Where("created_at >= ? AND created_at <= ?", params.DateFrom, params.DateTo).
			Session(&gorm.Session{})

		err = qPOPaymentCreatedAt.Where("status = ?", enums.PurchaseOrderStatusPaid).Find(&result.POPaid).Error
		var qPoPaymentBefore, qPoPaymentAll int64
		err = qPOPaymentCreatedAt.Find(&qPoPaymentAll).Error
		err = qPOPayment.Where("created_at >= ? AND created_at <= ?", params.DateFrom-spacing, params.DateTo-spacing).
			Where("status = ?", enums.PurchaseOrderStatusPaid).
			Find(&qPoPaymentBefore).Error
		if qPoPaymentBefore != 0 {
			result.POPaidDiffPercentage = float64(result.POPaid) / float64(qPoPaymentBefore)
		}
		if qPoPaymentAll != 0 {
			result.POPaidTotalPercentage = float64(result.POPaid) * 100 / float64(qPoPaymentAll)
		}
	}()
	go func() {
		defer wg.Done()
		// Bulk PO First Payment
		qBPOFirstPayment := r.db.Model(&models.BulkPurchaseOrder{}).Select("count(1)").
			Session(&gorm.Session{})

		qBPOFirstPaymentCreatedAt := r.db.Model(&models.BulkPurchaseOrder{}).Select("count(1)").
			Where("created_at >= ? AND created_at <= ?", params.DateFrom, params.DateTo).
			Session(&gorm.Session{})

		err = qBPOFirstPaymentCreatedAt.Where("tracking_status LIKE ?", "%payment%").Find(&result.BulkPOPaid).Error
		var qBPOFirstPaymentBefore, qBPOFirstPaymentAll int64
		err = qBPOFirstPaymentCreatedAt.Find(&qBPOFirstPaymentAll).Error
		err = qBPOFirstPayment.Where("created_at >= ? AND created_at <= ?", params.DateFrom-spacing, params.DateTo-spacing).
			Where("tracking_status LIKE ?", "%payment%").
			Find(&qBPOFirstPaymentBefore).Error
		if qBPOFirstPaymentBefore != 0 {
			result.BulkPOPaidDiffPercentage = float64(result.BulkPOPaid) / float64(qBPOFirstPaymentBefore)
		}
		if qBPOFirstPaymentAll != 0 {
			result.BulkPOPaidTotalPercentage = float64(result.BulkPOPaid) * 100 / float64(qBPOFirstPaymentAll)
		}
	}()
	wg.Wait()
	return
}

type BuyerDataAnalyticRFQParams struct {
	models.PaginationParams
	models.JwtClaimsInfo
}

func (r *DataAnalyticRepo) BuyerDataAnalyticRFQ(params BuyerDataAnalyticRFQParams) (result *models.BuyerDataAnalyticRFQ, err error) {
	userID := params.GetUserID()
	result = &models.BuyerDataAnalyticRFQ{}
	r.db.Model(&models.Inquiry{}).Select("count(1)").
		Where("user_id = ?", userID).
		Find(&result.Total)
	r.db.Model(&models.Inquiry{}).Select("count(1)").
		Where("user_id = ?", userID).
		Where("buyer_quotation_status = ?", enums.InquiryBuyerStatusApproved).
		Find(&result.Approved)
	r.db.Model(&models.Inquiry{}).Select("count(1)").
		Where("user_id = ?", userID).
		Where("buyer_quotation_status = ?", enums.InquiryBuyerStatusRejected).
		Find(&result.Rejected)
	r.db.Model(&models.PurchaseOrder{}).Select("count(1)").Where("user_id = ?", userID).Find(&result.PurchaseOrder)
	r.db.Model(&models.BulkPurchaseOrder{}).Select("count(1)").Where("user_id = ?", userID).Find(&result.BulkPurchaseOrder)

	return
}

type BuyerDataAnalyticPendingTasksParams struct {
	models.PaginationParams
	models.JwtClaimsInfo
}

func (r *DataAnalyticRepo) BuyerDataAnalyticPendingTasks(params BuyerDataAnalyticPendingTasksParams) (result *query.Pagination, err error) {
	userID := params.GetUserID()
	result = &query.Pagination{}
	if params.Page <= 0 {
		params.Page = 1
	}
	if params.Limit <= 0 {
		params.Limit = 10
	}

	type Alias struct {
		ID          string              `json:"id,omitempty"`
		Title       string              `json:"title,omitempty"`
		Quantity    *int64              `json:"quantity,omitempty"`
		Attachments *models.Attachments `json:"attachments,omitempty"`
		ReferenceID string              `json:"reference_id,omitempty"`
		Status      string              `json:"status,omitempty"`
	}
	var total int64
	subQueryTitle := r.db.Model(&models.Inquiry{}).Select("title").Where("id = inquiry_id")
	subQueryQuantity := r.db.Model(&models.Inquiry{}).Select("quantity").Where("id = inquiry_id")
	subQueryAttachments := r.db.Model(&models.Inquiry{}).Select("attachments").Where("id = inquiry_id")
	err = r.db.Raw("select count(1) from (? UNION ? UNION ?) as ipb",
		r.db.Select("id", "title", "quantity", "attachments", "reference_id", "status", "updated_at").Model(&models.Inquiry{}).Where("user_id = ?", userID),
		r.db.Select("id,(?),(?),(?),reference_id,status,updated_at", subQueryTitle, subQueryQuantity, subQueryAttachments).Model(&models.PurchaseOrder{}).Where("user_id = ?", userID),
		r.db.Select("id,(?),(?),(?),reference_id,status,updated_at", subQueryTitle, subQueryQuantity, subQueryAttachments).Model(&models.BulkPurchaseOrder{}).Where("user_id = ?", userID),
	).Find(&total).Error
	if err != nil {
		return
	}

	var tasks []Alias
	err = r.db.Raw("? UNION ? UNION ? ORDER BY updated_at DESC LIMIT ? OFFSET ?",
		r.db.Select("id", "title", "quantity", "attachments", "reference_id", "status", "updated_at").Model(&models.Inquiry{}).Where("user_id = ?", userID),
		r.db.Select("id,(?),(?),(?),reference_id,status,updated_at", subQueryTitle, subQueryQuantity, subQueryAttachments).Model(&models.PurchaseOrder{}).Where("user_id = ?", userID),
		r.db.Select("id,(?),(?),(?),reference_id,status,updated_at", subQueryTitle, subQueryQuantity, subQueryAttachments).Model(&models.BulkPurchaseOrder{}).Where("user_id = ?", userID),
		params.Limit,
		(params.Page-1)*params.Limit,
	).Find(&tasks).Error
	if err != nil {
		return
	}
	var records []models.BuyerDataAnalyticPendingTask
	for _, v := range tasks {
		records = append(records, models.BuyerDataAnalyticPendingTask{
			ID:          v.ID,
			ProductName: v.Title,
			Quantity:    v.Quantity,
			Attachments: v.Attachments,
			OrderID:     v.ReferenceID,
			Status:      v.Status,
		})
	}
	result = &query.Pagination{
		HasNext:            total > int64(params.Page*params.Limit),
		HasPrev:            false,
		PerPage:            params.Limit,
		NextPage:           params.Page + 1,
		Page:               params.Page,
		PrevPage:           params.Page - 1,
		Offset:             (params.Page - 1) * params.Limit,
		Records:            records,
		TotalRecord:        int(total),
		TotalPage:          int(total / int64(params.Limit)),
		TotalCurrentRecord: 0,
	}
	return
}

type BuyerDataAnalyticPendingPaymentsParams struct {
	models.PaginationParams
	models.JwtClaimsInfo
}

func (r *DataAnalyticRepo) BuyerDataAnalyticPendingPayments(params BuyerDataAnalyticPendingPaymentsParams) (result *query.Pagination, err error) {
	result = &query.Pagination{}
	if params.Page <= 0 {
		params.Page = 1
	}
	if params.Limit <= 0 {
		params.Limit = 10
	}
	type transactionAlias struct {
		UserID              string                 `json:"user_id"`
		Title               string                 `json:"title"`
		Quantity            int64                  `json:"quantity"`
		Attachments         *models.Attachments    `json:"iq_attachments,omitempty"`
		ReferenceID         string                 `json:"reference_id"`
		TotalAmount         float64                `json:"total_amount"`
		PaidAmount          float64                `json:"paid_amount"`
		Milestone           enums.PaymentMilestone `json:"milestone"`
		Currency            enums.Currency         `json:"currency"`
		MetaInquiryID       string                 `json:"meta_inquiry_id"`
		IQReferenceID       string                 `json:"iq_reference_id"`
		BulkReferenceID     string                 `json:"bulk_reference_id"`
		POReferenceID       string                 `json:"po_reference_id"`
		PurchaseOrderID     string                 `json:"purchase_order_id"`
		BulkPurchaseOrderID string                 `json:"bulk_purchase_order_id"`
	}
	userID := params.GetUserID()

	var total int64
	err = r.db.Model(&models.PaymentTransaction{}).
		Joins("JOIN inquiries ON (metadata ->> 'inquiry_id' != '' AND inquiries.id = metadata ->> 'inquiry_id')").
		Where("payment_transactions.user_id = ?", userID).
		Where("milestone != ?", enums.PaymentMilestoneFinalPayment).
		Where("total_amount is not null").
		Where("inquiry_id is not null OR purchase_order_id is not null OR bulk_purchase_order_id is not null").
		Count(&total).Error
	if err != nil {
		return
	}

	var transactions []transactionAlias
	err = r.db.Model(&models.PaymentTransaction{}).
		Joins("JOIN inquiries ON (metadata ->> 'inquiry_id' != '' AND inquiries.id = metadata ->> 'inquiry_id')").
		Select("*,inquiries.attachments as iq_attachments,title,quantity,metadata ->> 'inquiry_id' as meta_inquiry_id, metadata ->> 'inquiry_reference_id' as iq_reference_id, metadata ->> 'bulk_purchase_order_reference_id' as bulk_reference_id, metadata ->> 'purchase_order_reference_id' as po_reference_id").
		Where("payment_transactions.user_id = ?", userID).
		Where("milestone != ?", enums.PaymentMilestoneFinalPayment).
		Where("total_amount is not null").
		Where("inquiry_id is not null or purchase_order_id is not null or bulk_purchase_order_id is not null").
		Limit(params.Limit).
		Offset((params.Page - 1) * params.Limit).
		Find(&transactions).Error
	if err != nil {
		return
	}
	var records []models.BuyerDataAnalyticPendingPayment
	for _, trx := range transactions {
		pending := models.BuyerDataAnalyticPendingPayment{
			Amount:      trx.TotalAmount - trx.PaidAmount,
			Milestone:   trx.Milestone,
			Currency:    trx.Currency,
			ProductName: trx.Title,
			Quantity:    trx.Quantity,
			Attachments: trx.Attachments,
		}
		if trx.BulkPurchaseOrderID != "" {
			pending.ID = trx.BulkPurchaseOrderID
			pending.OrderID = trx.BulkReferenceID
		} else if trx.PurchaseOrderID != "" {
			pending.ID = trx.PurchaseOrderID
			pending.OrderID = trx.POReferenceID
		} else if trx.MetaInquiryID != "" {
			pending.ID = trx.MetaInquiryID
			pending.OrderID = trx.IQReferenceID
		}
		records = append(records, pending)
	}

	result = &query.Pagination{
		HasNext:            total > int64(params.Page*params.Limit),
		HasPrev:            false,
		PerPage:            params.Limit,
		NextPage:           params.Page + 1,
		Page:               params.Page,
		PrevPage:           params.Page - 1,
		Offset:             (params.Page - 1) * params.Limit,
		Records:            records,
		TotalRecord:        int(total),
		TotalPage:          int(total / int64(params.Limit)),
		TotalCurrentRecord: 0,
	}
	return
}

func (r *DataAnalyticRepo) BuyerDataAnalyticPendingPaymentsV2(params BuyerDataAnalyticPendingPaymentsParams) (result *query.Pagination, err error) {
	var builder = queryfunc.NewBuyerDashboardPaymentTransactionBuilder(queryfunc.PaymentTransactionBuilderOptions{
		QueryBuilderOptions: queryfunc.QueryBuilderOptions{
			Role: params.GetRole(),
		},
	})
	result = query.New(r.db, builder).
		WhereFunc(func(builder *query.Builder) {
			builder.Where("user_id = ?", params.GetUserID())
			builder.Where("status != ?", enums.PaymentStatusPaid)
			builder.Where("milestone IN ?", []enums.PaymentMilestone{enums.PaymentMilestoneFirstPayment, enums.PaymentMilestoneFinalPayment})
		}).
		OrderBy("created_at DESC").
		Page(params.Page).
		Limit(params.Limit).
		PagingFunc()
	return
}

type BuyerDataAnalyticTotalStyleProduceParams struct {
	models.JwtClaimsInfo
	DateFrom int64 `json:"date_from" query:"date_from" form:"date_from" param:"date_from"`
	DateTo   int64 `json:"date_to" query:"date_to" form:"date_to" param:"date_to" `
}

func (r *DataAnalyticRepo) BuyerDataAnalyticTotalStyleProduce(params BuyerDataAnalyticTotalStyleProduceParams) (result *models.BuyerDataAnalyticTotalStyleProduce, err error) {
	userID := params.GetUserID()
	var pos []models.PurchaseOrder
	r.db.Model(&models.PurchaseOrder{}).Select("id", "created_at").
		Where("user_id = ?", userID).
		Where("created_at >= ? AND created_at <= ?", params.DateFrom, params.DateTo).
		Find(&pos)
	var bulkPos []models.BulkPurchaseOrder
	r.db.Model(&models.BulkPurchaseOrder{}).Select("id", "created_at").
		Where("user_id = ?", userID).
		Where("created_at >= ? AND created_at <= ?", params.DateFrom, params.DateTo).
		Find(&bulkPos)
	var deliveredBulkPO []models.BulkPurchaseOrder
	r.db.Model(&models.BulkPurchaseOrder{}).Select("id", "created_at").
		Where("user_id = ? and tracking_status = ?", userID, enums.BulkPoTrackingStatusDelivered).
		Where("created_at >= ? AND created_at <= ?", params.DateFrom, params.DateTo).
		Find(&deliveredBulkPO)

	m := make(map[string]*models.TotalStyleProduceChart)
	const layout = "02/01/2006"
	for _, v := range pos {
		date := time.Unix(v.CreatedAt, 0).Format(layout)
		if _, ok := m[date]; !ok {
			m[date] = &models.TotalStyleProduceChart{}
		}
		m[date].PurchaseOrder += 1
		m[date].Date = date
		m[date].Timestamp = v.CreatedAt
	}

	for _, v := range bulkPos {
		date := time.Unix(v.CreatedAt, 0).Format(layout)
		if _, ok := m[date]; !ok {
			m[date] = &models.TotalStyleProduceChart{}
		}
		m[date].BulkPurchaseOrder += 1
		m[date].Date = date
		m[date].Timestamp = v.CreatedAt
	}

	for _, v := range deliveredBulkPO {
		date := time.Unix(v.CreatedAt, 0).Format(layout)
		if _, ok := m[date]; !ok {
			m[date] = &models.TotalStyleProduceChart{}
		}
		m[date].Shipped += 1
		m[date].Date = date
		m[date].Timestamp = v.CreatedAt
	}

	var charts []models.TotalStyleProduceChart
	for _, v := range m {
		charts = append(charts, *v)
	}
	sort.Slice(charts, func(i, j int) bool {
		return charts[i].Date < charts[j].Date
	})

	result = &models.BuyerDataAnalyticTotalStyleProduce{
		TotalPurchaseOrder:      int64(len(pos)),
		TotalBulkPurchaserOrder: int64(len(bulkPos)),
		TotalShipped:            int64(len(deliveredBulkPO)),
		Charts:                  charts,
	}
	return
}

type GetBestAnalyticProductParams struct {
	models.PaginationParams
	models.JwtClaimsInfo
}

func (r *DataAnalyticRepo) GetBest(params GetBestAnalyticProductParams) (mainProduct models.AnalyticProduct, others []models.AnalyticProduct, err error) {
	var mainProductId = "cmeh5hbb2hjc375m021g"
	var otherProductIds = []string{"cnnsdejb2hjdqovgvjj0", "cnuierjb2hjb11aiam5g", "clpj3e7skmp0brpd0500", "cobrl5jb2hjb11ajbmo0"}

	var prepareIDs = []string{mainProductId}
	prepareIDs = append(prepareIDs, otherProductIds...)

	var products []models.AnalyticProduct
	if err = r.adb.Model(&models.AnalyticProduct{}).Where("id IN ?", prepareIDs).Find(&products).Error; err != nil {
		return
	}
	for _, p := range products {
		if p.ID == mainProductId {
			mainProduct = p
		} else {
			others = append(others, p)
		}
	}
	return
}
