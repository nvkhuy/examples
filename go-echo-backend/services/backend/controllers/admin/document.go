package controllers

import (
	"github.com/engineeringinflow/inflow-backend/pkg/models"
	"github.com/engineeringinflow/inflow-backend/pkg/repo"
	"github.com/labstack/echo/v4"
	"github.com/rotisserie/eris"
)

// @Tags Admin-Document
// @Summary Create document
// @Description This API allows admin to create document
// @Accept  json
// @Produce  json
// @Param data body models.CreateDocumentRequest true
// @Success 200 {object} models.Document
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/documents [post]
func AdminCreateDocument(c echo.Context) error {
	var cc, ok = c.(*models.CustomContext)
	if !ok {
		return eris.New("invalid custom context")
	}
	var req models.CreateDocumentRequest

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, "get jwt claim info error")
	}

	err = cc.BindAndValidate(&req)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}
	req.JwtClaimsInfo = claims
	req.UserID = claims.GetUserID()
	resp, err := repo.NewDocumentRepo(cc.App.DB).CreateDocument(&req)
	if err != nil {
		return err
	}

	return cc.Success(resp)
}

// @Tags Admin-Document
// @Summary Update document
// @Description This API allows admin to update existing document
// @Accept  json
// @Produce  json
// @Param data body models.UpdateDocumentRequest true
// @Param document_id path string true
// @Success 200 {object} models.Document
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/documents/:document_id [put]
func AdminUpdateDocument(c echo.Context) error {
	var cc, ok = c.(*models.CustomContext)
	if !ok {
		return eris.New("invalid custom context")
	}
	var req models.UpdateDocumentRequest

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, "get jwt claim info error")
	}

	err = cc.BindAndValidate(&req)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}
	req.JwtClaimsInfo = claims
	resp, err := repo.NewDocumentRepo(cc.App.DB).UpdateDocument(&req)
	if err != nil {
		return err
	}

	return cc.Success(resp)
}

// @Tags Admin-Document
// @Summary Get document list
// @Description This API allows admin to list document with pagination
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
// @Router /api/v1/admin/documents/ [get]
func AdminGetDocumentList(c echo.Context) error {
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
	resp := repo.NewDocumentRepo(cc.App.DB).GetDocumentList(&params)

	return cc.Success(resp)
}

// @Tags Admin-Document
// @Summary Get document detail
// @Description This API allows admin to retrieve document detail
// @Accept  json
// @Produce  json
// @Param slug path string true
// @Success 200 {object} models.Document
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/documents/:document_id [get]
func AdminGetDocument(c echo.Context) error {
	var cc, ok = c.(*models.CustomContext)
	if !ok {
		return eris.New("invalid custom context")
	}
	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, "get jwt claim info error")
	}
	var params models.GetDocumentParams
	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}
	params.JwtClaimsInfo = claims
	resp, err := repo.NewDocumentRepo(cc.App.DB).GetDocument(&params)
	if err != nil {
		return err
	}

	return cc.Success(resp)
}

// @Tags Admin-Document
// @Summary Delete document
// @Description This API allows admin to archive document
// @Accept  json
// @Produce  json
// @Param document_id path string true
// @Success 200 {object} struct{message string}
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/documents/:document_id [delete]
func AdminDeleteDocument(c echo.Context) error {
	var cc, ok = c.(*models.CustomContext)
	if !ok {
		return eris.New("invalid custom context")
	}
	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, "get jwt claim info error")
	}
	var params models.DeleteDocumentParams
	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}
	params.JwtClaimsInfo = claims
	err = repo.NewDocumentRepo(cc.App.DB).DeleteDocument(&params)
	if err != nil {
		return err
	}

	return cc.Success("Archived")
}
