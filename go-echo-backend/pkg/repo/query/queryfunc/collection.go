package queryfunc

import (
	"text/template"

	"gorm.io/plugin/dbresolver"

	"github.com/engineeringinflow/inflow-backend/pkg/db"
	"github.com/engineeringinflow/inflow-backend/pkg/helper"
	"github.com/engineeringinflow/inflow-backend/pkg/models"
	"github.com/engineeringinflow/inflow-backend/pkg/repo/query"
)

type CollectionAlias struct {
	*models.Collection
}

type CollectionBuilderOptions struct {
	QueryBuilderOptions
	IsConsistentRead bool
}

func NewCollectionBuilder(options CollectionBuilderOptions) *Builder {
	var rawSQL = `
	SELECT /* {{Description}} */ c.*

	FROM collections c
	`
	var countSQL = `
	SELECT /* {{Description}} */ 1

	FROM collections c
	`

	b := NewBuilder(rawSQL, countSQL).
		WithOptions(options, template.FuncMap{
			"Description": func() string {
				return helper.JoinNonEmptyStrings(
					"-",
					GetCaller(),
					options.Role.DisplayName(),
				)
			},
		}).
		WithOrderBy("c.order ASC").
		WithPaginationFunc(func(db, rawSQL *db.DB) (interface{}, error) {
			var records = make([]*models.Collection, rawSQL.RowsAffected)

			rows, err := rawSQL.Rows()
			if err != nil {
				db.CustomLogger.Errorf("Scan rows error", err)
				return nil, err
			}
			defer rows.Close()

			var collectionIDs []string
			for rows.Next() {
				var alias CollectionAlias
				err = db.ScanRows(rows, &alias)
				if err != nil {
					db.CustomLogger.Errorf("Scan rows error", err)
					continue
				}

				if alias.Collection.ID != "" && !helper.StringContains(collectionIDs, alias.Collection.ID) {
					collectionIDs = append(collectionIDs, alias.Collection.ID)
					var products []*models.Product
					_ = db.Where("id IN ?", []string(alias.Collection.ProductIDs)).Order("created_at desc").Find(&products)
					for _, product := range products {
						var variants []*models.Variant
						_ = db.Where("product_id = ?", product.ID).Find(&variants)
						product.Variants = variants
						var category models.Category
						_ = db.First(&category, "id = ?", product.CategoryID)
						product.Category = &category
					}
					alias.Collection.Products = products
				}
				records = append(records, alias.Collection)
			}

			if len(collectionIDs) > 0 {

				var collectionGroups []*models.CollectionProductGroup
				_ = query.New(db, NewCollectionProductGroupBuilder(CollectionProductGroupBuilderOptions{
					IncludeProducts: true,
				})).FindFunc(&collectionGroups)

				for _, collectionGroup := range collectionGroups {
					for _, record := range records {
						if collectionGroup.CollectionID == record.ID {
							record.ProductGroups = append(record.ProductGroups, collectionGroup)
						}
					}
				}

			}
			return &records, nil
		})
	if options.IsConsistentRead {
		b.WithClauses(dbresolver.Write)
	}
	return b
}
