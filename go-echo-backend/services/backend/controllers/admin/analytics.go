package controllers

import (
	"github.com/engineeringinflow/inflow-backend/pkg/models"
	"github.com/engineeringinflow/inflow-backend/pkg/repo"
	"github.com/labstack/echo/v4"
	"github.com/rotisserie/eris"
)

// PaginatePotentialOverdueInquiries
// @Tags Admin-Analytics
// @Summary PaginatePotentialOverdueInquiries
// @Description PaginateBlogCategory
// @Accept  json
// @Produce  json
// @Param keyword query string false "Keyword"
// @Param name query string false "Name"
// @Success 200 {object} models.BlogCategory
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/analytics/inquiries/potential_overdue [get]
func PaginatePotentialOverdueInquiries(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.PaginatePotentialOverdueInquiriesParams

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	params.JwtClaimsInfo = claims
	var result = repo.NewAnalyticsRepo(cc.App.DB).PaginatePotentialOverdueInquiries(params)
	return cc.Success(result)
}

// PaginateInquiriesTimeline
// @Tags Admin-Analytics
// @Summary PaginateInquiriesTimeline
// @Description PaginateBlogCategory
// @Accept  json
// @Produce  json
// @Param keyword query string false "Keyword"
// @Param name query string false "Name"
// @Success 200 {object} models.BlogCategory
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/analytics/inquiries/potential_overdue [get]
func PaginateInquiriesTimeline(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.PaginatePotentialOverdueInquiriesParams

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	params.JwtClaimsInfo = claims
	var result = repo.NewAnalyticsRepo(cc.App.DB).PaginateInquiriesTimeline(params)
	return cc.Success(result)
}
