package queryfunc

import (
	"text/template"

	"github.com/engineeringinflow/inflow-backend/pkg/db"
	"github.com/engineeringinflow/inflow-backend/pkg/helper"
	"github.com/engineeringinflow/inflow-backend/pkg/models"
)

type BlogCategoryAlias struct {
	*models.BlogCategory
	TotalPost int `json:"total_post,omitempty"`
}

type BlogCategoryBuilderOptions struct {
	QueryBuilderOptions
}

func NewBlogCategoryBuilder(options BlogCategoryBuilderOptions) *Builder {
	var rawSQL = `
	SELECT /* {{Description}} */ cate.*, (select count(1) from posts where category_id = cate.id) as total_post

	FROM blog_categories cate
	`
	var countSQL = `
	SELECT /* {{Description}}  */ 1

	FROM blog_categories cate
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
		WithOrderBy("cate.updated_at DESC").
		WithPaginationFunc(func(db, rawSQL *db.DB) (interface{}, error) {
			var records = make([]*models.BlogCategory, rawSQL.RowsAffected)

			rows, err := rawSQL.Rows()
			if err != nil {
				return nil, err

			}
			defer rows.Close()

			for rows.Next() {
				var alias BlogCategoryAlias
				err = db.ScanRows(rows, &alias)
				if err != nil {
					db.CustomLogger.Errorf("Scan rows error", err)
					continue
				}
				alias.BlogCategory.TotalPost = alias.TotalPost
				records = append(records, alias.BlogCategory)
			}

			return &records, nil
		})
}
