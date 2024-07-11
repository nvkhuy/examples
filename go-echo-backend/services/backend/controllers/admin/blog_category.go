package controllers

import (
	"github.com/engineeringinflow/inflow-backend/pkg/models"
	"github.com/engineeringinflow/inflow-backend/pkg/repo"
	"github.com/labstack/echo/v4"

	"github.com/rotisserie/eris"
)

// PaginateBlogCategory
// @Tags Admin-Blog
// @Summary PaginateBlogCategory
// @Description PaginateBlogCategory
// @Accept  json
// @Produce  json
// @Param keyword query string false "Keyword"
// @Param name query string false "Name"
// @Success 200 {object} models.BlogCategory
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/blog/categories [get]
func PaginateBlogCategory(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.PaginateBlogCategoryParams

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	params.JwtClaimsInfo = claims
	var result = repo.NewBlogCategoryRepo(cc.App.DB).PaginateBlogCategory(params)
	return cc.Success(result)
}

// CreateBlogCategory CreateFromPayload Post
// @Tags Admin-Blog
// @Summary CreateBlogCategory
// @Description CreateBlogCategory
// @Accept  json
// @Produce  json
// @Param data body models.BlogCategoryCreateForm true "Form"
// @Success 200 {object} models.BlogCategory
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/blog/categories [post]
func CreateBlogCategory(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var form models.BlogCategoryCreateForm

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	err = cc.BindAndValidate(&form)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	form.JwtClaimsInfo = claims
	result, err := repo.NewBlogCategoryRepo(cc.App.DB).CreateBlogCategory(form)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	return cc.Success(result)
}

// BlogCategoryUpdate Update category
// @Tags Admin-Blog
// @Summary Update category
// @Description Update category
// @Accept  json
// @Produce  json
// @Param blog_category_id path string true "ID"
// @Param data body models.BlogCategoryUpdateForm true "Form"
// @Success 200 {object} models.BlogCategory
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/blog/categories/{blog_category_id} [put]
func UpdateBlogCategory(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var form models.BlogCategoryUpdateForm

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	err = cc.BindAndValidate(&form)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	form.JwtClaimsInfo = claims
	Post, err := repo.NewBlogCategoryRepo(cc.App.DB).UpdateBlogCategory(form)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	return cc.Success(Post)
}

// DeleteBlogCategory
// @Tags Admin-Blog
// @Summary Delete category
// @Description Delete category
// @Accept  json
// @Produce  json
// @Param blog_category_id path string true "ID"
// @Success 200 {object} models.M
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/blog/categories/{blog_category_id}/delete [delete]
func DeleteBlogCategory(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.DeleteBlogCategoryParams

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	params.JwtClaimsInfo = claims
	err = repo.NewBlogCategoryRepo(cc.App.DB).DeleteBlogCategory(params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	return cc.Success("Archived")
}
