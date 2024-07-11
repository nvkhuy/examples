package controllers

import (
	"fmt"

	"github.com/engineeringinflow/inflow-backend/pkg/customerio"
	"github.com/engineeringinflow/inflow-backend/pkg/models"
	"github.com/engineeringinflow/inflow-backend/pkg/models/enums"
	"github.com/engineeringinflow/inflow-backend/pkg/repo"
	"github.com/engineeringinflow/inflow-backend/services/consumer/tasks"
	"github.com/labstack/echo/v4"
	"github.com/rotisserie/eris"
	"github.com/samber/lo"
)

// PaginatePurchaseOrders
// @Tags Marketplace-PO
// @Summary Purchase orders
// @Description Purchase orders
// @Accept  json
// @Produce  json
// @Param keyword query string false "Keyword"
// @Param page query int false "Page number"
// @Param limit query int false "Limit"
// @Success 200 {object} query.Pagination{records=[]models.PurchaseOrder}
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/buyer/purchase_orders [get]
func PaginatePurchaseOrders(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	var params repo.PaginatePurchaseOrdersParams
	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}
	params.JwtClaimsInfo = claims
	params.IncludeUnpaidWithoutInquiry = true
	params.IncludeItems = true
	params.IncludeCollection = true
	var result = repo.NewPurchaseOrderRepo(cc.App.DB).PaginatePurchaseOrders(params)

	return cc.Success(result)
}

// GetPurchaseOrder
// @Tags Marketplace-PO
// @Summary Get purchase order
// @Description Get purchase order
// @Accept  json
// @Produce  json
// @Param purchase_order_id query string true "ID"
// @Success 200 {object} models.PurchaseOrder
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/buyer/purchase_orders/{purchase_order_id} [get]
func GetPurchaseOrder(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, "")
	}

	var params repo.GetPurchaseOrderParams
	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, "")
	}

	params.JwtClaimsInfo = claims
	params.IncludeSampleMaker = true
	params.UserID = claims.GetUserID()
	result, err := repo.NewPurchaseOrderRepo(cc.App.DB).GetPurchaseOrder(params)
	if err != nil {
		return eris.Wrap(err, "")
	}

	return cc.Success(result)
}

// PurchaseOrderFeedback
// @Tags Marketplace-PO
// @Summary Get purchase order
// @Description Get purchase order
// @Accept  json
// @Produce  json
// @Param purchase_order_id query string true "ID"
// @Success 200 {object} models.PurchaseOrder
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/buyer/purchase_orders/{purchase_order_id}/feedback [get]
func PurchaseOrderFeedback(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	_, err := cc.GetJwtClaims()
	if err != nil {
		return eris.Wrap(err, "")
	}

	var params repo.PurchaseOrderFeedbackParams
	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, "")
	}

	err = repo.NewPurchaseOrderRepo(cc.App.DB).BuyerGivePurchaseOrderFeedback(params)
	if err != nil {
		return eris.Wrap(err, "")
	}

	return cc.Success("Received")
}

// UpdatePurchaseOrderTrackingStatus
// @Tags Marketplace-PO
// @Summary update purchase order
// @Description update purchase order
// @Accept  json
// @Produce  json
// @Param purchase_order_id query string true "ID"
// @Success 200 {object} models.PurchaseOrder
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/buyer/purchase_orders/{purchase_order_id}/update_tracking_status [put]
func UpdatePurchaseOrderTrackingStatus(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	claims, err := cc.GetJwtClaims()
	if err != nil {
		return eris.Wrap(err, "")
	}

	var params repo.UpdatePurchaseOrderTrackingStatusParams
	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, "")
	}

	params.UserID = claims.ID

	result, err := repo.NewPurchaseOrderRepo(cc.App.DB).UpdatePurchaseOrderTrackingStatus(params)
	if err != nil {
		return eris.Wrap(err, "")
	}

	return cc.Success(result)
}

// UpdatePurchaseOrderTrackingStatus
// @Tags Marketplace-PO
// @Summary update purchase order
// @Description update purchase order
// @Accept  json
// @Produce  json
// @Param purchase_order_id query string true "ID"
// @Success 200 {object} models.PurchaseOrder
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/buyer/purchase_orders/{purchase_order_id}/logs [get]
func PaginatePurchaseOrderTracking(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, "")
	}

	var params repo.PaginatePurchaseOrderTrackingParams
	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, "")
	}

	params.JwtClaimsInfo = claims
	params.UserGroup = enums.PoTrackingUserGroupBuyer
	params.UserID = claims.GetUserID()
	var result = repo.NewPurchaseOrderTrackingRepo(cc.App.DB).PaginatePurchaseOrderTrackings(params)

	return cc.Success(result)
}

// BuyerApprovePurchaseOrderDesign
// @Tags Marketplace-Inquiry
// @Summary Approve inquiry quotation
// @Description Approve inquiry quotation
// @Accept  json
// @Produce  json
// @Param data body models.BuyerApproveInquiryQuotationForm true "Form"
// @Success 200 {object} models.Inquiry
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/buyer/purchase_orders/{purchase_order_id}/approve_design [put]
func BuyerApprovePurchaseOrderDesign(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.BuyerApproveDesignParams

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, "")
	}

	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, "")
	}

	params.JwtClaimsInfo = claims
	purchaseOrder, err := repo.NewPurchaseOrderRepo(cc.App.DB).BuyerApproveDesign(params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	for _, assigneeID := range purchaseOrder.AssigneeIDs {
		tasks.TrackCustomerIOTask{
			UserID: assigneeID,
			Event:  customerio.EventPoBuyerApproveDesign,
			Data:   purchaseOrder.GetCustomerIOMetadata(nil),
		}.Dispatch(c.Request().Context())

	}

	return cc.Success("Approved")
}

// BuyerApprovePurchaseOrderDesign
// @Tags Marketplace-Inquiry
// @Summary Approve inquiry quotation
// @Description Approve inquiry quotation
// @Accept  json
// @Produce  json
// @Param data body models.BuyerApproveInquiryQuotationForm true "Form"
// @Success 200 {object} models.Inquiry
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/buyer/purchase_orders/{purchase_order_id}/reject_design [put]
func BuyerRejectPurchaseOrderDesign(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.BuyerRejectDesignParams

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, "")
	}

	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, "")
	}

	params.JwtClaimsInfo = claims
	purchaseOrder, err := repo.NewPurchaseOrderRepo(cc.App.DB).BuyerRejectDesign(params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	for _, assigneeID := range purchaseOrder.AssigneeIDs {
		tasks.TrackCustomerIOTask{
			UserID: assigneeID,
			Event:  customerio.EventPoBuyerRejectDesign,
			Data:   purchaseOrder.GetCustomerIOMetadata(nil),
		}.Dispatch(c.Request().Context())
	}

	return cc.Success("Rejected")
}

// BuyerApprovePurchaseOrderDesign
// @Tags Marketplace-Inquiry
// @Summary Approve inquiry quotation
// @Description Approve inquiry quotation
// @Accept  json
// @Produce  json
// @Param data body models.BuyerApproveInquiryQuotationForm true "Form"
// @Success 200 {object} models.Inquiry
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/buyer/purchase_orders/{purchase_order_id}/confirm_delivered [post]
func BuyerPurchaseOrderConfirmDelivered(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.BuyerConfirmDeliveredParams

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, "")
	}

	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, "")
	}

	params.JwtClaimsInfo = claims
	purchaseOrder, err := repo.NewPurchaseOrderRepo(cc.App.DB).BuyerConfirmDelivered(params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	for _, assigneeID := range purchaseOrder.AssigneeIDs {
		tasks.TrackCustomerIOTask{
			UserID: assigneeID,
			Event:  customerio.EventPoBuyerConfirmDelivered,
			Data:   purchaseOrder.GetCustomerIOMetadata(nil),
		}.Dispatch(c.Request().Context())
	}

	return cc.Success("Confirm Delivered")
}

// BuyerPurchaseOrderAddDesignComments
// @Tags Marketplace-Inquiry
// @Summary Add design comments
// @Description Add design comments
// @Accept  json
// @Produce  json
// @Param data body models.Comment true "Form"
// @Success 200 {object} models.Inquiry
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/buyer/purchase_orders/{purchase_order_id}/design_comments [post]
func BuyerPurchaseOrderAddDesignComments(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params models.CommentCreateForm

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, "")
	}
	params.TargetType = enums.CommentTargetTypePurchaseOrderDesign
	var orderID = cc.GetPathParamString("purchase_order_id")
	params.TargetID = orderID
	params.JwtClaimsInfo = claims

	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, "")
	}

	purchaseOrder, err := repo.NewPurchaseOrderRepo(cc.App.DB).GetPurchaseOrderShortInfo(orderID)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	result, err := repo.NewCommentRepo(cc.App.DB).CreateComment(params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	// CMS notification
	tasks.CreateCmsNotificationTask{
		Message:          fmt.Sprintf("New comment on sample purchase order %s", purchaseOrder.ReferenceID),
		NotificationType: enums.CmsNotificationTypePoDesignNewComment,
		Metadata: &models.NotificationMetadata{
			CommentID:                result.ID,
			PurchaseOrderID:          purchaseOrder.ID,
			PurchaseOrderReferenceID: purchaseOrder.ReferenceID,
		},
	}.Dispatch(c.Request().Context())

	if userInfo, err := repo.NewUserRepo(cc.App.DB).GetShortUserInfo(claims.GetUserID()); err == nil {
		for _, userID := range purchaseOrder.AssigneeIDs {
			tasks.TrackCustomerIOTask{
				UserID: userID,
				Event:  customerio.EventBuyerNewDesignComment,
				Data: result.GetCustomerIOMetadata(map[string]interface{}{
					"sender":       userInfo.GetCustomerIOMetadata(nil),
					"admin_po_url": fmt.Sprintf("%s/samples/%s/customer?open_design_comments=true", cc.App.Config.AdminPortalBaseURL, purchaseOrder.ID),
					"brand_po_url": fmt.Sprintf("%s/samples/%s?open_design_comments=true", cc.App.Config.BrandPortalBaseURL, purchaseOrder.ID),
				}),
			}.Dispatch(c.Request().Context())
		}
	}

	return cc.Success(result)
}

// BuyerPurchaseOrderPaginateDesignComments
// @Tags Marketplace-Inquiry
// @Summary Get design comments
// @Description Get design comments
// @Accept  json
// @Produce  json
// @Param data body models.Comment true "Form"
// @Success 200 {object} models.Inquiry
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/buyer/purchase_orders/{purchase_order_id}/design_comments [get]
func BuyerPurchaseOrderPaginateDesignComments(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var orderID = cc.GetPathParamString("purchase_order_id")
	var params repo.PaginateCommentsParams

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, "")
	}

	params.JwtClaimsInfo = claims
	params.TargetID = orderID
	params.TargetType = enums.CommentTargetTypePurchaseOrderDesign

	var results = repo.NewCommentRepo(cc.App.DB).PaginateComment(params)

	return cc.Success(results)
}

// BuyerPurchaseOrderDesignCommentMarkSeen
// @Tags Marketplace-PO
// @Summary Add design comments
// @Description Add design comments
// @Accept  json
// @Produce  json
// @Param data body models.Comment true "Form"
// @Success 200 {object} models.Inquiry
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/buyer/purchase_orders/{purchase_order_id}/design_comments/mark_seen [put]
func BuyerPurchaseOrderDesignCommentMarkSeen(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.PoCommentMarkSeenParams

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, "")
	}
	params.JwtClaimsInfo = claims

	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, "")
	}

	err = repo.NewPurchaseOrderRepo(cc.App.DB).DesignCommentMarkSeen(params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	return cc.Success("Seen")
}

// BuyerPurchaseOrderDesignCommentStatusCount
// @Tags Marketplace-PO
// @Summary Design comment status count
// @Description Mark seen comments
// @Accept  json
// @Produce  json
// @Param data body models.Comment true "Form"
// @Success 200 {object} models.Inquiry
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/buyer/purchase_orders/{purchase_order_id}/design_comments/status_count [get]
func BuyerPurchaseOrderDesignCommentStatusCount(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.PoCommentStatusCountParams

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, "")
	}
	params.JwtClaimsInfo = claims

	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, "")
	}

	result, err := repo.NewPurchaseOrderRepo(cc.App.DB).DesignCommentStatusCount(params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	return cc.Success(result)
}

// PurchaseOrderBuyerApproveRawMaterial
// @Tags Marketplace-PO
// @Summary update purchase order
// @Description update purchase order
// @Accept  json
// @Produce  json
// @Param purchase_order_id query string true "ID"
// @Success 200 {object} models.PurchaseOrder
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/buyer/purchase_orders/{purchase_order_id}/approve_raw_material [post]
func PurchaseOrderBuyerApproveRawMaterial(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.PurchaseApproveRawMaterialParams

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, "")
	}

	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, "")
	}

	params.JwtClaimsInfo = claims
	order, err := repo.NewPurchaseOrderRepo(cc.App.DB).PurchaseOrderApproveRawMaterial(params)
	if err != nil {
		return eris.Wrap(err, "")
	}

	for _, assigneeID := range order.AssigneeIDs {
		var data = order.GetCustomerIOMetadata(nil)
		if order.PoRawMaterials != nil && len(*order.PoRawMaterials) > 0 {
			var list models.PoRawMaterialMetas = lo.Filter(*order.PoRawMaterials, func(item *models.PoRawMaterialMeta, index int) bool {
				return lo.Contains(params.ItemIDs, item.ReferenceID)
			})

			data["approved_raw_materials"] = list.GenerateFileURL()
		}
		tasks.TrackCustomerIOTask{
			UserID: assigneeID,
			Event:  customerio.EventPoBuyerApproveRawMaterial,
			Data:   data,
		}.Dispatch(c.Request().Context())
	}

	return cc.Success(order)
}

// UpdatePurchaseOrder
// @Tags Marketplace-PO
// @Summary update purchase order
// @Description update purchase order
// @Accept  json
// @Produce  json
// @Param purchase_order_id query string true "ID"
// @Success 200 {object} models.PurchaseOrder
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/buyer/purchase_orders/{purchase_order_id} [put]
func UpdatePurchaseOrder(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.UpdatePurchaseOrderParams

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, "")
	}

	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, "")
	}

	params.JwtClaimsInfo = claims
	params.UserID = claims.GetUserID()
	order, err := repo.NewPurchaseOrderRepo(cc.App.DB).UpdatePurchaseOrder(params)
	if err != nil {
		return eris.Wrap(err, "")
	}

	for _, assigneeID := range order.AssigneeIDs {
		var data = order.GetCustomerIOMetadata(nil)
		_, _ = tasks.TrackCustomerIOTask{
			UserID: assigneeID,
			Event:  customerio.EventPoBuyerUpdated,
			Data:   data,
		}.Dispatch(c.Request().Context())
	}

	return cc.Success(order)
}

// BuyerCreateBulkPurchaseOrderFromSample
// @Tags Marketplace-Inquiry
// @Summary Create bulk
// @Description Create bulk
// @Accept  json
// @Produce  json
// @Param data body repo.InquiryCheckoutParams true "Form"
// @Success 200 {object} models.PurchaseOrder
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router api/v1/buyer/purchase_orders/{purchase_order_id}/bulk_purchase_order [post]
func BuyerCreateBulkPurchaseOrderFromSample(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.CreateBulkPurchaseOrderParams

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, "")
	}

	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	params.JwtClaimsInfo = claims

	result, err := repo.NewPurchaseOrderRepo(cc.App.DB).CreateBulkPurchaseOrder(params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}
	_, _ = tasks.CreateChatRoomTask{
		UserID:              claims.GetUserID(),
		Role:                claims.GetRole(),
		BulkPurchaseOrderID: result.ID,
		BuyerID:             result.UserID,
	}.Dispatch(c.Request().Context())

	return cc.Success(result)
}

// GetPurchaseOrderInvoice
// @Tags Marketplace-PO
// @Summary Get purchase order
// @Description Get purchase order
// @Accept  json
// @Produce  json
// @Param purchase_order_id query string true "ID"
// @Success 200 {object} models.PurchaseOrder
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/buyer/purchase_orders/{purchase_order_id}/invoice [get]
func GetPurchaseOrderInvoice(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, "")
	}

	var params repo.GetPurchaseOrderInvoiceParams
	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, "")
	}

	params.JwtClaimsInfo = claims
	result, err := repo.NewPurchaseOrderRepo(cc.App.DB).GetPurchaseOrderInvoice(params)
	if err != nil {
		return eris.Wrap(err, "")
	}

	return cc.Success(result)
}

// BuyerUpdatePurchaseOrderLogs
// @Tags Admin-PurchaseOrder
// @Summary PurchaseOrder cancel
// @Description PurchaseOrder cancel
// @Accept  json
// @Produce  json
// @Success 200 {object} models.PurchaseOrder
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/buyer/purchase_orders/{purchase_order_id}/logs [put]
func BuyerUpdatePurchaseOrderLogs(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var form repo.UpdatePurchaseOrderLogsParams

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, "")
	}

	err = cc.BindAndValidate(&form)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	form.JwtClaimsInfo = claims
	result, err := repo.NewPurchaseOrderRepo(cc.App.DB).UpdatePurchaseOrderLogs(form)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	return cc.Success(result)
}

// BuyerCreatePurchaseOrderLogs
// @Tags Admin-PurchaseOrder
// @Summary PurchaseOrder cancel
// @Description PurchaseOrder cancel
// @Accept  json
// @Produce  json
// @Success 200 {object} models.PurchaseOrder
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/buyer/purchase_orders/{purchase_order_id}/logs [delete]
func BuyerDeletePurchaseOrderLogs(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var form repo.DeletePurchaseOrderLogsParams

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, "")
	}

	err = cc.BindAndValidate(&form)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	form.JwtClaimsInfo = claims
	result, err := repo.NewPurchaseOrderRepo(cc.App.DB).DeletePurchaseOrderLogs(form)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	return cc.Success(result)
}

// BuyerPurchaseOrderUpdateDesign
// @Tags Buyer-PO
// @Summary update purchase order
// @Description update purchase order
// @Accept  json
// @Produce  json
// @Param purchase_order_id query string true "ID"
// @Success 200 {object} models.PurchaseOrder
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/buyer/purchase_orders/{purchase_order_id}/update_design [put]
func BuyerPurchaseOrderUpdateDesign(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, "")
	}

	var params repo.BuyerUpdatePurchaseOrderDesignParams
	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, "")
	}

	params.JwtClaimsInfo = claims

	purchaseOrder, err := repo.NewPurchaseOrderRepo(cc.App.DB).BuyerUpdatePurchaseOrderDesign(params)
	if err != nil {
		return eris.Wrap(err, "")
	}

	for _, userID := range purchaseOrder.AssigneeIDs {
		tasks.CreateUserNotificationTask{
			UserID:           userID,
			Message:          fmt.Sprintf("New design was updated for %s", purchaseOrder.ReferenceID),
			NotificationType: enums.UserNotificationTypePoUpdateDesign,
			Metadata: &models.UserNotificationMetadata{
				AdminID:                  claims.GetUserID(),
				PurchaseOrderID:          purchaseOrder.ID,
				PurchaseOrderReferenceID: purchaseOrder.ReferenceID,
				InquiryID: func() string {
					if purchaseOrder.Inquiry != nil {
						return purchaseOrder.Inquiry.ID
					}
					return ""
				}(),
				InquiryReferenceID: func() string {
					if purchaseOrder.Inquiry != nil {
						return purchaseOrder.Inquiry.ReferenceID
					}
					return ""
				}(),
			},
		}.Dispatch(c.Request().Context())

		var eventData = purchaseOrder.GetCustomerIOMetadata(map[string]interface{}{
			"admin_po_url": fmt.Sprintf("%s/samples/%s/customer?open_design_comments=true", cc.App.Config.AdminPortalBaseURL, purchaseOrder.ID),
			"brand_po_url": fmt.Sprintf("%s/samples/%s?open_design_comments=true", cc.App.Config.BrandPortalBaseURL, purchaseOrder.ID),
		})

		if params.TechpackAttachments != nil {
			eventData["updated_techpack_attachments"] = params.TechpackAttachments.GenerateFileURL()
		}

		if params.Attachments != nil {
			eventData["updated_attachments"] = params.Attachments.GenerateFileURL()
		}
		tasks.TrackCustomerIOTask{
			UserID: userID,
			Event:  customerio.EventAdminPoUpdateDesign,
			Data:   eventData,
		}.Dispatch(c.Request().Context())
	}

	_, _ = tasks.UpdateUserProductClassesTask{
		UserID:          claims.GetUserID(),
		PurchaseOrderID: purchaseOrder.ID,
	}.Dispatch(c.Request().Context())

	return cc.Success(purchaseOrder)
}

// BuyerMakeAnotherPurchaseOrder
// @Tags Admin-Inquiry
// @Summary Inquiry confirm
// @Description Inquiry confirm
// @Accept  json
// @Produce  json
// @Success 200 {object} models.Inquiry
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/buyer/purchase_orders/{purchase_order_id}/make_another [post]
func BuyerMakeAnotherPurchaseOrder(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.BuyerMakeAnotherSampleParams

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, "")
	}

	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	params.JwtClaimsInfo = claims
	result, err := repo.NewPurchaseOrderRepo(cc.App.DB).BuyerMakeAnotherSample(params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	return cc.Success(result)
}
