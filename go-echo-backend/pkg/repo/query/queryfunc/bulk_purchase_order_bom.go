package queryfunc

import (
	"text/template"

	"github.com/engineeringinflow/inflow-backend/pkg/db"
	"github.com/engineeringinflow/inflow-backend/pkg/helper"
	"github.com/engineeringinflow/inflow-backend/pkg/models"
)

type BulkPurchaseOrderBomAlias struct {
	*models.Bom
}

type BulkPurchaseOrderBomBuilderOptions struct {
	QueryBuilderOptions
}

func NewBulkPurchaseOrderBomBuilder(options BulkPurchaseOrderBomBuilderOptions) *Builder {
	var rawSQL = `
	SELECT /* {{Description}} */ b.*

	FROM boms b
	`
	var countSQL = `
	SELECT /* {{Description}} */ 1

	FROM boms b
	`

	var orderBy = "b.updated_at DESC"

	return NewBuilder(rawSQL, countSQL).
		WithOrderBy(orderBy).
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
			var records = make([]*models.Bom, rawSQL.RowsAffected)

			rows, err := rawSQL.Rows()
			if err != nil {
				return nil, err

			}
			defer rows.Close()

			for rows.Next() {
				var alias BulkPurchaseOrderBomAlias
				err = db.ScanRows(rows, &alias)
				if err != nil {
					db.CustomLogger.Errorf("Scan rows error", err)
					continue
				}

				records = append(records, alias.Bom)
			}

			return records, nil
		})
}
