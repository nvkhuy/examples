package controllers

import (
	"github.com/engineeringinflow/inflow-backend/pkg/models"
	"github.com/engineeringinflow/inflow-backend/pkg/repo"
	"github.com/labstack/echo/v4"
	"github.com/rotisserie/eris"
)

// GetDataAnalyticRFQ
// @Tags Admin-PO
// @Summary Delete Setting Size
// @Description Delete Setting Size
// @Accept  json
// @Produce  json
// @Param keyword query string false "Keyword"
// @Param page query int false "Page number"
// @Param limit query int false "Limit"
// @Success 200 {object} query.Pagination{records=[]models.SettingSize}
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/buyer/dashboard/rfq [put]
func GetDataAnalyticRFQ(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.BuyerDataAnalyticRFQParams

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	params.JwtClaimsInfo = claims
	result, err := repo.NewDataAnalyticRepo(cc.App.AnalyticDB).WithDB(cc.App.DB).BuyerDataAnalyticRFQ(params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}
	return cc.Success(result)
}

// GetDataAnalyticPendingTasks
// @Tags Admin-PO
// @Summary Delete Setting Size
// @Description Delete Setting Size
// @Accept  json
// @Produce  json
// @Param keyword query string false "Keyword"
// @Param page query int false "Page number"
// @Param limit query int false "Limit"
// @Success 200 {object} query.Pagination{records=[]models.SettingSize}
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/buyer/dashboard/pending_tasks [put]
func GetDataAnalyticPendingTasks(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.BuyerDataAnalyticPendingTasksParams

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	params.JwtClaimsInfo = claims
	result, err := repo.NewDataAnalyticRepo(cc.App.AnalyticDB).WithDB(cc.App.DB).BuyerDataAnalyticPendingTasks(params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}
	return cc.Success(result)
}

// GetDataAnalyticPendingPayments
// @Tags Admin-PO
// @Summary Delete Setting Size
// @Description Delete Setting Size
// @Accept  json
// @Produce  json
// @Param keyword query string false "Keyword"
// @Param page query int false "Page number"
// @Param limit query int false "Limit"
// @Success 200 {object} query.Pagination{records=[]models.SettingSize}
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/buyer/dashboard/pending_tasks [put]
func GetDataAnalyticPendingPayments(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.BuyerDataAnalyticPendingPaymentsParams

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	params.JwtClaimsInfo = claims
	result, err := repo.NewDataAnalyticRepo(cc.App.AnalyticDB).WithDB(cc.App.DB).BuyerDataAnalyticPendingPaymentsV2(params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}
	return cc.Success(result)
}

// GetDataAnalyticTotalStyleProduced
// @Tags Admin-PO
// @Summary Delete Setting Size
// @Description Delete Setting Size
// @Accept  json
// @Produce  json
// @Param keyword query string false "Keyword"
// @Param page query int false "Page number"
// @Param limit query int false "Limit"
// @Success 200 {object} query.Pagination{records=[]models.SettingSize}
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/buyer/dashboard/pending_tasks [put]
func GetDataAnalyticTotalStyleProduced(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.BuyerDataAnalyticTotalStyleProduceParams

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	params.JwtClaimsInfo = claims
	result, err := repo.NewDataAnalyticRepo(cc.App.AnalyticDB).WithDB(cc.App.DB).BuyerDataAnalyticTotalStyleProduce(params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}
	return cc.Success(result)
}
