package queryfunc

import (
	"text/template"

	"github.com/engineeringinflow/inflow-backend/pkg/db"
	"github.com/engineeringinflow/inflow-backend/pkg/helper"
	"github.com/engineeringinflow/inflow-backend/pkg/models"
)

type SysNotificationAlias struct {
	*models.SysNotification
}

type SysNotificationBuilderOptions struct {
	QueryBuilderOptions
}

func NewSysNotificationBuilder(options SysNotificationBuilderOptions) *Builder {
	var rawSQL = `
	SELECT /* {{Description}} */ sn.*

	FROM sys_notifications sn
	`
	var countSQL = `
	SELECT /* {{Description}} */ 1

	FROM sys_notifications sn
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
		WithOrderBy("sn.created_at DESC").
		WithPaginationFunc(func(db, rawSQL *db.DB) (interface{}, error) {
			var records = make([]*models.SysNotification, rawSQL.RowsAffected)

			rows, err := rawSQL.Rows()
			if err != nil {
				return nil, err

			}
			defer rows.Close()

			for rows.Next() {
				var alias SysNotificationAlias
				err = db.ScanRows(rows, &alias)
				if err != nil {
					db.CustomLogger.Errorf("Scan rows error", err)
					continue
				}

				records = append(records, alias.SysNotification)
			}

			return &records, nil
		})
}
