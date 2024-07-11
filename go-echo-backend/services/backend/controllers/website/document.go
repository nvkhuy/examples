package controllers

import (
	"github.com/engineeringinflow/inflow-backend/pkg/models"
	"github.com/engineeringinflow/inflow-backend/pkg/models/enums"
	"github.com/engineeringinflow/inflow-backend/pkg/repo"
	"github.com/labstack/echo/v4"
	"github.com/rotisserie/eris"
)

// GetDocuments
// @Tags Marketplace-Documents
// @Summary Get Documents
// @Description Get Documents
// @Accept  json
// @Produce  json
// @Param page query int false
// @Success 200 {object} models.Documents
// @Failure 404 {object} errs.Error
// @Router /api/v1/documents [get]
func GetDocumentList(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params models.GetDocumentListParams

	var err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}
	params.Roles = []string{enums.RoleSeller.String(), enums.RoleClient.String()}
	params.Statuses = []enums.DocumentStatus{enums.DocumentStatusPublished}

	documents := repo.NewDocumentRepo(cc.App.DB).GetDocumentList(&params)

	return cc.Success(documents)
}

// @Tags Website-Document
// @Summary Get document detail
// @Description This API allows admin to retrieve document detail
// @Accept  json
// @Produce  json
// @Param slug path string true
// @Success 200 {object} models.Document
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/documents/:slug [get]
func GetDocument(c echo.Context) error {
	var cc, ok = c.(*models.CustomContext)
	if !ok {
		return eris.New("invalid custom context")
	}
	var params models.GetDocumentParams
	var err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}
	resp, err := repo.NewDocumentRepo(cc.App.DB).GetDocument(&params)
	if err != nil {
		return err
	}

	return cc.Success(resp)
}

// GetDocumentCategoryList
// @Tags Website-DocumentCategory
// @Summary Get Document category list
// @Description Get Document category list
// @Accept  json
// @Produce  json
// @Param page query int false
// @Success 200 {object} []models.DocumentCategory
// @Failure 404 {object} errs.Error
// @Router /api/v1/document_categories [get]
func GetDocumentCategoryList(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params models.GetDocumentCategoryListParams

	var err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	documents := repo.NewDocumentCategoryRepo(cc.App.DB).GetDocumentCategoryList(&params)

	return cc.Success(documents)
}
