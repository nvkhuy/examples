package controllers

import (
	"github.com/engineeringinflow/inflow-backend/pkg/models"
	"github.com/engineeringinflow/inflow-backend/pkg/repo"
	"github.com/labstack/echo/v4"
	"github.com/rotisserie/eris"
)

// AdminNotificationList
// @Tags Admin-Notification
// @Summary Paginate notification
// @Description Paginate notification
// @Accept  json
// @Produce  json
// @Param page query int false "Page number"
// @Success 200 {object} query.Pagination{records=[]models.CmsNotification}
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/notifications [get]
func AdminNotificationList(c echo.Context) error {
	var cc = c.(*models.CustomContext)

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	var params repo.PaginateCmsNotificationsParams
	err = cc.Bind(&params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	params.JwtClaimsInfo = claims
	var result = repo.NewCmsNotificationRepo(cc.App.DB).PaginateCmsNotifications(params)

	return cc.Success(result)
}

// AdminNotificationMarkSeen
// @Tags Admin-Notification
// @Summary Paginate notification
// @Description Paginate notification
// @Accept  json
// @Produce  json
// @Param page query int false "Page number"
// @Success 200 {object} query.Pagination{records=[]models.CmsNotification}
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/notifications/{notification_id}/mark_seen [put]
func AdminNotificationMarkSeen(c echo.Context) error {
	var cc = c.(*models.CustomContext)

	var err = repo.NewCmsNotificationRepo(cc.App.DB).MarkSeen(cc.GetPathParamString("notification_id"))

	if err != nil {
		return eris.Wrap(err, err.Error())
	}
	return cc.Success("Seen")
}

// Admin Notification
// @Tags Admin-Notification
// @Summary Paginate notification
// @Description Paginate notification
// @Accept  json
// @Produce  json
// @Param page query int false "Page number"
// @Success 200 {object} query.Pagination{records=[]models.CmsNotification}
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/notifications/mark_seen_all [put]
func AdminNotificationMarkSeenAll(c echo.Context) error {
	var cc = c.(*models.CustomContext)

	var err = repo.NewCmsNotificationRepo(cc.App.DB).MarkSeenAll()

	if err != nil {
		return eris.Wrap(err, err.Error())
	}
	return cc.Success("Seen")
}
