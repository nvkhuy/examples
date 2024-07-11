package routes

import (
	// docs

	"os"

	"github.com/engineeringinflow/inflow-backend/pkg/config"
	_ "github.com/engineeringinflow/inflow-backend/services/backend/docs"
	"github.com/mvrilo/go-redoc"

	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
)

// SetupDocRoute setup document route
func (router *Router) SetupDocRoute(g *echo.Group) {
	g.GET("/*", echoSwagger.EchoWrapHandler(func(c *echoSwagger.Config) {
		c.URL = "./swagger/doc.json"
	}))

	router.GET("/api/docs/swagger/doc.json", echoSwagger.EchoWrapHandler(func(c *echoSwagger.Config) {
		c.URL = "./swagger/doc.json"
	}))
}

func (router *Router) SetupReDocRoute(g *echo.Group) {
	var specFile = "services/backend/docs/swagger.json"
	if !router.App.DB.Configuration.IsLocal() {
		specFile = "./docs/swagger.json"
	}

	if _, err := os.Stat(specFile); err == nil {
		var doc = redoc.Redoc{
			Title:       router.App.DB.Configuration.GetServerName(config.ServiceBackend),
			Description: router.App.DB.Configuration.GetDocDescription(),
			SpecFile:    specFile,
			SpecPath:    "/api/docs/swagger/doc.json",
			DocsPath:    "/api/docs",
		}

		var handle = doc.Handler()
		g.GET(doc.DocsPath, func(c echo.Context) error {
			handle(c.Response(), c.Request())
			return nil
		})
	}

}
