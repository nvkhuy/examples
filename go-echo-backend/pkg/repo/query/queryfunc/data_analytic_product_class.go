package queryfunc

import (
	"github.com/engineeringinflow/inflow-backend/pkg/db"
	"github.com/engineeringinflow/inflow-backend/pkg/helper"
	"github.com/engineeringinflow/inflow-backend/pkg/models"
	"text/template"
)

type AnalyticProductClassAlias struct {
	models.AnalyticProductClass
	Total int64 `json:"total,omitempty"`
}

type AnalyticProductClassBuilderOptions struct {
	QueryBuilderOptions
	IsGroupByClass bool
}

func NewAnalyticProductClassBuilder(options AnalyticProductClassBuilderOptions) *Builder {
	var rawSQL = `
	SELECT /* {{Description}} */ pc.*
	FROM product_classes pc
	`
	var countSQL = `
	SELECT /* {{Description}} */ 1
	FROM product_classes pc
	`

	if options.IsGroupByClass {
		rawSQL = `
	SELECT /* {{Description}} */ class,count(1) as total

	FROM product_classes pc
	`
		countSQL = `
	SELECT /* {{Description}} */ 1

	FROM product_classes pc
	`
	}

	builder := NewBuilder(rawSQL, countSQL).
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
			var records = make([]AnalyticProductClassAlias, rawSQL.RowsAffected)

			rows, err := rawSQL.Rows()
			if err != nil {
				return nil, err
			}
			defer rows.Close()

			for rows.Next() {
				var alias AnalyticProductClassAlias
				err = db.ScanRows(rows, &alias)
				if err != nil {
					db.CustomLogger.Errorf("Scan rows error", err)
					continue
				}
				records = append(records, alias)
			}

			return &records, nil
		})
	return builder
}
