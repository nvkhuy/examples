package queryfunc

import (
	"text/template"

	"github.com/engineeringinflow/inflow-backend/pkg/db"
	"github.com/engineeringinflow/inflow-backend/pkg/helper"
	"github.com/engineeringinflow/inflow-backend/pkg/models"
)

type AdsVideoAlias struct {
	*models.AdsVideo
}

type AdsVideoBuilderOptions struct {
	QueryBuilderOptions
}

func NewAdsVideoBuilder(options AdsVideoBuilderOptions) *Builder {
	var rawSQL = `
	SELECT /* {{Description}} */ ads.*

	FROM ads_videos ads
	`
	var countSQL = `
	SELECT /* {{Description}}  */ 1

	FROM ads_videos ads
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
		WithOrderBy("ads.created_at DESC").
		WithPaginationFunc(func(db, rawSQL *db.DB) (interface{}, error) {
			var records = make([]*models.AdsVideo, rawSQL.RowsAffected)

			rows, err := rawSQL.Rows()
			if err != nil {
				return nil, err

			}
			defer rows.Close()

			for rows.Next() {
				var alias AdsVideoAlias
				err = db.ScanRows(rows, &alias)
				if err != nil {
					db.CustomLogger.Errorf("Scan rows error", err)
					continue
				}

				records = append(records, alias.AdsVideo)
			}

			return &records, nil
		})
}
