package controllers

import (
	"github.com/engineeringinflow/inflow-backend/pkg/models"
	"github.com/engineeringinflow/inflow-backend/pkg/models/enums"
	"github.com/engineeringinflow/inflow-backend/pkg/repo"
	"github.com/labstack/echo/v4"
	"github.com/rotisserie/eris"
)

// PaginateTrendings
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
// @Router /api/v1/buyer/trendings [get]
func PaginateTrendings(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.PaginateTrendingParams

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}
	params.Statuses = []enums.TrendingStatus{enums.TrendingStatusPublished}

	params.JwtClaimsInfo = claims
	result := repo.NewTrendingRepo(cc.App.DB).
		WithADB(cc.App.AnalyticDB).PaginateTrendings(params)

	return cc.Success(result)
}

// GetTrending
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
// @Router /api/v1/buyer/trendings/:id [get]
func GetTrending(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.GetTrendingParams

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	params.JwtClaimsInfo = claims
	result, err := repo.NewTrendingRepo(cc.App.DB).WithADB(cc.App.AnalyticDB).Get(params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	return cc.Success(result)
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
// @Router /api/v1/buyer/analytics/product_trendings/:product_trending_id [get]
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
// @Router /api/v1/buyer/analytics/product_trendings/{product_trending_id}/chart [get]
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

// GetProductTrendingTags
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
// @Router /api/v1/buyer/analytics/product_trendings/tags [get]
func GetProductTrendingTags(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.PaginateAnalyticGrowingTagsParams

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	params.JwtClaimsInfo = claims
	result, err := repo.NewAnalyticGrowingTagRepo(cc.App.AnalyticDB).Paginate(params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}
	return cc.Success(result)
}
