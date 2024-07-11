package queryfunc

import (
	"fmt"
	"text/template"

	"github.com/engineeringinflow/inflow-backend/pkg/db"
	"github.com/engineeringinflow/inflow-backend/pkg/helper"
	"github.com/engineeringinflow/inflow-backend/pkg/models"
)

type AddressAlias struct {
	*models.Address

	Coordinate *models.Coordinate `gorm:"embedded;embeddedPrefix:c__"`
}

type AddressBuilderOptions struct {
	QueryBuilderOptions
}

func NewAddressBuilder(options AddressBuilderOptions) *Builder {
	var rawSQL = `
	SELECT a.*,

	/* Location */
	c.id AS c__id,
	c.address_number AS c__address_number,
	c.lat AS c__lat,
	c.lng AS c__lng,
	c.formatted_address AS c__formatted_address,
	c.street AS c__street,
	c.level_1 AS c__level_1,
	c.level_2 AS c__level_2,
	c.level_3 AS c__level_3,
	c.level_4 AS c__level_4,
	c.postal_code AS c__postal_code,
	c.country_code AS c__country_code,
	c.timezone_name AS c__timezone_name,
	c.timezone_offset AS c__timezone_offset,
	c.place_id AS c__place_id

	FROM addresses a
	LEFT JOIN coordinates c ON c.id = a.coordinate_id
	`
	var countSQL = `
	SELECT 1

	FROM addresses a
	LEFT JOIN coordinates c ON c.id = a.coordinate_id
	`

	return NewBuilder(rawSQL, countSQL).
		WithOptions(options, template.FuncMap{
			"Description": func() string {
				return fmt.Sprintf(
					"%s:%s:%s - %s",
					helper.GetFuncName(16),
					helper.GetFuncName(17),
					helper.GetFuncName(18),
					options.Role.DisplayName(),
				)
			},
		}).
		WithPaginationFunc(func(db, rawSQL *db.DB) (interface{}, error) {
			var records = make([]*models.Address, rawSQL.RowsAffected)

			rows, err := rawSQL.Rows()
			if err != nil {
				return nil, err

			}
			defer rows.Close()

			for rows.Next() {
				var alias AddressAlias
				err = db.ScanRows(rows, &alias)
				if err != nil {
					db.CustomLogger.Errorf("Scan rows error", err)
					continue
				}

				alias.Address.Coordinate = alias.Coordinate
				records = append(records, alias.Address)
			}

			return &records, nil
		})
}
