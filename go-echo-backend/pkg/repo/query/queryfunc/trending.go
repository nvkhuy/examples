package queryfunc

import (
	"github.com/engineeringinflow/inflow-backend/pkg/db"
	"github.com/engineeringinflow/inflow-backend/pkg/helper"
	"github.com/engineeringinflow/inflow-backend/pkg/models"
	"text/template"
)

type TrendingAlias struct {
	*models.Trending
}

type TrendingBuilderOptions struct {
	QueryBuilderOptions
	Adb *db.DB
}

func NewTrendingBuilder(options TrendingBuilderOptions) *Builder {
	var rawSQL = `
	SELECT /* {{Description}} */ t.*
	FROM trendings t
	`

	var countSQL = `
	SELECT /* {{Description}} */ 1

	FROM trendings t
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
			var records = make([]*models.Trending, rawSQL.RowsAffected)

			rows, err := rawSQL.Rows()
			if err != nil {
				return nil, err

			}
			defer rows.Close()

			var (
				productTrendingIds []string
			)
			for rows.Next() {
				var alias TrendingAlias
				err = db.ScanRows(rows, &alias)
				if err != nil {
					db.CustomLogger.Errorf("Scan rows error", err)
					continue
				}
				if alias.Trending != nil && len(alias.Trending.ProductTrendingIDs) > 0 {
					productTrendingIds = append(productTrendingIds, alias.Trending.ProductTrendingIDs...)
				}
				records = append(records, alias.Trending)
			}
			if len(productTrendingIds) > 0 {
				var products []*models.AnalyticProductTrending
				options.Adb.Model(&models.AnalyticProductTrending{}).
					Where("id IN ?", productTrendingIds).Find(&products)
				if len(products) > 0 {
					m := make(map[string]*models.AnalyticProductTrending)
					for _, p := range products {
						m[p.ID] = p
					}
					for _, record := range records {
						for _, id := range record.ProductTrendingIDs {
							if p, ok := m[id]; ok && p != nil {
								record.ProductTrendings = append(record.ProductTrendings, *p)
							}
						}
					}
				}
			}

			return records, nil
		})
	return builder
}
