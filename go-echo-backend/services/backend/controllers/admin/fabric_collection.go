package controllers

import (
	"github.com/engineeringinflow/inflow-backend/pkg/models"
	"github.com/engineeringinflow/inflow-backend/pkg/repo"
	"github.com/labstack/echo/v4"
	"github.com/rotisserie/eris"
)

// CreateFabricCollection create category
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
// @Router /api/v1/admin/fabric_collections [post]
func CreateFabricCollection(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.CreateFabricCollectionParams

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	params.JwtClaimsInfo = claims
	result, err := repo.NewFabricCollectionRepo(cc.App.DB).Create(params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	return cc.Success(result)
}

// UpdateFabricCollection Update category
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
func UpdateFabricCollection(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.UpdateFabricCollectionParams

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	params.JwtClaimsInfo = claims
	result, err := repo.NewFabricCollectionRepo(cc.App.DB).Update(params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	return cc.Success(result)
}

// PaginateFabric
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
func PaginateFabric(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.PaginateFabricParams

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	params.JwtClaimsInfo = claims
	result := repo.NewFabricRepo(cc.App.DB).Paginate(params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	return cc.Success(result)
}

// DetailsFabric
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
func DetailsFabric(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.DetailsFabricParams

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	params.JwtClaimsInfo = claims
	result, err := repo.NewFabricRepo(cc.App.DB).Details(params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	return cc.Success(result)
}

// DeleteFabric
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
// @Router /api/v1/admin/fabrics/{id} [delete]
func DeleteFabric(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.DeleteFabricParams

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	params.JwtClaimsInfo = claims
	err = repo.NewFabricRepo(cc.App.DB).Delete(params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	return cc.Success("deleted")
}

// AddFabricToCollection create category
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
// @Router /api/v1/admin//fabric_collections/:id/remove_fabric [post]
func AddFabricToCollection(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.AddFabricToCollectionParams

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	params.JwtClaimsInfo = claims
	err = repo.NewFabricCollectionRepo(cc.App.DB).AddFabric(params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	return cc.Success("added")
}

// RemoveFabricFromCollection create category
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
// @Router /api/v1/admin/fabric_collections/:id/remove_fabric [delete]
func RemoveFabricFromCollection(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.RemoveFabricFromCollectionParams

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	params.JwtClaimsInfo = claims
	_, err = repo.NewFabricCollectionRepo(cc.App.DB).RemoveFabric(params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	return cc.Success("removed")
}
