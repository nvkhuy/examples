package middlewares

import (
	"github.com/engineeringinflow/inflow-backend/pkg/models"
	"github.com/getsentry/sentry-go"
	sentryecho "github.com/getsentry/sentry-go/echo"
	"github.com/labstack/echo/v4"
)

// Recover re panic
func (m *Middleware) UseSentry() {
	_ = sentry.Init(sentry.ClientOptions{
		Dsn:              m.App.Config.SentryDsn,
		Debug:            true,
		AttachStacktrace: true,
		EnableTracing:    true,
		Environment:      m.App.Config.BuildEnv,
	})

	m.Use(sentryecho.New(sentryecho.Options{
		Repanic: true,
	}))

	m.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if hub := sentryecho.GetHubFromContext(c); hub != nil {
				if cc, ok := c.(*models.CustomContext); ok {
					if info, err := cc.GetJwtClaimsInfo(); err == nil {
						hub.Scope().SetUser(sentry.User{
							ID: info.GetUserID(),
						})
						hub.Scope().SetTag("user_id", info.GetUserID())
						hub.Scope().SetTag("role", info.GetRole().String())

					}
				}
			}
			return next(c)
		}
	})
}
