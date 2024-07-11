package queryfunc

import (
	"sync"
	"text/template"

	"github.com/engineeringinflow/inflow-backend/pkg/db"
	"github.com/engineeringinflow/inflow-backend/pkg/helper"
	"github.com/engineeringinflow/inflow-backend/pkg/models"
)

type InquiryCartAlias struct {
	*models.Inquiry
}

type InquiryCartBuilderOptions struct {
	QueryBuilderOptions
}

func NewInquiryCartBuilder(options InquiryCartBuilderOptions) *Builder {
	var rawSQL = `
	SELECT /* {{Description}} */ iq.*

	FROM inquiries iq
	INNER JOIN inquiry_cart_items c ON iq.id = c.inquiry_id
	LEFT JOIN purchase_orders po ON po.inquiry_id = iq.id

	`
	var countSQL = `
	SELECT /* {{Description}} */ 1

	FROM inquiries iq
	INNER JOIN inquiry_cart_items c ON iq.id = c.inquiry_id
	LEFT JOIN purchase_orders po ON po.inquiry_id = iq.id
	`

	if options.Role.IsAdmin() {
		rawSQL = `
		SELECT /* {{Description}} */ iq.*

		FROM inquiries iq
		INNER JOIN inquiry_cart_items c ON iq.id = c.inquiry_id
		LEFT JOIN purchase_orders po ON po.inquiry_id = iq.id

		`
		countSQL = `
		SELECT /* {{Description}} */ 1

		FROM inquiries iq
		INNER JOIN inquiry_cart_items c ON iq.id = c.inquiry_id
		LEFT JOIN purchase_orders po ON po.inquiry_id = iq.id
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
		WithOrderBy("iq.updated_at DESC").
		WithGroupBy("iq.id").
		WithPaginationFunc(func(db, rawSQL *db.DB) (interface{}, error) {
			var records = make([]*models.Inquiry, rawSQL.RowsAffected)

			rows, err := rawSQL.Rows()
			if err != nil {
				return nil, err

			}
			defer rows.Close()
			var inquiryIDs []string

			for rows.Next() {
				var alias InquiryCartAlias
				err = db.ScanRows(rows, &alias)
				if err != nil {
					db.CustomLogger.Errorf("Scan rows error", err)
					continue
				}

				if !helper.StringContains(inquiryIDs, alias.ID) {
					inquiryIDs = append(inquiryIDs, alias.ID)
				}

				alias.Inquiry.QuotedPrice = alias.Inquiry.GetQuotedPrice().ToPtr()
				records = append(records, alias.Inquiry)
			}

			var wg sync.WaitGroup

			if len(inquiryIDs) > 0 {
				wg.Add(1)
				go func() {
					defer wg.Done()
					var items []*models.InquiryCartItem
					db.Find(&items, "inquiry_id IN ?", inquiryIDs)

					for _, cartItem := range items {
						for _, record := range records {
							if record.ID == cartItem.InquiryID {
								record.CartItems = append(record.CartItems, cartItem)
							}
						}
					}
				}()

				wg.Add(1)
				go func() {
					defer wg.Done()
					var items []*models.PurchaseOrder
					db.Select("ID", "InquiryID", "SubTotal", "ShippingFee", "TransactionFee", "Tax", "TotalPrice", "TaxPercentage").Find(&items, "inquiry_id IN ?", inquiryIDs)

					for _, item := range items {
						for _, record := range records {
							if record.ID == item.InquiryID {
								record.PurchaseOrder = item
							}
						}
					}
				}()
			}

			wg.Wait()
			return &records, nil
		})
}
