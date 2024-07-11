package queryfunc

import (
	"text/template"

	"github.com/engineeringinflow/inflow-backend/pkg/db"
	"github.com/engineeringinflow/inflow-backend/pkg/helper"
	"github.com/engineeringinflow/inflow-backend/pkg/models"
)

type CommentAlias struct {
	*models.Comment

	User *models.User `gorm:"embedded;embeddedPrefix:u__"`
}

type CommentBuilderOptions struct {
	QueryBuilderOptions
}

func NewCommentBuilder(options CommentBuilderOptions) *Builder {
	var rawSQL = `
	SELECT /* {{Description}} */ c.*,

	u.id AS u__id,
	u.name AS u__name,
	u.avatar AS u__avatar,
	u.first_name AS u__first_name,
	u.last_name AS u__last_name,
	u.email AS u__email,
	u.avatar AS u__avatar

	FROM comments c
	JOIN users u ON u.id = c.user_id
	`
	var countSQL = `
	SELECT /* {{Description}} */ 1

	FROM comments c
	JOIN users u ON u.id = c.user_id
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
		WithOrderBy("c.created_at ASC").
		WithPaginationFunc(func(db, rawSQL *db.DB) (interface{}, error) {
			var records = make([]*models.Comment, rawSQL.RowsAffected)

			rows, err := rawSQL.Rows()
			if err != nil {
				db.CustomLogger.Errorf("Scan rows error", err)
				return nil, err
			}
			defer rows.Close()

			for rows.Next() {
				var alias CommentAlias
				err = db.ScanRows(rows, &alias)
				if err != nil {
					db.CustomLogger.Errorf("Scan rows error", err)
					continue
				}
				alias.Comment.User = alias.User
				records = append(records, alias.Comment)
			}

			return &records, nil
		})
}
