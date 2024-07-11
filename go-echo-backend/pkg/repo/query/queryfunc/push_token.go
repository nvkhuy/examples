package queryfunc

import (
	"text/template"

	"github.com/engineeringinflow/inflow-backend/pkg/db"
	"github.com/engineeringinflow/inflow-backend/pkg/helper"
	"github.com/engineeringinflow/inflow-backend/pkg/models"
)

type PushTokenAlias struct {
	*models.PushToken
}

type PushTokenBuilderOptions struct {
	QueryBuilderOptions
}

func NewPushTokenBuilder(options PushTokenBuilderOptions) *Builder {
	var rawSQL = `
	SELECT /* {{Description}} */ pt.*
	FROM push_tokens pt
	`
	var countSQL = `
	SELECT /* {{Description}} */ 1
	FROM push_tokens pt
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
		WithOrderBy("pt.last_used DESC").
		WithPaginationFunc(func(db, rawSQL *db.DB) (interface{}, error) {
			var records = make([]*models.PushToken, rawSQL.RowsAffected)

			rows, err := rawSQL.Rows()
			if err != nil {
				db.CustomLogger.ErrorAny(err)
				return nil, err

			}
			defer rows.Close()

			for rows.Next() {
				var copy PushTokenAlias
				err = db.ScanRows(rows, &copy)
				if err != nil {
					continue
				}

				records = append(records, copy.PushToken)
			}

			return &records, nil
		})
}
