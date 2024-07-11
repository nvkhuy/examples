package queryfunc

import (
	"text/template"

	"github.com/engineeringinflow/inflow-backend/pkg/db"
	"github.com/engineeringinflow/inflow-backend/pkg/helper"
	"github.com/engineeringinflow/inflow-backend/pkg/models"
)

type CmsNotificationAlias struct {
	*models.CmsNotification
}

type CmsNotificationBuilderOptions struct {
	QueryBuilderOptions
}

func NewCmsNotificationBuilder(options CmsNotificationBuilderOptions) *Builder {
	var rawSQL = `
	SELECT /* {{Description}} */ cn.*

	FROM cms_notifications cn
	`
	var countSQL = `
	SELECT /* {{Description}} */ 1

	FROM cms_notifications cn
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
		WithOrderBy("cn.created_at DESC").
		WithPaginationFunc(func(db, rawSQL *db.DB) (interface{}, error) {
			var records = make([]*models.CmsNotification, rawSQL.RowsAffected)

			rows, err := rawSQL.Rows()
			if err != nil {
				return nil, err

			}
			defer rows.Close()

			for rows.Next() {
				var alias CmsNotificationAlias
				err = db.ScanRows(rows, &alias)
				if err != nil {
					db.CustomLogger.Errorf("Scan rows error", err)
					continue
				}

				records = append(records, alias.CmsNotification)
			}

			return &records, nil
		})
}
