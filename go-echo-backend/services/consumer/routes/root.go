package routes

import (
	"github.com/engineeringinflow/inflow-backend/pkg/config"
	"github.com/engineeringinflow/inflow-backend/pkg/models"
	"github.com/labstack/echo/v4"
)

// SetupAppInfoRoute setup root's routes
func (router *Router) SetupAppInfoRoute(g *echo.Group) {
	g.GET("/config", configHandler)
	g.GET("/info", infoHandler)

}

// SetupRootRoute setup root's routes
func (router *Router) SetupRootRoute(g *echo.Group) {
	g.GET("/", rootHandler)
	g.GET("/health_check", healthCheckHandler)
}

func healthCheckHandler(c echo.Context) error {
	var response = models.M{
		"status": "ok",
	}
	return c.JSON(200, response)
}

func rootHandler(c echo.Context) error {
	var response = models.M{
		"status": "ok",
	}
	return c.JSON(200, response)

}

func infoHandler(c echo.Context) error {
	var cc = c.(*models.CustomContext)

	var response = cc.App.Config.GetServerCommonInfo(config.ServiceConsumer)

	return cc.Success(response)

}
func configHandler(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	return cc.Success(cc.App.Config)

}
