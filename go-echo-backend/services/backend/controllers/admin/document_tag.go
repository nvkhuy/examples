package controllers

import (
	"github.com/engineeringinflow/inflow-backend/pkg/models"
	"github.com/engineeringinflow/inflow-backend/pkg/repo"
	"github.com/labstack/echo/v4"
	"github.com/rotisserie/eris"
)

// @Tags Admin-DocumentTag
// @Summary Create document tag
// @Description This API allows admin to create document tag
// @Accept  json
// @Produce  json
// @Param data body models.CreateDocumentTagRequest true
// @Success 200 {object} models.DocumentTag
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/document_tags [post]
func AdminCreateDocumentTag(c echo.Context) error {
	var cc, ok = c.(*models.CustomContext)
	if !ok {
		return eris.New("invalid custom context")
	}
	var req models.CreateDocumentTagRequest

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, "get jwt claim info error")
	}

	err = cc.BindAndValidate(&req)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}
	req.JwtClaimsInfo = claims
	resp, err := repo.NewDocumentTagRepo(cc.App.DB).CreateDocumentTag(&req)
	if err != nil {
		return err
	}

	return cc.Success(resp)
}

// @Tags Admin-DocumentTag
// @Summary Update document tag
// @Description This API allows admin to update document Tag
// @Accept  json
// @Produce  json
// @Param data body models.UpdateDocumentTagRequest true
// @Param document_tag_id path string true
// @Success 200 {object} models.DocumentTag
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/document_tags/:document_tag_id [put]
func AdminUpdateDocumentTag(c echo.Context) error {
	var cc, ok = c.(*models.CustomContext)
	if !ok {
		return eris.New("invalid custom context")
	}
	var req models.UpdateDocumentTagRequest

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, "get jwt claim info error")
	}

	err = cc.BindAndValidate(&req)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}
	req.JwtClaimsInfo = claims
	resp, err := repo.NewDocumentTagRepo(cc.App.DB).UpdateDocumentTag(&req)
	if err != nil {
		return err
	}

	return cc.Success(resp)
}

// @Tags Admin-DocumentTag
// @Summary Get document tag list
// @Description This API allows admin to list document tag with pagination
// @Accept  json
// @Produce  json
// @Param page query int false
// @Param limit query int false
// @Param keyword query string false
// @Success 200 {object} query.Pagination{Records=[]models.DocumentTag}
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/document_tags [get]
func AdminGetDocumentTagList(c echo.Context) error {
	var cc, ok = c.(*models.CustomContext)
	if !ok {
		return eris.New("invalid custom context")
	}
	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, "get jwt claim info error")
	}
	var params models.GetDocumentTagListParams
	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}
	params.JwtClaimsInfo = claims
	resp := repo.NewDocumentTagRepo(cc.App.DB).GetDocumentTagList(&params)

	return cc.Success(resp)
}

// @Tags Admin-DocumentTag
// @Summary Delete document tag
// @Description This API allows admin to archive document tag
// @Accept  json
// @Produce  json
// @Param document_tag_id path string true
// @Success 200 {object} struct{message string}
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/document_tags/:document_tag_id [delete]
func AdminDeleteDocumentTag(c echo.Context) error {
	var cc, ok = c.(*models.CustomContext)
	if !ok {
		return eris.New("invalid custom context")
	}
	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, "get jwt claim info error")
	}
	var params models.DeleteDocumentTagRequest
	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}
	params.JwtClaimsInfo = claims
	err = repo.NewDocumentTagRepo(cc.App.DB).DeleteDocumentTag(&params)
	if err != nil {
		return err
	}

	return cc.Success("Archived")
}
