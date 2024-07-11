package queryfunc

import (
	"text/template"

	"github.com/engineeringinflow/inflow-backend/pkg/db"
	"github.com/engineeringinflow/inflow-backend/pkg/helper"
	"github.com/engineeringinflow/inflow-backend/pkg/models"
)

type UserRoleAlias struct {
	Team  string `json:"team"`
	Count int    `json:"count"`
}

type UserRoleBuilderOptions struct {
	QueryBuilderOptions
}

func NewUserRoleBuilder(options UserRoleBuilderOptions) *Builder {
	var rawSQL = `
	SELECT /* {{Description}} */ u.team, count(1) as count

	FROM users u
	`
	var countSQL = `
	SELECT /* {{Description}} */ 1

	FROM users u
	`
	var groupBy = `team`

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
		WithGroupBy(groupBy).
		WithOrderBy("u.team").
		WithPaginationFunc(func(db, rawSQL *db.DB) (interface{}, error) {
			var records = make([]models.UserRoleStat, rawSQL.RowsAffected)

			rows, err := rawSQL.Rows()
			if err != nil {
				return nil, err

			}
			defer func() {
				_ = rows.Close()
			}()

			for rows.Next() {
				var alias UserRoleAlias
				err = db.ScanRows(rows, &alias)
				if err != nil {
					db.CustomLogger.Errorf("Scan rows error", err)
					continue
				}
				records = append(records, models.UserRoleStat{
					Team:  alias.Team,
					Count: alias.Count,
				})
			}
			return &records, nil
		})
}
