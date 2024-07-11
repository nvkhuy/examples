package controllers

import (
	"github.com/engineeringinflow/inflow-backend/pkg/models"
	"github.com/engineeringinflow/inflow-backend/pkg/models/enums"
	"github.com/engineeringinflow/inflow-backend/pkg/repo"
	"github.com/labstack/echo/v4"
	"github.com/rotisserie/eris"
)

// StatsSuppliers Stats suppliers
// @Tags Admin-Stats
// @Summary Stats suppliers
// @Description Stats suppliers
// @Accept  json
// @Produce  json
// @Param from_time query int false "From"
// @Param to_time query int false "To"
// @Success 200 {object} models.StatsSuppliers
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/stats/suppliers [get]
func StatsSuppliers(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.StatsSuppliersParams

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	err = cc.Bind(&params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	params.JwtClaimsInfo = claims

	var result = repo.NewStatsRepo(cc.App.DB).StatsSuppliers(params)

	return cc.Success(result)
}

// StatBuyers Stats buyers
// @Tags Admin-Stats
// @Summary Stats buyers
// @Description Stats buyers
// @Accept  json
// @Produce  json
// @Param from_time query int false "From"
// @Param to_time query int false "To"
// @Success 200 {object} models.StatsBuyers
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/stats/buyers [get]
func StatBuyers(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.StatsBuyersParams
	var err = cc.Bind(&params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	params.ForRole = enums.RoleSuperAdmin

	var result = repo.NewStatsRepo(cc.App.DB).StatsBuyers(params)

	return cc.Success(result)
}

// StatsProducts Stats products
// @Tags Admin-Stats
// @Summary Stats products
// @Description Stats products
// @Accept  json
// @Produce  json
// @Param from_time query int false "From"
// @Param to_time query int false "To"
// @Success 200 {object} models.StatsProducts
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/stats/products [get]
func StatsProducts(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.StatsProductsParams
	var err = cc.Bind(&params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	params.ForRole = enums.RoleSuperAdmin

	var result = repo.NewStatsRepo(cc.App.DB).StatsProducts(params)

	return cc.Success(result)
}
