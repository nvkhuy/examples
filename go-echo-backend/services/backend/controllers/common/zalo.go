package controllers

import (
	"github.com/engineeringinflow/inflow-backend/pkg/models"
	"github.com/engineeringinflow/inflow-backend/pkg/repo"
	"github.com/labstack/echo/v4"
	"github.com/rotisserie/eris"
)

// SendZNS Check exists
// @Tags Common
// @Summary Check exists
// @Description Check exists
// @Accept  json
// @Produce  json
// @Param data body models.CheckExistsForm true "Form"
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/common/zns/send [post]
func SendZNS(c echo.Context) error {
	var cc = c.(*models.CustomContext)

	var params repo.SendZNSParams
	var err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, "")
	}

	result, err := repo.NewCommonRepo(cc.App.DB).SendZNS(params)
	if err != nil {
		return eris.Wrap(err, "")
	}

	return cc.Success(result)
}
