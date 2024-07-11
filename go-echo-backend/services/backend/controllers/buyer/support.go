package controllers

import (
	"github.com/engineeringinflow/inflow-backend/pkg/models"
	"github.com/engineeringinflow/inflow-backend/pkg/models/enums"
	"github.com/engineeringinflow/inflow-backend/pkg/repo"
	"github.com/labstack/echo/v4"
	"github.com/rotisserie/eris"
)

// CreateFactoryTour create FactoryTour
// @Tags Marketplace-Support
// @Summary create FactoryTour
// @Description create FactoryTour
// @Accept  json
// @Produce  json
// @Param data body models.FactoryTourUpdateForm true "Form"
// @Success 200 {object} models.FactoryTour
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/support/factory_tours/create [post]
func CreateFactoryTour(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var form models.FactoryTourUpdateForm

	err := cc.BindAndValidate(&form)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	form.ForRole = enums.RoleClient

	cate, err := repo.NewFactoryTourRepo(cc.App.DB).CreateFactoryTour(form)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	return cc.Success(cate)
}
