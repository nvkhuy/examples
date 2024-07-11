package controllers

import (
	"github.com/engineeringinflow/inflow-backend/pkg/models"
	"github.com/engineeringinflow/inflow-backend/pkg/repo"
	"github.com/labstack/echo/v4"
	"github.com/rotisserie/eris"
)

// SubscribeUpdates Subscribe updates
// @Tags Marketplace-Support
// @Summary Subscribe updates
// @Description Subscribe updates
// @Accept  json
// @Produce  json
// @Param data body models.FactoryTourUpdateForm true "Form"
// @Success 200 {object} models.FactoryTour
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/subscribe [post]
func SubscribeUpdates(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.SubscribeUpdatesParams

	err := cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	err = repo.NewCommonRepo(cc.App.DB).SubscribeUpdates(params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	return cc.Success("Subscribed")
}
