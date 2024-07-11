package controllers

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/engineeringinflow/inflow-backend/pkg/models"
	"github.com/engineeringinflow/inflow-backend/pkg/models/enums"
	"github.com/engineeringinflow/inflow-backend/pkg/repo"
	"github.com/labstack/echo/v4"
	"github.com/rotisserie/eris"
)

// PaginateBlogCategory
// @Tags Marketplace-Blog
// @Summary BlogCategoryList
// @Description BlogCategoryList
// @Accept  json
// @Produce  json
// @Success 200 {object} models.Records{records=[]models.BlogCategory}
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/blog/categories [get]
func PaginateBlogCategory(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.PaginateBlogCategoryParams

	var err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	params.TotalPost = aws.Int(0)
	var result = repo.NewBlogCategoryRepo(cc.App.DB).PaginateBlogCategory(params)

	return cc.Success(result)
}

// GetBlogCategory
// @Tags Marketplace-Blog
// @Summary BlogCategoryDetail
// @Description BlogCategoryDetail
// @Accept  json
// @Produce  json
// @Success 200 {object} models.BlogCategory
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/blog/categories/{blog_category_id} [get]
func GetBlogCategory(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.GetBlogCategoryParams

	var err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	cate, err := repo.NewBlogCategoryRepo(cc.App.DB).GetBlogCategory(params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	var result = repo.NewPostRepo(cc.App.DB).PaginatePost(repo.PaginatePostParams{
		JwtClaimsInfo: params.JwtClaimsInfo,
		CategoryIDs:   []string{params.BlogCategoryID},
		Language:      cc.GetLanguage(),
		Statuses: []enums.PostStatus{
			enums.PostStatusPublished,
		},
	})
	result.Metadata = cate

	return cc.Success(result)
}
