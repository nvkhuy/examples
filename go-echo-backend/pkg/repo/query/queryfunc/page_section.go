package queryfunc

import (
	"sync"
	"text/template"

	"github.com/engineeringinflow/inflow-backend/pkg/db"
	"github.com/engineeringinflow/inflow-backend/pkg/helper"
	"github.com/engineeringinflow/inflow-backend/pkg/models"
	"github.com/engineeringinflow/inflow-backend/pkg/repo/query"
	"github.com/samber/lo"
)

type PageSectionAlias struct {
	*models.PageSection
}

type PageSectionBuilderOptions struct {
	QueryBuilderOptions
}

func NewPageSectionBuilder(options PageSectionBuilderOptions) *Builder {
	var rawSQL = `
	SELECT /* {{Description}} */ p.*

	FROM page_sections p
	`
	var countSQL = `
	SELECT /* {{Description}} */ 1

	FROM page_sections p
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
		WithOrderBy("p.order ASC").
		WithPaginationFunc(func(db, rawSQL *db.DB) (interface{}, error) {
			var records = make([]*models.PageSection, rawSQL.RowsAffected)

			rows, err := rawSQL.Rows()
			if err != nil {
				return nil, err

			}
			defer rows.Close()

			var categoryIDs []string
			var productIDs []string
			var collectionIDs []string

			for rows.Next() {
				var alias PageSectionAlias
				err = db.ScanRows(rows, &alias)
				if err != nil {
					db.CustomLogger.Errorf("Scan rows error", err)
					continue
				}

				if len(alias.CategoryIds) > 0 {
					categoryIDs = append(categoryIDs, alias.CategoryIds...)
				}
				if len(alias.ProductIds) > 0 {
					productIDs = append(productIDs, alias.ProductIds...)
				}
				if len(alias.CollectionIds) > 0 {
					collectionIDs = append(collectionIDs, alias.CollectionIds...)
				}

				records = append(records, alias.PageSection)
			}

			var wg sync.WaitGroup

			if len(categoryIDs) > 0 {
				wg.Add(1)
				go func() {
					defer wg.Done()
					var categories []*models.Category
					if err := query.New(db, NewCategoryBuilder(CategoryBuilderOptions{})).
						WhereFunc(func(builder *query.Builder) {
							builder.Where("cate.id IN ?", categoryIDs)
						}).
						FindFunc(&categories); err != nil {
						return
					}

					for _, cate := range categories {
						for _, record := range records {
							if len(record.CategoryIds) > 0 {
								if lo.Contains(record.CategoryIds, cate.ID) {
									record.Categories = append(record.Categories, cate)
								}
							}
						}
					}
				}()
			}
			if len(productIDs) > 0 {
				wg.Add(1)
				go func() {
					defer wg.Done()
					var products []*models.Product
					if err := query.New(db, NewProductBuilder(ProductBuilderOptions{})).
						WhereFunc(func(builder *query.Builder) {
							builder.Where("p.id IN ?", productIDs)
						}).
						FindFunc(&products); err != nil {
						return
					}

					for _, product := range products {
						for _, record := range records {
							if len(record.ProductIds) > 0 {
								if lo.Contains(record.ProductIds, product.ID) {
									record.Products = append(record.Products, product)
								}
							}
						}
					}
				}()
			}
			if len(collectionIDs) > 0 {
				wg.Add(1)
				go func() {
					defer wg.Done()
					var collections []*models.Collection
					if err := query.New(db, NewCollectionBuilder(CollectionBuilderOptions{})).
						WhereFunc(func(builder *query.Builder) {
							builder.Where("c.id IN ?", collectionIDs)
						}).
						FindFunc(&collections); err != nil {
						return
					}

					for _, collection := range collections {
						for _, record := range records {
							if len(record.CollectionIds) > 0 {
								if lo.Contains(record.CollectionIds, collection.ID) {
									record.Collections = append(record.Collections, collection)
								}
							}
						}
					}
				}()
			}

			wg.Wait()

			return &records, nil
		})
}
