package controllers

import (
	"github.com/engineeringinflow/inflow-backend/pkg/models"
	"github.com/engineeringinflow/inflow-backend/pkg/repo"
	"github.com/labstack/echo/v4"
	"github.com/rotisserie/eris"
	"github.com/thaitanloi365/go-utils/values"
)

// GetSettingTaxes Get settings taxes
// @Tags Common
// @Summary Get settings taxes
// @Description Get settings taxes
// @Accept  json
// @Produce  json
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/common/settings/taxes [get]
func GetSettingTaxes(c echo.Context) error {
	var cc = c.(*models.CustomContext)

	var params repo.GetSettingTaxParams
	var err = cc.Bind(&params)
	if err != nil {
		return eris.Wrap(err, "")
	}

	records, err := repo.NewSettingTaxRepo(cc.App.DB).GetSettingTaxes(params)
	if err != nil {
		return eris.Wrap(err, "")
	}

	return cc.Success(records)
}

// GetSettingSizes Get settings sizes
// @Tags Common
// @Summary Get settings sizes
// @Description Get settings sizes
// @Accept  json
// @Produce  json
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/common/settings/sizes [get]
func GetSettingSizes(c echo.Context) error {
	var cc = c.(*models.CustomContext)

	var params repo.GetSettingSizesParams
	var err = cc.Bind(&params)
	if err != nil {
		return eris.Wrap(err, "")
	}

	records, err := repo.NewSettingSizeRepo(cc.App.DB).GetSettingSizes(params)
	if err != nil {
		return eris.Wrap(err, "")
	}

	return cc.Success(records)
}

// GetSettingBankInfos Get settings taxes
// @Tags Common
// @Summary Get settings taxes
// @Description Get settings taxes
// @Accept  json
// @Produce  json
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/common/settings/bank_infos [get]
func GetSettingBankInfos(c echo.Context) error {
	var cc = c.(*models.CustomContext)

	var params repo.GetSettingBankInfosParams
	var err = cc.Bind(&params)
	if err != nil {
		return eris.Wrap(err, "")
	}

	params.IsDisabled = values.Bool(false)
	records, err := repo.NewSettingBankRepo(cc.App.DB).GetSettingBankInfos(params)
	if err != nil {
		return eris.Wrap(err, "")
	}

	return cc.Success(records)
}
