package controllers

import (
	"net/http"

	"github.com/engineeringinflow/inflow-backend/pkg/models"
	"github.com/engineeringinflow/inflow-backend/pkg/repo"
	"github.com/labstack/echo/v4"
	"github.com/rotisserie/eris"
)

// GetAttachment Get attachment
// @Tags Common
// @Summary Get attachment
// @Description Get attachment
// @Accept  json
// @Produce  json
// @Param data body models.CheckExistsForm true "Form"
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/common/attachments/{file_key} [get]
func GetAttachment(c echo.Context) error {
	var cc = c.(*models.CustomContext)

	var params repo.GetAttachmentParams
	var err = cc.Bind(&params)
	if err != nil {
		return eris.Wrap(err, "")
	}

	var url = repo.NewCommonRepo(cc.App.DB).GetAttachment(cc.App.S3Client, params)

	cc.Response().Header().Add("Cache-Control", "no-cache, no-store, must-revalidate")
	cc.Response().Header().Add("Pragma", "no-cache")
	cc.Response().Header().Add("Expires", "0")

	return cc.Redirect(http.StatusPermanentRedirect, url)
}

// GetThumbnailAttachment Get attachment
// @Tags Common
// @Summary Get attachment
// @Description Get attachment
// @Accept  json
// @Produce  json
// @Param data body models.CheckExistsForm true "Form"
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/common/thumbnail_attachments/{file_key} [get]
func GetThumbnailAttachment(c echo.Context) error {
	var cc = c.(*models.CustomContext)

	var params repo.GetAttachmentParams
	var err = cc.Bind(&params)
	if err != nil {
		return eris.Wrap(err, "")
	}

	var url = repo.NewCommonRepo(cc.App.DB).GetThumbnailAttachment(cc.App.S3Client, params)

	cc.Response().Header().Add("Cache-Control", "no-cache, no-store, must-revalidate")
	cc.Response().Header().Add("Pragma", "no-cache")
	cc.Response().Header().Add("Expires", "0")

	return cc.Redirect(http.StatusPermanentRedirect, url)
}

// GetBlurAttachment Get blur attachment
// @Tags Common
// @Summary Get blur attachment
// @Description Get blur attachment
// @Accept  json
// @Produce  json
// @Param data body models.CheckExistsForm true "Form"
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/common/blur_attachments/{file_key} [get]
func GetBlurAttachment(c echo.Context) error {
	var cc = c.(*models.CustomContext)

	var params repo.GetAttachmentParams
	var err = cc.Bind(&params)
	if err != nil {
		return eris.Wrap(err, "")
	}

	var url = repo.NewCommonRepo(cc.App.DB).GetBlurAttachment(cc.App.S3Client, params)

	cc.Response().Header().Add("Cache-Control", "no-cache, no-store, must-revalidate")
	cc.Response().Header().Add("Pragma", "no-cache")
	cc.Response().Header().Add("Expires", "0")

	return cc.Redirect(http.StatusPermanentRedirect, url)
}

// GetBlurAttachment Get blur attachment
// @Tags Common
// @Summary Get blur attachment
// @Description Get blur attachment
// @Accept  json
// @Produce  json
// @Param data body models.CheckExistsForm true "Form"
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/common/blur_attachments/{file_key} [get]
func GetBlurDataURL(c echo.Context) error {
	var cc = c.(*models.CustomContext)

	var params repo.GetAttachmentParams
	var err = cc.Bind(&params)
	if err != nil {
		return eris.Wrap(err, "")
	}

	var url = repo.NewCommonRepo(cc.App.DB).GetBlurAttachment(cc.App.S3Client, params)

	cc.Response().Header().Add("Cache-Control", "no-cache, no-store, must-revalidate")
	cc.Response().Header().Add("Pragma", "no-cache")
	cc.Response().Header().Add("Expires", "0")

	return cc.Redirect(http.StatusPermanentRedirect, url)
}

// GetAttachmentDownloadLink Get attachment
// @Tags Common
// @Summary Get attachment
// @Description Get attachment
// @Accept  json
// @Produce  json
// @Param data body models.CheckExistsForm true "Form"
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/common/download_url/{file_key} [get]
func GetAttachmentDownloadLink(c echo.Context) error {
	var cc = c.(*models.CustomContext)

	var params repo.GetDownloadLinkParams
	var err = cc.Bind(&params)
	if err != nil {
		return eris.Wrap(err, "")
	}

	var url = repo.NewCommonRepo(cc.App.DB).GetDownloadLink(params)

	return cc.Success(map[string]interface{}{
		"download_url": url,
	})
}
