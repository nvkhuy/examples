package queryfunc

import (
	"text/template"

	"github.com/engineeringinflow/inflow-backend/pkg/db"
	"github.com/engineeringinflow/inflow-backend/pkg/helper"
	"github.com/engineeringinflow/inflow-backend/pkg/models"
)

type UserCustomerIOBuilderOptions struct {
	QueryBuilderOptions

	FromTime int
	ToTime   int
}

func NewUserCustomerIOBuilder(options UserCustomerIOBuilderOptions) *Builder {
	var rawSQL = `
	SELECT /* {{Description}} */ 
	u.id,
	u.role, 
	u.account_status, 
	u.created_at, 
	u.phone_number, 
	u.email,
	u.timezone, 
	u.country_code, 
	u.last_login,
	u.last_name,
	u.first_name,
	u.name,
	u.stripe_customer_id,
	u.avatar

	FROM users u
	`
	var countSQL = `
	SELECT /* {{Description}} */ 1
	
	FROM users u
	`

	var groupBy = "u.id"

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
			var records []*models.User
			var err = rawSQL.Find(&records).Error

			return records, err
		}).
		WithGroupBy(groupBy)

}
