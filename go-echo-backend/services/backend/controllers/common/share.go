package controllers

import (
	"net/http"

	"github.com/engineeringinflow/inflow-backend/pkg/models"
	"github.com/engineeringinflow/inflow-backend/pkg/repo"
	"github.com/labstack/echo/v4"
	"github.com/rotisserie/eris"
)

// GetShareLink Get share link
// @Tags Common
// @Summary Get share link
// @Description Get share link
// @Accept  json
// @Produce  json
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/common/share/link/{link_id} [get]
func GetShareLink(c echo.Context) error {
	var cc = c.(*models.CustomContext)

	var params repo.GetShareLinkParams
	var err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, "")
	}

	var link = repo.NewCommonRepo(cc.App.DB).GetShareLink(params)
	if err != nil {
		return eris.Wrap(err, "")
	}

	return cc.Redirect(http.StatusPermanentRedirect, link)
}

// CreateShareLink Create share link
// @Tags Common
// @Summary Create share link
// @Description Create share link
// @Accept  json
// @Produce  json
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/common/share/link/{reference_id} [post]
func CreateShareLink(c echo.Context) error {
	var cc = c.(*models.CustomContext)

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, "")
	}
	var params repo.CreateShareLinkParams
	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, "")
	}

	params.JwtClaimsInfo = claims
	records, err := repo.NewCommonRepo(cc.App.DB).CreateShareLink(params)
	if err != nil {
		return eris.Wrap(err, "")
	}

	return cc.Success(records)
}

// CreateCheckoutShareLink Create share link
// @Tags Common
// @Summary Create share link
// @Description Create share link
// @Accept  json
// @Produce  json
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/common/share/checkout_link/{reference_id} [post]
func CreateCheckoutShareLink(c echo.Context) error {
	var cc = c.(*models.CustomContext)

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, "")
	}
	var params repo.CreateShareLinkParams
	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, "")
	}

	params.JwtClaimsInfo = claims
	records, err := repo.NewCommonRepo(cc.App.DB).CreateCheckoutShareLink(params)
	if err != nil {
		return eris.Wrap(err, "")
	}

	return cc.Success(records)
}

// CreateCheckoutShareLink Create share link
// @Tags Common
// @Summary Create share link
// @Description Create share link
// @Accept  json
// @Produce  json
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/common/share/checkout_link/{reference_id} [get]
func GetCheckoutInfo(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.GetCheckoutInfoParams
	var err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, "")
	}

	records, err := repo.NewCommonRepo(cc.App.DB).GetCheckoutInfo(params)
	if err != nil {
		return eris.Wrap(err, "")
	}

	return cc.Success(records)
}
