package controllers

import (
	"github.com/engineeringinflow/inflow-backend/pkg/models"
	"github.com/engineeringinflow/inflow-backend/pkg/repo"
	"github.com/labstack/echo/v4"
	"github.com/rotisserie/eris"
)

// PaginateUserNotifications
// @Tags Marketplace-Notification
// @Summary Paginate notification
// @Description Paginate notification
// @Accept  json
// @Produce  json
// @Param page query int false "Page number"
// @Success 200 {object} query.Pagination{records=[]models.UserNotification}
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/buyer/notifications [get]
func PaginateUserNotifications(c echo.Context) error {
	var cc = c.(*models.CustomContext)

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	var params repo.PaginateUserNotificationsParams
	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	params.JwtClaimsInfo = claims
	var result = repo.NewUserNotificationRepo(cc.App.DB).PaginateUserNotifications(params)

	return cc.Success(result)
}

// NotificationMarkSeenAll
// @Tags Marketplace-Notification
// @Summary NotificationMarkSeenAll
// @Description NotificationMarkSeenAll
// @Accept  json
// @Produce  json
// @Param page query int false "Page number"
// @Success 200 {object} query.Pagination{records=[]models.UserNotification}
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/buyer/notifications/mark_seen_all [put]
func NotificationMarkSeenAll(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	err = repo.NewUserNotificationRepo(cc.App.DB).MarkSeenAll(claims)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}
	return cc.Success("Seen")
}

// UserNotificationMarkSeen
// @Tags Marketplace-Notification
// @Summary UserNotificationMarkSeen
// @Description UserNotificationMarkSeen
// @Accept  json
// @Produce  json
// @Param page query int false "Page number"
// @Success 200 {object} query.Pagination{records=[]models.CmNotification}
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/buyer/notifications/{notification_id}/mark_seen [put]
func UserNotificationMarkSeen(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.MarkSeenNotificationParams

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	params.JwtClaimsInfo = claims
	err = repo.NewUserNotificationRepo(cc.App.DB).MarkSeen(params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}
	return cc.Success("Seen")
}

// PaginateSysNotifications
// @Tags Marketplace-Notification
// @Summary Paginate notification
// @Description Paginate notification
// @Accept  json
// @Produce  json
// @Param page query int false "Page number"
// @Success 200 {object} query.Pagination{records=[]models.UserNotification}
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/buyer/sys/notifications [get]
func PaginateSysNotifications(c echo.Context) error {
	var cc = c.(*models.CustomContext)

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	var params repo.PaginateSysNotificationsParams
	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	params.JwtClaimsInfo = claims
	var result = repo.NewSysNotificationRepo(cc.App.DB).Paginate(params)

	return cc.Success(result)
}

// SysNotificationMarkSeenAll
// @Tags Marketplace-Notification
// @Summary NotificationMarkSeenAll
// @Description NotificationMarkSeenAll
// @Accept  json
// @Produce  json
// @Param page query int false "Page number"
// @Success 200 {object} query.Pagination{records=[]models.UserNotification}
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/buyer/sys/notifications/mark_seen_all [put]
func SysNotificationMarkSeenAll(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.MarkSeenAllSysNotificationParams

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	params.JwtClaimsInfo = claims
	err = repo.NewSysNotificationRepo(cc.App.DB).MarkSeenAll(params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}
	return cc.Success("Seen")
}

// SysNotificationMarkSeen
// @Tags Marketplace-Notification
// @Summary UserNotificationMarkSeen
// @Description UserNotificationMarkSeen
// @Accept  json
// @Produce  json
// @Param page query int false "Page number"
// @Success 200 {object} query.Pagination{records=[]models.CmNotification}
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/buyer/notifications/{notification_id}/mark_seen [put]
func SysNotificationMarkSeen(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.MarkSeenSysNotificationParams

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	params.JwtClaimsInfo = claims
	err = repo.NewSysNotificationRepo(cc.App.DB).MarkSeen(params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}
	return cc.Success("Seen")
}
