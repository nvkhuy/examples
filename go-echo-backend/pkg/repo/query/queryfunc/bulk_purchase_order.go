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
)

type BulkPurchaseOrderAlias struct {
	*models.BulkPurchaseOrder

	Seller        *models.User                             `gorm:"embedded;embeddedPrefix:s__" json:"seller,omitempty"`
	BulkQuotation *models.BulkPurchaseOrderSellerQuotation `gorm:"embedded;embeddedPrefix:bposq__" json:"bulk_quotation,omitempty"`
	Inquiry       *models.Inquiry                          `gorm:"embedded;embeddedPrefix:iq__" json:"inquiry,omitempty"`
	PurchaseOrder *models.PurchaseOrder                    `gorm:"embedded;embeddedPrefix:po__" json:"purchase_order,omitempty"`
}

type BulkPurchaseOrderBuilderOptions struct {
	QueryBuilderOptions

	IncludeShippingAddress     bool
	IncludePaymentTransactions bool
	IncludeItems               bool
	IncludeUser                bool
	IncludeSellerQuotation     bool
	IncludeAssignee            bool
	IncludeInvoice             bool
	IncludeTrackings           bool
	IncludeCollection          bool
}

func NewBulkPurchaseOrderBuilder(options BulkPurchaseOrderBuilderOptions) *Builder {
	var rawSQL = `
	SELECT /* {{Description}} */ bpo.*,
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
	iq.tax_percentage AS iq__tax_percentage,
	iq.product_weight AS iq__product_weight,
	iq.expected_price AS iq__expected_price,

	po.id AS po__id,
	po.created_at AS po__created_at,
	po.updated_at AS po__updated_at,
	po.deleted_at AS po__deleted_at,
	po.reference_id AS po__reference_id,
	po.client_reference_id AS po__client_reference_id,
	po.checkout_session_id AS po__checkout_session_id,
	po.user_id AS po__user_id,
	po.inquiry_id AS po__inquiry_id,
	po.status AS po__status,
	po.payment_intent_id AS po__payment_intent_id,
	po.charge_id AS po__charge_id,
	po.receipt_url AS po__receipt_url,
	po.payment_type AS po__payment_type,
	po.transfered_at AS po__transfered_at,
	po.mark_as_paid_at AS po__mark_as_paid_at,
	po.mark_as_unpaid_at AS po__mark_as_unpaid_at,
	po.transaction_ref_id AS po__transaction_ref_id,
	po.transaction_attachment AS po__transaction_attachment,
	po.payout_transaction_ref_id AS po__payout_transaction_ref_id,
	po.payout_transaction_attachment AS po__payout_transaction_attachment,
	po.payout_transfered_at AS po__payout_transfered_at,
	po.payout_mark_as_paid_at AS po__payout_mark_as_paid_at,
	po.attachments AS po__attachments,
	po.document AS po__document,
	po.design AS po__design,
	po.fabric_attachments AS po__fabric_attachments,
	po.techpack_attachments AS po__techpack_attachments,
	po.sample_attachments AS po__sample_attachments,
	po.approve_reject_meta AS po__approve_reject_meta,
	po.tracking_status AS po__tracking_status,
	po.po_raw_materials AS po__po_raw_materials,
	po.logistic_info AS po__logistic_info,
	po.making_info AS po__making_info,
	po.submit_info AS po__submit_info,
	po.receiver_confirmed_at AS po__receiver_confirmed_at,
	po.delivery_started_at AS po__delivery_started_at,
	po.assignee_ids AS po__assignee_ids,
	po.currency AS po__currency,
	po.feedback AS po__feedback,
	po.sample_maker_id AS po__sample_maker_id,
	po.seller_po_attachment AS po__seller_po_attachment,
	po.seller_design AS po__seller_design,
	po.seller_techpack_attachments AS po__seller_techpack_attachments,
	po.seller_est_making_at AS po__seller_est_making_at,
	po.seller_est_delivery_at AS po__seller_est_delivery_at,
	po.seller_submit_info AS po__seller_submit_info,
	po.seller_logistic_info AS po__seller_logistic_info,
	po.seller_delivery_started_at AS po__seller_delivery_started_at,
	po.seller_delivery_confirmed_at AS po__seller_delivery_confirmed_at,
	po.seller_delivery_feedback AS po__seller_delivery_feedback,
	po.invoice_number AS po__invoice_number,
	po.payment_transaction_reference_id AS po__payment_transaction_reference_id,
	po.refund_reason AS po__refund_reason,
	po.sub_total AS po__sub_total,
	po.sub_total_after_deduction AS po__sub_total_after_deduction,
	po.shipping_fee AS po__shipping_fee,
	po.transaction_fee AS po__transaction_fee,
	po.tax AS po__tax,
	po.total_price AS po__total_price,
	po.tax_percentage AS po__tax_percentage,
	po.seller_total_price AS po__seller_total_price

	FROM bulk_purchase_orders bpo
	LEFT JOIN inquiries iq ON bpo.inquiry_id = iq.id
	LEFT JOIN purchase_orders po ON bpo.purchase_order_id = po.id
	`
	var countSQL = `
	SELECT /* {{Description}} */ 1

	FROM bulk_purchase_orders bpo
	LEFT JOIN inquiries iq ON bpo.inquiry_id = iq.id
	LEFT JOIN purchase_orders po ON bpo.purchase_order_id = po.id
	`

	var orderBy = "bpo.updated_at DESC"

	if options.Role.IsSeller() {
		rawSQL = `
		SELECT /* {{Description}} */ bpo.*,
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
		iq.tax_percentage AS iq__tax_percentage,
		iq.product_weight AS iq__product_weight,
		iq.expected_price AS iq__expected_price,

		po.id AS po__id,
		po.created_at AS po__created_at,
		po.updated_at AS po__updated_at,
		po.deleted_at AS po__deleted_at,
		po.reference_id AS po__reference_id,
		po.client_reference_id AS po__client_reference_id,
		po.checkout_session_id AS po__checkout_session_id,
		po.user_id AS po__user_id,
		po.inquiry_id AS po__inquiry_id,
		po.status AS po__status,
		po.payment_intent_id AS po__payment_intent_id,
		po.charge_id AS po__charge_id,
		po.receipt_url AS po__receipt_url,
		po.payment_type AS po__payment_type,
		po.transfered_at AS po__transfered_at,
		po.mark_as_paid_at AS po__mark_as_paid_at,
		po.mark_as_unpaid_at AS po__mark_as_unpaid_at,
		po.transaction_ref_id AS po__transaction_ref_id,
		po.transaction_attachment AS po__transaction_attachment,
		po.payout_transaction_ref_id AS po__payout_transaction_ref_id,
		po.payout_transaction_attachment AS po__payout_transaction_attachment,
		po.payout_transfered_at AS po__payout_transfered_at,
		po.payout_mark_as_paid_at AS po__payout_mark_as_paid_at,
		po.attachments AS po__attachments,
		po.document AS po__document,
		po.design AS po__design,
		po.fabric_attachments AS po__fabric_attachments,
		po.techpack_attachments AS po__techpack_attachments,
		po.sample_attachments AS po__sample_attachments,
		po.approve_reject_meta AS po__approve_reject_meta,
		po.tracking_status AS po__tracking_status,
		po.po_raw_materials AS po__po_raw_materials,
		po.logistic_info AS po__logistic_info,
		po.making_info AS po__making_info,
		po.submit_info AS po__submit_info,
		po.receiver_confirmed_at AS po__receiver_confirmed_at,
		po.delivery_started_at AS po__delivery_started_at,
		po.assignee_ids AS po__assignee_ids,
		po.currency AS po__currency,
		po.feedback AS po__feedback,
		po.sample_maker_id AS po__sample_maker_id,
		po.seller_po_attachment AS po__seller_po_attachment,
		po.seller_design AS po__seller_design,
		po.seller_techpack_attachments AS po__seller_techpack_attachments,
		po.seller_est_making_at AS po__seller_est_making_at,
		po.seller_est_delivery_at AS po__seller_est_delivery_at,
		po.seller_submit_info AS po__seller_submit_info,
		po.seller_logistic_info AS po__seller_logistic_info,
		po.seller_delivery_started_at AS po__seller_delivery_started_at,
		po.seller_delivery_confirmed_at AS po__seller_delivery_confirmed_at,
		po.seller_delivery_feedback AS po__seller_delivery_feedback,
		po.invoice_number AS po__invoice_number,
		po.payment_transaction_reference_id AS po__payment_transaction_reference_id,
		po.refund_reason AS po__refund_reason,
		po.sub_total AS po__sub_total,
		po.sub_total_after_deduction AS po__sub_total_after_deduction,
		po.shipping_fee AS po__shipping_fee,
		po.transaction_fee AS po__transaction_fee,
		po.tax AS po__tax,
		po.total_price AS po__total_price,
		po.tax_percentage AS po__tax_percentage,
		po.seller_total_price AS po__seller_total_price,

		bposq.id AS bposq__id,
		bposq.created_at AS bposq__created_at,
		bposq.updated_at AS bposq__updated_at,
		bposq.due_day AS bposq__due_day,
		bposq.status AS bposq__status,
		bposq.delivery_date AS bposq__delivery_date,
		bposq.currency AS bposq__currency,
		bposq.order_type AS bposq__order_type,
		bposq.offer_price AS bposq__offer_price,
		bposq.offer_remark AS bposq__offer_remark,
		bposq.variance_amount AS bposq__variance_amount,
		bposq.variance_percentage AS bposq__variance_percentage,
		bposq.fabric_cost AS bposq__fabric_cost,
		bposq.decoration_cost AS bposq__decoration_cost,
		bposq.making_cost AS bposq__making_cost,
		bposq.other_cost AS bposq__other_cost,
		bposq.sample_unit_price AS bposq__sample_unit_price,
		bposq.sample_lead_time AS bposq__sample_lead_time,
		bposq.seller_remark AS bposq__seller_remark,
		bposq.admin_sent_at AS bposq__admin_sent_at,
		bposq.quotation_at AS bposq__quotation_at,
		bposq.bulk_quotations AS bposq__bulk_quotations,
		bposq.expected_start_production_date AS bposq__expected_start_production_date,
		bposq.start_production_date AS bposq__start_production_date,
		bposq.capacity_per_day AS bposq__capacity_per_day,
		bposq.reject_reason AS bposq__reject_reason,
		bposq.expected_price AS bposq__expected_price,
		bposq.note AS bposq__note,
		bposq.admin_reject_reason AS bposq__admin_reject_reason

		FROM bulk_purchase_orders bpo
		LEFT JOIN inquiries iq ON bpo.inquiry_id = iq.id
		LEFT JOIN purchase_orders po ON bpo.purchase_order_id = po.id
		LEFT JOIN bulk_purchase_order_seller_quotations bposq ON bposq.bulk_purchase_order_id = bpo.id AND bposq.user_id = @user_id
		`
		countSQL = `
		SELECT /* {{Description}} */ 1

		FROM bulk_purchase_orders bpo
		LEFT JOIN inquiries iq ON bpo.inquiry_id = iq.id
		LEFT JOIN purchase_orders po ON bpo.purchase_order_id = po.id
		LEFT JOIN bulk_purchase_order_seller_quotations bposq ON bposq.bulk_purchase_order_id = bpo.id AND bposq.user_id = @user_id
		`
	}

	if options.Role.IsAdmin() {
		rawSQL = `
		SELECT /* {{Description}} */ bpo.*,
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
		iq.tax_percentage AS iq__tax_percentage,
		iq.product_weight AS iq__product_weight,
		iq.expected_price AS iq__expected_price,

		po.id AS po__id,
		po.created_at AS po__created_at,
		po.updated_at AS po__updated_at,
		po.deleted_at AS po__deleted_at,
		po.reference_id AS po__reference_id,
		po.client_reference_id AS po__client_reference_id,
		po.checkout_session_id AS po__checkout_session_id,
		po.user_id AS po__user_id,
		po.inquiry_id AS po__inquiry_id,
		po.status AS po__status,
		po.payment_intent_id AS po__payment_intent_id,
		po.charge_id AS po__charge_id,
		po.receipt_url AS po__receipt_url,
		po.payment_type AS po__payment_type,
		po.transfered_at AS po__transfered_at,
		po.mark_as_paid_at AS po__mark_as_paid_at,
		po.mark_as_unpaid_at AS po__mark_as_unpaid_at,
		po.transaction_ref_id AS po__transaction_ref_id,
		po.transaction_attachment AS po__transaction_attachment,
		po.payout_transaction_ref_id AS po__payout_transaction_ref_id,
		po.payout_transaction_attachment AS po__payout_transaction_attachment,
		po.payout_transfered_at AS po__payout_transfered_at,
		po.payout_mark_as_paid_at AS po__payout_mark_as_paid_at,
		po.attachments AS po__attachments,
		po.document AS po__document,
		po.design AS po__design,
		po.fabric_attachments AS po__fabric_attachments,
		po.techpack_attachments AS po__techpack_attachments,
		po.sample_attachments AS po__sample_attachments,
		po.approve_reject_meta AS po__approve_reject_meta,
		po.tracking_status AS po__tracking_status,
		po.po_raw_materials AS po__po_raw_materials,
		po.logistic_info AS po__logistic_info,
		po.making_info AS po__making_info,
		po.submit_info AS po__submit_info,
		po.receiver_confirmed_at AS po__receiver_confirmed_at,
		po.delivery_started_at AS po__delivery_started_at,
		po.assignee_ids AS po__assignee_ids,
		po.currency AS po__currency,
		po.feedback AS po__feedback,
		po.sample_maker_id AS po__sample_maker_id,
		po.seller_po_attachment AS po__seller_po_attachment,
		po.seller_design AS po__seller_design,
		po.seller_techpack_attachments AS po__seller_techpack_attachments,
		po.seller_est_making_at AS po__seller_est_making_at,
		po.seller_est_delivery_at AS po__seller_est_delivery_at,
		po.seller_submit_info AS po__seller_submit_info,
		po.seller_logistic_info AS po__seller_logistic_info,
		po.seller_delivery_started_at AS po__seller_delivery_started_at,
		po.seller_delivery_confirmed_at AS po__seller_delivery_confirmed_at,
		po.seller_delivery_feedback AS po__seller_delivery_feedback,
		po.invoice_number AS po__invoice_number,
		po.payment_transaction_reference_id AS po__payment_transaction_reference_id,
		po.refund_reason AS po__refund_reason,
		po.sub_total AS po__sub_total,
		po.sub_total_after_deduction AS po__sub_total_after_deduction,
		po.shipping_fee AS po__shipping_fee,
		po.transaction_fee AS po__transaction_fee,
		po.tax AS po__tax,
		po.total_price AS po__total_price,
		po.tax_percentage AS po__tax_percentage,
		po.seller_total_price AS po__seller_total_price,

		s.id AS s__id,
		s.name AS s__name,
		s.company_name AS s__company_name,
		s.payment_terms AS s__payment_terms,
		s.coordinate_id AS s__coordinate_id,
		s.country_code AS s__country_code,
		s.phone_number AS s__phone_number

		FROM bulk_purchase_orders bpo
		LEFT JOIN inquiries iq ON bpo.inquiry_id = iq.id
		LEFT JOIN purchase_orders po ON bpo.purchase_order_id = po.id
		LEFT JOIN users s ON s.id = bpo.seller_id
		`
		countSQL = `
		SELECT /* {{Description}} */ 1

		FROM bulk_purchase_orders bpo
		LEFT JOIN inquiries iq ON bpo.inquiry_id = iq.id
		LEFT JOIN purchase_orders po ON bpo.purchase_order_id = po.id
		LEFT JOIN users s ON s.id = bpo.seller_id
		`
	}

	return NewBuilder(rawSQL, countSQL).
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
			var records = make([]*models.BulkPurchaseOrder, rawSQL.RowsAffected)

			rows, err := rawSQL.Rows()
			if err != nil {
				db.CustomLogger.Errorf("Scan rows error", err)
				return nil, err
			}
			defer rows.Close()

			var bulkPurchaseOrderIDs []string
			var inquiryIDs []string
			var addressIDs []string
			var userIDs []string
			var assigneeIDs []string
			var invoiceNumbers []int
			var orderGroupIDs []string

			for rows.Next() {
				var alias BulkPurchaseOrderAlias
				err = db.ScanRows(rows, &alias)
				if err != nil {
					db.CustomLogger.Errorf("Scan rows error", err)
					continue
				}

				if !helper.StringContains(inquiryIDs, alias.BulkPurchaseOrder.InquiryID) {
					inquiryIDs = append(inquiryIDs, alias.BulkPurchaseOrder.InquiryID)
				}

				if !helper.StringContains(bulkPurchaseOrderIDs, alias.BulkPurchaseOrder.ID) {
					bulkPurchaseOrderIDs = append(bulkPurchaseOrderIDs, alias.BulkPurchaseOrder.ID)
				}

				if alias.ShippingAddressID != "" && !helper.StringContains(addressIDs, alias.ShippingAddressID) {
					addressIDs = append(addressIDs, alias.ShippingAddressID)
				}

				if alias.UserID != "" && !helper.StringContains(userIDs, alias.UserID) {
					userIDs = append(userIDs, alias.UserID)
				}
				if alias.OrderGroupID != "" && !helper.StringContains(orderGroupIDs, alias.OrderGroupID) {
					orderGroupIDs = append(orderGroupIDs, alias.OrderGroupID)
				}

				if alias.FirstPaymentInvoiceNumber != 0 {
					invoiceNumbers = append(invoiceNumbers, alias.FirstPaymentInvoiceNumber)
				}

				if alias.FinalPaymentInvoiceNumber != 0 {
					invoiceNumbers = append(invoiceNumbers, alias.FinalPaymentInvoiceNumber)
				}

				for _, v := range alias.AssigneeIDs {
					if !helper.StringContains(assigneeIDs, v) {
						assigneeIDs = append(assigneeIDs, v)
					}
				}

				if alias.Inquiry.ID != "" {
					alias.BulkPurchaseOrder.Inquiry = alias.Inquiry
				}

				if alias.PurchaseOrder.ID != "" {
					alias.BulkPurchaseOrder.PurchaseOrder = alias.PurchaseOrder
				}

				if alias.BulkQuotation != nil && alias.BulkQuotation.ID != "" {
					alias.BulkPurchaseOrder.SellerBulkQuotation = alias.BulkQuotation
				}

				if alias.Seller != nil && alias.Seller.ID != "" {
					alias.BulkPurchaseOrder.Seller = alias.Seller
				}

				alias.BulkPurchaseOrder.QuotedPrice = alias.BulkPurchaseOrder.GetQuotedPrice().ToPtr()
				alias.BulkPurchaseOrder.FinalPaymentTotal = alias.BulkPurchaseOrder.GetFinalPaymentAmount().ToPtr()
				alias.BulkPurchaseOrder.QuotationLeadTime = alias.BulkPurchaseOrder.GetQuotationLeadTime()

				records = append(records, alias.BulkPurchaseOrder)
			}

			var wg sync.WaitGroup

			if len(bulkPurchaseOrderIDs) > 0 && options.IncludeSellerQuotation {
				wg.Add(1)
				go func() {
					defer wg.Done()
					var items []*models.BulkPurchaseOrderSellerQuotation
					db.Find(&items, "bulk_purchase_order_id IN ?", bulkPurchaseOrderIDs)

					for _, item := range items {
						for _, record := range records {
							if record.ID == item.BulkPurchaseOrderID {
								record.SellerBulkQuotations = append(record.SellerBulkQuotations, item)
								if record.SellerID == item.UserID {
									record.SellerBulkQuotation = item
								}
							}
						}
					}

				}()
			}

			if len(bulkPurchaseOrderIDs) > 0 && options.IncludeItems {
				wg.Add(1)
				go func() {
					defer wg.Done()
					var items []*models.BulkPurchaseOrderItem
					db.Find(&items, "purchase_order_id IN ?", bulkPurchaseOrderIDs)

					for _, item := range items {
						for _, record := range records {
							if record.ID == item.PurchaseOrderID {
								record.Items = append(record.Items, item)
							}
						}
					}

				}()

				wg.Add(1)
				go func() {
					defer wg.Done()

					var orderCartItems []*models.OrderCartItem
					db.Find(&orderCartItems, "bulk_purchase_order_id IN ?", bulkPurchaseOrderIDs)

					for _, item := range orderCartItems {
						for _, record := range records {
							if record.ID == item.BulkPurchaseOrderID {
								record.OrderCartItems = append(record.OrderCartItems, item)
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

			if len(bulkPurchaseOrderIDs) > 0 && options.IncludeShippingAddress {
				wg.Add(1)
				go func() {
					defer wg.Done()
					var items []*models.PaymentTransaction
					db.Find(&items, "bulk_purchase_order_id IN ?", bulkPurchaseOrderIDs).Order("created_at desc")

					for _, item := range items {
						for _, record := range records {
							if record.ID == item.PurchaseOrderID {
								record.PaymentTransactions = append(record.PaymentTransactions, item)
							}
						}
					}
				}()
			}

			if len(userIDs) > 0 && options.IncludeUser {
				wg.Add(1)
				go func() {
					defer wg.Done()
					var users []*models.User
					db.Select("ID", "Name", "Avatar", "Email", "IsOffline", "LastOnlineAt").Find(&users, "id IN ?", userIDs)

					for _, user := range users {
						for _, record := range records {
							if record.UserID == user.ID {
								record.User = user
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
							if record.FirstPaymentInvoiceNumber == inv.InvoiceNumber {
								record.FirstPaymentInvoice = inv
							}

							if record.FinalPaymentInvoiceNumber == inv.InvoiceNumber {
								record.FinalPaymentInvoice = inv
							}
						}
					}
				}()
			}

			if options.IncludeTrackings && len(bulkPurchaseOrderIDs) > 0 {
				wg.Add(1)
				go func() {
					defer wg.Done()
					var trackings []*models.BulkPurchaseOrderTracking
					db.Find(&trackings, "purchase_order_id IN ?", bulkPurchaseOrderIDs)

					for _, tracking := range trackings {
						for _, record := range records {
							if record.ID == tracking.PurchaseOrderID {
								record.Trackings = append(record.Trackings, tracking)
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
			if options.IncludePaymentTransactions && len(bulkPurchaseOrderIDs) > 0 {
				wg.Add(1)
				go func() {
					defer wg.Done()

					var items []*models.PaymentTransaction
					if err := db.Find(&items, "count_elements(bulk_purchase_order_ids,?) >=1", pq.StringArray(bulkPurchaseOrderIDs)).Error; err != nil {
						return
					}

					for _, item := range items {
						for _, record := range records {
							if len(item.BulkPurchaseOrderIDs) > 0 && lo.Contains(item.BulkPurchaseOrderIDs, record.ID) {
								record.PaymentTransactions = append(record.PaymentTransactions, item)
							}
						}
					}

				}()

			}

			wg.Wait()

			return records, nil
		})
}
