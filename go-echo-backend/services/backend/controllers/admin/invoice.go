package controllers

import (
	"github.com/engineeringinflow/inflow-backend/pkg/models"
	"github.com/engineeringinflow/inflow-backend/pkg/repo"
	"github.com/labstack/echo/v4"
	"github.com/rotisserie/eris"
)

// PaginateInvoice
// @Tags Admin-Product
// @Summary create product
// @Description create product
// @Accept  json
// @Produce  json
// @Param data body models.ProductCreateForm true "Form"
// @Success 200 {object} models.Product
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/invoice [get]
func PaginateInvoices(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.PaginateInvoicesParams

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	params.JwtClaimsInfo = claims
	resp := repo.NewInvoiceRepo(cc.App.DB).PaginateInvoices(params)
	return cc.Success(resp)
}

// DetailsInvoice
// @Tags Admin-Product
// @Summary create product
// @Description create product
// @Accept  json
// @Produce  json
// @Param data body models.ProductCreateForm true "Form"
// @Success 200 {object} models.Product
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/invoice/{invoice_number} [get]
func DetailsInvoice(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.InvoiceDetailsPrams

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	params.JwtClaimsInfo = claims
	var result models.Invoice
	result, err = repo.NewInvoiceRepo(cc.App.DB).Details(params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}
	return cc.Success(result)
}

// CreateInvoice
// @Tags Admin-Product
// @Summary create product
// @Description create product
// @Accept  json
// @Produce  json
// @Param data body models.ProductCreateForm true "Form"
// @Success 200 {object} models.Product
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/invoice [post]
func CreateInvoice(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params models.CreateInvoiceParams

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	params.JwtClaimsInfo = claims
	resp, err := repo.NewInvoiceRepo(cc.App.DB).CreateInvoice(params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}
	return cc.Success(resp)
}

// UpdateInvoice
// @Tags Admin-Product
// @Summary create product
// @Description create product
// @Accept  json
// @Produce  json
// @Param data body models.ProductCreateForm true "Form"
// @Success 200 {object} models.Product
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/invoice/:invoice_number [post]
func UpdateInvoice(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params models.UpdateInvoiceParams

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	params.JwtClaimsInfo = claims
	var resp *models.Invoice
	resp, err = repo.NewInvoiceRepo(cc.App.DB).UpdateInvoice(params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}
	return cc.Success(resp)
}

// ExitsInvoiceNumber
// @Tags Admin-Product
// @Summary create product
// @Description create product
// @Accept  json
// @Produce  json
// @Param data body models.ProductCreateForm true "Form"
// @Success 200 {object} models.Product
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/invoice/{invoice_number}/exits [get]
func ExitsInvoiceNumber(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.InvoiceDetailsPrams

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	params.JwtClaimsInfo = claims
	var result bool
	result, err = repo.NewInvoiceRepo(cc.App.DB).IsExitsInvoiceNumber(params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}
	return cc.Success(result)
}

// NextInvoiceNumber
// @Tags Admin-Product
// @Summary create product
// @Description create product
// @Accept  json
// @Produce  json
// @Param data body models.ProductCreateForm true "Form"
// @Success 200 {object} models.Product
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/invoice/invoice_number/next [get]
func NextInvoiceNumber(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.NextInvoiceNumberParams

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	params.JwtClaimsInfo = claims
	var result int
	result, err = repo.NewInvoiceRepo(cc.App.DB).NextInvoiceNumber()
	if err != nil {
		return eris.Wrap(err, err.Error())
	}
	return cc.Success(result)
}

// NextInvoiceNumber
// @Tags Admin-Product
// @Summary create product
// @Description create product
// @Accept  json
// @Produce  json
// @Param data body models.ProductCreateForm true "Form"
// @Success 200 {object} models.Product
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/invoice/{invoice_number}/attachment [get]
func GetInvoiceAttachment(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.GetInvoiceAttachmentParams

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	params.JwtClaimsInfo = claims
	result, err := repo.NewInvoiceRepo(cc.App.DB).GetInvoiceAttachment(params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}
	return cc.Success(result)
}
