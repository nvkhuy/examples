package queryfunc

import (
	"text/template"

	"github.com/engineeringinflow/inflow-backend/pkg/db"
	"github.com/engineeringinflow/inflow-backend/pkg/helper"
	"github.com/engineeringinflow/inflow-backend/pkg/models"
)

type AsFeaturedInAlias struct {
	*models.AsFeaturedIn
}

type AsFeaturedInBuilderOptions struct {
	QueryBuilderOptions
}

func NewAsFeaturedInBuilder(options AsFeaturedInBuilderOptions) *Builder {
	var rawSQL = `
	SELECT /* {{Description}} */ afi.*

	FROM as_featured_ins afi
	`
	var countSQL = `
	SELECT /* {{Description}}  */ 1

	FROM as_featured_ins afi
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
		WithOrderBy("afi.created_at DESC").
		WithPaginationFunc(func(db, rawSQL *db.DB) (interface{}, error) {
			var records = make([]*models.AsFeaturedIn, rawSQL.RowsAffected)

			rows, err := rawSQL.Rows()
			if err != nil {
				return nil, err

			}
			defer rows.Close()

			for rows.Next() {
				var alias AsFeaturedInAlias
				err = db.ScanRows(rows, &alias)
				if err != nil {
					db.CustomLogger.Errorf("Scan rows error", err)
					continue
				}

				records = append(records, alias.AsFeaturedIn)
			}

			return &records, nil
		})
}
