package queryfunc

import (
	"text/template"

	"github.com/engineeringinflow/inflow-backend/pkg/helper"
)

type StatsSuppliersBuilderOptions struct {
	QueryBuilderOptions
}

func NewStatsSuppliersBuilder(options StatsSuppliersBuilderOptions) *Builder {
	var rawSQL = `
	SELECT /* {{Description}} */ COUNT(1) AS total_records

	FROM shops s
	`

	return NewBuilder(rawSQL).
		WithOptions(options, template.FuncMap{
			"Description": func() string {
				return helper.JoinNonEmptyStrings(
					"-",
					GetCaller(),
					options.Role.DisplayName(),
				)
			},
		})
}
