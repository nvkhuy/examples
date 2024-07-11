package queryfunc

import (
	"text/template"

	"github.com/engineeringinflow/inflow-backend/pkg/db"
	"github.com/engineeringinflow/inflow-backend/pkg/helper"
	"github.com/engineeringinflow/inflow-backend/pkg/models"
)

type VariantAlias struct {
	*models.Variant
}

type VariantBuilderOptions struct {
	QueryBuilderOptions
}

func NewVariantBuilder(options VariantBuilderOptions) *Builder {
	var rawSQL = `
	SELECT /* {{Description}} */ v.*

	FROM variants v
	`
	var countSQL = `
	SELECT /* {{Description}} */ 1

	FROM variants v
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
		WithOptions(options, template.FuncMap{
			"Description": func() string {
				return helper.JoinNonEmptyStrings(
					"-",
					GetCaller(),
					options.Role.DisplayName(),
				)
			},
		}).
		WithOrderBy("v.updated_at DESC").
		WithPaginationFunc(func(db, rawSQL *db.DB) (interface{}, error) {
			var records = make([]*models.Variant, rawSQL.RowsAffected)

			rows, err := rawSQL.Rows()
			if err != nil {
				return nil, err

			}
			defer rows.Close()

			for rows.Next() {
				var copy VariantAlias
				err = db.ScanRows(rows, &copy)
				if err != nil {
					continue
				}
				records = append(records, copy.Variant)
			}

			return &records, nil
		})
}
