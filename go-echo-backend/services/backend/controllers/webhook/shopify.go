package webhook

// import (
// 	goshopify "github.com/bold-commerce/go-shopify/v3"
// 	"github.com/engineeringinflow/inflow-backend/pkg/models"
// 	"github.com/engineeringinflow/inflow-backend/pkg/repo"
// 	"github.com/labstack/echo/v4"
// )

// // ShopifyProductCreate Shopify product create
// // @Tags Webhook
// // @Summary Shopify product create
// // @Description Shopify product create
// // @Accept  json
// // @Produce  json
// // @Param keyword query string false "Keyword" default(1)
// // @Param page query int false "Page index" default(1)
// // @Param limit query int false "Size of page" default(20)
// // @Header 200 {string} Bearer YOUR_TOKEN
// // @Security ApiKeyAuth
// // @Failure 404 {object} errs.Error
// // @Router /api/v1/webhook/shopify/{shop_name}/product/create [post]
// func ShopifyProductCreate(c echo.Context) error {
// 	var cc = c.(*models.CustomContext)

// 	var product goshopify.Product
// 	var err = cc.Bind(&product)
// 	if err != nil {
// 		return err
// 	}

// 	_, err = repo.NewShopifyMappingRepo(cc.App.DB).UpdateOrCreateProduct(repo.UpdateOrCreateProductParams{
// 		ShopName:       cc.GetPathParamString("shop_name"),
// 		ShopifyProduct: &product,
// 	})
// 	if err != nil {
// 		return err
// 	}

// 	return cc.Success("")
// }

// // ShopifyProductUpdate Shopify product update
// // @Tags Webhook
// // @Summary Shopify product update
// // @Description Shopify product update
// // @Accept  json
// // @Produce  json
// // @Param keyword query string false "Keyword" default(1)
// // @Param page query int false "Page index" default(1)
// // @Param limit query int false "Size of page" default(20)
// // @Header 200 {string} Bearer YOUR_TOKEN
// // @Security ApiKeyAuth
// // @Failure 404 {object} errs.Error
// // @Router /api/v1/webhook/shopify/{shop_name}/product/update [post]
// func ShopifyProductUpdate(c echo.Context) error {
// 	var cc = c.(*models.CustomContext)

// 	var product goshopify.Product
// 	var err = cc.Bind(&product)
// 	if err != nil {
// 		return err
// 	}

// 	_, err = repo.NewShopifyMappingRepo(cc.App.DB).UpdateOrCreateProduct(repo.UpdateOrCreateProductParams{
// 		ShopName:       cc.GetPathParamString("shop_name"),
// 		ShopifyProduct: &product,
// 	})
// 	if err != nil {
// 		return err
// 	}

// 	return cc.Success("")
// }

// // ShopifyShopUpdate Shopify shop update
// // @Tags Webhook
// // @Summary Shopify shop update
// // @Description Shopify shop update
// // @Accept  json
// // @Produce  json
// // @Param keyword query string false "Keyword" default(1)
// // @Param page query int false "Page index" default(1)
// // @Param limit query int false "Size of page" default(20)
// // @Header 200 {string} Bearer YOUR_TOKEN
// // @Security ApiKeyAuth
// // @Failure 404 {object} errs.Error
// // @Router /api/v1/webhook/shopify/{shop_name}/shop/update [post]
// func ShopifyShopUpdate(c echo.Context) error {
// 	var cc = c.(*models.CustomContext)

// 	return cc.Success("")
// }

// // ShopifyInventoryItemsCreate Shopify inventory item create
// // @Tags Webhook
// // @Summary Shopify inventory item create
// // @Description Shopify inventory item create
// // @Accept  json
// // @Produce  json
// // @Param keyword query string false "Keyword" default(1)
// // @Param page query int false "Page index" default(1)
// // @Param limit query int false "Size of page" default(20)
// // @Header 200 {string} Bearer YOUR_TOKEN
// // @Security ApiKeyAuth
// // @Failure 404 {object} errs.Error
// // @Router /api/v1/webhook/shopify/{shop_name}/inventory_items/create [post]
// func ShopifyInventoryItemsCreate(c echo.Context) error {
// 	var cc = c.(*models.CustomContext)

// 	return cc.Success("")
// }

// // ShopifyInventoryItemsCreate Shopify inventory item update
// // @Tags Webhook
// // @Summary Shopify inventory item update
// // @Description Shopify inventory item update
// // @Accept  json
// // @Produce  json
// // @Param keyword query string false "Keyword" default(1)
// // @Param page query int false "Page index" default(1)
// // @Param limit query int false "Size of page" default(20)
// // @Header 200 {string} Bearer YOUR_TOKEN
// // @Security ApiKeyAuth
// // @Failure 404 {object} errs.Error
// // @Router /api/v1/webhook/shopify/{shop_name}/inventory_items/update [post]
// func ShopifyInventoryItemsUpdate(c echo.Context) error {
// 	var cc = c.(*models.CustomContext)

// 	return cc.Success("")
// }

// // ShopifyInventoryItemsCreate Shopify inventory item delete
// // @Tags Webhook
// // @Summary Shopify inventory item delete
// // @Description Shopify inventory item delete
// // @Accept  json
// // @Produce  json
// // @Param keyword query string false "Keyword" default(1)
// // @Param page query int false "Page index" default(1)
// // @Param limit query int false "Size of page" default(20)
// // @Header 200 {string} Bearer YOUR_TOKEN
// // @Security ApiKeyAuth
// // @Failure 404 {object} errs.Error
// // @Router /api/v1/webhook/shopify/{shop_name}/inventory_items/delete [post]
// func ShopifyInventoryItemsDelete(c echo.Context) error {
// 	var cc = c.(*models.CustomContext)

// 	return cc.Success("")
// }
