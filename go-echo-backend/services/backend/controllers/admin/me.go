package controllers

import (
	"github.com/engineeringinflow/inflow-backend/pkg/models"
	"github.com/engineeringinflow/inflow-backend/pkg/repo"
	"github.com/engineeringinflow/inflow-backend/services/consumer/tasks"
	"github.com/labstack/echo/v4"
	"github.com/rotisserie/eris"
)

// TrackActivity Track login
// @Tags Me
// @Summary Track login
// @Description Track login
// @Accept  json
// @Produce  json
// @Success 200 {object} models.User
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/me/track_activity [post]
func TrackActivity(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var form models.UserTrackActivityForm

	claims, err := cc.GetJwtClaims()
	if err != nil {
		return eris.Wrap(err, "")
	}

	err = cc.BindAndValidate(&form)
	if err != nil {
		return eris.Wrap(err, "")
	}

	tasks.TrackActivityTask{
		UserID:                claims.ID,
		UserTrackActivityForm: form,
	}.Dispatch(c.Request().Context())

	return cc.Success("Tracker")
}

// GetMe get user's profile
// @Tags Me
// @Summary Get me
// @Description Get me
// @Accept  json
// @Produce  json
// @Success 200 {object} models.User
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/me [get]
func GetMe(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.GetMeParams

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, "")
	}

	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, "")
	}

	params.JwtClaimsInfo = claims
	u, err := repo.NewUserRepo(cc.App.DB).GetMe(params)
	if err != nil {
		return eris.Wrap(err, "")
	}

	return cc.Success(u)
}

// UpdateMe Update me
// @Tags Me
// @Summary Update me
// @Description Update me
// @Accept  json
// @Produce  json
// @Success 200 {object} models.User
// @Param data body models.UserUpdateForm true "User Update Form"
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/me [put]
func UpdateMe(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var form models.UserUpdateForm

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, "")
	}

	err = cc.BindAndValidate(&form)
	if err != nil {
		return eris.Wrap(err, "")
	}

	form.JwtClaimsInfo = claims
	form.UserID = claims.GetUserID()
	u, err := repo.NewUserRepo(cc.App.DB).UpdateUserByID(form)
	if err != nil {
		return eris.Wrap(err, "")
	}

	return cc.Success(u)
}

// Logout logout
// @Security ApiKeyAuth
// @Tags Me
// @Summary Logout
// @Description Logout
// @Accept  json
// @Produce  json
// @Header 200 {string} Bearer YOUR_TOKEN
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/me/logout [delete]
func Logout(c echo.Context) error {
	var cc = c.(*models.CustomContext)

	claim, err := cc.GetJwtClaims()
	if err != nil {
		return eris.Wrap(err, "")
	}

	repo.NewUserRepo(cc.App.DB).Logout(claim.ID)

	return cc.Success("Logout successfully")

}
