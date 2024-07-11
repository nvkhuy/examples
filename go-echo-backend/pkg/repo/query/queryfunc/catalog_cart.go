package queryfunc

import (
	"sync"
	"text/template"

	"github.com/engineeringinflow/inflow-backend/pkg/db"
	"github.com/engineeringinflow/inflow-backend/pkg/helper"
	"github.com/engineeringinflow/inflow-backend/pkg/models"
	"github.com/engineeringinflow/inflow-backend/pkg/repo/query"
)

type CatalogCartAlias struct {
	*models.PurchaseOrder
}

type CatalogCartBuilderOptions struct {
	QueryBuilderOptions
}

func NewCatalogCartBuilder(options CatalogCartBuilderOptions) *Builder {
	var rawSQL = `
	SELECT /* {{Description}} */ po.*

	FROM purchase_orders po

	`
	var countSQL = `
	SELECT /* {{Description}} */ 1

	FROM purchase_orders po
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
		WithOrderBy("po.updated_at DESC").
		WithPaginationFunc(func(db, rawSQL *db.DB) (interface{}, error) {
			var records = make([]*models.PurchaseOrder, rawSQL.RowsAffected)

			rows, err := rawSQL.Rows()
			if err != nil {
				return nil, err

			}
			defer rows.Close()

			var purchaseOrderIDs []string
			for rows.Next() {
				var alias CatalogCartAlias
				err = db.ScanRows(rows, &alias)
				if err != nil {
					db.CustomLogger.Errorf("Scan rows error", err)
					continue
				}

				purchaseOrderIDs = append(purchaseOrderIDs, alias.ID)
				records = append(records, alias.PurchaseOrder)
			}

			var wg sync.WaitGroup

			if len(purchaseOrderIDs) > 0 {
				wg.Add(1)
				go func() {
					defer wg.Done()
					var cartItems []*models.PurchaseOrderItem
					query.New(db, NewPurchaseOrderItemBuilder(PurchaseOrderItemBuilderOptions{
						IncludeProduct: true,
						IncludeFabric:  true,
						IncludeVariant: true,
					})).
						WhereFunc(func(builder *query.Builder) {
							builder.Where("purchase_order_id IN ?", purchaseOrderIDs)
						}).
						FindFunc(&cartItems)

					for _, cartItem := range cartItems {
						for _, record := range records {
							if cartItem.PurchaseOrderID == record.ID {
								record.Items = append(record.Items, cartItem)
							}
						}
					}
				}()
			}

			wg.Wait()

			return &records, nil
		})
}
