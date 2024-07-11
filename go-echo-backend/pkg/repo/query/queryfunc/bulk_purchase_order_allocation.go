package queryfunc

import (
	"text/template"

	"github.com/engineeringinflow/inflow-backend/pkg/db"
	"github.com/engineeringinflow/inflow-backend/pkg/helper"
	"github.com/engineeringinflow/inflow-backend/pkg/models"
)

type BulkPurchaseOrderAllocationAlias struct {
	*models.User
}

type BulkPurchaseOrderAllocationBuilderOptions struct {
	QueryBuilderOptions
}

func NewBulkPurchaseOrderAllocationBuilder(options BulkPurchaseOrderAllocationBuilderOptions) *Builder {
	var rawSQL = `
	SELECT /* {{Description}} */ u.*

	FROM users u
	JOIN bulk_purchase_order_seller_quotations bposq ON bposq.user_id = u.id
	`
	var countSQL = `
	SELECT /* {{Description}} */ u.*

	FROM users u
	JOIN bulk_purchase_order_seller_quotations bposq ON bposq.user_id = u.id
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
		WithOrderBy("u.created_at DESC").
		WithPaginationFunc(func(db, rawSQL *db.DB) (interface{}, error) {
			var records = make([]*models.User, rawSQL.RowsAffected)

			rows, err := rawSQL.Rows()
			if err != nil {
				return nil, err

			}
			defer rows.Close()

			for rows.Next() {
				var alias BulkPurchaseOrderAllocationAlias
				err = db.ScanRows(rows, &alias)
				if err != nil {
					db.CustomLogger.Errorf("Scan rows error", err)
					continue
				}

				records = append(records, alias.User)
			}

			return &records, nil
		})
}
