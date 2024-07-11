package queryfunc

import (
	"sync"
	"text/template"

	"github.com/engineeringinflow/inflow-backend/pkg/db"
	"github.com/engineeringinflow/inflow-backend/pkg/helper"
	"github.com/engineeringinflow/inflow-backend/pkg/models"
)

type ChatMessageBuilderOptions struct {
	QueryBuilderOptions
}

func NewChatMessageBuilder(options ChatMessageBuilderOptions) *Builder {
	var rawSQL = `
	SELECT /* {{Description}} */ 
	c.id,
	c.created_at,
	c.updated_at,
	c.receiver_id,
	c.sender_id,
	c.message,
	c.message_type,
	c.attachments,
	c.seen_at
	
	FROM chat_messages c
	`

	var countSQL = `
	SELECT /* {{Description}} */ 1

	FROM chat_messages c
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
		WithOrderBy("c.created_at DESC").
		WithPaginationFunc(func(db, rawSQL *db.DB) (interface{}, error) {
			var records = make([]*models.ChatMessage, rawSQL.RowsAffected)

			rows, err := rawSQL.Rows()
			if err != nil {
				db.CustomLogger.Errorf("Scan rows error", err)
				return nil, err
			}
			defer rows.Close()

			var userIds []string

			for rows.Next() {
				var copy models.ChatMessage
				err = db.ScanRows(rows, &copy)
				if err != nil {
					continue
				}

				userIds = append(userIds, copy.SenderID)
				records = append(records, &copy)
			}
			var wg sync.WaitGroup
			wg.Add(1)
			go func() {
				defer wg.Done()
				if len(userIds) > 0 {
					var users []*models.User
					err = db.Select("ID", "Name", "Avatar").Find(&users, "id IN ?", userIds).Error
					if err != nil {
						return
					}

					for _, user := range users {
						for _, record := range records {
							if record.SenderID == user.ID {
								record.Sender = user
							}
						}
					}
				}
			}()
			wg.Wait()

			return &records, nil
		})
}
