package queryfunc

import (
	"text/template"

	"github.com/engineeringinflow/inflow-backend/pkg/db"
	"github.com/engineeringinflow/inflow-backend/pkg/helper"
	"github.com/engineeringinflow/inflow-backend/pkg/models"
)

type ReleaseNoteAlias struct {
	*models.ReleaseNote
}

type ReleaseNoteBuilderOptions struct {
	QueryBuilderOptions
}

func NewReleaseNoteBuilder(options ReleaseNoteBuilderOptions) *Builder {
	var rawSQL = `
	SELECT /* {{Description}} */ r.*

	FROM release_notes r
	`
	var countSQL = `
	SELECT /* {{Description}} */ 1

	FROM release_notes r
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
		WithOrderBy("r.release_date DESC").
		WithPaginationFunc(func(db, rawSQL *db.DB) (interface{}, error) {
			var records = make([]*models.ReleaseNote, rawSQL.RowsAffected)

			rows, err := rawSQL.Rows()
			if err != nil {
				return nil, err

			}
			defer rows.Close()

			for rows.Next() {
				var alias ReleaseNoteAlias
				err = db.ScanRows(rows, &alias)
				if err != nil {
					db.CustomLogger.Errorf("Scan rows error", err)
					continue
				}

				records = append(records, alias.ReleaseNote)
			}

			return &records, nil
		})
}
