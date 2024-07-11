package routes

import (
	"github.com/engineeringinflow/inflow-backend/pkg/models/enums"
	controllers "github.com/engineeringinflow/inflow-backend/services/backend/controllers/seller"
	"github.com/labstack/echo/v4"
)

func (router *Router) SetupSellerRoutes(g *echo.Group) {
	var authorizedGroup = g.Group("", router.Middlewares.IsAuthorized())
	var authorizedWithRoleGroup = authorizedGroup.Group("", router.Middlewares.CheckRole(enums.RoleSeller))
	var authorizedWithUserGroup = authorizedGroup.Group("", router.Middlewares.CheckTokenExpiredAndAttachUserInfo())

	authorizedWithUserGroup.POST("/me/track_activity", controllers.TrackActivity)

	authorizedWithRoleGroup.PUT("/me", controllers.UpdateMe)
	authorizedWithRoleGroup.GET("/me", controllers.GetMe)
	authorizedWithRoleGroup.GET("/me/teams", controllers.GetTeams)
	authorizedWithRoleGroup.DELETE("/me/logout", controllers.Logout)
	authorizedWithRoleGroup.PUT("/me/update_password", controllers.UpdatePassword)
	authorizedWithRoleGroup.POST("/me/onboarding_submit", controllers.OnboardingSubmit)

	authorizedWithRoleGroup.GET("/banks", controllers.PaginateUserBanks)
	authorizedWithRoleGroup.POST("/banks", controllers.CreateUserBanks)
	authorizedWithRoleGroup.PUT("/banks", controllers.UpdateUserBanks)
	authorizedWithRoleGroup.DELETE("/banks/:id", controllers.DeleteUserBanks)
	authorizedWithRoleGroup.DELETE("/banks/countries/:country_code", controllers.DeleteUserBanksByCountryCode)

	// Inquiry
	authorizedWithRoleGroup.GET("/inquiry_quotations", controllers.SellerInquiryList)
	authorizedWithRoleGroup.GET("/inquiry_quotations/:inquiry_seller_id", controllers.SellerInquiryDetail)
	authorizedWithRoleGroup.POST("/inquiry_quotations/:inquiry_seller_id/submit_quotation", controllers.SellerSubmitQuotation)
	authorizedWithRoleGroup.POST("/inquiry_quotations/submit_multiple_quotations", controllers.SellerSubmitMultipleQuotations)

	authorizedWithRoleGroup.POST("/inquiry_quotations/:inquiry_seller_id/approve_offer", controllers.SellerApproveOffer)
	authorizedWithRoleGroup.DELETE("/inquiry_quotations/:inquiry_seller_id/reject_offer", controllers.SellerRejectOffer)

	authorizedWithRoleGroup.POST("/inquiry_quotations/:inquiry_seller_id/design_comments", controllers.SellerInquiryCreateComment)
	authorizedWithRoleGroup.GET("/inquiry_quotations/:inquiry_seller_id/design_comments", controllers.SellerInquiryCommentList)
	authorizedWithRoleGroup.PUT("/inquiry_quotations/:inquiry_seller_id/design_comments/mark_seen", controllers.SellerInquiryCommentMarkSeen)
	authorizedWithRoleGroup.GET("/inquiry_quotations/:inquiry_seller_id/design_comments/status_count", controllers.SellerInquiryCommentStatusCount)

	// Purchase order
	authorizedWithRoleGroup.GET("/purchase_orders", controllers.SellerPaginatePurchaseOrders)
	authorizedWithRoleGroup.GET("/purchase_orders/:purchase_order_id", controllers.SellerGetPurchaseOrder)
	authorizedWithRoleGroup.POST("/purchase_orders/:purchase_order_id/approve_design", controllers.SellerApprovePurchaseOrderDesign)
	authorizedWithRoleGroup.POST("/purchase_orders/:purchase_order_id/approve_po", controllers.SellerPurchaseOrderApprovePo)
	authorizedWithRoleGroup.DELETE("/purchase_orders/:purchase_order_id/reject_po", controllers.SellerPurchaseOrderRejectPo)

	authorizedWithRoleGroup.GET("/purchase_orders/:purchase_order_id/po_upload_comments", controllers.SellerPurchaseOrderPaginatePoUploadComments)
	authorizedWithRoleGroup.POST("/purchase_orders/:purchase_order_id/po_upload_comments", controllers.SellerPurchaseOrderAddPoUploadComments)
	authorizedWithRoleGroup.PUT("/purchase_orders/:purchase_order_id/po_upload_comments/mark_seen", controllers.SellerPurchaseOrderPoUploadCommentMarkSeen)
	authorizedWithRoleGroup.GET("/purchase_orders/:purchase_order_id/po_upload_comments/status_count", controllers.SellerPurchasePoUploadCommentStatusCount)

	authorizedWithRoleGroup.GET("/purchase_orders/:purchase_order_id/design_comments", controllers.SellerPurchaseOrderPaginateDesignComments)
	authorizedWithRoleGroup.POST("/purchase_orders/:purchase_order_id/design_comments", controllers.SellerPurchaseOrderDesignCommentCreate)
	authorizedWithRoleGroup.PUT("/purchase_orders/:purchase_order_id/design_comments/mark_seen", controllers.SellerPurchaseOrderDesignCommentMarkSeen)
	authorizedWithRoleGroup.GET("/purchase_orders/:purchase_order_id/design_comments/status_count", controllers.SellerPurchasePoUploadCommentStatusCount)

	// After admin approve final design, seller will see and comment
	authorizedWithRoleGroup.GET("/purchase_orders/:purchase_order_id/final_design_comments", controllers.SellerPurchaseOrderPaginateFinalDesignComments)
	authorizedWithRoleGroup.POST("/purchase_orders/:purchase_order_id/final_design_comments", controllers.SellerPurchaseOrderFinalDesignCommentCreate)
	authorizedWithRoleGroup.PUT("/purchase_orders/:purchase_order_id/final_design_comments/mark_seen", controllers.SellerPurchaseOrderFinalDesignCommentMarkSeen)
	authorizedWithRoleGroup.GET("/purchase_orders/:purchase_order_id/final_design_comments/status_count", controllers.SellerPurchaseOrderFinalDesignCommentStatusCount)

	authorizedWithRoleGroup.PUT("/purchase_orders/:purchase_order_id/mark_raw_material", controllers.SellerPurchaseOrderMarkRawMaterial)
	authorizedWithRoleGroup.PUT("/purchase_orders/:purchase_order_id/update_raw_material", controllers.SellerPurchaseOrderUpdateRawMaterial)

	authorizedWithRoleGroup.GET("/purchase_orders/:purchase_order_id/raw_material_comments", controllers.SellerPurchaseOrderPaginateRawMaterialComments)
	authorizedWithRoleGroup.POST("/purchase_orders/:purchase_order_id/raw_material_comments", controllers.SellerPurchaseOrderRawMaterialCommentCreate)
	authorizedWithRoleGroup.PUT("/purchase_orders/:purchase_order_id/raw_material_comments/mark_seen", controllers.SellerPurchaseOrderRawMaterialCommentMarkSeen)
	authorizedWithRoleGroup.GET("/purchase_orders/:purchase_order_id/raw_material_comments/status_count", controllers.SellerPurchaseOrderRawMaterialStatusCount)

	authorizedWithRoleGroup.PUT("/purchase_orders/:purchase_order_id/mark_making", controllers.SellerPurchaseOrderMarkMaking)
	authorizedWithRoleGroup.PUT("/purchase_orders/:purchase_order_id/mark_submit", controllers.SellerPurchaseOrderMarkSubmit)
	authorizedWithRoleGroup.PUT("/purchase_orders/:purchase_order_id/mark_delivering", controllers.SellerPurchaseOrderMarkDelivering)
	authorizedWithRoleGroup.GET("/purchase_orders/:purchase_order_id/logs", controllers.SellerPaginatePurchaseOrderTracking)

	authorizedWithRoleGroup.PUT("/purchase_orders/:purchase_order_id/update_raw_material", controllers.PurchaseOrderUpdateRawMaterial)
	authorizedWithRoleGroup.PUT("/purchase_orders/:purchase_order_id/receive_payment", controllers.SellerPurchaseOrderReceivePayment)
	authorizedWithRoleGroup.PUT("/purchase_orders/:purchase_order_id/skip_raw_material", controllers.SellerPurchaseOrderSkipRawMaterial)

	authorizedWithRoleGroup.GET("/bulk_purchase_orders", controllers.SellerPaginateBulkPurchaseOrders)
	authorizedWithRoleGroup.GET("/bulk_purchase_orders/:bulk_purchase_order_id", controllers.SellerGetBulkPurchaseOrder)
	authorizedWithRoleGroup.PUT("/bulk_purchase_orders/:bulk_purchase_order_id/update_raw_material", controllers.SellerBulkPurchaseOrderUpdateRawMaterial)
	authorizedWithRoleGroup.PUT("/bulk_purchase_orders/:bulk_purchase_order_id/update_pps", controllers.SellerBulkPurchaseOrderUpdatePps)
	authorizedWithRoleGroup.PUT("/bulk_purchase_orders/:bulk_purchase_order_id/update_production", controllers.SellerBulkPurchaseOrderUpdateProduction)
	authorizedWithRoleGroup.POST("/bulk_purchase_orders/:bulk_purchase_order_id/qc_report", controllers.SellerBulkPurchaseOrderCreateQcReport)
	authorizedWithRoleGroup.PUT("/bulk_purchase_orders/:bulk_purchase_order_id/mark_delivering", controllers.SellerBulkPurchaseOrderMarkDelivering)
	authorizedWithRoleGroup.PUT("/bulk_purchase_orders/:bulk_purchase_order_id/confirm_receive_final_payment", controllers.SellerBulkPurchaseOrderConfirmReceiveFinalPayment)

	authorizedWithRoleGroup.PUT("/bulk_purchase_orders/:bulk_purchase_order_id/product_photo", controllers.SellerUpdateBulkPurchaseOrderProductPhoto)
	authorizedWithRoleGroup.PUT("/bulk_purchase_orders/:bulk_purchase_order_id/techpack", controllers.SellerUpdateBulkPurchaseOrderTechpack)
	authorizedWithRoleGroup.PUT("/bulk_purchase_orders/:bulk_purchase_order_id/bill_of_material", controllers.SellerUpdateBulkPurchaseOrderBillOfMaterial)
	authorizedWithRoleGroup.PUT("/bulk_purchase_orders/:bulk_purchase_order_id/size_chart", controllers.SellerUpdateBulkPurchaseOrderSizeChart)
	authorizedWithRoleGroup.PUT("/bulk_purchase_orders/:bulk_purchase_order_id/size_spec", controllers.SellerUpdateBulkPurchaseOrderSizeSpec)
	authorizedWithRoleGroup.PUT("/bulk_purchase_orders/:bulk_purchase_order_id/size_grading", controllers.SellerUpdateBulkPurchaseOrderSizeGrading)

	authorizedWithRoleGroup.PUT("/bulk_purchase_orders/:bulk_purchase_order_id/approve_po", controllers.SellerBulkPurchaseOrderApprovePO)
	authorizedWithRoleGroup.DELETE("/bulk_purchase_orders/:bulk_purchase_order_id/reject_po", controllers.SellerBulkPurchaseOrderRejectPO)
	authorizedWithRoleGroup.PUT("/bulk_purchase_orders/:bulk_purchase_order_id/start_without_first_payment", controllers.SellerBulkPurchaseOrderStartWithoutFirstPayment)
	authorizedWithRoleGroup.PUT("/bulk_purchase_orders/:bulk_purchase_order_id/confirm_receive_first_payment", controllers.SellerBulkPurchaseOrderConfirmReceiveFirstPayment)
	authorizedWithRoleGroup.GET("/bulk_purchase_orders/:bulk_purchase_order_id/logs", controllers.SellerPaginateBulkPurchaseOrderTracking)

	authorizedWithRoleGroup.PUT("/bulk_purchase_orders/:bulk_purchase_order_id/feedback", controllers.SellerBulkPurchaseOrderFeedback)
	authorizedWithRoleGroup.GET("/bulk_purchase_orders/:bulk_purchase_order_id/first_payment_invoice", controllers.SellerBulkPurchaseOrderFirstPaymentInvoice)
	authorizedWithRoleGroup.GET("/bulk_purchase_orders/:bulk_purchase_order_id/final_payment_invoice", controllers.SellerBulkPurchaseOrderFinalPaymentInvoice)
	authorizedWithRoleGroup.PUT("/bulk_purchase_orders/:bulk_purchase_order_id/mark_raw_material", controllers.SellerBulkPurchaseOrderMarkRawMaterial)
	authorizedWithRoleGroup.PUT("/bulk_purchase_orders/:bulk_purchase_order_id/mark_production", controllers.SellerBulkPurchaseOrderMarkProduction)
	authorizedWithRoleGroup.PUT("/bulk_purchase_orders/:bulk_purchase_order_id/mark_inspection", controllers.SellerBulkPurchaseOrderMarkInspection)

	// Bulk PO
	authorizedWithRoleGroup.POST("/bulk_quotations/:seller_quotation_id/submit_quotation", controllers.SellerBulkPOSubmitQuotation)
	authorizedWithRoleGroup.POST("/bulk_quotations/:seller_quotation_id/re_submit_quotation", controllers.SellerBulkPOReSubmitQuotation)
	authorizedWithRoleGroup.POST("/bulk_quotations/submit_multiple_quotations", controllers.SellerBulkPOSubmitMultipleQuotation)
	authorizedWithRoleGroup.GET("/bulk_quotations/:seller_quotation_id", controllers.SellerGetBulkPOQuotationDetails)

	// Payment Transactions
	authorizedWithRoleGroup.GET("/payment_transactions", controllers.PaginatePaymentTransaction)
	authorizedWithRoleGroup.GET("/payment_transactions/:payment_transaction_id", controllers.GetPaymentTransaction)

	// Dashboard
	authorizedWithRoleGroup.GET("/dashboard/revenue", controllers.SellerDashboardRevenue)
	authorizedWithRoleGroup.GET("/dashboard", controllers.SellerDashboard)
	// Chat
	authorizedWithRoleGroup.POST("/chat_messages", controllers.SellerCreateChatMessage)
	authorizedWithRoleGroup.GET("/chat_messages", controllers.SellerGetChatMessageList)

	// Chat Room
	authorizedWithRoleGroup.GET("/chat_rooms/relevant_stage", controllers.SellerGetChatUserRelevantStage)
	authorizedWithRoleGroup.POST("/chat_rooms", controllers.SellerCreateChatRoom)
	authorizedWithRoleGroup.GET("/chat_rooms", controllers.SellerGetChatRoomList)
	authorizedWithRoleGroup.PUT("/chat_rooms/:chat_room_id/seen_messages", controllers.SellerMarkSeenChatRoomMessage)
	authorizedWithRoleGroup.GET("/chat_rooms/unseen_message", controllers.SellerCountUnseenChatMessageOnRoom)

}
