package routes

import (
	"github.com/engineeringinflow/inflow-backend/pkg/models/enums"
	controllers "github.com/engineeringinflow/inflow-backend/services/backend/controllers/buyer"
	"github.com/labstack/echo/v4"
)

func (router *Router) SetupBuyerRoutes(g *echo.Group) {
	var authorizedGroup = g.Group("", router.Middlewares.IsAuthorized())
	var authorizedWithRoleGroup = authorizedGroup.Group("", router.Middlewares.CheckRole(enums.RoleClient))
	var authorizedWithUserGroup = authorizedGroup.Group("", router.Middlewares.CheckTokenExpiredAndAttachUserInfo())

	authorizedWithUserGroup.POST("/me/track_activity", controllers.TrackActivity)

	authorizedWithRoleGroup.GET("/products/get_category_tree", controllers.GetCategoryTree)
	authorizedWithRoleGroup.GET("/collections", controllers.PaginateCollection)

	authorizedWithRoleGroup.PUT("/me", controllers.UpdateMe)
	authorizedWithRoleGroup.GET("/me", controllers.GetMe)
	authorizedWithRoleGroup.GET("/me/teams", controllers.GetTeams)
	authorizedWithRoleGroup.DELETE("/me/logout", controllers.Logout)
	authorizedWithRoleGroup.PUT("/me/update_password", controllers.UpdatePassword)
	authorizedWithRoleGroup.GET("/me/last_shipping_address", controllers.GetLastShippingAddress)
	authorizedWithRoleGroup.PUT("/me/complete_inquiry_tutorial", controllers.CompleteInquiryTutorial)

	authorizedWithRoleGroup.GET("/documents", controllers.BuyerGetDocumentList)

	authorizedWithRoleGroup.GET("/products", controllers.PaginateProducts)
	authorizedWithRoleGroup.GET("/products/recommend", controllers.ProductRecommend)
	authorizedWithRoleGroup.GET("/products/get_category_tree", controllers.GetCategoryTree)
	authorizedWithRoleGroup.GET("/products/categories/:parent_category_id/children", controllers.PaginateProductCategoriesChildren)
	authorizedWithRoleGroup.GET("/products/categories", controllers.GetProductCategories)
	authorizedWithRoleGroup.POST("/products/:product_id/cart", controllers.CreateCatalogCart)

	authorizedWithRoleGroup.POST("/payment_methods/create", controllers.AddPaymentMethod)
	authorizedWithRoleGroup.DELETE("/payment_methods/:payment_method_id/detact", controllers.RemovePaymentMethod)
	authorizedWithRoleGroup.GET("/payment_methods/:payment_method_id", controllers.GetPaymentMethod)
	authorizedWithRoleGroup.PUT("/payment_methods/:payment_method_id/default", controllers.MarkDefaultPaymentMethod)
	authorizedWithRoleGroup.GET("/payment_methods", controllers.GetPaymentMethods)

	// Collection product
	authorizedGroup.GET("/collections", controllers.PaginateCollection)
	authorizedGroup.GET("/collections/:collection_id/get", controllers.CollectionDetail)
	authorizedGroup.GET("/collections/get_product", controllers.CollectionGetProduct)
	authorizedGroup.GET("/collections/:collection_id/get_product", controllers.CollectionGetProduct)

	authorizedWithRoleGroup.POST("/setup_payment", controllers.SetupPayment)

	// Inquiry
	authorizedWithRoleGroup.GET("/inquiries", controllers.BuyerPaginateInquiry)
	authorizedWithRoleGroup.POST("/inquiries/create", controllers.BuyerCreateInquiry)
	authorizedWithRoleGroup.POST("/inquiries/create_multiple", controllers.BuyerCreateMultipleInquiries)
	authorizedWithRoleGroup.GET("/inquiries/:inquiry_id", controllers.BuyerGetInquiry)
	authorizedWithRoleGroup.PUT("/inquiries/:inquiry_id", controllers.BuyerUpdateInquiry)
	authorizedWithRoleGroup.PUT("/inquiries/:inquiry_id/edit_timeout", controllers.BuyerExtendInquiryEditTimeout)
	authorizedWithRoleGroup.POST("/inquiries/:inquiry_id/clone", controllers.BuyerCloneInquiry)

	authorizedWithRoleGroup.POST("/inquiries/:inquiry_id/approve_quotation", controllers.BuyerApproveInquiryQuotation)
	authorizedWithRoleGroup.POST("/inquiries/:inquiry_id/reject_quotation", controllers.BuyerRejectInquiryQuotation)
	authorizedWithRoleGroup.GET("/inquiries/:inquiry_id/quotation_history", controllers.BuyerInquiryQuotationHistory)
	authorizedWithRoleGroup.GET("/inquiries/:inquiry_id/logs", controllers.BuyerInquiryLogs)
	authorizedWithRoleGroup.POST("/inquiries/:inquiry_id/confirm", controllers.BuyerConfirmInquiry)
	authorizedWithRoleGroup.POST("/inquiries/approve_multiple_quotations", controllers.BuyerApproveMultipleInquiryQuotations)
	authorizedWithRoleGroup.POST("/inquiries/reject_multiple_quotations", controllers.BuyerRejectMultipleInquiryQuotations)

	authorizedWithRoleGroup.GET("/inquiries/collections", controllers.BuyerInquiryCollections)
	authorizedWithRoleGroup.POST("/inquiries/collections", controllers.BuyerInquiryCollectionCreate)

	// Cart with multi inquiry
	authorizedWithRoleGroup.GET("/inquiry_carts/cart", controllers.BuyerInquiryCart)
	authorizedWithRoleGroup.POST("/inquiries/:inquiry_id/update_cart_items", controllers.BuyerInquiryUpdateCartItems)
	authorizedWithRoleGroup.GET("/inquiries/:inquiry_id/cart_items", controllers.BuyerInquiryCartItems)
	authorizedWithRoleGroup.DELETE("/inquiries/:inquiry_id/remove_items", controllers.BuyerInquiryRemoveItems)
	authorizedWithRoleGroup.POST("/inquiry_carts/preview_checkout", controllers.BuyerMultiInquiryPreviewCheckout)
	authorizedWithRoleGroup.POST("/inquiry_carts/checkout", controllers.BuyerMultiInquiryCheckout)
	authorizedWithRoleGroup.GET("/inquiry_carts/checkout_info", controllers.BuyerMultiInquiryCheckoutInfo)
	authorizedWithRoleGroup.PUT("/inquiries/:inquiry_id/attachments", controllers.BuyerInquiryUpdateAttachments)
	authorizedWithRoleGroup.PUT("/inquiries/:inquiry_id/close", controllers.BuyerCloseInquiry)
	authorizedWithRoleGroup.DELETE("/inquiries/:inquiry_id/cancel", controllers.BuyerCancelInquiry)
	authorizedWithRoleGroup.POST("/inquiries/:inquiry_id/preview_checkout", controllers.BuyerInquiryPreviewCheckout)
	authorizedWithRoleGroup.POST("/inquiries/:inquiry_id/checkout", controllers.BuyerInquiryCheckout)

	authorizedWithRoleGroup.PUT("/inquiries/:inquiry_id/logs", controllers.BuyerUpdateInquiryLogs)
	authorizedWithRoleGroup.DELETE("/inquiries/:inquiry_id/logs", controllers.BuyerDeleteInquiryLogs)

	authorizedWithRoleGroup.GET("/catalog_carts", controllers.BuyerPaginateCatalogCarts)
	authorizedWithRoleGroup.PUT("/catalog_carts", controllers.BuyerUpdateCatalogCarts)
	authorizedWithRoleGroup.POST("/catalog_carts/place_orders", controllers.BuyerCreateCatalogCartsOrders)
	authorizedWithRoleGroup.POST("/catalog_carts/checkout_info", controllers.BuyerMultiCatalogCartCheckoutInfo)
	authorizedWithRoleGroup.POST("/catalog_carts/checkout", controllers.BuyerMultiCatalogCartCheckout)

	authorizedWithRoleGroup.GET("/purchase_orders", controllers.PaginatePurchaseOrders)
	authorizedWithRoleGroup.PUT("/purchase_orders/:purchase_order_id", controllers.UpdatePurchaseOrder)
	authorizedWithRoleGroup.GET("/purchase_orders/:purchase_order_id", controllers.GetPurchaseOrder)
	authorizedWithRoleGroup.GET("/purchase_orders/:purchase_order_id/invoice", controllers.GetPurchaseOrderInvoice)
	authorizedWithRoleGroup.PUT("/purchase_orders/:purchase_order_id/feedback", controllers.PurchaseOrderFeedback)
	authorizedWithRoleGroup.PUT("/purchase_orders/:purchase_order_id/update_tracking_status", controllers.UpdatePurchaseOrderTrackingStatus)
	authorizedWithRoleGroup.GET("/purchase_orders/:purchase_order_id/logs", controllers.PaginatePurchaseOrderTracking)
	authorizedWithRoleGroup.POST("/purchase_orders/:purchase_order_id/approve_design", controllers.BuyerApprovePurchaseOrderDesign)
	authorizedWithRoleGroup.POST("/purchase_orders/:purchase_order_id/reject_design", controllers.BuyerRejectPurchaseOrderDesign)
	authorizedWithRoleGroup.POST("/purchase_orders/:purchase_order_id/confirm_delivered", controllers.BuyerPurchaseOrderConfirmDelivered)
	authorizedWithRoleGroup.PUT("/purchase_orders/:purchase_order_id/update_design", controllers.BuyerPurchaseOrderUpdateDesign)

	authorizedWithRoleGroup.POST("/purchase_orders/:purchase_order_id/design_comments", controllers.BuyerPurchaseOrderAddDesignComments)
	authorizedWithRoleGroup.GET("/purchase_orders/:purchase_order_id/design_comments", controllers.BuyerPurchaseOrderPaginateDesignComments)
	authorizedWithRoleGroup.PUT("/purchase_orders/:purchase_order_id/design_comments/mark_seen", controllers.BuyerPurchaseOrderDesignCommentMarkSeen)
	authorizedWithRoleGroup.GET("/purchase_orders/:purchase_order_id/design_comments/status_count", controllers.BuyerPurchaseOrderDesignCommentStatusCount)

	authorizedWithRoleGroup.POST("/purchase_orders/:purchase_order_id/approve_raw_material", controllers.PurchaseOrderBuyerApproveRawMaterial)
	authorizedWithRoleGroup.POST("/purchase_orders/:purchase_order_id/bulk_purchase_order", controllers.BuyerCreateBulkPurchaseOrderFromSample)

	authorizedWithRoleGroup.PUT("/purchase_orders/:purchase_order_id/logs", controllers.BuyerUpdatePurchaseOrderLogs)
	authorizedWithRoleGroup.DELETE("/purchase_orders/:purchase_order_id/logs", controllers.BuyerDeletePurchaseOrderLogs)
	authorizedWithRoleGroup.POST("/purchase_orders/:purchase_order_id/make_another", controllers.BuyerMakeAnotherPurchaseOrder)

	authorizedWithRoleGroup.PUT("/bulk_purchase_orders/:bulk_purchase_order_id", controllers.UpdateBulkPurchaseOrder)
	authorizedWithRoleGroup.GET("/bulk_purchase_orders/:bulk_purchase_order_id", controllers.GetBulkPurchaseOrder)
	authorizedWithRoleGroup.POST("/bulk_purchase_orders/:bulk_purchase_order_id/submit", controllers.SubmitBulkPurchaseOrder)
	authorizedWithRoleGroup.GET("/bulk_purchase_orders", controllers.PaginateBulkPurchaseOrder)
	authorizedWithRoleGroup.POST("/bulk_purchase_orders/preview_checkout", controllers.BulkPurchaseOrdersPreviewCheckout)
	authorizedWithRoleGroup.POST("/bulk_purchase_orders/:bulk_purchase_order_id/preview_checkout", controllers.BulkPurchaseOrderPreviewCheckout) //deprecated
	authorizedWithRoleGroup.POST("/bulk_purchase_orders/:bulk_purchase_order_id/checkout", controllers.BulkPurchaseOrderCheckout)                //deprecated
	authorizedWithRoleGroup.GET("/bulk_purchase_orders/:bulk_purchase_order_id/logs", controllers.PaginateBulkPurchaseOrderTracking)
	authorizedWithRoleGroup.POST("/bulk_purchase_orders/:bulk_purchase_order_id/approve_qc", controllers.BulkPurchaseBuyerApproveQc)
	authorizedWithRoleGroup.POST("/bulk_purchase_orders/:bulk_purchase_order_id/approve_raw_material", controllers.BulkPurchaseBuyerApproveRawMaterial)
	authorizedWithRoleGroup.POST("/bulk_purchase_orders/:bulk_purchase_order_id/confirm_delivered", controllers.BuyerBulkPurchaseOrderConfirmDelivered)
	authorizedWithRoleGroup.GET("/bulk_purchase_orders/:bulk_purchase_order_id/sample_po", controllers.BulkPurchaseOrderGetSamplePO)
	authorizedWithRoleGroup.PUT("/bulk_purchase_orders/:bulk_purchase_order_id/update_design", controllers.BulkPurchaseOrderUpdateDesign)
	authorizedWithRoleGroup.PUT("/bulk_purchase_orders/:bulk_purchase_order_id/feedback", controllers.BulkPurchaseOrderFeedback)

	authorizedWithRoleGroup.PUT("/bulk_purchase_orders/:bulk_purchase_order_id/logs", controllers.BuyerUpdateBulkPurchaseOrderLogs)
	authorizedWithRoleGroup.DELETE("/bulk_purchase_orders/:bulk_purchase_order_id/logs", controllers.BuyerDeleteBulkPurchaseOrderLogs)

	authorizedWithRoleGroup.GET("/bulk_purchase_orders/:bulk_purchase_order_id/first_payment_invoice", controllers.GetBulkPurchaseOrderFirstPaymentInvoice)
	authorizedWithRoleGroup.GET("/bulk_purchase_orders/:bulk_purchase_order_id/final_payment_invoice", controllers.GetBulkPurchaseOrderFinalPaymentInvoice)
	authorizedWithRoleGroup.GET("/bulk_purchase_orders/:bulk_purchase_order_id/debit_notes", controllers.GetBulkPurchaseOrderDebitNotes)
	authorizedWithRoleGroup.GET("/bulk_purchase_orders/:bulk_purchase_order_id/invoice", controllers.GetBulkPurchaseOrderInvoice)
	authorizedWithRoleGroup.PUT("/bulk_purchase_orders/:bulk_purchase_order_id/update_pps", controllers.BulkPurchaseOrderUpdatePps)

	// Notification
	authorizedWithRoleGroup.GET("/notifications", controllers.PaginateUserNotifications)
	authorizedWithRoleGroup.PUT("/notifications/:notification_id/mark_seen", controllers.UserNotificationMarkSeen)
	authorizedWithRoleGroup.PUT("/notifications/mark_seen_all", controllers.NotificationMarkSeenAll)

	// Sys Notification
	authorizedWithRoleGroup.GET("/sys/notifications", controllers.PaginateSysNotifications)
	authorizedWithRoleGroup.PUT("/sys/notifications/:notification_id/mark_seen", controllers.SysNotificationMarkSeen)
	authorizedWithRoleGroup.PUT("/sys/notifications/mark_seen_all", controllers.SysNotificationMarkSeenAll)

	// Team Invite
	authorizedWithRoleGroup.GET("/teams/members", controllers.BrandTeamMembers)
	authorizedWithRoleGroup.POST("/teams/members/invite", controllers.BrandTeamInvite)
	authorizedWithRoleGroup.PUT("/teams/members/:member_id/actions", controllers.UpdateBrandTeamMemberActions)
	authorizedWithRoleGroup.DELETE("/teams/members/:member_id", controllers.DeleteBrandTeamMember)

	// Ads videos
	authorizedWithRoleGroup.GET("/ads_videos", controllers.PaginateAdsVideos)

	// Docs Agreement
	authorizedWithRoleGroup.GET("/docs/:type", controllers.GetSettingDoc)
	authorizedWithRoleGroup.GET("/docs/:setting_doc_type/agreement", controllers.GetDocsAgreement)
	authorizedWithRoleGroup.POST("/docs/:setting_doc_type/agreement", controllers.CreateDocsAgreement)

	// Dashboard
	authorizedWithRoleGroup.GET("/data_analytics/rfq", controllers.GetDataAnalyticRFQ)
	authorizedWithRoleGroup.GET("/data_analytics/pending_tasks", controllers.GetDataAnalyticPendingTasks)
	authorizedWithRoleGroup.GET("/data_analytics/pending_payments", controllers.GetDataAnalyticPendingPayments)
	authorizedWithRoleGroup.GET("/data_analytics/total_styles_produced", controllers.GetDataAnalyticTotalStyleProduced)

	// Time And Action
	authorizedWithRoleGroup.GET("/tnas", controllers.PaginateTNA)

	// Chat
	authorizedWithRoleGroup.POST("/chat_messages", controllers.BuyerCreateChatMessage)
	authorizedWithRoleGroup.GET("/chat_messages", controllers.BuyerGetChatMessageList)

	// Chat Room
	authorizedWithRoleGroup.GET("/chat_rooms/relevant_stage", controllers.BuyerGetChatUserRelevantStage)
	authorizedWithRoleGroup.POST("/chat_rooms", controllers.BuyerCreateChatRoom)
	authorizedWithRoleGroup.GET("/chat_rooms", controllers.BuyerGetChatRoomList)
	authorizedWithRoleGroup.PUT("/chat_rooms/:chat_room_id/seen_messages", controllers.BuyerMarkSeenChatRoomMessage)
	authorizedWithRoleGroup.GET("/chat_rooms/unseen_message", controllers.BuyerCountUnseenChatMessageOnRoom)
	// Order Group
	authorizedWithRoleGroup.GET("/order_groups", controllers.BuyerGetOrderGroupList)
	authorizedWithRoleGroup.POST("/order_groups", controllers.BuyerCreateOrderGroup)
	authorizedWithRoleGroup.GET("/order_groups/:order_group_id", controllers.BuyerGetOrderGroupDetail)
	authorizedWithRoleGroup.POST("/order_groups/assign", controllers.BuyerAssignOrderGroup)
	// Order Cart
	authorizedWithRoleGroup.GET("/order_cart", controllers.BuyerGetOrderCart)
	authorizedWithRoleGroup.POST("/order_cart/preview_checkout", controllers.BuyerGetOrderCartPreviewCheckout)
	authorizedWithRoleGroup.POST("/order_cart/checkout", controllers.BuyerOrderCartCheckout)
	authorizedWithRoleGroup.POST("/order_cart/checkout_info", controllers.BuyerOrderCartGetCheckoutInfo)

	// Analytics
	authorizedWithRoleGroup.GET("/analytics/products", controllers.PaginateAnalyticProduct)
	authorizedWithRoleGroup.GET("/analytics/products/group", controllers.PaginateAnalyticProductGroup)
	authorizedWithRoleGroup.GET("/analytics/products/:product_id", controllers.GetAnalyticProductDetails)
	authorizedWithRoleGroup.GET("/analytics/products/recommend", controllers.RecommendAnalyticProducts)
	authorizedWithRoleGroup.GET("/analytics/products/:product_id/chart", controllers.GetAnalyticProductChart)
	authorizedWithRoleGroup.GET("/analytics/products/product_classes/group", controllers.GetAnalyticProductClassGroup)
	authorizedWithRoleGroup.GET("/analytics/products/trending/group", controllers.GetAnalyticProductTrendingGroup)

	authorizedWithRoleGroup.GET("/trendings", controllers.PaginateTrendings)
	authorizedWithRoleGroup.GET("/trendings/:id", controllers.GetTrending)
	authorizedWithRoleGroup.GET("/analytics/product_trendings/:product_trending_id", controllers.GetProductTrending)
	authorizedWithRoleGroup.GET("/analytics/product_trendings/:product_trending_id/chart", controllers.GetProductTrendingChart)
	authorizedWithRoleGroup.GET("/analytics/product_trendings/tags", controllers.GetProductTrendingTags)

}
