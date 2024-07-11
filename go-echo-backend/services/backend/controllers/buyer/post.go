package controllers

import (
	"github.com/engineeringinflow/inflow-backend/pkg/models"
	"github.com/engineeringinflow/inflow-backend/pkg/models/enums"
	"github.com/engineeringinflow/inflow-backend/pkg/repo"
	"github.com/labstack/echo/v4"
	"github.com/rotisserie/eris"
)

// PaginatePost
// @Tags Marketplace-Post
// @Summary Search Post
// @Description Search Post
// @Accept  json
// @Produce  json
// @Param page query int false "Page number"
// @Success 200 {object} models.Post
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/posts [get]
func PaginatePost(c echo.Context) error {
	var cc = c.(*models.CustomContext)

	var params repo.PaginatePostParams
	var err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	params.Statuses = []enums.PostStatus{
		enums.PostStatusPublished,
	}
	var result = repo.NewPostRepo(cc.App.DB).PaginatePost(params)

	return cc.Success(result)
}

// GetPost
// @Tags Marketplace-Post
// @Summary Search Post
// @Description Search Post
// @Accept  json
// @Produce  json
// @Param page query int false "Page number"
// @Success 200 {object} models.Post
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/posts/{slug} [get]
func GetPost(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.GetPostParams

	var err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	post, err := repo.NewPostRepo(cc.App.DB).GetPost(params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	return cc.Success(post)
}
