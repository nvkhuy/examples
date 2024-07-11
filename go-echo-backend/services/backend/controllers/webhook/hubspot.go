package webhook

import (
	"io"

	"github.com/engineeringinflow/inflow-backend/pkg/helper"
	"github.com/engineeringinflow/inflow-backend/pkg/models"
	"github.com/labstack/echo/v4"
)

// HubspotWebhook Hubspot product create
// @Tags Webhook
// @Summary Hubspot product create
// @Description Hubspot product create
// @Accept  json
// @Produce  json
// @Param keyword query string false "Keyword" default(1)
// @Param page query int false "Page index" default(1)
// @Param limit query int false "Size of page" default(20)
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/webhook/hubspot [post]
func HubspotWebhook(c echo.Context) error {
	var cc = c.(*models.CustomContext)

	data, err := io.ReadAll(cc.Request().Body)
	if err != nil {
		return err
	}

	helper.PrintJSONBytes(data)

	return cc.Success("Success")
}
