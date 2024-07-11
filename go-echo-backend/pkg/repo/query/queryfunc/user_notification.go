package queryfunc

import (
	"text/template"

	"github.com/engineeringinflow/inflow-backend/pkg/db"
	"github.com/engineeringinflow/inflow-backend/pkg/helper"
	"github.com/engineeringinflow/inflow-backend/pkg/models"
)

type UserNotificationAlias struct {
	*models.UserNotification
}

type UserNotificationBuilderOptions struct {
	QueryBuilderOptions
}

func NewUserNotificationBuilder(options UserNotificationBuilderOptions) *Builder {
	var rawSQL = `
	SELECT /* {{Description}} */ un.*

	FROM user_notifications un
	`
	var countSQL = `
	SELECT /* {{Description}} */ 1

	FROM user_notifications un
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
		WithOrderBy("un.created_at DESC").
		WithPaginationFunc(func(db, rawSQL *db.DB) (interface{}, error) {
			var records = make([]*models.UserNotification, rawSQL.RowsAffected)

			rows, err := rawSQL.Rows()
			if err != nil {
				return nil, err

			}
			defer rows.Close()

			for rows.Next() {
				var alias UserNotificationAlias
				err = db.ScanRows(rows, &alias)
				if err != nil {
					db.CustomLogger.Errorf("Scan rows error", err)
					continue
				}

				records = append(records, alias.UserNotification)
			}

			return &records, nil
		})
}
