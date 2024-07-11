package controllers

import (
	"github.com/engineeringinflow/inflow-backend/pkg/models"
	"github.com/engineeringinflow/inflow-backend/pkg/repo"
	"github.com/labstack/echo/v4"

	"github.com/rotisserie/eris"
)

// PaginateAsFeaturedIns
// @Tags Admin-Ads
// @Summary PaginateAsFeaturedIns
// @Description PaginateAsFeaturedIns
// @Accept  json
// @Produce  json
// @Param keyword query string false "Keyword"
// @Param name query string false "Name"
// @Success 200 {object} models.AsFeaturedIn
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/as_featured_ins [get]
func PaginateAsFeaturedIns(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.PaginateAsFeaturedInParams

	var err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	var result = repo.NewAsFeaturedInRepo(cc.App.DB).PaginateAsFeaturedIn(params)
	return cc.Success(result)
}
