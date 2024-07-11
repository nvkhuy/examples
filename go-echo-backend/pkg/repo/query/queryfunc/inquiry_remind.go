package queryfunc

import (
	"sync"
	"text/template"

	"github.com/engineeringinflow/inflow-backend/pkg/db"
	"github.com/engineeringinflow/inflow-backend/pkg/helper"
	"github.com/engineeringinflow/inflow-backend/pkg/models"
)

type InquiryRemindAlias struct {
	*models.Inquiry
}

type InquiryRemindBuilderOptions struct {
	QueryBuilderOptions

	IncludeAssignee bool
	IncludeUser     bool
}

func NewInquiryRemindBuilder(options InquiryRemindBuilderOptions) *Builder {
	var rawSQL = `
	SELECT /* {{Description}} */ iq.*

	FROM inquiries iq
	`
	var countSQL = `
	SELECT /* {{Description}} */ 1

	FROM FROM inquiries iq
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
		WithOrderBy("iq.created_at DESC").
		WithPaginationFunc(func(db, rawSQL *db.DB) (interface{}, error) {
			var records = make([]*models.Inquiry, rawSQL.RowsAffected)

			rows, err := rawSQL.Rows()
			if err != nil {
				return nil, err

			}
			defer rows.Close()

			var assigneeIDs []string
			var userIDs []string
			for rows.Next() {
				var alias InquiryRemindAlias
				err = db.ScanRows(rows, &alias)
				if err != nil {
					db.CustomLogger.Errorf("Scan rows error", err)
					continue
				}

				for _, v := range alias.AssigneeIDs {
					if !helper.StringContains(assigneeIDs, v) {
						assigneeIDs = append(assigneeIDs, v)
					}
				}

				if !helper.StringContains(userIDs, alias.UserID) {
					userIDs = append(userIDs, alias.UserID)
				}
				records = append(records, alias.Inquiry)
			}

			var wg sync.WaitGroup

			if len(assigneeIDs) > 0 && options.IncludeAssignee {
				wg.Add(1)

				go func() {
					defer wg.Done()
					var items []*models.User
					db.Select("ID", "Email", "Name", "Avatar").Find(&items, "id IN ?", assigneeIDs)

					for _, item := range items {
						for _, record := range records {
							if helper.StringContains(record.AssigneeIDs, item.ID) {
								record.Assignees = append(record.Assignees, item)
							}
						}
					}
				}()
			}

			if len(userIDs) > 0 && options.IncludeUser {
				wg.Add(1)

				go func() {
					defer wg.Done()
					var items []*models.User
					db.Select("ID", "Email", "Name", "Avatar").Find(&items, "id IN ?", userIDs)

					for _, item := range items {
						for _, record := range records {
							if record.UserID == item.ID {
								record.User = item
							}
						}
					}
				}()
			}

			wg.Wait()
			return &records, nil
		})
}
