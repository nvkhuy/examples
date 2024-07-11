package queryfunc

import (
	"text/template"

	"github.com/engineeringinflow/inflow-backend/pkg/db"
	"github.com/engineeringinflow/inflow-backend/pkg/helper"
	"github.com/engineeringinflow/inflow-backend/pkg/models"
)

type DocumentTagBuilderOptions struct {
	QueryBuilderOptions
}

func NewDocumentTagBuilder(options DocumentTagBuilderOptions) *Builder {
	var rawSQL = `
	SELECT /* {{Description}} */ t.*

	FROM document_tags t
	`

	var countSQL = `
	SELECT /* {{Description}} */ 1

	FROM document_tags t
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
		WithOrderBy("t.updated_at DESC").
		WithPaginationFunc(func(db, rawSQL *db.DB) (interface{}, error) {
			var records = make([]*models.DocumentTag, rawSQL.RowsAffected)

			rows, err := rawSQL.Rows()
			if err != nil {
				db.CustomLogger.Errorf("Scan rows error", err)
				return nil, err
			}
			defer rows.Close()

			for rows.Next() {
				var copy models.DocumentTag
				err = db.ScanRows(rows, &copy)
				if err != nil {
					continue
				}
				records = append(records, &copy)
			}

			return &records, nil
		})
}
