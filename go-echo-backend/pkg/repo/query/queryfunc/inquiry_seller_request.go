package queryfunc

import (
	"github.com/engineeringinflow/inflow-backend/pkg/repo/query"
	"sync"
	"text/template"

	"github.com/engineeringinflow/inflow-backend/pkg/db"
	"github.com/engineeringinflow/inflow-backend/pkg/helper"
	"github.com/engineeringinflow/inflow-backend/pkg/models"
	"github.com/engineeringinflow/inflow-backend/pkg/models/enums"
	"github.com/thaitanloi365/go-utils/values"
)

type InquirySellerRequestAlias struct {
	*models.InquirySeller

	Seller          *models.User            `gorm:"embedded;embeddedPrefix:u__"`
	BusinessProfile *models.BusinessProfile `gorm:"embedded;embeddedPrefix:bu__"`
}

type InquirySellerRequestBuilderOptions struct {
	QueryBuilderOptions

	IncludeInquiry            bool
	IncludeUnseenCommentCount bool
	CurrentUserID             string
}

func NewInquirySellerRequestBuilder(options InquirySellerRequestBuilderOptions) *Builder {
	var rawSQL = `
	SELECT /* {{Description}} */ rq.*

	FROM inquiry_sellers rq
	JOIN inquiries iq ON iq.id = rq.inquiry_id
	`
	var countSQL = `
	SELECT /* {{Description}} */ rq.*

	FROM inquiry_sellers rq
	JOIN inquiries iq ON iq.id = rq.inquiry_id
	`

	if options.Role.IsAdmin() {
		rawSQL = `
		SELECT /* {{Description}} */ rq.*,
		u.id AS u__id,
		u.name AS u__name,
		u.company_name AS u__company_name,
		u.payment_terms AS u__payment_terms,
		u.coordinate_id AS u__coordinate_id,
		u.country_code AS u__country_code,
		u.phone_number AS u__phone_number,
		bu.order_types AS bu__order_types

		FROM inquiry_sellers rq
		JOIN inquiries iq ON iq.id = rq.inquiry_id
		JOIN users u ON u.id = rq.user_id
		LEFT JOIN business_profiles bu ON bu.user_id = u.id
		`

		countSQL = `
		SELECT /* {{Description}} */ 1

		FROM inquiry_sellers rq
		JOIN inquiries iq ON iq.id = rq.inquiry_id
		JOIN users u ON u.id = rq.user_id
		LEFT JOIN business_profiles bu ON bu.user_id = u.id
		`
	}

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
		WithOrderBy("rq.created_at DESC").
		WithPaginationFunc(func(db, rawSQL *db.DB) (interface{}, error) {
			var records = make([]*models.InquirySeller, rawSQL.RowsAffected)

			rows, err := rawSQL.Rows()
			if err != nil {
				return nil, err

			}
			defer rows.Close()

			var inquiryIDs []string
			var coordinateIDs []string

			for rows.Next() {
				var alias InquirySellerRequestAlias
				err = db.ScanRows(rows, &alias)
				if err != nil {
					db.CustomLogger.Errorf("Scan rows error", err)
					continue
				}

				if !helper.StringContains(inquiryIDs, alias.InquiryID) {
					inquiryIDs = append(inquiryIDs, alias.InquiryID)
				}

				if alias.Seller != nil {
					alias.InquirySeller.User = alias.Seller
				}
				if alias.InquirySeller.User != nil && alias.BusinessProfile != nil {
					alias.InquirySeller.User.BusinessProfile = alias.BusinessProfile
				}

				if alias.Seller != nil && alias.Seller.CoordinateID != "" && !helper.StringContains(coordinateIDs, alias.Seller.CoordinateID) {
					coordinateIDs = append(coordinateIDs, alias.Seller.CoordinateID)
				}

				records = append(records, alias.InquirySeller)
			}

			var wg sync.WaitGroup

			if len(inquiryIDs) > 0 && options.IncludeInquiry {
				wg.Add(1)
				go func() {
					defer wg.Done()
					var inquiries []*models.Inquiry
					var builder = NewInquiryBuilder(InquiryBuilderOptions{
						IncludeAssignee: true,
						QueryBuilderOptions: QueryBuilderOptions{
							Role: options.Role,
						},
					})
					err = query.New(db, builder).WhereFunc(func(builder *query.Builder) {
						builder.Where("iq.id IN ?", inquiryIDs)
					}).FindFunc(&inquiries)
					if err != nil {
						return
					}
					for _, inquiry := range inquiries {
						for _, record := range records {
							if record.InquiryID == inquiry.ID {
								record.Inquiry = inquiry
							}
						}
					}
				}()
			}

			if len(coordinateIDs) > 0 {
				wg.Add(1)
				go func() {
					defer wg.Done()
					var items []*models.Coordinate
					db.Find(&items, "id IN ?", coordinateIDs)

					for _, record := range records {
						for _, item := range items {
							if record.User.CoordinateID == item.ID {
								record.User.Coordinate = item
							}
						}
					}
				}()
			}

			if options.CurrentUserID != "" && options.IncludeUnseenCommentCount {
				wg.Add(1)
				go func() {
					defer wg.Done()

					for _, record := range records {
						var unseenCount int64
						db.Model(&models.Comment{}).Where("user_id != ? AND target_type = ? AND target_id = ? AND seen_at IS NULL", options.CurrentUserID, enums.CommentTargetTypeInquirySellerRequest, record.ID).Count(&unseenCount)
						record.UnseenCommentCount = values.Int64(unseenCount)
					}
				}()
			}

			wg.Wait()
			return &records, nil
		})
}
