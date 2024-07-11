package controllers

import (
	"strconv"

	"github.com/engineeringinflow/inflow-backend/pkg/models"
	"github.com/engineeringinflow/inflow-backend/pkg/repo"
	"github.com/engineeringinflow/inflow-backend/services/consumer/tasks"
	"github.com/labstack/echo/v4"
	"github.com/rotisserie/eris"
)

func UploadProductFile(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	formSiteName := cc.FormValue("site_name")
	formScrapeDate := cc.FormValue("scrape_date")
	scrapeDate, err := strconv.Atoi(formScrapeDate)
	if err != nil {
		return err
	}
	formFile, err := cc.FormFile("file")
	if err != nil {
		return err
	}

	result, err := repo.NewProductFileUploadRepo(cc.App.DB).UploadFile(&models.UploadProductFileRequest{
		JwtClaimsInfo: claims,
		SiteName:      formSiteName,
		ScrapeDate:    scrapeDate,
		File:          formFile,
	})
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	tasks.UploadProductFileTask{
		UploadID:   result.ID,
		SiteName:   result.SiteName,
		FileKey:    result.Attachment.FileKey,
		ScrapeDate: result.ScrapeDate,
	}.Dispatch(cc.Request().Context())

	return cc.Success(result)
}

func GetProductFileUploadInfoList(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, err.Error())
	}
	var params models.GetProductFileListRequest
	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}
	params.JwtClaimsInfo = claims
	result := repo.NewProductFileUploadRepo(cc.App.DB).GetProductFileList(&params)

	return cc.Success(result)
}
