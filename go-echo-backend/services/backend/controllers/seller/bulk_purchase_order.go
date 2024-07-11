package controllers

import (
	"github.com/engineeringinflow/inflow-backend/pkg/customerio"
	"github.com/engineeringinflow/inflow-backend/pkg/models"
	"github.com/engineeringinflow/inflow-backend/pkg/models/enums"
	"github.com/engineeringinflow/inflow-backend/pkg/repo"
	"github.com/engineeringinflow/inflow-backend/services/consumer/tasks"
	"github.com/labstack/echo/v4"
	"github.com/rotisserie/eris"
)

// SellerPaginateBulkPurchaseOrders
// @Tags Seller-BPO
// @Summary Purchase orders
// @Description Purchase orders
// @Accept  json
// @Produce  json
// @Param keyword query string false "Keyword"
// @Param page query int false "Page number"
// @Param limit query int false "Limit"
// @Success 200 {object} query.Pagination{records=[]models.BulkPurchaseOrder}
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/seller/bulk_purchase_orders [get]
func SellerPaginateBulkPurchaseOrders(c echo.Context) error {
	var cc = c.(*models.CustomContext)

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	var params repo.PaginateBulkPurchaseOrderParams
	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	params.JwtClaimsInfo = claims
	var result = repo.NewBulkPurchaseOrderRepo(cc.App.DB).PaginateBulkPurchaseOrder(params)

	return cc.Success(result)
}

// GetBulkPurchaseOrder
// @Tags Seller-BPO
// @Summary Get purchase order
// @Description Get purchase order
// @Accept  json
// @Produce  json
// @Param purchase_order_id query string true "ID"
// @Success 200 {object} models.BulkPurchaseOrder
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/seller/bulk_purchase_orders/{bulk_purchase_order_id} [get]
func SellerGetBulkPurchaseOrder(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, "")
	}

	var params repo.GetBulkPurchaseOrderParams
	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, "")
	}

	params.JwtClaimsInfo = claims
	params.IncludeItems = true
	result, err := repo.NewBulkPurchaseOrderRepo(cc.App.DB).GetBulkPurchaseOrder(params)
	if err != nil {
		return eris.Wrap(err, "")
	}

	return cc.Success(result)
}

// SellerBulkPurchaseOrderUpdateRawMaterial
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
// @Router /api/v1/seller/bulk_purchase_orders/{bulk_purchase_order_id}/update_raw_material [put]
func SellerBulkPurchaseOrderUpdateRawMaterial(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, "")
	}

	var params repo.SellerBulkPurchaseOrderUpdateRawMaterialParams
	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, "")
	}

	params.JwtClaimsInfo = claims
	bulkPO, err := repo.NewSellerBulkPurchaseOrderRepo(cc.App.DB).SellerBulkPurchaseOrderUpdateRawMaterial(params)
	if err != nil {
		return eris.Wrap(err, "")
	}

	for _, userID := range bulkPO.AssigneeIDs {
		tasks.TrackCustomerIOTask{
			UserID: userID,
			Event:  customerio.EventSellerBulkPoUpdateRawMaterial,
			Data:   bulkPO.GetCustomerIOMetadata(nil),
		}.Dispatch(c.Request().Context())
	}

	return cc.Success(bulkPO)
}

// SellerBulkPurchaseOrderUpdatePps
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
// @Router /api/v1/seller/bulk_purchase_orders/{bulk_purchase_order_id}/update_pps [put]
func SellerBulkPurchaseOrderUpdatePps(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, "")
	}

	var params repo.SellerBulkPurchaseOrderUpdatePpsParams
	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, "")
	}

	params.JwtClaimsInfo = claims
	bulkPO, err := repo.NewSellerBulkPurchaseOrderRepo(cc.App.DB).SellerBulkPurchaseOrderUpdatePps(params)
	if err != nil {
		return eris.Wrap(err, "")
	}

	for _, userID := range bulkPO.AssigneeIDs {
		tasks.TrackCustomerIOTask{
			UserID: userID,
			Event:  customerio.EventSellerBulkPoUpdatePps,
			Data:   bulkPO.GetCustomerIOMetadata(nil),
		}.Dispatch(c.Request().Context())
	}

	return cc.Success(bulkPO)
}

// SellerBulkPurchaseOrderUpdateProduction
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
// @Router /api/v1/seller/bulk_purchase_orders/{bulk_purchase_order_id}/update_production [put]
func SellerBulkPurchaseOrderUpdateProduction(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, "")
	}

	var params repo.SellerBulkPurchaseOrderUpdateProductionParams
	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, "")
	}

	params.JwtClaimsInfo = claims
	bulkPO, err := repo.NewSellerBulkPurchaseOrderRepo(cc.App.DB).SellerBulkPurchaseOrderUpdateProduction(params)
	if err != nil {
		return eris.Wrap(err, "")
	}

	for _, userID := range bulkPO.AssigneeIDs {
		tasks.TrackCustomerIOTask{
			UserID: userID,
			Event:  customerio.EventSellerBulkPoUpdateProduction,
			Data:   bulkPO.GetCustomerIOMetadata(nil),
		}.Dispatch(c.Request().Context())
	}

	return cc.Success(bulkPO)
}

// SellerBulkPurchaseOrderCreateQcReport
// @Tags Admin-PO
// @Summary BulkPurchaseOrderCreateQcReport
// @Description BulkPurchaseOrderCreateQcReport
// @Accept  json
// @Produce  json
// @Param purchase_order_id query string true "ID"
// @Success 200 {object} models.PurchaseOrder
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/seller/bulk_purchase_orders/{bulk_purchase_order_id}/qc_report [post]
func SellerBulkPurchaseOrderCreateQcReport(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.SellerCreateQcReportParams

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, "")
	}

	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, "")
	}

	params.JwtClaimsInfo = claims
	order, err := repo.NewSellerBulkPurchaseOrderRepo(cc.App.DB).SellerBulkPurchaseOrderCreateQcReport(params)
	if err != nil {
		return eris.Wrap(err, "")
	}

	for _, userID := range order.AssigneeIDs {
		tasks.TrackCustomerIOTask{
			UserID: userID,
			Event:  customerio.EventBulkPoCreateQcReport,
			Data:   order.GetCustomerIOMetadata(nil),
		}.Dispatch(c.Request().Context())
	}

	return cc.Success(order)
}

// SellerPaginateInquirySellerRequests Seller update product photo
// @Tags Seller-Inquiry
// @Summary Seller update product photo
// @Description Seller update product photo
// @Accept  json
// @Produce  json
// @Success 200 {object} models.Inquiry
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/seller/bulk_purchase_orders/{bulk_purchase_order_id}/product_photo [put]
func SellerUpdateBulkPurchaseOrderProductPhoto(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.UpdateBulkPurchaseOrderProductPhotoParams

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	params.JwtClaimsInfo = claims
	result, err := repo.NewSellerBulkPurchaseOrderRepo(cc.App.DB).UpdateBulkPurchaseOrderProductPhoto(params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	return cc.Success(result)
}

// SellerUpdateBulkPurchaseOrderTechpack Seller update techpack
// @Tags Seller-Inquiry
// @Summary Seller update techpack
// @Description Seller update techpack
// @Accept  json
// @Produce  json
// @Success 200 {object} models.Inquiry
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/seller/bulk_purchase_orders/{bulk_purchase_order_id}/techpack [put]
func SellerUpdateBulkPurchaseOrderTechpack(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.UpdateBulkPurchaseOrderTechpackParams

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	params.JwtClaimsInfo = claims
	result, err := repo.NewSellerBulkPurchaseOrderRepo(cc.App.DB).UpdateBulkPurchaseOrderTechpack(params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	return cc.Success(result)
}

// SellerPaginateInquirySellerRequests Seller bill of material
// @Tags Seller-Inquiry
// @Summary Seller bill of material
// @Description Seller bill of material
// @Accept  json
// @Produce  json
// @Success 200 {object} models.Inquiry
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/seller/bulk_purchase_orders/{bulk_purchase_order_id}/bill_of_material [put]
func SellerUpdateBulkPurchaseOrderBillOfMaterial(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.UpdateBulkPurchaseOrderBillOfMaterialParams

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	params.JwtClaimsInfo = claims
	result, err := repo.NewSellerBulkPurchaseOrderRepo(cc.App.DB).UpdateBulkPurchaseOrderBillOfMaterial(params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	return cc.Success(result)
}

// SellerPaginateInquirySellerRequests Seller update size chart
// @Tags Seller-Inquiry
// @Summary Seller update size chart
// @Description Seller update size chart
// @Accept  json
// @Produce  json
// @Success 200 {object} models.Inquiry
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/seller/bulk_purchase_orders/{bulk_purchase_order_id}/size_chart [put]
func SellerUpdateBulkPurchaseOrderSizeChart(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.UpdateBulkPurchaseOrderSizeChartParams

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	params.JwtClaimsInfo = claims
	result, err := repo.NewSellerBulkPurchaseOrderRepo(cc.App.DB).UpdateBulkPurchaseOrderSizeChart(params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	return cc.Success(result)
}

// SellerPaginateInquirySellerRequests Seller size spec
// @Tags Seller-Inquiry
// @Summary Seller size spec
// @Description Seller size spec
// @Accept  json
// @Produce  json
// @Success 200 {object} models.Inquiry
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/seller/bulk_purchase_orders/{bulk_purchase_order_id}/size_spec [put]
func SellerUpdateBulkPurchaseOrderSizeSpec(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.UpdateBulkPurchaseOrderSizeSpecParams

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	params.JwtClaimsInfo = claims
	result, err := repo.NewSellerBulkPurchaseOrderRepo(cc.App.DB).UpdateBulkPurchaseOrderSizeSpec(params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	return cc.Success(result)
}

// SellerPaginateInquirySellerRequests Seller paginate inquiry quotation
// @Tags Seller-Inquiry
// @Summary Seller paginate inquiry quotation
// @Description Seller paginate inquiry quotation
// @Accept  json
// @Produce  json
// @Success 200 {object} models.Inquiry
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/seller/bulk_purchase_orders/{bulk_purchase_order_id}/size_grading [put]
func SellerUpdateBulkPurchaseOrderSizeGrading(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.UpdateBulkPurchaseOrderSizeGradingParams

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	params.JwtClaimsInfo = claims
	result, err := repo.NewSellerBulkPurchaseOrderRepo(cc.App.DB).UpdateBulkPurchaseOrderSizeGrading(params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	return cc.Success(result)
}

// SellerBulkPurchaseOrderApprovePO Seller paginate inquiry quotation
// @Tags Seller-Inquiry
// @Summary Seller paginate inquiry quotation
// @Description Seller paginate inquiry quotation
// @Accept  json
// @Produce  json
// @Success 200 {object} models.Inquiry
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/seller/bulk_purchase_orders/{bulk_purchase_order_id}/approve_po [put]
func SellerBulkPurchaseOrderApprovePO(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.BulkPurchaseOrderApprovePOParams

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	params.JwtClaimsInfo = claims
	result, err := repo.NewSellerBulkPurchaseOrderRepo(cc.App.DB).BulkPurchaseOrderApprovePO(params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	return cc.Success(result)
}

// SellerBulkPurchaseOrderRejectPO Seller paginate inquiry quotation
// @Tags Seller-Inquiry
// @Summary Seller paginate inquiry quotation
// @Description Seller paginate inquiry quotation
// @Accept  json
// @Produce  json
// @Success 200 {object} models.Inquiry
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/seller/bulk_purchase_orders/{bulk_purchase_order_id}/reject_po [delete]
func SellerBulkPurchaseOrderRejectPO(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.BulkPurchaseOrderRejectPOParams

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	params.JwtClaimsInfo = claims
	result, err := repo.NewSellerBulkPurchaseOrderRepo(cc.App.DB).BulkPurchaseOrderRejectPO(params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	return cc.Success(result)
}

// SellerBulkPurchaseOrderStartWithoutFirstPayment Seller paginate inquiry quotation
// @Tags Seller-Inquiry
// @Summary Seller paginate inquiry quotation
// @Description Seller paginate inquiry quotation
// @Accept  json
// @Produce  json
// @Success 200 {object} models.Inquiry
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/seller/bulk_purchase_orders/{bulk_purchase_order_id}/start_without_first_payment [put]
func SellerBulkPurchaseOrderStartWithoutFirstPayment(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.BulkPurchaseOrderStartWithoutFirstPaymentParams

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	params.JwtClaimsInfo = claims
	result, err := repo.NewSellerBulkPurchaseOrderRepo(cc.App.DB).BulkPurchaseOrderStartWithoutFirstPayment(params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	return cc.Success(result)
}

// SellerBulkPurchaseOrderStartWithoutFirstPayment Seller paginate inquiry quotation
// @Tags Seller-Inquiry
// @Summary Seller paginate inquiry quotation
// @Description Seller paginate inquiry quotation
// @Accept  json
// @Produce  json
// @Success 200 {object} models.Inquiry
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/seller/bulk_purchase_orders/{bulk_purchase_order_id}/confirm_receive_first_payment [put]
func SellerBulkPurchaseOrderConfirmReceiveFirstPayment(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.BulkPurchaseOrderConfirmReceiveFirstPaymentParams

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	params.JwtClaimsInfo = claims
	result, err := repo.NewSellerBulkPurchaseOrderRepo(cc.App.DB).BulkPurchaseOrderConfirmReceiveFirstPayment(params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	return cc.Success(result)
}

// SellerBulkPurchaseOrderStartWithoutFirstPayment Seller paginate inquiry quotation
// @Tags Seller-Inquiry
// @Summary Seller paginate inquiry quotation
// @Description Seller paginate inquiry quotation
// @Accept  json
// @Produce  json
// @Success 200 {object} models.Inquiry
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/seller/bulk_purchase_orders/{bulk_purchase_order_id}/confirm_receive_final_payment [put]
func SellerBulkPurchaseOrderConfirmReceiveFinalPayment(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.BulkPurchaseOrderConfirmReceiveFinalPaymentParams

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	params.JwtClaimsInfo = claims
	result, err := repo.NewSellerBulkPurchaseOrderRepo(cc.App.DB).BulkPurchaseOrderConfirmReceiveFinalPayment(params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	return cc.Success(result)
}

// SellerBulkPurchaseOrderMarkDelivering Seller paginate inquiry quotation
// @Tags Seller-Inquiry
// @Summary Seller paginate inquiry quotation
// @Description Seller paginate inquiry quotation
// @Accept  json
// @Produce  json
// @Success 200 {object} models.Inquiry
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/seller/bulk_purchase_orders/{bulk_purchase_order_id}/mark_delivering [put]
func SellerBulkPurchaseOrderMarkDelivering(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.SellerBulkPurchaseOrderMarkDeliveringParams

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	params.JwtClaimsInfo = claims
	result, err := repo.NewSellerBulkPurchaseOrderRepo(cc.App.DB).SellerBulkPurchaseOrderMarkDelivering(params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	return cc.Success(result)
}

// PaginateBulkPurchaseOrderTracking
// @Tags Admin-PO
// @Summary PaginateBulkPurchaseOrderTracking
// @Description PaginateBulkPurchaseOrderTracking
// @Accept  json
// @Produce  json
// @Param purchase_order_id query string true "ID"
// @Success 200 {object} models.PurchaseOrder
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/seller/bulk_purchase_orders/{purchase_order_id}/logs [get]
func SellerPaginateBulkPurchaseOrderTracking(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, "")
	}

	var params repo.PaginateBulkPurchaseOrderTrackingParams
	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, "")
	}

	params.JwtClaimsInfo = claims
	params.IncludeTrackings = true
	params.UserGroup = enums.PoTrackingUserGroupSeller
	var result = repo.NewBulkPurchaseOrderTrackingRepo(cc.App.DB).PaginateBulkPurchaseOrderTrackings(params)

	return cc.Success(result)
}

// PaginateBulkPurchaseOrderTracking
// @Tags Admin-PO
// @Summary PaginateBulkPurchaseOrderTracking
// @Description PaginateBulkPurchaseOrderTracking
// @Accept  json
// @Produce  json
// @Param purchase_order_id query string true "ID"
// @Success 200 {object} models.PurchaseOrder
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/seller/bulk_purchase_orders/{purchase_order_id}/feedback [put]
func SellerBulkPurchaseOrderFeedback(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, "")
	}

	var params repo.SellerBulkPurchaseOrderFeedbackParams
	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, "")
	}

	params.JwtClaimsInfo = claims
	result, err := repo.NewSellerBulkPurchaseOrderRepo(cc.App.DB).SellerBulkPurchaseOrderFeedback(params)
	if err != nil {
		return eris.Wrap(err, "")
	}

	return cc.Success(result)
}

// SellerBulkPurchaseOrderFirstPaymentInvoice
// @Tags Admin-PO
// @Summary SellerBulkPurchaseOrderFirstPaymentInvoice
// @Description SellerBulkPurchaseOrderFirstPaymentInvoice
// @Accept  json
// @Produce  json
// @Param purchase_order_id query string true "ID"
// @Success 200 {object} models.PurchaseOrder
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/seller/bulk_purchase_orders/{purchase_order_id}/first_payment_invoice [get]
func SellerBulkPurchaseOrderFirstPaymentInvoice(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, "")
	}

	var params repo.SellerBulkPurchaseOrderFirstPaymentInvoiceParams
	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, "")
	}

	params.JwtClaimsInfo = claims
	result, err := repo.NewSellerBulkPurchaseOrderRepo(cc.App.DB).SellerBulkPurchaseOrderFirstPaymentInvoice(params)
	if err != nil {
		return eris.Wrap(err, "")
	}

	return cc.Success(result)
}

// SellerBulkPurchaseOrderFirstPaymentInvoice
// @Tags Admin-PO
// @Summary SellerBulkPurchaseOrderFirstPaymentInvoice
// @Description SellerBulkPurchaseOrderFirstPaymentInvoice
// @Accept  json
// @Produce  json
// @Param purchase_order_id query string true "ID"
// @Success 200 {object} models.PurchaseOrder
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/seller/bulk_purchase_orders/{purchase_order_id}/final_payment_invoice [get]
func SellerBulkPurchaseOrderFinalPaymentInvoice(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, "")
	}

	var params repo.SellerBulkPurchaseOrderFinalPaymentInvoiceParams
	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, "")
	}

	params.JwtClaimsInfo = claims
	result, err := repo.NewSellerBulkPurchaseOrderRepo(cc.App.DB).SellerBulkPurchaseOrderFinalPaymentInvoice(params)
	if err != nil {
		return eris.Wrap(err, "")
	}

	return cc.Success(result)
}

// SellerBulkPurchaseOrderMarkRawMaterial Admin paginate inquiry quotation
// @Tags Admin-Inquiry
// @Summary Admin paginate inquiry quotation
// @Description Admin paginate inquiry quotation
// @Accept  json
// @Produce  json
// @Success 200 {object} models.Inquiry
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/seller/bulk_purchase_orders/{bulk_purchase_order_id}/mark_raw_material [put]
func SellerBulkPurchaseOrderMarkRawMaterial(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.SellerBulkPurchaseOrderUpdateTrackingStatusParams

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	params.TrackingAction = enums.BulkPoTrackingActionSellerMarkRawMaterial
	params.SellerTrackingStatus = enums.SellerBulkPoTrackingStatusRawMaterial
	params.JwtClaimsInfo = claims
	bulkPO, err := repo.NewSellerBulkPurchaseOrderRepo(cc.App.DB).SellerBulkPurchaseOrderUpdateTrackingStatus(params)
	if err != nil {
		return eris.Wrap(err, "")
	}

	for _, assigneeID := range bulkPO.AssigneeIDs {
		tasks.TrackCustomerIOTask{
			UserID: assigneeID,
			Event:  customerio.EventSellerBulkPoMarkRawMaterial,
			Data:   bulkPO.GetCustomerIOMetadata(nil),
		}.Dispatch(c.Request().Context())
	}

	return cc.Success(bulkPO)
}

// BulkPurchaseOrderMarkProduction
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
// @Router /api/v1/seller/bulk_purchase_orders/{purchase_order_id}/mark_production [put]
func SellerBulkPurchaseOrderMarkProduction(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.SellerBulkPurchaseOrderMarkProductionParams
	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, "")
	}

	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, "")
	}

	params.JwtClaimsInfo = claims
	bulkPO, err := repo.NewSellerBulkPurchaseOrderRepo(cc.App.DB).SellerBulkPurchaseOrderMarkProduction(params)
	if err != nil {
		return eris.Wrap(err, "")
	}

	for _, assigneeID := range bulkPO.AssigneeIDs {
		tasks.TrackCustomerIOTask{
			UserID: assigneeID,
			Event:  customerio.EventSellerBulkPoMarkProduction,
			Data:   bulkPO.GetCustomerIOMetadata(nil),
		}.Dispatch(c.Request().Context())
	}

	return cc.Success(bulkPO)
}

// BulkPurchaseOrderMarkInspection
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
// @Router /api/v1/seller/bulk_purchase_orders/{purchase_order_id}/mark_inspection [put]
func SellerBulkPurchaseOrderMarkInspection(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.SellerBulkPurchaseOrderMarkInspectionParams
	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, "")
	}

	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, "")
	}

	params.JwtClaimsInfo = claims
	bulkPO, err := repo.NewSellerBulkPurchaseOrderRepo(cc.App.DB).SellerBulkPurchaseOrderMarkInspection(params)
	if err != nil {
		return eris.Wrap(err, "")
	}

	for _, assigneeID := range bulkPO.AssigneeIDs {
		tasks.TrackCustomerIOTask{
			UserID: assigneeID,
			Event:  customerio.EventSellerBulkPoMarkInspection,
			Data:   bulkPO.GetCustomerIOMetadata(nil),
		}.Dispatch(c.Request().Context())
	}

	return cc.Success(bulkPO)
}
