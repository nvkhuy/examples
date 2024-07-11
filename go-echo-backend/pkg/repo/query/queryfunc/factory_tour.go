package queryfunc

import (
	"text/template"

	"github.com/engineeringinflow/inflow-backend/pkg/db"
	"github.com/engineeringinflow/inflow-backend/pkg/helper"
	"github.com/engineeringinflow/inflow-backend/pkg/models"
)

type FactoryTourAlias struct {
	*models.FactoryTour
}

type FactoryTourBuilderOptions struct {
	QueryBuilderOptions
}

func NewFactoryTourBuilder(options FactoryTourBuilderOptions) *Builder {
	var rawSQL = `
	SELECT /* {{Description}} */ t.*

	FROM factory_tours t
	`
	var countSQL = `
	SELECT /* {{Description}} */ 1

	FROM factory_tours t
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
		WithOrderBy("t.updated_at DESC").
		WithPaginationFunc(func(db, rawSQL *db.DB) (interface{}, error) {
			var records = make([]*models.FactoryTour, rawSQL.RowsAffected)

			rows, err := rawSQL.Rows()
			if err != nil {
				return nil, err

			}
			defer rows.Close()

			for rows.Next() {
				var alias FactoryTourAlias
				err = db.ScanRows(rows, &alias)
				if err != nil {
					db.CustomLogger.Errorf("Scan rows error", err)
					continue
				}

				records = append(records, alias.FactoryTour)
			}

			return &records, nil
		})
}
