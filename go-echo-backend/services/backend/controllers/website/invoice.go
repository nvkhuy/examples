package controllers

import (
	"github.com/engineeringinflow/inflow-backend/pkg/models"
	"github.com/engineeringinflow/inflow-backend/pkg/repo"
	"github.com/labstack/echo/v4"
	"github.com/rotisserie/eris"
)

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

	err := cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	var result models.Invoice
	result, err = repo.NewInvoiceRepo(cc.App.DB).Details(params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}
	return cc.Success(result)
}
