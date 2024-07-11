package controllers

import (
	"fmt"
	"time"

	"github.com/hibiken/asynq"

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
// @Tags Admin-PO
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
// @Router /api/v1/admin/purchase_orders [post]
func CreatePurchaseOrder(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	var params repo.CreatePurchaseOrderParams
	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	params.JwtClaimsInfo = claims
	result, err := repo.NewPurchaseOrderRepo(cc.App.DB).CreatePurchaseOrder(params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	return cc.Success(result)
}

// PaginatePurchaseOrders
// @Tags Admin-PO
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
// @Router /api/v1/admin/purchase_orders [get]
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
	params.IncludeAssignee = true
	params.IncludeSampleMaker = true
	params.IncludeTrackings = true
	params.IncludeUnpaidWithoutInquiry = true
	params.IncludeItems = true
	params.IncludeCollection = true
	var result = repo.NewPurchaseOrderRepo(cc.App.DB).PaginatePurchaseOrders(params)

	return cc.Success(result)
}

// ExportPurchaseOrders
// @Tags Admin-PO
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
// @Router /api/v1/admin/purchase_orders/export [get]
func ExportPurchaseOrders(c echo.Context) error {
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
	result, err := repo.NewPurchaseOrderRepo(cc.App.DB).ExportExcel(params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}
	return cc.Success(result)
}

// GetPurchaseOrder
// @Tags Admin-PO
// @Summary Get purchase order
// @Description Get purchase order
// @Accept  json
// @Produce  json
// @Param purchase_order_id query string true "ID"
// @Success 200 {object} models.PurchaseOrder
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/purchase_orders/{purchase_order_id} [get]
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
	params.IncludeAssignee = true
	params.IncludeSampleMaker = true
	params.IncludeInquirySeller = true
	params.IncludeUsers = true
	params.IncludePaymentTransaction = true
	result, err := repo.NewPurchaseOrderRepo(cc.App.DB).GetPurchaseOrder(params)
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
// @Router /api/v1/admin/purchase_orders/{purchase_order_id}/update_tracking_status [put]
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

// AdminUpdatePurchaseOrderDesign
// @Tags Admin-PO
// @Summary update purchase order
// @Description update purchase order
// @Accept  json
// @Produce  json
// @Param purchase_order_id query string true "ID"
// @Success 200 {object} models.PurchaseOrder
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/purchase_orders/{purchase_order_id}/update_design [put]
func AdminUpdatePurchaseOrderDesign(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, "")
	}

	var params repo.AdminUpdatePurchaseOrderDesignParams
	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, "")
	}

	params.JwtClaimsInfo = claims

	purchaseOrder, err := repo.NewPurchaseOrderRepo(cc.App.DB).AdminUpdatePurchaseOrderDesign(params)
	if err != nil {
		return eris.Wrap(err, "")
	}

	tasks.CreateUserNotificationTask{
		UserID:           purchaseOrder.UserID,
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

	tasks.TrackCustomerIOTask{
		UserID: purchaseOrder.UserID,
		Event:  customerio.EventPoUpdateDesign,
		Data: purchaseOrder.GetCustomerIOMetadata(map[string]interface{}{
			"updated_techpack_attachments": params.TechpackAttachments.GenerateFileURL(),
			"admin_po_url":                 fmt.Sprintf("%s/samples/%s/customer?open_design_comments=true", cc.App.Config.AdminPortalBaseURL, purchaseOrder.ID),
			"brand_po_url":                 fmt.Sprintf("%s/samples/%s?open_design_comments=true", cc.App.Config.BrandPortalBaseURL, purchaseOrder.ID),
		}),
	}.Dispatch(c.Request().Context())

	_, _ = tasks.UpdateUserProductClassesTask{
		UserID:          purchaseOrder.UserID,
		PurchaseOrderID: purchaseOrder.ID,
	}.Dispatch(c.Request().Context())

	return cc.Success(purchaseOrder)
}

// AdminUpdatePurchaseOrderRawMaterial
// @Tags Admin-PO
// @Summary update purchase order
// @Description update purchase order
// @Accept  json
// @Produce  json
// @Param purchase_order_id query string true "ID"
// @Success 200 {object} models.PurchaseOrder
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/purchase_orders/{purchase_order_id}/update_raw_material [put]
func AdminUpdatePurchaseOrderRawMaterial(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, "")
	}

	var params repo.AdminUpdatePurchaseOrderRawMaterialParams
	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, "")
	}

	params.JwtClaimsInfo = claims

	purchaseOrder, err := repo.NewPurchaseOrderRepo(cc.App.DB).AdminUpdatePurchaseOrderRawMaterial(params)
	if err != nil {
		return eris.Wrap(err, "")
	}

	if params.ApproveRawMaterialAt != nil {
		_, _ = tasks.PurchaseOrderRawMaterialApproveTask{
			JwtClaimsInfo:   claims,
			PurchaseOrderID: params.PurchaseOrderID,
		}.DispatchAt(time.Unix(*params.ApproveRawMaterialAt, 0), asynq.MaxRetry(0))
	}

	tasks.CreateUserNotificationTask{
		UserID:           purchaseOrder.UserID,
		Message:          fmt.Sprintf("New raw material was updated for %s", purchaseOrder.ReferenceID),
		NotificationType: enums.UserNotificationTypePoUpdateRawMaterial,
		Metadata: &models.UserNotificationMetadata{
			AdminID:                  claims.GetUserID(),
			PurchaseOrderID:          purchaseOrder.ID,
			PurchaseOrderReferenceID: purchaseOrder.ReferenceID,
		},
	}.Dispatch(c.Request().Context())

	if len(params.PoRawMaterials) > 0 {
		tasks.TrackCustomerIOTask{
			UserID: purchaseOrder.UserID,
			Event:  customerio.EventPoUpdateMaterial,
			Data: purchaseOrder.GetCustomerIOMetadata(map[string]interface{}{
				"updated_po_raw_materials": params.PoRawMaterials.GenerateFileURL(),
			}),
		}.Dispatch(c.Request().Context())
	}

	return cc.Success(purchaseOrder)
}

// AdmiPurchaseOrderMarkMaking
// @Tags Admin-PO
// @Summary update purchase order
// @Description update purchase order
// @Accept  json
// @Produce  json
// @Param purchase_order_id query string true "ID"
// @Success 200 {object} models.PurchaseOrder
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/purchase_orders/{purchase_order_id}/mark_making [put]
func AdminPurchaseOrderMarkMaking(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, "")
	}

	var params repo.AdminPurchaseOrderMarkMakingParams
	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, "")
	}

	params.JwtClaimsInfo = claims

	purchaseOrder, err := repo.NewPurchaseOrderRepo(cc.App.DB).AdminPurchaseOrderMarkMaking(params)
	if err != nil {
		return eris.Wrap(err, "")
	}

	tasks.CreateUserNotificationTask{
		UserID:           purchaseOrder.UserID,
		Message:          fmt.Sprintf("Order %s is on making", purchaseOrder.ReferenceID),
		NotificationType: enums.UserNotificationTypePoMarkMaking,
		Metadata: &models.UserNotificationMetadata{
			AdminID:                  claims.GetUserID(),
			PurchaseOrderID:          purchaseOrder.ID,
			PurchaseOrderReferenceID: purchaseOrder.ReferenceID,
		},
	}.Dispatch(c.Request().Context())

	tasks.TrackCustomerIOTask{
		UserID: purchaseOrder.UserID,
		Event:  customerio.EventPoMarkMaking,
		Data:   purchaseOrder.GetCustomerIOMetadata(nil),
	}.Dispatch(c.Request().Context())

	return cc.Success(purchaseOrder)
}

// AdminPurchaseOrderMarkMakingWithoutRawMaterial
// @Tags Admin-PO
// @Summary update purchase order
// @Description update purchase order
// @Accept  json
// @Produce  json
// @Param purchase_order_id query string true "ID"
// @Success 200 {object} models.PurchaseOrder
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/purchase_orders/{purchase_order_id}/mark_making_without_raw_material [put]
func AdminPurchaseOrderMarkMakingWithoutRawMaterial(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, "")
	}

	var params repo.AdminPurchaseOrderMarkMakingWithoutRawMaterialParams
	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, "")
	}

	params.JwtClaimsInfo = claims

	purchaseOrder, err := repo.NewPurchaseOrderRepo(cc.App.DB).MarkMakingWithoutRawMaterial(params)
	if err != nil {
		return eris.Wrap(err, "")
	}

	tasks.CreateUserNotificationTask{
		UserID:           purchaseOrder.UserID,
		Message:          fmt.Sprintf("Order %s is on making", purchaseOrder.ReferenceID),
		NotificationType: enums.UserNotificationTypePoMarkMaking,
		Metadata: &models.UserNotificationMetadata{
			AdminID:                  claims.GetUserID(),
			PurchaseOrderID:          purchaseOrder.ID,
			PurchaseOrderReferenceID: purchaseOrder.ReferenceID,
		},
	}.Dispatch(c.Request().Context())

	tasks.TrackCustomerIOTask{
		UserID: purchaseOrder.UserID,
		Event:  customerio.EventPoMarkMaking,
		Data:   purchaseOrder.GetCustomerIOMetadata(nil),
	}.Dispatch(c.Request().Context())

	return cc.Success(purchaseOrder)
}

// AdminPurchaseOrderMarkSubmit
// @Tags Admin-PO
// @Summary update purchase order
// @Description update purchase order
// @Accept  json
// @Produce  json
// @Param purchase_order_id query string true "ID"
// @Success 200 {object} models.PurchaseOrder
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/purchase_orders/{purchase_order_id}/mark_submit [put]
func AdminPurchaseOrderMarkSubmit(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, "")
	}

	var params repo.AdminPurchaseOrderMarkSubmitParams
	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, "")
	}

	params.JwtClaimsInfo = claims

	purchaseOrder, err := repo.NewPurchaseOrderRepo(cc.App.DB).AdminPurchaseOrderMarkSubmit(params)
	if err != nil {
		return eris.Wrap(err, "")
	}

	tasks.CreateUserNotificationTask{
		UserID:           purchaseOrder.UserID,
		Message:          fmt.Sprintf("Order %s is on submit", purchaseOrder.ReferenceID),
		NotificationType: enums.UserNotificationTypePoMarkSubmit,
		Metadata: &models.UserNotificationMetadata{
			AdminID:                  claims.GetUserID(),
			PurchaseOrderID:          purchaseOrder.ID,
			PurchaseOrderReferenceID: purchaseOrder.ReferenceID,
		},
	}.Dispatch(c.Request().Context())

	tasks.TrackCustomerIOTask{
		UserID: purchaseOrder.UserID,
		Event:  customerio.EventPoMarkSubmit,
		Data:   purchaseOrder.GetCustomerIOMetadata(nil),
	}.Dispatch(c.Request().Context())

	return cc.Success(purchaseOrder)
}

// AdmiPurchaseOrderMarkDelivering
// @Tags Admin-PO
// @Summary update purchase order
// @Description update purchase order
// @Accept  json
// @Produce  json
// @Param purchase_order_id query string true "ID"
// @Success 200 {object} models.PurchaseOrder
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/purchase_orders/{purchase_order_id}/mark_delivering [put]
func AdmiPurchaseOrderMarkDelivering(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, "")
	}

	var params repo.AdminPurchaseOrderMarkDeliveringParams
	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, "")
	}

	params.JwtClaimsInfo = claims

	purchaseOrder, err := repo.NewPurchaseOrderRepo(cc.App.DB).AdminPurchaseOrderMarkDelivering(params)
	if err != nil {
		return eris.Wrap(err, "")
	}

	tasks.CreateUserNotificationTask{
		UserID:           purchaseOrder.UserID,
		Message:          fmt.Sprintf("Order %s is on delivering", purchaseOrder.ReferenceID),
		NotificationType: enums.UserNotificationTypePoMarkDelivering,
		Metadata: &models.UserNotificationMetadata{
			AdminID:                  claims.GetUserID(),
			PurchaseOrderID:          purchaseOrder.ID,
			PurchaseOrderReferenceID: purchaseOrder.ReferenceID,
		},
	}.Dispatch(c.Request().Context())

	tasks.TrackCustomerIOTask{
		UserID: purchaseOrder.UserID,
		Event:  customerio.EventPoMarkDelivering,
		Data:   purchaseOrder.GetCustomerIOMetadata(nil),
	}.Dispatch(c.Request().Context())

	return cc.Success(purchaseOrder)
}

// AdmiPurchaseOrderMarkDelivered
// @Tags Admin-PO
// @Summary update purchase order
// @Description update purchase order
// @Accept  json
// @Produce  json
// @Param purchase_order_id query string true "ID"
// @Success 200 {object} models.PurchaseOrder
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/purchase_orders/{purchase_order_id}/mark_delivered [put]
func AdmiPurchaseOrderMarkDelivered(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.UpdatePurchaseOrderTrackingStatusParams
	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, "")
	}

	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, "")
	}

	params.JwtClaimsInfo = claims
	params.TrackingAction = enums.PoTrackingActionConfirmDelivered
	params.TrackingStatus = enums.PoTrackingStatusDeliveryConfirmed
	purchaseOrder, err := repo.NewPurchaseOrderRepo(cc.App.DB).AdminPurchaseOrderConfirmDelivered(params)
	if err != nil {
		return eris.Wrap(err, "")
	}

	tasks.TrackCustomerIOTask{
		UserID: purchaseOrder.UserID,
		Event:  customerio.EventPoConfirmDelivered,
		Data:   purchaseOrder.GetCustomerIOMetadata(nil),
	}.Dispatch(c.Request().Context())

	return cc.Success(purchaseOrder)
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
// @Router /api/v1/admin/purchase_orders/{purchase_order_id}/logs [get]
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
	var result = repo.NewPurchaseOrderTrackingRepo(cc.App.DB).PaginatePurchaseOrderTrackings(params)

	return cc.Success(result)
}

// AdminApprovePurchaseOrderDesign
// @Tags Marketplace-PurchaseOrder
// @Summary Approve PurchaseOrder quotation
// @Description Approve PurchaseOrder quotation
// @Accept  json
// @Produce  json
// @Param data body repo.BuyerApproveDesignParams true "Form"
// @Success 200 {object} models.PurchaseOrder
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/purchase_orders/{purchase_order_id}/approve_design [put]
func AdminApprovePurchaseOrderDesign(c echo.Context) error {
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

	tasks.TrackCustomerIOTask{
		UserID: purchaseOrder.UserID,
		Event:  customerio.EventPoApproveDesign,
		Data:   purchaseOrder.GetCustomerIOMetadata(nil),
	}.Dispatch(c.Request().Context())

	return cc.Success("Approved")
}

// AdminPurchaseOrderAddDesignComments
// @Tags Admin-PO
// @Summary Add design comments
// @Description Add design comments
// @Accept  json
// @Produce  json
// @Param data body models.Comment true "Form"
// @Success 200 {object} models.PurchaseOrder
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/purchase_orders/{purchase_order_id}/design_comments [post]
func AdminPurchaseOrderAddDesignComments(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params models.CommentCreateForm
	var orderID = cc.GetPathParamString("purchase_order_id")

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, "")
	}
	params.TargetType = enums.CommentTargetTypePurchaseOrderDesign
	params.TargetID = orderID
	params.JwtClaimsInfo = claims

	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, "")
	}

	comment, err := repo.NewCommentRepo(cc.App.DB).CreateComment(params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	if purchaseOrder, err := repo.NewPurchaseOrderRepo(cc.App.DB).GetPurchaseOrderShortInfo(orderID); err == nil {
		comment.PurchaseOrder = purchaseOrder
		tasks.CreateUserNotificationTask{
			UserID:           purchaseOrder.UserID,
			Message:          fmt.Sprintf("New comment on sample purchase order %s", purchaseOrder.ReferenceID),
			NotificationType: enums.UserNotificationTypePoDesignNewComment,
			Metadata: &models.UserNotificationMetadata{
				CommentID:                comment.ID,
				PurchaseOrderID:          purchaseOrder.ID,
				PurchaseOrderReferenceID: purchaseOrder.ReferenceID,
			},
		}.Dispatch(c.Request().Context())

		if userInfo, err := repo.NewUserRepo(cc.App.DB).GetShortUserInfo(claims.GetUserID()); err == nil {
			tasks.TrackCustomerIOTask{
				UserID: purchaseOrder.UserID,
				Event:  customerio.EventAdminNewDesignComment,
				Data: comment.GetCustomerIOMetadata(map[string]interface{}{
					"sender":       userInfo.GetCustomerIOMetadata(nil),
					"admin_po_url": fmt.Sprintf("%s/samples/%s/customer?open_design_comments=true", cc.App.Config.AdminPortalBaseURL, purchaseOrder.ID),
					"brand_po_url": fmt.Sprintf("%s/samples/%s?open_design_comments=true", cc.App.Config.BrandPortalBaseURL, purchaseOrder.ID),
				}),
			}.Dispatch(c.Request().Context())
		}

	}

	return cc.Success(comment)
}

// AdminPurchaseOrderPaginateDesignComments
// @Tags Admin-PO
// @Summary Paginate design comments
// @Description Paginate design comments
// @Accept  json
// @Produce  json
// @Param data body models.Comment true "Form"
// @Success 200 {object} models.PurchaseOrder
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/purchase_orders/{purchase_order_id}/design_comments [get]
func AdminPurchaseOrderPaginateDesignComments(c echo.Context) error {
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

// AdminPurchaseOrderDesignCommentMarkSeen
// @Tags Admin-PO
// @Summary Mark seen comments
// @Description Mark seen comments
// @Accept  json
// @Produce  json
// @Param data body models.Comment true "Form"
// @Success 200 {object} models.Inquiry
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/purchase_orders/{purchase_order_id}/design_comments/mark_seen [put]
func AdminPurchaseOrderDesignCommentMarkSeen(c echo.Context) error {
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

// AdminPurchaseOrderDesignCommentMarkSeen
// @Tags Admin-PO
// @Summary Mark seen comments
// @Description Mark seen comments
// @Accept  json
// @Produce  json
// @Param data body models.Comment true "Form"
// @Success 200 {object} models.Inquiry
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/purchase_orders/{purchase_order_id}/design_comments/status_count [get]
func AdminPurchaseOrderDesignCommentStatusCount(c echo.Context) error {
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

// AdminPurchaseOrderAssignPIC
// @Tags Admin-PO
// @Summary Assign PIC
// @Description Assign PIC
// @Accept  json
// @Produce  json
// @Success 200 {object} models.User
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/purchase_orders/{purchase_order_id}/assign_pic [put]
func AdminPurchaseOrderAssignPIC(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params models.PurchaseOrderAssignPICParam

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	params.JwtClaimsInfo = claims
	result, err := repo.NewPurchaseOrderRepo(cc.App.DB).PurchaseOrderAssignPIC(params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	for _, userID := range result.AssigneeIDs {
		tasks.AssignPurchaseOrderPICTask{
			AssignerID:      claims.GetUserID(),
			AssigneeID:      userID,
			PurchaseOrderID: params.PurchaseOrderID,
		}.Dispatch(c.Request().Context())
	}
	return cc.Success(result)
}

// AdminArchivePurchaseOrder Admin archive PurchaseOrder
// @Tags Admin-PO
// @Summary Admin archive PurchaseOrder
// @Description Admin archive PurchaseOrder
// @Accept  json
// @Produce  json
// @Success 200 {object} models.PurchaseOrder
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/purchase_orders/{purchase_order_id}/archive [delete]
func AdminArchivePurchaseOrder(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.AdminArchivePurchaseOrderParams

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, "")
	}

	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	params.JwtClaimsInfo = claims
	err = repo.NewPurchaseOrderRepo(cc.App.DB).AdminArchivePurchaseOrder(params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	return cc.Success("Archived")
}

// AdminUnarchivePurchaseOrder Admin unarchive PurchaseOrder
// @Tags Admin-PO
// @Summary Admin unarchive PurchaseOrder
// @Description Admin unarchive PurchaseOrder
// @Accept  json
// @Produce  json
// @Success 200 {object} models.PurchaseOrder
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/purchase_orders/{purchase_order_id}/unarchive [put]
func AdminUnarchivePurchaseOrder(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.AdminUnarchivePurchaseOrderParams

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, "")
	}

	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, "")
	}

	params.JwtClaimsInfo = claims
	err = repo.NewPurchaseOrderRepo(cc.App.DB).AdminUnarchivePurchaseOrder(params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	return cc.Success("Unarchived")
}

// AdminDeletePurchaseOrder Admin archive PurchaseOrder
// @Tags Admin-PO
// @Summary Admin archive PurchaseOrder
// @Description Admin archive PurchaseOrder
// @Accept  json
// @Produce  json
// @Success 200 {object} models.PurchaseOrder
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/purchase_orders/{purchase_order_id}/delete [delete]
func AdminDeletePurchaseOrder(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.AdminDeletePurchaseOrderParams

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, "")
	}

	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	params.JwtClaimsInfo = claims
	err = repo.NewPurchaseOrderRepo(cc.App.DB).AdminDeletePurchaseOrder(params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	return cc.Success("Deleted")
}

// AdminPurchaseOrderApproveRawMaterial
// @Tags Admin-PO
// @Summary update purchase order
// @Description update purchase order
// @Accept  json
// @Produce  json
// @Param purchase_order_id query string true "ID"
// @Success 200 {object} models.PurchaseOrder
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/purchase_orders/{purchase_order_id}/approve_raw_material [post]
func AdminPurchaseOrderApproveRawMaterial(c echo.Context) error {
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
			Event:  customerio.EventPoApproveRawMaterial,
			Data:   data,
		}.Dispatch(c.Request().Context())
	}

	return cc.Success(order)
}

// PurchaseOrderAddNote
// @Tags Admin-PurchaseOrder
// @Summary Add design comments
// @Description Add design comments
// @Accept  json
// @Produce  json
// @Param data body models.Comment true "Form"
// @Success 200 {object} models.PurchaseOrder
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/purchase_orders/{purchase_order_id}/notes [post]
func PurchaseOrderAddNote(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params models.CommentCreateForm

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, "")
	}
	params.TargetType = enums.CommentTargetTypePOInternalNotes
	params.TargetID = cc.GetPathParamString("purchase_order_id")
	params.JwtClaimsInfo = claims

	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, "")
	}

	result, err := repo.NewCommentRepo(cc.App.DB).CreateComment(params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	tasks.NewPONotesTask{
		UserID:          claims.GetUserID(),
		PurchaseOrderID: params.TargetID,
		MentionUserIDs:  params.MentionUserIDs,
		Message:         params.Message,
		Attachments:     params.Attachments,
	}.Dispatch(c.Request().Context())

	return cc.Success(result)
}

// PaginatePurchaseOrderNotes
// @Tags Admin-PurchaseOrder
// @Summary Get design comments
// @Description Get design comments
// @Accept  json
// @Produce  json
// @Param data body models.Comment true "Form"
// @Success 200 {object} models.PurchaseOrder
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/purchase_orders/{purchase_order_id}/notes [get]
func PaginatePurchaseOrderNotes(c echo.Context) error {
	var cc = c.(*models.CustomContext)
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
	params.TargetID = cc.GetPathParamString("purchase_order_id")
	params.TargetType = enums.CommentTargetTypePOInternalNotes
	params.OrderByQuery = "c.created_at DESC"

	var results = repo.NewCommentRepo(cc.App.DB).PaginateComment(params)

	return cc.Success(results)
}

// PurchaseOrderNoteMarkSeen
// @Tags Marketplace-PO
// @Summary Add design comments
// @Description Add design comments
// @Accept  json
// @Produce  json
// @Param data body models.Comment true "Form"
// @Success 200 {object} models.PurchaseOrder
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/purchase_orders/{purchase_order_id}/notes/mark_seen [put]
func PurchaseOrderNoteMarkSeen(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.MarkSeenParams

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, "")
	}
	params.JwtClaimsInfo = claims
	params.TargetID = cc.GetPathParamString("purchase_order_id")
	params.TargetType = enums.CommentTargetTypePOInternalNotes

	err = repo.NewCommentRepo(cc.App.DB).MarkSeen(params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	return cc.Success("Seen")
}

// PurchaseOrderNoteStatusCount
// @Tags Marketplace-PO
// @Summary Design comment status count
// @Description Mark seen comments
// @Accept  json
// @Produce  json
// @Param data body models.Comment true "Form"
// @Success 200 {object} models.PurchaseOrder
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/purchase_orders/{purchase_order_id}/notes/unread_count [get]
func PurchaseOrderNoteUnreadCount(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.GetUnreadCountParams

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, "")
	}
	params.JwtClaimsInfo = claims
	params.TargetID = cc.GetPathParamString("purchase_order_id")
	params.TargetType = enums.CommentTargetTypePOInternalNotes

	var result = repo.NewCommentRepo(cc.App.DB).GetUnreadCount(params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	return cc.Success(result)
}

// AdminPurchaseOrderCommentDelete
// @Tags Admin-PurchaseOrder
// @Summary delete inquiry comment
// @Description delete inquiry comment
// @Accept  json
// @Produce  json
// @Param data body models.ContentCommentCreateForm true "Form"
// @Success 200 {object} models.Comment
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/purchase_orders/{purchase_order_id}/notes/{comment_id} [delete]
func AdminPurchaseOrderCommentDelete(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.DeleteCommentParams

	var err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	err = repo.NewCommentRepo(cc.App.DB).DeleteComment(params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	return cc.Success("Deleted")
}

// AdminPurchaseOrderRefund
// @Tags Admin-PurchaseOrder
// @Summary delete inquiry comment
// @Description delete inquiry comment
// @Accept  json
// @Produce  json
// @Param data body models.ContentCommentCreateForm true "Form"
// @Success 200 {object} models.Comment
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/purchase_orders/{purchase_order_id}/cancel [delete]
func AdminPurchaseOrderRefund(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.RefundPurchaseOrderParams

	var err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	result, err := repo.NewPurchaseOrderRepo(cc.App.DB).RefundPurchaseOrder(params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	return cc.Success(result)
}

// AdmiPurchaseOrderStageCommentsCreate
// @Tags Admin-PO
// @Summary create stage comment
// @Description create stage comment
// @Accept  json
// @Produce  json
// @Param purchase_order_id query string true "ID"
// @Success 200 {object} models.PurchaseOrder
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/purchase_orders/{purchase_order_id}/stage_comments [post]
func AdmiPurchaseOrderStageCommentsCreate(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, "")
	}

	var params repo.PurchaseOrderStageCommentsParams
	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, "")
	}

	params.JwtClaimsInfo = claims

	purchaseOrder, err := repo.NewPurchaseOrderRepo(cc.App.DB).StageCommentsCreate(params)
	if err != nil {
		return eris.Wrap(err, "")
	}

	tasks.TrackCustomerIOTask{
		UserID: purchaseOrder.UserID,
		Event:  customerio.EventPoNewComment,
		Data: purchaseOrder.GetCustomerIOMetadata(map[string]interface{}{
			"comment": params.Comment,
			"comment_attachments": func() models.Attachments {
				return params.Attachments.GenerateFileURL()
			}(),
		}),
	}.Dispatch(c.Request().Context())

	return cc.Success(purchaseOrder)
}

// AdminUpdatePurchaseOrder
// @Tags Admin-PO
// @Summary update purchase order
// @Description update purchase order
// @Accept  json
// @Produce  json
// @Param purchase_order_id query string true "ID"
// @Success 200 {object} models.PurchaseOrder
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/purchase_orders/{purchase_order_id} [put]
func AdminUpdatePurchaseOrder(c echo.Context) error {
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

// AdminCreateBulkPurchaseOrderFromSample
// @Tags Admin-PO
// @Summary Create bulk from sample
// @Description Create bulk from sample
// @Accept  json
// @Produce  json
// @Param data body repo.InquiryCheckoutParams true "Form"
// @Success 200 {object} models.PurchaseOrder
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/purchase_orders/{purchase_order_id}/bulk_purchase_order [post]
func AdminCreateBulkPurchaseOrderFromSample(c echo.Context) error {
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

	if len(result.AssigneeIDs) > 0 {
		for _, userID := range result.AssigneeIDs {
			tasks.TrackCustomerIOTask{
				UserID: userID,
				Event:  customerio.EventBulkPoCreated,
				Data:   result.GetCustomerIOMetadata(nil),
			}.Dispatch(c.Request().Context())
		}
	} else {
		tasks.TrackCustomerIOTask{
			UserID: cc.App.Config.InflowMerchandiseGroupEmail,
			Event:  customerio.EventBulkPoCreated,
			Data:   result.GetCustomerIOMetadata(nil),
		}.Dispatch(c.Request().Context())
	}

	return cc.Success(result)
}

// AdminCreatePurchaseOrder
// @Tags Admin-PO
// @Summary update purchase order
// @Description update purchase order
// @Accept  json
// @Produce  json
// @Param purchase_order_id query string true "ID"
// @Success 200 {object} models.PurchaseOrder
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/purchase_orders/{purchase_order_id} [put]
func AdminCreatePurchaseOrder(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.CreatePurchaseOrderParams

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, "")
	}

	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, "")
	}

	params.JwtClaimsInfo = claims
	order, err := repo.NewPurchaseOrderRepo(cc.App.DB).CreatePurchaseOrder(params)
	if err != nil {
		return eris.Wrap(err, "")
	}

	for _, assigneeID := range order.AssigneeIDs {
		var data = order.GetCustomerIOMetadata(nil)
		tasks.TrackCustomerIOTask{
			UserID: assigneeID,
			Event:  customerio.EventPoCreated,
			Data:   data,
		}.Dispatch(c.Request().Context())
	}

	return cc.Success(order)
}

// AdminPurchaseOrderCreatePaymentLink Admin create payment link
// @Tags Admin-PurchaseOrder
// @Summary Admin create payment link
// @Description Admin create payment link
// @Accept  json
// @Produce  json
// @Success 200 {object} models.PurchaseOrder
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/purchase_orders/payment_link [post]
func AdminMultiPurchaseOrderCreatePaymentLink(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.CreateMultiPurchaseOrderPaymentLinkParams

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, "")
	}

	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	params.JwtClaimsInfo = claims
	resp, err := repo.NewPurchaseOrderRepo(cc.App.DB).MultiPurchaseOrderCreatePaymentLink(params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	return cc.Success(resp)
}

// AdminPurchaseOrderConfirm Admin confirm PO
// @Tags Admin-PurchaseOrder
// @Summary Admin confirm PO
// @Description Admin confirm PO
// @Accept  json
// @Produce  json
// @Success 200 {object} models.PurchaseOrder
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/purchase_orders/{purchase_order_id}/confirm [put]
func AdminPurchaseOrderConfirm(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.AdminConfirmPOParams

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, "")
	}

	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	params.JwtClaimsInfo = claims
	resp, err := repo.NewPurchaseOrderRepo(cc.App.DB).AdminConfirmPO(params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	userAdmin, err := repo.NewUserRepo(cc.App.DB).GetShortUserInfo(params.GetUserID())
	if err != nil {
		return err
	}

	tasks.TrackCustomerIOTask{
		UserID: resp.UserID,
		Event:  customerio.EventPoConfirmed,
		Data: resp.GetCustomerIOMetadata(map[string]interface{}{
			"confirmed_by": userAdmin.GetCustomerIOMetadata(nil),
		}),
	}.Dispatch(c.Request().Context())

	return cc.Success(resp)
}

// AdminPurchaseOrderConfirm Admin confirm PO
// @Tags Admin-PurchaseOrder
// @Summary Admin confirm PO
// @Description Admin confirm PO
// @Accept  json
// @Produce  json
// @Success 200 {object} models.PurchaseOrder
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/purchase_orders/{purchase_order_id}/cancel [put]
func AdminPurchaseOrderCancel(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.AdminCancelPOParams

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, "")
	}

	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	params.JwtClaimsInfo = claims
	resp, err := repo.NewPurchaseOrderRepo(cc.App.DB).AdminCancelPO(params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	userAdmin, err := repo.NewUserRepo(cc.App.DB).GetShortUserInfo(params.GetUserID())
	if err != nil {
		return err
	}

	tasks.TrackCustomerIOTask{
		UserID: resp.UserID,
		Event:  customerio.EventAdminPoCanceled,
		Data: resp.GetCustomerIOMetadata(map[string]interface{}{
			"canceled_by": userAdmin.GetCustomerIOMetadata(nil),
		}),
	}.Dispatch(c.Request().Context())

	return cc.Success(resp)
}

// AdminPurchaseOrderMarkAsPaid
// @Tags Admin-PurchaseOrder
// @Summary Mark as paid
// @Description Mark as paid
// @Accept  json
// @Produce  json
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/purchase_orders/{purchase_order_id}/mark_as_paid [post]
func AdminPurchaseOrderMarkAsPaid(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params models.PurchaseOrderIDParam

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	params.JwtClaimsInfo = claims

	// Handle multi inquiry same payment
	order, err := repo.NewPurchaseOrderRepo(cc.App.DB).GetPurchaseOrderShortInfo(params.PurchaseOrderID)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	if order.CheckoutSessionID != "" {
		purchaseOrders, err := repo.NewPurchaseOrderRepo(cc.App.DB).MultiPurchaseOrderMarkAsPaid(repo.MultiPurchaseOrderParams{
			CheckoutSessionID: order.CheckoutSessionID,
			Note:              params.Note,
		})
		if err != nil {
			return eris.Wrap(err, err.Error())
		}
		for _, purchaseOrder := range purchaseOrders {
			_, _ = tasks.PurchaseOrderBankTransferConfirmedTask{
				ApprovedByUserID: claims.GetUserID(),
				PurchaseOrderID:  purchaseOrder.ID,
			}.Dispatch(c.Request().Context())
		}

	} else {
		purchaseOrder, err := repo.NewPurchaseOrderRepo(cc.App.DB).PurchaseOrderMarkAsPaid(params)
		if err != nil {
			return eris.Wrap(err, err.Error())
		}
		_, _ = tasks.PurchaseOrderBankTransferConfirmedTask{
			ApprovedByUserID: claims.GetUserID(),
			PurchaseOrderID:  purchaseOrder.ID,
		}.Dispatch(c.Request().Context())
	}

	return cc.Success("Confirmed")
}

// AdminPurchaseOrderMarkAsUnpaid
// @Tags Admin-PurchaseOrder
// @Summary Mark as paid
// @Description Mark as paid
// @Accept  json
// @Produce  json
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/purchase_orders/{purchase_order_id}/mark_as_unpaid [post]
func AdminPurchaseOrderMarkAsUnpaid(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params models.PurchaseOrderIDParam

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	params.JwtClaimsInfo = claims

	// Handle multi inquiry same payment
	order, err := repo.NewPurchaseOrderRepo(cc.App.DB).GetPurchaseOrderShortInfo(params.PurchaseOrderID)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	if order.CheckoutSessionID != "" {
		purchaseOrders, err := repo.NewPurchaseOrderRepo(cc.App.DB).MultiPurchaseOrderMarkAsUnpaid(repo.MultiPurchaseOrderParams{
			CheckoutSessionID: order.CheckoutSessionID,
			Note:              params.Note,
		})
		if err != nil {
			return eris.Wrap(err, err.Error())
		}
		for _, purchaseOrder := range purchaseOrders {
			_, _ = tasks.PurchaseOrderBankTransferConfirmedTask{
				ApprovedByUserID: claims.GetUserID(),
				PurchaseOrderID:  purchaseOrder.ID,
			}.Dispatch(c.Request().Context())
		}

	} else {
		purchaseOrder, err := repo.NewPurchaseOrderRepo(cc.App.DB).PurchaseOrderMarkAsUnpaid(params)
		if err != nil {
			return eris.Wrap(err, err.Error())
		}
		_, _ = tasks.PurchaseOrderBankTransferConfirmedTask{
			ApprovedByUserID: claims.GetUserID(),
			PurchaseOrderID:  purchaseOrder.ID,
		}.Dispatch(c.Request().Context())
	}

	return cc.Success("Confirmed")
}

// AdminPurchaseOrderSkipDesign
// @Tags Admin-PurchaseOrder
// @Summary Mark as paid
// @Description Mark as paid
// @Accept  json
// @Produce  json
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/purchase_orders/{purchase_order_id}/skip_design [put]
func AdminPurchaseOrderSkipDesign(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.AdminSkipDesignParams

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	params.JwtClaimsInfo = claims

	order, err := repo.NewPurchaseOrderRepo(cc.App.DB).AdminSkipDesign(params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	return cc.Success(order)
}

// AdminPurchaseOrderApproveRound
// @Tags Admin-PurchaseOrder
// @Summary Mark as paid
// @Description Mark as paid
// @Accept  json
// @Produce  json
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/purchase_orders/{purchase_order_id}/rounds/{round_id}/approve [put]
func AdminPurchaseOrderApproveRound(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.AdminPurchaseOrderApproveRoundParams

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	params.JwtClaimsInfo = claims
	err = repo.NewPurchaseOrderRepo(cc.App.DB).AdminPurchaseOrderApproveRound(params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	return cc.Success("Approved")
}

// AdminPurchaseOrderRejectRound
// @Tags Admin-PurchaseOrder
// @Summary Mark as paid
// @Description Mark as paid
// @Accept  json
// @Produce  json
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/purchase_orders/{purchase_order_id}/rounds/{round_id}/reject [delete]
func AdminPurchaseOrderRejectRound(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.AdminPurchaseOrderRejectRoundParams

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	params.JwtClaimsInfo = claims
	err = repo.NewPurchaseOrderRepo(cc.App.DB).AdminPurchaseOrderRejectRound(params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	return cc.Success("Rejected")
}
