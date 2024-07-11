package controllers

import (
	"github.com/engineeringinflow/inflow-backend/pkg/models"
	"github.com/engineeringinflow/inflow-backend/pkg/repo"
	"github.com/labstack/echo/v4"

	"github.com/rotisserie/eris"
)

// PaginatePost
// @Tags Admin-Post
// @Summary Post List
// @Description Post List
// @Accept  json
// @Produce  json
// @Param keyword query string false "Keyword"
// @Param Post query int false "Post number"
// @Success 200 {object} models.Post
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/posts [get]
func PaginatePost(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.PaginatePostParams

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	params.JwtClaimsInfo = claims
	var result = repo.NewPostRepo(cc.App.DB).PaginatePost(params)
	return cc.Success(result)
}

// CreatePost CreateFromPayload Post
// @Tags Admin-Post
// @Summary CreateFromPayload Post
// @Description CreateFromPayload Post
// @Accept  json
// @Produce  json
// @Param data body models.PostCreateForm true "Form"
// @Success 200 {object} models.Post
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/posts/create [post]
func CreatePost(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var form models.PostCreateForm

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, "")
	}

	err = cc.BindAndValidate(&form)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	form.JwtClaimsInfo = claims
	form.UserID = claims.GetUserID()
	post, err := repo.NewPostRepo(cc.App.DB).CreatePost(form)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	return cc.Success(post)
}

// UpdatePost Update Post
// @Tags Admin-Post
// @Summary Update Post
// @Description Update Post
// @Accept  json
// @Produce  json
// @Param user_id path string true "ID"
// @Param data body models.PostUpdateForm true "Form"
// @Success 200 {object} models.Post
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/posts/{post_id} [put]
func UpdatePost(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var form models.PostUpdateForm

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	err = cc.BindAndValidate(&form)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	form.JwtClaimsInfo = claims
	post, err := repo.NewPostRepo(cc.App.DB).UpdatePost(form)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	return cc.Success(post)
}

// GetPost
// @Tags Admin-Post
// @Summary GetPost
// @Description GetPost
// @Accept  json
// @Produce  json
// @Param post_id query string true "PostID"
// @Success 200 {object} models.Post
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/posts/{slug} [get]
func GetPost(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.GetPostParams

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	params.JwtClaimsInfo = claims
	post, err := repo.NewPostRepo(cc.App.DB).GetPost(params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	return cc.Success(post)
}

// DeletePost
// @Tags Admin-Post
// @Summary Delete Post
// @Description Delete Post
// @Accept  json
// @Produce  json
// @Param id path string true "ID"
// @Success 200 {object} models.M
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/posts/{post_id}/delete [delete]
func DeletePost(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.DeletePostParams

	var err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	err = repo.NewPostRepo(cc.App.DB).DeletePost(params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	return cc.Success("Archived")
}

// GeneratePostSlug
// @Tags Admin-Post
// @Summary Delete Post
// @Description Delete Post
// @Accept  json
// @Produce  json
// @Param id path string true "ID"
// @Success 200 {object} models.M
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/posts/slug/generate [patch]
func GeneratePostSlug(c echo.Context) error {
	var cc = c.(*models.CustomContext)

	_, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	err = repo.NewPostRepo(cc.App.DB).GenerateSlug()
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	return cc.Success("Generated!")
}
