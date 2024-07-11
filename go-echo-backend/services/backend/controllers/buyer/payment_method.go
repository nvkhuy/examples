package controllers

import (
	"github.com/engineeringinflow/inflow-backend/pkg/models"
	"github.com/engineeringinflow/inflow-backend/pkg/repo"
	"github.com/labstack/echo/v4"
	"github.com/rotisserie/eris"
)

// AddPaymentMethod Add payment method
// @Tags Marketplace-Payment-Method
// @Summary Add payment method
// @Description Add payment method
// @Accept  json
// @Produce  json
// @Param data body models.UserPaymentMethodCreateForm true "CreateFromPayload form"
// @Header 200 {string} Bearer YOUR_TOKEN
// @Failure 404 {object} errs.Error
// @Security ApiKeyAuth
// @Router /api/v1/payment_methods/create [post]
func AddPaymentMethod(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var form models.UserPaymentMethodCreateForm

	claims, err := cc.GetJwtClaims()
	if err != nil {
		return eris.Wrap(err, "")
	}

	err = cc.BindAndValidate(&form)
	if err != nil {
		return eris.Wrap(err, "")
	}

	pm, err := repo.NewPaymentMethodRepo(cc.App.DB).AttachPaymentMethod(claims.ID, form)
	if err != nil {
		return eris.Wrap(err, "")
	}

	return cc.Success(pm)
}

// RemovePaymentMethod Remove payment method
// @Tags Marketplace-Payment-Method
// @Summary Remove payment method
// @Description Remove payment method
// @Accept  json
// @Produce  json
// @Param id path string true "Payment method id"
// @Header 200 {string} Bearer YOUR_TOKEN
// @Failure 404 {object} errs.Error
// @Security ApiKeyAuth
// @Router /api/v1/payment_methods/{id}/detact [delete]
func RemovePaymentMethod(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.DetactPaymentMethodParams

	claims, err := cc.GetJwtClaims()
	if err != nil {
		return eris.Wrap(err, "")
	}

	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, "")
	}
	params.UserID = claims.ID
	pm, err := repo.NewPaymentMethodRepo(cc.App.DB).DetachPaymentMethod(params)
	if err != nil {
		return eris.Wrap(err, "")
	}

	return cc.Success(pm)
}

// GetPaymentMethod Get payment method
// @Tags Marketplace-Payment-Method
// @Summary Get payment method
// @Description Get payment method
// @Accept  json
// @Produce  json
// @Param id path string true "Payment method id"
// @Header 200 {string} Bearer YOUR_TOKEN
// @Failure 404 {object} errs.Error
// @Security ApiKeyAuth
// @Router /api/v1/payment_methods/{id} [get]
func GetPaymentMethod(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.GetPaymentMethodParams

	claims, err := cc.GetJwtClaims()
	if err != nil {
		return eris.Wrap(err, "")
	}

	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, "")
	}

	params.UserID = claims.ID

	result, err := repo.NewPaymentMethodRepo(cc.App.DB).GetPaymentMethod(params)
	if err != nil {
		return eris.Wrap(err, "")
	}

	return cc.Success(result)
}

// GetPaymentMethods Get payment methods
// @Tags Marketplace-Payment-Method
// @Summary Get payment methods
// @Description Get payment methods
// @Accept  json
// @Produce  json
// @Header 200 {string} Bearer YOUR_TOKEN
// @Failure 404 {object} errs.Error
// @Security ApiKeyAuth
// @Router /api/v1/payment_methods [get]
func GetPaymentMethods(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.GetPaymentMethodsParams
	claims, err := cc.GetJwtClaims()
	if err != nil {
		return eris.Wrap(err, "")
	}

	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, "")
	}

	params.UserID = claims.ID

	list, err := repo.NewPaymentMethodRepo(cc.App.DB).GetPaymentMethods(params)
	if err != nil {
		return eris.Wrap(err, "")
	}

	return cc.Success(list)
}

// GetPaymentMethods Get payment methods
// @Tags Marketplace-Payment-Method
// @Summary Get payment methods
// @Description Get payment methods
// @Accept  json
// @Produce  json
// @Header 200 {string} Bearer YOUR_TOKEN
// @Failure 404 {object} errs.Error
// @Security ApiKeyAuth
// @Router /api/v1/payment_methods/{id}/default [put]
func MarkDefaultPaymentMethod(c echo.Context) error {
	var cc = c.(*models.CustomContext)

	var params repo.MarkDefaultPaymentMethodParams

	claims, err := cc.GetJwtClaims()
	if err != nil {
		return eris.Wrap(err, "")
	}

	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, "")
	}

	params.UserID = claims.ID
	list, err := repo.NewPaymentMethodRepo(cc.App.DB).MarkDefaultPaymentMethod(params)
	if err != nil {
		return eris.Wrap(err, "")
	}

	return cc.Success(list)
}
