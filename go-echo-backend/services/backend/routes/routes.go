package routes

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/engineeringinflow/inflow-backend/pkg/app"
	"github.com/engineeringinflow/inflow-backend/pkg/config"
	"github.com/engineeringinflow/inflow-backend/services/backend/docs"

	"github.com/engineeringinflow/inflow-backend/pkg/middlewares"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// Router Echo
type Router struct {
	*echo.Echo
	Middlewares *middlewares.Middleware
	App         *app.App
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
		if err := router.Start(serverAddress); err != nil && !errors.Is(err, http.ErrServerClosed) {
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

	router.SetupDocRoute(docGroup)

	router.SetupReDocRoute(docGroup)

	router.SetupAppInfoRoute(appInfoGroup)

	// Auth
	router.SetupAuthRoutes(apiV1Routes)

	router.SetupBuyerRoutes(apiV1Routes.Group("/buyer"))

	// Admin routes
	router.SetupAdminRoutes(apiV1Routes.Group("/admin"))

	// Common routes
	router.SetupCommonRoutes(apiV1Routes.Group("/common"))

	// Marketplace routes
	router.SetupWebsiteRoutes(apiV1Routes.Group(""))

	// Seller routes
	router.SetupSellerRoutes(apiV1Routes.Group("/seller"))

	// App routes
	router.SetupAppRoutes(apiV1Routes.Group("/app"))

	// Webhook routes
	router.SetupWebhookRoutes(apiV1Routes.Group("/webhook"))

	router.SetupCallbackRoutes(apiV1Routes.Group("/callback"))

	router.SetupFilesRoutes(router.Group(""))

	router.setupDocConfig()

}

func (router *Router) setupDocConfig() {
	var segments = strings.Split(router.App.Config.ServerBaseURL, "://")
	docs.SwaggerInfo.Title = fmt.Sprintf("%s api", router.App.Config.GetServerName(config.ServiceBackend))
	docs.SwaggerInfo.Description = router.App.Config.GetDocDescription()
	docs.SwaggerInfo.Version = router.App.Config.BuildVersion
	docs.SwaggerInfo.Host = segments[1]
	docs.SwaggerInfo.Schemes = []string{"http", "https"}

	if strings.HasPrefix(segments[0], "https") {
		docs.SwaggerInfo.Schemes = []string{"https", "http"}
	}
}
