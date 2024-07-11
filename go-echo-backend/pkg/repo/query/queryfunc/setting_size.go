package queryfunc

import (
	"text/template"

	"github.com/engineeringinflow/inflow-backend/pkg/db"
	"github.com/engineeringinflow/inflow-backend/pkg/helper"
	"github.com/engineeringinflow/inflow-backend/pkg/models"
	"github.com/lib/pq"
)

type SettingSizeAlias struct {
	*models.SettingSize

	SizeIDs pq.StringArray `gorm:"type:varchar(255)[]"`
}

type SettingSizeBuilderOptions struct {
	QueryBuilderOptions
}

func NewSettingSizeBuilder(options SettingSizeBuilderOptions) *Builder {
	var rawSQL = `
	SELECT /* {{Description}} */ ss.type, array_agg(ss2.id) AS size_ids
	
	FROM setting_sizes ss
	JOIN setting_sizes ss2 ON ss.type = ss2.type
	`
	var countSQL = `
	SELECT /* {{Description}} */ 1

	FROM setting_sizes ss
	JOIN setting_sizes ss2 ON ss.type = ss2.type
	`

	var groupBy = "ss.type"

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
			var records = make([]*models.SettingSize, rawSQL.RowsAffected)

			rows, err := rawSQL.Rows()
			if err != nil {
				db.CustomLogger.Errorf("Scan rows error", err)
				return nil, err
			}
			defer rows.Close()

			var sizeIDs []string
			for rows.Next() {
				var alias SettingSizeAlias
				err = db.ScanRows(rows, &alias)
				if err != nil {
					db.CustomLogger.Errorf("Scan rows error", err)
					continue
				}

				for _, id := range alias.SizeIDs {
					if !helper.StringContains(sizeIDs, id) {
						sizeIDs = append(sizeIDs, id)
					}
				}

				records = append(records, alias.SettingSize)
			}

			if len(sizeIDs) > 0 {
				var sizes []*models.SettingSize
				db.Order("updated_at DESC").Find(&sizes, "id IN ?", sizeIDs)

				for _, size := range sizes {
					for _, record := range records {
						if size.Type == record.Type {
							record.Sizes = append(record.Sizes, size)
						}
					}
				}
			}

			return &records, nil
		})
}
