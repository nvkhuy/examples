package queryfunc

import (
	"sync"
	"text/template"

	"github.com/engineeringinflow/inflow-backend/pkg/db"
	"github.com/engineeringinflow/inflow-backend/pkg/helper"
	"github.com/engineeringinflow/inflow-backend/pkg/models"
)

type PurchaseOrderItemAlias struct {
	*models.PurchaseOrderItem
}

type PurchaseOrderItemBuilderOptions struct {
	QueryBuilderOptions

	IncludeProduct bool
	IncludeFabric  bool
	IncludeVariant bool
}

func NewPurchaseOrderItemBuilder(options PurchaseOrderItemBuilderOptions) *Builder {
	var rawSQL = `
	SELECT /* {{Description}} */ poi.*

	FROM purchase_order_items poi
	`
	var countSQL = `
	SELECT /* {{Description}} */ 1

	FROM purchase_order_items poi
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
		WithPaginationFunc(func(db, rawSQL *db.DB) (interface{}, error) {
			var records = make([]*models.PurchaseOrderItem, rawSQL.RowsAffected)

			rows, err := rawSQL.Rows()
			if err != nil {
				return nil, err

			}
			defer rows.Close()

			var productIDs []string
			var fabricIDs []string
			var variantIDs []string

			for rows.Next() {
				var alias PurchaseOrderItemAlias
				err = db.ScanRows(rows, &alias)
				if err != nil {
					db.CustomLogger.Errorf("Scan rows error", err)
					continue
				}

				if alias.ProductID != "" {
					productIDs = append(productIDs, alias.ProductID)
				}

				if alias.FabricID != "" {
					fabricIDs = append(fabricIDs, alias.FabricID)
				}

				if alias.VariantID != "" {
					variantIDs = append(variantIDs, alias.VariantID)
				}
				records = append(records, alias.PurchaseOrderItem)
			}

			var wg sync.WaitGroup

			if options.IncludeProduct && len(productIDs) > 0 {
				wg.Add(1)
				go func() {
					defer wg.Done()
					var products []*models.Product
					db.Find(&products, "id IN ?", productIDs)

					for _, product := range products {
						for _, record := range records {
							if record.ProductID == product.ID {
								record.Product = product
							}
						}
					}
				}()
			}

			if options.IncludeFabric && len(fabricIDs) > 0 {
				wg.Add(1)
				go func() {
					defer wg.Done()
					var fabrics []*models.Fabric
					db.Find(&fabrics, "id IN ?", fabricIDs)

					for _, fabric := range fabrics {
						for _, record := range records {
							if record.FabricID == fabric.ID {
								record.Fabric = fabric
							}
						}
					}
				}()
			}

			if options.IncludeVariant && len(variantIDs) > 0 {
				wg.Add(1)
				go func() {
					defer wg.Done()
					var variants []*models.Variant
					db.Find(&variants, "id IN ?", variantIDs)

					for _, variant := range variants {
						for _, record := range records {
							if record.VariantID == variant.ID {
								record.Variant = variant
							}
						}
					}
				}()
			}

			wg.Wait()

			return &records, nil
		})
}
