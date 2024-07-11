package controllers

import (
	"github.com/engineeringinflow/inflow-backend/pkg/models"
	"github.com/engineeringinflow/inflow-backend/pkg/repo"
	"github.com/labstack/echo/v4"
	"github.com/rotisserie/eris"
)

// SellerBulkPOSubmitQuotation
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
// @Router /api/v1/seller/bulk_quotations/{seller_quotation_id}/submit_quotation [post]
func SellerBulkPOSubmitQuotation(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var form repo.SellerSubmitBulkQuotationParams

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, "")
	}

	err = cc.BindAndValidate(&form)
	if err != nil {
		return eris.Wrap(err, "")
	}

	form.JwtClaimsInfo = claims
	result, err := repo.NewSellerBulkPurchaseOrderRepo(cc.App.DB).SellerSubmitBulkQuotation(form)
	if err != nil {
		return eris.Wrap(err, "")
	}

	return cc.Success(result)
}

// SellerBulkPOReSubmitQuotation
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
// @Router /api/v1/seller/bulk_quotations/{seller_quotation_id}/re_submit_quotation [post]
func SellerBulkPOReSubmitQuotation(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.SellerSubmitBulkQuotationParams

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, "")
	}

	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, "")
	}

	params.JwtClaimsInfo = claims
	result, err := repo.NewSellerBulkPurchaseOrderRepo(cc.App.DB).SellerReSubmitBulkQuotation(params)
	if err != nil {
		return eris.Wrap(err, "")
	}

	return cc.Success(result)
}

// SellerGetBulkPOQuotationDetails
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
// @Router /api/v1/seller/bulk_quotations/{seller_quotation_id} [get]
func SellerGetBulkPOQuotationDetails(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.GetBulkPurchaseOrderSellerQuotationParams

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, "")
	}

	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, "")
	}

	params.JwtClaimsInfo = claims
	result, err := repo.NewSellerBulkPurchaseOrderRepo(cc.App.DB).GetBulkPurchaseOrderSellerQuotation(params)
	if err != nil {
		return eris.Wrap(err, "")
	}

	return cc.Success(result)
}

// SellerBulkPOSubmitMultipleQuotation
// @Tags Seller-BPO
// @Summary Seller submit multiple bulk quotations
// @Description Seller submit multiple bulk quotations
// @Accept  json
// @Produce  json
// @Param data body models.SellerSubmitMultipleBulkQuotationsParams true "Form"
// @Success 200 {object} models.BulkPurchaseOrderSellerQuotations
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/seller/bulk_quotations/submit_multiple_quotations [post]
func SellerBulkPOSubmitMultipleQuotation(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var form repo.SellerSubmitMultipleBulkQuotationsParams

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, "")
	}

	err = cc.BindAndValidate(&form)
	if err != nil {
		return eris.Wrap(err, "")
	}

	form.JwtClaimsInfo = claims
	result, err := repo.NewSellerBulkPurchaseOrderRepo(cc.App.DB).SubmitMultipleBulkQuotations(&form)
	if err != nil {
		return eris.Wrap(err, "")
	}

	return cc.Success(result)
}
