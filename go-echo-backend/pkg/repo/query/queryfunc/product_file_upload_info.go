package queryfunc

import (
	"text/template"

	"github.com/engineeringinflow/inflow-backend/pkg/db"
	"github.com/engineeringinflow/inflow-backend/pkg/helper"
	"github.com/engineeringinflow/inflow-backend/pkg/models"
)

type ProductFileUploadInfoOptions struct {
	QueryBuilderOptions
}

func NewProductFileUploadInfoBuilder(options ProductFileUploadInfoOptions) *Builder {
	var rawSQL = `
	SELECT /* {{Description}} */ p.*

	FROM product_file_upload_infos p
	`
	var countSQL = `
	SELECT /* {{Description}}  */ 1

	FROM product_file_upload_infos p
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
		WithOrderBy("p.created_at DESC").
		WithPaginationFunc(func(db, rawSQL *db.DB) (interface{}, error) {
			var records = make([]*models.ProductFileUploadInfo, rawSQL.RowsAffected)

			rows, err := rawSQL.Rows()
			if err != nil {
				return nil, err
			}
			defer rows.Close()

			for rows.Next() {
				var alias models.ProductFileUploadInfo
				err = db.ScanRows(rows, &alias)
				if err != nil {
					db.CustomLogger.Errorf("Scan rows error", err)
					continue
				}

				records = append(records, &alias)
			}

			return &records, nil
		})
}
