package controllers

import (
	"github.com/engineeringinflow/inflow-backend/pkg/models"
	"github.com/engineeringinflow/inflow-backend/pkg/models/enums"
	"github.com/engineeringinflow/inflow-backend/pkg/repo"
	"github.com/labstack/echo/v4"
	"github.com/rotisserie/eris"
)

// GetSEOTranslations
// @Summary Paginate setting bank
// @Description Paginate setting bank
// @Accept  json
// @Produce  json
// @Param keyword query string false "Keyword"
// @Param page query int false "Page number"
// @Param limit query int false "Limit"
// @Success 200 {object} query.Pagination{records=[]models.SettingBank}
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/common/seo/translations[get]
func GetSEOTranslations(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var err error

	var params models.GetSEOTranslationForm
	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, "")
	}

	params.Domain = enums.DomainWebsite
	result, err := repo.NewSeoTranslationRepo(cc.App.DB).GetSEOTranslation(params)
	if err != nil {
		return eris.Wrap(err, "")
	}
	return cc.Success(result)
}

// GetBuyerSEOTranslations
// @Summary Paginate setting bank
// @Description Paginate setting bank
// @Accept  json
// @Produce  json
// @Param keyword query string false "Keyword"
// @Param page query int false "Page number"
// @Param limit query int false "Limit"
// @Success 200 {object} query.Pagination{records=[]models.SettingBank}
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/common/buyer/seo/translations[get]
func GetBuyerSEOTranslations(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var err error

	var params models.GetSEOTranslationForm
	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, "")
	}

	params.Domain = enums.DomainBuyer
	result, err := repo.NewSeoTranslationRepo(cc.App.DB).GetSEOTranslation(params)
	if err != nil {
		return eris.Wrap(err, "")
	}
	return cc.Success(result)
}

// GetSellerSEOTranslations
// @Summary Paginate setting bank
// @Description Paginate setting bank
// @Accept  json
// @Produce  json
// @Param keyword query string false "Keyword"
// @Param page query int false "Page number"
// @Param limit query int false "Limit"
// @Success 200 {object} query.Pagination{records=[]models.SettingBank}
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/common/seller/seo/translations[get]
func GetSellerSEOTranslations(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var err error

	var params models.GetSEOTranslationForm
	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, "")
	}

	params.Domain = enums.DomainSeller
	result, err := repo.NewSeoTranslationRepo(cc.App.DB).GetSEOTranslation(params)
	if err != nil {
		return eris.Wrap(err, "")
	}
	return cc.Success(result)
}

// GetAdminSEOTranslations
// @Summary Paginate setting bank
// @Description Paginate setting bank
// @Accept  json
// @Produce  json
// @Param keyword query string false "Keyword"
// @Param page query int false "Page number"
// @Param limit query int false "Limit"
// @Success 200 {object} query.Pagination{records=[]models.SettingBank}
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/common/admin/seo/translations[get]
func GetAdminSEOTranslations(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var err error

	var params models.GetSEOTranslationForm
	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, "")
	}

	params.Domain = enums.DomainAdmin
	result, err := repo.NewSeoTranslationRepo(cc.App.DB).GetSEOTranslation(params)
	if err != nil {
		return eris.Wrap(err, "")
	}
	return cc.Success(result)
}

// GetSEOTranslations
// @Tags Admin-PO
// @Summary Paginate setting bank
// @Description Paginate setting bank
// @Accept  json
// @Produce  json
// @Param keyword query string false "Keyword"
// @Param page query int false "Page number"
// @Param limit query int false "Limit"
// @Success 200 {object} query.Pagination{records=[]models.SettingBank}
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /settings/seo/route/{route_name} [get]
func GetSettingSEOByRouteName(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var err error

	var params repo.GetSettingSEOByRouteNameParams
	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, "")
	}

	result, err := repo.NewSettingSEORepo(cc.App.DB).GetSettingSEOByRouteName(params)
	if err != nil {
		return eris.Wrap(err, "")
	}
	return cc.Success(result)
}
