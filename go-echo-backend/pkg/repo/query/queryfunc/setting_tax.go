package queryfunc

import (
	"text/template"

	"github.com/engineeringinflow/inflow-backend/pkg/db"
	"github.com/engineeringinflow/inflow-backend/pkg/helper"
	"github.com/engineeringinflow/inflow-backend/pkg/models"
	"github.com/lib/pq"
)

type SettingTaxAlias struct {
	*models.SettingTax

	TaxIDs pq.StringArray `gorm:"type:varchar(255)[]"`
}

type SettingTaxBuilderOptions struct {
	QueryBuilderOptions
}

func NewSettingTaxBuilder(options SettingTaxBuilderOptions) *Builder {
	var rawSQL = `
	SELECT /* {{Description}} */ st.country_code, array_agg(st2.id) AS tax_ids

	FROM setting_taxes st
	JOIN setting_taxes st2 ON st2.country_code = st.country_code
	`
	var countSQL = `
	SELECT /* {{Description}} */ 1

	FROM setting_taxes st
	JOIN setting_taxes st2 ON st2.country_code = st.country_code
	`

	var groupBy = "st.country_code"

	return NewBuilder(rawSQL, countSQL).
		WithGroupBy(groupBy).
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
			var records = make([]*models.SettingTax, rawSQL.RowsAffected)

			rows, err := rawSQL.Rows()
			if err != nil {
				db.CustomLogger.Errorf("Scan rows error", err)
				return nil, err
			}
			defer rows.Close()

			var taxIDs []string
			for rows.Next() {
				var alias SettingTaxAlias
				err = db.ScanRows(rows, &alias)
				if err != nil {
					db.CustomLogger.Errorf("Scan rows error", err)
					continue
				}

				for _, id := range alias.TaxIDs {
					if !helper.StringContains(taxIDs, id) {
						taxIDs = append(taxIDs, id)
					}
				}

				records = append(records, alias.SettingTax)
			}

			if len(taxIDs) > 0 {
				var taxes []*models.SettingTax
				db.Find(&taxes, "id IN ?", taxIDs)

				for _, tax := range taxes {
					for _, record := range records {
						if tax.CountryCode == record.CountryCode {
							record.Taxes = append(record.Taxes, tax)
						}
					}
				}
			}
			return &records, nil
		})
}
