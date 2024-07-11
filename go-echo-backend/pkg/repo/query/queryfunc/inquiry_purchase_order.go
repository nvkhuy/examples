package queryfunc

import (
	"sync"
	"text/template"

	"github.com/engineeringinflow/inflow-backend/pkg/db"
	"github.com/engineeringinflow/inflow-backend/pkg/helper"
	"github.com/engineeringinflow/inflow-backend/pkg/models"
	"github.com/samber/lo"
)

type InquiryPurchaseOrderAlias struct {
	*models.Inquiry
}

type InquiryPurchaseOrderBuilderOptions struct {
	QueryBuilderOptions

	IncludeInquiryBuyer    bool
	IncludePurchaseOrder   bool
	IncludeShippingAddress bool
	IncludeCollection      bool
	IncludeAuditLog        bool
}

func NewInquiryPurchaseOrderBuilder(options InquiryPurchaseOrderBuilderOptions) *Builder {
	var rawSQL = `
	SELECT /* {{Description}} */ iq.*

	FROM inquiries iq
	LEFT JOIN purchase_orders po ON iq.id = po.inquiry_id
	`
	var countSQL = `
	SELECT /* {{Description}} */ 1

	FROM inquiries iq
	LEFT JOIN purchase_orders po ON iq.id = po.inquiry_id
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
		WithOrderBy("iq.created_at DESC").
		WithPaginationFunc(func(db, rawSQL *db.DB) (interface{}, error) {
			var records = make([]*models.Inquiry, rawSQL.RowsAffected)

			rows, err := rawSQL.Rows()
			if err != nil {
				return nil, err

			}
			defer rows.Close()

			var inquiryIDs []string
			var userIDs []string

			for rows.Next() {
				var alias InquiryAlias
				err = db.ScanRows(rows, &alias)
				if err != nil {
					db.CustomLogger.Errorf("Scan rows error", err)
					continue
				}

				if !helper.StringContains(inquiryIDs, alias.ID) {
					inquiryIDs = append(inquiryIDs, alias.ID)
				}

				if alias.UserID != "" && !helper.StringContains(userIDs, alias.UserID) {
					userIDs = append(userIDs, alias.UserID)
				}

				records = append(records, alias.Inquiry)
			}

			var wg sync.WaitGroup

			if len(inquiryIDs) > 0 && options.IncludePurchaseOrder {
				wg.Add(1)

				go func() {
					defer wg.Done()

					var items []*models.PurchaseOrder
					db.Find(&items, "inquiry_id IN ? AND user_id IN ?", inquiryIDs, userIDs)

					var cartItemsIDs []string
					for _, item := range items {
						for _, record := range records {
							if record.ID == item.InquiryID && record.UserID == item.UserID {
								record.PurchaseOrder = item

								cartItemsIDs = append(cartItemsIDs, item.CartItemIDs...)
							}
						}
					}

					if len(cartItemsIDs) > 0 {
						var cartItems []*models.InquiryCartItem
						db.Find(&cartItems, "id IN ? AND inquiry_id IN ?", cartItemsIDs, inquiryIDs)

						for _, cardItem := range cartItems {
							for _, record := range records {
								if record.PurchaseOrder != nil && lo.Contains(record.PurchaseOrder.CartItemIDs, cardItem.ID) {
									record.PurchaseOrder.CartItems = append(record.PurchaseOrder.CartItems, cardItem)
								}
							}
						}
					}
				}()

			}

			wg.Wait()
			return &records, nil
		})
}
