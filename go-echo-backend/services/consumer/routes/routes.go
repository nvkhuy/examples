package routes

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/engineeringinflow/inflow-backend/services/consumer/docs"

	"github.com/engineeringinflow/inflow-backend/pkg/app"
	"github.com/engineeringinflow/inflow-backend/pkg/config"
	"github.com/engineeringinflow/inflow-backend/pkg/middlewares"
	"github.com/engineeringinflow/inflow-backend/pkg/ws"
	"github.com/engineeringinflow/inflow-backend/services/consumer"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// Router router
type Router struct {
	*echo.Echo
	Middlewares *middlewares.Middleware
	App         *app.App

	Host string
}

func NewRouter(app *app.App) *Router {
	var e = echo.New()
	var m = middlewares.New(e, app)

	var router = &Router{
		Echo:        e,
		App:         app,
		Middlewares: m,
	}

	return router
}

func (router *Router) ListenAndServe() {
	router.SetupRoutes()

	router.startServer()
}

func (router *Router) startServer() {
	go func() {
		var serverAddress = fmt.Sprintf(":%s", router.App.Config.ServerPort)
		if err := router.Start(serverAddress); err != nil && err != http.ErrServerClosed {
			router.Logger.Fatal("Shutting down the server")
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with a timeout of 10 seconds.
	// Use a buffered channel to avoid missing signals as recommended for signal.Notify
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := router.Shutdown(ctx); err != nil {
		router.Logger.Fatal(err)
	}
}

func (router *Router) SetupRoutes() {
	var apiV1Routes = router.Group("/api/v1")
	apiV1Routes.Use(middleware.Gzip())

	var appInfoGroup = router.Group("/api/app")
	appInfoGroup.Use(router.Middlewares.IsBasicAuth())

	var docGroup = router.Group("/api/docs")
	docGroup.Use(router.Middlewares.IsBasicAuth())

	router.SetupRootRoute(router.Group(""))

	router.SetupAppInfoRoute(appInfoGroup)

	router.setupDocConfig()

	ws.New(router.Echo, router.App, consumer.WSMessageHandler)

}

func (router *Router) setupDocConfig() {
	urlInfo, _ := url.Parse(router.App.Config.ServerBaseURL)

	docs.SwaggerInfo.Title = fmt.Sprintf("%s api", router.App.Config.GetServerName(config.ServiceConsumer))
	docs.SwaggerInfo.Description = router.App.Config.GetDocDescription()
	docs.SwaggerInfo.Version = router.App.Config.BuildVersion
	docs.SwaggerInfo.Host = urlInfo.Host
	docs.SwaggerInfo.Schemes = []string{"http", "https"}

	if urlInfo.Port() != router.App.Config.ServerPort {
		docs.SwaggerInfo.Host = strings.ReplaceAll(docs.SwaggerInfo.Host, urlInfo.Port(), router.App.Config.ServerPort)
	}
	if urlInfo.Scheme == "https" {
		docs.SwaggerInfo.Schemes = []string{"https", "http"}
	}
}
