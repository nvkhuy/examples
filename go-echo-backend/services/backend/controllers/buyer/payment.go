package controllers

import (
	"github.com/engineeringinflow/inflow-backend/pkg/models"
	"github.com/engineeringinflow/inflow-backend/pkg/repo"
	"github.com/labstack/echo/v4"
	"github.com/rotisserie/eris"
)

// SetupPayment Setup payment
// @Tags Marketplace-Payment-Method
// @Summary Setup payment
// @Description Setup payment
// @Accept  json
// @Produce  json
// @Header 200 {string} Bearer YOUR_TOKEN
// @Failure 404 {object} errs.Error
// @Security ApiKeyAuth
// @Router /api/v1/setup_payment [post]
func SetupPayment(c echo.Context) error {
	var cc = c.(*models.CustomContext)

	claims, err := cc.GetJwtClaims()
	if err != nil {
		return eris.Wrap(err, "")
	}

	result, err := repo.NewUserRepo(cc.App.DB).SetupIntentForClientSecert(claims.ID)
	if err != nil {
		return eris.Wrap(err, "")
	}

	return cc.Success(result)
}
