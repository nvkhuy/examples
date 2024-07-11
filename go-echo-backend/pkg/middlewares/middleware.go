package middlewares

import (
	"github.com/engineeringinflow/inflow-backend/pkg/app"
	"github.com/engineeringinflow/inflow-backend/pkg/validation"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// Middleware struct
type Middleware struct {
	*echo.Echo
	App *app.App
}

// New init
func New(echo *echo.Echo, app *app.App) *Middleware {
	var m = &Middleware{
		Echo: echo,
		App:  app,
	}
	m.Use(middleware.RequestID())
	m.Use(middleware.CORS())
	m.Use(m.RegisterCustomContext())
	m.Use(m.Logger())
	m.Use(m.Recover())
	m.UseSentry()

	m.Validator = validation.RegisterValidation()

	m.HTTPErrorHandler = m.CustomErrorHandler
	return m
}
