package queryfunc

import (
	"fmt"
	"math"
	"sort"
	"strings"
	"text/template"
	"time"

	"github.com/engineeringinflow/inflow-backend/pkg/ai"
	"github.com/engineeringinflow/inflow-backend/pkg/models/enums"
	"github.com/jinzhu/copier"
	"github.com/lib/pq"
	"github.com/samber/lo"

	"github.com/engineeringinflow/inflow-backend/pkg/db"
	"github.com/engineeringinflow/inflow-backend/pkg/helper"
	"github.com/engineeringinflow/inflow-backend/pkg/models"
)

type AnalyticProductAlias struct {
	models.AnalyticProduct
}

type AnalyticProductBuilderOptions struct {
	QueryBuilderOptions
	IncludeProductClass bool
}

func NewAnalyticProductBuilder(options AnalyticProductBuilderOptions) *Builder {
	var rawSQL = `
	SELECT /* {{Description}} */ p.*
	FROM products p
	LEFT JOIN product_classes pc ON p.id = pc.product_id 
	`
	var countSQL = `
	SELECT /* {{Description}} */ 1
	FROM products p
	LEFT JOIN product_classes pc ON p.id = pc.product_id
	`

	builder := NewBuilder(rawSQL, countSQL).
		WithOptions(options, template.FuncMap{
			"Description": func() string {
				return helper.JoinNonEmptyStrings(
					"-",
					GetCaller(),
					options.Role.DisplayName(),
				)
			},
		}).
		WithPaginationFunc(func(db, rawSQL *db.DB) (interface{}, error) {
			var records = make([]*models.AnalyticProduct, rawSQL.RowsAffected)

			rows, err := rawSQL.Rows()
			if err != nil {
				return nil, err
			}
			defer rows.Close()
			var productIDs []string
			var productURls []string
			for rows.Next() {
				var alias AnalyticProductAlias
				err = db.ScanRows(rows, &alias)
				if err != nil {
					db.CustomLogger.Errorf("Scan rows error", err)
					continue
				}
				alias.AnalyticProduct.Metadata = nil
				if alias.AnalyticProduct.ID != "" {
					productIDs = append(productIDs, alias.AnalyticProduct.ID)
					productURls = append(productURls, alias.AnalyticProduct.URL)
					alias.AnalyticProduct.Images = helper.FixProductImages(alias.AnalyticProduct.Images)
					records = append(records, &alias.AnalyticProduct)
				}
			}

			if options.IncludeProductClass && len(productIDs) > 0 {
				var classes []*models.ProductClass
				db.Where("product_id IN ?", productIDs).Find(&classes)
				for _, record := range records {
					for _, class := range classes {
						if record.ID == class.ProductID {
							record.ProductClasses = append(record.ProductClasses, class)
						}
					}
				}
			}

			return &records, nil
		})
	return builder
}

func NewAnalyticProductGroupBuilder(options AnalyticProductBuilderOptions) *Builder {
	// product_groups is material view
	var rawSQL = `
	SELECT /* {{Description}} */ p.*, random() as ordering
	FROM product_groups p
	LEFT JOIN product_classes pc ON p.id = pc.product_id 
	`
	var countSQL = `
	SELECT /* {{Description}} */ 1
	FROM product_groups p
	LEFT JOIN product_classes pc ON p.id = pc.product_id
	`

	builder := NewBuilder(rawSQL, countSQL).
		WithOptions(options, template.FuncMap{
			"Description": func() string {
				return helper.JoinNonEmptyStrings(
					"-",
					GetCaller(),
					options.Role.DisplayName(),
				)
			},
		}).
		WithPaginationFunc(func(db, rawSQL *db.DB) (interface{}, error) {
			var records = make([]*models.AnalyticProduct, rawSQL.RowsAffected)

			rows, err := rawSQL.Rows()
			if err != nil {
				return nil, err
			}
			defer rows.Close()
			var productIDs []string
			for rows.Next() {
				var alias AnalyticProductAlias
				err = db.ScanRows(rows, &alias)
				if err != nil {
					db.CustomLogger.Errorf("Scan rows error", err)
					continue
				}

				if !options.Role.IsAdmin() {
					alias.AnalyticProduct.URL = ""
					alias.AnalyticProduct.Domain = ""
				}

				alias.AnalyticProduct.Metadata = nil
				if alias.AnalyticProduct.ID != "" {
					productIDs = append(productIDs, alias.AnalyticProduct.ID)
					alias.AnalyticProduct.Images = helper.FixProductImages(alias.AnalyticProduct.Images)
					records = append(records, &alias.AnalyticProduct)
				}

			}

			if options.IncludeProductClass && len(productIDs) > 0 {
				var classes []*models.ProductClass
				db.Where("product_id IN ?", productIDs).Find(&classes)
				for _, record := range records {
					for _, class := range classes {
						if record.ID == class.ProductID {
							record.ProductClasses = append(record.ProductClasses, class)
						}
					}
				}
			}

			return &records, nil
		})
	return builder
}

type RankDAProductAlias struct {
	models.AnalyticProduct
}

type RankDAProductBuilderOptions struct {
	models.PaginationParams
	QueryBuilderOptions
	OrderBy  string `json:"order_by" validate:"required,oneof=sold"`
	DateFrom int64  `json:"date_from"`
	DateTo   int64  `json:"date_to"`
}

func NewRankDAProductBuilder(options RankDAProductBuilderOptions) *Builder {
	var rawSQL = `
	SELECT /* {{Description}} */ {{fields}}, pc.score AS score
	FROM products p
	JOIN (
		SELECT product_id, MAX({{order_by}}) as score
		FROM product_changes
		GROUP BY product_id
	) pc ON p.id = pc.product_id
	`

	var countSQL = `
	SELECT /* {{Description}} */ 1
	FROM products p
	JOIN (
		SELECT product_id, MAX({{order_by}}) as score
		FROM product_changes
		GROUP BY product_id
	) pc ON p.id = pc.product_id
	`

	return NewBuilder(rawSQL, countSQL).
		WithOptions(options, template.FuncMap{
			"Description": func() string {
				return helper.JoinNonEmptyStrings(
					"-",
					GetCaller(),
					options.Role.DisplayName(),
				)
			},
			"fields": func() string {
				return "p.id, p.created_at, p.updated_at, p.url , p.country_code,p.domain ,p.name , p.description , p.category , p.sub_category , p.images , p.private_images , p.price , p.sold , p.stock , p.trending"
			},
			"order_by": func() string {
				return options.OrderBy
			},
		}).
		WithOrderBy("pc.score desc, p.id desc").
		WithPaginationFunc(func(db, rawSQL *db.DB) (interface{}, error) {
			var records = make([]*models.AnalyticProduct, rawSQL.RowsAffected)

			rows, err := rawSQL.Rows()
			if err != nil {
				return nil, err

			}
			defer rows.Close()

			for rows.Next() {
				var alias AnalyticProductAlias
				err = db.ScanRows(rows, &alias)
				if err != nil {
					db.CustomLogger.Errorf("Scan rows error", err)
					continue
				}
				alias.AnalyticProduct.Metadata = nil
				if options.OrderBy == "sold" {
					alias.AnalyticProduct.Sold = int64(alias.AnalyticProduct.Score)
				}
				alias.AnalyticProduct.Score = 0
				records = append(records, &alias.AnalyticProduct)
			}

			return &records, nil
		})
}

type RankDACategoryAlias struct {
	models.DataAnalyticCategory
}

type RankDACategoryBuilderOptions struct {
	models.PaginationParams
	QueryBuilderOptions
	Select   []string `json:"group_by,omitempty"`
	OrderBy  string   `json:"order_by,omitempty"`
	DateFrom int64    `json:"date_from"`
	DateTo   int64    `json:"date_to"`
}

func NewRankDACategoryBuilder(options RankDACategoryBuilderOptions) *Builder {
	var rawSQL = `
	SELECT /* {{Description}} */  {{group_by}}, sum(sold) as sold, sum(sold * price) as rev, json_agg(distinct domain) as domains
	FROM products p
	`

	var countSQL = `
	SELECT /* {{Description}} */  1
	FROM products p
	`

	return NewBuilder(rawSQL, countSQL).
		WithOptions(options, template.FuncMap{
			"Description": func() string {
				return helper.JoinNonEmptyStrings(
					"-",
					GetCaller(),
					options.Role.DisplayName(),
				)
			},
			"group_by": func() string {
				return strings.Join(options.Select, ",")
			},
		}).
		WithGroupBy(strings.Join(options.Select, ",")).
		WithOrderBy(fmt.Sprintf("%s DESC", options.OrderBy)).
		WithPaginationFunc(func(db, rawSQL *db.DB) (interface{}, error) {
			var records = make([]*models.DataAnalyticCategory, rawSQL.RowsAffected)

			rows, err := rawSQL.Rows()
			if err != nil {
				db.CustomLogger.Errorf("Scan rows error", err)
				return nil, err
			}
			defer rows.Close()

			rank := (options.Page-1)*options.Limit + 1
			for rows.Next() {
				var alias RankDACategoryAlias
				err = db.ScanRows(rows, &alias)
				if err != nil {
					db.CustomLogger.Errorf("Scan rows error", err)
					continue
				}
				alias.DataAnalyticCategory.Rank = rank
				records = append(records, &alias.DataAnalyticCategory)
				rank += 1
			}

			return &records, nil
		})
}

type ChartDAProductAlias struct {
	models.AnalyticProductChanges
}

type ChartDAProductBuilderOptions struct {
	models.PaginationParams
	QueryBuilderOptions
	PredictNext int
	PredictOn   enums.PredictChartOn
}

func NewChartDAProductBuilder(options ChartDAProductBuilderOptions) *Builder {
	var rawSQL = `
	SELECT /* {{Description}} */  pc.*, pc.price * pc.sold as rev
	FROM product_changes pc
	`

	var countSQL = `
	SELECT /* {{Description}} */  1
	FROM product_changes pc
	`
	if options.PredictNext > 0 {
		rawSQL = `
		SELECT max(id) as id, url, scrape_date, avg(price) as price, avg(sold) as sold, avg(stock) as stock,max(created_at) as created_at
		from product_changes pc
		`
		countSQL = `
		SELECT 1
		from product_changes pc
		`
	}

	return NewBuilder(rawSQL, countSQL).
		WithOptions(options, template.FuncMap{
			"Description": func() string {
				return helper.JoinNonEmptyStrings(
					"-",
					GetCaller(),
					options.Role.DisplayName(),
				)
			},
		}).
		WithOrderBy("pc.created_at").
		WithPaginationFunc(func(db, rawSQL *db.DB) (interface{}, error) {
			var records = make([]*models.AnalyticProductChanges, rawSQL.RowsAffected)

			rows, err := rawSQL.Rows()
			if err != nil {
				return nil, err

			}
			defer rows.Close()

			for rows.Next() {
				var alias ChartDAProductAlias
				err = db.ScanRows(rows, &alias)
				if err != nil {
					db.CustomLogger.Errorf("Scan rows error", err)
					continue
				}

				if !options.Role.IsAdmin() {
					alias.AnalyticProductChanges.URL = ""
					alias.AnalyticProductChanges.Domain = ""
				}
				records = append(records, &alias.AnalyticProductChanges)
			}

			if options.PredictNext > 0 && len(records) > 1 {
				var (
					prices  []float64
					sales   []float64
					stocks  []float64
					monthly int
					daily   int = 1
				)
				if options.PredictOn == enums.PredictOnWeek {
					daily = 7
				}
				if options.PredictOn == enums.PredictOnMonth {
					monthly = 1
				}
				sort.Slice(records, func(i, j int) bool {
					return records[i].ScrapeDate < records[j].ScrapeDate
				})
				for _, record := range records {
					prices = append(prices, record.Price)
					sales = append(sales, record.Sold)
					stocks = append(stocks, record.Stock)
				}
				predictedPrices := ai.PredictSeries(prices, options.PredictNext)
				predictedSales := ai.PredictSeries(sales, options.PredictNext)
				predictedStocks := ai.PredictSeries(stocks, options.PredictNext)
				last := time.Now()
				for i := 0; i < options.PredictNext; i++ {
					var changes models.AnalyticProductChanges
					if err = copier.Copy(&changes, &records[0]); err != nil {
						continue
					}
					changes.ScrapeDate = last.Format("2006-01-02")
					changes.IsPrediction = true
					last = last.AddDate(0, monthly, daily)
					if i < len(predictedPrices) {
						changes.Price = lo.Max([]float64{predictedPrices[i], 1})
					}
					if i < len(predictedSales) {
						changes.Sold = math.Round(lo.Max([]float64{predictedSales[i], 0}))
					}
					if i < len(predictedStocks) {
						changes.Stock = math.Round(lo.Max([]float64{predictedStocks[i], 0}))
					}
					records = append(records, &changes)
				}
			}

			return &records, nil
		})
}

type TopMoverDaProductAlias struct {
	models.AnalyticProduct
}

type TopTopMoverProductBuilderOptions struct {
	models.PaginationParams
	QueryBuilderOptions
}

func NewTopMoverProductBuilder(options TopTopMoverProductBuilderOptions) *Builder {
	var rawSQL = `
	SELECT /* {{Description}} */ {{fields}}, pc.average_sales_increase_percentage
	FROM products p
	JOIN (
		SELECT product_id, average_sales_increase_percentage
		FROM products_with_positive_monthly_growth
		WHERE month_filter = @month
				AND year_filter = @year
				ORDER BY average_sales_increase_percentage desc) pc ON p.id = pc.product_id
	`
	var countSQL = `
	SELECT /* {{Description}} */ 1
	FROM products p
	JOIN (
		SELECT product_id, average_sales_increase_percentage
		FROM products_with_positive_monthly_growth
		WHERE month_filter = @month
				AND year_filter = @year
				ORDER BY average_sales_increase_percentage desc) pc ON p.id = pc.product_id
	`

	return NewBuilder(rawSQL, countSQL).
		WithOptions(options, template.FuncMap{
			"Description": func() string {
				return helper.JoinNonEmptyStrings(
					"-",
					GetCaller(),
					options.Role.DisplayName(),
				)
			},
			"fields": func() string {
				return "p.id,p.created_at,p.updated_at,p.url,p.country_code,p.domain ,p.name , p.description , p.category , p.sub_category , p.images , p.private_images , p.price , p.sold , p.stock , p.trending"
			},
		}).
		WithPaginationFunc(func(db, rawSQL *db.DB) (interface{}, error) {
			var records = make([]*models.AnalyticProduct, rawSQL.RowsAffected)

			rows, err := rawSQL.Rows()
			if err != nil {
				return nil, err

			}
			defer rows.Close()

			for rows.Next() {
				var alias AnalyticProductAlias
				err = db.ScanRows(rows, &alias)
				if err != nil {
					db.CustomLogger.Errorf("Scan rows error", err)
					continue
				}
				records = append(records, &alias.AnalyticProduct)
			}

			return &records, nil
		})
}

type AnalyticProductTrendingAlias struct {
	ID                string         `json:"id,omitempty"`
	Trending          string         `json:"trending,omitempty"`
	Total             int64          `json:"total,omitempty"`
	Images            pq.StringArray `gorm:"type:varchar(200)[]" json:"images"`
	OverallGrowthRate float64        `json:"overall_growth_rate,omitempty"`
	TopSubCategories  pq.StringArray `gorm:"type:varchar(200)[]" json:"top_sub_categories"`
}

type AnalyticProductTrendingBuilderOptions struct {
	QueryBuilderOptions
}

func NewAnalyticProductTrendingBuilder(options AnalyticProductTrendingBuilderOptions) *Builder {
	var rawSQL = `
	WITH trending_products AS (
		SELECT id, trending, 
		overall_growth_rate, 
		sub_category, 
		images 
		FROM products 
		WHERE trending IS NOT NULL AND trending != '' AND overall_growth_rate IS NOT NULL
	),
	trendings AS (
		SELECT /* {{Description}} */ MAX(tp.id) as id, 
		tp.trending,
		MAX(tp.overall_growth_rate) as overall_growth_rate,
		COUNT(1) as total,
		MAX(tp.images) as images,
		ARRAY (
				SELECT sub_category
				FROM trending_products
				WHERE trending = tp.trending
				GROUP BY sub_category
				ORDER BY count(1) desc
				LIMIT 3
		) AS top_sub_categories
		FROM trending_products tp
		GROUP BY tp.trending
	)
	SELECT * 
	FROM trendings tp
	`

	var countSQL = `
	WITH trending_products AS (
		SELECT id, trending, 
		overall_growth_rate, 
		sub_category, 
		images 
		FROM products 
		WHERE trending IS NOT NULL AND trending != '' AND overall_growth_rate IS NOT NULL
	),
	trendings AS (
		SELECT /* {{Description}} */ MAX(tp.id) as id, 
		tp.trending,
		MAX(tp.overall_growth_rate) as overall_growth_rate,
		COUNT(1) as total,
		MAX(tp.images) as images,
		ARRAY (
				SELECT sub_category
				FROM trending_products
				WHERE trending = tp.trending
				GROUP BY sub_category
				ORDER BY count(1) desc
				LIMIT 3
		) AS top_sub_categories
		FROM trending_products tp
		GROUP BY tp.trending
	)
	SELECT 1 
	FROM trendings tp
	`

	builder := NewBuilder(rawSQL, countSQL).
		WithOptions(options, template.FuncMap{
			"Description": func() string {
				return helper.JoinNonEmptyStrings(
					"-",
					GetCaller(),
					options.Role.DisplayName(),
				)
			},
		}).
		WithPaginationFunc(func(db, rawSQL *db.DB) (interface{}, error) {
			var records = make([]AnalyticProductTrendingAlias, rawSQL.RowsAffected)

			rows, err := rawSQL.Rows()
			if err != nil {
				return nil, err
			}
			defer rows.Close()

			for rows.Next() {
				var alias AnalyticProductTrendingAlias
				err = db.ScanRows(rows, &alias)
				if err != nil {
					db.CustomLogger.Errorf("Scan rows error", err)
					continue
				}
				alias.Images = helper.FixProductImages(alias.Images)
				records = append(records, alias)
			}

			return &records, nil
		})
	return builder
}
