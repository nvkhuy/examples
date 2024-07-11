package controllers

import (
	"github.com/engineeringinflow/inflow-backend/pkg/models"
	"github.com/engineeringinflow/inflow-backend/pkg/repo"
	"github.com/labstack/echo/v4"
	"github.com/rotisserie/eris"
)

// CreateTrending
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
// @Router /api/v1/admin/trendings [post]
func CreateTrending(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.CreateTrendingsParams

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	params.JwtClaimsInfo = claims
	result, err := repo.NewTrendingRepo(cc.App.DB).WithADB(cc.App.AnalyticDB).Create(params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	return cc.Success(result)
}

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
// @Router /api/v1/admin/trendings [post]
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

	params.JwtClaimsInfo = claims
	result := repo.NewTrendingRepo(cc.App.DB).
		WithADB(cc.App.AnalyticDB).PaginateTrendings(params)

	return cc.Success(result)
}

// ListTrendingStats
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
// @Router /api/v1/admin/trendings/dropdown [post]
func ListTrendingStats(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.ListTrendingStatsParams

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	params.JwtClaimsInfo = claims
	result := repo.NewTrendingRepo(cc.App.DB).
		WithADB(cc.App.AnalyticDB).ListTrendingStats(params)

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
// @Router /api/v1/admin/trendings/:id [get]
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

// UpdateTrending
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
// @Router /api/v1/admin/trendings [put]
func UpdateTrending(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.UpdateTrendingParams

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	params.JwtClaimsInfo = claims
	result, err := repo.NewTrendingRepo(cc.App.DB).WithADB(cc.App.AnalyticDB).Update(params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	return cc.Success(result)
}

// DeleteTrending
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
// @Router /api/v1/admin/trendings [delete]
func DeleteTrending(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.DeleteTrendingParams

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	params.JwtClaimsInfo = claims
	err = repo.NewTrendingRepo(cc.App.DB).
		WithADB(cc.App.AnalyticDB).Delete(params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	return cc.Success("Deleted")
}

// AddProductToTrending
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
// @Router /api/v1/admin/trendings/product_trendings [put]
func AddProductToTrending(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.AddProductToTrendingParams

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	params.JwtClaimsInfo = claims
	result, err := repo.NewTrendingRepo(cc.App.DB).WithADB(cc.App.AnalyticDB).AddProductToTrending(params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	return cc.Success(result)
}

// RemoveProductFromTrending
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
// @Router /api/v1/admin/trendings/product_trendings [delete]
func RemoveProductFromTrending(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.RemoveProductFromTrendingParams

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	params.JwtClaimsInfo = claims
	result, err := repo.NewTrendingRepo(cc.App.DB).WithADB(cc.App.AnalyticDB).
		RemoveProductFromTrending(params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	return cc.Success(result)
}
