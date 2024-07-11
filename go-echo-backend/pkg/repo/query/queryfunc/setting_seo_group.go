package queryfunc

import (
	"encoding/json"
	"text/template"

	"github.com/engineeringinflow/inflow-backend/pkg/db"
	"github.com/engineeringinflow/inflow-backend/pkg/helper"
	"github.com/engineeringinflow/inflow-backend/pkg/models"
	"github.com/engineeringinflow/inflow-backend/pkg/models/enums"
	"gorm.io/datatypes"
)

type SettingSEOLanguageAlias struct {
	Route string         `json:"route"`
	Data  datatypes.JSON `json:"data"`
}

func NewSettingSEOLanguageGroupBuilder(options SettingSEOBuilderOptions) *Builder {
	var rawSQL = `
	SELECT /* {{Description}} */ ss.route, json_agg(ss.*) as data

	FROM setting_seos ss
	`
	var countSQL = `
	SELECT /* {{Description}} */ 1

	FROM setting_seos ss
	`

	var groupBy = "route"

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
			var records = make([]*models.SettingSEOLanguageGroup, rawSQL.RowsAffected)

			rows, err := rawSQL.Rows()
			if err != nil {
				db.CustomLogger.Errorf("Scan rows error", err)
				return nil, err
			}
			defer rows.Close()

			for rows.Next() {
				var alias SettingSEOLanguageAlias
				err = db.ScanRows(rows, &alias)
				if err != nil {
					db.CustomLogger.Errorf("Scan rows error", err)
					continue
				}
				var seoSlice models.SettingSEOSlice
				var b []byte
				b, err = alias.Data.MarshalJSON()
				if err != nil {
					continue
				}
				if err = json.Unmarshal(b, &seoSlice); err != nil {
					continue
				}
				l := models.SettingSEOLanguageGroup{
					Route: alias.Route,
				}
				for _, v := range seoSlice {
					switch enums.LanguageCode(v.LanguageCode) {
					case enums.LanguageCodeEnglish:
						l.EN = v
					case enums.LanguageCodeVietnam:
						l.VI = v
					}
				}
				records = append(records, &l)
			}

			return &records, nil
		})
}
