package controllers

import (
	"github.com/engineeringinflow/inflow-backend/pkg/models"
	"github.com/engineeringinflow/inflow-backend/pkg/repo"
	"github.com/labstack/echo/v4"
	"github.com/rotisserie/eris"
)

// GetDocs Get user tokens
// @Tags Common
// @Summary Get user tokens
// @Description Get user tokens
// @Accept  json
// @Produce  json
// @Param keyword query string false "Keyword" default(1)
// @Param page query int false "Page index" default(1)
// @Param limit query int false "Size of page" default(20)
// @Success 200 {object} query.Pagination{records=[]models.PushToken}
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/common/docs [post]
func GetDocs(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.GetDocsParams

	var err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	var result = repo.NewCommonRepo(cc.App.DB).GetDocs(params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	return cc.Success(result)
}
