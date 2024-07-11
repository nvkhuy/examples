package controllers

import (
	"github.com/engineeringinflow/inflow-backend/pkg/geo"
	"github.com/engineeringinflow/inflow-backend/pkg/models"
	"github.com/labstack/echo/v4"
)

// SearchAddresses Search addresses
// @Tags Common
// @Summary Search addresses
// @Description Search addresses
// @Param keyword query string false "Keyword" default(1)
// @Param limit query int false "Size of page" default(20)
// @Accept  json
// @Produce  json
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/common/search_addresses [get]
func SearchAddresses(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params geo.SearchPlaceIndexForTextParams
	var err = cc.BindAndValidate(&params)
	if err != nil {
		return err
	}

	result, err := geo.GetInstance().SearchPlaceIndexForText(params)
	if err != nil {
		return err
	}

	return cc.Success(result)
}
