package controllers

import (
	"github.com/engineeringinflow/inflow-backend/pkg/models"
	"github.com/engineeringinflow/inflow-backend/pkg/repo"
	"github.com/engineeringinflow/inflow-backend/services/consumer/tasks"
	"github.com/labstack/echo/v4"
	"github.com/rotisserie/eris"
)

// GetCategoryTree
// @Tags Marketplace-Product
// @Summary GetCategoryTree
// @Description GetCategoryTree
// @Accept  json
// @Produce  json
// @Success 200 {object} models.User
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/buyer/products/get_category_tree [get]
func GetCategoryTree(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.PaginateCategoriesParams
	var err = cc.Bind(&params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	var result = repo.NewCategoryRepo(cc.App.DB).GetCategoryTree(params)

	return cc.Success(result)
}

// GetProductCategories
// @Tags Marketplace-Product
// @Summary GetCategoryTree
// @Description GetCategoryTree
// @Accept  json
// @Produce  json
// @Success 200 {object} models.User
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/buyer/products/categories [get]
func GetProductCategories(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.PaginateCategoriesParams
	var err = cc.Bind(&params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	var result = repo.NewCategoryRepo(cc.App.DB).GetCategories(params)

	return cc.Success(result)
}

// PaginateProductCategoriesChildren
// @Tags Marketplace-Product
// @Summary PaginateProductCategoriesChildren
// @Description PaginateProductCategoriesChildren
// @Accept  json
// @Produce  json
// @Success 200 {object} models.User
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/buyer/products/categories/{parent_category_id}/children [get]
func PaginateProductCategoriesChildren(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.PaginateCategoriesParams
	var err = cc.Bind(&params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	var result = repo.NewCategoryRepo(cc.App.DB).PaginateCategoriesParams(params)

	return cc.Success(result)
}

// CreateCatalogCart
// @Tags Marketplace-Product
// @Summary CreateProductCart
// @Description CreateProductCart
// @Accept  json
// @Produce  json
// @Success 200 {object} models.User
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/buyer/products/{product_id}/cart [post]
func CreateCatalogCart(c echo.Context) error {
	var cc = c.(*models.CustomContext)

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	var form repo.CreateCatalogCartForm
	err = cc.BindAndValidate(&form)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	form.JwtClaimsInfo = claims
	result, err := repo.NewCatalogCartRepo(cc.App.DB).CreateCatalogCart(form)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}
	_, _ = tasks.CreateChatRoomTask{
		UserID:          claims.GetUserID(),
		Role:            claims.GetRole(),
		PurchaseOrderID: result.ID,
		BuyerID:         result.UserID,
	}.Dispatch(c.Request().Context())

	return cc.Success(result)
}
