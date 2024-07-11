package queryfunc

import (
	"sync"
	"text/template"

	"github.com/engineeringinflow/inflow-backend/pkg/db"
	"github.com/engineeringinflow/inflow-backend/pkg/helper"
	"github.com/engineeringinflow/inflow-backend/pkg/models"
)

type InquirySellerMatchingAlias struct {
	*models.User

	BusinessProfile *models.BusinessProfile `gorm:"embedded;embeddedPrefix:bu__"`
}

type InquirySellerMatchingBuilderOptions struct {
	QueryBuilderOptions
}

func NewInquirySellerMatchingBuilder(options InquirySellerMatchingBuilderOptions) *Builder {
	var rawSQL = `
	SELECT /* {{Description}} */ u.*,
	bu.id AS bu__id,
	bu.order_types AS bu__order_types,
	bu.product_types AS bu__product_types,
	bu.product_groups AS bu__product_groups,
	bu.product_groups AS bu__product_groups,
	bu.excepted_fabric_types AS bu__excepted_fabric_types,
	bu.flat_mill_fabric_types AS bu__flat_mill_fabric_types,
	bu.mill_fabric_types AS bu__mill_fabric_types

	FROM users u
	LEFT JOIN business_profiles bu ON bu.user_id = u.id
	`
	var countSQL = `
	SELECT /* {{Description}} */ u.*

	FROM users u
	LEFT JOIN business_profiles bu ON bu.user_id = u.id
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
		WithOrderBy("u.created_at DESC").
		WithPaginationFunc(func(db, rawSQL *db.DB) (interface{}, error) {
			var records = make([]*models.User, rawSQL.RowsAffected)

			rows, err := rawSQL.Rows()
			if err != nil {
				return nil, err

			}
			defer rows.Close()

			var coordinateIDs []string

			for rows.Next() {
				var alias InquirySellerMatchingAlias
				err = db.ScanRows(rows, &alias)
				if err != nil {
					db.CustomLogger.Errorf("Scan rows error", err)
					continue
				}

				if alias.CoordinateID != "" && !helper.StringContains(coordinateIDs, alias.CoordinateID) {
					coordinateIDs = append(coordinateIDs, alias.CoordinateID)
				}

				alias.User.BusinessProfile = alias.BusinessProfile
				records = append(records, alias.User)
			}

			var wg sync.WaitGroup

			if len(coordinateIDs) > 0 {
				wg.Add(1)
				go func() {
					defer wg.Done()
					var items []*models.Coordinate
					db.Find(&items, "id IN ?", coordinateIDs)

					for _, record := range records {
						for _, item := range items {
							if record.CoordinateID == item.ID {
								record.Coordinate = item
							}
						}
					}
				}()
			}

			wg.Wait()

			return &records, nil
		})
}
