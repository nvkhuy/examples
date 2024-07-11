package controllers

import (
	"github.com/engineeringinflow/inflow-backend/pkg/models"
	"github.com/engineeringinflow/inflow-backend/pkg/repo"
	"github.com/labstack/echo/v4"

	"github.com/rotisserie/eris"
)

// AdminGetCategoryTree
// @Tags Admin-Category
// @Summary GetCategoryTree
// @Description GetCategoryTree
// @Accept  json
// @Produce  json
// @Success 200 {object} models.Category
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admins/categories/get_category_tree [get]
func AdminGetCategoryTree(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.PaginateCategoriesParams

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	err = cc.Bind(&params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	params.JwtClaimsInfo = claims
	var result = repo.NewCategoryRepo(cc.App.DB).GetCategoryTree(params)

	return cc.Success(result)
}

// CreateCategory create category
// @Tags Admin-Category
// @Summary CreateFromPayload category
// @Description CreateFromPayload category
// @Accept  json
// @Produce  json
// @Param data body models.CategoryUpdateForm true "Form"
// @Success 200 {object} models.Category
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/categories [post]
func CreateCategory(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var form models.CategoryCreateForm

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	err = cc.BindAndValidate(&form)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	form.JwtClaimsInfo = claims
	cate, err := repo.NewCategoryRepo(cc.App.DB).CreateCategory(form)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	return cc.Success(cate)
}

// UpdateCategory Update category
// @Tags Admin-Category
// @Summary Update category
// @Description Update category
// @Accept  json
// @Produce  json
// @Param user_id path string true "ID"
// @Param data body models.CategoryUpdateForm true "Form"
// @Success 200 {object} models.Category
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/categories/{category_id} [put]
func UpdateCategory(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var form models.CategoryUpdateForm

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	err = cc.BindAndValidate(&form)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	form.JwtClaimsInfo = claims
	cate, err := repo.NewCategoryRepo(cc.App.DB).UpdateCategoryByID(form)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	return cc.Success(cate)
}

// DeleteCategory Delete category
// @Tags Admin-Category
// @Summary Delete category
// @Description Delete category
// @Accept  json
// @Produce  json
// @Param category_id path string true "ID"
// @Success 200 {object} models.M
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/categories/{category_id}/delete [delete]
func DeleteCategory(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.DeleteCategoryParams

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	params.JwtClaimsInfo = claims
	err = repo.NewCategoryRepo(cc.App.DB).DeleteCategory(params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	return cc.Success("Deleted")
}

// GenerateCategorySlug
// @Tags Admin-Category
// @Summary Delete category
// @Description Delete category
// @Accept  json
// @Produce  json
// @Param category_id path string true "ID"
// @Success 200 {object} models.M
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/categories/slug/generate [patch]
func GenerateCategorySlug(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.GenerateCategorySlugParams

	_, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	err = repo.NewCategoryRepo(cc.App.DB).GenerateSlug()
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	return cc.Success("Generated")
}
