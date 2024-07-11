package controllers

import (
	"github.com/engineeringinflow/inflow-backend/pkg/models"
	"github.com/engineeringinflow/inflow-backend/pkg/repo"
	"github.com/labstack/echo/v4"
	"github.com/rotisserie/eris"
)

// @Tags Admin-DocumentCategory
// @Summary Create document category
// @Description This API allows admin to create document category
// @Accept  json
// @Produce  json
// @Param data body models.CreateDocumentCategoryRequest true
// @Success 200 {object} models.DocumentCategory
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/document_categories [post]
func AdminCreateDocumentCategory(c echo.Context) error {
	var cc, ok = c.(*models.CustomContext)
	if !ok {
		return eris.New("invalid custom context")
	}
	var req models.CreateDocumentCategoryRequest

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, "get jwt claim info error")
	}

	err = cc.BindAndValidate(&req)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}
	req.JwtClaimsInfo = claims
	resp, err := repo.NewDocumentCategoryRepo(cc.App.DB).CreateDocumentCategory(&req)
	if err != nil {
		return err
	}

	return cc.Success(resp)
}

// @Tags Admin-DocumentCategory
// @Summary Update document category
// @Description This API allows admin to update document category
// @Accept  json
// @Produce  json
// @Param data body models.UpdateDocumentCategoryRequest true
// @Param document_category_id path string true
// @Success 200 {object} models.DocumentCategory
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/document_categories/:document_category_id [put]
func AdminUpdateDocumentCategory(c echo.Context) error {
	var cc, ok = c.(*models.CustomContext)
	if !ok {
		return eris.New("invalid custom context")
	}
	var req models.UpdateDocumentCategoryRequest

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, "get jwt claim info error")
	}

	err = cc.BindAndValidate(&req)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}
	req.JwtClaimsInfo = claims
	resp, err := repo.NewDocumentCategoryRepo(cc.App.DB).UpdateDocumentCategory(&req)
	if err != nil {
		return err
	}

	return cc.Success(resp)
}

// @Tags Admin-DocumentCategory
// @Summary Get document category list
// @Description This API allows admin to list document category with pagination
// @Accept  json
// @Produce  json
// @Param page query int false
// @Param limit query int false
// @Param keyword query string false
// @Success 200 {object} query.Pagination{Records=[]models.DocumentCategory}
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/document_categories [get]
func AdminGetDocumentCategoryList(c echo.Context) error {
	var cc, ok = c.(*models.CustomContext)
	if !ok {
		return eris.New("invalid custom context")
	}
	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, "get jwt claim info error")
	}
	var params models.GetDocumentCategoryListParams
	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}
	params.JwtClaimsInfo = claims
	resp := repo.NewDocumentCategoryRepo(cc.App.DB).GetDocumentCategoryList(&params)

	return cc.Success(resp)
}

// @Tags Admin-DocumentCategory
// @Summary Delete document category
// @Description This API allows admin to archive document category
// @Accept  json
// @Produce  json
// @Param document_category_id path string true
// @Success 200 {object} struct{message string}
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/document_categories/:document_category_id [delete]
func AdminDeleteDocumentCategory(c echo.Context) error {
	var cc, ok = c.(*models.CustomContext)
	if !ok {
		return eris.New("invalid custom context")
	}
	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, "get jwt claim info error")
	}
	var params models.GetDocumentCategoryParams
	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}
	params.JwtClaimsInfo = claims
	err = repo.NewDocumentCategoryRepo(cc.App.DB).DeleteDocumentCategory(&params)
	if err != nil {
		return err
	}

	return cc.Success("Archived")
}
