package main

import (
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"leaderboard-app/common"
	"leaderboard-app/controller"
	"leaderboard-app/middleware"
)

func main() {
	// Create a new Echo instance
	e := echo.New()
	e.Validator = &common.CustomValidator{Validator: validator.New()}

	e.Use(middleware.UseCORS(), middleware.UseCustomContext)

	e.GET("/users/rankings", controller.UserRankings)
	e.PUT("/users/:id", controller.UpdateUser)

	e.Logger.Fatal(e.Start(":8080"))
}
