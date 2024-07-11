package controllers

import (
	"github.com/engineeringinflow/inflow-backend/pkg/models"
	"github.com/engineeringinflow/inflow-backend/pkg/repo"
	"github.com/labstack/echo/v4"

	"github.com/rotisserie/eris"
)

// PaginateAdsVideos
// @Tags Marketplace-Ads
// @Summary PaginateAdsVideos
// @Description PaginateAdsVideos
// @Accept  json
// @Produce  json
// @Param keyword query string false "Keyword"
// @Param name query string false "Name"
// @Success 200 {object} models.AdsVideo
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/ads_videos [get]
func PaginateAdsVideos(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.PaginateAdsVideoParams

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	params.JwtClaimsInfo = claims
	var result = repo.NewAdsVideoRepo(cc.App.DB).PaginateAdsVideo(params)
	return cc.Success(result)
}
