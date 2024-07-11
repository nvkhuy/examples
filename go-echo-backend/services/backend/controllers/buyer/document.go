package controllers

import (
	"github.com/engineeringinflow/inflow-backend/pkg/models"
	"github.com/engineeringinflow/inflow-backend/pkg/models/enums"
	"github.com/engineeringinflow/inflow-backend/pkg/repo"
	"github.com/labstack/echo/v4"
	"github.com/rotisserie/eris"
)

// @Tags Buyer-Document
// @Summary Get document list
// @Description This API allows Buyer to list document with pagination
// @Accept  json
// @Produce  json
// @Param page query int false
// @Param limit query int false
// @Param keyword query string false
// @Param status query string false
// @Param category query string false
// @Success 200 {object} query.Pagination{Records=[]models.Document}
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/buyer/documents/ [get]
func BuyerGetDocumentList(c echo.Context) error {
	var cc, ok = c.(*models.CustomContext)
	if !ok {
		return eris.New("invalid custom context")
	}
	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, "get jwt claim info error")
	}
	var params models.GetDocumentListParams
	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}
	params.JwtClaimsInfo = claims
	if claims.GetRole().IsBuyer() {
		params.Roles = []string{enums.RoleClient.String()}
	}
	resp := repo.NewDocumentRepo(cc.App.DB).GetDocumentList(&params)

	return cc.Success(resp)
}
