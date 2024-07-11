package controllers

import (
	"github.com/engineeringinflow/inflow-backend/pkg/models"
	"github.com/engineeringinflow/inflow-backend/pkg/repo"
	"github.com/labstack/echo/v4"
)

// GetConstants Generate sitemap
// @Tags Common
// @Summary Generate sitemap
// @Description Generate sitemap
// @Accept  json
// @Produce  json
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/common/sitemap [get]
func GetSitemap(c echo.Context) error {
	var cc = c.(*models.CustomContext)

	var resp = repo.NewCommonRepo(cc.App.DB).GenerateSitemap()

	return cc.Success(resp)
}
