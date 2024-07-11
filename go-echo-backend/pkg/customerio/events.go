package customerio

type Event string

var (
	EventTrackActivity             Event = "track_activity"
	EventConfirmMail               Event = "confirm_mail"
	EventResetPassword             Event = "reset_password"
	EventInviteBrandMember         Event = "invite_brand_member"
	EventStaffInvitation           Event = "staff_invitation"
	EventWelcomeToBoard            Event = "welcome_to_board"
	EventOnboardSeller             Event = "onboard_seller"
	EventNotifyUserApproved        Event = "notify_user_approved"
	EventNotifyUserRejected        Event = "notify_user_rejected"
	EventNotifySellerApproved      Event = "notify_seller_approved"
	EventNotifySellerRejected      Event = "notify_seller_reject"
	EventBuyerApproveSkuQuotation  Event = "buyer_approve_sku_quotation"
	EventBuyerRejectSkuQuotation   Event = "buyer_reject_sku_quotation"
	EventAdminSentQuotationToBuyer Event = "admin_sent_quotation_to_buyer"
	EventEmailVerification         Event = "email_verification"

	EventSellerQuotationApproved Event = "seller_quotation_approved"
	EventSellerQuotationRejected Event = "seller_quotation_rejected"
	EventNewSubscriber           Event = "new_subscriber"

	EventNewInquiry                      Event = "new_inquiry"
	EventNewInquiryRemindAdmin           Event = "new_inquiry_remind_admin"
	EventAdminMarkAsPaid                 Event = "admin_mark_as_paid"
	EventAdminMarkAsUnpaid               Event = "admin_mark_as_unpaid"
	EventAdminInviteNewUser              Event = "admin_invite_new_user"
	EventAdminInviteClient               Event = "admin_invite_client"
	EventAdminInquiryAssignPIC           Event = "admin_inquiry_assign_pic"
	EventAdminPurchaseOrderAssignPIC     Event = "admin_purchase_order_assign_pic"
	EventAdminBulkPurchaseOrderAssignPIC Event = "admin_bulk_purchase_order_assign_pic"
	EventAdminSendTNANotification        Event = "admin_send_tna_notification"

	EventSellerNewRFQRequest                 Event = "seller_new_rfq_request"
	EventSellerNewBulkPurchaseOrderQuotation Event = "seller_new_bulk_purchase_order_quotation"
	// Event for sample po
	EventPoBuyerPaymentSucceeded                   Event = "po_buyer_payment_succeeded"
	EventPoBuyerPaymentFailed                      Event = "po_buyer_payment_failed"
	EventPoWaitingConfirmBankTransfer              Event = "po_waiting_confirm_bank_transfer"
	EventPoMultipleItemsWaitingConfirmBankTransfer Event = "po_multiple_items_waiting_confirm_bank_transfer"

	EventPoApproveDesign      Event = "po_approve_design"
	EventPoBuyerApproveDesign Event = "po_buyer_approve_design"

	EventPurchaseOrderBankTransferConfirmed Event = "po_bank_transfer_confirmed"

	EventPoBuyerUpdated            Event = "po_buyer_updated"
	EventPoBuyerApproveRawMaterial Event = "po_buyer_approve_raw_material"
	EventPoCreated                 Event = "po_created"
	EventPoWorkingOnDesign         Event = "po_working_on_design"
	EventPoRejectDesign            Event = "po_reject_design"
	EventPoUpdateDesign            Event = "po_update_design"
	EventPoUpdateMaterial          Event = "po_update_raw_material"
	EventPoApproveRawMaterial      Event = "po_approve_raw_material"
	EventPoMarkMaking              Event = "po_mark_making"
	EventPoMarkSubmit              Event = "po_mark_submit"
	EventPoMarkDelivering          Event = "po_mark_delivering"
	EventPoConfirmDelivered        Event = "po_confirm_delivered"
	EventPoConfirmed               Event = "po_confirmed"
	EventPoBuyerRejectDesign       Event = "po_buyer_reject_design"
	EventPoBuyerConfirmDelivered   Event = "po_buyer_confirm_delivered"
	EventPoNewComment              Event = "po_new_comment"

	EventAdminPoUpdateDesign Event = "admin_po_update_design"

	// Event for bulk po
	EventBulkPoSubmitOrder                            Event = "bulk_po_submit_order"
	EventBulkPoSubmitQuotation                        Event = "bulk_po_submit_quotation"
	EventBulkPoMakeFirstPayment                       Event = "bulk_po_make_first_payment"
	EventBulkPoMakeFinalPayment                       Event = "bulk_po_make_final_payment"
	EventBulkPoCreateQcReport                         Event = "bulk_po_create_qc_report"
	EventBulkPoUpdateMaterial                         Event = "bulk_po_update_raw_material"
	EventBulkPoMarkPps                                Event = "bulk_po_mark_pps"
	EventBulkPoUpdatePps                              Event = "bulk_po_update_pps"
	EventBulkPoMarkProduction                         Event = "bulk_po_mark_production"
	EventBulkPoMarkQc                                 Event = "bulk_po_mark_qc"
	EventBulkPoFirstPaymentWaitingConfirmBankTransfer Event = "bulk_po_first_payment_waiting_confirm_bank_transfer"
	EventBulkPoMarkFirstPayment                       Event = "bulk_po_mark_first_payment"
	EventBulkPoMarkFinalPayment                       Event = "bulk_po_mark_final_payment"
	EventBulkPoDelivered                              Event = "bulk_po_delivered"
	EventBulkPoMarkMaking                             Event = "bulk_po_mark_making"
	EventBulkPoMarkDelivering                         Event = "bulk_po_mark_delivering"
	EventBulkPoMarkRawMaterial                        Event = "bulk_po_mark_raw_material"
	EventBulkPoBuyerWaitingForQuotation               Event = "bulk_buyer_waiting_for_quotation"

	EventBulkPoConfirmQCReport                        Event = "bulk_po_final_confirm_qc_report"
	EventBulkPoFinalPaymentWaitingConfirmBankTransfer Event = "bulk_po_final_payment_waiting_confirm_bank_transfer"

	EventBulkPoBuyerSubmitOrder        Event = "bulk_po_buyer_submit_order"
	EventBulkPoBuyerApproveRawMaterial Event = "bulk_po_buyer_approve_raw_material"
	EventBulkPoBuyerApproveQc          Event = "bulk_po_buyer_approve_qc"
	EventBulkPoBuyerConfirmDelivered   Event = "bulk_po_buyer_confirm_delivered"

	EventBulkPoBuyerFirstPaymentSucceeded Event = "bulk_po_buyer_first_payment_succeeded"
	EventBulkPoBuyerFinalPaymentSucceeded Event = "bulk_po_buyer_final_payment_succeeded"
	EventBulkPoBuyerFinalPaymentDenied    Event = "bulk_po_buyer_final_payment_denied"
	EventBulkPoBuyerFirstPaymentDenied    Event = "bulk_po_buyer_first_payment_denied"

	EventBulkPoFirstPaymentSucceeded Event = "bulk_po_first_payment_succeeded"
	EventBulkPoFinalPaymentSucceeded Event = "bulk_po_final_payment_succeeded"

	EventBulkPoCreated Event = "bulk_po_created"

	EventBulkPoBuyerDepositSucceeded Event = "bulk_po_buyer_deposit_succeeded"
	EventBulkPoDepositSucceeded      Event = "bulk_po_deposit_succeeded"

	EventAdminInquiryNewNotes Event = "admin_inquiry_new_notes"
	EventAdminPONewNotes      Event = "admin_po_new_notes"
	EventAdminPoCanceled      Event = "admin_po_canceled"

	EventAdminBulkPONewNotes Event = "admin_bulk_po_new_notes"

	EventBuyerUpdateInquiry Event = "buyer_update_inquiry"

	EventAdminNewDesignComment  Event = "admin_new_design_comment"
	EventBuyerNewDesignComment  Event = "buyer_new_design_comment"
	EventSellerNewDesignComment Event = "seller_new_design_comment"

	EventBulkPoNewComment Event = "bulk_po_new_comment"

	EventSellerSubmitQuotation Event = "seller_submit_quotation"

	EventAdminCommentOnSellerRequest       Event = "admin_comment_on_seller_request"
	EventAdminDesignCommentOnSellerRequest Event = "admin_design_comment_on_seller_request"
	EventSellerCommentOnSellerRequest      Event = "seller_comment_on_seller_request"
	EventNewChatMessage                    Event = "new_chat_message"
	EventRemindUnseenChatMessage           Event = "remind_unseen_chat_message"

	EventSellerBulkPoUpdateRawMaterial Event = "seller_bulk_po_update_raw_material"

	EventSellerBulkPoUpdatePps Event = "seller_bulk_po_update_pps"

	EventSellerBulkPoUpdateProduction Event = "seller_bulk_po_update_production"

	EventSellerBulkPoMarkRawMaterial Event = "seller_bulk_po_mark_raw_material"
	EventSellerBulkPoMarkProduction  Event = "seller_bulk_po_mark_production"
	EventSellerBulkPoMarkInspection  Event = "seller_bulk_po_mark_inspection"

	EventBuyerCheckoutThroughBankTransfer Event = "buyer_checkout_through_bank_transfer"
	EventAdminConfirmPaymentReceived      Event = "admin_confirm_payment_received"
)
