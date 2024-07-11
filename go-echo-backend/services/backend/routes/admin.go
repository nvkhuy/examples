package routes

import (
	"github.com/engineeringinflow/inflow-backend/pkg/models/enums"
	controllers "github.com/engineeringinflow/inflow-backend/services/backend/controllers/admin"
	"github.com/labstack/echo/v4"
)

func (router *Router) SetupAdminRoutes(g *echo.Group) {
	var authorizedGroup = g.Group("", router.Middlewares.IsAuthorized())
	var authorizedWithRoleGroup = authorizedGroup.Group("", router.Middlewares.CheckRole(enums.RoleSuperAdmin, enums.RoleLeader, enums.RoleStaff))
	var authorizedWithUserGroup = authorizedWithRoleGroup.Group("", router.Middlewares.CheckTokenExpiredAndAttachUserInfo())

	authorizedWithUserGroup.POST("/me/track_activity", controllers.TrackActivity)

	authorizedWithRoleGroup.DELETE("/me/logout", controllers.Logout)

	authorizedWithRoleGroup.POST("/task/dispatch", controllers.DispatchTask)

	// Me
	authorizedWithRoleGroup.PUT("/me", controllers.UpdateMe)
	authorizedWithRoleGroup.GET("/me", controllers.GetMe)
	authorizedWithRoleGroup.PUT("/me/change_password", controllers.ChangePassword)

	// authorizedWithRoleGroup.GET("/inventories", controllers.AdminPaginateInventories)
	// authorizedWithRoleGroup.POST("/inventories/restock", controllers.AdminInventoryRestock)

	// User
	authorizedWithRoleGroup.GET("/users/search", controllers.SearchUsers)
	authorizedWithRoleGroup.GET("/users", controllers.PaginateUsers)
	authorizedWithRoleGroup.GET("/users/recent", controllers.PaginateUserRecent)
	authorizedWithRoleGroup.GET("/users/:user_id", controllers.GetUser)
	authorizedWithRoleGroup.PUT("/users/:user_id/assign_owners", controllers.AssignContactOwners)
	authorizedWithRoleGroup.POST("/users/:user_id/access_token", controllers.GetAccessToken)
	authorizedWithRoleGroup.PUT("/users/:user_id/change_password", controllers.ChangePassword)
	authorizedWithRoleGroup.DELETE("/users/:user_id/archive", controllers.ArchiveUser)
	authorizedWithRoleGroup.PUT("/users/:user_id/unarchive", controllers.UnarchiveUser)
	authorizedWithRoleGroup.PUT("/users/:user_id", controllers.UpdateUser)
	authorizedWithRoleGroup.POST("/users/:user_id/approve", controllers.ApproveUser)
	authorizedWithRoleGroup.POST("/users/invite", controllers.InviteUser)
	authorizedWithRoleGroup.POST("/users/client", controllers.CreateClient)
	authorizedWithRoleGroup.DELETE("/users/:user_id/reject", controllers.RejectUser)
	authorizedWithRoleGroup.DELETE("/users/:user_id/delete", controllers.DeleteUser)
	authorizedWithRoleGroup.GET("/users/:user_id/payment_methods", controllers.GetUserPaymentMethods)
	authorizedWithRoleGroup.GET("/users/:user_id/banks", controllers.GetUserBanks)
	authorizedWithRoleGroup.POST("/users/:user_id/upload_bulks", controllers.UploadBulkPurchaseOrder)
	authorizedWithRoleGroup.GET("/users/:user_id/activities", controllers.GetActivities)

	// Category
	authorizedWithRoleGroup.PUT("/categories/:category_id", controllers.UpdateCategory)
	authorizedWithRoleGroup.PATCH("/categories/slug/generate", controllers.GenerateCategorySlug)
	authorizedWithRoleGroup.POST("/categories", controllers.CreateCategory)
	authorizedWithRoleGroup.DELETE("/categories/:category_id/delete", controllers.DeleteCategory)
	authorizedWithRoleGroup.GET("/categories/get_category_tree", controllers.AdminGetCategoryTree)

	// Product
	authorizedWithRoleGroup.GET("/products", controllers.PaginateProduct)
	authorizedWithRoleGroup.GET("/products/export", controllers.ExportProducts)
	authorizedWithRoleGroup.GET("/products/search", controllers.SearchProduct)
	authorizedWithRoleGroup.PATCH("/products/slug/generate", controllers.GenerateProductSlug)
	authorizedWithRoleGroup.GET("/products/types/price", controllers.PaginateProductTypesPrice)
	authorizedWithRoleGroup.PATCH("/products/types/price", controllers.FetchProductTypesPrice)
	authorizedWithRoleGroup.PATCH("/products/fabric/price", controllers.FetchFabricTypesPrice)
	authorizedWithRoleGroup.PATCH("/products/types/price/images", controllers.FetchProductTypesPriceImagesURL)
	authorizedWithRoleGroup.GET("/products/types/price/vine", controllers.PaginateProductTypesPriceVine)
	authorizedWithRoleGroup.GET("/products/fabric/price/vine", controllers.PaginateRWDFabricPriceVine)
	authorizedWithRoleGroup.GET("/products/types/price/quote", controllers.PaginateProductTypesPriceQuote)
	authorizedWithRoleGroup.POST("/products/create", controllers.CreateProduct)
	authorizedWithRoleGroup.PUT("/products/:product_id", controllers.UpdateProduct)
	authorizedWithRoleGroup.GET("/products/:product_id/get", controllers.GetProduct)
	authorizedWithRoleGroup.DELETE("/products/:product_id/delete", controllers.AdminDeleteProduct)
	authorizedWithRoleGroup.GET("/products/qr_code", controllers.GetProductQRCode)

	// Subscriber
	authorizedWithRoleGroup.GET("/subscribers/search", controllers.SearchSubscribers)

	// Collection product
	authorizedWithRoleGroup.POST("/collections/create", controllers.AdminCreateCollection)
	authorizedWithRoleGroup.PUT("/collections/:collection_id", controllers.AdminUpdateCollection)
	authorizedWithRoleGroup.DELETE("/collections/:collection_id/delete", controllers.AdminDeleteCollection)
	authorizedWithRoleGroup.GET("/collections", controllers.AdminCollectionList)
	authorizedWithRoleGroup.GET("/collections/:collection_id/get", controllers.AdminCollectionDetail)
	authorizedWithRoleGroup.GET("/collections/get_product", controllers.AdminCollectionList)
	authorizedWithRoleGroup.POST("/collections/:collection_id/add_product", controllers.AdminCollectionAddProduct)
	authorizedWithRoleGroup.DELETE("/collections/:collection_id/delete_product", controllers.AdminCollectionDeleteProduct)
	authorizedWithRoleGroup.GET("/collections/:collection_id/get_product", controllers.AdminCollectionGetProduct)

	// Page
	authorizedWithRoleGroup.GET("/pages", controllers.AdminPageList)
	authorizedWithRoleGroup.GET("/pages/:page_id", controllers.AdminPageDetail)
	authorizedWithRoleGroup.PUT("/pages/:page_id", controllers.AdminUpdatePage)

	authorizedWithRoleGroup.POST("/pages/create", controllers.AdminCreatePage)
	authorizedWithRoleGroup.POST("/pages/:id/add_section", controllers.AdminAddPageSection)
	authorizedWithRoleGroup.PUT("/pages/sections/:id/update", controllers.AdminUpdatePageSection)
	authorizedWithRoleGroup.DELETE("/pages/sections/:id/delete", controllers.AdminDeletePageSection)

	authorizedWithRoleGroup.GET("/stats/products", controllers.StatsProducts)
	authorizedWithRoleGroup.GET("/stats/buyers", controllers.StatBuyers)
	authorizedWithRoleGroup.GET("/stats/suppliers", controllers.StatsSuppliers)

	// News on homepage
	authorizedWithRoleGroup.POST("/posts", controllers.CreatePost)
	authorizedWithRoleGroup.PUT("/posts/:post_id", controllers.UpdatePost)
	authorizedWithRoleGroup.GET("/posts", controllers.PaginatePost)
	authorizedWithRoleGroup.GET("/posts/:slug", controllers.GetPost)
	authorizedWithRoleGroup.DELETE("/posts/:post_id/delete", controllers.DeletePost)
	authorizedWithRoleGroup.PATCH("/posts/slug/generate", controllers.GeneratePostSlug)

	// Blog category
	authorizedWithRoleGroup.GET("/blog/categories", controllers.PaginateBlogCategory)
	authorizedWithRoleGroup.POST("/blog/categories", controllers.CreateBlogCategory)
	authorizedWithRoleGroup.PUT("/blog/categories/:blog_category_id", controllers.UpdateBlogCategory)
	authorizedWithRoleGroup.DELETE("/blog/categories/:blog_category_id", controllers.DeleteBlogCategory)

	// Inquiries
	authorizedWithRoleGroup.GET("/inquiries", controllers.AdminPaginateInquiry)
	authorizedWithRoleGroup.GET("/inquiries/export", controllers.ExportInquiries)
	authorizedWithRoleGroup.POST("/inquiries", controllers.AdminCreateInquiry)
	authorizedWithRoleGroup.GET("/inquiries/for_creating_order", controllers.AdminInquiryListForCreatingOrder)
	authorizedWithRoleGroup.GET("/inquiries/:inquiry_id", controllers.AdminInquiryDetail)
	authorizedWithRoleGroup.PUT("/inquiries/:inquiry_id", controllers.AdminInquiryUpdate)
	authorizedWithRoleGroup.DELETE("/inquiries/:inquiry_id/archive", controllers.AdminArchiveInquiry)
	authorizedWithRoleGroup.PUT("/inquiries/:inquiry_id/unarchive", controllers.AdminUnarchiveInquiry)
	authorizedWithRoleGroup.DELETE("/inquiries/:inquiry_id/delete", controllers.AdminDeleteInquiry)
	authorizedWithRoleGroup.PUT("/inquiries/:inquiry_id/close", controllers.AdminCloseInquiry)
	authorizedWithRoleGroup.POST("/inquiries/:inquiry_id/notes", controllers.InquiryAddNote)
	authorizedWithRoleGroup.GET("/inquiries/:inquiry_id/notes", controllers.PaginateInquiryNotes)
	authorizedWithRoleGroup.PUT("/inquiries/:inquiry_id/notes/mark_seen", controllers.InquiryNoteMarkSeen)
	authorizedWithRoleGroup.GET("/inquiries/:inquiry_id/notes/unread_count", controllers.InquiryNoteUnreadCount)
	authorizedWithRoleGroup.DELETE("/inquiries/:inquiry_id/notes/:comment_id", controllers.AdminInquiryCommentDelete)
	authorizedWithRoleGroup.POST("/inquiries/:inquiry_id/preview_checkout", controllers.AdminInquiryPreviewCheckout)

	authorizedWithRoleGroup.POST("/inquiries/:inquiry_id/send_to_buyer", controllers.AdminSendInquiryToBuyer)
	authorizedWithRoleGroup.POST("/inquiries/:inquiry_id/submit_quotation", controllers.AdminSubmitInquiryQuotation)
	authorizedWithRoleGroup.POST("/inquiries/submit_multiple_quotations", controllers.AdminSubmitMultipleInquiryQuotations)
	authorizedWithRoleGroup.POST("/inquiries/:inquiry_id/internal_approve_quotation", controllers.AdminInquiryInternalApproveQuotation)
	authorizedWithRoleGroup.POST("/inquiries/approve_multiple_quotations", controllers.AdminApproveMultipleInquiryQuotations)

	authorizedWithRoleGroup.GET("/inquiries/:inquiry_id/quotation_history", controllers.AdminInquiryQuotationHistory)
	authorizedWithRoleGroup.GET("/inquiries/:inquiry_id/logs", controllers.AdminInquiryLogs)

	authorizedWithRoleGroup.PUT("/inquiries/:inquiry_id/mark_seen", controllers.AdminInquiryMarkSeen)
	authorizedWithRoleGroup.POST("/inquiries/:inquiry_id/clone", controllers.AdminCloneInquiry)

	authorizedWithRoleGroup.POST("/inquiries/:inquiry_id/mark_as_paid", controllers.AdminInquiryMarkAsPaid)
	authorizedWithRoleGroup.POST("/inquiries/:inquiry_id/mark_as_unpaid", controllers.AdminInquiryMarkAsUnpaid)

	authorizedWithRoleGroup.PUT("/inquiries/:inquiry_id/assign_pic", controllers.AdminInquiryAssignPIC)
	authorizedWithRoleGroup.POST("/inquiries/:inquiry_id/payment_link", controllers.AdminInquiryCreatePaymentLink) //legacy
	authorizedWithRoleGroup.POST("/inquiries/payment_link", controllers.AdminCreateBuyerPaymentLink)

	// Sync old inquiry already paid from customer
	// Will be bypass some steps payment from user
	authorizedWithRoleGroup.POST("/inquiries/:inquiry_id/sync_sample", controllers.AdminInquirySyncSample)

	// Inquiry seller
	authorizedWithRoleGroup.GET("/inquiries/:inquiry_id/seller_requests", controllers.AdminPaginateInquirySellerRequests)
	authorizedWithRoleGroup.GET("/inquiries/:inquiry_id/matching_sellers", controllers.AdminPaginateMatchingSellers)
	authorizedWithRoleGroup.POST("/inquiries/:inquiry_id/send_to_seller", controllers.AdminSendInquiryToSeller)
	authorizedWithRoleGroup.PUT("/inquiries/:inquiry_id/update_costing", controllers.AdminUpdateInquiryCosting)
	authorizedWithRoleGroup.GET("/inquiries/:inquiry_id/seller_requests/:inquiry_seller_id/quotation_comments", controllers.AdminInquirySellerRequestPaginateComments)
	authorizedWithRoleGroup.POST("/inquiries/:inquiry_id/seller_requests/:inquiry_seller_id/quotation_comments", controllers.AdminInquirySellerRequestCreateComment)
	authorizedWithRoleGroup.PUT("/inquiries/:inquiry_id/seller_requests/:inquiry_seller_id/quotation_comments/mark_seen", controllers.AdminInquirySellerRequestCommentMarkSeen)
	authorizedWithRoleGroup.GET("/inquiries/:inquiry_id/seller_requests/status_count", controllers.AdminInquirySellerStatusCount)

	authorizedWithRoleGroup.POST("/inquiry_sellers/:inquiry_seller_id/design_comments", controllers.AdminInquirySellerCreateDesignComment)
	authorizedWithRoleGroup.GET("/inquiry_sellers/:inquiry_seller_id/design_comments", controllers.AdminInquirySellerPaginateDesignComments)
	authorizedWithRoleGroup.PUT("/inquiry_sellers/:inquiry_seller_id/design_comments/mark_seen", controllers.AdminInquirySellerPaginateDesignCommentsMarkSeen)

	authorizedWithRoleGroup.POST("/inquiry_sellers/:inquiry_seller_id/approve", controllers.AdminInquirySellerApproveQuotation)
	authorizedWithRoleGroup.DELETE("/inquiry_sellers/:inquiry_seller_id/reject", controllers.AdminInquirySellerRejectQuotation)

	authorizedWithRoleGroup.PUT("/inquiries/:inquiry_id/attachments", controllers.AdminInquiryUpdateAttachments)
	authorizedWithRoleGroup.GET("/inquiries/:inquiry_id/seller_allocations", controllers.AdminInquirySellerAllocationSearchSeller)

	// Seller PO
	authorizedWithRoleGroup.POST("/seller_purchase_orders/:purchase_order_id/preview_payout", controllers.AdminSellerPurchaseOrderPreviewPayout)
	authorizedWithRoleGroup.POST("/seller_purchase_orders/:purchase_order_id/assign_maker", controllers.AdminPurchaseOrderAssignMaker)
	authorizedWithRoleGroup.PUT("/seller_purchase_orders/:purchase_order_id/upload_po", controllers.AdminSellerPurchaseOrderUploadPo)

	authorizedWithRoleGroup.GET("/seller_purchase_orders/:purchase_order_id/po_upload_comments", controllers.AdminSellerPurchaseOrderPaginatePoUploadComments)
	authorizedWithRoleGroup.POST("/seller_purchase_orders/:purchase_order_id/po_upload_comments", controllers.AdminSellerPurchaseOrderAddPoUploadComments)
	authorizedWithRoleGroup.PUT("/seller_purchase_orders/:purchase_order_id/po_upload_comments/mark_seen", controllers.AdminSellerPurchaseOrderPoUploadCommentMarkSeen)
	authorizedWithRoleGroup.GET("/seller_purchase_orders/:purchase_order_id/po_upload_comments/status_count", controllers.AdminSellerPurchasePoUploadCommentStatusCount)

	authorizedWithRoleGroup.PUT("/seller_purchase_orders/:purchase_order_id/mark_design_approval", controllers.AdminSellerPurchaseOrderMarkDesignApproval)
	authorizedWithRoleGroup.PUT("/seller_purchase_orders/:purchase_order_id/update_design", controllers.AdminSellerPurchaseOrderUpdateDesign)

	authorizedWithRoleGroup.GET("/seller_purchase_orders/:purchase_order_id/design_comments", controllers.AdminSellerPurchaseOrderPaginateDesignComments)
	authorizedWithRoleGroup.POST("/seller_purchase_orders/:purchase_order_id/design_comments", controllers.AdminSellerPurchaseOrderDesignCommentCreate)
	authorizedWithRoleGroup.PUT("/seller_purchase_orders/:purchase_order_id/design_comments/mark_seen", controllers.AdminSellerPurchaseOrderDesignCommentMarkSeen)
	authorizedWithRoleGroup.GET("/seller_purchase_orders/:purchase_order_id/design_comments/status_count", controllers.AdminSellerPurchaseDesignCommentStatusCount)

	// Design cloned from buyer after buyer approved
	authorizedWithRoleGroup.PUT("/seller_purchase_orders/:purchase_order_id/update_final_design", controllers.AdminSellerPurchaseOrderUpdateFinalDesign)
	authorizedWithRoleGroup.POST("/seller_purchase_orders/:purchase_order_id/approve_final_design", controllers.AdminSellerPurchaseOrderApproveFinalDesign)

	authorizedWithRoleGroup.GET("/seller_purchase_orders/:purchase_order_id/final_design_comments", controllers.AdminSellerPurchaseOrderPaginateFinalDesignComments)
	authorizedWithRoleGroup.POST("/seller_purchase_orders/:purchase_order_id/final_design_comments", controllers.AdminSellerPurchaseOrderFinalDesignCommentCreate)
	authorizedWithRoleGroup.PUT("/seller_purchase_orders/:purchase_order_id/final_design_comments/mark_seen", controllers.AdminSellerPurchaseOrderFinalDesignCommentMarkSeen)
	authorizedWithRoleGroup.GET("/seller_purchase_orders/:purchase_order_id/final_design_comments/status_count", controllers.AdminSellerPurchaseFinalDesignCommentStatusCount)

	authorizedWithRoleGroup.POST("/seller_purchase_orders/:purchase_order_id/raw_material/send_to_buyer", controllers.AdminSellerPurchaseOrderRawMaterialSendToBuyer)

	authorizedWithRoleGroup.GET("/seller_purchase_orders/:purchase_order_id/raw_material_comments", controllers.AdminSellerPurchaseOrderPaginateRawMaterialComments)
	authorizedWithRoleGroup.POST("/seller_purchase_orders/:purchase_order_id/raw_material_comments", controllers.AdminSellerPurchaseOrderRawMaterialCommentCreate)
	authorizedWithRoleGroup.PUT("/seller_purchase_orders/:purchase_order_id/raw_material_comments/mark_seen", controllers.AdminSellerPurchaseOrderRawMaterialCommentMarkSeen)
	authorizedWithRoleGroup.GET("/seller_purchase_orders/:purchase_order_id/raw_material_comments/status_count", controllers.AdminSellerPurchaseRawMaterialCommentStatusCount)

	authorizedWithRoleGroup.PUT("/seller_purchase_orders/:purchase_order_id/mark_delivered", controllers.AdminSellerPurchaseOrderMarkDelivered)
	authorizedWithRoleGroup.GET("/seller_purchase_orders/:purchase_order_id/logs", controllers.PaginateSellerPurchaseOrderTracking)
	authorizedWithRoleGroup.PUT("/seller_purchase_orders/:purchase_order_id/delivery_feedback", controllers.AdminSellerPurchaseOrderDeliveryFeedback)
	authorizedWithRoleGroup.POST("/seller_purchase_orders/:purchase_order_id/payout", controllers.AdminSellerPurchaseOrderPayout)
	authorizedWithRoleGroup.PUT("/seller_purchase_orders/:purchase_order_id/mark_making", controllers.AdminPurchaseOrderMarkMaking)

	// Seller Bulk PO
	authorizedWithRoleGroup.POST("/seller_purchase_orders/:purchase_order_id/approve_raw_material", controllers.AdminSellerApproveRawMaterials)

	// Notification
	authorizedWithRoleGroup.GET("/notifications", controllers.AdminNotificationList)
	authorizedWithRoleGroup.PUT("/notifications/:notification_id/mark_seen", controllers.AdminNotificationMarkSeen)
	authorizedWithRoleGroup.PUT("/notifications/mark_seen_all", controllers.AdminNotificationMarkSeenAll)

	// Purchase Order ( Sample )
	authorizedWithRoleGroup.GET("/purchase_orders", controllers.PaginatePurchaseOrders)
	authorizedWithRoleGroup.GET("/purchase_orders/export", controllers.ExportPurchaseOrders)
	authorizedWithRoleGroup.GET("/purchase_orders/:purchase_order_id", controllers.GetPurchaseOrder)
	authorizedWithRoleGroup.POST("/purchase_orders", controllers.AdminCreatePurchaseOrder)
	authorizedWithRoleGroup.PUT("/purchase_orders/:purchase_order_id", controllers.AdminUpdatePurchaseOrder)
	authorizedWithRoleGroup.POST("/purchase_orders/:purchase_order_id/approve_design", controllers.AdminApprovePurchaseOrderDesign)
	authorizedWithRoleGroup.PUT("/purchase_orders/:purchase_order_id/update_tracking_status", controllers.UpdatePurchaseOrderTrackingStatus)
	authorizedWithRoleGroup.PUT("/purchase_orders/:purchase_order_id/update_design", controllers.AdminUpdatePurchaseOrderDesign)
	authorizedWithRoleGroup.PUT("/purchase_orders/:purchase_order_id/update_raw_material", controllers.AdminUpdatePurchaseOrderRawMaterial)
	authorizedWithRoleGroup.POST("/purchase_orders/:purchase_order_id/approve_raw_material", controllers.AdminPurchaseOrderApproveRawMaterial)
	authorizedWithRoleGroup.PUT("/purchase_orders/:purchase_order_id/mark_making_without_raw_material", controllers.AdminPurchaseOrderMarkMakingWithoutRawMaterial) // change order to making without raw material
	authorizedWithRoleGroup.PUT("/purchase_orders/:purchase_order_id/mark_making", controllers.AdminPurchaseOrderMarkMaking)
	authorizedWithRoleGroup.PUT("/purchase_orders/:purchase_order_id/mark_submit", controllers.AdminPurchaseOrderMarkSubmit)
	authorizedWithRoleGroup.PUT("/purchase_orders/:purchase_order_id/mark_delivering", controllers.AdmiPurchaseOrderMarkDelivering)
	authorizedWithRoleGroup.PUT("/purchase_orders/:purchase_order_id/mark_delivered", controllers.AdmiPurchaseOrderMarkDelivered)
	authorizedWithRoleGroup.GET("/purchase_orders/:purchase_order_id/logs", controllers.PaginatePurchaseOrderTracking)
	authorizedWithRoleGroup.POST("/purchase_orders/:purchase_order_id/stage_comments", controllers.AdmiPurchaseOrderStageCommentsCreate)

	authorizedWithRoleGroup.GET("/purchase_orders/:purchase_order_id/design_comments", controllers.AdminPurchaseOrderPaginateDesignComments)
	authorizedWithRoleGroup.POST("/purchase_orders/:purchase_order_id/design_comments", controllers.AdminPurchaseOrderAddDesignComments)
	authorizedWithRoleGroup.PUT("/purchase_orders/:purchase_order_id/design_comments/mark_seen", controllers.AdminPurchaseOrderDesignCommentMarkSeen)
	authorizedWithRoleGroup.GET("/purchase_orders/:purchase_order_id/design_comments/status_count", controllers.AdminPurchaseOrderDesignCommentStatusCount)
	authorizedWithRoleGroup.POST("/purchase_orders/payment_link", controllers.AdminMultiPurchaseOrderCreatePaymentLink)

	authorizedWithRoleGroup.PUT("/purchase_orders/:purchase_order_id/assign_pic", controllers.AdminPurchaseOrderAssignPIC)
	authorizedWithRoleGroup.DELETE("/purchase_orders/:purchase_order_id/archive", controllers.AdminArchivePurchaseOrder)
	authorizedWithRoleGroup.DELETE("/purchase_orders/:purchase_order_id/delete", controllers.AdminDeletePurchaseOrder)
	authorizedWithRoleGroup.PUT("/purchase_orders/:purchase_order_id/unarchive", controllers.AdminUnarchivePurchaseOrder)
	authorizedWithRoleGroup.POST("/purchase_orders/:purchase_order_id/notes", controllers.PurchaseOrderAddNote)
	authorizedWithRoleGroup.GET("/purchase_orders/:purchase_order_id/notes", controllers.PaginatePurchaseOrderNotes)
	authorizedWithRoleGroup.PUT("/purchase_orders/:purchase_order_id/notes/mark_seen", controllers.PurchaseOrderNoteMarkSeen)
	authorizedWithRoleGroup.GET("/purchase_orders/:purchase_order_id/notes/unread_count", controllers.PurchaseOrderNoteUnreadCount)
	authorizedWithRoleGroup.DELETE("/purchase_orders/:purchase_order_id/notes/:comment_id", controllers.AdminPurchaseOrderCommentDelete)
	authorizedWithRoleGroup.POST("/purchase_orders/:purchase_order_id/refund", controllers.AdminPurchaseOrderRefund)
	authorizedWithRoleGroup.POST("/purchase_orders/:purchase_order_id/bulk_purchase_order", controllers.AdminCreateBulkPurchaseOrderFromSample)
	authorizedWithRoleGroup.PUT("/purchase_orders/:purchase_order_id/confirm", controllers.AdminPurchaseOrderConfirm)
	authorizedWithRoleGroup.DELETE("/purchase_orders/:purchase_order_id/cancel", controllers.AdminPurchaseOrderCancel)
	authorizedWithRoleGroup.POST("/purchase_orders/:purchase_order_id/mark_as_paid", controllers.AdminPurchaseOrderMarkAsPaid)
	authorizedWithRoleGroup.POST("/purchase_orders/:purchase_order_id/mark_as_unpaid", controllers.AdminPurchaseOrderMarkAsUnpaid)
	authorizedWithRoleGroup.PUT("/purchase_orders/:purchase_order_id/skip_design", controllers.AdminPurchaseOrderSkipDesign)

	authorizedWithRoleGroup.PUT("/purchase_orders/:purchase_order_id/rounds/:round_id/approve", controllers.AdminPurchaseOrderApproveRound)
	authorizedWithRoleGroup.DELETE("/purchase_orders/:purchase_order_id/rounds/:round_id/reject", controllers.AdminPurchaseOrderRejectRound)

	// Bulk purchase order
	authorizedWithRoleGroup.GET("/bulk_purchase_orders", controllers.AdminPaginateBulkPurchaseOrder)
	authorizedWithRoleGroup.GET("/bulk_purchase_orders/export", controllers.AdminExportBulkPurchaseOrder)
	authorizedWithRoleGroup.GET("/bulk_purchase_orders/:bulk_purchase_order_id", controllers.AdminGetBulkPurchaseOrder)
	authorizedWithRoleGroup.PUT("/bulk_purchase_orders/:bulk_purchase_order_id", controllers.AdminUpdateBulkPurchaseOrder)
	authorizedWithRoleGroup.POST("/bulk_purchase_orders/:bulk_purchase_order_id/send_quotation", controllers.AdminSendQuotationToBuyer) //deprecated
	authorizedWithRoleGroup.PUT("/bulk_purchase_orders/:bulk_purchase_order_id/mark_first_payment", controllers.AdminBulkPurchaseOrderMarkFirstPayment)
	authorizedWithRoleGroup.POST("/bulk_purchase_orders/:bulk_purchase_order_id/approve_raw_material", controllers.AdminBulkPurchaseApproveRawMaterial)
	authorizedWithRoleGroup.PUT("/bulk_purchase_orders/:bulk_purchase_order_id/update_raw_material", controllers.AdminBulkPurchaseOrderUpdateRawMaterial)
	authorizedWithRoleGroup.POST("/bulk_purchase_orders/:bulk_purchase_order_id/qc_report", controllers.AdminBulkPurchaseOrderCreateQcReport)
	authorizedWithRoleGroup.GET("/bulk_purchase_orders/:bulk_purchase_order_id/logs", controllers.PaginateBulkPurchaseOrderTracking)
	authorizedWithRoleGroup.PUT("/bulk_purchase_orders/:bulk_purchase_order_id/mark_production", controllers.AdminBulkPurchaseOrderMarkProduction)
	authorizedWithRoleGroup.PUT("/bulk_purchase_orders/:bulk_purchase_order_id/update_production", controllers.AdminBulkPurchaseOrderUpdateProduction)
	authorizedWithRoleGroup.PUT("/bulk_purchase_orders/:bulk_purchase_order_id/mark_raw_material", controllers.AdminBulkPurchaseOrderMarkRawMaterial)
	authorizedWithRoleGroup.PUT("/bulk_purchase_orders/:bulk_purchase_order_id/mark_pps", controllers.AdminBulkPurchaseOrderMarkPps)
	authorizedWithRoleGroup.PUT("/bulk_purchase_orders/:bulk_purchase_order_id/update_pps", controllers.AdminBulkPurchaseOrderUpdatePps)
	authorizedWithRoleGroup.PUT("/bulk_purchase_orders/:bulk_purchase_order_id/mark_qc", controllers.AdminBulkPurchaseOrderMarkQc) // auto change status
	authorizedWithRoleGroup.PUT("/bulk_purchase_orders/:bulk_purchase_order_id/confirm_qc_report", controllers.AdminBulkPurchaseOrderConfirmQcReport)
	authorizedWithRoleGroup.PUT("/bulk_purchase_orders/:bulk_purchase_order_id/mark_final_payment", controllers.AdminBulkPurchaseOrderMarkFinalPayment) // confirm received payment
	authorizedWithRoleGroup.PUT("/bulk_purchase_orders/:bulk_purchase_order_id/mark_delivering", controllers.AdminBulkPurchaseOrderMarkDelivering)
	authorizedWithRoleGroup.PUT("/bulk_purchase_orders/:bulk_purchase_order_id/mark_delivered", controllers.AdminBulkPurchaseOrderMarkDelivered)
	authorizedWithRoleGroup.POST("/bulk_purchase_orders/:bulk_purchase_order_id/approve_qc", controllers.AdminBulkPurchaseApproveQc)
	authorizedWithRoleGroup.PUT("/bulk_purchase_orders/:bulk_purchase_order_id/assign_pic", controllers.AdminBulkPurchaseOrderAssignPIC)
	authorizedWithRoleGroup.GET("/bulk_purchase_orders/:bulk_purchase_order_id/sample_po", controllers.AdminBulkPurchaseOrderGetSamplePO)
	authorizedWithRoleGroup.POST("/bulk_purchase_orders/:bulk_purchase_order_id/invoice", controllers.AdminCreateBulkPurchaseInvoice)
	authorizedWithRoleGroup.POST("/bulk_purchase_orders/:bulk_purchase_order_id/preview_checkout", controllers.AdminBulkPurchaseOrderPreviewCheckout)
	authorizedWithRoleGroup.PUT("/bulk_purchase_orders/:bulk_purchase_order_id/reset", controllers.AdminResetBulkPurchaseOrder)
	authorizedWithRoleGroup.POST("/bulk_purchase_orders/:bulk_purchase_order_id/notes", controllers.AdminBulkPurchaseOrderAddNote)
	authorizedWithRoleGroup.GET("/bulk_purchase_orders/:bulk_purchase_order_id/notes", controllers.PaginateBulkPurchaseOrderNotes)
	authorizedWithRoleGroup.PUT("/bulk_purchase_orders/:bulk_purchase_order_id/notes/mark_seen", controllers.AdminBulkPurchaseOrderNoteMarkSeen)
	authorizedWithRoleGroup.GET("/bulk_purchase_orders/:bulk_purchase_order_id/notes/unread_count", controllers.AdminBulkPurchaseOrderNoteUnreadCount)
	authorizedWithRoleGroup.DELETE("/bulk_purchase_orders/:bulk_purchase_order_id/notes/:comment_id", controllers.AdminBulkPurchaseOrderCommentDelete)
	authorizedWithRoleGroup.POST("/bulk_purchase_orders/:bulk_purchase_order_id/deposit", controllers.AdminBulkPurchaseOrderDeposit)
	authorizedWithRoleGroup.POST("/bulk_purchase_orders/:bulk_purchase_order_id/stage_comments", controllers.AdminBulkPurchaseOrderStageCommentsCreate)
	authorizedWithRoleGroup.POST("/bulk_purchase_orders/:bulk_purchase_order_id/upload_bom", controllers.AdminBulkPurchaseOrderUploadBOM)
	authorizedWithRoleGroup.GET("/bulk_purchase_orders/:bulk_purchase_order_id/bom", controllers.AdminBulkPurchaseOrderPaginateBOM)
	authorizedWithRoleGroup.POST("/bulk_purchase_orders/:bulk_purchase_order_id/submit", controllers.AdminBulkPurchaseOrderSubmit)
	authorizedWithRoleGroup.POST("/bulk_purchase_orders/create_multiple", controllers.AdminCreateMultipleBulks)
	authorizedWithRoleGroup.POST("/bulk_purchase_orders/submit_multiple_quotations", controllers.AdminSubmitMultipleBulkQuotations)
	authorizedWithRoleGroup.POST("/bulk_purchase_orders/upload_file", controllers.AdminUploadBulkPurchaseOrders)

	// Bulk PO - SEller
	authorizedWithRoleGroup.POST("/seller_bulk_purchase_orders/:bulk_purchase_order_id/first_payout", controllers.AdminSellerBulkPurchaseOrderFirstPayout)
	authorizedWithRoleGroup.POST("/seller_bulk_purchase_orders/:bulk_purchase_order_id/final_payout", controllers.AdminSellerBulkPurchaseOrderFinalPayout)
	authorizedWithRoleGroup.POST("/seller_bulk_purchase_orders/:bulk_purchase_order_id/assign_maker", controllers.AdminSellerBulkPurchaseOrderAssignMaker)
	authorizedWithRoleGroup.GET("/seller_bulk_purchase_orders/:bulk_purchase_order_id/seller_allocations", controllers.AdminSellerBulkPurchaseOrderAllocations)
	authorizedWithRoleGroup.GET("/seller_bulk_purchase_orders/:bulk_purchase_order_id/matching_sellers", controllers.AdminSellerPaginateBulkPurchaseOrderMatchingSellers)
	authorizedWithRoleGroup.POST("/seller_bulk_purchase_orders/:bulk_purchase_order_id/send_to_seller", controllers.AdminSellerBulkPurhcaseOrderSendToSeller)
	authorizedWithRoleGroup.GET("/seller_bulk_purchase_orders/:bulk_purchase_order_id/seller_quotations", controllers.AdminSellerPaginateBulkPurchaseOrderSellerQuotations)
	authorizedWithRoleGroup.PUT("/seller_bulk_purchase_orders/:bulk_purchase_order_id/upload_po", controllers.AdminSellerBulkPurchaseOrderUploadPo)
	authorizedWithRoleGroup.PUT("/seller_bulk_purchase_orders/:bulk_purchase_order_id/upload_po", controllers.AdminSellerBulkPurchaseOrderUploadPo)

	authorizedWithRoleGroup.PUT("/seller_bulk_purchase_orders/:bulk_purchase_order_id/product_photo", controllers.AdminSellerBulkPurchaseOrderUpdateProductPhoto)
	authorizedWithRoleGroup.PUT("/seller_bulk_purchase_orders/:bulk_purchase_order_id/techpack", controllers.AdminSellerBulkPurchaseOrderUpdateTechpack)
	authorizedWithRoleGroup.PUT("/seller_bulk_purchase_orders/:bulk_purchase_order_id/bill_of_material", controllers.AdminSellerBulkPurchaseOrderUpdateBillOfMaterial)
	authorizedWithRoleGroup.PUT("/seller_bulk_purchase_orders/:bulk_purchase_order_id/size_chart", controllers.AdminSellerBulkPurchaseOrderUpdateSizeChart)
	authorizedWithRoleGroup.PUT("/seller_bulk_purchase_orders/:bulk_purchase_order_id/size_spec", controllers.AdminSellerBulkPurchaseOrderUpdateSizeSpec)
	authorizedWithRoleGroup.PUT("/seller_bulk_purchase_orders/:bulk_purchase_order_id/size_grading", controllers.AdminSellerBulkPurchaseOrderUpdateSizeGrading)
	authorizedWithRoleGroup.PUT("/seller_bulk_purchase_orders/:bulk_purchase_order_id/inspection_procedure", controllers.AdminSellerBulkPurchaseOrderUpdateInspectionProcedure)
	authorizedWithRoleGroup.PUT("/seller_bulk_purchase_orders/:bulk_purchase_order_id/inspection_testing_requirements", controllers.AdminSellerBulkPurchaseUpdateInspectionTestingRequirements)
	authorizedWithRoleGroup.PUT("/seller_bulk_purchase_orders/:bulk_purchase_order_id/point_of_measurement", controllers.AdminSellerBulkPurchaseOrderUpdatePointOfMeasurement)
	authorizedWithRoleGroup.PUT("/seller_bulk_purchase_orders/:bulk_purchase_order_id/label_guide", controllers.AdminSellerBulkPurchaseOrderUpdateLabelGuide)
	authorizedWithRoleGroup.PUT("/seller_bulk_purchase_orders/:bulk_purchase_order_id/mark_delivered", controllers.AdminSellerBulkPurchaseOrderMarkDelivered)
	authorizedWithRoleGroup.PUT("/seller_bulk_purchase_orders/:bulk_purchase_order_id/mark_raw_material", controllers.AdminSellerBulkPurchaseOrderMarkRawMaterial)
	authorizedWithRoleGroup.PUT("/seller_bulk_purchase_orders/:bulk_purchase_order_id/update_pps", controllers.AdminSellerBulkPurchaseOrderUpdatePps)
	authorizedWithRoleGroup.POST("/seller_bulk_purchase_orders/:seller_quotation_id/approve", controllers.AdminSellerBulkPurchaseOrderApproveSellerQuotation)
	authorizedWithRoleGroup.DELETE("/seller_bulk_purchase_orders/:seller_quotation_id/reject", controllers.AdminSellerBulkPurchaseOrderRejectSellerQuotation)

	// Settings
	authorizedWithRoleGroup.GET("/settings/taxes", controllers.PaginateSettingTax)
	authorizedWithRoleGroup.GET("/settings/taxes/:tax_id", controllers.GetSettingTax)
	authorizedWithRoleGroup.POST("/settings/taxes", controllers.CreateSettingTax)
	authorizedWithRoleGroup.PUT("/settings/taxes/:tax_id", controllers.UpdateSettingTax)
	authorizedWithRoleGroup.DELETE("/settings/taxes/:tax_id", controllers.DeleteSettingTaxes)

	authorizedWithRoleGroup.GET("/settings/sizes", controllers.PaginateSettingSizes)
	authorizedWithRoleGroup.POST("/settings/sizes", controllers.CreateSettingSizes)
	authorizedWithRoleGroup.PUT("/settings/sizes", controllers.UpdateSettingSizes)
	authorizedWithRoleGroup.PUT("/settings/sizes/:size_id", controllers.UpdateSettingSize)
	authorizedWithRoleGroup.GET("/settings/sizes/:size_id", controllers.GetSettingSize)
	authorizedWithRoleGroup.DELETE("/settings/sizes/:size_id", controllers.DeleteSettingSize)
	authorizedWithRoleGroup.PUT("/settings/sizes/:size_id/sizes/:size_id", controllers.UpdateSettingSize)
	authorizedWithRoleGroup.DELETE("/settings/sizes/:size_id/sizes/:size_id", controllers.DeleteSettingSize)
	authorizedWithRoleGroup.GET("/settings/sizes/types/:type", controllers.GetSettingSizeType)
	authorizedWithRoleGroup.PUT("/settings/sizes/types/:type", controllers.UpdateSettingSizeType)
	authorizedWithRoleGroup.DELETE("/settings/sizes/types/:type", controllers.DeleteSettingSizeType)

	authorizedWithRoleGroup.GET("/settings/banks", controllers.PaginateSettingBanks)
	authorizedWithRoleGroup.POST("/settings/banks", controllers.CreateSettingBanks)
	authorizedWithRoleGroup.PUT("/settings/banks", controllers.UpdateSettingBanks)
	authorizedWithRoleGroup.DELETE("/settings/banks/:id", controllers.DeleteSettingBanks)
	authorizedWithRoleGroup.DELETE("/settings/banks/countries/:country_code", controllers.DeleteSettingBanksByCountryCode)

	authorizedWithRoleGroup.GET("/settings/seo", controllers.PaginateSettingSEO)
	authorizedWithRoleGroup.GET("/settings/seo/group", controllers.PaginateSettingRouteGroupSEO)
	authorizedWithRoleGroup.POST("/settings/seo", controllers.CreateSettingSEO)
	authorizedWithRoleGroup.PUT("/settings/seo/:id", controllers.UpdateSettingSEO)
	authorizedWithRoleGroup.DELETE("/settings/seo/:id", controllers.DeleteSettingSEO)
	authorizedWithRoleGroup.PATCH("/settings/seo/translations", controllers.PatchSEOTranslations)

	authorizedWithRoleGroup.POST("/settings/docs", controllers.CreateSettingDoc)
	authorizedWithRoleGroup.GET("/settings/docs/:type", controllers.GetSettingDoc)
	authorizedWithRoleGroup.PUT("/settings/docs/:type", controllers.UpdateSettingDoc)

	authorizedWithRoleGroup.POST("/settings/inquiries", controllers.CreateSettingInquiry)
	authorizedWithRoleGroup.GET("/settings/inquiries/:type", controllers.GetSettingInquiry)
	authorizedWithRoleGroup.PUT("/settings/inquiries/:type", controllers.UpdateSettingInquiry)

	// Resources management
	authorizedWithRoleGroup.GET("/resources", controllers.ResourcesManagement)
	authorizedWithRoleGroup.PATCH("/resources/policy", controllers.ReloadResourcesPolicy)
	authorizedWithRoleGroup.GET("/resources/policy", controllers.GetResourcesPolicy)

	// Payment Transactions
	authorizedWithRoleGroup.GET("/payment_transactions", controllers.PaginatePaymentTransaction)
	authorizedWithRoleGroup.GET("/payment_transactions/export", controllers.ExportPaymentTransactions)
	authorizedWithRoleGroup.GET("/payment_transactions/:payment_transaction_id", controllers.GetPaymentTransaction)
	authorizedWithRoleGroup.PUT("/payment_transactions/:payment_transaction_id/approve", controllers.ApprovePaymentTransaction)
	authorizedWithRoleGroup.PUT("/payment_transactions/:payment_transaction_id/reject", controllers.RejectPaymentTransaction)
	authorizedWithRoleGroup.GET("/payment_transactions/:payment_transaction_id/attachments", controllers.GetPaymentTransactionAttachments)

	// Ads videos
	authorizedWithRoleGroup.GET("/ads_videos", controllers.PaginateAsFeaturedIns)
	authorizedWithRoleGroup.POST("/ads_videos", controllers.CreateAdsVideo)
	authorizedWithRoleGroup.PUT("/ads_videos/:ads_video_id", controllers.UpdateAdsVideo)
	authorizedWithRoleGroup.DELETE("/ads_videos/:ads_video_id", controllers.DeleteAdsVideo)

	// As featured in
	authorizedWithRoleGroup.GET("/as_featured_ins", controllers.PaginateAsFeaturedIns)
	authorizedWithRoleGroup.POST("/as_featured_ins", controllers.CreateAsFeaturedIn)
	authorizedWithRoleGroup.PUT("/as_featured_ins/:as_featured_in_id", controllers.UpdateAsFeaturedIn)
	authorizedWithRoleGroup.DELETE("/as_featured_ins/:as_featured_in_id", controllers.DeleteAsFeaturedIn)

	// Invoice
	authorizedWithRoleGroup.GET("/invoices", controllers.PaginateInvoices)
	authorizedWithRoleGroup.GET("/invoices/:invoice_number", controllers.DetailsInvoice)
	authorizedWithRoleGroup.GET("/invoices/invoice_number/next", controllers.NextInvoiceNumber)
	authorizedWithRoleGroup.GET("/invoices/:invoice_number/exits", controllers.ExitsInvoiceNumber)
	authorizedWithRoleGroup.POST("/invoices", controllers.CreateInvoice)
	authorizedWithRoleGroup.PUT("/invoices/:invoice_number", controllers.UpdateInvoice)
	authorizedWithRoleGroup.GET("/invoices/:invoice_number/attachment", controllers.GetInvoiceAttachment)

	// Invoice
	authorizedWithRoleGroup.GET("/analytics/inquiries/potential_overdue", controllers.PaginatePotentialOverdueInquiries)
	authorizedWithRoleGroup.GET("/analytics/inquiries/timeline", controllers.PaginateInquiriesTimeline)

	// Data Analytic - Open Market
	authorizedWithRoleGroup.GET("/data_analytics/platforms/overview", controllers.OverviewDataAnalyticPlatform)
	authorizedWithRoleGroup.GET("/data_analytics/products/search", controllers.SearchDAProducts)
	authorizedWithRoleGroup.GET("/data_analytics/products/:product_id", controllers.GetDAProductDetails)
	authorizedWithRoleGroup.GET("/data_analytics/products/:product_id/chart", controllers.GetDAProductChart)
	authorizedWithRoleGroup.GET("/data_analytics/products/top", controllers.TopDAProducts)
	authorizedWithRoleGroup.GET("/data_analytics/products/top_movers", controllers.TopMoversDAProducts)
	authorizedWithRoleGroup.GET("/data_analytics/categories/top", controllers.TopDACategories)
	authorizedWithRoleGroup.GET("/data_analytics/sub_categories/top", controllers.TopDASubCategories)

	// Data Analytic - Product Trending
	authorizedWithRoleGroup.GET("/data_analytics/product_trendings", controllers.PaginateProductTrending)
	authorizedWithRoleGroup.POST("/data_analytics/product_trendings", controllers.CreateProductTrending)
	authorizedWithRoleGroup.PUT("/data_analytics/product_trendings", controllers.UpdateProductTrending)
	authorizedWithRoleGroup.DELETE("/data_analytics/product_trendings", controllers.DeleteProductTrending)
	authorizedWithRoleGroup.GET("/data_analytics/product_trendings/:product_trending_id", controllers.GetProductTrending)
	authorizedWithRoleGroup.GET("/data_analytics/product_trendings/:product_trending_id/chart", controllers.GetProductTrendingChart)
	authorizedWithRoleGroup.GET("/data_analytics/product_trendings/group", controllers.PaginateProductTrendingGroup)
	authorizedWithRoleGroup.GET("/data_analytics/product_trendings/domains", controllers.ListProductTrendingDomain)
	authorizedWithRoleGroup.GET("/data_analytics/product_trendings/categories", controllers.ListProductTrendingCategory)
	authorizedWithRoleGroup.GET("/data_analytics/product_trendings/sub_categories", controllers.ListProductTrendingSubCategory)

	// Data Analytic - Trending
	authorizedWithRoleGroup.GET("/data_analytics/products", controllers.PaginateAnalyticProduct)
	authorizedWithRoleGroup.GET("/data_analytics/products/group", controllers.PaginateAnalyticProductGroup)
	authorizedWithRoleGroup.GET("/data_analytics/products/product_classes/group", controllers.GetAnalyticProductClassGroup)
	authorizedWithRoleGroup.GET("/data_analytics/products/trending/group", controllers.GetAnalyticProductTrendingGroup)

	authorizedWithRoleGroup.POST("/trendings", controllers.CreateTrending)
	authorizedWithRoleGroup.GET("/trendings", controllers.PaginateTrendings)
	authorizedWithRoleGroup.GET("/trendings/stats", controllers.ListTrendingStats)
	authorizedWithRoleGroup.GET("/trendings/:id", controllers.GetTrending)
	authorizedWithRoleGroup.PUT("/trendings/:id", controllers.UpdateTrending)
	authorizedWithRoleGroup.DELETE("/trendings/:id", controllers.DeleteTrending)

	authorizedWithRoleGroup.PUT("/trendings/product_trendings/add", controllers.AddProductToTrending)
	authorizedWithRoleGroup.PUT("/trendings/product_trendings/remove", controllers.RemoveProductFromTrending)

	// Data Analytic - Inflow
	authorizedWithRoleGroup.GET("/data_analytics/users/new", controllers.DANewUsers)
	authorizedWithRoleGroup.GET("/data_analytics/catalog_products/new", controllers.DANewCatalogProducts)
	authorizedWithRoleGroup.GET("/data_analytics/inquiries", controllers.DAInquiries)
	authorizedWithRoleGroup.GET("/data_analytics/purchase_orders", controllers.DANewPurchaseOrders)
	authorizedWithRoleGroup.GET("/data_analytics/bulk_purchase_orders", controllers.DANewBulkPurchaseOrders)
	authorizedWithRoleGroup.GET("/data_analytics/ops_biz_performance", controllers.DAOpsBizPerformance)

	// Fabrics
	authorizedWithRoleGroup.POST("/fabrics", controllers.CreateFabric)
	authorizedWithRoleGroup.PUT("/fabrics/:id", controllers.UpdateFabric)
	authorizedWithRoleGroup.GET("/fabrics", controllers.PaginateFabric)
	authorizedWithRoleGroup.GET("/fabrics/:id", controllers.DetailsFabric)
	authorizedWithRoleGroup.DELETE("/fabrics/**//:id", controllers.DeleteFabric)
	authorizedWithRoleGroup.DELETE("/fabrics/:id/add_collection", controllers.DeleteFabric)

	authorizedWithRoleGroup.POST("/fabric_collections", controllers.CreateFabricCollection)
	authorizedWithRoleGroup.PUT("/fabric_collections/:id", controllers.UpdateFabricCollection)
	authorizedWithRoleGroup.GET("/fabric_collections", controllers.PaginateFabricCollection)
	authorizedWithRoleGroup.GET("/fabric_collections/:id", controllers.DetailsFabricCollection)
	authorizedWithRoleGroup.DELETE("/fabric_collections/:id", controllers.DeleteFabricCollection)
	authorizedWithRoleGroup.POST("/fabric_collections/:id/add_fabric", controllers.AddFabricToCollection)
	authorizedWithRoleGroup.DELETE("/fabric_collections/:id/remove_fabric", controllers.RemoveFabricFromCollection)

	// Time And Action
	authorizedWithRoleGroup.POST("/tnas", controllers.AdminCreateTNA)
	authorizedWithRoleGroup.PUT("/tnas/:id", controllers.AdminUpdateTNA)
	authorizedWithRoleGroup.GET("/tnas", controllers.AdminPaginateTNA)
	authorizedWithRoleGroup.DELETE("/tnas/:id", controllers.AdminDeleteTNA)

	// Release Notes
	authorizedWithRoleGroup.POST("/release_notes", controllers.CreateReleaseNote)
	authorizedWithRoleGroup.PUT("/release_notes/:id", controllers.UpdateReleaseNote)
	authorizedWithRoleGroup.GET("/release_notes", controllers.PaginateReleaseNote)
	authorizedWithRoleGroup.DELETE("/release_notes/:id", controllers.DeleteReleaseNote)

	// Document
	authorizedWithRoleGroup.POST("/documents", controllers.AdminCreateDocument)
	authorizedWithRoleGroup.GET("/documents", controllers.AdminGetDocumentList)
	authorizedWithRoleGroup.PUT("/documents/:document_id", controllers.AdminUpdateDocument)
	authorizedWithRoleGroup.GET("/documents/:slug", controllers.AdminGetDocument)
	authorizedWithRoleGroup.DELETE("/documents/:document_id", controllers.AdminDeleteDocument)
	// Document Category
	authorizedWithRoleGroup.POST("/document_categories", controllers.AdminCreateDocumentCategory)
	authorizedWithRoleGroup.GET("/document_categories", controllers.AdminGetDocumentCategoryList)
	authorizedWithRoleGroup.PUT("/document_categories/:document_category_id", controllers.AdminUpdateDocumentCategory)
	authorizedWithRoleGroup.DELETE("/document_categories/:document_category_id", controllers.AdminDeleteDocumentCategory)
	// Document Tag
	authorizedWithRoleGroup.POST("/document_tags", controllers.AdminCreateDocumentTag)
	authorizedWithRoleGroup.GET("/document_tags", controllers.AdminGetDocumentTagList)
	authorizedWithRoleGroup.PUT("/document_tags/:document_tag_id", controllers.AdminUpdateDocumentTag)
	authorizedWithRoleGroup.DELETE("/document_tags/:document_tag_id", controllers.AdminDeleteDocumentTag)
	// Chat
	authorizedWithRoleGroup.POST("/chat_messages", controllers.AdminCreateChatMessage)
	authorizedWithRoleGroup.GET("/chat_messages", controllers.AdminGetChatMessageList)
	authorizedWithRoleGroup.GET("/chat_messages/unseen_messages", controllers.AdminCountUnseenChatMessage)
	//Chat Room
	authorizedWithRoleGroup.GET("/chat_rooms/relevant_stage", controllers.AdminGetChatUserRelevantStage)
	authorizedWithRoleGroup.POST("/chat_rooms", controllers.AdminCreateChatRoom)
	authorizedWithRoleGroup.GET("/chat_rooms", controllers.AdminGetChatRoomList)
	authorizedWithRoleGroup.PUT("/chat_rooms/:chat_room_id/seen_messages", controllers.AdminMarkSeenChatRoomMessage)
	authorizedWithRoleGroup.GET("/chat_rooms/unseen_message", controllers.AdminCountUnseenChatMessageOnRoom)
	// Order Group
	authorizedWithRoleGroup.GET("/order_groups", controllers.AdminGetOrderGroupList)
	authorizedWithRoleGroup.POST("/order_groups", controllers.AdminCreateOrderGroup)
	authorizedWithRoleGroup.GET("/order_groups/:order_group_id", controllers.AdminGetOrderGroupDetail)
	authorizedWithRoleGroup.POST("/order_groups/assign", controllers.AdminAssignOrderGroup)
	// Order Cart
	authorizedWithRoleGroup.POST("/order_cart/:buyer_id/preview", controllers.GetBuyerOrderCartPreview)
	authorizedWithRoleGroup.POST("/order_cart/:buyer_id/create_payment_link", controllers.CreateBuyerPaymentLink)

	// Product File Upload Info
	authorizedWithRoleGroup.POST("/product_file_upload_infos/upload", controllers.UploadProductFile)
	authorizedWithRoleGroup.GET("/product_file_upload_infos", controllers.GetProductFileUploadInfoList)
}
