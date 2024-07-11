package models

import (
	"github.com/engineeringinflow/inflow-backend/pkg/app"
	"github.com/engineeringinflow/inflow-backend/pkg/logger"
	"github.com/labstack/echo/v4"
)

// CustomContext custom echo context
type CustomContext struct {
	echo.Context
	CustomLogger *logger.Logger
	App          *app.App
}
