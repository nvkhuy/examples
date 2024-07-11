package controllers

import (
	"github.com/engineeringinflow/inflow-backend/pkg/models"
	"github.com/engineeringinflow/inflow-backend/pkg/repo"
	"github.com/engineeringinflow/inflow-backend/services/consumer/tasks"
	"github.com/labstack/echo/v4"
	"github.com/rotisserie/eris"
)

// PaginateSettingTax
// @Tags Admin-PO
// @Summary Paginate setting tax
// @Description Paginate setting tax
// @Accept  json
// @Produce  json
// @Param keyword query string false "Keyword"
// @Param page query int false "Page number"
// @Param limit query int false "Limit"
// @Success 200 {object} query.Pagination{records=[]models.SettingTax}
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/settings/taxes [get]
func PaginateSettingTax(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	_, err := cc.GetJwtClaims()
	if err != nil {
		return eris.Wrap(err, "")
	}

	var params repo.PaginateSettingTaxParams
	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, "")
	}

	result := repo.NewSettingTaxRepo(cc.App.DB).PaginateSettingTaxes(params)
	return cc.Success(result)
}

// CreateSettingTax
// @Tags Admin-PO
// @Summary CreateFromPayload setting tax
// @Description CreateFromPayload setting tax
// @Accept  json
// @Produce  json
// @Param keyword query string false "Keyword"
// @Param page query int false "Page number"
// @Param limit query int false "Limit"
// @Success 200 {object} []models.SettingTax
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/settings/taxes [post]
func CreateSettingTax(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	_, err := cc.GetJwtClaims()
	if err != nil {
		return eris.Wrap(err, "CreateSettingTaxes - Jwt Claims Error")
	}

	var form models.CreateSettingTaxForm
	err = cc.BindAndValidate(&form)
	if err != nil {
		return eris.Wrap(err, "")
	}

	result, err := repo.NewSettingTaxRepo(cc.App.DB).CreateSettingTax(form)
	if err != nil {
		return eris.Wrap(err, "CreateSettingTaxes - Repo - UpsertSettingTax Error")
	}
	return cc.Success(result)
}

// CreateSettingTax
// @Tags Admin-PO
// @Summary CreateFromPayload setting tax
// @Description CreateFromPayload setting tax
// @Accept  json
// @Produce  json
// @Param keyword query string false "Keyword"
// @Param page query int false "Page number"
// @Param limit query int false "Limit"
// @Success 200 {object} []models.SettingTax
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/settings/taxes/{tax_id} [get]
func GetSettingTax(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	_, err := cc.GetJwtClaims()
	if err != nil {
		return eris.Wrap(err, "CreateSettingTaxes - Jwt Claims Error")
	}

	var form models.GetSettingTaxForm
	err = cc.BindAndValidate(&form)
	if err != nil {
		return eris.Wrap(err, "")
	}

	result, err := repo.NewSettingTaxRepo(cc.App.DB).GetSettingTax(form)
	if err != nil {
		return eris.Wrap(err, "CreateSettingTaxes - Repo - UpsertSettingTax Error")
	}
	return cc.Success(result)
}

// UpdateSettingTax
// @Tags Admin-PO
// @Summary Update setting tax
// @Description Update setting tax
// @Accept  json
// @Produce  json
// @Param keyword query string false "Keyword"
// @Param page query int false "Page number"
// @Param limit query int false "Limit"
// @Success 200 {object} []models.SettingTax
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/settings/taxes/{tax_id} [put]
func UpdateSettingTax(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	_, err := cc.GetJwtClaims()
	if err != nil {
		return eris.Wrap(err, "UpdateSettingTaxes - Jwt Claims Error")
	}

	var form models.UpdateSettingTaxForm
	err = cc.BindAndValidate(&form)
	if err != nil {
		return eris.Wrap(err, "")
	}

	result, err := repo.NewSettingTaxRepo(cc.App.DB).UpdateSettingTax(form)
	if err != nil {
		return eris.Wrap(err, "UpdateSettingTaxes - Repo - UpsertSettingTax Error")
	}
	return cc.Success(result)
}

// DeleteSettingTaxes
// @Tags Admin-PO
// @Summary Delete Setting Taxes
// @Description Delete Setting Taxes
// @Accept  json
// @Produce  json
// @Param keyword query string false "Keyword"
// @Param page query int false "Page number"
// @Param limit query int false "Limit"
// @Success 200 {object} query.Pagination{records=[]models.SettingTax}
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/settings/taxes [delete]
func DeleteSettingTaxes(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	_, err := cc.GetJwtClaims()
	if err != nil {
		return eris.Wrap(err, "DeleteSettingTaxes - Jwt Claims Error")
	}

	var form models.DeleteSettingTaxForm
	err = cc.BindAndValidate(&form)

	if err != nil {
		return eris.Wrap(err, "")
	}

	err = repo.NewSettingTaxRepo(cc.App.DB).DeleteSettingTax(form)
	if err != nil {
		return eris.Wrap(err, "DeleteSettingTaxes - Repo - DeleteSettingSize Error")
	}
	return cc.Success("Deleted")
}

// PaginateSettingSizes
// @Tags Admin-PO
// @Summary Paginate Setting Sizes
// @Description Paginate Setting Sizes
// @Accept  json
// @Produce  json
// @Param keyword query string false "Keyword"
// @Param page query int false "Page number"
// @Param limit query int false "Limit"
// @Success 200 {object} query.Pagination{records=[]models.SettingSize}
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/settings/sizes [get]
func PaginateSettingSizes(c echo.Context) error {
	var cc = c.(*models.CustomContext)

	claimsInfo, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, "")
	}

	var params repo.PaginateSettingSizesParams
	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, "")
	}

	params.JwtClaimsInfo = claimsInfo

	var result = repo.NewSettingSizeRepo(cc.App.DB).PaginateSettingSizes(params)

	return cc.Success(result)
}

// CreateSettingSizes
// @Tags Admin-PO
// @Summary CreateFromPayload Setting Sizes
// @Description CreateFromPayload Setting Sizes
// @Accept  json
// @Produce  json
// @Param keyword query string false "Keyword"
// @Param page query int false "Page number"
// @Param limit query int false "Limit"
// @Success 200 {object} query.Pagination{records=[]models.SettingSize}
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/settings/sizes [post]
func CreateSettingSizes(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	_, err := cc.GetJwtClaims()
	if err != nil {
		return eris.Wrap(err, "CreateSettingSize - Jwt Claims Error")
	}

	var form models.SettingSizeCreateForm
	err = cc.BindAndValidate(&form)

	if err != nil {
		return eris.Wrap(err, "")
	}

	result, err := repo.NewSettingSizeRepo(cc.App.DB).CreateSettingSizes(form)
	if err != nil {
		return eris.Wrap(err, "CreateSettingSize - Repo - CreateSettingSize Error")
	}
	return cc.Success(result)
}

// UpdateSettingSize
// @Tags Admin-PO
// @Summary Update Setting Sizes
// @Description Update Setting Sizes
// @Accept  json
// @Produce  json
// @Param keyword query string false "Keyword"
// @Param page query int false "Page number"
// @Param limit query int false "Limit"
// @Success 200 {object} query.Pagination{records=[]models.SettingSize}
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/settings/sizes/{size_chart_id} [put]
func UpdateSettingSize(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	_, err := cc.GetJwtClaims()
	if err != nil {
		return eris.Wrap(err, "UpdateSettingSize - Jwt Claims Error")
	}

	var form models.SettingSizeUpdateForm
	err = cc.BindAndValidate(&form)

	if err != nil {
		return eris.Wrap(err, "")
	}

	result, err := repo.NewSettingSizeRepo(cc.App.DB).UpdateSettingSize(form)
	if err != nil {
		return eris.Wrap(err, "UpdateSettingSize - Repo - UpdateSettingSize Error")
	}

	return cc.Success(result)
}

// DeleteSettingSize
// @Tags Admin-PO
// @Summary Update Setting Sizes
// @Description Update Setting Sizes
// @Accept  json
// @Produce  json
// @Param keyword query string false "Keyword"
// @Param page query int false "Page number"
// @Param limit query int false "Limit"
// @Success 200 {object} query.Pagination{records=[]models.SettingSize}
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/settings/sizes/{size_chart_id} [delete]
func DeleteSettingSize(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	_, err := cc.GetJwtClaims()
	if err != nil {
		return eris.Wrap(err, "DeleteSettingSize - Jwt Claims Error")
	}

	var form models.SettingSizeDeleteForm
	err = cc.BindAndValidate(&form)

	if err != nil {
		return eris.Wrap(err, "")
	}

	err = repo.NewSettingSizeRepo(cc.App.DB).DeleteSettingSize(form)
	if err != nil {
		return eris.Wrap(err, "DeleteSettingSize - Repo - DeleteSettingSize Error")
	}
	return cc.Success("Updated")
}

// GetSettingSize
// @Tags Admin-PO
// @Summary Update Setting Sizes
// @Description Update Setting Sizes
// @Accept  json
// @Produce  json
// @Param keyword query string false "Keyword"
// @Param page query int false "Page number"
// @Param limit query int false "Limit"
// @Success 200 {object} query.Pagination{records=[]models.SettingSize}
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/settings/sizes/{size_chart_id} [delete]
func GetSettingSize(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	_, err := cc.GetJwtClaims()
	if err != nil {
		return eris.Wrap(err, "GetSettingSize - Jwt Claims Error")
	}

	var form models.SettingSizeIDForm
	err = cc.BindAndValidate(&form)

	if err != nil {
		return eris.Wrap(err, "")
	}

	result, err := repo.NewSettingSizeRepo(cc.App.DB).GetSettingSize(form)
	if err != nil {
		return eris.Wrap(err, "GetSettingSize - Repo - GetSettingSize Error")
	}
	return cc.Success(result)
}

// GetSettingSizeType
// @Tags Admin-PO
// @Summary Delete Setting Size
// @Description Delete Setting Size
// @Accept  json
// @Produce  json
// @Param keyword query string false "Keyword"
// @Param page query int false "Page number"
// @Param limit query int false "Limit"
// @Success 200 {object} query.Pagination{records=[]models.SettingSize}
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/settings/sizes/types/{type} [get]
func GetSettingSizeType(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	_, err := cc.GetJwtClaims()
	if err != nil {
		return eris.Wrap(err, "GetSettingSizeType - Jwt Claims Error")
	}

	var form models.GetSettingSizeTypeForm
	err = cc.BindAndValidate(&form)
	if err != nil {
		return eris.Wrap(err, "")
	}

	result, err := repo.NewSettingSizeRepo(cc.App.DB).GetSettingSizeType(form)
	if err != nil {
		return eris.Wrap(err, "GetSettingSizeType - Repo - GetSettingSizeType Error")
	}
	return cc.Success(result)
}

// UpdateSettingSizeType
// @Tags Admin-PO
// @Summary Delete Setting Size
// @Description Delete Setting Size
// @Accept  json
// @Produce  json
// @Param keyword query string false "Keyword"
// @Param page query int false "Page number"
// @Param limit query int false "Limit"
// @Success 200 {object} query.Pagination{records=[]models.SettingSize}
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/settings/sizes/types/{type} [put]
func UpdateSettingSizeType(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	_, err := cc.GetJwtClaims()
	if err != nil {
		return eris.Wrap(err, "UpdateSettingSizeType - Jwt Claims Error")
	}

	var form models.SettingSizeUpdateTypeForm
	err = cc.BindAndValidate(&form)
	if err != nil {
		return eris.Wrap(err, "")
	}
	var result models.SettingSize
	result, err = repo.NewSettingSizeRepo(cc.App.DB).UpdateSettingSizeType(form)
	if err != nil {
		return eris.Wrap(err, "UpdateSettingSizeType - Repo - UpdateSettingSizeType Error")
	}
	return cc.Success(result)
}

// DeleteSettingSizeType
// @Tags Admin-PO
// @Summary Delete Setting Size
// @Description Delete Setting Size
// @Accept  json
// @Produce  json
// @Param keyword query string false "Keyword"
// @Param page query int false "Page number"
// @Param limit query int false "Limit"
// @Success 200 {object} query.Pagination{records=[]models.SettingSize}
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/settings/sizes/types/{type} [delete]
func DeleteSettingSizeType(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	_, err := cc.GetJwtClaims()
	if err != nil {
		return eris.Wrap(err, "DeleteSettingSizeType - Jwt Claims Error")
	}

	var form models.SettingSizeDeleteTypeForm
	err = cc.BindAndValidate(&form)
	if err != nil {
		return eris.Wrap(err, "")
	}

	err = repo.NewSettingSizeRepo(cc.App.DB).DeleteSettingSizeType(form)
	if err != nil {
		return eris.Wrap(err, "DeleteSettingSizeType - Repo - DeleteSettingSizeType Error")
	}
	return cc.Success("Deleted")
}

// UpdateSettingSizes
// @Tags Admin-PO
// @Summary Delete Setting Size
// @Description Delete Setting Size
// @Accept  json
// @Produce  json
// @Param keyword query string false "Keyword"
// @Param page query int false "Page number"
// @Param limit query int false "Limit"
// @Success 200 {object} query.Pagination{records=[]models.SettingSize}
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/settings/sizes [put]
func UpdateSettingSizes(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	_, err := cc.GetJwtClaims()
	if err != nil {
		return eris.Wrap(err, "UpdateSettingSizes - Jwt Claims Error")
	}

	var form models.SettingSizesUpdateForm
	err = cc.BindAndValidate(&form)
	if err != nil {
		return eris.Wrap(err, "")
	}

	result, err := repo.NewSettingSizeRepo(cc.App.DB).UpdateSettingSizes(form)
	if err != nil {
		return eris.Wrap(err, "UpdateSettingSizes - Repo - UpdateSettingSizes Error")
	}
	return cc.Success(result)
}

// PaginateSettingBanks
// @Tags Admin-PO
// @Summary Paginate setting bank
// @Description Paginate setting bank
// @Accept  json
// @Produce  json
// @Param keyword query string false "Keyword"
// @Param page query int false "Page number"
// @Param limit query int false "Limit"
// @Success 200 {object} query.Pagination{records=[]models.SettingBank}
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/settings/banks [get]
func PaginateSettingBanks(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	_, err := cc.GetJwtClaims()
	if err != nil {
		return eris.Wrap(err, "")
	}

	var params repo.PaginateSettingBankParams
	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, "")
	}

	result := repo.NewSettingBankRepo(cc.App.DB).PaginateSettingBank(params)
	return cc.Success(result)
}

// CreateSettingBanks
// @Tags Admin-PO
// @Summary CreateFromPayload setting bank
// @Description CreateFromPayload setting bank
// @Accept  json
// @Produce  json
// @Param keyword query string false "Keyword"
// @Param page query int false "Page number"
// @Param limit query int false "Limit"
// @Success 200 {object} models.SettingBank
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/settings/banks [post]
func CreateSettingBanks(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	_, err := cc.GetJwtClaims()
	if err != nil {
		return eris.Wrap(err, "CreateSettingBanks - Jwt Claims Error")
	}

	var form models.SettingBanksForm
	err = cc.BindAndValidate(&form)

	if err != nil {
		return eris.Wrap(err, "")
	}

	result, err := repo.NewSettingBankRepo(cc.App.DB).CreateSettingBanks(form)
	if err != nil {
		return eris.Wrap(err, "CreateSettingBanks - Repo - UpsertSettingBank Error")
	}
	return cc.Success(result)
}

// UpdateSettingBanks
// @Tags Admin-PO
// @Summary Update setting bank
// @Description Update setting bank
// @Accept  json
// @Produce  json
// @Param keyword query string false "Keyword"
// @Param page query int false "Page number"
// @Param limit query int false "Limit"
// @Success 200 {object} models.SettingBank
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/settings/banks [put]
func UpdateSettingBanks(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	_, err := cc.GetJwtClaims()
	if err != nil {
		return eris.Wrap(err, "UpdateSettingBanks - Jwt Claims Error")
	}

	var form models.SettingBanksForm
	err = cc.BindAndValidate(&form)

	if err != nil {
		return eris.Wrap(err, "")
	}

	result, err := repo.NewSettingBankRepo(cc.App.DB).UpdateSettingBanks(form)
	if err != nil {
		return eris.Wrap(err, "UpdateSettingBanks - Repo - UpsertSettingBank Error")
	}
	return cc.Success(result)
}

// DeleteSettingBanks
// @Tags Admin-PO
// @Summary Delete Setting Banks
// @Description Delete Setting Banks
// @Accept  json
// @Produce  json
// @Param keyword query string false "Keyword"
// @Param page query int false "Page number"
// @Param limit query int false "Limit"
// @Success 200 {object} query.Pagination{records=[]models.SettingBank}
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/settings/banks/{id} [delete]
func DeleteSettingBanks(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	_, err := cc.GetJwtClaims()
	if err != nil {
		return eris.Wrap(err, "DeleteSettingBanks - Jwt Claims Error")
	}

	var form models.DeleteSettingBanksForm
	err = cc.BindAndValidate(&form)

	if err != nil {
		return eris.Wrap(err, "")
	}

	err = repo.NewSettingBankRepo(cc.App.DB).DeleteSettingBanks(form)
	if err != nil {
		return eris.Wrap(err, "DeleteSettingBanks - Repo - DeleteSettingBanks Error")
	}
	return cc.Success("Deleted")
}

// DeleteSettingBanksByCountryCode
// @Tags Admin-PO
// @Summary Delete Setting Banks By Country Code
// @Description Delete Setting Banks By Country Code
// @Accept  json
// @Produce  json
// @Param keyword query string false "Keyword"
// @Param page query int false "Page number"
// @Param limit query int false "Limit"
// @Success 200 {object} query.Pagination{records=[]models.SettingBank}
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/settings/banks/countries/{country_code} [delete]
func DeleteSettingBanksByCountryCode(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	_, err := cc.GetJwtClaims()
	if err != nil {
		return eris.Wrap(err, "DeleteSettingBanksByCountryCode - Jwt Claims Error")
	}

	var form models.DeleteSettingBanksByCountryCodeForm
	err = cc.BindAndValidate(&form)

	if err != nil {
		return eris.Wrap(err, "")
	}

	result, err := repo.NewSettingBankRepo(cc.App.DB).DeleteSettingBanksByCountryCode(form)
	if err != nil {
		return eris.Wrap(err, "Repo - DeleteSettingBanksByCountryCode Error")
	}
	return cc.Success(result)
}

// PaginateSettingSEO
// @Tags Admin-PO
// @Summary Paginate setting bank
// @Description Paginate setting bank
// @Accept  json
// @Produce  json
// @Param keyword query string false "Keyword"
// @Param page query int false "Page number"
// @Param limit query int false "Limit"
// @Success 200 {object} query.Pagination{records=[]models.SettingBank}
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/settings/seo [get]
func PaginateSettingSEO(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	_, err := cc.GetJwtClaims()
	if err != nil {
		return eris.Wrap(err, "")
	}

	var params repo.PaginateSettingSEOParams
	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, "")
	}

	result := repo.NewSettingSEORepo(cc.App.DB).PaginateSettingSEO(params)
	return cc.Success(result)
}

// PaginateSettingRouteGroupSEO
// @Tags Admin-PO
// @Summary Paginate setting bank
// @Description Paginate setting bank
// @Accept  json
// @Produce  json
// @Param keyword query string false "Keyword"
// @Param page query int false "Page number"
// @Param limit query int false "Limit"
// @Success 200 {object} query.Pagination{records=[]models.SettingBank}
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/settings/seo/route [get]
func PaginateSettingRouteGroupSEO(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	_, err := cc.GetJwtClaims()
	if err != nil {
		return eris.Wrap(err, "")
	}

	var params repo.PaginateSettingSEOParams
	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, "")
	}

	result := repo.NewSettingSEORepo(cc.App.DB).PaginateSettingSEOLanguageGroup(params)
	return cc.Success(result)
}

// CreateSettingSEO
// @Tags Admin-PO
// @Summary CreateFromPayload setting bank
// @Description CreateFromPayload setting bank
// @Accept  json
// @Produce  json
// @Param keyword query string false "Keyword"
// @Param page query int false "Page number"
// @Param limit query int false "Limit"
// @Success 200 {object} models.SettingBank
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/settings/seo [post]
func CreateSettingSEO(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	_, err := cc.GetJwtClaims()
	if err != nil {
		return eris.Wrap(err, "CreateSettingSEO - Jwt Claims Error")
	}

	var form models.CreateSettingSEOForm
	err = cc.BindAndValidate(&form)

	if err != nil {
		return eris.Wrap(err, "")
	}

	result, err := repo.NewSettingSEORepo(cc.App.DB).CreateSettingSEO(form)
	if err != nil {
		return eris.Wrap(err, "CreateSettingSEO - Repo - UpsertSettingBank Error")
	}
	return cc.Success(result)
}

// UpdateSettingSEO
// @Tags Admin-PO
// @Summary Update setting bank
// @Description Update setting bank
// @Accept  json
// @Produce  json
// @Param keyword query string false "Keyword"
// @Param page query int false "Page number"
// @Param limit query int false "Limit"
// @Success 200 {object} models.SettingBank
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/settings/seo [put]
func UpdateSettingSEO(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	_, err := cc.GetJwtClaims()
	if err != nil {
		return eris.Wrap(err, "UpdateSettingSEO - Jwt Claims Error")
	}

	var form models.UpdateSettingSEOForm
	err = cc.BindAndValidate(&form)

	if err != nil {
		return eris.Wrap(err, "")
	}

	result, err := repo.NewSettingSEORepo(cc.App.DB).UpdateSettingSEO(form)
	if err != nil {
		return eris.Wrap(err, "UpdateSettingSEO - Repo - UpsertSettingBank Error")
	}
	return cc.Success(result)
}

// DeleteSettingSEO
// @Tags Admin-PO
// @Summary Delete Setting SEO
// @Description Delete Setting SEO
// @Accept  json
// @Produce  json
// @Param keyword query string false "Keyword"
// @Param page query int false "Page number"
// @Param limit query int false "Limit"
// @Success 200 {object} query.Pagination{records=[]models.SettingBank}
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/settings/seo/{id} [delete]
func DeleteSettingSEO(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	_, err := cc.GetJwtClaims()
	if err != nil {
		return eris.Wrap(err, "DeleteSettingSEO - Jwt Claims Error")
	}

	var form models.DeleteSettingSEOForm
	err = cc.BindAndValidate(&form)

	if err != nil {
		return eris.Wrap(err, "")
	}

	err = repo.NewSettingSEORepo(cc.App.DB).DeleteSettingSEOs(form)
	if err != nil {
		return eris.Wrap(err, "DeleteSettingSEO - Repo - DeleteSettingSEO Error")
	}
	return cc.Success("Deleted")
}

// GetSettingSEOByRouteName
// @Tags Admin-PO
// @Summary Delete Setting SEO
// @Description Delete Setting SEO
// @Accept  json
// @Produce  json
// @Param keyword query string false "Keyword"
// @Param page query int false "Page number"
// @Param limit query int false "Limit"
// @Success 200 {object} query.Pagination{records=[]models.SettingBank}
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/settings/seo/route/{route_name} [get]
func GetSettingSEOByRouteName(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	_, err := cc.GetJwtClaims()
	if err != nil {
		return eris.Wrap(err, "GetSettingSEOByRouteName - Jwt Claims Error")
	}

	var params repo.GetSettingSEOByRouteNameParams
	err = cc.BindAndValidate(&params)

	if err != nil {
		return eris.Wrap(err, "")
	}

	result, err := repo.NewSettingSEORepo(cc.App.DB).GetSettingSEOByRouteName(params)
	if err != nil {
		return eris.Wrap(err, "GetSettingSEOByRouteName - Repo - GetSettingSEOByRouteName Error")
	}
	return cc.Success(result)
}

// PatchSEOTranslations
// @Tags Admin-PO
// @Summary Delete Setting SEO
// @Description Delete Setting SEO
// @Accept  json
// @Produce  json
// @Param keyword query string false "Keyword"
// @Param page query int false "Page number"
// @Param limit query int false "Limit"
// @Success 200 {object} query.Pagination{records=[]models.SettingBank}
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/settings/seo/translations [patch]
func PatchSEOTranslations(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	_, err := cc.GetJwtClaims()
	if err != nil {
		return eris.Wrap(err, "DeleteSettingSEO - Jwt Claims Error")
	}

	var form models.FetchSeoTranslationParams
	err = cc.BindAndValidate(&form)

	if err != nil {
		return eris.Wrap(err, "")
	}

	result, err := repo.NewSeoTranslationRepo(cc.App.DB).WithSheetAPI(cc.App.SheetAPI).FetchSeoTranslation(&form)
	if err != nil {
		return eris.Wrap(err, "PatchSEOTranslations - Repo - DeleteSettingSEO Error")
	}
	return cc.Success(result)
}

// GetSettingDoc
// @Tags Admin-PO
// @Summary Delete Setting Size
// @Description Delete Setting Size
// @Accept  json
// @Produce  json
// @Param keyword query string false "Keyword"
// @Param page query int false "Page number"
// @Param limit query int false "Limit"
// @Success 200 {object} query.Pagination{records=[]models.SettingSize}
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/settings/docs/{type} [put]
func GetSettingDoc(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, err.Error())
	}
	var params repo.GetSettingDocParams
	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	params.JwtClaimsInfo = claims
	var result *models.SettingDoc
	result, err = repo.NewSettingDocRepo(cc.App.DB).Get(params)
	if err != nil {
		return eris.Wrap(err, "")
	}
	return cc.Success(result)
}

// CreateSettingDoc
// @Tags Admin-PO
// @Summary Delete Setting Size
// @Description Delete Setting Size
// @Accept  json
// @Produce  json
// @Param keyword query string false "Keyword"
// @Param page query int false "Page number"
// @Param limit query int false "Limit"
// @Success 200 {object} query.Pagination{records=[]models.SettingSize}
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/settings/doc [put]
func CreateSettingDoc(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, err.Error())
	}
	var params repo.SettingDocCreateParams
	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	params.JwtClaimsInfo = claims
	var result *models.SettingDoc
	result, err = repo.NewSettingDocRepo(cc.App.DB).Create(params)
	if err != nil {
		return eris.Wrap(err, "")
	}
	return cc.Success(result)
}

// UpdateSettingDoc
// @Tags Admin-PO
// @Summary Delete Setting Size
// @Description Delete Setting Size
// @Accept  json
// @Produce  json
// @Param keyword query string false "Keyword"
// @Param page query int false "Page number"
// @Param limit query int false "Limit"
// @Success 200 {object} query.Pagination{records=[]models.SettingSize}
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/settings/doc/{type} [put]
func UpdateSettingDoc(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, err.Error())
	}
	var params repo.SettingDocUpdateParams
	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	params.JwtClaimsInfo = claims
	var result *models.SettingDoc
	result, err = repo.NewSettingDocRepo(cc.App.DB).Update(params)
	if err != nil {
		return eris.Wrap(err, "")
	}
	_, _ = tasks.DeleteUserDocAgreementTask{
		SettingDocType: params.Type,
	}.Dispatch(c.Request().Context())

	return cc.Success(result)
}

// CreateSettingInquiry
// @Tags Admin-PO
// @Summary Delete Setting Size
// @Description Delete Setting Size
// @Accept  json
// @Produce  json
// @Param keyword query string false "Keyword"
// @Param page query int false "Page number"
// @Param limit query int false "Limit"
// @Success 200 {object} query.Pagination{records=[]models.SettingSize}
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/settings/inquiries [put]
func CreateSettingInquiry(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, err.Error())
	}
	var params repo.SettingInquiryCreateParams
	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	params.JwtClaimsInfo = claims
	var result *models.SettingInquiry
	result, err = repo.NewSettingInquiryRepo(cc.App.DB).Create(params)
	if err != nil {
		return eris.Wrap(err, "")
	}
	return cc.Success(result)
}

// GetSettingInquiry
// @Tags Admin-PO
// @Summary Delete Setting Size
// @Description Delete Setting Size
// @Accept  json
// @Produce  json
// @Param keyword query string false "Keyword"
// @Param page query int false "Page number"
// @Param limit query int false "Limit"
// @Success 200 {object} query.Pagination{records=[]models.SettingSize}
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/settings/inquiries [get]
func GetSettingInquiry(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, err.Error())
	}
	var params repo.GetSettingInquiryParams
	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	params.JwtClaimsInfo = claims
	var result *models.SettingInquiry
	result, err = repo.NewSettingInquiryRepo(cc.App.DB).Get(params)
	if err != nil {
		return eris.Wrap(err, "")
	}
	return cc.Success(result)
}

// UpdateSettingInquiry
// @Tags Admin-PO
// @Summary Delete Setting Size
// @Description Delete Setting Size
// @Accept  json
// @Produce  json
// @Param keyword query string false "Keyword"
// @Param page query int false "Page number"
// @Param limit query int false "Limit"
// @Success 200 {object} query.Pagination{records=[]models.SettingSize}
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/settings/doc/{type} [put]
func UpdateSettingInquiry(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, err.Error())
	}
	var params repo.SettingInquiryUpdateParams
	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	params.JwtClaimsInfo = claims
	var result *models.SettingInquiry
	result, err = repo.NewSettingInquiryRepo(cc.App.DB).Update(params)
	if err != nil {
		return eris.Wrap(err, "")
	}

	return cc.Success(result)
}
