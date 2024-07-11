package routes

import (
	controllers "github.com/engineeringinflow/inflow-backend/services/backend/controllers/webhook"
	"github.com/labstack/echo/v4"
)

func (router *Router) SetupWebhookRoutes(g *echo.Group) {
	// g.POST("/shopify/:shop_name/products/create", controllers.ShopifyProductCreate)
	// g.POST("/shopify/:shop_name/products/update", controllers.ShopifyProductUpdate)
	// g.POST("/shopify/:shop_name/shop/update", controllers.ShopifyShopUpdate)
	// g.POST("/shopify/:shop_name/inventory_items/create", controllers.ShopifyInventoryItemsCreate)
	// g.POST("/shopify/:shop_name/inventory_items/delete", controllers.ShopifyInventoryItemsDelete)
	// g.POST("/shopify/:shop_name/inventory_items/update", controllers.ShopifyInventoryItemsUpdate)

	g.POST("/stripe", controllers.StripeWebhook)
	g.POST("/hubspot", controllers.HubspotWebhook)

}
