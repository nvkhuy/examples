package queryfunc

import (
	"text/template"

	"github.com/engineeringinflow/inflow-backend/pkg/db"
	"github.com/engineeringinflow/inflow-backend/pkg/helper"
	"github.com/engineeringinflow/inflow-backend/pkg/models"
)

type DocumentCategoryBuilderOptions struct {
	QueryBuilderOptions
}

func NewDocumentCategoryBuilder(options DocumentCategoryBuilderOptions) *Builder {
	var rawSQL = `
	SELECT /* {{Description}} */ cate.*

	FROM document_categories cate
	`

	var countSQL = `
	SELECT /* {{Description}} */ 1

	FROM document_categories cate
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
		WithOrderBy("cate.updated_at DESC").
		WithPaginationFunc(func(db, rawSQL *db.DB) (interface{}, error) {
			var records = make([]*models.DocumentCategory, rawSQL.RowsAffected)

			rows, err := rawSQL.Rows()
			if err != nil {
				db.CustomLogger.Errorf("Scan rows error", err)
				return nil, err
			}
			defer rows.Close()

			for rows.Next() {
				var copy models.DocumentCategory
				err = db.ScanRows(rows, &copy)
				if err != nil {
					continue
				}
				records = append(records, &copy)
			}

			return &records, nil
		})
}
