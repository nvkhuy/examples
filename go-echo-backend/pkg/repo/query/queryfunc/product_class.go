package queryfunc

import (
	"github.com/engineeringinflow/inflow-backend/pkg/db"
	"github.com/engineeringinflow/inflow-backend/pkg/helper"
	"github.com/engineeringinflow/inflow-backend/pkg/models"
	"text/template"
)

type ProductClassAlias struct {
	*models.ProductClass
}

type ProductClassBuilderOptions struct {
	QueryBuilderOptions
}

func NewProductClassBuilder(options ProductClassBuilderOptions) *Builder {
	var rawSQL = `
	SELECT /* {{Description}} */ pc.*

	FROM product_classes pc
	`
	var countSQL = `
	SELECT /* {{Description}} */ 1

	FROM product_classes pc
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
		WithOrderBy("pc.conf DESC").
		WithPaginationFunc(func(db, rawSQL *db.DB) (interface{}, error) {
			var records = make([]*models.ProductClass, rawSQL.RowsAffected)

			rows, err := rawSQL.Rows()
			if err != nil {
				return nil, err

			}
			defer rows.Close()

			for rows.Next() {
				var alias ProductClassAlias
				err = db.ScanRows(rows, &alias)
				if err != nil {
					continue
				}

				records = append(records, alias.ProductClass)
			}

			return &records, nil
		})
}
