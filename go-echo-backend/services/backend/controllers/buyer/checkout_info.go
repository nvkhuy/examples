package controllers

import (
	"github.com/engineeringinflow/inflow-backend/pkg/models"
	"github.com/engineeringinflow/inflow-backend/pkg/repo"
	"github.com/labstack/echo/v4"
	"github.com/rotisserie/eris"
)

// GetBulksPreviewCheckout
// @Tags Marketplace-Product
// @Summary GetCategoryTree
// @Description GetCategoryTree
// @Accept  json
// @Produce  json
// @Success 200 {object} models.User
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/buyer/{link_id}/checkout_info [get]
func GetCheckoutInfo(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.GetCheckoutInfoParams
	var err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	result, err := repo.NewCommonRepo(cc.App.DB).GetCheckoutInfo(params)
	if err != nil {
		return eris.Wrap(err, "")
	}

	return cc.Success(result)
}
