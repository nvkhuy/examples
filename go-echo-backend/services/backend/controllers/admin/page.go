package controllers

import (
	"github.com/engineeringinflow/inflow-backend/pkg/models"
	"github.com/engineeringinflow/inflow-backend/pkg/models/enums"
	"github.com/engineeringinflow/inflow-backend/pkg/repo"
	"github.com/labstack/echo/v4"

	"github.com/rotisserie/eris"
)

// PageList
// @Tags Admin-Page
// @Summary Page List
// @Description Page List
// @Accept  json
// @Produce  json
// @Success 200 {object} models.Page
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/pages [get]
func AdminPageList(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.SearchPageParams
	var err = cc.Bind(&params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	params.ForRole = enums.RoleSuperAdmin
	var result = repo.NewPageRepo(cc.App.DB).ListPage(params)
	return cc.Success(result)
}

// AdminPageDetail
// @Tags Admin-Page
// @Summary Page List
// @Description Page List
// @Accept  json
// @Produce  json
// @Success 200 {object} models.Page
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/pages/{page_id} [get]
func AdminPageDetail(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.GetPageDetailByIDParams

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	params.JwtClaimsInfo = claims
	result, err := repo.NewPageRepo(cc.App.DB).GetPageDetailByID(params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}
	return cc.Success(result)
}

// CreatePage create Page
// @Tags Admin-Page
// @Summary CreateFromPayload Page
// @Description CreateFromPayload Page
// @Accept  json
// @Produce  json
// @Param data body models.PageUpdateForm true "Form"
// @Success 200 {object} models.Page
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/pages [post]
func AdminCreatePage(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var form models.PageUpdateForm

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	err = cc.BindAndValidate(&form)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	form.JwtClaimsInfo = claims
	cate, err := repo.NewPageRepo(cc.App.DB).CreatePage(form)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	return cc.Success(cate)
}

// AdminUpdatePage
// @Tags Admin-Page
// @Summary Update Page
// @Description Update Page
// @Accept  json
// @Produce  json
// @Param user_id path string true "ID"
// @Param data body models.PageWithSectionUpdateForm true "Form"
// @Success 200 {object} models.Page
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/pages/{page_id} [put]
func AdminUpdatePage(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var id = cc.GetPathParamString("page_id")
	var form models.PageWithSectionUpdateForm

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	err = cc.BindAndValidate(&form)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	form.JwtClaimsInfo = claims
	cate, err := repo.NewPageRepo(cc.App.DB).UpdatePageByID(id, form)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	return cc.Success(cate)
}

// AddPageSection
// @Tags Admin-Page
// @Summary Add PageSection
// @Description Add PageSection
// @Accept  json
// @Produce  json
// @Param id path string true "ID"
// @Param data body models.PageSectionCreateForm true "Form"
// @Success 200 {object} models.Page
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/pages/{page_id}/add_section [post]
func AdminAddPageSection(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var id = cc.GetPathParamString("id")
	var form models.PageSectionCreateForm

	err := cc.BindAndValidate(&form)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	form.ForRole = enums.RoleSuperAdmin
	form.PageID = id
	cate, err := repo.NewPageSectionRepo(cc.App.DB).CreatePageSection(form)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	return cc.Success(cate)
}

// AdminUpdatePageSection
// @Tags Admin-Page
// @Summary Add PageSection
// @Description Add PageSection
// @Accept  json
// @Produce  json
// @Param id path string true "ID"
// @Param data body models.PageSectionUpdateForm true "Form"
// @Success 200 {object} models.Page
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/pages/sections/{section_id}/update [put]
func AdminUpdatePageSection(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var form models.PageSectionUpdateForm
	var id = cc.GetPathParamString("id")

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	err = cc.BindAndValidate(&form)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	form.JwtClaimsInfo = claims
	cate, err := repo.NewPageSectionRepo(cc.App.DB).UpdatePageSectionByID(id, form)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	return cc.Success(cate)
}

// AdminDeletePageSection
// @Tags Admin-Page
// @Summary Delete PageSection
// @Description Delete PageSection
// @Accept  json
// @Produce  json
// @Param id path string true "ID"
// @Param data body models.PageSectionUpdateForm true "Form"
// @Success 200 {object} models.M
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/pages/sections/{section_id}/delete [delete]
func AdminDeletePageSection(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var id = cc.GetPathParamString("id")

	err := repo.NewPageSectionRepo(cc.App.DB).DeletePageSectionByID(id)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	return cc.Success("Archived")
}
