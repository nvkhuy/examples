package models

import (
	"github.com/engineeringinflow/inflow-backend/pkg/models/enums"
	"github.com/engineeringinflow/inflow-backend/pkg/models/price"
	"github.com/lib/pq"
	"github.com/stripe/stripe-go/v74"
)

type BulkPurchaseOrder struct {
	Model

	ProductName string `gorm:"sze:1000" json:"product_name,omitempty"`
	Note        string `gorm:"sze:3000" json:"note,omitempty"`

	Items          []*BulkPurchaseOrderItem `gorm:"-" json:"items,omitempty"` //legacy
	OrderCartItems []*OrderCartItem         `gorm:"-" json:"order_cart_items,omitempty"`

	AdditionalItems BulkPurchaseOrderAdditionalItems `json:"additional_items,omitempty"`

	ReferenceID string `gorm:"unique;size:100" json:"reference_id"`

	User   *User  `gorm:"-" json:"user,omitempty"`
	UserID string `gorm:"size:100" json:"user_id,omitempty"`

	PurchaseOrderID string         `gorm:"size:100" json:"purchase_order_id,omitempty"`
	PurchaseOrder   *PurchaseOrder `gorm:"-" json:"purchase_order,omitempty"`

	InquiryID string   `gorm:"size:100" json:"inquiry_id,omitempty"`
	Inquiry   *Inquiry `gorm:"-" json:"inquiry,omitempty"`

	ShippingMethod      enums.ShippingMethod `json:"shipping_method,omitempty"`
	ShippingAttachments *Attachments         `json:"shipping_attachments,omitempty"`
	ShippingNote        *string              `gorm:"size:2000" json:"shipping_note,omitempty"`

	PackingNote             *string      `gorm:"size:2000" json:"packing_note,omitempty"`
	PackingAttachments      *Attachments `json:"packing_attachments,omitempty"`
	Attachments             *Attachments `json:"attachments,omitempty"`
	TechpackAttachments     *Attachments `json:"techpack_attachments,omitempty"`
	SizeAttachments         *Attachments `json:"size_attachments,omitempty"`
	TableColorSizeNote      *string      `json:"table_color_size_note,omitempty"`
	UploadPOAttachments     *Attachments `json:"upload_po_note_attachments,omitempty"`
	UploadPOAttachmentsNote *string      `gorm:"size:2000" json:"upload_po_attachments_note,omitempty"`

	AdditionalRequirements string `gorm:"size:2000" json:"additional_requirements,omitempty"`

	Currency enums.Currency             `gorm:"default:'USD'" json:"currency"`
	Feedback *BulkPurchaseOrderFeedback `json:"feedback,omitempty"`

	Pricing

	ShippingAddressID string   `gorm:"size:100" json:"shipping_address_id,omitempty"`
	ShippingAddress   *Address `gorm:"-" json:"shipping_address,omitempty"`

	DepositInvoiceNumber                 int          `json:"deposit_invoice_number,omitempty"`
	DepositInvoice                       *Invoice     `gorm:"-" json:"deposit_invoice,omitempty"`
	DepositAmount                        *price.Price `gorm:"type:decimal(20,4);default:0.0" json:"deposit_amount,omitempty"`
	DepositPaidAmount                    *price.Price `gorm:"type:decimal(20,4);default:0.0" json:"deposit_paid_amount,omitempty"`
	DepositTransactionFee                *price.Price `gorm:"type:decimal(20,4);default:0.0" json:"deposit_transaction_fee,omitempty"`
	DepositNote                          string       `json:"deposit_note,omitempty"`
	DepositPaymentIntentID               string       `gorm:"size:100" json:"deposit_payment_intent_id,omitempty"`
	DepositChargeID                      string       `gorm:"size:100" json:"deposit_charge_id,omitempty"`
	DepositTxnID                         string       `gorm:"size:100" json:"deposit_txn_id,omitempty"`
	DepositTransactionRefID              string       `gorm:"size:100" json:"deposit_transaction_ref_id,omitempty"`
	DepositTransactionAttachment         *Attachment  `json:"deposit_transaction_attachment,omitempty"`
	DepositPaymentLink                   string       `gorm:"size:2000" json:"deposit_payment_link,omitempty"`
	DepositPaymentLinkID                 string       `gorm:"size:100" json:"deposit_payment_link_id,omitempty"`
	DepositTransferedAt                  *int64       `json:"deposit_transfered_at,omitempty"`
	DepositMarkAsPaidAt                  *int64       `json:"deposit_mark_as_paid_at,omitempty"`
	DepositReceiptURL                    string       `gorm:"size:2000" json:"deposit_receipt_url,omitempty"`
	DepositPaymentTransactionReferenceID string       `gorm:"size:100" json:"deposit_paymnet_transaction_reference_id,omitempty"`

	FirstPaymentInvoiceNumber          int               `json:"first_payment_invoice_number,omitempty"`
	FirstPaymentInvoice                *Invoice          `gorm:"-" json:"first_payment_invoice,omitempty"`
	FirstPaymentType                   enums.PaymentType `gorm:"size:50;not null;default:'bank_transfer'" json:"first_payment_type,omitempty"`
	FirstPaymentTransactionFee         *price.Price      `gorm:"type:decimal(20,4);default:0.0" json:"first_payment_transaction_fee,omitempty"`
	FirstPaymentTax                    *price.Price      `gorm:"type:decimal(20,4);default:0.0" json:"first_payment_tax,omitempty"`
	FirstPaymentReceivedAt             *int64            `json:"first_payment_received_at,omitempty"`
	FirstPaymentSubTotal               *price.Price      `gorm:"type:decimal(20,4);default:0.0" json:"first_payment_sub_total,omitempty"`
	FirstPaymentTotal                  *price.Price      `gorm:"type:decimal(20,4);default:0.0" json:"first_payment_total,omitempty"`
	FirstPaymentPercentage             *float64          `gorm:"default:40.0" json:"first_payment_percentage,omitempty"`
	FirstPaymentTransferedAt           *int64            `json:"first_payment_transfered_at,omitempty"`
	FirstPaymentMarkAsPaidAt           *int64            `json:"first_payment_mark_as_paid_at,omitempty"`
	FirstPaymentMarkAsUnpaidAt         *int64            `json:"first_payment_mark_as_unpaid_at,omitempty"`
	FirstPaymentIntentID               string            `gorm:"size:100" json:"first_payment_intent_id,omitempty"`
	FirstPaymentChargeID               string            `gorm:"size:100" json:"first_payment_charge_id,omitempty"`
	FirstPaymentTxnID                  string            `gorm:"size:100" json:"first_payment_txn_id,omitempty"`
	FirstPaymentReceiptURL             string            `gorm:"size:2000" json:"first_payment_receipt_url,omitempty"`
	FirstPaymentTransactionRefID       string            `gorm:"size:100" json:"first_payment_transaction_ref_id,omitempty"`
	FirstPaymentTransactionAttachment  *Attachment       `json:"first_payment_transaction_attachment,omitempty"`
	FirstPaymentLink                   string            `gorm:"size:2000" json:"first_payment_link,omitempty"`
	FirstPaymentLinkID                 string            `gorm:"size:100" json:"first_payment_link_id,omitempty"`
	FirstPaymentTransactionReferenceID string            `gorm:"size:100" json:"first_payment_transaction_reference_id,omitempty"`
	FirstPaymentCheckoutSessionID      string            `gorm:"size:100" json:"first_payment_checkout_session_id,omitempty"`

	SecondPaymentInvoiceNumber          int               `json:"second_payment_invoice_number,omitempty"`
	SecondPaymentInvoice                *Invoice          `gorm:"-" json:"second_payment_invoice,omitempty"`
	SecondPaymentType                   enums.PaymentType `gorm:"size:50;not null;default:'bank_transfer'" json:"second_payment_type,omitempty"`
	SecondPaymentTransactionFee         *price.Price      `gorm:"type:decimal(20,4);default:0.0" json:"second_payment_transaction_fee,omitempty"`
	SecondPaymentTax                    *price.Price      `gorm:"type:decimal(20,4);default:0.0" json:"second_payment_tax,omitempty"`
	SecondPaymentReceivedAt             *int64            `json:"second_payment_received_at,omitempty"`
	SecondPaymentSubTotal               *price.Price      `gorm:"type:decimal(20,4);default:0.0" json:"second_payment_sub_total,omitempty"`
	SecondPaymentTotal                  *price.Price      `gorm:"type:decimal(20,4);default:0.0" json:"second_payment_total,omitempty"`
	SecondPaymentPercentage             float64           `gorm:"default:0.0" json:"second_payment_percentage,omitempty"`
	SecondPaymentTransferedAt           *int64            `json:"second_payment_transfered_at,omitempty"`
	SecondPaymentMarkAsPaidAt           *int64            `json:"second_payment_mark_as_paid_at,omitempty"`
	SecondPaymentMarkAsUnpaidAt         *int64            `json:"second_payment_mark_as_unpaid_at,omitempty"`
	SecondPaymentIntentID               string            `gorm:"size:100" json:"second_payment_intent_id,omitempty"`
	SecondPaymentChargeID               string            `gorm:"size:100" json:"second_payment_charge_id,omitempty"`
	SecondPaymentTxnID                  string            `gorm:"size:100" json:"second_payment_txn_id,omitempty"`
	SecondPaymentReceiptURL             string            `gorm:"size:2000" json:"second_payment_receipt_url,omitempty"`
	SecondPaymentTransactionRefID       string            `gorm:"size:100" json:"second_payment_transaction_ref_id,omitempty"`
	SecondPaymentTransactionAttachment  *Attachment       `json:"second_payment_transaction_attachment,omitempty"`
	SecondPaymentLink                   string            `gorm:"size:2000" json:"second_payment_link,omitempty"`
	SecondPaymentLinkID                 string            `gorm:"size:100" json:"second_payment_link_id,omitempty"`
	SecondPaymentTransactionReferenceID string            `gorm:"size:100" json:"second_payment_transaction_reference_id,omitempty"`

	FinalPaymentInvoiceNumber          int               `json:"final_payment_invoice_number,omitempty"`
	FinalPaymentInvoice                *Invoice          `gorm:"-" json:"final_payment_invoice,omitempty"`
	FinalPaymentType                   enums.PaymentType `gorm:"size:50;not null;default:'bank_transfer'" json:"final_payment_type,omitempty"`
	FinalPaymentTransactionFee         *price.Price      `gorm:"type:decimal(20,4);default:0.0" json:"final_payment_transaction_fee,omitempty"`
	FinalPaymentTax                    *price.Price      `gorm:"type:decimal(20,4);default:0.0" json:"final_payment_tax,omitempty"`
	FinalPaymentSubTotal               *price.Price      `gorm:"type:decimal(20,4);default:0.0" json:"final_payment_sub_total,omitempty"`
	FinalPaymentIntentID               string            `gorm:"size:100" json:"final_payment_intent_id,omitempty"`
	FinalPaymentChargeID               string            `json:"final_payment_charge_id,omitempty"`
	FinalPaymentTxnID                  string            `gorm:"size:100" json:"final_payment_txn_id,omitempty"`
	FinalPaymentReceiptURL             string            `gorm:"size:2000" json:"final_payment_receipt_url,omitempty"`
	FinalPaymentReceivedAt             *int64            `json:"final_payment_received_at,omitempty"`
	FinalPaymentTotal                  *price.Price      `json:"final_payment_total,omitempty"`
	FinalPaymentTransferedAt           *int64            `json:"final_payment_transfered_at,omitempty"`
	FinalPaymentMarkAsPaidAt           *int64            `json:"final_payment_mark_as_paid_at,omitempty"`
	FinalPaymentMarkAsUnpaidAt         *int64            `json:"final_payment_mark_as_unpaid_at,omitempty"`
	FinalPaymentTransactionRefID       string            `gorm:"size:100" json:"final_payment_transaction_ref_id"`
	FinalPaymentTransactionAttachment  *Attachment       `json:"final_payment_transaction_attachment,omitempty"`
	FinalPaymentDeductionAmount        *price.Price      `gorm:"type:decimal(20,4);default:0.0" json:"final_payment_deduction_amount,omitempty"`
	FinalPaymentLink                   string            `gorm:"size:2000" json:"final_payment_link,omitempty"`
	FinalPaymentLinkID                 string            `gorm:"size:100" json:"final_payment_link_id,omitempty"`
	FinalPaymentTransactionReferenceID string            `gorm:"size:100" json:"final_payment_transaction_reference_id,omitempty"`
	FinalPaymentCheckoutSessionID      string            `gorm:"size:100" json:"final_payment_checkout_session_id,omitempty"`

	SampleDeductionAmount *price.Price `gorm:"type:decimal(20,4);default:0.0" json:"sample_deduction_amount,omitempty"`

	Status         enums.BulkPurchaseOrderStatus `gorm:"size:50;default:'new'" json:"status,omitempty"`
	TrackingStatus enums.BulkPoTrackingStatus    `gorm:"size:50;default:'new'" json:"tracking_status,omitempty"`

	QuotedPrice              *price.Price          `gorm:"-" json:"quoted_price,omitempty"`
	QuotationLeadTime        *int64                `gorm:"-"  json:"quotation_lead_time,omitempty"`
	QuotationAt              *int64                `json:"quotation_at,omitempty"`
	QuotationNote            string                `json:"quotation_note,omitempty"`
	QuotationNoteAttachments *Attachments          `json:"quotation_note_attachments,omitempty"`
	AdminQuotations          InquiryQuotationItems `json:"admin_quotations,omitempty"`

	StartDate      *int64 `json:"start_date,omitempty"`
	LeadTime       int    `json:"lead_time,omitempty"`
	CompletionDate *int64 `json:"completion_date,omitempty"`

	PoRawMaterials       *PoRawMaterialMetas `json:"po_raw_materials,omitempty"`
	PpsInfo              *PoPpsMetas         `json:"pps_info,omitempty"`        // pre production step
	ProductionInfo       *PoProductionMeta   `json:"production_info,omitempty"` //  production step
	PoQcReports          *PoReportMetas      `json:"po_qc_reports,omitempty"`
	ApproveRawMaterialAt *int64              `json:"approve_raw_material_at,omitempty"`
	ApproveQCAt          *int64              `json:"approve_qc_at,omitempty"`

	PaymentTransactions []*PaymentTransaction `gorm:"-" json:"payment_transactions,omitempty"`

	LogisticInfo        *PoLogisticMeta `json:"logistic_info,omitempty"`
	ReceiverConfirmedAt *int64          `json:"receiver_confirmed_at,omitempty"`
	DeliveryStartedAt   *int64          `json:"delivery_started_at,omitempty"`
	SubmittedAt         *int64          `json:"submitted_at,omitempty"`
	DeliveredAt         *int64          `json:"delivered_at,omitempty"`

	AssigneeIDs pq.StringArray `gorm:"type:varchar(200)[]" json:"assignee_ids,omitempty"`
	Assignees   Users          `gorm:"-" json:"assignees,omitempty"`

	PaymentIntentNextAction   *stripe.PaymentIntentNextAction `gorm:"-" json:"payment_intent_next_action,omitempty"`
	PaymentIntentClientSecret string                          `gorm:"-" json:"payment_intent_client_secret,omitempty"`

	DebitNoteAttachment         *Attachment                         `json:"debit_note_attachment,omitempty"`
	CommercialInvoiceAttachment *Attachment                         `json:"commercial_invoice_attachment,omitempty"`
	CommercialInvoice           *BulkPurchaseOrderCommercialInvoice `json:"commercial_invoice,omitempty"`

	Trackings []*BulkPurchaseOrderTracking `gorm:"-" json:"trackings,omitempty"`

	OrderGroupID string      `json:"order_group_id,omitempty"`
	OrderGroup   *OrderGroup `gorm:"-" json:"order_group,omitempty"`

	HubspotDealID string `gorm:"size:100"  json:"hubspot_deal_id,omitempty"`

	// Seller
	SellerTrackingStatus enums.SellerBulkPoTrackingStatus `gorm:"size:200" json:"seller_tracking_status,omitempty"`
	Seller               *User                            `gorm:"-" json:"seller,omitempty"`
	SellerID             string                           `gorm:"size:200" json:"seller_id,omitempty"`

	SellerPoRawMaterials *PoRawMaterialMetas `json:"seller_po_raw_materials,omitempty"`
	SellerPpsInfo        *PoPpsMetas         `json:"seller_pps_info,omitempty"`        // pre production step
	SellerProductionInfo *PoProductionMeta   `json:"seller_production_info,omitempty"` //  production step
	SellerPoQcReports    *PoReportMetas      `json:"seller_po_qc_reports,omitempty"`

	SellerQuotationRejectReason string                              `gorm:"size:2000" json:"seller_quotation_reject_reason,omitempty"`
	SellerQuotationSubmittedAt  *int64                              `json:"seller_quotation_submitted_at,omitempty"`
	SellerQuotationApprovedAt   *int64                              `json:"seller_quotation_approved_at,omitempty"`
	SellerQuotationStatus       enums.BulkPurchaseOrderSellerStatus `json:"seller_quotation_status,omitempty"`

	SellerFirstPayoutPercentage             *float64     `json:"seller_first_payout_percentage,omitempty"`
	SellerFirstPayoutTotalAmount            *price.Price `json:"seller_first_payout_total_amount,omitempty"`
	SellerFirstPayoutTransactionRefID       string       `gorm:"size:200" json:"seller_first_payout_transaction_ref_id,omitempty"`
	SellerFirstPayoutTransactionReferenceID string       `gorm:"size:200" json:"seller_first_payout_transaction_reference_id,omitempty"`
	SellerFirstPayoutTransactionAttachment  *Attachment  `json:"seller_first_payout_transaction_attachment,omitempty"`
	SellerFirstPayoutTransferedAt           *int64       `json:"seller_first_payout_transfered_at,omitempty"`
	SellerFirstPayoutMarkAsPaidAt           *int64       `json:"seller_first_payout_mark_as_paid_at,omitempty"`

	SellerFinalPayoutTransactionRefID       string       `gorm:"size:200" json:"seller_final_payout_transaction_ref_id,omitempty"`
	SellerFinalPayoutTotalAmount            *price.Price `json:"seller_final_payout_total_amount,omitempty"`
	SellerFinalPayoutTransactionReferenceID string       `gorm:"size:200" json:"seller_final_payout_transaction_reference_id,omitempty"`
	SellerFinalPayoutTransactionAttachment  *Attachment  `json:"seller_final_payout_transaction_attachment,omitempty"`
	SellerFinalPayoutTransferedAt           *int64       `json:"seller_final_payout_transfered_at,omitempty"`
	SellerFinalPayoutMarkAsPaidAt           *int64       `json:"seller_final_payout_mark_as_paid_at,omitempty"`

	SellerPayoutTotalAmount *price.Price `json:"seller_payout_total_amount,omitempty"`

	SellerSizeChartAttachments                     *Attachments         `json:"seller_size_chart_attachments,omitempty"`
	SellerSizeSpecAttachments                      *Attachments         `json:"seller_size_spec_attachments,omitempty"`
	SellerSizeGradingAttachments                   *Attachments         `json:"seller_size_grading_attachments,omitempty"`
	SellerLabelGuideAttachments                    *Attachments         `json:"seller_label_guide_attachments,omitempty"`
	SellerPointOfMeasurementAttachments            *Attachments         `json:"seller_point_of_measurement_attachments,omitempty"`
	SellerInspectionProcedureAttachments           *Attachments         `json:"seller_inspection_procedure,omitempty"`
	SellerInspectionProcedureNote                  string               `json:"seller_inspection_procedure_note,omitempty"`
	SellerInspectionTestingRequirementsAttachments *Attachments         `json:"seller_inspection_testing_requirements,omitempty"`
	SellerInspectionTestingRequirementsNote        string               `json:"seller_inspection_testing_requirements_note,omitempty"`
	SellerPackingNote                              string               `json:"seller_packing_note,omitempty"`
	SellerPackingAttachments                       *Attachments         `json:"seller_packing_attachments,omitempty"`
	SellerShippngMethod                            enums.ShippingMethod `json:"seller_shippng_method,omitempty"`
	SellerShippingAttachments                      *Attachments         `json:"seller_shipping_attachments,omitempty"`
	SellerShippingNote                             *string              `json:"seller_shipping_note,omitempty"`

	SellerPoAttachments             *PoAttachments `json:"seller_po_attachments,omitempty"`
	SellerProductAttachments        *Attachments   `json:"seller_product_attachments,omitempty"`
	SellerTechpackAttachments       *Attachments   `json:"seller_techpack_attachments,omitempty"`
	SellerBillOfMaterialAttachments *Attachments   `json:"seller_bill_of_material,omitempty"`

	SellerBulkQuotation *BulkPurchaseOrderSellerQuotation `gorm:"-" json:"seller_bulk_quotation,omitempty"`

	SellerBulkQuotations []*BulkPurchaseOrderSellerQuotation `gorm:"-" json:"seller_bulk_quotations,omitempty"`
	SellerFeedback       *BulkPurchaseOrderFeedback          `json:"seller_feedback,omitempty"`

	SellerLogisticInfo      *PoLogisticMeta `json:"seller_logistic_info,omitempty"`
	SellerDeliveryStartedAt *int64          `json:"seller_delivery_started_at,omitempty"`
	SellerDeliveredAt       *int64          `json:"seller_delivered_at,omitempty"`
}

type BulkPurchaseOrders []*BulkPurchaseOrder

type BulkPurchaseOrderUpdateForm struct {
	JwtClaimsInfo

	BulkPurchaseOrderID string `param:"bulk_purchase_order_id" validate:"required"`

	UserID string `json:"user_id,omitempty"`

	Items []*BulkPurchaseOrderItem `gorm:"-" json:"items,omitempty"`

	ShippingMethod      string       `json:"shipping_method,omitempty"`
	ShippingAttachments *Attachments `json:"shipping_attachments,omitempty"`
	ShippingNote        *string      `json:"shipping_note"`

	PackingNote        *string      `json:"packing_note"`
	PackingAttachments *Attachments `json:"packing_attachments,omitempty"`

	AdditionalItems        BulkPurchaseOrderAdditionalItems `json:"additional_items"`
	AdditionalRequirements string                           `json:"additional_requirements"`

	Attachments             *Attachments `json:"attachments,omitempty"`
	TechpackAttachments     *Attachments `json:"techpack_attachments"`
	SizeAttachments         *Attachments `json:"size_attachments"`
	TableColorSizeNote      *string      `json:"table_color_size_note"`
	UploadPOAttachments     *Attachments `json:"upload_po_note_attachments"`
	UploadPOAttachmentsNote *string      `json:"upload_po_attachments_note"`

	BankTransferInfos BankTransferInfos `gorm:"-" json:"bank_transfer_infos,omitempty"`

	ShippingAddress   *Address `json:"shipping_address,omitempty"`
	ShippingAddressID string   `gorm:"size:100" json:"shipping_address_id,omitempty"`

	PoRawMaterials *PoRawMaterialMetas `json:"po_raw_materials,omitempty"`

	LogisticInfo        *PoLogisticMeta `json:"logistic_info,omitempty"`
	ReceiverConfirmedAt *int64          `json:"receiver_confirmed_at,omitempty"`
	DeliveryStartedAt   *int64          `json:"delivery_started_at,omitempty"`

	SellerSizeChartAttachments              *Attachments         `json:"seller_size_chart_attachments"`
	SellerSizeSpecAttachments               *Attachments         `json:"seller_size_spec_attachments"`
	SellerSizeGradingAttachments            *Attachments         `json:"seller_size_grading_attachments"`
	SellerLabelGuideAttachments             *Attachments         `json:"seller_label_guide_attachments"`
	SellerPointOfMeasurementAttachments     *Attachments         `json:"seller_point_of_measurement_attachments"`
	SellerInspectionProcedure               *Attachments         `json:"seller_inspection_procedure"`
	SellerInspectionProcedureNote           string               `json:"seller_inspection_procedure_note"`
	SellerInspectionTestingRequirements     *Attachments         `json:"seller_inspection_testing_requirements"`
	SellerInspectionTestingRequirementsNote string               `json:"seller_inspection_testing_requirements_note"`
	SellerPackingNote                       *string              `json:"seller_packing_note"`
	SellerPackingAttachments                *Attachments         `json:"seller_packing_attachments,omitempty"`
	SellerShippngMethod                     enums.ShippingMethod `json:"seller_shippng_method,omitempty"`
	SellerShippingAttachments               *Attachments         `json:"seller_shipping_attachments,omitempty"`
	SellerShippingNote                      *string              `gorm:"size:2000" json:"seller_shipping_note"`
}

type BulkPurchaseOrderCreateForm struct {
	JwtClaimsInfo
	InquiryID string `json:"inquiry_id,omitempty" param:"inquiry_id" query:"inquiry_id" validate:"required"`
}

type BulkPurchaseOrderIDParam struct {
	JwtClaimsInfo

	BulkPurchaseOrderID string `json:"bulk_purchase_order_id" query:"bulk_purchase_order_id" param:"bulk_purchase_order_id" validate:"required"`
}

type BulkPurchaseOrderMarkAsPaidParams struct {
	JwtClaimsInfo

	BulkPurchaseOrderID string                 `json:"bulk_purchase_order_id" query:"bulk_purchase_order_id" param:"bulk_purchase_order_id" validate:"required"`
	Milestone           enums.PaymentMilestone `json:"milestone" query:"milestone" param:"milestone"`
}

type BulkPurchaseOrderMarkAsUnpaidParams struct {
	JwtClaimsInfo

	BulkPurchaseOrderID string                 `json:"bulk_purchase_order_id" query:"bulk_purchase_order_id" param:"bulk_purchase_order_id" validate:"required"`
	Milestone           enums.PaymentMilestone `json:"milestone" query:"milestone" param:"milestone"`
	Reason              string                 `json:"reason"`
}

type BulkPurchaseOrderAssignPICParam struct {
	JwtClaimsInfo
	BulkPurchaseOrderID string   `json:"bulk_purchase_order_id" query:"bulk_purchase_order_id" param:"bulk_purchase_order_id" validate:"required"`
	AssigneeIDs         []string `json:"assignee_ids" query:"assignee_ids" param:"assignee_ids" validate:"required"`
}

type CreateBulkPurchaseOrderRequest struct {
	ProductName         string                `json:"product_name,omitempty" validate:"required"`
	Note                string                `json:"note,omitempty"`
	Attachments         *Attachments          `json:"attachments,omitempty" validate:"required"`
	Items               OrderCartItems        `json:"items,omitempty" validate:"required"`
	TechpackAttachments *Attachments          `json:"techpack_attachments,omitempty" validate:"required"`
	SizeAttachments     *Attachments          `json:"size_attachments,omitempty" validate:"required"`
	OrderGroupID        string                `json:"order_group_id,omitempty"`
	Quotations          InquiryQuotationItems `json:"quotations" validate:"required"`
	// purchase orders
	PurchaseOrderClientReferenceID string         `json:"purchase_order_client_reference_id,omitempty"`
	PurchaseOrderItems             OrderCartItems `json:"purchase_order_items,omitempty"`
	PurchaseOrderTaxPercentage     float64        `json:"purchase_order_tax_percentage" validate:"min=0,max=100"`
	PurchaseOrderShippingFee       price.Price    `json:"purchase_order_shipping_fee"`
}

type SendBulkPurchaseOrderQuotationParams struct {
	JwtClaimsInfo

	FirstPaymentPercentage float64               `json:"first_payment_percentage" validate:"min=0,max=100"`
	AdminQuotations        InquiryQuotationItems `json:"admin_quotations" param:"admin_quotations" query:"admin_quotations" validate:"required"`

	QuotationNote            string       `json:"quotation_note" param:"quotation_note" query:"quotation_note" validate:"required"`
	QuotationNoteAttachments *Attachments `json:"quotation_note_attachments" param:"quotation_note_attachments" query:"quotation_note_attachments"`
	BulkPurchaseOrderID      string       `json:"bulk_purchase_order_id" param:"bulk_purchase_order_id" query:"bulk_purchase_order_id" validate:"required"`
}
type CreateMultipleBulkPurchaseOrdersRequest struct {
	JwtClaimsInfo
	UserID                 string                            `json:"user_id,omitempty" validate:"required"`
	ShippingMethod         enums.ShippingMethod              `json:"shipping_method,omitempty" validate:"required,oneof=fob cif exw"`
	ShippingAttachments    Attachments                       `json:"shipping_attachments,omitempty" validate:"required"`
	ShippingNote           string                            `json:"shipping_note,omitempty"`
	PackingAttachments     Attachments                       `json:"packing_attachments,omitempty" validate:"required"`
	PackingNote            string                            `json:"packing_note,omitempty"`
	AddressCoordinate      *Coordinate                       `json:"address_coordinate,omitempty"`
	Currency               enums.Currency                    `json:"currency,omitempty" validate:"required,oneof=USD VND"`
	TaxPercentage          float64                           `json:"tax_percentage" validate:"min=0,max=100"`
	ShippingFee            price.Price                       `json:"shipping_fee"`
	FirstPaymentPercentage float64                           `json:"first_payment_percentage" validate:"min=0,max=100"`
	Bulks                  []*CreateBulkPurchaseOrderRequest `json:"bulks,omitempty" validate:"required"`
}

type SubmitMultipleBulkQuotationsRequest struct {
	JwtClaimsInfo
	Quotations []*SendBulkPurchaseOrderQuotationParams `json:"quotations" validate:"required"`
}

type UploadBulksRequest struct {
	JwtClaimsInfo
	FileKey string `json:"file_key" validate:"required"`
}

type UploadBulksResponse struct {
	CreateMultipleBulkPurchaseOrdersRequest
}
