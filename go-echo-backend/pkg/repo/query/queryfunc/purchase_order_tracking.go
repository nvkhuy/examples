package queryfunc

import (
	"text/template"

	"github.com/engineeringinflow/inflow-backend/pkg/db"
	"github.com/engineeringinflow/inflow-backend/pkg/helper"
	"github.com/engineeringinflow/inflow-backend/pkg/models"
)

type PurchaseOrderTrackingAlias struct {
	*models.PurchaseOrderTracking
}

type PurchaseOrderTrackingBuilderOptions struct {
	QueryBuilderOptions
}

func NewPurchaseOrderTrackingBuilder(options PurchaseOrderTrackingBuilderOptions) *Builder {
	var rawSQL = `
	SELECT /* {{Description}} */ pot.*

	FROM purchase_order_trackings pot
	`
	var countSQL = `
	SELECT /* {{Description}} */ 1

	FROM purchase_order_trackings pot
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
		WithOrderBy("pot.created_at DESC").
		WithPaginationFunc(func(db, rawSQL *db.DB) (interface{}, error) {
			var records = make([]*models.PurchaseOrderTracking, rawSQL.RowsAffected)

			rows, err := rawSQL.Rows()
			if err != nil {
				return nil, err

			}
			defer rows.Close()

			for rows.Next() {
				var alias PurchaseOrderTrackingAlias
				err = db.ScanRows(rows, &alias)
				if err != nil {
					db.CustomLogger.Errorf("Scan rows error", err)
					continue
				}

				records = append(records, alias.PurchaseOrderTracking)
			}

			return &records, nil
		})
}
