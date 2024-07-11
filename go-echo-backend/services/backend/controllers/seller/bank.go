package controllers

import (
	"github.com/engineeringinflow/inflow-backend/pkg/models"
	"github.com/engineeringinflow/inflow-backend/pkg/repo"
	"github.com/labstack/echo/v4"
	"github.com/rotisserie/eris"
)

// PaginateUserBanks
// @Tags Seller-Bank
// @Summary Paginate setting bank
// @Description Paginate setting bank
// @Accept  json
// @Produce  json
// @Param keyword query string false "Keyword"
// @Param page query int false "Page number"
// @Param limit query int false "Limit"
// @Success 200 {object} query.Pagination{records=[]models.UserBank}
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/seller/banks [get]
func PaginateUserBanks(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, "")
	}
	var params repo.PaginateUserBankParams
	params.UserID = claims.GetUserID()
	params.JwtClaimsInfo = claims

	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, "")
	}

	result := repo.NewUserBankRepo(cc.App.DB).PaginateUserBank(params)
	return cc.Success(result)
}

// CreateUserBanks
// @Tags Seller-Bank
// @Summary CreateFromPayload setting bank
// @Description CreateFromPayload setting bank
// @Accept  json
// @Produce  json
// @Param keyword query string false "Keyword"
// @Param page query int false "Page number"
// @Param limit query int false "Limit"
// @Success 200 {object} models.UserBank
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/seller/banks [post]
func CreateUserBanks(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, "")
	}
	var form models.UserBanksForm
	form.JwtClaimsInfo = claims
	err = cc.BindAndValidate(&form)

	if err != nil {
		return eris.Wrap(err, "")
	}

	result, err := repo.NewUserBankRepo(cc.App.DB).CreateUserBanks(form)
	if err != nil {
		return eris.Wrap(err, "CreateUserBanks - Repo - UpsertUserBank Error")
	}
	return cc.Success(result)
}

// UpdateUserBanks
// @Tags Seller-Bank
// @Summary Update setting bank
// @Description Update setting bank
// @Accept  json
// @Produce  json
// @Param keyword query string false "Keyword"
// @Param page query int false "Page number"
// @Param limit query int false "Limit"
// @Success 200 {object} models.UserBank
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/seller/banks [put]
func UpdateUserBanks(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, "")
	}

	var form models.UserBanksForm
	form.JwtClaimsInfo = claims
	err = cc.BindAndValidate(&form)

	if err != nil {
		return eris.Wrap(err, "")
	}

	result, err := repo.NewUserBankRepo(cc.App.DB).UpdateUserBanks(form)
	if err != nil {
		return eris.Wrap(err, "UpdateUserBanks - Repo - UpsertUserBank Error")
	}
	return cc.Success(result)
}

// DeleteUserBanks
// @Tags Seller-Bank
// @Summary Delete Setting Banks
// @Description Delete Setting Banks
// @Accept  json
// @Produce  json
// @Param keyword query string false "Keyword"
// @Param page query int false "Page number"
// @Param limit query int false "Limit"
// @Success 200 {object} query.Pagination{records=[]models.UserBank}
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/seller/banks/{id} [delete]
func DeleteUserBanks(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, "")
	}

	var form models.DeleteUserBanksForm
	form.JwtClaimsInfo = claims
	err = cc.BindAndValidate(&form)

	if err != nil {
		return eris.Wrap(err, "")
	}

	err = repo.NewUserBankRepo(cc.App.DB).DeleteUserBanks(form)
	if err != nil {
		return eris.Wrap(err, "DeleteUserBanks - Repo - DeleteUserBanks Error")
	}
	return cc.Success("Deleted")
}

// DeleteUserBanksByCountryCode
// @Tags Seller-Bank
// @Summary Delete Setting Banks By Country Code
// @Description Delete Setting Banks By Country Code
// @Accept  json
// @Produce  json
// @Param keyword query string false "Keyword"
// @Param page query int false "Page number"
// @Param limit query int false "Limit"
// @Success 200 {object} query.Pagination{records=[]models.UserBank}
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/seller/banks/countries/{country_code} [delete]
func DeleteUserBanksByCountryCode(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, "DeleteUserBanksByCountryCode - Jwt Claims Error")
	}

	var form models.DeleteUserBanksByCountryCodeForm
	err = cc.BindAndValidate(&form)
	if err != nil {
		return eris.Wrap(err, "")
	}

	form.JwtClaimsInfo = claims
	result, err := repo.NewUserBankRepo(cc.App.DB).DeleteUserBanksByCountryCode(form)
	if err != nil {
		return eris.Wrap(err, "Repo - DeleteUserBanksByCountryCode Error")
	}
	return cc.Success(result)
}
