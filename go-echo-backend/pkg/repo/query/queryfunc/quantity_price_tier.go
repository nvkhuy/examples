package queryfunc

import (
	"text/template"

	"github.com/engineeringinflow/inflow-backend/pkg/db"
	"github.com/engineeringinflow/inflow-backend/pkg/helper"
	"github.com/engineeringinflow/inflow-backend/pkg/models"
)

type QuantityPriceTierAlias struct {
	*models.QuantityPriceTier
}

type QuantityPriceTierBuilderOptions struct {
	QueryBuilderOptions
}

func NewQuantityPriceTierBuilder(options QuantityPriceTierBuilderOptions) *Builder {
	var rawSQL = `
	SELECT /* {{Description}} */ tier.*

	FROM quantity_price_tiers tier
	`
	var countSQL = `
	SELECT /* QuantityPriceTierBuilder - %s */ 1

	FROM quantity_price_tiers tier
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
		WithOrderBy("tier.min_quantity ASC").
		WithPaginationFunc(func(db, rawSQL *db.DB) (interface{}, error) {
			var records = make([]*models.QuantityPriceTier, rawSQL.RowsAffected)

			rows, err := rawSQL.Rows()
			if err != nil {
				return nil, err

			}
			defer rows.Close()

			for rows.Next() {
				var alias QuantityPriceTierAlias
				err = db.ScanRows(rows, &alias)
				if err != nil {
					db.CustomLogger.Errorf("Scan rows error", err)
					continue
				}

				records = append(records, alias.QuantityPriceTier)
			}

			return &records, nil
		})
}
