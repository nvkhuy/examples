package controllers

import (
	"github.com/engineeringinflow/inflow-backend/pkg/models"
	"github.com/engineeringinflow/inflow-backend/pkg/repo"
	"github.com/engineeringinflow/inflow-backend/services/consumer/tasks"
	"github.com/labstack/echo/v4"
	"github.com/rotisserie/eris"
)

// PaginatePaymentTransaction
// @Tags Admin-Order
// @Summary Order payment transaction list
// @Description Order payment transaction list
// @Accept  json
// @Produce  json
// @Success 200 {object} []models.PaymentTransaction
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/payment_transactions [get]
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
	params.IncludeDetails = true
	var result = repo.NewPaymentTransactionRepo(cc.App.DB).PaginatePaymentTransactions(params)
	return cc.Success(result)
}

// GetPaymentTransaction
// @Tags Admin-Order
// @Summary Order payment transaction list
// @Description Order payment transaction list
// @Accept  json
// @Produce  json
// @Success 200 {object} []models.PaymentTransaction
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/payment_transactions/{payment_transactions_id} [get]
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
	params.IncludeDetails = true
	params.IncludeInvoice = true
	result, err := repo.NewPaymentTransactionRepo(cc.App.DB).GetPaymentTransaction(params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}
	return cc.Success(result)
}

// ApprovePaymentTransaction
// @Tags Admin-Order
// @Summary Order payment transaction list
// @Description Order payment transaction list
// @Accept  json
// @Produce  json
// @Success 200 {object} []models.PaymentTransaction
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/payment_transactions/{payment_transactions_id}/approve [put]
func ApprovePaymentTransaction(c echo.Context) error {
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
	transaction, err := repo.NewPaymentTransactionRepo(cc.App.DB).ApprovePaymentTransactions(params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	tasks.CreatePaymentInvoiceTask{
		PaymentTransactionID: transaction.ID,
		ApprovedByUserID:     claims.GetUserID(),
	}.Dispatch(cc.Request().Context())

	// if len(transaction.PurchaseOrders) > 0 {
	// 	for _, po := range transaction.PurchaseOrders {
	// _, _ = tasks.CreatePOPaymentInvoiceTask{
	// 	PurchaseOrderID:  po.ID,
	// 	ApprovedByUserID: claims.GetUserID(),
	// }.Dispatch(cc.Request().Context())
	// 	}
	// }
	// if len(transaction.BulkPurchaseOrders) > 0 {
	// 	for _, bpo := range transaction.BulkPurchaseOrders {
	// 		if bpo.TrackingStatus == enums.BulkPoTrackingStatusFirstPaymentConfirmed {
	// 			tasks.CreateBulkPoFirstPaymentInvoiceTask{
	// 				ApprovedByUserID:    claims.GetUserID(),
	// 				BulkPurchaseOrderID: bpo.ID,
	// 			}.Dispatch(cc.Request().Context())
	// 		}
	// 		if bpo.TrackingStatus == enums.BulkPoTrackingStatusFinalPaymentConfirmed {
	// 			tasks.CreateBulkPoFinalPaymentInvoiceTask{
	// 				ApprovedByUserID:    claims.GetUserID(),
	// 				BulkPurchaseOrderID: bpo.ID,
	// 			}.Dispatch(cc.Request().Context())
	// 		}
	// 	}
	// }

	return cc.Success("Approved")
}

// RejectPaymentTransaction
// @Tags Admin-Order
// @Summary Order payment transaction list
// @Description Order payment transaction list
// @Accept  json
// @Produce  json
// @Success 200 {object} []models.PaymentTransaction
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/payment_transactions/{payment_transactions_id}/reject [put]
func RejectPaymentTransaction(c echo.Context) error {
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
	result, err := repo.NewPaymentTransactionRepo(cc.App.DB).RejectPaymentTransactions(params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}
	return cc.Success(result)
}

// GetPaymentTransactionAttachments
// @Tags Admin-Order
// @Summary View payment transaction attachments
// @Description View payment transaction attachments
// @Accept  json
// @Produce  json
// @Success 200 {object} []models.PaymentTransaction
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/payment_transactions/{payment_transactions_id}/attachments [get]
func GetPaymentTransactionAttachments(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.GetPaymentTransactionAttachmentsParams

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	params.JwtClaimsInfo = claims
	result, err := repo.NewPaymentTransactionRepo(cc.App.DB).GetPaymentTransactionAttachments(params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}
	return cc.Success(result)
}

// ExportPaymentTransactions
// @Tags Admin-Order
// @Summary Export payment transactions
// @Description Export payment transactions
// @Accept  json
// @Produce  json
// @Success 200 {object} []models.PaymentTransaction
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/payment_transactions/export [get]
func ExportPaymentTransactions(c echo.Context) error {
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
	result, err := repo.NewPaymentTransactionRepo(cc.App.DB).ExportExcel(params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	return cc.Success(result)
}
