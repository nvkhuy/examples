package middlewares

import (
	"github.com/labstack/echo/v4"

	"github.com/engineeringinflow/inflow-backend/pkg/logger"
	"github.com/engineeringinflow/inflow-backend/pkg/models"
)

// RegisterCustomContext register custom context
func (m *Middleware) RegisterCustomContext() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			var cc = &models.CustomContext{
				Context:      c,
				App:          m.App,
				CustomLogger: logger.New("api"),
			}
			return next(cc)
		}
	}
}
