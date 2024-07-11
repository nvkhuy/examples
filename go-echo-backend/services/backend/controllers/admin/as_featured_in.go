package controllers

import (
	"github.com/engineeringinflow/inflow-backend/pkg/models"
	"github.com/engineeringinflow/inflow-backend/pkg/repo"
	"github.com/labstack/echo/v4"

	"github.com/rotisserie/eris"
)

// PaginateAsFeaturedIns
// @Tags Admin-Ads
// @Summary PaginateAsFeaturedIns
// @Description PaginateAsFeaturedIns
// @Accept  json
// @Produce  json
// @Param keyword query string false "Keyword"
// @Param name query string false "Name"
// @Success 200 {object} models.AsFeaturedIn
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/as_featured_ins [get]
func PaginateAsFeaturedIns(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.PaginateAsFeaturedInParams

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	params.JwtClaimsInfo = claims
	var result = repo.NewAsFeaturedInRepo(cc.App.DB).PaginateAsFeaturedIn(params)
	return cc.Success(result)
}

// CreateAsFeaturedIn CreateFromPayload ads video
// @Tags Admin-Ads
// @Summary CreateFromPayload ads video
// @Description CreateFromPayload ads video
// @Accept  json
// @Produce  json
// @Param data body models.AsFeaturedIn true "Form"
// @Success 200 {object} models.AsFeaturedIn
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/as_featured_ins [post]
func CreateAsFeaturedIn(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var form models.AsFeaturedInCreateForm

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	err = cc.BindAndValidate(&form)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	form.JwtClaimsInfo = claims
	result, err := repo.NewAsFeaturedInRepo(cc.App.DB).CreateAsFeaturedIn(form)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	return cc.Success(result)
}

// UpdateAsFeaturedIn Update ads video
// @Tags Admin-Blog
// @Summary Update ads video
// @Description Update ads video
// @Accept  json
// @Produce  json
// @Param blog_category_id path string true "ID"
// @Param data body models.Ad true "Form"
// @Success 200 {object} models.AsFeaturedInUpdateForm
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/as_featured_ins/{as_featured_in_id} [put]
func UpdateAsFeaturedIn(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var form models.AsFeaturedInUpdateForm

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	err = cc.BindAndValidate(&form)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	form.JwtClaimsInfo = claims
	Post, err := repo.NewAsFeaturedInRepo(cc.App.DB).UpdateAsFeaturedIn(form)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	return cc.Success(Post)
}

// DeleteAsFeaturedIn
// @Tags Admin-Blog
// @Summary Delete category
// @Description Delete category
// @Accept  json
// @Produce  json
// @Param blog_category_id path string true "ID"
// @Success 200 {object} models.M
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/as_featured_ins/{as_featured_in_id} [delete]
func DeleteAsFeaturedIn(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.DeleteAsFeaturedInParams

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	params.JwtClaimsInfo = claims
	err = repo.NewAsFeaturedInRepo(cc.App.DB).DeleteAsFeaturedIn(params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	return cc.Success("Archived")
}
