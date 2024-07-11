package enums

type UserNotificationType string

var (
	UserNotificationTypeInquirySubmitQuotation  UserNotificationType = "inquiry_submit_quotation"
	UserNotificationTypeInquirySellerNewComment UserNotificationType = "inquiry_seller_new_comment"

	// Sample po
	UserNotificationTypePoCreated UserNotificationType = "po_created"
	// UserNotificationTypePoApproveDesign    UserNotificationType = "po_approve_design"
	// UserNotificationTypePoRejectDesign     UserNotificationType = "po_reject_design"
	UserNotificationTypePoUpdateDesign      UserNotificationType = "po_update_design"
	UserNotificationTypePoUpdateRawMaterial UserNotificationType = "po_update_raw_material"
	UserNotificationTypePoMarkMaking        UserNotificationType = "po_mark_making"
	UserNotificationTypePoMarkSubmit        UserNotificationType = "po_mark_submit"
	UserNotificationTypePoMarkDelivering    UserNotificationType = "po_mark_delivering"
	// UserNotificationTypePoConfirmDelivered  UserNotificationType = "po_confirm_delivered"
	UserNotificationTypePoDesignNewComment UserNotificationType = "po_design_new_comment"

	// Bulk po
	// UserNotificationTypeBulkPoSubmitOrder             UserNotificationType = "bulk_po_submit_order"
	UserNotificationTypeBulkPoSubmitQuotation UserNotificationType = "bulk_po_submit_quotation"
	// UserNotificationTypeBulkPoMakeFirstPayment        UserNotificationType = "bulk_po_make_first_payment"
	// UserNotificationTypeBulkPoMakeFinalPayment        UserNotificationType = "bulk_po_make_final_payment"
	UserNotificationTypeBulkPoUpdateRawMaterial UserNotificationType = "bulk_po_update_raw_material"
	UserNotificationTypeBulkPoMarkPps           UserNotificationType = "bulk_po_mark_pps"
	UserNotificationTypeBulkPoUpdatePps         UserNotificationType = "bulk_po_update_pps"
	UserNotificationTypeBulkPoMarkProduction    UserNotificationType = "bulk_po_mark_production"
	UserNotificationTypeBulkPoMarkQc            UserNotificationType = "bulk_po_mark_qc"
	UserNotificationTypeBulkPoCreateQcReport    UserNotificationType = "bulk_po_create_qc_report"
	UserNotificationTypeBulkPoMarkFirstPayment  UserNotificationType = "bulk_po_mark_first_payment"
	UserNotificationTypeBulkPoMarkFinalPayment  UserNotificationType = "bulk_po_mark_final_payment"
	// UserNotificationTypeBulkPoBuyerApproveRawMaterial UserNotificationType = "bulk_po_buyer_approve_raw_material"
	// UserNotificationTypeBulkPoBuyerApproveQc          UserNotificationType = "bulk_po_buyer_approve_qc"
	// UserNotificationTypeBulkPoConfirmDelivered      UserNotificationType = "bulk_po_confirm_delivered"
	// UserNotificationTypeBulkPoDelivered             UserNotificationType = "bulk_po_delivered"
	UserNotificationTypeBulkPoMarkDelivering  UserNotificationType = "bulk_po_mark_delivering"
	UserNotificationTypeBulkPoMarkRawMaterial UserNotificationType = "bulk_po_mark_raw_material"
	// UserNotificationTypeBulkPoFirstPaymentConfirmed UserNotificationType = "bulk_po_first_payment_confirmed"
	// UserNotificationTypeBulkPoFinalPaymentConfirmed UserNotificationType = "bulk_po_final_payment_confirmed"
)
