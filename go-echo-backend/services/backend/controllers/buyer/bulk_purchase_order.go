package controllers

import (
	"github.com/engineeringinflow/inflow-backend/pkg/customerio"
	"github.com/engineeringinflow/inflow-backend/pkg/errs"
	"github.com/engineeringinflow/inflow-backend/pkg/models"
	"github.com/engineeringinflow/inflow-backend/pkg/models/enums"
	"github.com/engineeringinflow/inflow-backend/pkg/repo"
	"github.com/engineeringinflow/inflow-backend/services/consumer/tasks"
	"github.com/labstack/echo/v4"
	"github.com/rotisserie/eris"
	"github.com/samber/lo"
)

// PaginateBulkPurchaseOrder
// @Tags Buyer-PO
// @Summary PaginateBulkPurchaseOrder
// @Description PaginateBulkPurchaseOrder
// @Accept  json
// @Produce  json
// @Success 200 {object} models.Records{records=[]models.BulkPurchaseOrder}
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/bulk_purchase_orders [get]
func PaginateBulkPurchaseOrder(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.PaginateBulkPurchaseOrderParams

	err := cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, "")
	}

	params.JwtClaimsInfo = claims
	params.IncludeTrackings = true
	params.IncludeCollection = true
	var result = repo.NewBulkPurchaseOrderRepo(cc.App.DB).PaginateBulkPurchaseOrder(params)

	return cc.Success(result)
}

// UpdateBulkPurchaseOrder
// @Tags Buyer-PO
// @Summary Get bulk purchase order
// @Description Get bulk purchase order
// @Accept  json
// @Produce  json
// @Param purchase_order_id query string true "ID"
// @Success 200 {object} models.PurchaseOrder
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/bulk_purchase_orders/{bulk_purchase_order_id} [put]
func UpdateBulkPurchaseOrder(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var form models.BulkPurchaseOrderUpdateForm

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, "")
	}

	err = cc.BindAndValidate(&form)
	if err != nil {
		return eris.Wrap(err, "")
	}

	form.UserID = claims.GetUserID()
	form.JwtClaimsInfo = claims
	result, err := repo.NewBulkPurchaseOrderRepo(cc.App.DB).UpdateBulkPurchaseOrder(form)
	if err != nil {
		return eris.Wrap(err, "")
	}

	if result.TrackingStatus == enums.BulkPoTrackingStatusWaitingForQuotation {
		if len(result.AssigneeIDs) > 0 {
			for _, userID := range result.AssigneeIDs {
				tasks.TrackCustomerIOTask{
					UserID: userID,
					Event:  customerio.EventBulkPoBuyerWaitingForQuotation,
					Data:   result.GetCustomerIOMetadata(nil),
				}.Dispatch(c.Request().Context())
			}

		} else {
			tasks.TrackCustomerIOTask{
				UserID: cc.App.Config.InflowMerchandiseGroupEmail,
				Event:  customerio.EventBulkPoBuyerWaitingForQuotation,
				Data:   result.GetCustomerIOMetadata(nil),
			}.Dispatch(c.Request().Context())
		}

	}

	return cc.Success(result)
}

// SubmitBulkPurchaseOrder
// @Tags Buyer-PO
// @Summary Get bulk purchase order
// @Description Get bulk purchase order
// @Accept  json
// @Produce  json
// @Param purchase_order_id query string true "ID"
// @Success 200 {object} models.PurchaseOrder
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/bulk_purchase_orders/{bulk_purchase_order_id}/submit [post]
func SubmitBulkPurchaseOrder(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var form models.BulkPurchaseOrderUpdateForm

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, "")
	}

	err = cc.BindAndValidate(&form)
	if err != nil {
		return eris.Wrap(err, "")
	}

	form.JwtClaimsInfo = claims
	form.UserID = claims.GetUserID()
	order, err := repo.NewBulkPurchaseOrderRepo(cc.App.DB).SubmitBulkPurchaseOrder(form)
	if err != nil {
		return eris.Wrap(err, "")
	}

	if len(order.AssigneeIDs) > 0 {
		for _, assigneeID := range order.AssigneeIDs {
			tasks.TrackCustomerIOTask{
				UserID: assigneeID,
				Event:  customerio.EventBulkPoBuyerSubmitOrder,
				Data:   order.GetCustomerIOMetadata(nil),
			}.Dispatch(c.Request().Context())
		}
	} else {
		tasks.TrackCustomerIOTask{
			UserID: cc.App.Config.InflowMerchandiseGroupEmail,
			Event:  customerio.EventBulkPoBuyerSubmitOrder,
			Data:   order.GetCustomerIOMetadata(nil),
		}.Dispatch(c.Request().Context())
	}

	return cc.Success(order)
}

// GetBulkPurchaseOrder
// @Tags Buyer-PO
// @Summary Get bulk purchase order
// @Description Get bulk purchase order
// @Accept  json
// @Produce  json
// @Param purchase_order_id query string true "ID"
// @Success 200 {object} models.PurchaseOrder
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/bulk_purchase_orders/{bulk_purchase_order_id} [put]
func GetBulkPurchaseOrder(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.GetBulkPurchaseOrderParams

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, "")
	}

	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, "")
	}

	params.JwtClaimsInfo = claims
	params.IncludeShippingAddress = true
	params.IncludePaymentTransactions = true
	params.IncludeInvoice = true
	params.IncludeItems = true
	params.UserID = claims.GetUserID()
	result, err := repo.NewBulkPurchaseOrderRepo(cc.App.DB).GetBulkPurchaseOrder(params)
	if err != nil {
		return eris.Wrap(err, "")
	}

	return cc.Success(result)
}

// BulkPurchaseOrdersPreviewCheckout
// @Tags Buyer-PO
// @Summary Preview checkout bulk purchase order
// @Description Preview checkout bulk purchase order
// @Accept  json
// @Produce  json
// @Param data body repo.BulkPurchaseOrderPreviewCheckoutParams true "Form"
// @Success 200 {object} models.PurchaseOrder
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/bulk_purchase_orders/preview_checkout [post]
func BulkPurchaseOrdersPreviewCheckout(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.BulkPurchaseOrdersPreviewCheckoutParams

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, "")
	}

	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	params.JwtClaimsInfo = claims
	params.UpdatePricing = true
	result, err := repo.NewBulkPurchaseOrderRepo(cc.App.DB).BulkPurchaseOrdersPreviewCheckout(params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	return cc.Success(result)
}

// BulkPurchaseOrderPreviewCheckout
// @Tags Buyer-PO
// @Summary Preview checkout bulk purchase order
// @Description Preview checkout bulk purchase order
// @Accept  json
// @Produce  json
// @Param data body repo.BulkPurchaseOrderPreviewCheckoutParams true "Form"
// @Success 200 {object} models.PurchaseOrder
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/bulk_purchase_orders/{bulk_purchase_order_id}/preview_checkout [post]
func BulkPurchaseOrderPreviewCheckout(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.BulkPurchaseOrderPreviewCheckoutParams

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, "")
	}

	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	params.JwtClaimsInfo = claims
	params.UpdatePricing = true
	result, err := repo.NewBulkPurchaseOrderRepo(cc.App.DB).BulkPurchaseOrderPreviewCheckout(params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	return cc.Success(result)
}

// GetBulkPurchaseOrder
// @Tags Buyer-PO
// @Summary Get bulk purchase order
// @Description Get bulk purchase order
// @Accept  json
// @Produce  json
// @Param purchase_order_id query string true "ID"
// @Success 200 {object} models.PurchaseOrder
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/bulk_purchase_orders/{bulk_purchase_order_id}/checkout [post]
func BulkPurchaseOrderCheckout(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var form repo.BulkPurchaseOrderCheckoutParams

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, "")
	}

	err = cc.BindAndValidate(&form)
	if err != nil {
		return eris.Wrap(err, "")
	}

	form.JwtClaimsInfo = claims
	result, err := repo.NewBulkPurchaseOrderRepo(cc.App.DB).BulkPurchaseOrderCheckout(form)
	if err != nil {
		return eris.Wrap(err, "")
	}

	if form.PaymentType == enums.PaymentTypeBankTransfer {
		var event = customerio.EventBulkPoFirstPaymentWaitingConfirmBankTransfer
		if form.Milestone == enums.PaymentMilestoneFinalPayment {
			event = customerio.EventBulkPoFinalPaymentWaitingConfirmBankTransfer

		}

		var eventData = result.GetCustomerIOMetadata(nil)
		if form.TransactionAttachment != nil {
			eventData["note"] = form.Note
			eventData["transaction_ref_id"] = form.TransactionRefID
			eventData["transaction_attachment"] = form.TransactionAttachment.GenerateFileURL()

		}
		for _, assigneeID := range result.AssigneeIDs {
			tasks.TrackCustomerIOTask{
				Event:  event,
				UserID: assigneeID,
				Data:   eventData,
			}.Dispatch(c.Request().Context())
		}

	}

	return cc.Success(result)
}

// UpdatePurchaseOrderTrackingStatus
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
// @Router /api/v1/bulk_bulk_purchase_orders/{bulk_purchase_order_id}/update_tracking_status [put]
func UpdateBulkPurchaseOrderTrackingStatus(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.UpdateBulkPurchaseOrderTrackingStatusParams

	claims, err := cc.GetJwtClaims()
	if err != nil {
		return eris.Wrap(err, "")
	}

	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, "")
	}

	params.UserID = claims.ID
	result, err := repo.NewBulkPurchaseOrderRepo(cc.App.DB).UpdatePurchaseOrderTrackingStatus(params)
	if err != nil {
		return eris.Wrap(err, "")
	}

	return cc.Success(result)
}

// PaginateBulkPurchaseOrderTracking
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
// @Router /api/v1/bulk_bulk_purchase_orders/{bulk_purchase_order_id}/logs [get]
func PaginateBulkPurchaseOrderTracking(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.PaginateBulkPurchaseOrderTrackingParams

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, "")
	}

	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, "")
	}

	params.JwtClaimsInfo = claims

	var result = repo.NewBulkPurchaseOrderTrackingRepo(cc.App.DB).PaginateBulkPurchaseOrderTrackings(params)

	return cc.Success(result)
}

// BulkPurchaseBuyerApproveQc
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
// @Router /api/v1/bulk_bulk_purchase_orders/{bulk_purchase_order_id}/approve_qc [put]
func BulkPurchaseBuyerApproveQc(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.BulkPurchaseOrderUpdateTrackingStatusParams

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, "")
	}

	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, "")
	}

	params.TrackingAction = enums.BulkPoTrackingActionBuyerApproveQc
	params.TrackingStatus = enums.BulkPoTrackingStatusSubmit
	params.JwtClaimsInfo = claims
	order, err := repo.NewBulkPurchaseOrderRepo(cc.App.DB).BulkPurchaseOrderUpdateTrackingStatus(params)
	if err != nil {
		return eris.Wrap(err, "")
	}

	for _, assigneeID := range order.AssigneeIDs {
		tasks.TrackCustomerIOTask{
			UserID: assigneeID,
			Event:  customerio.EventBulkPoBuyerApproveQc,
			Data:   order.GetCustomerIOMetadata(nil),
		}.Dispatch(c.Request().Context())
	}

	return cc.Success(order)
}

// BulkPurchaseBuyerApproveRawMaterial
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
// @Router /api/v1/admin/bulk_bulk_purchase_orders/{bulk_purchase_order_id}/approve_raw_material [post]
func BulkPurchaseBuyerApproveRawMaterial(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.BulkPurchaseBuyerApproveRawMaterialParams

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, "")
	}

	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, "")
	}

	params.JwtClaimsInfo = claims
	order, err := repo.NewBulkPurchaseOrderRepo(cc.App.DB).BulkPurchaseOrderBuyerApproveRawMaterial(params)
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
			Event:  customerio.EventBulkPoBuyerApproveRawMaterial,
			Data:   data,
		}.Dispatch(c.Request().Context())
	}

	return cc.Success(order)
}

// BuyerBulkPurchaseOrderConfirmDelivered
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
// @Router /api/v1/bulk_bulk_purchase_orders/{bulk_purchase_order_id}/confirm_delivered [post]
func BuyerBulkPurchaseOrderConfirmDelivered(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.BulkOrderConfirmDeliveredParams

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, "")
	}

	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, "")
	}

	params.JwtClaimsInfo = claims
	order, err := repo.NewBulkPurchaseOrderRepo(cc.App.DB).BuyerConfirmDelivered(params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	for _, assigneeID := range order.AssigneeIDs {
		tasks.TrackCustomerIOTask{
			UserID: assigneeID,
			Event:  customerio.EventBulkPoBuyerConfirmDelivered,
			Data:   order.GetCustomerIOMetadata(nil),
		}.Dispatch(c.Request().Context())
	}

	return cc.Success("Confirm Delivered")
}

// BulkPurchaseOrderGetSamplePO
// @Tags Buyer-PO
// @Summary BulkPurchaseOrderGetSamplePO
// @Description BulkPurchaseOrderGetSamplePO
// @Accept  json
// @Produce  json
// @Success 200 {object} models.PurchaseOrder
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/bulk_bulk_purchase_orders/{bulk_purchase_order_id}/sample_po [get]
func BulkPurchaseOrderGetSamplePO(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.BulkPurchaseOrderGetSamplePOParams

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, "")
	}

	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	params.JwtClaimsInfo = claims
	result, err := repo.NewBulkPurchaseOrderRepo(cc.App.DB).BulkPurchaseOrderGetSamplePO(params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	return cc.Success(result)
}

// BulkPurchaseOrderUpdateDesign UpdateDesign
// @Tags Buyer-PO
// @Summary BulkPurchaseOrderGetSamplePO
// @Description BulkPurchaseOrderGetSamplePO
// @Accept  json
// @Produce  json
// @Success 200 {object} models.PurchaseOrder
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/bulk_bulk_purchase_orders/{bulk_purchase_order_id}/update_design [get]
func BulkPurchaseOrderUpdateDesign(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.BulkPurchaseOrderUpdateDesignParams

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, "")
	}

	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	params.JwtClaimsInfo = claims
	result, err := repo.NewBulkPurchaseOrderRepo(cc.App.DB).UpdateDesign(params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	_, _ = tasks.UpdateUserProductClassesTask{
		UserID:              claims.GetUserID(),
		BulkPurchaseOrderID: result.ID,
	}.Dispatch(c.Request().Context())

	return cc.Success(result)
}

// BuyerUpdateBulkPurchaseOrderLog
// @Tags Buyer-BulkPurchaseOrder
// @Summary BulkPurchaseOrder cancel
// @Description BulkPurchaseOrder cancel
// @Accept  json
// @Produce  json
// @Success 200 {object} models.BulkPurchaseOrder
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/buyer/bulk_purchase_orders/{bulk_purchase_order_id}/logs [put]
func BuyerUpdateBulkPurchaseOrderLogs(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var form repo.UpdateBulkPurchaseOrderLogsParams

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, "")
	}

	err = cc.BindAndValidate(&form)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	form.JwtClaimsInfo = claims
	result, err := repo.NewBulkPurchaseOrderRepo(cc.App.DB).UpdateBulkPurchaseOrderLogs(form)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	return cc.Success(result)
}

// BuyerCreateBulkPurchaseOrderLogs
// @Tags Buyer-BulkPurchaseOrder
// @Summary BulkPurchaseOrder cancel
// @Description BulkPurchaseOrder cancel
// @Accept  json
// @Produce  json
// @Success 200 {object} models.BulkPurchaseOrder
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/buyer/bulk_purchase_orders/{bulk_purchase_order_id}/logs [delete]
func BuyerDeleteBulkPurchaseOrderLogs(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var form repo.DeleteBulkPurchaseOrderLogsParams

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, "")
	}

	err = cc.BindAndValidate(&form)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	form.JwtClaimsInfo = claims
	result, err := repo.NewBulkPurchaseOrderRepo(cc.App.DB).DeleteBulkPurchaseOrderLogs(form)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	return cc.Success(result)
}

// GetBulkPurchaseOrderInvoice
// @Tags Buyer-PO
// @Summary Get purchase order
// @Description Get purchase order
// @Accept  json
// @Produce  json
// @Param purchase_order_id query string true "ID"
// @Success 200 {object} models.PurchaseOrder
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/buyer/bulk_purchase_orders/{bulk_purchase_order_id}/invoice [get]
func GetBulkPurchaseOrderInvoice(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, "")
	}

	var params repo.GetBulkPurchaseOrderInvoiceParams
	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, "")
	}

	params.JwtClaimsInfo = claims
	result, err := repo.NewBulkPurchaseOrderRepo(cc.App.DB).GetBulkPurchaseOrderInvoice(params)
	if err != nil {
		if eris.Is(err, errs.ErrBulkPoInvoiceAlreadyGenerated) {
			return cc.Success(result)
		}
		return eris.Wrap(err, "")
	}

	return cc.Success(result)
}

// GetBulkPurchaseOrderFirstPaymentInvoice
// @Tags Buyer-PO
// @Summary Get purchase order
// @Description Get purchase order
// @Accept  json
// @Produce  json
// @Param bulk_purchase_order_id query string true "ID"
// @Success 200 {object} models.BulkPurchaseOrder
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/buyer/bulk_purchase_orders/{bulk_purchase_order_id}/first_payment_invoice [get]
func GetBulkPurchaseOrderFirstPaymentInvoice(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, "")
	}

	var params repo.CreateBulkFirstPaymentInvoiceParams
	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, "")
	}

	params.JwtClaimsInfo = claims
	result, err := repo.NewInvoiceRepo(cc.App.DB).CreateBulkFirstPaymentInvoice(params)
	if err != nil {
		if eris.Is(err, errs.ErrBulkPoInvoiceAlreadyGenerated) {
			return cc.Success(result)
		}
		return eris.Wrap(err, "")
	}

	return cc.Success(result)
}

// GetBulkPurchaseOrderFirstPaymentInvoice
// @Tags Buyer-PO
// @Summary Get purchase order
// @Description Get purchase order
// @Accept  json
// @Produce  json
// @Param bulk_purchase_order_id query string true "ID"
// @Success 200 {object} models.BulkPurchaseOrder
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/buyer/bulk_purchase_orders/{bulk_purchase_order_id}/final_payment_invoice [get]
func GetBulkPurchaseOrderFinalPaymentInvoice(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, "")
	}

	var params repo.CreateBulkFinalPaymentInvoiceParams
	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, "")
	}

	params.JwtClaimsInfo = claims
	result, err := repo.NewInvoiceRepo(cc.App.DB).CreateBulkFinalPaymentInvoice(params)
	if err != nil {
		if eris.Is(err, errs.ErrBulkPoInvoiceAlreadyGenerated) {
			return cc.Success(result)
		}
		return eris.Wrap(err, "")
	}

	return cc.Success(result)
}

// GetBulkPurchaseOrderDebitNotes
// @Tags Buyer-PO
// @Summary Get purchase order
// @Description Get purchase order
// @Accept  json
// @Produce  json
// @Param purchase_order_id query string true "ID"
// @Success 200 {object} models.PurchaseOrder
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/buyer/bulk_purchase_orders/{bulk_purchase_order_id}/debit_notes [get]
func GetBulkPurchaseOrderDebitNotes(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, "")
	}

	var params repo.CreateBulkDebitNotesParams
	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, "")
	}

	params.JwtClaimsInfo = claims
	result, err := repo.NewInvoiceRepo(cc.App.DB).CreateBulkDebitNotes(params)
	if err != nil {
		if eris.Is(err, errs.ErrBulkPoInvoiceAlreadyGenerated) {
			return cc.Success(result)
		}
		return eris.Wrap(err, "")
	}

	return cc.Success(result)
}

// BulkPurchaseOrderFeedback
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
// @Router /api/v1/buyer/bulk_purchase_orders/{bulk_purchase_order_id}/feedback [get]
func BulkPurchaseOrderFeedback(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	_, err := cc.GetJwtClaims()
	if err != nil {
		return eris.Wrap(err, "")
	}

	var params repo.BulkPurchaseOrderFeedbackParams
	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, "")
	}

	err = repo.NewBulkPurchaseOrderRepo(cc.App.DB).BuyerGiveBulkPurchaseOrderFeedback(params)
	if err != nil {
		return eris.Wrap(err, "")
	}

	return cc.Success("Received")
}

// BulkPurchaseOrderUpdatePps
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
// @Router /api/v1/buyer/bulk_purchase_orders/{purchase_order_id}/update_pps [put]
func BulkPurchaseOrderUpdatePps(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, "")
	}

	var params repo.BulkPurchaseOrderUpdatePpsParams
	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, "")
	}

	params.JwtClaimsInfo = claims
	bulkPO, err := repo.NewBulkPurchaseOrderRepo(cc.App.DB).BulkPurchaseOrderUpdatePps(params)
	if err != nil {
		return eris.Wrap(err, "")
	}

	return cc.Success(bulkPO)
}
