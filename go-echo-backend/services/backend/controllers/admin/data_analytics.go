package controllers

import (
	"github.com/engineeringinflow/inflow-backend/pkg/models"
	"github.com/engineeringinflow/inflow-backend/pkg/repo"
	"github.com/labstack/echo/v4"
	"github.com/rotisserie/eris"
)

// OverviewDataAnalyticPlatform
// @Tags Admin-Blog
// @Summary PaginateAIProducts
// @Description PaginateAIProducts
// @Accept  json
// @Produce  json
// @Param keyword query string false "Keyword"
// @Param name query string false "Name"
// @Success 200 {object} models.BlogCategory
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/analytics/platform/overview [get]
func OverviewDataAnalyticPlatform(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.DataAnalyticPlatformOverviewParam

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	params.JwtClaimsInfo = claims
	result := repo.NewDataAnalyticRepo(cc.App.AnalyticDB).Overview(params)

	return cc.Success(result)
}

// SearchDAProducts
// @Tags Admin-Blog
// @Summary PaginateAIProducts
// @Description PaginateAIProducts
// @Accept  json
// @Produce  json
// @Param keyword query string false "Keyword"
// @Param name query string false "Name"
// @Success 200 {object} models.BlogCategory
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/analytics/products/search [get]
func SearchDAProducts(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.DataAnalyticSearchProductsParam

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	params.JwtClaimsInfo = claims
	result := repo.NewDataAnalyticRepo(cc.App.AnalyticDB).SearchProducts(params)
	return cc.Success(result)
}

// GetDAProductDetails
// @Tags Admin-Blog
// @Summary PaginateAIProducts
// @Description PaginateAIProducts
// @Accept  json
// @Produce  json
// @Param keyword query string false "Keyword"
// @Param name query string false "Name"
// @Success 200 {object} models.BlogCategory
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/analytics/products/{product_id} [get]
func GetDAProductDetails(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.DataAnalyticGetProductParam

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	params.JwtClaimsInfo = claims
	result, err := repo.NewDataAnalyticRepo(cc.App.AnalyticDB).GetProduct(params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}
	return cc.Success(result)
}

// GetDAProductChart
// @Tags Admin-Blog
// @Summary PaginateAIProducts
// @Description PaginateAIProducts
// @Accept  json
// @Produce  json
// @Param keyword query string false "Keyword"
// @Param name query string false "Name"
// @Success 200 {object} models.BlogCategory
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/analytics/products/{product_id}/chart [get]
func GetDAProductChart(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.DataAnalyticGetDAProductChartParam

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	params.JwtClaimsInfo = claims
	result := repo.NewDataAnalyticRepo(cc.App.AnalyticDB).GetProductChart(params)
	return cc.Success(result)
}

// TopDAProducts
// @Tags Admin-Blog
// @Summary PaginateAIProducts
// @Description PaginateAIProducts
// @Accept  json
// @Produce  json
// @Param keyword query string false "Keyword"
// @Param name query string false "Name"
// @Success 200 {object} models.BlogCategory
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/analytics/products/top [get]
func TopDAProducts(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.DataAnalyticTopProductsParam

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	params.JwtClaimsInfo = claims
	result := repo.NewDataAnalyticRepo(cc.App.AnalyticDB).TopProducts(params)
	return cc.Success(result)
}

// TopMoversDAProducts
// @Tags Admin-Blog
// @Summary PaginateAIProducts
// @Description PaginateAIProducts
// @Accept  json
// @Produce  json
// @Param keyword query string false "Keyword"
// @Param name query string false "Name"
// @Success 200 {object} models.BlogCategory
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/analytics/products/top [get]
func TopMoversDAProducts(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.DataAnalyticTopMoverProductsParam

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	params.JwtClaimsInfo = claims
	result := repo.NewDataAnalyticRepo(cc.App.AnalyticDB).TopMoverProducts(params)
	return cc.Success(result)
}

// TopDACategories
// @Tags Admin-Blog
// @Summary PaginateAIProducts
// @Description PaginateAIProducts
// @Accept  json
// @Produce  json
// @Param keyword query string false "Keyword"
// @Param name query string false "Name"
// @Success 200 {object} models.BlogCategory
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/analytics/categories/top [get]
func TopDACategories(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.DataAnalyticTopCategoriesParam

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	params.JwtClaimsInfo = claims
	params.Select = []string{"category"}
	result := repo.NewDataAnalyticRepo(cc.App.AnalyticDB).TopCategories(params)
	return cc.Success(result)
}

// TopDASubCategories
// @Tags Admin-Blog
// @Summary PaginateAIProducts
// @Description PaginateAIProducts
// @Accept  json
// @Produce  json
// @Param keyword query string false "Keyword"
// @Param name query string false "Name"
// @Success 200 {object} models.BlogCategory
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/analytics/sub_categories/top [get]
func TopDASubCategories(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.DataAnalyticTopCategoriesParam

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	params.JwtClaimsInfo = claims
	params.Select = []string{"sub_category"}
	result := repo.NewDataAnalyticRepo(cc.App.AnalyticDB).TopCategories(params)
	return cc.Success(result)
}

// DANewUsers
// @Tags Admin-Blog
// @Summary PaginateAIProducts
// @Description PaginateAIProducts
// @Accept  json
// @Produce  json
// @Param keyword query string false "Keyword"
// @Param name query string false "Name"
// @Success 200 {object} models.BlogCategory
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/analytics/users/new [get]
func DANewUsers(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.DataAnalyticNewUsersParam

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	params.JwtClaimsInfo = claims
	result, err := repo.NewDataAnalyticRepo(cc.App.AnalyticDB).WithDB(cc.App.DB).PaginateNewUsers(params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}
	return cc.Success(result)
}

// DANewCatalogProducts
// @Tags Admin-Blog
// @Summary PaginateAIProducts
// @Description PaginateAIProducts
// @Accept  json
// @Produce  json
// @Param keyword query string false "Keyword"
// @Param name query string false "Name"
// @Success 200 {object} models.BlogCategory
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/analytics/catalog_products/new [get]
func DANewCatalogProducts(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.DataAnalyticNewCatalogProductParam

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	params.JwtClaimsInfo = claims
	result, err := repo.NewDataAnalyticRepo(cc.App.AnalyticDB).WithDB(cc.App.DB).PaginateNewCatalogProduct(params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}
	return cc.Success(result)
}

// DAInquiries
// @Tags Admin-Blog
// @Summary PaginateAIProducts
// @Description PaginateAIProducts
// @Accept  json
// @Produce  json
// @Param keyword query string false "Keyword"
// @Param name query string false "Name"
// @Success 200 {object} models.BlogCategory
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/analytics/inquiries/new [get]
func DAInquiries(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.DataAnalyticInquiriesParam

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	params.JwtClaimsInfo = claims
	result, err := repo.NewDataAnalyticRepo(cc.App.AnalyticDB).WithDB(cc.App.DB).PaginateInquiries(params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}
	return cc.Success(result)
}

// DANewPurchaseOrders
// @Tags Admin-Blog
// @Summary PaginateAIProducts
// @Description PaginateAIProducts
// @Accept  json
// @Produce  json
// @Param keyword query string false "Keyword"
// @Param name query string false "Name"
// @Success 200 {object} models.BlogCategory
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/analytics/purchase_orders/new [get]
func DANewPurchaseOrders(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.DataAnalyticPurchaseOrdersParam

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	params.JwtClaimsInfo = claims
	result, err := repo.NewDataAnalyticRepo(cc.App.AnalyticDB).WithDB(cc.App.DB).PaginatePurchaseOrders(params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}
	return cc.Success(result)
}

// DANewBulkPurchaseOrders
// @Tags Admin-Blog
// @Summary PaginateAIProducts
// @Description PaginateAIProducts
// @Accept  json
// @Produce  json
// @Param keyword query string false "Keyword"
// @Param name query string false "Name"
// @Success 200 {object} models.BlogCategory
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/analytics/bulk_purchase_orders/new [get]
func DANewBulkPurchaseOrders(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.DataAnalyticBulkPurchaseOrdersParam

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	params.JwtClaimsInfo = claims
	result, err := repo.NewDataAnalyticRepo(cc.App.AnalyticDB).WithDB(cc.App.DB).PaginateBulkPurchaseOrders(params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}
	return cc.Success(result)
}

// DAOpsBizPerformance
// @Tags Admin-Blog
// @Summary PaginateAIProducts
// @Description PaginateAIProducts
// @Accept  json
// @Produce  json
// @Param keyword query string false "Keyword"
// @Param name query string false "Name"
// @Success 200 {object} models.BlogCategory
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/analytics/ops_biz_perfomance [get]
func DAOpsBizPerformance(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.DataAnalyticOpsBizPerformance

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	params.JwtClaimsInfo = claims
	result, err := repo.NewDataAnalyticRepo(cc.App.AnalyticDB).WithDB(cc.App.DB).GetOpsBizPerformance(params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}
	return cc.Success(result)
}

// PaginateProductTrending
// @Tags Admin-Blog
// @Summary PaginateAIProducts
// @Description PaginateAIProducts
// @Accept  json
// @Produce  json
// @Param keyword query string false "Keyword"
// @Param name query string false "Name"
// @Success 200 {object} models.BlogCategory
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/analytics/product_trendings [get]
func PaginateProductTrending(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.PaginateProductTrendingParams

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	params.JwtClaimsInfo = claims
	result := repo.NewProductTrending(cc.App.AnalyticDB).WithDB(cc.App.DB).PaginateProductTrendings(params)

	return cc.Success(result)
}

// CreateProductTrending
// @Tags Admin-Blog
// @Summary PaginateAIProducts
// @Description PaginateAIProducts
// @Accept  json
// @Produce  json
// @Param keyword query string false "Keyword"
// @Param name query string false "Name"
// @Success 200 {object} models.BlogCategory
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/analytics/product_trendings [post]
func CreateProductTrending(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.CreateProductTrendingParams

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	params.JwtClaimsInfo = claims
	result, err := repo.NewProductTrending(cc.App.AnalyticDB).WithDB(cc.App.DB).CreateProductTrending(params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	return cc.Success(result)
}

// UpdateProductTrending
// @Tags Admin-Blog
// @Summary PaginateAIProducts
// @Description PaginateAIProducts
// @Accept  json
// @Produce  json
// @Param keyword query string false "Keyword"
// @Param name query string false "Name"
// @Success 200 {object} models.BlogCategory
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/analytics/product_trendings [put]
func UpdateProductTrending(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.UpdateProductTrendingParams

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	params.JwtClaimsInfo = claims
	result, err := repo.NewProductTrending(cc.App.AnalyticDB).WithDB(cc.App.DB).UpdateProductTrending(params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	return cc.Success(result)
}

// DeleteProductTrending
// @Tags Admin-Blog
// @Summary PaginateAIProducts
// @Description PaginateAIProducts
// @Accept  json
// @Produce  json
// @Param keyword query string false "Keyword"
// @Param name query string false "Name"
// @Success 200 {object} models.BlogCategory
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/analytics/product_trendings [delete]
func DeleteProductTrending(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.DeleteProductTrendingParams

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	params.JwtClaimsInfo = claims
	err = repo.NewProductTrending(cc.App.AnalyticDB).WithDB(cc.App.DB).DeleteProductTrending(params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	return cc.Success("done")
}

// GetProductTrending
// @Tags Admin-Blog
// @Summary PaginateAIProducts
// @Description PaginateAIProducts
// @Accept  json
// @Produce  json
// @Param keyword query string false "Keyword"
// @Param name query string false "Name"
// @Success 200 {object} models.BlogCategory
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/analytics/product_trendings/:product_trending_id [get]
func GetProductTrending(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.GetProductTrendingParam

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	params.JwtClaimsInfo = claims
	result, err := repo.NewProductTrending(cc.App.AnalyticDB).WithDB(cc.App.DB).Get(params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	return cc.Success(result)
}

// GetProductTrendingChart
// @Tags Admin-Blog
// @Summary PaginateAIProducts
// @Description PaginateAIProducts
// @Accept  json
// @Produce  json
// @Param keyword query string false "Keyword"
// @Param name query string false "Name"
// @Success 200 {object} models.BlogCategory
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/analytics/product_trendings/{product_trending_id}/chart [get]
func GetProductTrendingChart(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.PaginateProductTrendingChartParams

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	params.JwtClaimsInfo = claims
	result := repo.NewProductTrending(cc.App.AnalyticDB).Chart(params)
	return cc.Success(result)
}

// PaginateProductTrendingGroup
// @Tags Admin-Blog
// @Summary PaginateAIProducts
// @Description PaginateAIProducts
// @Accept  json
// @Produce  json
// @Param keyword query string false "Keyword"
// @Param name query string false "Name"
// @Success 200 {object} models.BlogCategory
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/analytics/product_trendings/group [get]
func PaginateProductTrendingGroup(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.PaginateProductTrendingGroupParams

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	params.JwtClaimsInfo = claims
	result := repo.NewProductTrending(cc.App.AnalyticDB).WithDB(cc.App.DB).PaginateProductTrendingGroup(params)

	return cc.Success(result)
}

// ListProductTrendingDomain
// @Tags Admin-Blog
// @Summary PaginateAIProducts
// @Description PaginateAIProducts
// @Accept  json
// @Produce  json
// @Param keyword query string false "Keyword"
// @Param name query string false "Name"
// @Success 200 {object} models.BlogCategory
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/analytics/product_trendings/domains [get]
func ListProductTrendingDomain(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.ListProductTrendingDomainParams

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	params.JwtClaimsInfo = claims
	result, err := repo.NewProductTrending(cc.App.AnalyticDB).WithDB(cc.App.DB).ListProductTrendingDomain(params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}
	return cc.Success(result)
}

// ListProductTrendingCategory
// @Tags Admin-Blog
// @Summary PaginateAIProducts
// @Description PaginateAIProducts
// @Accept  json
// @Produce  json
// @Param keyword query string false "Keyword"
// @Param name query string false "Name"
// @Success 200 {object} models.BlogCategory
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/analytics/product_trendings/categories [get]
func ListProductTrendingCategory(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.ListProductTrendingCategoryParams

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	params.JwtClaimsInfo = claims
	result, err := repo.NewProductTrending(cc.App.AnalyticDB).WithDB(cc.App.DB).ListProductTrendingCategory(params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}
	return cc.Success(result)
}

// ListProductTrendingSubCategory
// @Tags Admin-Blog
// @Summary PaginateAIProducts
// @Description PaginateAIProducts
// @Accept  json
// @Produce  json
// @Param keyword query string false "Keyword"
// @Param name query string false "Name"
// @Success 200 {object} models.BlogCategory
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/analytics/product_trendings/sub_categories [get]
func ListProductTrendingSubCategory(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.ListProductTrendingCategoryParams

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	params.JwtClaimsInfo = claims
	result, err := repo.NewProductTrending(cc.App.AnalyticDB).WithDB(cc.App.DB).ListProductTrendingSubCategory(params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}
	return cc.Success(result)
}

// PaginateAnalyticProduct
// @Tags Admin-Blog
// @Summary PaginateAIProducts
// @Description PaginateAIProducts
// @Accept  json
// @Produce  json
// @Param keyword query string false "Keyword"
// @Param name query string false "Name"
// @Success 200 {object} models.BlogCategory
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/data_analytics/products [get]
func PaginateAnalyticProduct(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.DataAnalyticSearchProductsParam

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, "")
	}

	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}
	params.JwtClaimsInfo = claims

	result := repo.NewDataAnalyticRepo(cc.App.AnalyticDB).SearchProducts(params)
	return cc.Success(result)
}

// PaginateAnalyticProductGroup
// @Tags Admin-Blog
// @Summary PaginateAIProducts
// @Description PaginateAIProducts
// @Accept  json
// @Produce  json
// @Param keyword query string false "Keyword"
// @Param name query string false "Name"
// @Success 200 {object} models.BlogCategory
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/data_analytics/products/group [get]
func PaginateAnalyticProductGroup(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.DataAnalyticSearchProductsParam

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, "")
	}

	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}
	params.JwtClaimsInfo = claims

	result := repo.NewDataAnalyticRepo(cc.App.AnalyticDB).ProductGroupURL(params)
	return cc.Success(result)
}

// GetAnalyticProductClassGroup
// @Tags Admin-Blog
// @Summary PaginateAIProducts
// @Description PaginateAIProducts
// @Accept  json
// @Produce  json
// @Param keyword query string false "Keyword"
// @Param name query string false "Name"
// @Success 200 {object} models.BlogCategory
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/data_analytics/products/product_classes/group [get]
func GetAnalyticProductClassGroup(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.GetAnalyticProductClassGroupParam

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, "")
	}

	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}
	params.JwtClaimsInfo = claims

	result := repo.NewDataAnalyticRepo(cc.App.AnalyticDB).GetProductClassGroup(params)
	return cc.Success(result)
}

// GetAnalyticProductTrendingGroup
// @Tags Admin-Blog
// @Summary PaginateAIProducts
// @Description PaginateAIProducts
// @Accept  json
// @Produce  json
// @Param keyword query string false "Keyword"
// @Param name query string false "Name"
// @Success 200 {object} models.BlogCategory
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/data_analytics/products/trending/group [get]
func GetAnalyticProductTrendingGroup(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.GetAnalyticProductTrendingGroupParam

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, "")
	}

	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}
	params.JwtClaimsInfo = claims

	result := repo.NewDataAnalyticRepo(cc.App.AnalyticDB).GetProductTrendingGroup(params)
	return cc.Success(result)
}
