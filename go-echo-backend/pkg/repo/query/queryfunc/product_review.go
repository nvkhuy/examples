package queryfunc

import (
	"text/template"

	"github.com/engineeringinflow/inflow-backend/pkg/db"
	"github.com/engineeringinflow/inflow-backend/pkg/helper"
	"github.com/engineeringinflow/inflow-backend/pkg/models"
)

type ProductReviewAlias struct {
	*models.ProductReview
}

type ProductReviewBuilderOptions struct {
	QueryBuilderOptions
}

func NewProductReviewBuilder(options ProductReviewBuilderOptions) *Builder {
	var rawSQL = `
	SELECT /* {{Description}} */ pr.*

	FROM product_reviews pr
	`
	var countSQL = `
	SELECT /* {{Description}} */ 1

	FROM product_reviews pr
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
		WithOrderBy("pr.updated_at DESC").
		WithPaginationFunc(func(db, rawSQL *db.DB) (interface{}, error) {
			var records = make([]*models.ProductReview, rawSQL.RowsAffected)

			rows, err := rawSQL.Rows()
			if err != nil {
				return nil, err

			}
			defer rows.Close()

			var userIds []string
			for rows.Next() {
				var copy ProductReviewAlias

				err = db.ScanRows(rows, &copy)
				if err != nil {
					continue
				}

				userIds = append(userIds, copy.UserID)
				records = append(records, copy.ProductReview)
			}

			if len(userIds) > 0 {
				var users []*models.User
				err = db.Find(&users, "id IN ?", userIds).Error
				if err != nil {
					db.CustomLogger.ErrorAny(err)
					return nil, err
				}

				for _, user := range users {
					for _, record := range records {
						if record.UserID == user.ID {
							record.User = user
						}
					}
				}
			}

			var sampleAttachments = models.Attachments{
				&models.Attachment{FileKey: "product_image.png", FileURL: "https://dev-static.joininflow.io/sample/product_image.png"},
				&models.Attachment{FileKey: "product_image.png", FileURL: "https://dev-static.joininflow.io/sample/product_image.png"},
				&models.Attachment{FileKey: "product_image.png", FileURL: "https://dev-static.joininflow.io/sample/product_image.png"},
			}

			for _, record := range records {
				record.Attachments = &sampleAttachments
			}

			return &records, nil
		})
}
