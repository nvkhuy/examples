package queryfunc

import (
	"sync"
	"text/template"

	"github.com/engineeringinflow/inflow-backend/pkg/db"
	"github.com/engineeringinflow/inflow-backend/pkg/helper"
	"github.com/engineeringinflow/inflow-backend/pkg/models"
)

type BulkPurchaseOrderSellerQuotationAlias struct {
	*models.BulkPurchaseOrderSellerQuotation

	Seller          *models.User            `gorm:"embedded;embeddedPrefix:u__"`
	BusinessProfile *models.BusinessProfile `gorm:"embedded;embeddedPrefix:bu__"`
}

type BulkPurchaseOrderSellerQuotationBuilderOptions struct {
	QueryBuilderOptions

	IncludeBulk   bool
	CurrentUserID string
}

func NewBulkPurchaseOrderSellerQuotationBuilder(options BulkPurchaseOrderSellerQuotationBuilderOptions) *Builder {
	var rawSQL = `
	SELECT /* {{Description}} */ rq.*

	FROM bulk_purchase_order_seller_quotations rq
	`
	var countSQL = `
	SELECT /* {{Description}} */ 1

	FROM bulk_purchase_order_seller_quotations rq
	`

	if options.Role.IsAdmin() {
		rawSQL = `
		SELECT /* {{Description}} */ rq.*, 
		rq.id AS bulk_quotation_id,
		u.id AS u__id,
		u.name AS u__name,
		u.company_name AS u__company_name,
		u.payment_terms AS u__payment_terms,
		u.coordinate_id AS u__coordinate_id,
		u.country_code AS u__country_code,
		u.phone_number AS u__phone_number,
		bu.order_types AS bu__order_types

		FROM bulk_purchase_order_seller_quotations rq
		JOIN users u ON u.id = rq.user_id
		LEFT JOIN business_profiles bu ON bu.user_id = u.id
		`

		countSQL = `
		SELECT /* {{Description}} */ 1

		FROM bulk_purchase_order_seller_quotations rq
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
			var records = make([]*models.BulkPurchaseOrderSellerQuotation, rawSQL.RowsAffected)

			rows, err := rawSQL.Rows()
			if err != nil {
				return nil, err

			}
			defer rows.Close()

			var bulkPurchaseOrderIDs []string
			var coordinateIDs []string

			for rows.Next() {
				var alias BulkPurchaseOrderSellerQuotationAlias
				err = db.ScanRows(rows, &alias)
				if err != nil {
					db.CustomLogger.Errorf("Scan rows error", err)
					continue
				}

				if !helper.StringContains(bulkPurchaseOrderIDs, alias.BulkPurchaseOrderID) {
					bulkPurchaseOrderIDs = append(bulkPurchaseOrderIDs, alias.BulkPurchaseOrderID)
				}

				if alias.Seller != nil {
					alias.BulkPurchaseOrderSellerQuotation.User = alias.Seller
				}
				if alias.BulkPurchaseOrderSellerQuotation.User != nil && alias.BusinessProfile != nil {
					alias.BulkPurchaseOrderSellerQuotation.User.BusinessProfile = alias.BusinessProfile
				}

				if alias.Seller != nil && alias.Seller.CoordinateID != "" && !helper.StringContains(coordinateIDs, alias.Seller.CoordinateID) {
					coordinateIDs = append(coordinateIDs, alias.Seller.CoordinateID)
				}

				records = append(records, alias.BulkPurchaseOrderSellerQuotation)
			}

			var wg sync.WaitGroup

			if len(bulkPurchaseOrderIDs) > 0 && options.IncludeBulk {
				wg.Add(1)
				go func() {
					defer wg.Done()
					var bulks []*models.BulkPurchaseOrder
					db.Find(&bulks, "id IN ?", bulkPurchaseOrderIDs)
					for _, bulk := range bulks {
						for _, record := range records {
							if record.BulkPurchaseOrderID == bulk.ID {
								record.BulkPurchaseOrder = bulk
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

			wg.Wait()
			return &records, nil
		})
}
