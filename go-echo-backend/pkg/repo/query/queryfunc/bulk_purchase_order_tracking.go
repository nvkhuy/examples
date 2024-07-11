package queryfunc

import (
	"text/template"

	"github.com/engineeringinflow/inflow-backend/pkg/db"
	"github.com/engineeringinflow/inflow-backend/pkg/helper"
	"github.com/engineeringinflow/inflow-backend/pkg/models"
)

type BulkPurchaseOrderTrackingAlias struct {
	*models.BulkPurchaseOrderTracking
}

type BulkPurchaseOrderTrackingBuilderOptions struct {
	QueryBuilderOptions
}

func NewBulkPurchaseOrderTrackingBuilder(options BulkPurchaseOrderTrackingBuilderOptions) *Builder {
	var rawSQL = `
	SELECT /* {{Description}} */ pot.*

	FROM bulk_purchase_order_trackings pot
	`
	var countSQL = `
	SELECT /* {{Description}} */ 1

	FROM bulk_purchase_order_trackings pot
	`

	if options.Role.IsSeller() {
		rawSQL = `
		SELECT /* {{Description}} */ pot.*

		FROM bulk_purchase_order_trackings pot
		LEFT JOIN bulk_purchase_order_seller_quotations bposq ON bposq.bulk_purchase_order_id = pot.purchase_order_id AND bposq.user_id = @user_id
		`
		countSQL = `
		SELECT /* {{Description}} */ 1

		FROM bulk_purchase_order_trackings pot
		LEFT JOIN bulk_purchase_order_seller_quotations bposq ON bposq.bulk_purchase_order_id = pot.purchase_order_id AND bposq.user_id = @user_id
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
		WithOrderBy("pot.created_at DESC").
		WithPaginationFunc(func(db, rawSQL *db.DB) (interface{}, error) {
			var records = make([]*models.BulkPurchaseOrderTracking, rawSQL.RowsAffected)

			rows, err := rawSQL.Rows()
			if err != nil {
				return nil, err

			}
			defer rows.Close()

			for rows.Next() {
				var alias BulkPurchaseOrderTrackingAlias
				err = db.ScanRows(rows, &alias)
				if err != nil {
					db.CustomLogger.Errorf("Scan rows error", err)
					continue
				}
				records = append(records, alias.BulkPurchaseOrderTracking)
			}

			return &records, nil
		})
}
