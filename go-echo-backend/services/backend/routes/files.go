package routes

import (
	"net/http"

	"github.com/engineeringinflow/inflow-backend/pkg/middlewares"
	controllers "github.com/engineeringinflow/inflow-backend/services/backend/controllers/files"
	"github.com/labstack/echo/v4"
)

func (router *Router) SetupFilesRoutes(g *echo.Group) {
	var fs = http.FileServer(http.Dir(router.App.Config.EFSPath))

	router.Echo.GET("/files/*", echo.WrapHandler(fs), middlewares.IsAuthorizedWithQueryToken(router.App.Config.JWTAssetSecret), router.Middlewares.ResponseDownloadableHeader())
	router.Echo.GET("/fs", controllers.WalkFiles, middlewares.IsBasicAuth())
	router.Echo.GET("/fs/:file_name/delete", controllers.DeleteFile, middlewares.IsBasicAuth())

}
