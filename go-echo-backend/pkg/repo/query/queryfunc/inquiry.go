package queryfunc

import (
	"sync"
	"text/template"

	"github.com/engineeringinflow/inflow-backend/pkg/db"
	"github.com/engineeringinflow/inflow-backend/pkg/helper"
	"github.com/engineeringinflow/inflow-backend/pkg/models"
	"github.com/engineeringinflow/inflow-backend/pkg/models/enums"
	"github.com/engineeringinflow/inflow-backend/pkg/repo/query"
	"github.com/lib/pq"
	"github.com/samber/lo"
)

type InquiryAlias struct {
	*models.Inquiry
	Assignee *models.User `gorm:"embedded;embeddedPrefix:ua__"`
}

type InquiryBuilderOptions struct {
	QueryBuilderOptions

	IncludePurchaseOrder   bool
	IncludeShippingAddress bool
	IncludeCollection      bool
	IncludeAuditLog        bool
	IncludeAssignee        bool
	IncludeUser            bool
}

func NewInquiryBuilder(options InquiryBuilderOptions) *Builder {
	var rawSQL = `
	SELECT /* {{Description}} */ iq.*
	FROM inquiries iq
	`
	var countSQL = `
	SELECT /* {{Description}} */ 1

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
		WithOrderBy("iq.created_at DESC").
		WithPaginationFunc(func(db, rawSQL *db.DB) (interface{}, error) {
			var records = make([]*models.Inquiry, rawSQL.RowsAffected)

			rows, err := rawSQL.Rows()
			if err != nil {
				return nil, err

			}
			defer rows.Close()

			var categoryIDs []string
			var userIDs []string
			var inquiryIDs []string
			var addressIDs []string
			var assigneeIDs []string
			var orderGroupIDs []string

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

				if alias.CategoryID != "" && !helper.StringContains(categoryIDs, alias.CategoryID) {
					categoryIDs = append(categoryIDs, alias.CategoryID)
				}

				if alias.UserID != "" && !helper.StringContains(userIDs, alias.UserID) {
					userIDs = append(userIDs, alias.UserID)
				}

				if alias.OrderGroupID != "" && !helper.StringContains(orderGroupIDs, alias.OrderGroupID) {
					orderGroupIDs = append(orderGroupIDs, alias.OrderGroupID)
				}

				if alias.ShippingAddressID != "" && !helper.StringContains(addressIDs, alias.ShippingAddressID) {
					addressIDs = append(addressIDs, alias.ShippingAddressID)
				}

				for _, v := range alias.AssigneeIDs {
					if !helper.StringContains(assigneeIDs, v) {
						assigneeIDs = append(assigneeIDs, v)
					}
				}

				alias.Inquiry.QuotedPrice = alias.Inquiry.GetQuotedPrice().ToPtr()
				records = append(records, alias.Inquiry)
			}

			var wg sync.WaitGroup

			if len(categoryIDs) > 0 {
				wg.Add(1)
				go func() {
					defer wg.Done()
					var categories []*models.Category
					db.Find(&categories, "id IN ?", categoryIDs)

					for _, cate := range categories {
						for _, record := range records {
							if record.CategoryID == cate.ID {
								record.Category = cate
							}
						}
					}
				}()

			}

			if len(userIDs) > 0 && options.IncludeUser {
				wg.Add(1)

				go func() {
					defer wg.Done()
					var items []*models.User
					db.Select("ID", "Email", "Name", "Avatar", "IsOffline", "LastOnlineAt").Find(&items, "id IN ?", userIDs)

					for _, item := range items {
						for _, record := range records {
							if record.UserID == item.ID {
								record.User = item
							}
						}
					}
				}()
			}

			if len(assigneeIDs) > 0 && options.IncludeAssignee {
				wg.Add(1)

				go func() {
					defer wg.Done()
					var items []*models.User
					db.Select("ID", "Email", "Name", "Avatar").Find(&items, "id IN ?", assigneeIDs)

					for _, item := range items {
						for _, record := range records {
							if helper.StringContains(record.AssigneeIDs, item.ID) {
								record.Assignees = append(record.Assignees, item)
							}
						}
					}
				}()
			}

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

					var purchaseOrders models.PurchaseOrders
					if err := db.Find(&purchaseOrders, "inquiry_id IN ? AND user_id IN ?", inquiryIDs, userIDs).Error; err != nil {
						return
					}

					for _, item := range purchaseOrders {
						for _, record := range records {
							if record.ID == item.InquiryID && record.UserID == item.UserID {
								record.PurchaseOrder = item
							}
						}
					}

					var orderCartItems []*models.OrderCartItem
					if err := db.Find(&orderCartItems, "purchase_order_id IN ?", purchaseOrders.IDs()).Error; err != nil {
						return
					}
					if len(orderCartItems) > 0 {
						for _, item := range orderCartItems {
							for _, record := range records {
								if record.PurchaseOrder != nil && record.PurchaseOrder.ID == item.PurchaseOrderID {
									record.PurchaseOrder.OrderCartItems = append(record.PurchaseOrder.OrderCartItems, item)
								}
							}
						}
					}

					var payments []*models.PaymentTransaction
					if err := db.Find(&payments, "count_elements(purchase_order_ids,?) >=1", pq.StringArray(purchaseOrders.IDs())).Error; err != nil {
						return
					}
					for _, payment := range payments {
						for _, record := range records {
							if len(payment.PurchaseOrderIDs) > 0 && record.PurchaseOrder != nil {
								if lo.Contains(payment.PurchaseOrderIDs, record.PurchaseOrder.ID) {
									record.PurchaseOrder.PaymentTransaction = payment
								}
							}
						}
					}
				}()

			}

			if len(addressIDs) > 0 && options.IncludeShippingAddress {
				wg.Add(1)
				go func() {
					defer wg.Done()
					var items []*models.Address

					query.New(db, NewAddressBuilder(AddressBuilderOptions{})).
						Where("a.id IN ?", addressIDs).
						FindFunc(&items)

					for _, item := range items {
						for _, record := range records {
							if record.ShippingAddressID == item.ID {
								record.ShippingAddress = item
							}
						}
					}
				}()
			}

			if len(inquiryIDs) > 0 && options.IncludeAuditLog {
				wg.Add(1)
				go func() {
					defer wg.Done()
					var items []*models.InquiryAudit
					var availableActions = []enums.AuditActionType{
						enums.AuditActionTypeInquiryCreated,
						enums.AuditActionTypeInquiryBuyerApproveQuotation,
						enums.AuditActionTypeInquirySamplePoCreated,
					}

					if options.Role.IsAdmin() {
						db.Find(&items, "inquiry_id IN ?", inquiryIDs)
					} else {
						db.Find(&items, "inquiry_id IN ? and action_type IN ?", inquiryIDs, availableActions)
					}

					for _, item := range items {
						for _, record := range records {
							if record.ID == item.InquiryID {
								record.AuditLogs = append(record.AuditLogs, item)
							}
						}
					}
				}()
			}

			wg.Wait()
			return records, nil
		})
}
