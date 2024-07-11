package controllers

import (
	"github.com/engineeringinflow/inflow-backend/pkg/models"
	"github.com/engineeringinflow/inflow-backend/pkg/repo"
	"github.com/engineeringinflow/inflow-backend/pkg/repo/query/queryfunc"
	"github.com/labstack/echo/v4"
	"github.com/rotisserie/eris"
)

// PageCatalog
// @Tags Marketplace-Page
// @Summary PageCatalog
// @Description PageCatalog
// @Accept  json
// @Produce  json
// @Success 200 {object} models.PageDetailResponse
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/pages/catalog [get]
func PageCatalog(c echo.Context) error {
	var cc = c.(*models.CustomContext)

	page, err := repo.NewPageRepo(cc.App.DB).GetPageByType("our_catalog", queryfunc.PageBuilderOptions{})
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	result, err := repo.NewPageRepo(cc.App.DB).PageCatalog(page.ID)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	var response = &models.PageDetailResponse{
		ID:      page.ID,
		Title:   page.Title,
		Url:     page.Url,
		Content: result,
	}

	return cc.Success(response)
}

// PageHome
// @Tags Marketplace-Page
// @Summary PageHome
// @Description PageHome
// @Accept  json
// @Produce  json
// @Success 200 {object} models.PageDetailResponse
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/pages/home [get]
func PageHome(c echo.Context) error {
	var cc = c.(*models.CustomContext)

	page, err := repo.NewPageRepo(cc.App.DB).GetPageByType("home", queryfunc.PageBuilderOptions{})
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	result, err := repo.NewPageRepo(cc.App.DB).PageHome(page.ID)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	var response = &models.PageDetailResponse{
		Content: result,
	}

	return cc.Success(response)
}
