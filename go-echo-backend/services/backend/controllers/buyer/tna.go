package controllers

import (
	"github.com/engineeringinflow/inflow-backend/pkg/models"
	"github.com/engineeringinflow/inflow-backend/pkg/repo"
	"github.com/labstack/echo/v4"
	"github.com/rotisserie/eris"
)

// CreateTNA create category
// @Tags Admin-Fabric
// @Summary CreateFromPayload category
// @Description CreateFromPayload category
// @Accept  json
// @Produce  json
// @Param data body models.FabricUpdateForm true "Form"
// @Success 200 {object} models.Fabric
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/tnas [post]
func CreateTNA(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.CreateTNAsParams

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	params.JwtClaimsInfo = claims
	result, err := repo.NewTNARepo(cc.App.DB).Create(params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	return cc.Success(result)
}

// UpdateTNA update category
// @Tags Admin-Fabric
// @Summary CreateFromPayload category
// @Description CreateFromPayload category
// @Accept  json
// @Produce  json
// @Param data body models.FabricUpdateForm true "Form"
// @Success 200 {object} models.Fabric
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/tnas/:id [post]
func UpdateTNA(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.UpdateTNAsParams

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	params.JwtClaimsInfo = claims
	result, err := repo.NewTNARepo(cc.App.DB).Update(params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	return cc.Success(result)
}

// PaginateTNA
// @Tags Admin-Blog
// @Summary PaginateBlogCategory
// @Description PaginateBlogCategory
// @Accept  json
// @Produce  json
// @Param keyword query string false "Keyword"
// @Param name query string false "Name"
// @Success 200 {object} models.BlogCategory
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/tnas [get]
func PaginateTNA(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.PaginateTNAsParams

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	params.JwtClaimsInfo = claims
	result := repo.NewTNARepo(cc.App.DB).Paginate(params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	return cc.Success(result)
}

// DeleteTNA update category
// @Tags Admin-Fabric
// @Summary CreateFromPayload category
// @Description CreateFromPayload category
// @Accept  json
// @Produce  json
// @Param data body models.FabricUpdateForm true "Form"
// @Success 200 {object} models.Fabric
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/tnas/:id [delete]
func DeleteTNA(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.DeleteTNAsParams

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	params.JwtClaimsInfo = claims
	err = repo.NewTNARepo(cc.App.DB).Delete(params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	return cc.Success("deleted")
}
