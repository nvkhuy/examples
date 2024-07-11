package queryfunc

import (
	"text/template"

	"github.com/engineeringinflow/inflow-backend/pkg/db"
	"github.com/engineeringinflow/inflow-backend/pkg/helper"
	"github.com/engineeringinflow/inflow-backend/pkg/models"
)

type InquiryAuditAlias struct {
	*models.InquiryAudit
}

type InquiryAuditBuilderOptions struct {
	QueryBuilderOptions
}

func NewInquiryAuditBuilder(options InquiryAuditBuilderOptions) *Builder {
	var rawSQL = `
	SELECT /* {{Description}} */ ia.*

	FROM inquiry_audits ia
	`
	var countSQL = `
	SELECT /* {{Description}} */ 1

	FROM inquiry_audits ia
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
		WithOrderBy("ia.created_at DESC").
		WithPaginationFunc(func(db, rawSQL *db.DB) (interface{}, error) {
			var records = make([]*models.InquiryAudit, rawSQL.RowsAffected)

			rows, err := rawSQL.Rows()
			if err != nil {
				return nil, err

			}
			defer rows.Close()

			for rows.Next() {
				var alias InquiryAuditAlias
				err = db.ScanRows(rows, &alias)
				if err != nil {
					db.CustomLogger.Errorf("Scan rows error", err)
					continue
				}
				records = append(records, alias.InquiryAudit)
			}

			return &records, nil
		})
}
