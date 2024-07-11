package queryfunc

import (
	"text/template"

	"github.com/engineeringinflow/inflow-backend/pkg/db"
	"github.com/engineeringinflow/inflow-backend/pkg/helper"
	"github.com/engineeringinflow/inflow-backend/pkg/models"
)

type InquiryOrderAlias struct {
	*models.Inquiry
}

type InquiryOrderBuilderOptions struct {
	QueryBuilderOptions
}

func NewInquiryOrderBuilder(options InquiryOrderBuilderOptions) *Builder {
	var rawSQL = `
	SELECT /* {{Description}} */ iq.*

	FROM inquiries iq
	LEFT JOIN orders o ON iq.id = o.inquiry_id
	`
	var countSQL = `
	SELECT /* {{Description}} */ 1

	FROM inquiries iq
	LEFT JOIN orders o ON iq.id = o.inquiry_id
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
		WithOrderBy("iq.updated_at DESC").
		WithPaginationFunc(func(db, rawSQL *db.DB) (interface{}, error) {
			var records = make([]*models.Inquiry, rawSQL.RowsAffected)

			rows, err := rawSQL.Rows()
			if err != nil {
				return nil, err

			}
			defer rows.Close()

			for rows.Next() {
				var alias InquiryAlias
				err = db.ScanRows(rows, &alias)
				if err != nil {
					db.CustomLogger.Errorf("Scan rows error", err)
					continue
				}

				records = append(records, alias.Inquiry)
			}

			return &records, nil
		})
}
