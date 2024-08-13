package controller

import (
	"github.com/labstack/echo/v4"
	"leaderboard-app/common"
	"leaderboard-app/handler"
	"leaderboard-app/models"
	"net/http"
)

func UserRankings(c echo.Context) error {
	cc := c.(*common.CustomContext)
	var params models.PaginateUserParams
	if err := cc.BindAndValidate(&params); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	rankings, err := handler.UserRankings()
	if err != nil {
		cc.Logger().Error("Failed to fetch user rankings: ", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Internal server error")
	}
	return cc.OkPagination(rankings, params.PaginateParams)
}

// UpdateUser update user
func UpdateUser(c echo.Context) error {
	cc := c.(*common.CustomContext)

	var payload models.UpdateUserPayload
	if err := cc.BindAndValidate(&payload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	user, err := handler.UpdateUser(payload)
	if err != nil {
		cc.Logger().Error("Failed to fetch user rankings: ", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Internal server error")
	}
	// Respond with the updated user data
	return cc.JSON(http.StatusOK, user)
}
