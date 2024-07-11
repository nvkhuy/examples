package queryfunc

import (
	"sync"
	"text/template"

	"github.com/engineeringinflow/inflow-backend/pkg/db"
	"github.com/engineeringinflow/inflow-backend/pkg/helper"
	"github.com/engineeringinflow/inflow-backend/pkg/models"
	"github.com/engineeringinflow/inflow-backend/pkg/repo/query"
	"github.com/lib/pq"
	"github.com/samber/lo"
	"gorm.io/plugin/dbresolver"
)

type PurchaseOrderAlias struct {
	*models.PurchaseOrder

	Inquiry *models.Inquiry `gorm:"embedded;embeddedPrefix:iq__" json:"inquiry,omitempty"`
}

type PurchaseOrderBuilderOptions struct {
	QueryBuilderOptions
	IncludeCartItems          bool
	IncludeUsers              bool
	IncludeAssignee           bool
	IncludeSampleMaker        bool
	IncludeInquirySeller      bool
	IncludeInvoice            bool
	IsConsistentRead          bool
	IncludePaymentTransaction bool
	IncludeTrackings          bool
	IncludeItems              bool
	IncludeCollection         bool
}

func NewPurchaseOrderBuilder(options PurchaseOrderBuilderOptions) *Builder {
	var rawSQL = `
	SELECT /* {{Description}} */ po.*,
	iq.id AS iq__id,
	iq.reference_id AS iq__reference_id,
	iq.title AS iq__title,
	iq.sku_note AS iq__sku_note,
	iq.user_id AS iq__user_id,
	iq.expired_date AS iq__expired_date,
	iq.delivery_date AS iq__delivery_date,
	iq.status AS iq__status,
	iq.attachments AS iq__attachments,
	iq.document AS iq__document,
	iq.design AS iq__design,
	iq.fabric_attachments AS iq__fabric_attachments,
	iq.techpack_attachments AS iq__techpack_attachments,
	iq.currency AS iq__currency,
	iq.category_id AS iq__category_id,
	iq.buyer_quotation_status AS iq__buyer_quotation_status,
	iq.new_seen_at AS iq__new_seen_at,
	iq.update_seen_at AS iq__update_seen_at,
	iq.quantity AS iq__quantity,
	iq.price_type AS iq__price_type,
	iq.quotation_at AS iq__quotation_at,
	iq.admin_quotations AS iq__admin_quotations,
	iq.approve_reject_meta AS iq__approve_reject_meta,
	iq.collection_id AS iq__collection_id,
	iq.shipping_address_id AS iq__shipping_address_id,
	iq.shipping_fee AS iq__shipping_fee,
	iq.size_list AS iq__size_list,
	iq.size_chart AS iq__size_chart,
	iq.composition AS iq__composition,
	iq.style_no AS iq__style_no,
	iq.fabric_name AS iq__fabric_name,
	iq.fabric_weight AS iq__fabric_weight,
	iq.color_list AS iq__color_list,
	iq.assignee_ids AS iq__assignee_ids,
	iq.product_weight AS iq__product_weight

	FROM purchase_orders po
	LEFT JOIN inquiries iq ON po.inquiry_id = iq.id
	`
	var countSQL = `
	SELECT /* {{Description}} */ 1

	FROM purchase_orders po
	LEFT JOIN inquiries iq ON po.inquiry_id = iq.id
	`

	var orderBy = "po.updated_at DESC"

	b := NewBuilder(rawSQL, countSQL).
		WithOrderBy(orderBy).
		WithOptions(options, template.FuncMap{
			"Description": func() string {
				return helper.JoinNonEmptyStrings(
					"-",
					GetCaller(),
					options.Role.DisplayName(),
				)
			},
		}).
		WithPaginationFunc(func(db, rawSQL *db.DB) (interface{}, error) {
			var records = make([]*models.PurchaseOrder, rawSQL.RowsAffected)

			rows, err := rawSQL.Rows()
			if err != nil {
				return nil, err

			}
			defer rows.Close()

			var purchaseOrderIDs []string
			var inquiryIDs []string
			var userIDs []string
			var assigneeIDs []string
			var cartIemIDs []string
			var sellerIDs []string
			var invoiceNumbers []int
			var checkoutSessionIDs []string
			var orderGroupIDs []string

			for rows.Next() {
				var alias PurchaseOrderAlias
				err = db.ScanRows(rows, &alias)
				if err != nil {
					db.CustomLogger.Errorf("Scan rows error", err)
					continue
				}

				if !helper.StringContains(purchaseOrderIDs, alias.PurchaseOrder.ID) {
					purchaseOrderIDs = append(purchaseOrderIDs, alias.PurchaseOrder.ID)
				}

				if !helper.StringContains(inquiryIDs, alias.PurchaseOrder.InquiryID) {
					inquiryIDs = append(inquiryIDs, alias.PurchaseOrder.InquiryID)
				}

				if !helper.StringContains(userIDs, alias.PurchaseOrder.UserID) {
					userIDs = append(userIDs, alias.PurchaseOrder.UserID)
				}

				if !helper.StringContains(userIDs, alias.CheckoutSessionID) && alias.CheckoutSessionID != "" {
					checkoutSessionIDs = append(checkoutSessionIDs, alias.CheckoutSessionID)
				}

				if !helper.StringContains(sellerIDs, alias.PurchaseOrder.SampleMakerID) && alias.PurchaseOrder.SampleMakerID != "" {
					if alias.PurchaseOrder.SampleMakerID != "inflow" {
						sellerIDs = append(sellerIDs, alias.PurchaseOrder.SampleMakerID)

					} else {
						alias.PurchaseOrder.SampleMaker = &models.User{
							Name: "Inflow Sample Room",
						}
					}
				}
				if alias.OrderGroupID != "" && !helper.StringContains(orderGroupIDs, alias.OrderGroupID) {
					orderGroupIDs = append(orderGroupIDs, alias.OrderGroupID)
				}

				for _, v := range alias.AssigneeIDs {
					if !helper.StringContains(assigneeIDs, v) {
						assigneeIDs = append(assigneeIDs, v)
					}
				}

				for _, v := range alias.PurchaseOrder.CartItemIDs {
					if !helper.StringContains(cartIemIDs, v) {
						cartIemIDs = append(cartIemIDs, v)
					}
				}

				if !lo.Contains(invoiceNumbers, alias.InvoiceNumber) && alias.InvoiceNumber > 0 {
					invoiceNumbers = append(invoiceNumbers, alias.InvoiceNumber)
				}

				if alias.Inquiry != nil && alias.Inquiry.ID != "" {
					alias.PurchaseOrder.Inquiry = alias.Inquiry
				}

				records = append(records, alias.PurchaseOrder)
			}

			var wg sync.WaitGroup

			if options.IncludeItems && len(purchaseOrderIDs) > 0 {
				wg.Add(1)
				go func() {
					defer wg.Done()

					var cartItems []*models.PurchaseOrderItem
					query.New(db, NewPurchaseOrderItemBuilder(PurchaseOrderItemBuilderOptions{
						IncludeProduct: true,
						IncludeFabric:  true,
						IncludeVariant: true,
					})).
						WhereFunc(func(builder *query.Builder) {
							builder.Where("purchase_order_id IN ?", purchaseOrderIDs)
						}).
						FindFunc(&cartItems)

					for _, cartItem := range cartItems {
						for _, record := range records {
							if cartItem.PurchaseOrderID == record.ID {
								record.Items = append(record.Items, cartItem)
							}
						}
					}

				}()

			}

			if options.IncludeCartItems {
				wg.Add(1)
				go func() {
					defer wg.Done()

					if len(cartIemIDs) > 0 {
						var items []*models.InquiryCartItem
						db.Find(&items, "id IN ?", cartIemIDs)

						for _, item := range items {
							for _, record := range records {
								if helper.StringContains(record.CartItemIDs, item.ID) {
									record.CartItems = append(record.CartItems, item)
								}
							}
						}
					}

					var orderCartItems []*models.OrderCartItem
					if err := db.Find(&orderCartItems, "purchase_order_id IN ?", purchaseOrderIDs).Error; err != nil {
						return
					}
					if len(orderCartItems) > 0 {
						for _, item := range orderCartItems {
							for _, record := range records {
								if record.ID == item.PurchaseOrderID {
									record.OrderCartItems = append(record.OrderCartItems, item)
								}
							}
						}
					}

				}()

			}

			if options.IncludePaymentTransaction && len(purchaseOrderIDs) > 0 {
				wg.Add(1)
				go func() {
					defer wg.Done()

					var items []*models.PaymentTransaction
					if err := db.Find(&items, "count_elements(purchase_order_ids,?) >=1", pq.StringArray(purchaseOrderIDs)).Error; err != nil {
						return
					}

					for _, item := range items {
						for _, record := range records {
							if len(item.PurchaseOrderIDs) > 0 && lo.Contains(item.PurchaseOrderIDs, record.ID) {
								record.PaymentTransaction = item
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

			if len(sellerIDs) > 0 && options.IncludeSampleMaker {
				wg.Add(1)

				go func() {
					defer wg.Done()
					var items []*models.User
					db.Select("ID", "Email", "Name", "CompanyName").Find(&items, "id IN ?", sellerIDs)

					for _, item := range items {
						for _, record := range records {
							if record.SampleMakerID == item.ID {
								record.SampleMaker = item
							}
						}
					}
				}()
			}

			if options.IncludeUsers && len(userIDs) > 0 {
				wg.Add(1)
				go func() {
					defer wg.Done()

					var items []*models.User
					db.Select("ID", "Name", "Email", "Avatar", "IsOffline", "LastOnlineAt").Find(&items, "id IN ?", userIDs)

					for _, item := range items {
						for _, record := range records {
							if item.ID == record.UserID {
								record.User = item
							}
						}
					}
				}()

			}

			if options.IncludeInquirySeller && len(sellerIDs) > 0 && len(inquiryIDs) > 0 {
				wg.Add(1)
				go func() {
					defer wg.Done()

					var items []*models.InquirySeller
					db.Select("*").Find(&items, "user_id IN ? AND inquiry_id IN ?", sellerIDs, inquiryIDs)

					for _, item := range items {
						for _, record := range records {
							if item.InquiryID == record.InquiryID && item.UserID == record.SampleMakerID {
								record.InquirySeller = item
							}
						}
					}
				}()

			}

			if len(invoiceNumbers) > 0 && options.IncludeInvoice {
				wg.Add(1)
				go func() {
					defer wg.Done()
					var invoices []*models.Invoice
					db.Find(&invoices, "invoice_number IN ?", invoiceNumbers)

					for _, inv := range invoices {
						for _, record := range records {
							if record.InvoiceNumber == inv.InvoiceNumber {
								record.Invoice = inv
							}
						}
					}
				}()
			}

			if options.IncludeTrackings && len(purchaseOrderIDs) > 0 {
				wg.Add(1)
				go func() {
					defer wg.Done()
					var trackings []*models.PurchaseOrderTracking
					db.Find(&trackings, "purchase_order_id IN ?", purchaseOrderIDs)

					for _, tracking := range trackings {
						for _, record := range records {
							if record.ID == tracking.PurchaseOrderID {
								record.Trackings = append(record.Trackings, tracking)
							}
						}
					}
				}()
			}
			if options.IncludeCollection && len(orderGroupIDs) > 0 {
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

			wg.Wait()

			return records, nil
		})

	if options.IsConsistentRead {
		b.WithClauses(dbresolver.Write)
	}
	return b
}
