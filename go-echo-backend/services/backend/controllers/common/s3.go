package controllers

import (
	"github.com/engineeringinflow/inflow-backend/pkg/models"
	"github.com/engineeringinflow/inflow-backend/pkg/repo"
	"github.com/labstack/echo/v4"
	"github.com/rotisserie/eris"
)

// GetS3Signatures Generate s3 signatures
// @Tags Common
// @Summary Generate s3 signatures
// @Description Generate s3 signatures
// @Accept  json
// @Produce  json
// @Param data body models.S3SignatureForms true "Form"
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/common/s3_signatures [post]
func GetS3Signatures(c echo.Context) error {
	var cc = c.(*models.CustomContext)

	var form models.S3SignatureForms
	var err = cc.BindAndValidate(&form)
	if err != nil {
		return eris.Wrap(err, "")
	}

	claims, _ := cc.GetJwtClaimsInfo()

	var result = repo.NewCommonRepo(cc.App.DB).GetS3Signatures(repo.GetS3SignaturesParams{
		JwtClaimsInfo: claims,
		Forms:         form.Records,
	})

	return cc.Success(result)
}

// GetS3Signature Generate s3 signatures
// @Tags Common
// @Summary Generate s3 signatures
// @Description Generate s3 signatures
// @Accept  json
// @Produce  json
// @Param data body models.S3SignatureForm true "Form"
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/common/s3_signature [post]
func GetS3Signature(c echo.Context) error {
	var cc = c.(*models.CustomContext)

	var form models.S3SignatureForm
	var err = cc.BindAndValidate(&form)
	if err != nil {
		return eris.Wrap(err, "")
	}

	claims, _ := cc.GetJwtClaimsInfo()

	var result = repo.NewCommonRepo(cc.App.DB).GetS3Signature(repo.GetS3SignatureParams{
		JwtClaimsInfo: claims,
		Form:          &form,
	})

	return cc.Success(result)
}
