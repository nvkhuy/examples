package controllers

import (
	"fmt"

	"github.com/engineeringinflow/inflow-backend/pkg/models"
	"github.com/engineeringinflow/inflow-backend/pkg/models/enums"
	"github.com/engineeringinflow/inflow-backend/pkg/repo"
	"github.com/engineeringinflow/inflow-backend/services/consumer/tasks"
	"github.com/labstack/echo/v4"
	"github.com/rotisserie/eris"
)

// BuyerGetOrderCart
// @Tags Marketplace-OrderCart
// @Summary Get order cart
// @Description This API allows admin to list chat room with pagination
// @Accept  json
// @Produce  json
// @Success 200 {object} models.GetOrderCartResponse
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router api/v1/buyer/order_cart/preview_checkout [get]
func BuyerGetOrderCart(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var req models.GetOrderCartRequest

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, "")
	}

	err = cc.BindAndValidate(&req)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	req.JwtClaimsInfo = claims
	result, err := repo.NewOrderCartRepo(cc.App.DB).GetOrderCart(&req)

	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	return cc.Success(result)
}

// BuyerGetOrderCartPreviewCheckout
// @Tags Marketplace-OrderCart
// @Summary Get order cart
// @Description This API allows admin to get order cart
// @Accept  json
// @Produce  json
// @Param data body models.OrderCartPreviewCheckoutRequest true
// @Success 200 {object} models.GetOrderCartResponse
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router api/v1/buyer/order_cart/preview_checkout [post]
func BuyerGetOrderCartPreviewCheckout(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var req models.OrderCartPreviewCheckoutRequest

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, "")
	}

	err = cc.BindAndValidate(&req)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	req.JwtClaimsInfo = claims
	result, err := repo.NewOrderCartRepo(cc.App.DB).OrderCartPreviewCheckout(&req)

	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	return cc.Success(result)
}

// BuyerOrderCartCheckout
// @Tags Marketplace-OrderCart
// @Summary Get order cart
// @Description This API allows admin to checkout order cart
// @Accept  json
// @Produce  json
// @Param data body models.OrderCartCheckoutRequest true
// @Success 200 {object} models.OrderCartCheckoutResponse
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router api/v1/buyer/order_cart/checkout [post]
func BuyerOrderCartCheckout(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var req models.OrderCartCheckoutRequest

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, "")
	}

	err = cc.BindAndValidate(&req)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	req.JwtClaimsInfo = claims
	result, err := repo.NewOrderCartRepo(cc.App.DB).CheckoutOrders(&req)

	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	if len(result.PurchaseOrders) > 0 {
		for _, po := range result.PurchaseOrders {
			tasks.CreateInquiryAuditTask{
				Form: models.InquiryAuditCreateForm{
					InquiryID:       po.InquiryID,
					ActionType:      enums.AuditActionTypeInquirySamplePoCreated,
					UserID:          po.UserID,
					Description:     fmt.Sprintf("New sample PO %s has been created for inquiry", po.ReferenceID),
					PurchaseOrderID: po.ID,
				},
			}.Dispatch(c.Request().Context())
		}
	}

	if result.PaymentTransaction.PaymentType == enums.PaymentTypeBankTransfer {
		tasks.NotifyAdminConfirmPaymentTask{
			PaymentTransactionID: result.PaymentTransaction.ID,
		}.Dispatch(c.Request().Context())
	}

	return cc.Success(result)
}

// BuyerOrderCartCheckout
// @Tags Marketplace-OrderCart
// @Summary Get order cart
// @Description This API allows admin to get order cart checkout info
// @Accept  json
// @Produce  json
// @Param checkout_session_id query string true
// @Success 200 {object} models.GetOrderCartResponse
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router api/v1/buyer/order_cart/checkout_info [get]
func BuyerOrderCartGetCheckoutInfo(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var req models.OrderCartGetCheckoutInfoRequest

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, "")
	}

	err = cc.BindAndValidate(&req)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	req.JwtClaimsInfo = claims
	result, err := repo.NewOrderCartRepo(cc.App.DB).GetOrderCartCheckoutInfo(&req)

	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	return cc.Success(result)
}
