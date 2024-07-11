package controllers

import (
	"github.com/engineeringinflow/inflow-backend/pkg/models"
	"github.com/engineeringinflow/inflow-backend/pkg/repo"
	"github.com/labstack/echo/v4"

	"github.com/rotisserie/eris"
)

// SearchSubscribers Search subscribers
// @Tags Admin-Subscribers
// @Summary Search subscribers
// @Description Search subscribers
// @Accept  json
// @Produce  json
// @Param keyword query string false "Keyword" default(1)
// @Param limit query int false "Size of page" default(20)
// @Success 200 {object} []models.Subscriber
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/subscribers/search [get]
func SearchSubscribers(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.PaginateSubscribersParams

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	err = cc.Bind(&params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	params.JwtClaimsInfo = claims
	var result = repo.NewSubscriberRepo(cc.App.DB).PaginateSubscribers(params)

	return cc.Success(result)
}

// UpdateSubscriber Update Subscriber
// @Tags Admin-Subscriber
// @Summary Update Subscriber
// @Description Update Subscriber
// @Accept  json
// @Produce  json
// @Param Subscriber_id path string true "ID"
// @Param data body models.SubscribeByEmailForm true "Form"
// @Success 200 {object} models.Subscriber
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/subscribers/{subscriber_id} [put]
func UpdateSubscriber(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var id = cc.GetPathParamString("subscriber_id")
	var form models.SubscribeByEmailForm

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	err = cc.BindAndValidate(&form)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	form.JwtClaimsInfo = claims
	u, err := repo.NewSubscriberRepo(cc.App.DB).UpdateSubscriberByID(id, form)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	return cc.Success(u)
}

// ArchiveSubscriber Archive Subscriber
// @Tags Admin-Subscriber
// @Summary Archive Subscriber
// @Description Archive Subscriber
// @Accept  json
// @Produce  json
// @Param id path string true "ID"
// @Success 200 {object} models.M
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/subscribers/{subscriber_id}/archive [delete]
func ArchiveSubscriber(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var id = cc.GetPathParamString("subscriber_id")

	err := repo.NewSubscriberRepo(cc.App.DB).ArchiveSubscriberByID(id)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	return cc.Success("Archived")
}
