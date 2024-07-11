package middlewares

import (
	"fmt"
	"net/http"
	"runtime"
	"strings"

	"github.com/engineeringinflow/inflow-backend/pkg/models"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

// Recover re panic
func (m *Middleware) Recover() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			var cc = c.(*models.CustomContext)
			defer func() {
				if r := recover(); r != nil {
					err, ok := r.(error)
					if !ok {
						err = fmt.Errorf("%v", r)
					}

					var msg = ""
					var reqInfo = cc.GetRequestInfo()
					reqInfo.StatusCode = http.StatusInternalServerError
					if _, file, line, ok := runtime.Caller(4); ok {
						msg = fmt.Sprintf("%s:%d", strings.TrimPrefix(file, "/app/"), line)
					}
					err = fmt.Errorf("%v (Panic Root cause %s)", err, msg)

					cc.CustomLogger.WithRequestInfo(reqInfo).Error(msg, zap.Error(err), zap.StackSkip("stack", 4))
					c.Error(err)
				}
			}()
			return next(c)
		}
	}
}
