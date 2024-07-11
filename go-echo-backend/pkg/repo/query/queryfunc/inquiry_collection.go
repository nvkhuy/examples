package queryfunc

import (
	"text/template"

	"github.com/engineeringinflow/inflow-backend/pkg/db"
	"github.com/engineeringinflow/inflow-backend/pkg/helper"
	"github.com/engineeringinflow/inflow-backend/pkg/models"
)

type InquiryCollectionAlias struct {
	*models.InquiryCollection
}

type InquiryCollectionBuilderOptions struct {
	QueryBuilderOptions
}

func NewInquiryCollectionBuilder(options InquiryCollectionBuilderOptions) *Builder {
	var rawSQL = `
	SELECT /* {{Description}} */ c.*

	FROM inquiry_collections c
	`
	var countSQL = `
	SELECT /* {{Description}} */ 1

	FROM inquiry_collections c
	`

	if options.Role.IsAdmin() {
		rawSQL = `
		SELECT /* {{Description}} */ c.*

		FROM inquiry_collections c
		`
		countSQL = `
		SELECT /* {{Description}} */ 1

		FROM inquiry_collections c
		`
	}

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
		WithOrderBy("c.name ASC").
		WithPaginationFunc(func(db, rawSQL *db.DB) (interface{}, error) {
			var records = make([]*models.InquiryCollection, rawSQL.RowsAffected)

			rows, err := rawSQL.Rows()
			if err != nil {
				return nil, err

			}
			defer rows.Close()

			for rows.Next() {
				var alias InquiryCollectionAlias
				err = db.ScanRows(rows, &alias)
				if err != nil {
					db.CustomLogger.Errorf("Scan rows error", err)
					continue
				}

				records = append(records, alias.InquiryCollection)
			}

			return &records, nil
		})
}
