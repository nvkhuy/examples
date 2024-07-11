package controllers

import (
	"github.com/engineeringinflow/inflow-backend/pkg/models"
	"github.com/engineeringinflow/inflow-backend/pkg/repo"
	"github.com/engineeringinflow/inflow-backend/services/consumer/tasks"
	"github.com/labstack/echo/v4"
	"github.com/rotisserie/eris"
)

// UpdateMe Update me
// @Tags Buyer-Me
// @Summary Update me
// @Description Update me
// @Accept  json
// @Produce  json
// @Success 200 {object} models.User
// @Param data body models.UserUpdateForm true "User Update Form"
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/buyer/me/last_shipping_address [get]
func GetLastShippingAddress(c echo.Context) error {
	var cc = c.(*models.CustomContext)

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, "")
	}

	u, err := repo.NewUserRepo(cc.App.DB).GetLastShippingAddress(claims)
	if err != nil {
		return eris.Wrap(err, "")
	}

	return cc.Success(u)
}

// TrackActivity Track login
// @Tags Buyer-Me
// @Summary Track login
// @Description Track login
// @Accept  json
// @Produce  json
// @Success 200 {object} models.User
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/buyer/me/track_activity [post]
func TrackActivity(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var form models.UserTrackActivityForm
	var user models.User

	var err = cc.BindAndValidate(&form)
	if err != nil {
		return eris.Wrap(err, "")
	}

	err = cc.GetUserFromContext(&user)
	if err != nil {
		return eris.Wrap(err, "")
	}

	tasks.TrackActivityTask{
		UserID:                user.ID,
		UserTrackActivityForm: form,
	}.Dispatch(c.Request().Context())

	return cc.Success(user)
}

// GetMe get user's profile
// @Tags Buyer-Me
// @Summary Get me
// @Description Get me
// @Accept  json
// @Produce  json
// @Success 200 {object} models.User
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/buyer/me [get]
func GetMe(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.GetMeParams

	var err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	params.JwtClaimsInfo = claims
	u, err := repo.NewUserRepo(cc.App.DB).GetMe(params)
	if err != nil {
		return eris.Wrap(err, "")
	}

	if u.StripeCustomerID == "" {
		_ = u.CreateStripeCustomer(cc.App.DB, nil, false)
	}

	return cc.Success(u)
}

// UpdateMe Update me
// @Tags Buyer-Me
// @Summary Update me
// @Description Update me
// @Accept  json
// @Produce  json
// @Success 200 {object} models.User
// @Param data body models.UserUpdateForm true "User Update Form"
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/buyer/me [put]
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
	u, err := repo.NewUserRepo(cc.App.DB).UpdateUserByID(form)
	if err != nil {
		return eris.Wrap(err, "")
	}

	tasks.SyncCustomerIOUserTask{
		UserID: u.ID,
	}.Dispatch(c.Request().Context())

	return cc.Success(u)
}

// LogoutClient logout
// @Security ApiKeyAuth
// @Tags Buyer-Me
// @Summary Logout
// @Description Logout
// @Accept  json
// @Produce  json
// @Header 200 {string} Bearer YOUR_TOKEN
// @Failure 404 {object} errs.Error
// @Router /api/v1/buyer/me/logout [delete]
func Logout(c echo.Context) error {
	var cc = c.(*models.CustomContext)

	claim, err := cc.GetJwtClaims()
	if err != nil {
		return eris.Wrap(err, "")
	}

	repo.NewUserRepo(cc.App.DB).Logout(claim.ID)

	return cc.Success("Logout successfully")

}

// UpdatePassword Update password
// @Tags Buyer-Me
// @Summary Update password
// @Description Update password
// @Accept  json
// @Produce  json
// @Header 200 {string} Bearer YOUR_TOKEN
// @Failure 404 {object} errs.Error
// @Router /api/v1/buyer/me/update_password [put]
func UpdatePassword(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var form models.UpdateMyPasswordForm

	claims, err := cc.GetJwtClaims()
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	err = cc.BindAndValidate(&form)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	err = repo.NewUserRepo(cc.App.DB).UpdatePassword(claims.ID, form)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	return cc.Success("Updated")
}

// GetTeams Teams
// @Tags Buyer-Me
// @Summary Teams
// @Description Teams
// @Accept  json
// @Produce  json
// @Header 200 {string} Bearer YOUR_TOKEN
// @Failure 404 {object} errs.Error
// @Router /api/v1/buyer/me/teams [put]
func GetTeams(c echo.Context) error {
	var cc = c.(*models.CustomContext)

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	var result = repo.NewUserRepo(cc.App.DB).GetTeams(claims)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	return cc.Success(result)
}

// UpdateMe Update me
// @Tags Buyer-Me
// @Summary Update me
// @Description Update me
// @Accept  json
// @Produce  json
// @Success 200 {object} models.User
// @Param data body models.UserUpdateForm true "User Update Form"
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/buyer/me/complete_inquiry_tutorial [put]
func CompleteInquiryTutorial(c echo.Context) error {
	var cc = c.(*models.CustomContext)

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, "")
	}

	err = repo.NewUserRepo(cc.App.DB).CompleteInquiryTutorial(claims)
	if err != nil {
		return eris.Wrap(err, "")
	}

	return cc.Success("Success")
}
