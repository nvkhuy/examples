package controllers

import (
	"github.com/engineeringinflow/inflow-backend/pkg/models"
	"github.com/engineeringinflow/inflow-backend/pkg/repo"
	"github.com/labstack/echo/v4"
	"github.com/rotisserie/eris"
)

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
// @Router /api/v1/analytics/products/ [get]
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
// @Router /api/v1/analytics/products/group [get]
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
// @Router /api/v1/analytics/products/{product_id} [get]
func GetAnalyticProductDetails(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.DataAnalyticGetProductParam

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, "")
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
// @Router /api/v1/analytics/products/recommend [get]
func RecommendAnalyticProducts(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.DataAnalyticRecommendProductsParam

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, "")
	}

	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}
	params.JwtClaimsInfo = claims

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
// @Router /api/v1/analytics/products/{product_id}/chart [get]
func GetAnalyticProductChart(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.DataAnalyticGetDAProductChartParam

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, "")
	}

	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}
	params.JwtClaimsInfo = claims

	result := repo.NewDataAnalyticRepo(cc.App.AnalyticDB).GetProductChart(params)
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
// @Router /api/v1/analytics/products/product_classes/group [get]
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
// @Router /api/v1/analytics/products/trending/group [get]
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
