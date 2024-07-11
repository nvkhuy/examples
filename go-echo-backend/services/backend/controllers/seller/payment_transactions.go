package controllers

import (
	"github.com/engineeringinflow/inflow-backend/pkg/models"
	"github.com/engineeringinflow/inflow-backend/pkg/repo"
	"github.com/labstack/echo/v4"
	"github.com/rotisserie/eris"
)

// PaginatePaymentTransaction
// @Tags Seller-Payment
// @Summary Order payment transaction list
// @Description Order payment transaction list
// @Accept  json
// @Produce  json
// @Success 200 {object} []models.PaymentTransaction
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/seller/payment_transactions [get]
func PaginatePaymentTransaction(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.PaginatePaymentTransactionsParams

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	params.JwtClaimsInfo = claims
	var result = repo.NewPaymentTransactionRepo(cc.App.DB).PaginatePaymentTransactions(params)
	return cc.Success(result)
}

// GetPaymentTransaction
// @Tags Seller-Payment
// @Summary Order payment transaction list
// @Description Order payment transaction list
// @Accept  json
// @Produce  json
// @Success 200 {object} []models.PaymentTransaction
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/seller/payment_transactions/{payment_transactions_id} [get]
func GetPaymentTransaction(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.GetPaymentTransactionsParams

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	params.JwtClaimsInfo = claims
	result, err := repo.NewPaymentTransactionRepo(cc.App.DB).GetPaymentTransaction(params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}
	return cc.Success(result)
}
