package controllers

import (
	"github.com/engineeringinflow/inflow-backend/pkg/models"
	"github.com/engineeringinflow/inflow-backend/pkg/repo"
	"github.com/labstack/echo/v4"
	"github.com/rotisserie/eris"
)

// ResourcesManagement
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
// @Router /api/v1/admin/resources [get]
func ResourcesManagement(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	_, err := cc.GetJwtClaims()
	if err != nil {
		return eris.Wrap(err, "")
	}

	var params repo.PaginateUsersRoleParams
	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, "")
	}

	result, err := repo.NewUserRepo(cc.App.DB).PaginateUsersRole(params)
	if err != nil {
		return eris.Wrap(err, "")
	}
	return cc.Success(result)
}

// ReloadResourcesPolicy
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
// @Router /api/v1/admin/resources/policy [patch]
func ReloadResourcesPolicy(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	_, err := cc.GetJwtClaims()
	if err != nil {
		return eris.Wrap(err, "")
	}
	err = cc.App.Enforcer.LoadPolicy()
	if err != nil {
		return eris.Wrap(err, "LoadPolicy Error")
	}
	return cc.Success("done")
}

// GetResourcesPolicy
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
// @Router /api/v1/admin/resources/policy [get]
func GetResourcesPolicy(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	_, err := cc.GetJwtClaims()
	if err != nil {
		return eris.Wrap(err, "")
	}
	all := cc.App.Enforcer.GetAllNamedSubjects("p")
	result := models.ListResources(cc.App.Enforcer, all...)
	return cc.Success(result)
}
