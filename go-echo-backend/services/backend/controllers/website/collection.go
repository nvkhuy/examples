package controllers

import (
	"github.com/engineeringinflow/inflow-backend/pkg/models"
	"github.com/engineeringinflow/inflow-backend/pkg/repo"
	"github.com/labstack/echo/v4"
	"github.com/rotisserie/eris"
)

// CollectionReadyDesign
// @Tags Marketplace-Collection
// @Summary Collection PreDesign
// @Description Collection PreDesign
// @Accept  json
// @Produce  json
// @Success 200 {object} models.Collection
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/collections/ready_design [get]
func CollectionReadyDesign(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.PaginateCollectionParams
	var err = cc.Bind(&params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	var result = repo.NewCollectionRepo(cc.App.DB).SearchCollection(params)

	return cc.Success(result)
}

// PaginateCollection
// @Tags Marketplace-Collection
// @Summary Collection PreDesign
// @Description Collection PreDesign
// @Accept  json
// @Produce  json
// @Success 200 {object} models.Collection
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/collections [get]
func PaginateCollection(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.PaginateCollectionParams
	var err = cc.Bind(&params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	var result = repo.NewCollectionRepo(cc.App.DB).PaginateCollection(params)

	return cc.Success(result)
}

// CollectionGetProduct
// @Tags Marketplace-Collection
// @Summary Collection GetProduct
// @Description Collection GetProduct
// @Accept  json
// @Produce  json
// @Success 200 {object} models.Product
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/collections/{id}/get_product [get]
func CollectionGetProduct(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.PaginateCollectionProductParams

	var err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	var result = repo.NewProductRepo(cc.App.DB).PaginateCollectionProduct(params)

	return cc.Success(result)
}

// CollectionDetail
// @Tags Marketplace-Collection
// @Summary Collection Detail
// @Description Collection Detail
// @Accept  json
// @Produce  json
// @Param collection_id query string true "CollectionID"
// @Success 200 {object} models.Collection
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/collections/get [get]
func CollectionDetail(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.GetCollectionDetailByIDParams

	var err = cc.BindAndValidate(&params)
	if err != nil {
		return err
	}

	result, err := repo.NewCollectionRepo(cc.App.DB).GetCollectionDetailByID(params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	return cc.Success(result)
}
