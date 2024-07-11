package queryfunc

import (
	"text/template"

	"github.com/engineeringinflow/inflow-backend/pkg/db"
	"github.com/engineeringinflow/inflow-backend/pkg/helper"
	"github.com/engineeringinflow/inflow-backend/pkg/models"
)

type UserBankAlias struct {
	*models.UserBank
}

type UserBankBuilderOptions struct {
	QueryBuilderOptions
}

func NewUserBankBuilder(options UserBankBuilderOptions) *Builder {
	var rawSQL = `
	SELECT /* {{Description}} */ ub.*

	FROM user_banks ub
	`
	var countSQL = `
	SELECT /* {{Description}} */ 1

	FROM user_banks ub
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
			var records = make([]*models.UserBank, rawSQL.RowsAffected)

			rows, err := rawSQL.Rows()
			if err != nil {
				db.CustomLogger.Errorf("Scan rows error", err)
				return nil, err
			}
			defer rows.Close()

			for rows.Next() {
				var alias UserBankAlias
				err = db.ScanRows(rows, &alias)
				if err != nil {
					db.CustomLogger.Errorf("Scan rows error", err)
					continue
				}

				records = append(records, alias.UserBank)
			}

			return &records, nil
		})
}
