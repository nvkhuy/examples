package queryfunc

import (
	"github.com/engineeringinflow/inflow-backend/pkg/db"
	"github.com/engineeringinflow/inflow-backend/pkg/helper"
	"github.com/engineeringinflow/inflow-backend/pkg/models"
	"github.com/lib/pq"
	"github.com/samber/lo"
	"text/template"
)

type ProductTrendingGroupAlias struct {
	models.AnalyticProductTrending
	ProductTrendingIds pq.StringArray `gorm:"type:varchar(255)[]" json:"product_trending_ids"`
}

type ProductTrendingGroupBuilderOptions struct {
	QueryBuilderOptions
	DB *db.DB
}

func NewProductTrendingGroupBuilder(options ProductTrendingGroupBuilderOptions) *Builder {
	// product_trending_groups is Material View
	var rawSQL = `
	SELECT /* {{Description}} */ *
	FROM product_trending_groups ptg
	`

	var countSQL = `
	SELECT /* {{Description}} */ 1

	FROM product_trending_groups ptg
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
		WithPaginationFunc(func(adb, rawSQL *db.DB) (interface{}, error) {
			var records = make([]*models.AnalyticProductTrending, rawSQL.RowsAffected)

			rows, err := rawSQL.Rows()
			if err != nil {
				return nil, err
			}
			defer rows.Close()

			var (
				urls []string
			)
			for rows.Next() {
				var alias ProductTrendingGroupAlias
				err = adb.ScanRows(rows, &alias)
				if err != nil {
					adb.CustomLogger.Errorf("Scan rows error", err)
					continue
				}
				if alias.AnalyticProductTrending.URL != "" {
					urls = append(urls, alias.AnalyticProductTrending.URL)
				}
				if len(alias.ProductTrendingIds) > 0 {
					alias.AnalyticProductTrending.ProductTrendingIds = alias.ProductTrendingIds
				}
				records = append(records, &alias.AnalyticProductTrending)
			}

			if len(urls) > 0 {
				for _, record := range records {
					if len(record.ProductTrendingIds) > 0 {
						record.ID = record.ProductTrendingIds[0]
					}

					var trendings []models.Trending
					options.DB.Model(&models.Trending{}).
						Select("id", "name").
						Where("count_elements(product_trending_ids,?) > 0", record.ProductTrendingIds).
						Find(&trendings)
					trendings = lo.UniqBy(trendings, func(item models.Trending) string {
						return item.ID
					})
					if len(trendings) > 0 {
						record.Trendings = trendings
					}
				}
			}

			return records, nil
		})
	return builder
}
