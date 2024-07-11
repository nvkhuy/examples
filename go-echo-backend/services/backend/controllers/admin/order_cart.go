package controllers

import (
	"github.com/engineeringinflow/inflow-backend/pkg/models"
	"github.com/engineeringinflow/inflow-backend/pkg/repo"

	"github.com/labstack/echo/v4"
	"github.com/rotisserie/eris"
)

// GetBuyerOrderCart
// @Tags Marketplace-OrderCart
// @Summary Get buyer order cart
// @Description This API allows admin to get order cart from buyer
// @Accept  json
// @Produce  json
// @Param buyer_id path string true "Buyer ID"
// @Param data body models.GetBuyerOrderCartRequest false
// @Success 200 {object} models.GetBuyerOrderCartResponse
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router api/v1/admin/order_cart/:buyer_id/preview [post]
func GetBuyerOrderCartPreview(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var req models.GetBuyerOrderCartRequest

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, "")
	}

	err = cc.BindAndValidate(&req)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	req.JwtClaimsInfo = claims
	result, err := repo.NewOrderCartRepo(cc.App.DB).GetBuyerOrderCartPreview(&req)

	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	return cc.Success(result)
}

// CreateBuyerPaymentLink
// @Tags Marketplace-OrderCart
// @Summary Create buyer payment link
// @Description This API allows admin to create payment link for buyer
// @Accept  json
// @Produce  json
// @Param buyer_id path string true "Buyer ID"
// @Param data body models.CreateBuyerPaymentLinkRequest false
// @Success 200 {object} models.CreateBuyerPaymentLinkResponse
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router api/v1/admin/order_cart/:buyer_id/create_payment_link [post]
func CreateBuyerPaymentLink(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var req models.CreateBuyerPaymentLinkRequest

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, "")
	}

	err = cc.BindAndValidate(&req)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	req.JwtClaimsInfo = claims
	result, err := repo.NewOrderCartRepo(cc.App.DB).CreateBuyerPaymentLink(&req)

	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	return cc.Success(result)
}
