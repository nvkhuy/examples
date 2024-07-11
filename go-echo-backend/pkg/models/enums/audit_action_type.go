package enums

type AuditActionType string

var (
	AuditActionTypeLabelCreated  AuditActionType = "label_created"
	AuditActionTypeLabelEdited   AuditActionType = "label_edited"
	AuditActionTypeLabelRejected AuditActionType = "label_rejected"
	AuditActionTypeLabelApproved AuditActionType = "label_approved"

	AuditActionTypeInquiryCreated               AuditActionType = "inquiry_created" // RFQ submitted
	AuditActionTypeInquiryEdited                AuditActionType = "inquiry_edited"
	AuditActionTypeInquiryBuyerApproveQuotation AuditActionType = "inquiry_buyer_approve_quotation" // Quotation Approved
	AuditActionTypeInquiryBuyerRejectQuotation  AuditActionType = "inquiry_buyer_reject_quotation"  // Quotation Approved

	AuditActionTypeInquiryAdminSendSellerQuotation AuditActionType = "inquiry_admin_send_seller_quotation"
	AuditActionTypeInquiryAdminSendBuyerQuotation  AuditActionType = "inquiry_admin_send_buyer_quotation"

	AuditActionTypeInquirySellerSendQuotation AuditActionType = "inquiry_seller_send_quotation"

	AuditActionTypeInquiryAdminApproveSellerQuotation AuditActionType = "inquiry_admin_approve_seller_quotation"
	AuditActionTypeInquiryAdminRejectSellerQuotation  AuditActionType = "inquiry_admin_reject_seller_quotation"
	AuditActionTypeInquiryAdminMarkAsPaid             AuditActionType = "inquiry_admin_mark_as_paid"
	AuditActionTypeInquiryAdminMarkAsUnPaid           AuditActionType = "inquiry_admin_mark_as_un_paid"

	AuditActionTypeInquirySamplePoCreated AuditActionType = "inquiry_sample_po_created"
	AuditActionTypeInquiryBulkPoCreated   AuditActionType = "inquiry_bulk_po_created"

	AuditActionTypeInquirySamplePoTrackingStatus AuditActionType = "inquiry_sample_po_tracking_status"
	AuditActionTypeInquiryBulkPoTrackingStatus   AuditActionType = "inquiry_bulk_po_tracking_status"

	AuditActionTypeBulkPoFirstPaymentAdminMarkAsPaid   AuditActionType = "bulk_po_first_payment_admin_mark_as_paid"
	AuditActionTypeBulkPoFinalPaymentAdminMarkAsUnPaid AuditActionType = "bulk_po_final_payment_as_un_paid"

	AuditActionTypeInquiryAdminRefund AuditActionType = "inquiry_admin_refund"

	AuditActionTypeInquiryAdminUpdateCosting AuditActionType = "inquiry_admin_update_costing"
)
