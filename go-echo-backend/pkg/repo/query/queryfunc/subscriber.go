package queryfunc

import (
	"text/template"

	"github.com/engineeringinflow/inflow-backend/pkg/db"
	"github.com/engineeringinflow/inflow-backend/pkg/helper"
	"github.com/engineeringinflow/inflow-backend/pkg/models"
)

type SubscriberAlias struct {
	*models.Subscriber
}

type SubscriberBuilderOptions struct {
	QueryBuilderOptions
}

func NewSubscriberBuilder(options SubscriberBuilderOptions) *Builder {
	var rawSQL = `
	SELECT /* {{Description}} */ sub.*

	FROM subscribers sub
	`
	var countSQL = `
	SELECT /* {{Description}} */ 1

	FROM subscribers sub
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
		WithOrderBy("sub.updated_at DESC").
		WithPaginationFunc(func(db, rawSQL *db.DB) (interface{}, error) {
			var records = make([]*models.Subscriber, rawSQL.RowsAffected)

			rows, err := rawSQL.Rows()
			if err != nil {
				return nil, err

			}
			defer rows.Close()

			for rows.Next() {
				var alias SubscriberAlias
				err = db.ScanRows(rows, &alias)
				if err != nil {
					db.CustomLogger.Errorf("Scan rows error", err)
					continue
				}

				records = append(records, alias.Subscriber)
			}

			return &records, nil
		})
}
