package queryfunc

import (
	"sync"
	"text/template"

	"github.com/engineeringinflow/inflow-backend/pkg/db"
	"github.com/engineeringinflow/inflow-backend/pkg/helper"
	"github.com/engineeringinflow/inflow-backend/pkg/models"
)

type CollectionProductGroupAlias struct {
	*models.CollectionProductGroup
}

type CollectionProductGroupBuilderOptions struct {
	QueryBuilderOptions

	IncludeProducts bool
}

func NewCollectionProductGroupBuilder(options CollectionProductGroupBuilderOptions) *Builder {
	var rawSQL = `
	SELECT /* {{Description}} */ g.*

	FROM collection_product_groups g
	`
	var countSQL = `
	SELECT /* {{Description}} */ 1

	FROM collection_product_groups g
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
		}).
		WithOrderBy("g.updated_at DESC").
		WithPaginationFunc(func(db, rawSQL *db.DB) (interface{}, error) {
			var records = make([]*models.CollectionProductGroup, rawSQL.RowsAffected)

			rows, err := rawSQL.Rows()
			if err != nil {
				return nil, err

			}
			defer rows.Close()

			var productIDs []string
			for rows.Next() {
				var alias CollectionProductGroupAlias
				err = db.ScanRows(rows, &alias)
				if err != nil {
					db.CustomLogger.Errorf("Scan rows error", err)
					continue
				}

				productIDs = append(productIDs, alias.ProductIDs...)

				records = append(records, alias.CollectionProductGroup)
			}

			var wg sync.WaitGroup

			if options.IncludeProducts {
				wg.Add(1)
				go func() {
					defer wg.Done()
					var products []*models.Product
					db.Find(&products, "id IN ?", productIDs)
					for _, product := range products {
						for _, record := range records {
							if helper.StringContains(record.ProductIDs, product.ID) {
								record.Products = append(record.Products, product)
							}
						}
					}
				}()

			}

			wg.Wait()
			return &records, nil
		})
}
