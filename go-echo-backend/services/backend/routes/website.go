package routes

import (
	controllers "github.com/engineeringinflow/inflow-backend/services/backend/controllers/website"
	"github.com/labstack/echo/v4"
)

func (router *Router) SetupWebsiteRoutes(g *echo.Group) {
	// var authorizedRSA = g.Group("", router.Middlewares.IsValidSignature())
	var authorizedRSA = g.Group("")

	authorizedRSA.POST("/subscribe", controllers.SubscribeUpdates)

	authorizedRSA.GET("/products/get_category_tree", controllers.GetCategoryTree)
	authorizedRSA.GET("/products/search", controllers.SearchProduct)
	authorizedRSA.GET("/products/get", controllers.ProductGetDetail)
	authorizedRSA.GET("/products/get_ratings", controllers.ProductGetRatings)
	authorizedRSA.GET("/products/best_selling", controllers.ProductGetBestSelling)
	authorizedRSA.GET("/products/just_for_you", controllers.ProductGetJustForYou)
	authorizedRSA.GET("/products/recommend", controllers.ProductRecommend)
	authorizedRSA.GET("/products/ready_to_ship", controllers.ProductReadyToShip)
	authorizedRSA.GET("/products/today_deals", controllers.ProductTodayDeals)
	authorizedRSA.GET("/products/get_filter", controllers.ProductGetFilter)

	authorizedRSA.POST("/support/factory_tours/create", controllers.CreateFactoryTour)

	authorizedRSA.GET("/collections", controllers.PaginateCollection)
	authorizedRSA.GET("/collections/ready_design", controllers.CollectionReadyDesign)
	authorizedRSA.GET("/collections/:collection_id/get_product", controllers.CollectionGetProduct)
	authorizedRSA.GET("/collections/:collection_id/get", controllers.CollectionDetail)

	authorizedRSA.GET("/pages/catalog", controllers.PageCatalog)
	authorizedRSA.GET("/pages/home", controllers.PageHome)

	authorizedRSA.GET("/posts", controllers.PaginatePost)
	authorizedRSA.GET("/posts/:slug", controllers.GetPost)
	authorizedRSA.GET("/posts/stats", controllers.GetPostStats)
	authorizedRSA.GET("/blog/categories/:blog_category_id", controllers.GetBlogCategory)
	authorizedRSA.GET("/blog/categories", controllers.PaginateBlogCategory)

	authorizedRSA.GET("/documents", controllers.GetDocumentList)
	authorizedRSA.GET("/documents/:slug", controllers.GetDocument)
	authorizedRSA.GET("/document_categories", controllers.GetDocumentCategoryList)

	authorizedRSA.GET("/fabric_collections/:id", controllers.DetailsFabricCollection)
	authorizedRSA.GET("/fabrics", controllers.PaginateFabric)
	authorizedRSA.GET("/fabrics/:id", controllers.DetailsFabric)

	authorizedRSA.GET("/as_featured_ins", controllers.PaginateAsFeaturedIns)

	authorizedRSA.GET("/analytics/products/:product_id", controllers.GetAnalyticProductDetails)
	authorizedRSA.GET("/analytics/products/recommend", controllers.RecommendAnalyticProducts)
	authorizedRSA.GET("/analytics/products/:product_id/chart", controllers.GetAnalyticProductChart)
	authorizedRSA.GET("/analytics/products/get_one", controllers.GetOneAnalyticProduct)
}
