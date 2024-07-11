package controllers

// import (
// 	"github.com/engineeringinflow/inflow-backend/pkg/models"
// 	"github.com/engineeringinflow/inflow-backend/pkg/repo"
// 	"github.com/labstack/echo/v4"
// 	"github.com/rotisserie/eris"
// )

// // AdminPaginateInventories
// // @Tags Admin-Inventories
// // @Summary Get inventories
// // @Description Get inventories
// // @Accept  json
// // @Produce  json
// // @Success 200 {object} models.Category
// // @Header 200 {string} Bearer YOUR_TOKEN
// // @Security ApiKeyAuth
// // @Failure 404 {object} errs.Error
// // @Router /api/v1/admin/inventories [get]
// func AdminPaginateInventories(c echo.Context) error {
// 	var cc = c.(*models.CustomContext)
// 	var params repo.PaginateInventoryParams
// 	var err = cc.BindAndValidate(&params)
// 	if err != nil {
// 		return eris.Wrap(err, err.Error())
// 	}

// 	var result = repo.NewInventoryRepo(cc.App.DB).PaginateInventory(params)

// 	return cc.Success(result)
// }

// // AdminInventoryRestock
// // @Tags Admin-Inventory
// // @Summary GetCategoryTree
// // @Description GetCategoryTree
// // @Accept  json
// // @Produce  json
// // @Success 200 {object} models.Category
// // @Header 200 {string} Bearer YOUR_TOKEN
// // @Security ApiKeyAuth
// // @Failure 404 {object} errs.Error
// // @Router /api/v1/admin/inventories/restock [post]
// func AdminInventoryRestock(c echo.Context) error {
// 	var cc = c.(*models.CustomContext)
// 	var params repo.RestockParams
// 	var err = cc.BindAndValidate(&params)
// 	if err != nil {
// 		return eris.Wrap(err, err.Error())
// 	}

// 	result, err := repo.NewInventoryRepo(cc.App.DB).Restock(params)
// 	if err != nil {
// 		return eris.Wrap(err, err.Error())
// 	}

// 	return cc.Success(result)
// }
