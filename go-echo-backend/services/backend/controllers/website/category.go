package controllers

import (
	"github.com/engineeringinflow/inflow-backend/pkg/models"
	"github.com/engineeringinflow/inflow-backend/pkg/repo"
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
// @Router /api/v1/products/get_category_tree [get]
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
