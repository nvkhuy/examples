package controllerss

import (
	"github.com/engineeringinflow/inflow-backend/pkg/models"
	"github.com/engineeringinflow/inflow-backend/pkg/repo"
	"github.com/labstack/echo/v4"
	"github.com/rotisserie/eris"
)

// Login by email
// @Tags App-Auth
// @Summary Login by email
// @Description Login by email
// @Accept  json
// @Produce  json
// @Param data body models.LoginEmailForm true "Login Form"
// @Success 200 {object} models.LoginResponse
// @Header 200 {string} Bearer YOUR_TOKEN
// @Failure 404 {object} errs.Error
// @Router /api/v1/app/auth/login_email [post]
func LoginEmail(c echo.Context) error {
	var form models.LoginEmailForm
	var cc = c.(*models.CustomContext)

	var err = cc.BindAndValidate(&form)

	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	form.IsSeller = true
	response, err := repo.NewAuthRepo(cc.App.DB).LoginEmail(form)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	return cc.Success(response)

}
