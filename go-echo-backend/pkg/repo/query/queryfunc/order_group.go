package queryfunc

import (
	"sync"
	"text/template"

	"github.com/engineeringinflow/inflow-backend/pkg/db"
	"github.com/engineeringinflow/inflow-backend/pkg/helper"
	"github.com/engineeringinflow/inflow-backend/pkg/models"
	"github.com/engineeringinflow/inflow-backend/pkg/models/enums"
	"github.com/engineeringinflow/inflow-backend/pkg/repo/query"
)

type OrderGroupBuilderOptions struct {
	QueryBuilderOptions
	WithOrderDetail bool
}

func NewOrderGroupBuilder(options OrderGroupBuilderOptions) *Builder {
	var rawSQL = `
		SELECT /* {{Description}} */ o.*
	
		FROM order_groups o
	`

	var countSQL = `
	SELECT /* {{Description}} */ 1

	FROM order_groups o
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
		WithOrderBy("o.updated_at DESC").
		WithPaginationFunc(func(db, rawSQL *db.DB) (interface{}, error) {
			var records = make([]*models.OrderGroup, rawSQL.RowsAffected)

			rows, err := rawSQL.Rows()
			if err != nil {
				db.CustomLogger.Errorf("Scan rows error", err)
				return nil, err
			}
			defer rows.Close()

			var orderGroupIDs []string
			var userIDs []string

			for rows.Next() {
				var alias models.OrderGroup
				err = db.ScanRows(rows, &alias)
				if err != nil {
					db.CustomLogger.Errorf("Scan rows error", err)
					continue
				}
				orderGroupIDs = append(orderGroupIDs, alias.ID)
				userIDs = append(userIDs, alias.UserID)
				records = append(records, &alias)
			}
			var wg sync.WaitGroup
			wg.Add(4)

			// get user belong to order group
			go func() {
				defer wg.Done()
				var users []*models.User
				if err := db.Select("ID", "Name", "Avatar").Find(&users, "id IN ?", userIDs).Error; err != nil {
					return
				}
				for _, user := range users {
					for _, record := range records {
						if user.ID == record.UserID {
							record.User = user
						}
					}
				}
			}()

			// get inquiries belong to order group
			go func() {
				defer wg.Done()
				var inquiries []*models.Inquiry
				query.New(db, NewOrderGroupInquiriesBuilder(OrderGroupInquiriesBuilderOptions{
					QueryBuilderOptions:     options.QueryBuilderOptions,
					WithAudit:               options.WithOrderDetail,
					WithPurchaseOrderStatus: options.WithOrderDetail,
					WithShippingAddress:     options.WithOrderDetail,
				})).
					WhereFunc(func(builder *query.Builder) {
						builder.Where("order_group_id IN ?", orderGroupIDs)
					}).
					FindFunc(&inquiries)

				for _, iq := range inquiries {
					for _, record := range records {
						if iq.OrderGroupID == record.ID {
							record.Inquiries = append(record.Inquiries, iq)
						}
					}
				}

			}()
			// get purchase orders belong to order group
			go func() {
				defer wg.Done()

				var purchaseOrders []*models.PurchaseOrder
				query.New(db, NewOrderGroupSamplesBuilder(OrderGroupSamplesBuilderOptions{
					QueryBuilderOptions:        options.QueryBuilderOptions,
					WithOrderCartItems:         options.WithOrderDetail,
					WithShippingAndUserAddress: options.WithOrderDetail,
					WithTrackingLogs:           options.WithOrderDetail,
				})).
					WhereFunc(func(builder *query.Builder) {
						builder.Where("order_group_id IN ?", orderGroupIDs)
						builder.Where("COALESCE(NULLIF(po.from_catalog,false),po.status = ?) = true", "paid")
					}).
					FindFunc(&purchaseOrders)

				for _, record := range records {
					for _, po := range purchaseOrders {
						if record.ID == po.OrderGroupID {
							record.Samples = append(record.Samples, po)
						}
					}
				}
			}()
			// get bulk purchase order belong to order group
			go func() {
				defer wg.Done()

				var bulkPurchaseOrders []*models.BulkPurchaseOrder
				query.New(db, NewOrderGroupBulksBuilder(OrderGroupBulksBuilderOptions{
					QueryBuilderOptions: options.QueryBuilderOptions,
					WithTrackingLogs:    options.WithOrderDetail,
					WithOrderCartItems:  options.WithOrderDetail,
					WithShippingAddress: options.WithOrderDetail,
					WithPurchaseOrder:   options.WithOrderDetail,
				})).
					WhereFunc(func(builder *query.Builder) {
						builder.Where("order_group_id IN ?", orderGroupIDs)
					}).
					FindFunc(&bulkPurchaseOrders)

				for _, record := range records {
					for _, bpo := range bulkPurchaseOrders {
						if record.ID == bpo.OrderGroupID {
							record.Bulks = append(record.Bulks, bpo)
						}
					}
				}
			}()

			wg.Wait()
			return &records, nil
		})
}

type OrderGroupInquiriesBuilderOptions struct {
	QueryBuilderOptions
	WithAudit               bool
	WithPurchaseOrderStatus bool
	WithShippingAddress     bool
}

func NewOrderGroupInquiriesBuilder(options OrderGroupInquiriesBuilderOptions) *Builder {
	var rawSQL = `
		SELECT /* {{Description}} */
		iq.id,
		iq.reference_id,
		iq.title,
		iq.sku_note,
		iq.attachments, 
		iq.status,
		iq.currency,
		iq.buyer_quotation_status,
		iq.techpack_attachments,
		iq.fabric_attachments,
		iq.quantity,      
		iq.expected_price,
		iq.size_list,
		iq.size_chart,
		iq.style_no,
		iq.fabric_name,
		iq.fabric_weight,
		iq.composition,
		iq.color_list,
		iq.created_at,
		iq.order_group_id,
		iq.shipping_address_id,
		iq.admin_quotations,
		iq.shipping_fee,
		iq.tax_percentage
		
		FROM inquiries iq
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
		WithOrderBy("iq.created_at DESC").
		WithPaginationFunc(func(db, rawSQL *db.DB) (interface{}, error) {
			var records = make([]*models.Inquiry, rawSQL.RowsAffected)

			rows, err := rawSQL.Rows()
			if err != nil {
				db.CustomLogger.Errorf("Scan rows error", err)
				return nil, err
			}
			defer rows.Close()

			var inquiryIDs []string
			var shippingAddressIDs []string

			for rows.Next() {
				var alias models.Inquiry
				err = db.ScanRows(rows, &alias)
				if err != nil {
					db.CustomLogger.Errorf("Scan rows error", err)
					continue
				}
				inquiryIDs = append(inquiryIDs, alias.ID)
				if alias.ShippingAddressID != "" {
					shippingAddressIDs = append(shippingAddressIDs, alias.ShippingAddressID)
				}

				records = append(records, &alias)
			}
			var wg sync.WaitGroup

			if options.WithAudit && len(inquiryIDs) > 0 {
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
			if options.WithPurchaseOrderStatus && len(inquiryIDs) > 0 {
				wg.Add(1)
				go func() {
					defer wg.Done()
					var purchaseOrders models.PurchaseOrders
					if err := db.Select("ID", "Status", "MarkAsPaidAt", "TransferedAt", "InquiryID").Find(&purchaseOrders, "inquiry_id IN ?", inquiryIDs).Error; err != nil {
						return
					}
					var orderCartItems models.OrderCartItems
					if err := db.Find(&orderCartItems, "purchase_order_id IN ?", purchaseOrders.IDs()).Error; err != nil {
						return
					}
					for _, item := range orderCartItems {
						for _, po := range purchaseOrders {
							if po.ID == item.PurchaseOrderID {
								po.OrderCartItems = append(po.OrderCartItems, item)
							}
						}
					}

					for _, po := range purchaseOrders {
						for _, record := range records {
							if record.ID == po.InquiryID {
								record.PurchaseOrder = po
							}
						}
					}
				}()
			}
			if options.WithShippingAddress && len(shippingAddressIDs) > 0 {
				wg.Add(1)
				go func() {
					defer wg.Done()
					var addresses []*models.Address
					query.New(db, NewAddressBuilder(AddressBuilderOptions{})).
						Where("a.id IN ?", shippingAddressIDs).
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

			wg.Wait()

			return &records, nil
		})
}

type OrderGroupSamplesBuilderOptions struct {
	QueryBuilderOptions
	WithOrderCartItems         bool
	WithShippingAndUserAddress bool
	WithTrackingLogs           bool
}

func NewOrderGroupSamplesBuilder(options OrderGroupSamplesBuilderOptions) *Builder {
	var rawSQL = `
		SELECT /* {{Description}} */
		po.id,
		po.product_name,
		po.reference_id,
		po.user_id,
		po.inquiry_id,
		po.cart_item_ids,
		po.status,
		po.currency,
		po.tracking_status,
		po.attachments,
		po.techpack_attachments,
		po.fabric_attachments,
		po.shipping_address_id,
		po.created_at,
		po.order_group_id,
		po.from_catalog
		
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
			var records = make([]*models.PurchaseOrder, rawSQL.RowsAffected)

			rows, err := rawSQL.Rows()
			if err != nil {
				db.CustomLogger.Errorf("Scan rows error", err)
				return nil, err
			}
			defer rows.Close()

			var purchaseOrderIDs []string
			var inquiryIDs []string
			var purchaseOrderIDsToCatalogCartItems []string
			var shippingAddressIDs []string
			var userIDsToAddress []string

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
				if alias.FromCatalog {
					purchaseOrderIDsToCatalogCartItems = append(purchaseOrderIDsToCatalogCartItems, alias.ID)
				}

				if alias.ShippingAddressID != "" {
					shippingAddressIDs = append(shippingAddressIDs, alias.ShippingAddressID)
				} else {
					userIDsToAddress = append(userIDsToAddress, alias.UserID)
				}

				records = append(records, &alias)
			}
			var wg sync.WaitGroup

			if len(inquiryIDs) > 0 {
				wg.Add(1)
				go func() {
					defer wg.Done()
					var inquiries []*models.Inquiry
					if err := db.Select("ID", "ReferenceID", "Title", "SkuNote").Find(&inquiries, "id IN ?", inquiryIDs).Error; err != nil {
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

			if options.WithOrderCartItems {
				wg.Add(1)
				go func() {
					defer wg.Done()
					var orderCartItems models.OrderCartItems
					if err := db.Find(&orderCartItems, "purchase_order_id IN ?", purchaseOrderIDs).Error; err != nil {
						return
					}
					for _, item := range orderCartItems {
						for _, record := range records {
							if record.ID == item.PurchaseOrderID {
								record.OrderCartItems = append(record.OrderCartItems, item)
							}
						}
					}
				}()
			}

			if len(purchaseOrderIDsToCatalogCartItems) > 0 {
				wg.Add(1)
				go func() {
					defer wg.Done()

					var cartItems []*models.PurchaseOrderItem
					query.New(db, NewPurchaseOrderItemBuilder(PurchaseOrderItemBuilderOptions{
						IncludeProduct: true,
					})).
						WhereFunc(func(builder *query.Builder) {
							builder.Where("purchase_order_id IN ?", purchaseOrderIDsToCatalogCartItems)
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
			if options.WithShippingAndUserAddress && len(shippingAddressIDs) > 0 {
				wg.Add(1)
				go func() {
					defer wg.Done()

					var addresses []*models.Address

					query.New(db, NewAddressBuilder(AddressBuilderOptions{})).
						Where("a.id IN ?", shippingAddressIDs).
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

			if options.WithShippingAndUserAddress && len(userIDsToAddress) > 0 {
				wg.Add(1)
				go func() {
					defer wg.Done()

					var users []*models.User
					if err := db.Select("ID", "Name", "PhoneNumber", "Email", "CoordinateID").Find(&users, "id IN ?", userIDsToAddress).Error; err != nil {
						return
					}
					var coordinateIDs []string
					for _, user := range users {
						coordinateIDs = append(coordinateIDs, user.CoordinateID)
					}
					var coordinates []*models.Coordinate
					if err := db.Find(&coordinates, "id IN ?", coordinateIDs).Error; err != nil {
						return
					}
					for _, user := range users {
						for _, coord := range coordinates {
							if coord.ID == user.CoordinateID {
								user.Coordinate = coord
							}
						}
					}
					for _, user := range users {
						for _, record := range records {
							if record.UserID == user.ID {
								record.User = user
							}
						}
					}
				}()
			}
			if options.WithTrackingLogs && len(purchaseOrderIDs) > 0 {
				wg.Add(1)
				go func() {
					defer wg.Done()

					var trackingLogs []*models.PurchaseOrderTracking
					if err := db.Select("ID", "PurchaseOrderID", "ActionType", "CreatedAt").Find(&trackingLogs, "purchase_order_id IN ?", purchaseOrderIDs).Error; err != nil {
						return
					}
					for _, log := range trackingLogs {
						for _, record := range records {
							if record.ID == log.PurchaseOrderID {
								record.Trackings = append(record.Trackings, log)
							}
						}
					}
				}()
			}

			wg.Wait()

			return &records, nil
		})
}

type OrderGroupBulksBuilderOptions struct {
	QueryBuilderOptions
	WithTrackingLogs    bool
	WithShippingAddress bool
	WithOrderCartItems  bool
	WithPurchaseOrder   bool
}

func NewOrderGroupBulksBuilder(options OrderGroupBulksBuilderOptions) *Builder {
	var rawSQL = `
		SELECT /* {{Description}} */
		bpo.id,
		bpo.product_name,
		bpo.note,
		bpo.attachments,
		bpo.reference_id,
		bpo.inquiry_id,
		bpo.purchase_order_id,
		bpo.status,
		bpo.currency,
		bpo.tracking_status,
		bpo.techpack_attachments,
		bpo.first_payment_mark_as_paid_at,
		bpo.first_payment_transfered_at,
		bpo.first_payment_received_at,
		bpo.first_payment_percentage,
		bpo.final_payment_mark_as_paid_at,
		bpo.final_payment_transfered_at,
		bpo.final_payment_received_at,
		bpo.size_attachments,
		bpo.quotation_note,
		bpo.quotation_note_attachments,
		bpo.admin_quotations,
		bpo.created_at,
		bpo.order_group_id,
		bpo.shipping_address_id
		
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
		WithOrderBy("bpo.created_at DESC").
		WithPaginationFunc(func(db, rawSQL *db.DB) (interface{}, error) {
			var records = make([]*models.BulkPurchaseOrder, rawSQL.RowsAffected)

			rows, err := rawSQL.Rows()
			if err != nil {
				db.CustomLogger.Errorf("Scan rows error", err)
				return nil, err
			}
			defer rows.Close()

			var bulkPurchaseOrderIDs []string
			var inquiryIDs []string
			var purchaseOrderIDs []string
			var purchaseOrderIDsToCatalogCartItems []string
			var shippingAddressIDs []string

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
				} else if alias.InquiryID == "" && alias.PurchaseOrderID != "" {
					purchaseOrderIDs = append(purchaseOrderIDs, alias.PurchaseOrderID)
					// purchaseOrderIDsToCatalogCartItems = append(purchaseOrderIDsToCatalogCartItems, alias.PurchaseOrderID)
				}
				if alias.ShippingAddressID != "" {
					shippingAddressIDs = append(shippingAddressIDs, alias.ShippingAddressID)
				}

				records = append(records, &alias)
			}
			var wg sync.WaitGroup

			if len(inquiryIDs) > 0 {
				wg.Add(1)
				go func() {
					defer wg.Done()
					var inquiries []*models.Inquiry
					if err := db.Select("ID", "Title", "ReferenceID", "SkuNote", "Attachments", "AdminQuotations").Find(&inquiries).Error; err != nil {
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
			if options.WithOrderCartItems {
				wg.Add(1)
				go func() {
					defer wg.Done()
					var orderCartItems models.OrderCartItems
					if err := db.Find(&orderCartItems, "bulk_purchase_order_id IN ?", bulkPurchaseOrderIDs).Error; err != nil {
						return
					}
					for _, item := range orderCartItems {
						for _, record := range records {
							if record.ID == item.BulkPurchaseOrderID {
								record.OrderCartItems = append(record.OrderCartItems, item)
							}
						}
					}
				}()
			}
			if options.WithShippingAddress && len(shippingAddressIDs) > 0 {
				wg.Add(1)
				go func() {
					defer wg.Done()
					var addresses []*models.Address
					query.New(db, NewAddressBuilder(AddressBuilderOptions{})).
						Where("a.id IN ?", shippingAddressIDs).
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
			if options.WithPurchaseOrder && len(purchaseOrderIDs) > 0 {
				wg.Add(1)
				go func() {
					defer wg.Done()
					var purchaseOrders = make(models.PurchaseOrders, 0, len(purchaseOrderIDs))
					if err := db.Select("ID", "Status").Find(&purchaseOrders, "id IN ?", purchaseOrderIDs).Error; err != nil {
						return
					}

					var orderCartItems models.OrderCartItems
					if err := db.Find(&orderCartItems, "purchase_order_id IN ?", purchaseOrders.IDs()).Error; err != nil {
						return
					}
					for _, item := range orderCartItems {
						for _, po := range purchaseOrders {
							if po.ID == item.PurchaseOrderID {
								po.OrderCartItems = append(po.OrderCartItems, item)
							}
						}
					}

					for _, po := range purchaseOrders {
						for _, record := range records {
							if record.PurchaseOrderID == po.ID {
								record.PurchaseOrder = po
							}
						}
					}

				}()
			}

			if len(purchaseOrderIDsToCatalogCartItems) > 0 {
				wg.Add(1)
				go func() {
					defer wg.Done()

					var purchaseOrders []*models.PurchaseOrder
					if err := db.Select("ID").Find(&purchaseOrders).Error; err != nil {
						return
					}

					var cartItems []*models.PurchaseOrderItem
					query.New(db, NewPurchaseOrderItemBuilder(PurchaseOrderItemBuilderOptions{
						IncludeProduct: true,
					})).
						WhereFunc(func(builder *query.Builder) {
							builder.Where("purchase_order_id IN ?", purchaseOrderIDsToCatalogCartItems)
						}).
						FindFunc(&cartItems)

					for _, po := range purchaseOrders {
						for _, cartItem := range cartItems {
							if cartItem.PurchaseOrderID == po.ID {
								po.Items = append(po.Items, cartItem)
							}
						}
					}

					for _, po := range purchaseOrders {
						for _, record := range records {
							if po.ID == record.PurchaseOrderID {
								record.PurchaseOrder = po
							}
						}
					}
				}()
			}
			if options.WithTrackingLogs && len(bulkPurchaseOrderIDs) > 0 {
				wg.Add(1)
				go func() {
					defer wg.Done()

					var trackingLogs []*models.BulkPurchaseOrderTracking
					if err := db.Select("ID", "PurchaseOrderID", "ActionType", "CreatedAt").Find(&trackingLogs, "purchase_order_id IN ?", bulkPurchaseOrderIDs).Error; err != nil {
						return
					}
					for _, log := range trackingLogs {
						for _, record := range records {
							if record.ID == log.PurchaseOrderID {
								record.Trackings = append(record.Trackings, log)
							}
						}
					}
				}()
			}

			wg.Wait()

			return &records, nil
		})
}
