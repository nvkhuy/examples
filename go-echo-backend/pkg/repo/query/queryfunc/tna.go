package queryfunc

import (
	"sync"
	"text/template"

	"github.com/engineeringinflow/inflow-backend/pkg/db"
	"github.com/engineeringinflow/inflow-backend/pkg/helper"
	"github.com/engineeringinflow/inflow-backend/pkg/models"
)

type TNAAlias struct {
	*models.TNA
}

type TNABuilderOptions struct {
	QueryBuilderOptions
}

func NewTNABuilder(options TNABuilderOptions) *Builder {
	var rawSQL = `
	SELECT /* {{Description}} */ t.*

	FROM tnas t
	`
	var countSQL = `
	SELECT /* {{Description}} */ 1

	FROM tnas t
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
		WithOrderBy("t.created_at DESC").
		WithPaginationFunc(func(db, rawSQL *db.DB) (interface{}, error) {
			var records = make([]*models.TNA, rawSQL.RowsAffected)

			rows, err := rawSQL.Rows()
			if err != nil {
				return nil, err

			}
			defer rows.Close()

			var (
				assigneeIds  []string
				referenceIds []string
			)
			for rows.Next() {
				var alias TNAAlias
				err = db.ScanRows(rows, &alias)
				if err != nil {
					db.CustomLogger.Errorf("Scan rows error", err)
					continue
				}
				assigneeIds = append(assigneeIds, alias.TNA.AssigneeIDs...)
				referenceIds = append(referenceIds, alias.TNA.ReferenceID)
				records = append(records, alias.TNA)
			}

			var wg sync.WaitGroup

			if len(assigneeIds) > 0 {
				wg.Add(1)

				go func() {
					defer wg.Done()
					var assignees []models.User
					db.Model(&models.User{}).Where("id IN ?", assigneeIds).Find(&assignees)
					for _, record := range records {
						for _, assigneeId := range record.AssigneeIDs {
							for _, user := range assignees {
								if assigneeId == user.ID {
									record.Assignees = append(record.Assignees, user)
								}
							}
						}
					}
				}()

			}
			if len(referenceIds) > 0 {
				var (
					iqReferenceIds  []string
					poReferenceIds  []string
					bpoReferenceIds []string
				)
				for _, ref := range referenceIds {
					switch ref[:2] {
					case "IQ":
						iqReferenceIds = append(iqReferenceIds, ref)
					case "PO":
						poReferenceIds = append(poReferenceIds, ref)
					case "BP":
						bpoReferenceIds = append(bpoReferenceIds, ref)
					}
				}

				if len(iqReferenceIds) > 0 {
					wg.Add(1)

					go func() {
						defer wg.Done()
						var iqs []models.Inquiry
						db.Model(&models.Inquiry{}).Where("reference_id IN ?", iqReferenceIds).Find(&iqs)
						for _, record := range records {
							for _, iq := range iqs {
								if iq.ReferenceID == record.ReferenceID {
									record.Inquiry = &iq
									record.UserID = iq.UserID
								}
							}
						}
					}()

				}
				if len(poReferenceIds) > 0 {
					wg.Add(1)
					go func() {
						defer wg.Done()
						var pos []models.PurchaseOrder
						db.Model(&models.PurchaseOrder{}).Where("reference_id IN ?", poReferenceIds).Find(&pos)
						for _, record := range records {
							for _, po := range pos {
								if po.ReferenceID == record.ReferenceID {
									record.PurchaseOrder = &po
									record.UserID = po.UserID
								}
							}
						}
					}()

				}

				if len(bpoReferenceIds) > 0 {
					wg.Add(1)

					go func() {
						defer wg.Done()
						var bpos []models.BulkPurchaseOrder
						db.Model(&models.BulkPurchaseOrder{}).Where("reference_id IN ?", bpoReferenceIds).Find(&bpos)
						for _, record := range records {
							for _, bpo := range bpos {
								if bpo.ReferenceID == record.ReferenceID {
									record.BulkPurchaseOrder = &bpo
									record.UserID = bpo.UserID
								}
							}
						}
					}()

				}
			}

			wg.Wait()

			return &records, nil
		})
}
