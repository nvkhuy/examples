package queryfunc

import (
	"github.com/engineeringinflow/inflow-backend/pkg/db"
	"github.com/engineeringinflow/inflow-backend/pkg/helper"
	"github.com/engineeringinflow/inflow-backend/pkg/models"
	"github.com/lib/pq"
	"text/template"
)

type ProductTrendingAlias struct {
	models.AnalyticProductTrending
}

type ProductTrendingBuilderOptions struct {
	QueryBuilderOptions
	DB *db.DB
}

func NewProductTrendingBuilder(options ProductTrendingBuilderOptions) *Builder {
	var rawSQL = `
	SELECT /* {{Description}} */ pt.*
	FROM product_trendings pt
	`

	var countSQL = `
	SELECT /* {{Description}} */ 1

	FROM product_trendings pt
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
				productTrendingIds []string
			)
			for rows.Next() {
				var alias ProductTrendingAlias
				err = adb.ScanRows(rows, &alias)
				if err != nil {
					adb.CustomLogger.Errorf("Scan rows error", err)
					continue
				}
				if alias.AnalyticProductTrending.ID != "" {
					productTrendingIds = append(productTrendingIds, alias.AnalyticProductTrending.ID)
				}
				records = append(records, &alias.AnalyticProductTrending)
			}

			if len(productTrendingIds) > 0 {
				for _, record := range records {
					var trendings []models.Trending
					options.DB.Model(&models.Trending{}).
						Select("id", "name").
						Where("count_elements(product_trending_ids,?) >= 1", pq.StringArray([]string{record.ID})).
						Find(&trendings)
					if len(trendings) > 0 {
						record.Trendings = trendings
					}
				}
			}

			return records, nil
		})
	return builder
}
