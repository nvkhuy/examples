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
)

// AdminSellerBulkPurchaseOrderFirstPayout Admin first payout
// @Tags Admin-PO
// @Summary Admin first payout
// @Description Admin first payout
// @Accept  json
// @Produce  json
// @Success 200 {object} models.PurchaseOrder
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/seller_bulk_purchase_orders/{bulk_purchase_order_id}/first_payout [post]
func AdminSellerBulkPurchaseOrderFirstPayout(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.AdminBulkPurchaseOrderFirstPayoutParams

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, "")
	}

	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	params.JwtClaimsInfo = claims
	result, err := repo.NewSellerBulkPurchaseOrderRepo(cc.App.DB).AdminBulkPurchaseOrderFirstPayout(params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	return cc.Success(result)
}

// AdminSellerBulkPurchaseOrderFinalPayout Admin final payout
// @Tags Admin-PO
// @Summary Admin final payout
// @Description Admin final payout
// @Accept  json
// @Produce  json
// @Success 200 {object} models.PurchaseOrder
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/seller_bulk_purchase_orders/{bulk_purchase_order_id}/final_payout [post]
func AdminSellerBulkPurchaseOrderFinalPayout(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.AdminBulkPurchaseOrderFinalPayoutParams

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, "")
	}

	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	params.JwtClaimsInfo = claims
	result, err := repo.NewSellerBulkPurchaseOrderRepo(cc.App.DB).AdminBulkPurchaseOrderFinalPayout(params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	return cc.Success(result)
}

// AdminSellerBulkPurchaseOrderAssignMaker Admin PO assign maker
// @Tags Admin-PO
// @Summary Admin PO assign maker
// @Description Admin PO assign maker
// @Accept  json
// @Produce  json
// @Success 200 {object} models.PurchaseOrder
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/seller_bulk_purchase_orders/{bulk_purchase_order_id}/assign_maker [post]
func AdminSellerBulkPurchaseOrderAssignMaker(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.AdminAssignBulkPOMakerParams

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, "")
	}

	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	params.JwtClaimsInfo = claims
	result, err := repo.NewBulkPurchaseOrderRepo(cc.App.DB).AdminAssignBulkPOMaker(params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	return cc.Success(result)
}

// AdminSellerBulkPurchaseOrderAllocations
// @Tags Admin-Inquiry
// @Summary Mark seen comments
// @Description Mark seen comments
// @Accept  json
// @Produce  json
// @Param data body models.Comment true "Form"
// @Success 200 {object} models.Inquiry
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/seller_bulk_purchase_orders/{bulk_purchase_order_id}/seller_allocations [get]
func AdminSellerBulkPurchaseOrderAllocations(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.PaginateBulkPurchaseOrderAllocationParams

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, "")
	}
	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, "")
	}

	params.JwtClaimsInfo = claims
	params.Statuses = append(params.Statuses, enums.InquirySellerStatusApproved)
	var result = repo.NewBulkPurchaseOrderRepo(cc.App.DB).PaginateBulkPurchaseOrderAllocation(params)

	return cc.Success(result)
}

// AdminSellerPaginateBulkPurchaseOrderMatchingSellers Admin paginate matching sellers
// @Tags Admin-Inquiry
// @Summary Admin paginate matching sellers
// @Description Admin paginate matching sellers
// @Accept  json
// @Produce  json
// @Success 200 {object} models.Inquiry
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/seller_bulk_purchase_orders/{bulk_purchase_order_id}/matching_sellers [get]
func AdminSellerPaginateBulkPurchaseOrderMatchingSellers(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.PaginateBulkPurchaseOrderMatchingSellersParams

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	params.JwtClaimsInfo = claims
	var result = repo.NewSellerBulkPurchaseOrderRepo(cc.App.DB).PaginateBulkPurchaseOrderMatchingSellers(params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	return cc.Success(result)
}

// AdminSellerBulkPurhcaseOrderSendToSeller Admin send inquiry to seller
// @Tags Admin-Inquiry
// @Summary Admin send inquiry to seller
// @Description Admin send inquiry to seller
// @Accept  json
// @Produce  json
// @Success 200 {object} models.Inquiry
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/seller_bulk_purchase_orders/{bulk_purchase_order_id}/send_to_seller [post]
func AdminSellerBulkPurhcaseOrderSendToSeller(c echo.Context) error {
	var cc = c.(*models.CustomContext)

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	var params repo.SendBulkPurchaseOrderToSellerParams
	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	params.JwtClaimsInfo = claims
	result, err := repo.NewSellerBulkPurchaseOrderRepo(cc.App.DB).SendBulkPurchaseOrderToSeller(params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	if len(result) > 0 {
		bulk, err := repo.NewBulkPurchaseOrderRepo(cc.App.DB).GetBulkPurchaseOrder(repo.GetBulkPurchaseOrderParams{
			BulkPurchaseOrderID: params.BulkPurchaseOrderID,
			JwtClaimsInfo:       params.JwtClaimsInfo,
		})
		if err == nil {
			for _, record := range result {
				tasks.TrackCustomerIOTask{
					UserID: record.UserID,
					Event:  customerio.EventSellerNewBulkPurchaseOrderQuotation,
					Data: bulk.GetCustomerIOMetadata(map[string]interface{}{
						"offer_price":  record.OfferPrice,
						"offer_remark": record.OfferRemark,
					}),
				}.Dispatch(cc.Request().Context())

				tasks.CreateChatRoomTask{
					UserID:              claims.GetUserID(),
					Role:                claims.GetRole(),
					BulkPurchaseOrderID: bulk.ID,
					SellerID:            record.UserID,
				}.Dispatch(c.Request().Context())
			}
		}

	}

	return cc.Success(result)
}

// AdminSellerPaginateBulkPurchaseOrderSellerQuotations Admin paginate inquiry quotation
// @Tags Admin-Inquiry
// @Summary Admin paginate inquiry quotation
// @Description Admin paginate inquiry quotation
// @Accept  json
// @Produce  json
// @Success 200 {object} models.Inquiry
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/seller_bulk_purchase_orders/{bulk_purchase_order_id}/seller_quotations [get]
func AdminSellerPaginateBulkPurchaseOrderSellerQuotations(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.PaginateBulkPurchaseOrderSellerQuotationsParams

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	params.JwtClaimsInfo = claims
	var result = repo.NewSellerBulkPurchaseOrderRepo(cc.App.DB).PaginateBulkPurchaseOrderSellerQuotations(params)

	return cc.Success(result)
}

// AdminSellerBulkPurchaseOrderUpdateProductPhoto Admin paginate inquiry quotation
// @Tags Admin-Inquiry
// @Summary Admin paginate inquiry quotation
// @Description Admin paginate inquiry quotation
// @Accept  json
// @Produce  json
// @Success 200 {object} models.Inquiry
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/seller_bulk_purchase_orders/{bulk_purchase_order_id}/product_photo [put]
func AdminSellerBulkPurchaseOrderUpdateProductPhoto(c echo.Context) error {
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

// AdminPaginateInquirySellerRequests Admin paginate inquiry quotation
// @Tags Admin-Inquiry
// @Summary Admin paginate inquiry quotation
// @Description Admin paginate inquiry quotation
// @Accept  json
// @Produce  json
// @Success 200 {object} models.Inquiry
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/seller_bulk_purchase_orders/{bulk_purchase_order_id}/techpack [put]
func AdminSellerBulkPurchaseOrderUpdateTechpack(c echo.Context) error {
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

// AdminSellerBulkPurchaseOrderUpdateBillOfMaterial Admin paginate inquiry quotation
// @Tags Admin-Inquiry
// @Summary Admin paginate inquiry quotation
// @Description Admin paginate inquiry quotation
// @Accept  json
// @Produce  json
// @Success 200 {object} models.Inquiry
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/seller_bulk_purchase_orders/{bulk_purchase_order_id}/bill_of_material [put]
func AdminSellerBulkPurchaseOrderUpdateBillOfMaterial(c echo.Context) error {
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

// AdminSellerBulkPurchaseOrderUpdateSizeChart Admin paginate inquiry quotation
// @Tags Admin-Inquiry
// @Summary Admin paginate inquiry quotation
// @Description Admin paginate inquiry quotation
// @Accept  json
// @Produce  json
// @Success 200 {object} models.Inquiry
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/seller_bulk_purchase_orders/{bulk_purchase_order_id}/size_chart [put]
func AdminSellerBulkPurchaseOrderUpdateSizeChart(c echo.Context) error {
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

// AdminSellerBulkPurchaseOrderUpdateSizeSpec Admin paginate inquiry quotation
// @Tags Admin-Inquiry
// @Summary Admin paginate inquiry quotation
// @Description Admin paginate inquiry quotation
// @Accept  json
// @Produce  json
// @Success 200 {object} models.Inquiry
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/seller_bulk_purchase_orders/{bulk_purchase_order_id}/size_spec [put]
func AdminSellerBulkPurchaseOrderUpdateSizeSpec(c echo.Context) error {
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

// AdminSellerBulkPurchaseOrderUpdateSizeGrading Admin paginate inquiry quotation
// @Tags Admin-Inquiry
// @Summary Admin paginate inquiry quotation
// @Description Admin paginate inquiry quotation
// @Accept  json
// @Produce  json
// @Success 200 {object} models.Inquiry
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/seller_bulk_purchase_orders/{bulk_purchase_order_id}/size_grading [put]
func AdminSellerBulkPurchaseOrderUpdateSizeGrading(c echo.Context) error {
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

// AdminSellerBulkPurchaseOrderUpdatePacking Admin paginate inquiry quotation
// @Tags Admin-Inquiry
// @Summary Admin paginate inquiry quotation
// @Description Admin paginate inquiry quotation
// @Accept  json
// @Produce  json
// @Success 200 {object} models.Inquiry
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/seller_bulk_purchase_orders/{bulk_purchase_order_id}/packing [put]
func AdminSellerBulkPurchaseOrderUpdatePacking(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.UpdateBulkPurchaseOrderPackingParams

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	params.JwtClaimsInfo = claims
	result, err := repo.NewSellerBulkPurchaseOrderRepo(cc.App.DB).UpdateBulkPurchaseOrderPacking(params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	return cc.Success(result)
}

// SellerBulkPurchaseOrderUpdateInspectionProcedure Admin paginate inquiry quotation
// @Tags Admin-Inquiry
// @Summary Admin paginate inquiry quotation
// @Description Admin paginate inquiry quotation
// @Accept  json
// @Produce  json
// @Success 200 {object} models.Inquiry
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/seller_bulk_purchase_orders/{bulk_purchase_order_id}/inspection_procedure [put]
func AdminSellerBulkPurchaseOrderUpdateInspectionProcedure(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.AdminSellerBulkPurchaseOrderUpdateInspectionProcedureParams

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	params.JwtClaimsInfo = claims
	result, err := repo.NewSellerBulkPurchaseOrderRepo(cc.App.DB).AdminSellerBulkPurchaseOrderUpdateInspectionProcedure(params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	return cc.Success(result)
}

// AdminSellerBulkPurchaseUpdateInspectionTestingRequirements Admin paginate inquiry quotation
// @Tags Admin-Inquiry
// @Summary Admin paginate inquiry quotation
// @Description Admin paginate inquiry quotation
// @Accept  json
// @Produce  json
// @Success 200 {object} models.Inquiry
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/seller_bulk_purchase_orders/{bulk_purchase_order_id}/inspection_testing_requirements [put]
func AdminSellerBulkPurchaseUpdateInspectionTestingRequirements(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.UpdateBulkPurchaseOrderTestingRequirementsParams

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	params.JwtClaimsInfo = claims
	result, err := repo.NewSellerBulkPurchaseOrderRepo(cc.App.DB).UpdateBulkPurchaseOrderTestingRequirements(params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	return cc.Success(result)
}

// AdminSellerBulkPurchaseOrderUpdateLabelGuide Admin paginate inquiry quotation
// @Tags Admin-Inquiry
// @Summary Admin paginate inquiry quotation
// @Description Admin paginate inquiry quotation
// @Accept  json
// @Produce  json
// @Success 200 {object} models.Inquiry
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/seller_bulk_purchase_orders/{bulk_purchase_order_id}/label_guide [put]
func AdminSellerBulkPurchaseOrderUpdateLabelGuide(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.UpdateBulkPurchaseOrderLabelGuideAttachmentsParams

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	params.JwtClaimsInfo = claims
	result, err := repo.NewSellerBulkPurchaseOrderRepo(cc.App.DB).UpdateBulkPurchaseOrderLabelGuideAttachments(params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	return cc.Success(result)
}

// AdminSellerBulkPurchaseOrderUpdatePointOfMeasurement Admin paginate inquiry quotation
// @Tags Admin-Inquiry
// @Summary Admin paginate inquiry quotation
// @Description Admin paginate inquiry quotation
// @Accept  json
// @Produce  json
// @Success 200 {object} models.Inquiry
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/seller_bulk_purchase_orders/{bulk_purchase_order_id}/point_of_measurement [put]
func AdminSellerBulkPurchaseOrderUpdatePointOfMeasurement(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.UpdateBulkPurchaseOrderPointOfMeasurementAttachmentsParams

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	params.JwtClaimsInfo = claims
	result, err := repo.NewSellerBulkPurchaseOrderRepo(cc.App.DB).UpdateBulkPurchaseOrderPointOfMeasurementAttachments(params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	return cc.Success(result)
}

// AdminSellerBulkPurchaseOrderApproveQuotation
// @Tags Admin-Inquiry
// @Summary Mark seen comments
// @Description Mark seen comments
// @Accept  json
// @Produce  json
// @Param data body models.Comment true "Form"
// @Success 200 {object} models.Inquiry
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/seller_bulk_purchase_orders/{seller_quotation_id}/reject [delete]
func AdminSellerBulkPurchaseOrderApproveQuotation(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.AdminApproveSellerBulkPurchaseOrderQuotationParams

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, "")
	}

	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, "")
	}

	params.JwtClaimsInfo = claims
	err = repo.NewSellerBulkPurchaseOrderRepo(cc.App.DB).AdminApproveSellerBulkPurchaseOrderQuotation(params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	tasks.AdminApproveSellerBulkPurchaseOrderQuotation{
		AdminID:           claims.GetUserID(),
		SellerQuotationID: params.SellerQuotationID,
	}.Dispatch(c.Request().Context())

	return cc.Success("Approved")
}

// AdminSellerBulkPurchaseOrderRejectQuotation
// @Tags Admin-Inquiry
// @Summary Mark seen comments
// @Description Mark seen comments
// @Accept  json
// @Produce  json
// @Param data body models.Comment true "Form"
// @Success 200 {object} models.Inquiry
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/seller_bulk_purchase_orders/{seller_quotation_id}/reject [delete]
func AdminSellerBulkPurchaseOrderRejectQuotation(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.AdminRejectSellerBulkPurchaseOrderQuotationParams

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, "")
	}

	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, "")
	}

	params.JwtClaimsInfo = claims
	err = repo.NewSellerBulkPurchaseOrderRepo(cc.App.DB).AdminRejectSellerBulkPurchaseOrderQuotation(params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	tasks.AdminRejectSellerBulkPurchaseOrderQuotation{
		AdminID:           claims.GetUserID(),
		SellerQuotationID: params.SellerQuotationID,
	}.Dispatch(c.Request().Context())

	return cc.Success("Rejected")
}

// AdminSellerPurchaseOrderUploadPo
// @Tags Admin-Seller-PO
// @Summary upload po seller purchase order
// @Description upload po purchase order
// @Accept  json
// @Produce  json
// @Param purchase_order_id query string true "ID"
// @Success 200 {object} models.PurchaseOrder
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/seller_bulk_purchase_orders/{bulk_purchase_order_id}/upload_po [put]
func AdminSellerBulkPurchaseOrderUploadPo(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, "")
	}

	var params repo.AdminSellerBulkPurchaseOrderUploadPoParams
	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, "")
	}

	params.JwtClaimsInfo = claims

	order, err := repo.NewBulkPurchaseOrderRepo(cc.App.DB).AdminSellerBulkPurchaseOrderUploadPo(params)
	if err != nil {
		return eris.Wrap(err, "")
	}

	return cc.Success(order)
}

// AdminSellerBulkPurchaseOrderApproveSellerQuotation
// @Tags Admin-BulkPurchaseOrder
// @Summary Approve seller quotation
// @Description Approve seller quotation
// @Accept  json
// @Produce  json
// @Success 200 {object} models.Comment
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/seller_bulk_purchase_orders/{seller_quotation_id}/approve [post]
func AdminSellerBulkPurchaseOrderApproveSellerQuotation(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.AdminApproveSellerBulkPurchaseOrderQuotationParams

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	params.JwtClaimsInfo = claims
	err = repo.NewSellerBulkPurchaseOrderRepo(cc.App.DB).AdminApproveSellerBulkPurchaseOrderQuotation(params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	return cc.Success("Approved")
}

// AdminSellerBulkPurchaseOrderRejectSellerQuotation
// @Tags Admin-BulkPurchaseOrder
// @Summary Reject seller quotation
// @Description Reject seller quotation
// @Accept  json
// @Produce  json
// @Success 200 {object} models.Comment
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/seller_bulk_purchase_orders/{seller_quotation_id}/reject [delete]
func AdminSellerBulkPurchaseOrderRejectSellerQuotation(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.AdminRejectSellerBulkPurchaseOrderQuotationParams

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	params.JwtClaimsInfo = claims
	err = repo.NewSellerBulkPurchaseOrderRepo(cc.App.DB).AdminRejectSellerBulkPurchaseOrderQuotation(params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	return cc.Success("Rejected")
}

// AdminSellerBulkPurchaseOrderUpdateLabelGuide Admin paginate inquiry quotation
// @Tags Admin-Inquiry
// @Summary Admin paginate inquiry quotation
// @Description Admin paginate inquiry quotation
// @Accept  json
// @Produce  json
// @Success 200 {object} models.Inquiry
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/seller_bulk_purchase_orders/{bulk_purchase_order_id}/mark_delivered [put]
func AdminSellerBulkPurchaseOrderMarkDelivered(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.AdminSellerBulkPoConfirmDeliveredParams

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	params.JwtClaimsInfo = claims
	result, err := repo.NewSellerBulkPurchaseOrderRepo(cc.App.DB).AdminSellerBulkPoConfirmDelivered(params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	return cc.Success(result)
}

// AdminSellerBulkPurchaseOrderMarkRawMaterial Admin paginate inquiry quotation
// @Tags Admin-Inquiry
// @Summary Admin paginate inquiry quotation
// @Description Admin paginate inquiry quotation
// @Accept  json
// @Produce  json
// @Success 200 {object} models.Inquiry
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/seller_bulk_purchase_orders/{bulk_purchase_order_id}/mark_raw_material [put]
func AdminSellerBulkPurchaseOrderMarkRawMaterial(c echo.Context) error {
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

	params.TrackingAction = enums.BulkPoTrackingActionMarkRawMaterial
	params.SellerTrackingStatus = enums.SellerBulkPoTrackingStatusRawMaterial
	params.JwtClaimsInfo = claims
	bulkPO, err := repo.NewSellerBulkPurchaseOrderRepo(cc.App.DB).SellerBulkPurchaseOrderUpdateTrackingStatus(params)
	if err != nil {
		return eris.Wrap(err, "")
	}

	tasks.TrackCustomerIOTask{
		UserID: bulkPO.SellerID,
		Event:  customerio.EventBulkPoMarkRawMaterial,
		Data:   bulkPO.GetCustomerIOMetadata(nil),
	}.Dispatch(c.Request().Context())

	return cc.Success(bulkPO)
}

// AdminSellerBulkPurchaseOrderUpdatePps
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
// @Router /api/v1/admin/seller_bulk_purchase_orders/{purchase_order_id}/update_pps [put]
func AdminSellerBulkPurchaseOrderUpdatePps(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, "")
	}

	var params repo.AdminSellerBulkPurchaseOrderUpdatePpsParams
	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, "")
	}

	params.JwtClaimsInfo = claims
	bulkPO, err := repo.NewSellerBulkPurchaseOrderRepo(cc.App.DB).AdminSellerBulkPurchaseOrderUpdatePps(params)
	if err != nil {
		return eris.Wrap(err, "")
	}

	tasks.CreateUserNotificationTask{
		UserID:           bulkPO.SellerID,
		Message:          fmt.Sprintf("Pre-Production info was updated for %s", bulkPO.ReferenceID),
		NotificationType: enums.UserNotificationTypeBulkPoUpdatePps,
		Metadata: &models.UserNotificationMetadata{
			AdminID:                      claims.GetUserID(),
			BulkPurchaseOrderID:          bulkPO.ID,
			BulkPurchaseOrderReferenceID: bulkPO.ReferenceID,
			InquiryID: func() string {
				if bulkPO.Inquiry != nil {
					return bulkPO.Inquiry.ID
				}
				return ""
			}(),
			InquiryReferenceID: func() string {
				if bulkPO.Inquiry != nil {
					return bulkPO.Inquiry.ReferenceID
				}
				return ""
			}(),
		},
	}.Dispatch(c.Request().Context())

	tasks.TrackCustomerIOTask{
		UserID: bulkPO.SellerID,
		Event:  customerio.EventBulkPoUpdatePps,
		Data:   bulkPO.GetCustomerIOMetadata(nil),
	}.Dispatch(c.Request().Context())
	return cc.Success(bulkPO)
}
