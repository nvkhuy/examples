package controllers

import (
	"github.com/engineeringinflow/inflow-backend/pkg/models"
	"github.com/engineeringinflow/inflow-backend/pkg/repo"
	"github.com/labstack/echo/v4"
	"github.com/rotisserie/eris"
)

// CreateReleaseNote create category
// @Tags Admin-Fabric
// @Summary CreateFromPayload category
// @Description CreateFromPayload category
// @Accept  json
// @Produce  json
// @Param data body models.FabricUpdateForm true "Form"
// @Success 200 {object} models.Fabric
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/release_notes [post]
func CreateReleaseNote(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.CreateReleaseNotesParams

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	params.JwtClaimsInfo = claims
	result, err := repo.NewReleaseNoteRepo(cc.App.DB).Create(params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	return cc.Success(result)
}

// UpdateReleaseNote update category
// @Tags Admin-Fabric
// @Summary CreateFromPayload category
// @Description CreateFromPayload category
// @Accept  json
// @Produce  json
// @Param data body models.FabricUpdateForm true "Form"
// @Success 200 {object} models.Fabric
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/release_notes/:id [post]
func UpdateReleaseNote(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.UpdateReleaseNotesParams

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	params.JwtClaimsInfo = claims
	result, err := repo.NewReleaseNoteRepo(cc.App.DB).Update(params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	return cc.Success(result)
}

// PaginateReleaseNote
// @Tags Admin-Blog
// @Summary PaginateBlogCategory
// @Description PaginateBlogCategory
// @Accept  json
// @Produce  json
// @Param keyword query string false "Keyword"
// @Param name query string false "Name"
// @Success 200 {object} models.BlogCategory
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/release_notes [get]
func PaginateReleaseNote(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.PaginateReleaseNotesParams

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	params.JwtClaimsInfo = claims
	result := repo.NewReleaseNoteRepo(cc.App.DB).Paginate(params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	return cc.Success(result)
}

// DeleteReleaseNote update category
// @Tags Admin-Fabric
// @Summary CreateFromPayload category
// @Description CreateFromPayload category
// @Accept  json
// @Produce  json
// @Param data body models.FabricUpdateForm true "Form"
// @Success 200 {object} models.Fabric
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/release_notes/:id [delete]
func DeleteReleaseNote(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.DeleteReleaseNotesParams

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	params.JwtClaimsInfo = claims
	err = repo.NewReleaseNoteRepo(cc.App.DB).Delete(params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	return cc.Success("deleted")
}
