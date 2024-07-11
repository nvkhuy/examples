package middlewares

import (
	"github.com/labstack/echo/v4"
)

func (m *Middleware) ResponseDownloadableHeader() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Response().Header().Set(echo.HeaderContentType, echo.MIMEOctetStream)
			c.Response().Header().Set(echo.HeaderContentDisposition, "attachment")
			return next(c)
		}
	}
}
