package controllers

import (
	"github.com/engineeringinflow/inflow-backend/pkg/models"
	"github.com/engineeringinflow/inflow-backend/pkg/repo"
	"github.com/engineeringinflow/inflow-backend/services/consumer/tasks"
	"github.com/labstack/echo/v4"
	"github.com/rotisserie/eris"
)

// TrackActivity Track login
// @Tags Seller-Me
// @Summary Track login
// @Description Track login
// @Accept  json
// @Produce  json
// @Success 200 {object} models.User
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/seller/me/track_activity [post]
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

	return cc.Success("Tracked")
}

// OnboardingSubmit
// @Tags Seller-Me
// @Summary Onboarding submit
// @Description Onboarding submit
// @Accept  json
// @Produce  json
// @Success 200 {object} models.User
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/seller/me/onboarding_submit [post]
func OnboardingSubmit(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var form models.BusinessProfileCreateForm

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, "")
	}
	form.JwtClaimsInfo = claims
	err = cc.BindAndValidate(&form)
	if err != nil {
		return eris.Wrap(err, "")
	}

	result, err := repo.NewUserRepo(cc.App.DB).OnboardingSubmit(cc.App.DB.DB, form)
	if err != nil {
		return eris.Wrap(err, "")
	}
	return cc.Success(result)
}

// GetMe get user's profile
// @Tags Seller-Me
// @Summary Get me
// @Description Get me
// @Accept  json
// @Produce  json
// @Success 200 {object} models.User
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/seller/me [get]
func GetMe(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.GetMeParams

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	params.JwtClaimsInfo = claims
	u, err := repo.NewUserRepo(cc.App.DB).GetMe(params)
	if err != nil {
		return eris.Wrap(err, "")
	}

	if u.StripeCustomerID == "" {
		u.CreateStripeCustomer(cc.App.DB, nil, false)
	}

	return cc.Success(u)
}

// UpdateMe Update me
// @Tags Seller-Me
// @Summary Update me
// @Description Update me
// @Accept  json
// @Produce  json
// @Success 200 {object} models.User
// @Param data body models.UserUpdateForm true "User Update Form"
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/seller/me [put]
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

	return cc.Success(u)
}

// Logout logout
// @Security ApiKeyAuth
// @Tags Seller-Me
// @Summary Logout
// @Description Logout
// @Accept  json
// @Produce  json
// @Header 200 {string} Bearer YOUR_TOKEN
// @Failure 404 {object} errs.Error
// @Router /api/v1/seller/me/logout [delete]
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
// @Tags Seller-Me
// @Summary Update password
// @Description Update password
// @Accept  json
// @Produce  json
// @Header 200 {string} Bearer YOUR_TOKEN
// @Failure 404 {object} errs.Error
// @Router /api/v1/seller/me/update_password [put]
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

	return cc.Success(result)
}
