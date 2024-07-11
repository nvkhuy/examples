package controllers

import (
	"github.com/engineeringinflow/inflow-backend/pkg/models"
	"github.com/engineeringinflow/inflow-backend/pkg/repo"
	"github.com/labstack/echo/v4"
	"github.com/rotisserie/eris"
)

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

// GetDocsAgreement
// @Tags Marketplace-Support
// @Summary Search SupportTicket
// @Description Search SupportTicket
// @Accept  json
// @Produce  json
// @Param keyword query string false "Keyword"
// @Param page query int false "Page number"
// @Success 200 {object} models.SupportTicket
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/buyer/docs/:setting_doc_type/agreement [put]
func GetDocsAgreement(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.GetUserDocAgreementParams

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	if err = cc.BindAndValidate(&params); err != nil {
		return eris.Wrap(err, err.Error())
	}

	params.JwtClaimsInfo = claims
	result, err := repo.NewUserDocAgreementRepo(cc.App.DB).Get(params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	return cc.Success(result)
}

// CreateDocsAgreement
// @Tags Marketplace-Support
// @Summary Search SupportTicket
// @Description Search SupportTicket
// @Accept  json
// @Produce  json
// @Param keyword query string false "Keyword"
// @Param page query int false "Page number"
// @Success 200 {object} models.SupportTicket
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/buyer/docs/:setting_doc_type/agreement [post]
func CreateDocsAgreement(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.CreateUserDocAgreementParams

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	if err = cc.BindAndValidate(&params); err != nil {
		return eris.Wrap(err, err.Error())
	}

	params.JwtClaimsInfo = claims
	result, err := repo.NewUserDocAgreementRepo(cc.App.DB).Create(params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	return cc.Success(result)
}
