package queryfunc

import (
	"sync"
	"text/template"

	"github.com/lib/pq"

	"gorm.io/plugin/dbresolver"

	"github.com/engineeringinflow/inflow-backend/pkg/db"
	"github.com/engineeringinflow/inflow-backend/pkg/helper"
	"github.com/engineeringinflow/inflow-backend/pkg/models"
	"github.com/engineeringinflow/inflow-backend/pkg/models/enums"
)

type UserAlias struct {
	*models.User
}

type UserBuilderOptions struct {
	QueryBuilderOptions

	IncludeAddress            bool
	IncludeBusinessProfile    bool
	IncludeContactOwners      bool
	IncludeBrandTeam          bool
	IsConsistentRead          bool
	IncludeAssignedInquiryIds bool
}

func NewUserBuilder(options UserBuilderOptions) *Builder {
	var rawSQL = `
	SELECT /* {{Description}} */ u.*

	FROM users u
	`
	var countSQL = `
	SELECT /* {{Description}} */ 1

	FROM users u
	`

	if options.IncludeAssignedInquiryIds {
		rawSQL = `
		WITH all_assignees AS (
			SELECT assignee_ids FROM inquiries WHERE assignee_ids is not null
			UNION ALL
			SELECT assignee_ids FROM purchase_orders WHERE assignee_ids is not null
			UNION ALL
			SELECT assignee_ids FROM bulk_purchase_orders WHERE assignee_ids is not null
		)
		SELECT /* {{Description}} */ u.*, COUNT(1) AS total_order_assigned
		FROM users u
		LEFT JOIN all_assignees
		ON u.id = ANY(all_assignees.assignee_ids)
		`
	}

	b := NewBuilder(rawSQL, countSQL).
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
			var userIDs []string
			var contactOwnerIDs []string

			for rows.Next() {
				var alias UserAlias
				err = db.ScanRows(rows, &alias)
				if err != nil {
					db.CustomLogger.Errorf("Scan rows error", err)
					continue
				}

				if alias.CoordinateID != "" && !helper.StringContains(coordinateIDs, alias.CoordinateID) {
					coordinateIDs = append(coordinateIDs, alias.CoordinateID)
				}

				if !helper.StringContains(userIDs, alias.ID) {
					userIDs = append(userIDs, alias.ID)
				}

				for _, v := range alias.ContactOwnerIDs {
					if !helper.StringContains(contactOwnerIDs, v) {
						contactOwnerIDs = append(contactOwnerIDs, v)
					}
				}

				records = append(records, alias.User)
			}

			var wg sync.WaitGroup

			if len(coordinateIDs) > 0 && options.IncludeAddress {
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

			if len(contactOwnerIDs) > 0 && options.IncludeContactOwners {
				wg.Add(1)
				go func() {
					defer wg.Done()
					var items []*models.User
					db.Select("ID", "Email", "Name", "Avatar").Find(&items, "id IN ?", contactOwnerIDs)

					for _, item := range items {
						for _, record := range records {
							if helper.StringContains(record.ContactOwnerIDs, item.ID) {
								record.ContactOwners = append(record.ContactOwners, item)
							}
						}
					}
				}()
			}

			if len(userIDs) > 0 && options.IncludeBusinessProfile {
				wg.Add(1)

				go func() {
					defer wg.Done()
					var items []*models.BusinessProfile
					db.Find(&items, "user_id IN ?", userIDs)

					for _, record := range records {
						for _, item := range items {
							if record.ID == item.UserID {
								record.BusinessProfile = item
							}
						}
					}
				}()
			}

			if len(userIDs) > 0 && options.IncludeBrandTeam {
				wg.Add(1)

				go func() {
					defer wg.Done()
					var items []*models.BrandTeam
					db.Find(&items, "user_id IN ?", userIDs)

					for _, record := range records {
						record.BrandTeam = &models.BrandTeam{
							Role: enums.BrandTeamRoleManager,
						}

						for _, item := range items {
							if record.ID == item.UserID {
								record.BrandTeam = item
							}
						}
					}
				}()
			}

			if len(userIDs) > 0 && options.IncludeAssignedInquiryIds {
				wg.Add(1)

				go func() {
					defer wg.Done()
					var iqs []*models.Inquiry
					db.Find(&iqs, "count_elements(assignee_ids,?) >= 1", pq.StringArray(userIDs))

					for _, record := range records {
						for _, iq := range iqs {
							for _, assigneeID := range iq.AssigneeIDs {
								if assigneeID == record.ID {
									record.AssignedInquiryIDs = append(record.AssignedInquiryIDs, iq.ID)
								}
							}
						}
					}
				}()
			}

			if len(userIDs) > 0 && options.IncludeAssignedInquiryIds {
				wg.Add(1)

				go func() {
					defer wg.Done()
					var pos []*models.PurchaseOrder
					db.Find(&pos, "count_elements(assignee_ids,?) >= 1", pq.StringArray(userIDs))

					for _, record := range records {
						for _, po := range pos {
							for _, assigneeID := range po.AssigneeIDs {
								if assigneeID == record.ID {
									record.AssignedPOIDs = append(record.AssignedPOIDs, po.ID)
								}
							}
						}
					}
				}()
			}

			if len(userIDs) > 0 && options.IncludeAssignedInquiryIds {
				wg.Add(1)

				go func() {
					defer wg.Done()
					var bulkPOs []*models.BulkPurchaseOrder
					db.Find(&bulkPOs, "count_elements(assignee_ids,?) >= 1", pq.StringArray(userIDs))

					for _, record := range records {
						for _, bulkPO := range bulkPOs {
							for _, assigneeID := range bulkPO.AssigneeIDs {
								if assigneeID == record.ID {
									record.AssignedBulkPOIDs = append(record.AssignedBulkPOIDs, bulkPO.ID)
								}
							}
						}
					}
				}()
			}

			wg.Wait()

			return &records, nil
		})

	if options.IsConsistentRead {
		b.WithClauses(dbresolver.Write)
	}

	if options.IncludeAssignedInquiryIds {
		b.WithGroupBy("u.id")
		b.WithOrderBy("total_order_assigned DESC")
	} else {
		b.WithOrderBy("u.created_at DESC")
	}
	return b
}
