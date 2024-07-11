package routes

import (
	controllers "github.com/engineeringinflow/inflow-backend/services/backend/controllers/app"

	"github.com/labstack/echo/v4"
)

func (router *Router) SetupAppRoutes(g *echo.Group) {
	var authGroup = g.Group("/auth")

	authGroup.POST("/login_email", controllers.LoginEmail)
}
