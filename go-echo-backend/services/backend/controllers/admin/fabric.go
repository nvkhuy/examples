package controllers

import (
	"github.com/engineeringinflow/inflow-backend/pkg/models"
	"github.com/engineeringinflow/inflow-backend/pkg/repo"
	"github.com/labstack/echo/v4"
	"github.com/rotisserie/eris"
)

// CreateFabric create category
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
// @Router /api/v1/admin/fabrics [post]
func CreateFabric(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.CreateFabricParams

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	params.JwtClaimsInfo = claims
	result, err := repo.NewFabricRepo(cc.App.DB).Create(params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	return cc.Success(result)
}

// UpdateFabric Update category
// @Tags Admin-Fabric
// @Summary Update category
// @Description Update category
// @Accept  json
// @Produce  json
// @Param user_id path string true "ID"
// @Param data body models.FabricUpdateForm true "Form"
// @Success 200 {object} models.Fabric
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/fabrics/{id} [put]
func UpdateFabric(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.UpdateFabricParams

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	params.JwtClaimsInfo = claims
	result, err := repo.NewFabricRepo(cc.App.DB).Update(params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	return cc.Success(result)
}

// PaginateFabricCollection
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
// @Router /api/v1/admin/fabrics [get]
func PaginateFabricCollection(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.PaginateFabricCollectionParams

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	params.JwtClaimsInfo = claims
	result := repo.NewFabricCollectionRepo(cc.App.DB).Paginate(params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	return cc.Success(result)
}

// DetailsFabricCollection
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
// @Router /api/v1/admin/fabrics/{id} [get]
func DetailsFabricCollection(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.DetailsFabricCollectionParams

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	params.JwtClaimsInfo = claims
	result, err := repo.NewFabricCollectionRepo(cc.App.DB).Details(params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	return cc.Success(result)
}

// DeleteFabricCollection
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
// @Router /api/v1/admin/fabric_collections/{id} [delete]
func DeleteFabricCollection(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.DeleteFabricCollectionParams

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	params.JwtClaimsInfo = claims
	err = repo.NewFabricCollectionRepo(cc.App.DB).Delete(params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	return cc.Success("deleted")
}
