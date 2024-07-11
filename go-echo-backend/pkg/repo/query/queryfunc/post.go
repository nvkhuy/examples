package queryfunc

import (
	"sync"
	"text/template"

	"github.com/engineeringinflow/inflow-backend/pkg/db"
	"github.com/engineeringinflow/inflow-backend/pkg/helper"
	"github.com/engineeringinflow/inflow-backend/pkg/models"
)

type PostAlias struct {
	*models.Post
	User *models.User `gorm:"embedded;embeddedPrefix:u__"`
}

type PostBuilderOptions struct {
	QueryBuilderOptions
}

func NewPostBuilder(options PostBuilderOptions) *Builder {
	var rawSQL = `
	SELECT /* {{Description}} */ 
	p.id,
	p.created_at,
	p.updated_at,
	p.title,
	p.short_description,
	p.slug,
	p.published_at,
	p.status,
	p.featured_image,
	p.category_id,
	p.user_id,
	p.vi,
	p.content,
	p.content_url,
	p.setting_seo_id

	FROM posts p
	`

	var countSQL = `
	SELECT /* {{Description}} */ 1

	FROM posts p
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
		WithOrderBy("p.created_at DESC").
		WithPaginationFunc(func(db, rawSQL *db.DB) (interface{}, error) {
			var records = make([]*models.Post, rawSQL.RowsAffected)

			rows, err := rawSQL.Rows()
			if err != nil {
				db.CustomLogger.Errorf("Scan rows error", err)
				return nil, err
			}
			defer rows.Close()

			var userIds []string
			var cateIds []string
			var settingSeoIDs []string

			for rows.Next() {
				var copy PostAlias
				err = db.ScanRows(rows, &copy)
				if err != nil {
					continue
				}

				userIds = append(userIds, copy.UserID)
				cateIds = append(cateIds, copy.CategoryID)
				settingSeoIDs = append(settingSeoIDs, copy.SettingSeoID)
				if copy.VI != nil {
					settingSeoIDs = append(settingSeoIDs, copy.VI.SettingSeoID)
				}
				records = append(records, copy.Post)
			}
			var wg sync.WaitGroup
			wg.Add(3)

			// match route
			go func() {
				defer wg.Done()
				if len(settingSeoIDs) > 0 {
					var seos []*models.SettingSEO
					err = db.Find(&seos, "id IN ?", settingSeoIDs).Error
					if err != nil {
						return
					}
					for _, seo := range seos {
						for _, record := range records {
							if record == nil {
								continue
							}
							if record.SettingSeoID != "" && seo.ID == record.SettingSeoID {
								record.SettingSEO = seo
							}
							if record.VI != nil && record.VI.SettingSeoID != "" {
								if seo.ID == record.VI.SettingSeoID {
									record.VI.SettingSEO = seo
								}
							}
						}
					}
				}
			}()
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
			// blog category
			go func() {
				defer wg.Done()
				if len(cateIds) > 0 {
					var cates []*models.BlogCategory
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
