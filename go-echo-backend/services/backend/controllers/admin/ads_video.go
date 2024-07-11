package controllers

import (
	"github.com/engineeringinflow/inflow-backend/pkg/models"
	"github.com/engineeringinflow/inflow-backend/pkg/repo"
	"github.com/labstack/echo/v4"

	"github.com/rotisserie/eris"
)

// PaginateAdsVideos
// @Tags Admin-Ads
// @Summary PaginateAdsVideos
// @Description PaginateAdsVideos
// @Accept  json
// @Produce  json
// @Param keyword query string false "Keyword"
// @Param name query string false "Name"
// @Success 200 {object} models.AdsVideo
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/ads_videos [get]
func PaginateAdsVideos(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.PaginateAdsVideoParams

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	params.JwtClaimsInfo = claims
	var result = repo.NewAdsVideoRepo(cc.App.DB).PaginateAdsVideo(params)
	return cc.Success(result)
}

// CreateAdsVideo CreateFromPayload ads video
// @Tags Admin-Ads
// @Summary CreateFromPayload ads video
// @Description CreateFromPayload ads video
// @Accept  json
// @Produce  json
// @Param data body models.AdsVideo true "Form"
// @Success 200 {object} models.AdsVideo
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/ads_videos [post]
func CreateAdsVideo(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var form models.AdsVideoCreateForm

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	err = cc.BindAndValidate(&form)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	form.JwtClaimsInfo = claims
	result, err := repo.NewAdsVideoRepo(cc.App.DB).CreateAdsVideo(form)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	return cc.Success(result)
}

// UpdateAdsVideo Update ads video
// @Tags Admin-Blog
// @Summary Update ads video
// @Description Update ads video
// @Accept  json
// @Produce  json
// @Param blog_category_id path string true "ID"
// @Param data body models.Ad true "Form"
// @Success 200 {object} models.AdsVideoUpdateForm
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/ads_videos/{ads_video_id} [put]
func UpdateAdsVideo(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var form models.AdsVideoUpdateForm

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	err = cc.BindAndValidate(&form)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	form.JwtClaimsInfo = claims
	Post, err := repo.NewAdsVideoRepo(cc.App.DB).UpdateAdsVideo(form)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	return cc.Success(Post)
}

// DeleteAdsVideo
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
// @Router /api/v1/admin/ads_videos/{ads_video_id} [delete]
func DeleteAdsVideo(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.DeleteAdsVideoParams

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	params.JwtClaimsInfo = claims
	err = repo.NewAdsVideoRepo(cc.App.DB).DeleteAdsVideo(params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	return cc.Success("Archived")
}
