package controllers

import (
	"fmt"
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
// @Router /api/v1/common/qrcode [get]
func GetQRCode(c echo.Context) error {
	var cc = c.(*models.CustomContext)

	var params repo.GenerateQRCodeParams
	var err = cc.Bind(&params)
	if err != nil {
		return eris.Wrap(err, "")
	}

	data, err := repo.NewCommonRepo(cc.App.DB).GenerateQRCode(params)
	if err != nil {
		return eris.Wrap(err, "")
	}

	cc.Response().Header().Add("Cache-Control", "no-cache, no-store, must-revalidate")
	cc.Response().Header().Add("Pragma", "no-cache")
	cc.Response().Header().Add("Expires", "0")

	return c.Blob(http.StatusOK, "image/png", data.Bytes())
}

// CreateQRCodes Create qrcodes
// @Tags Common
// @Summary Create qrcodes
// @Description Create qrcodes
// @Accept  json
// @Produce  json
// @Param data body models.CheckExistsForm true "Form"
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/common/qrcode/{content} [get]
func CreateQRCodes(c echo.Context) error {
	var cc = c.(*models.CustomContext)

	var params repo.CreateQRCodesParams
	var err = cc.Bind(&params)
	if err != nil {
		return eris.Wrap(err, "")
	}

	url, err := repo.NewCommonRepo(cc.App.DB).CreateQRCodes(params)
	if err != nil {
		return eris.Wrap(err, "")
	}

	cc.Response().Header().Add("Cache-Control", "no-cache, no-store, must-revalidate")
	cc.Response().Header().Add("Pragma", "no-cache")
	cc.Response().Header().Add("Expires", "0")

	return cc.Success(map[string]interface{}{
		"download_url": url,
	})
}

// DownloadQRCode Get attachment
// @Tags Common
// @Summary Get attachment
// @Description Get attachment
// @Accept  json
// @Produce  json
// @Param data body models.CheckExistsForm true "Form"
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/common/qrcode/download [get]
func DownloadQRCode(c echo.Context) error {
	var cc = c.(*models.CustomContext)

	var params repo.GenerateQRCodeParams
	var err = cc.Bind(&params)
	if err != nil {
		return eris.Wrap(err, "")
	}

	data, err := repo.NewCommonRepo(cc.App.DB).GenerateQRCode(params)
	if err != nil {
		return eris.Wrap(err, "")
	}

	file, err := repo.NewExportRepo(cc.App.DB).GenerateDownloadFileURL(fmt.Sprintf("%s.png", params.Content), data.Bytes())
	if err != nil {
		return eris.Wrap(err, "")
	}

	return cc.Success(map[string]interface{}{
		"download_url": file.DownloadURL,
	})
}
