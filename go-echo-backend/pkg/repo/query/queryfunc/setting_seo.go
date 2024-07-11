package queryfunc

import (
	"text/template"

	"github.com/engineeringinflow/inflow-backend/pkg/db"
	"github.com/engineeringinflow/inflow-backend/pkg/helper"
	"github.com/engineeringinflow/inflow-backend/pkg/models"
)

type SettingSEOAlias struct {
	*models.SettingSEO
}

type SettingSEOBuilderOptions struct {
	QueryBuilderOptions
}

func NewSettingSEOBuilder(options SettingSEOBuilderOptions) *Builder {
	var rawSQL = `
	SELECT /* {{Description}} */ ss.*

	FROM setting_seos ss
	`
	var countSQL = `
	SELECT /* {{Description}} */ 1

	FROM setting_seos ss
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
		WithPaginationFunc(func(db, rawSQL *db.DB) (interface{}, error) {
			var records = make([]*models.SettingSEO, rawSQL.RowsAffected)

			rows, err := rawSQL.Rows()
			if err != nil {
				db.CustomLogger.Errorf("Scan rows error", err)
				return nil, err
			}
			defer rows.Close()

			for rows.Next() {
				var alias SettingSEOAlias
				err = db.ScanRows(rows, &alias)
				if err != nil {
					db.CustomLogger.Errorf("Scan rows error", err)
					continue
				}

				records = append(records, alias.SettingSEO)
			}

			return &records, nil
		})
}
