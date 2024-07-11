package queryfunc

import (
	"text/template"

	"github.com/engineeringinflow/inflow-backend/pkg/db"
	"github.com/engineeringinflow/inflow-backend/pkg/helper"
	"github.com/engineeringinflow/inflow-backend/pkg/models"
)

type InquirySellerAllocationAlias struct {
	*models.User

	InquirySeller *models.InquirySeller `gorm:"embedded;embeddedPrefix:iqs__" json:"inquiry_seller,omitempty"`
}

type InquirySellerAllocationBuilderOptions struct {
	QueryBuilderOptions
}

func NewInquirySellerAllocationBuilder(options InquirySellerAllocationBuilderOptions) *Builder {
	var rawSQL = `
	SELECT /* {{Description}} */ u.*,
	iqs.id AS iqs__id,
	iqs.inquiry_id AS iqs__inquiry_id,
	iqs.order_type AS iqs__order_type,
	iqs.status AS iqs__status,
	iqs.bulk_quotations AS iqs__bulk_quotations,
	iqs.fabric_cost AS iqs__fabric_cost,
	iqs.decoration_cost AS iqs__decoration_cost,
	iqs.making_cost AS iqs__making_cost,
	iqs.other_cost AS iqs__other_cost,
	iqs.sample_unit_price AS iqs__sample_unit_price,
	iqs.sample_lead_time AS iqs__sample_lead_time


	FROM users u
	LEFT JOIN inquiry_sellers iqs ON iqs.user_id = u.id
	`
	var countSQL = `
	SELECT /* {{Description}} */ u.*

	FROM users u
	LEFT JOIN inquiry_sellers iqs ON iqs.user_id = u.id
	`

	if options.Role.IsAdmin() {
		rawSQL = `
		SELECT /* {{Description}} */ u.*,
		iqs.id AS iqs__id,
		iqs.inquiry_id AS iqs__inquiry_id,
		iqs.order_type AS iqs__order_type,
		iqs.status AS iqs__status,
		iqs.bulk_quotations AS iqs__bulk_quotations,
		iqs.fabric_cost AS iqs__fabric_cost,
		iqs.decoration_cost AS iqs__decoration_cost,
		iqs.making_cost AS iqs__making_cost,
		iqs.other_cost AS iqs__other_cost,
		iqs.sample_unit_price AS iqs__sample_unit_price,
		iqs.sample_lead_time AS iqs__sample_lead_time

		FROM users u
		LEFT JOIN inquiry_sellers iqs ON iqs.user_id = u.id
		`

		countSQL = `
		SELECT /* {{Description}} */ u.*

		FROM users u
		LEFT JOIN inquiry_sellers iqs ON iqs.user_id = u.id
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
		WithOrderBy("u.created_at DESC").
		WithPaginationFunc(func(db, rawSQL *db.DB) (interface{}, error) {
			var records = make([]*models.User, rawSQL.RowsAffected)

			rows, err := rawSQL.Rows()
			if err != nil {
				return nil, err

			}
			defer rows.Close()

			for rows.Next() {
				var alias InquirySellerAllocationAlias
				err = db.ScanRows(rows, &alias)
				if err != nil {
					db.CustomLogger.Errorf("Scan rows error", err)
					continue
				}

				if alias.InquirySeller != nil && alias.InquirySeller.ID != "" {
					alias.User.InquirySeller = alias.InquirySeller
				}

				records = append(records, alias.User)
			}

			return &records, nil
		})
}
