package controllers

import (
	"github.com/engineeringinflow/inflow-backend/pkg/models"
	"github.com/engineeringinflow/inflow-backend/pkg/repo"
	"github.com/labstack/echo/v4"
	"github.com/rotisserie/eris"
)

// GetAnalyticProductDetails
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
// @Router /api/v1/data_analytics/products/{product_id} [get]
func GetAnalyticProductDetails(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.DataAnalyticGetProductParam

	var err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	result, err := repo.NewDataAnalyticRepo(cc.App.AnalyticDB).GetProduct(params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}
	return cc.Success(result)
}

// RecommendAnalyticProducts
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
// @Router /api/v1/data-analytics/products/recommend [get]
func RecommendAnalyticProducts(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.DataAnalyticRecommendProductsParam

	err := cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	result := repo.NewDataAnalyticRepo(cc.App.AnalyticDB).WithDB(cc.App.DB).RecommendProducts(params)
	return cc.Success(result)
}

// GetAnalyticProductChart
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
// @Router /api/v1/data-analytics/products/{product_id}/chart [get]
func GetAnalyticProductChart(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.DataAnalyticGetDAProductChartParam

	var err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	result := repo.NewDataAnalyticRepo(cc.App.AnalyticDB).GetProductChart(params)
	return cc.Success(result)
}

// GetOneAnalyticProduct
// @Tags Website-Product
// @Summary GetOneAnalyticProductChart
// @Description GetOneAnalyticProductChart
// @Accept  json
// @Produce  json
// @Success 200 {object} {product: models.AnalyticProduct,related_products:models.Products }
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/data-analytics/products/get_one [get]
func GetOneAnalyticProduct(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.GetBestAnalyticProductParams

	var err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	main, others, err := repo.NewDataAnalyticRepo(cc.App.AnalyticDB).GetBest(params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	return cc.Success(map[string]interface{}{
		"product":          main,
		"related_products": others,
	})
}
