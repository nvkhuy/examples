package queryfunc

import (
	"sync"
	"text/template"

	"github.com/engineeringinflow/inflow-backend/pkg/db"
	"github.com/engineeringinflow/inflow-backend/pkg/helper"
	"github.com/engineeringinflow/inflow-backend/pkg/models"
	"github.com/samber/lo"
)

type InquiryBuyerAlias struct {
	*models.Inquiry

	User *models.User `gorm:"embedded;embeddedPrefix:u__"`
}

type InquiryBuyerBuilderOptions struct {
	QueryBuilderOptions

	IncludeInquiry         bool
	IncludeInquiryBuyer    bool
	IncludePurchaseOrder   bool
	IncludeShippingAddress bool
	IncludeCollection      bool
	IncludeAuditLog        bool
}

func NewInquiryBuyerBuilder(options InquiryBuyerBuilderOptions) *Builder {
	var rawSQL = `
	SELECT /* {{Description}} */ iq.*

	FROM inquiries iq
	`
	var countSQL = `
	SELECT /* {{Description}} */ iq.*

	FROM inquiries iq
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
		WithGroupBy("iq.id").
		WithOrderBy("iq.updated_at DESC").
		WithPaginationFunc(func(db, rawSQL *db.DB) (interface{}, error) {
			var records = make([]*models.Inquiry, rawSQL.RowsAffected)

			rows, err := rawSQL.Rows()
			if err != nil {
				return nil, err

			}
			defer rows.Close()

			var orderGroupIDs []string
			var inquiryIDs []string
			var userIDs []string

			for rows.Next() {
				var alias InquiryBuyerAlias
				err = db.ScanRows(rows, &alias)
				if err != nil {
					db.CustomLogger.Errorf("Scan rows error", err)
					continue
				}

				if alias.OrderGroupID != "" && !helper.StringContains(orderGroupIDs, alias.CollectionID) {
					orderGroupIDs = append(orderGroupIDs, alias.OrderGroupID)
				}

				if !helper.StringContains(inquiryIDs, alias.ID) {
					inquiryIDs = append(inquiryIDs, alias.ID)
				}

				if alias.UserID != "" && !helper.StringContains(userIDs, alias.UserID) {
					userIDs = append(userIDs, alias.UserID)
				}

				alias.Inquiry.QuotedPrice = alias.Inquiry.GetQuotedPrice().ToPtr()
				records = append(records, alias.Inquiry)
			}

			var wg sync.WaitGroup

			if len(orderGroupIDs) > 0 && options.IncludeCollection {
				wg.Add(1)

				go func() {
					defer wg.Done()
					var items []*models.OrderGroup
					db.Find(&items, "id IN ?", orderGroupIDs)

					for _, item := range items {
						for _, record := range records {
							if record.OrderGroupID == item.ID {
								record.OrderGroup = item
							}
						}
					}
				}()
			}

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
