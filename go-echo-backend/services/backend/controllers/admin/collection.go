package controllers

import (
	"fmt"

	"github.com/engineeringinflow/inflow-backend/pkg/models"
	"github.com/engineeringinflow/inflow-backend/pkg/models/enums"
	"github.com/engineeringinflow/inflow-backend/pkg/repo"
	"github.com/engineeringinflow/inflow-backend/services/consumer/tasks"
	"github.com/labstack/echo/v4"

	"github.com/rotisserie/eris"
)

// AdminCollectionList
// @Tags Admin-Collection
// @Summary Collection List
// @Description Collection List
// @Accept  json
// @Produce  json
// @Param keyword query string false "Keyword"
// @Param page query int false "Page number"
// @Success 200 {object} models.Collection
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/collections [get]
func AdminCollectionList(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.PaginateCollectionParams

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	params.JwtClaimsInfo = claims
	var result = repo.NewCollectionRepo(cc.App.DB).PaginateCollection(params)
	return cc.Success(result)
}

// AdminCreateCollection
// @Tags Admin-Collection
// @Summary CreateFromPayload Collection
// @Description CreateFromPayload Collection
// @Accept  json
// @Produce  json
// @Param data body models.CollectionUpdateForm true "Form"
// @Success 200 {object} models.Collection
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/collections/create [post]
func AdminCreateCollection(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var form models.CollectionCreateForm

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	err = cc.BindAndValidate(&form)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	form.JwtClaimsInfo = claims
	collection, err := repo.NewCollectionRepo(cc.App.DB).CreateCollection(form)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	if collection != nil {
		_, _ = tasks.CreateSysNotificationTask{
			SysNotification: models.SysNotification{
				Name:    fmt.Sprintf("New Collection - %s", collection.Name),
				Type:    enums.SysNotificationCreateCollectionType,
				Message: fmt.Sprintf("New Collection - %s", collection.Name),
			},
		}.Dispatch(c.Request().Context())
	}

	return cc.Success(collection)
}

// AdminUpdateCollection
// @Tags Admin-Collection
// @Summary Update Collection
// @Description Update Collection
// @Accept  json
// @Produce  json
// @Param id path string true "ID"
// @Param data body models.CollectionUpdateForm true "Form"
// @Success 200 {object} models.Collection
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/collections/{collection_id} [put]
func AdminUpdateCollection(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var form models.CollectionUpdateForm

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	err = cc.BindAndValidate(&form)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	form.JwtClaimsInfo = claims
	Collection, err := repo.NewCollectionRepo(cc.App.DB).UpdateCollectionByID(form)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	return cc.Success(Collection)
}

// CollectionDetail
// @Tags Admin-Collection
// @Summary Collection Detail
// @Description Collection Detail
// @Accept  json
// @Produce  json
// @Param collection_id query string true "CollectionID"
// @Success 200 {object} models.Post
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/collections/get [get]
func AdminCollectionDetail(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.GetCollectionDetailByIDParams

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	params.JwtClaimsInfo = claims
	result, err := repo.NewCollectionRepo(cc.App.DB).GetCollectionDetailByID(params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	return cc.Success(result)
}

// AdminDeleteCollection
// @Tags Admin-Collection
// @Summary Delete Collection
// @Description Delete Collection
// @Accept  json
// @Produce  json
// @Param id path string true "ID"
// @Success 200 {object} models.M
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/collections/{collection_id}/delete [delete]
func AdminDeleteCollection(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var id = cc.GetPathParamString("collection_id")

	err := repo.NewCollectionRepo(cc.App.DB).ArchiveCollectionByID(id)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	return cc.Success("Archived")
}

// Collection Add Product
// @Tags Admin-Collection
// @Summary Collection Add Product
// @Description Collection Add Product
// @Accept  json
// @Produce  json
// @Param data body models.CollectionProductIDsForm true "Form"
// @Success 200 {object} models.Collection
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/collections/{collection_id}/add_product [post]
func AdminCollectionAddProduct(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var id = cc.GetPathParamString("collection_id")
	var form models.CollectionProductIDsForm

	err := cc.BindAndValidate(&form)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	collection, err := repo.NewCollectionRepo(cc.App.DB).AddProduct(id, form.Products)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	return cc.Success(collection)
}

// AdminCollection Delete Product
// @Tags Admin-Collection
// @Summary Collection Delete Product
// @Description Collection Delete Product
// @Accept  json
// @Produce  json
// @Param id path string true "ID"
// @Param data body models.CollectionProductIDsForm true "Form"
// @Success 200 {object} models.M
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/collections/{collection_id}/delete_product [delete]
func AdminCollectionDeleteProduct(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var id = cc.GetPathParamString("collection_id")

	var form models.CollectionProductIDsForm

	err := cc.BindAndValidate(&form)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	err = repo.NewCollectionRepo(cc.App.DB).RemoveProduct(id, form.Products)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	return cc.Success("Deleted")
}

// CollectionGetProduct
// @Tags Admin-Collection
// @Summary Collection GetProduct
// @Description Collection GetProduct
// @Accept  json
// @Produce  json
// @Param id path string true "CollectionID"
// @Success 200 {object} models.Product
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/collections/{id}/get_product [get]
func AdminCollectionGetProduct(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.PaginateCollectionProductParams

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	params.JwtClaimsInfo = claims

	var result = repo.NewProductRepo(cc.App.DB).PaginateCollectionProduct(params)

	return cc.Success(result)
}
