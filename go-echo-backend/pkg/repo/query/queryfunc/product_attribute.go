package queryfunc

import (
	"text/template"

	"github.com/engineeringinflow/inflow-backend/pkg/db"
	"github.com/engineeringinflow/inflow-backend/pkg/helper"
	"github.com/engineeringinflow/inflow-backend/pkg/models"
)

type ProductAttributeAlias struct {
	*models.ProductAttribute
}

type ProductAttributeBuilderOptions struct {
	QueryBuilderOptions
}

func NewProductAttributeBuilder(options ProductAttributeBuilderOptions) *Builder {
	var rawSQL = `
	SELECT /* {{Description}} */ att.*

	FROM product_attributes att
	`
	var countSQL = `
	SELECT /* {{Description}} */ 1

	FROM product_attributes att
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
		WithOrderBy("att.order ASC").
		WithPaginationFunc(func(db, rawSQL *db.DB) (interface{}, error) {
			var records = make([]*models.ProductAttribute, rawSQL.RowsAffected)

			rows, err := rawSQL.Rows()
			if err != nil {
				return nil, err

			}
			defer rows.Close()

			for rows.Next() {
				var alias ProductAttributeAlias
				err = db.ScanRows(rows, &alias)
				if err != nil {
					db.CustomLogger.Errorf("Scan rows error", err)
					continue
				}

				records = append(records, alias.ProductAttribute)
			}

			return &records, nil
		})
}
