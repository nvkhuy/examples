package queryfunc

import (
	"sync"
	"text/template"

	"github.com/engineeringinflow/inflow-backend/pkg/db"
	"github.com/engineeringinflow/inflow-backend/pkg/helper"
	"github.com/engineeringinflow/inflow-backend/pkg/models"
	"github.com/engineeringinflow/inflow-backend/pkg/repo/query"
)

type OrderCartPurchaseOrderBuilderOptions struct {
	QueryBuilderOptions
	IncludeInquiry    bool
	IncludeAddress    bool
	IncludeCollection bool
}

func NewOrderCartPurchaseOrderBuilder(options OrderCartPurchaseOrderBuilderOptions) *Builder {
	var rawSQL = `
	SELECT /* {{Description}} */ 
	po.id,
	po.reference_id,
	po.status,
	po.product_name,
	po.attachments,
	po.inquiry_id,
	po.currency,
	po.sub_total,
	po.shipping_address_id,
	po.order_group_id

	FROM purchase_orders po
	`
	return NewBuilder(rawSQL).
		WithOptions(options, template.FuncMap{
			"Description": func() string {
				return helper.JoinNonEmptyStrings(
					"-",
					GetCaller(),
					options.Role.DisplayName(),
				)
			},
		}).
		WithOrderBy("po.updated_at DESC").
		WithPaginationFunc(func(db, rawSQL *db.DB) (interface{}, error) {
			var records []*models.PurchaseOrder

			rows, err := rawSQL.Rows()
			if err != nil {
				return nil, err

			}
			defer rows.Close()

			var purchaseOrderIDs []string
			var inquiryIDs []string
			var addressIDs []string
			var orderGroupIDs []string

			for rows.Next() {
				var alias models.PurchaseOrder
				err = db.ScanRows(rows, &alias)
				if err != nil {
					db.CustomLogger.Errorf("Scan rows error", err)
					continue
				}
				purchaseOrderIDs = append(purchaseOrderIDs, alias.ID)

				if alias.InquiryID != "" {
					inquiryIDs = append(inquiryIDs, alias.InquiryID)
				}
				if alias.ShippingAddressID != "" {
					addressIDs = append(addressIDs, alias.ShippingAddressID)
				}
				if alias.OrderGroupID != "" {
					orderGroupIDs = append(orderGroupIDs, alias.OrderGroupID)
				}

				records = append(records, &alias)
			}

			var wg sync.WaitGroup
			if len(purchaseOrderIDs) > 0 {
				wg.Add(1)
				go func() {
					defer wg.Done()
					var items []*models.OrderCartItem
					if err := db.Find(&items, "purchase_order_id IN ?", purchaseOrderIDs).Error; err != nil {
						return
					}
					for _, item := range items {
						for _, record := range records {
							if record.ID == item.PurchaseOrderID {
								record.OrderCartItems = append(record.OrderCartItems, item)
							}
						}
					}
				}()
			}

			if options.IncludeInquiry && len(inquiryIDs) > 0 {
				wg.Add(1)
				go func() {
					defer wg.Done()
					var inquiries = make([]*models.Inquiry, 0, len(inquiryIDs))
					if err := db.Select("ID", "ReferenceID", "Title", "Attachments", "Currency").Find(&inquiries, "id IN ?", inquiryIDs).Error; err != nil {
						return
					}

					for _, iq := range inquiries {
						for _, record := range records {
							if record.InquiryID == iq.ID {
								record.Inquiry = iq
							}
						}
					}
				}()
			}
			if options.IncludeAddress && len(addressIDs) > 0 {
				wg.Add(1)
				go func() {
					defer wg.Done()

					var addresses = make([]*models.Address, 0, len(addressIDs))

					query.New(db, NewAddressBuilder(AddressBuilderOptions{})).
						Where("a.id IN ?", addressIDs).
						FindFunc(&addresses)

					for _, adr := range addresses {
						for _, record := range records {
							if record.ShippingAddressID == adr.ID {
								record.ShippingAddress = adr
							}
						}
					}
				}()
			}
			if options.IncludeCollection && len(orderGroupIDs) > 0 {
				wg.Add(1)
				go func() {
					defer wg.Done()
					var orderGroups = make([]*models.OrderGroup, 0, len(orderGroupIDs))
					if err := db.Select("ID", "Name").Find(&orderGroups, "id IN ?", orderGroupIDs).Error; err != nil {
						return
					}

					for _, og := range orderGroups {
						for _, record := range records {
							if record.OrderGroupID == og.ID {
								record.OrderGroup = og
							}
						}
					}
				}()
			}

			wg.Wait()
			return &records, nil
		})
}

type OrderCartBulkBuilderOptions struct {
	QueryBuilderOptions
	IncludeInquiry    bool
	IncludeAddress    bool
	IncludeCollection bool
}

func NewOrderCartBulkBuilder(options OrderCartBulkBuilderOptions) *Builder {
	var rawSQL = `
	SELECT /* {{Description}} */ 
	bpo.id,
	bpo.reference_id,
	bpo.inquiry_id,
	bpo.purchase_order_id,
	bpo.product_name,
	bpo.attachments,
	bpo.status,
	bpo.tracking_status,
	bpo.currency,
	bpo.tax,
	bpo.total_price,
	bpo.shipping_fee,
	bpo.sub_total,
	bpo.total_price,
	bpo.transaction_fee,
	bpo.first_payment_total,
	bpo.first_payment_sub_total,
	bpo.first_payment_percentage,
	bpo.final_payment_total,
	bpo.final_payment_sub_total,
	bpo.final_payment_transaction_fee,
	bpo.final_payment_tax,
	bpo.shipping_address_id,
	bpo.order_group_id

	FROM bulk_purchase_orders bpo
	`
	return NewBuilder(rawSQL).
		WithOptions(options, template.FuncMap{
			"Description": func() string {
				return helper.JoinNonEmptyStrings(
					"-",
					GetCaller(),
					options.Role.DisplayName(),
				)
			},
		}).
		WithOrderBy("bpo.updated_at DESC").
		WithPaginationFunc(func(db, rawSQL *db.DB) (interface{}, error) {
			var records []*models.BulkPurchaseOrder

			rows, err := rawSQL.Rows()
			if err != nil {
				return nil, err

			}
			defer rows.Close()

			var bulkPurchaseOrderIDs []string
			var inquiryIDs []string
			var addressIDs []string
			var orderGroupIDs []string

			for rows.Next() {
				var alias models.BulkPurchaseOrder
				err = db.ScanRows(rows, &alias)
				if err != nil {
					db.CustomLogger.Errorf("Scan rows error", err)
					continue
				}
				bulkPurchaseOrderIDs = append(bulkPurchaseOrderIDs, alias.ID)

				if alias.InquiryID != "" {
					inquiryIDs = append(inquiryIDs, alias.InquiryID)
				}
				if alias.ShippingAddressID != "" {
					addressIDs = append(addressIDs, alias.ShippingAddressID)
				}
				if alias.OrderGroupID != "" {
					orderGroupIDs = append(orderGroupIDs, alias.OrderGroupID)
				}

				records = append(records, &alias)
			}

			var wg sync.WaitGroup
			if len(bulkPurchaseOrderIDs) > 0 {
				wg.Add(1)
				go func() {
					defer wg.Done()
					var items []*models.OrderCartItem
					if err := db.Find(&items, "bulk_purchase_order_id IN ?", bulkPurchaseOrderIDs).Error; err != nil {
						return
					}
					for _, item := range items {
						for _, record := range records {
							if record.ID == item.BulkPurchaseOrderID {
								record.OrderCartItems = append(record.OrderCartItems, item)
							}
						}
					}
				}()
			}

			if options.IncludeInquiry && len(inquiryIDs) > 0 {
				wg.Add(1)
				go func() {
					defer wg.Done()
					var inquiries []*models.Inquiry
					if err := db.Select("ID", "ReferenceID", "Title", "Attachments", "Currency").Find(&inquiries, "id IN ?", inquiryIDs).Error; err != nil {
						return
					}

					for _, iq := range inquiries {
						for _, record := range records {
							if record.InquiryID == iq.ID {
								record.Inquiry = iq
							}
						}
					}
				}()
			}
			if options.IncludeAddress && len(addressIDs) > 0 {
				wg.Add(1)
				go func() {
					defer wg.Done()

					var addresses = make([]*models.Address, 0, len(addressIDs))

					query.New(db, NewAddressBuilder(AddressBuilderOptions{})).
						Where("a.id IN ?", addressIDs).
						FindFunc(&addresses)

					for _, adr := range addresses {
						for _, record := range records {
							if record.ShippingAddressID == adr.ID {
								record.ShippingAddress = adr
							}
						}
					}
				}()
			}
			if options.IncludeCollection && len(orderGroupIDs) > 0 {
				wg.Add(1)
				go func() {
					defer wg.Done()
					var orderGroups = make([]*models.OrderGroup, 0, len(orderGroupIDs))
					if err := db.Select("ID", "Name").Find(&orderGroups, "id IN ?", orderGroupIDs).Error; err != nil {
						return
					}

					for _, og := range orderGroups {
						for _, record := range records {
							if record.OrderGroupID == og.ID {
								record.OrderGroup = og
							}
						}
					}
				}()
			}

			wg.Wait()
			return &records, nil
		})
}
