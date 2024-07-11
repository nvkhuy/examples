package queryfunc

import (
	"encoding/json"
	"sync"
	"text/template"

	"github.com/engineeringinflow/inflow-backend/pkg/db"
	"github.com/engineeringinflow/inflow-backend/pkg/helper"
	"github.com/engineeringinflow/inflow-backend/pkg/models"
	"gorm.io/datatypes"
)

type DocumentAlias struct {
	models.Document
	DocumentTagsJson datatypes.JSON `gorm:"column:document_tags_json"`
}

type DocumentBuilderOptions struct {
	QueryBuilderOptions
}

func NewDocumentBuilder(options DocumentBuilderOptions) *Builder {
	var rawSQL = `
	SELECT /* {{Description}} */ d.*,
	json_agg(json_build_object(
		'id',dt.id,
		'name',dt.name
	)) filter (where dt.id is not null)  AS document_tags_json

	FROM documents d
	LEFT JOIN tagged_documents td ON td.document_id = d.id
	LEFT JOIN document_tags dt ON dt.id = td.document_tag_id
	`

	var countSQL = `
	SELECT /* {{Description}} */ 1

	FROM documents d
	LEFT JOIN tagged_documents td ON td.document_id = d.id
	LEFT JOIN document_tags dt ON dt.id = td.document_tag_id

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
		WithGroupBy("d.id").
		WithOrderBy("d.updated_at DESC").
		WithPaginationFunc(func(db, rawSQL *db.DB) (interface{}, error) {
			var records = make([]*models.Document, rawSQL.RowsAffected)

			rows, err := rawSQL.Rows()
			if err != nil {
				db.CustomLogger.Errorf("Scan rows error", err)
				return nil, err
			}
			defer rows.Close()

			var userIds []string
			var cateIds []string

			for rows.Next() {
				var alias DocumentAlias
				err = db.ScanRows(rows, &alias)
				if err != nil {
					continue
				}
				if alias.DocumentTagsJson != nil {
					if err := json.Unmarshal(alias.DocumentTagsJson, &alias.Document.Tags); err != nil {
						continue
					}
				}
				userIds = append(userIds, alias.UserID)
				cateIds = append(cateIds, alias.CategoryID)

				records = append(records, &alias.Document)
			}
			var wg sync.WaitGroup
			wg.Add(2)
			// users
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
							if record.UserID == user.ID {
								record.User = user
							}
						}
					}

					for _, record := range records {
						record.UserID = ""
						if record.User != nil {
							record.User.ID = ""
						}
					}
				}
			}()
			// document category
			go func() {
				defer wg.Done()
				if len(cateIds) > 0 {
					var cates []*models.DocumentCategory
					err = db.Find(&cates, "id IN ?", cateIds).Error
					if err != nil {
						return
					}

					for _, cate := range cates {
						for _, record := range records {
							if record.CategoryID == cate.ID {
								record.Category = cate
							}
						}
					}
				}
			}()
			wg.Wait()

			return &records, nil
		})
}
