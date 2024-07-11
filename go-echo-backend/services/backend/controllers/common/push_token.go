package controllers

import (
	"github.com/engineeringinflow/inflow-backend/pkg/models"
	"github.com/engineeringinflow/inflow-backend/pkg/repo"
	"github.com/engineeringinflow/inflow-backend/services/consumer/tasks"
	"github.com/labstack/echo/v4"
	"github.com/rotisserie/eris"
)

// GetPushTokens Get user tokens
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
// @Router /api/v1/common/push_tokens [post]
func GetPushTokens(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.PaginatePushTokensParams

	var err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	var result = repo.NewPushTokenRepo(cc.App.DB).PaginatePushTokens(params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	return cc.Success(result)
}

// CreatePushToken CreateFromPayload user's device
// @Tags Common
// @Summary CreateFromPayload user's device
// @Description CreateFromPayload user's device
// @Accept  json
// @Produce  json
// @Param data body models.PushTokenCreateForm true "Form"
// @Success 200 {object} models.User
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/common/push_tokens [post]
func CreatePushToken(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var form models.PushTokenCreateForm

	var err = cc.BindAndValidate(&form)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	claim, err := cc.GetJwtClaims()
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	device, err := repo.NewPushTokenRepo(cc.App.DB).AddPushToken(claim.ID, form)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	tasks.AddCustomerIOUserDeviceTask{
		UserID: claim.ID,
		Device: device,
	}.Dispatch(c.Request().Context())

	return cc.Success(device)
}

// DeletePushToken Delete user's push token
// @Tags Common
// @Summary Delete user's push token
// @Description Delete user's push token
// @Accept  json
// @Produce  json
// @Param token path string true "User Push Token"
// @Success 200 {object} models.User
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/common/push_tokens/{token} [delete]
func DeletePushToken(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var deviceToken = cc.GetPathParamString("token")

	var err = repo.NewPushTokenRepo(cc.App.DB).DeletePushToken(deviceToken)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	return cc.Success("Push token is deleted")
}
