package controllers

import (
	"github.com/engineeringinflow/inflow-backend/pkg/models"
	"github.com/engineeringinflow/inflow-backend/pkg/repo"
	"github.com/engineeringinflow/inflow-backend/pkg/repo/query/queryfunc"
	"github.com/labstack/echo/v4"
	"github.com/rotisserie/eris"
)

// SearchProduct
// @Tags Marketplace-Product
// @Summary Search Product
// @Description Search Product
// @Accept  json
// @Produce  json
// @Param keyword query string false "Keyword"
// @Param category_id query string false "Category ID"
// @Param rating_star query int false "Rating start"
// @Param min_order query int false "Min order"
// @Param product_type query string false "Product type"
// @Param page query int false "Page number"
// @Success 200 {object} models.ProductResponse
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/products/search [get]
func SearchProduct(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.PaginateProductParams

	var err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	var result = repo.NewProductRepo(cc.App.DB).PaginateProducts(params)

	return cc.Success(result)
}

// ProductGetDetail
// @Tags Marketplace-Product
// @Summary Product Detail
// @Description Search Product
// @Accept  json
// @Produce  json
// @Param product_id query string true "ProductID"
// @Success 200 {object} models.Product
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/products/get [get]
func ProductGetDetail(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.GetProductParams

	var err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	product, err := repo.NewProductRepo(cc.App.DB).GetProduct(params)
	if err != nil {
		return eris.Wrap(err, "")
	}

	var response models.ProductDetailWithVariant
	response.Product = product

	variants, err := repo.NewVariantRepo(cc.App.DB).GetVariantsByProductID(product.ID, queryfunc.VariantBuilderOptions{})
	if err != nil {
		return eris.Wrap(err, "")
	}
	response.Variants = variants

	options, err := repo.NewProductAttributeRepo(cc.App.DB).GetProductAttributesByProductID(product.ID, queryfunc.ProductAttributeBuilderOptions{})
	if err != nil {
		return eris.Wrap(err, "")
	}
	response.Options = options

	priceTiers, err := repo.NewQuantityPriceTierRepo(cc.App.DB).GetQuantityPriceTierByProductID(product.ID, queryfunc.QuantityPriceTierBuilderOptions{})
	if err != nil {
		return eris.Wrap(err, "")
	}
	response.QuantityPriceTiers = priceTiers

	return cc.Success(response)
}

// ProductGetRating
// @Tags Marketplace-Product
// @Summary Product GetRating
// @Description Product GetRating
// @Accept  json
// @Produce  json
// @Param product_id query string true "ProductID"
// @Param page query int false "Page number"
// @Success 200 {object} models.Product
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/products/get_ratings [get]
func ProductGetRatings(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.PaginateProductReviewParams

	var err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	var result = repo.NewProductReviewRepo(cc.App.DB).PaginateProductReviews(params)
	return cc.Success(result)
}

// ProductGetBestSelling
// @Tags Marketplace-Product
// @Summary Product Get BestSelling
// @Description Product Get BestSelling
// @Accept  json
// @Produce  json
// @Param page query int false "Page number"
// @Success 200 {object} models.Product
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/products/best_selling [get]
func ProductGetBestSelling(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.PaginateProductParams

	var err = cc.Bind(&params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	var result = repo.NewProductRepo(cc.App.DB).PaginateBestSellingProducts(params)
	return cc.Success(result)
}

// ProductGetJustForYou
// @Tags Marketplace-Product
// @Summary Product GetJustForYou
// @Description Product Get GetJustForYou
// @Accept  json
// @Produce  json
// @Param page query int false "Page number"
// @Success 200 {object} models.Product
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/products/just_for_you [get]
func ProductGetJustForYou(c echo.Context) error {
	var cc = c.(*models.CustomContext)

	var params repo.PaginateProductParams

	var err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	var result = repo.NewProductRepo(cc.App.DB).PaginateJustForYouProducts(params)
	return cc.Success(result)
}

// ProductReadyToShip
// @Tags Marketplace-Product
// @Summary Product ReadyToShip
// @Description Product ReadyToShip
// @Accept  json
// @Produce  json
// @Param page query int false "Page number"
// @Success 200 {object} models.Product
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/products/ready_to_ship [get]
func ProductReadyToShip(c echo.Context) error {
	var cc = c.(*models.CustomContext)

	var params repo.PaginateProductParams
	var err = cc.Bind(&params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	params.ReadyToShip = true
	var result = repo.NewProductRepo(cc.App.DB).PaginateProducts(params)
	return cc.Success(result)
}

// ProductTodayDeals
// @Tags Marketplace-Product
// @Summary Product TodayDeals
// @Description Product TodayDeals
// @Accept  json
// @Produce  json
// @Param page query int false "Page number"
// @Success 200 {object} models.Product
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/products/today_deals [get]
func ProductTodayDeals(c echo.Context) error {
	var cc = c.(*models.CustomContext)

	var params repo.PaginateProductParams
	var err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	params.DailyDeal = true
	var result = repo.NewProductRepo(cc.App.DB).PaginateProducts(params)
	return cc.Success(result)
}

// ProductRecommend
// @Tags Marketplace-Product
// @Summary Product GetJustForYou
// @Description Product Get GetJustForYou
// @Accept  json
// @Produce  json
// @Param page query int false "Page number"
// @Success 200 {object} models.Product
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/products/recommendations [get]
func ProductRecommend(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.PaginateRecommendProductParams

	var err = cc.Bind(&params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}
	if params.Limit == 0 {
		params.Limit = 10
	}
	var result = repo.NewProductRepo(cc.App.DB).WithAnalyticDB(cc.App.AnalyticDB).PaginateRecommendations(params)
	return cc.Success(result)
}

// ProductGetFilter
// @Tags Marketplace-Product
// @Summary Product GetFilter
// @Description Product GetFilter
// @Accept  json
// @Produce  json
// @Success 200 {object} models.ProductFilter
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/products/get_filter [get]
func ProductGetFilter(c echo.Context) error {
	var cc = c.(*models.CustomContext)

	var constants = []models.ProductFilter{
		{
			Name: "Rating",
			Key:  "min_rating",
			Type: "single_select",
			Options: []models.ProductFilterOption{
				{
					Name:  "From 1 star",
					Value: 1,
				},
				{
					Name:  "From 2 star",
					Value: 2,
				},
				{
					Name:  "From 3 star",
					Value: 3,
				},
				{
					Name:  "From 4 star",
					Value: 4,
				},
				{
					Name:  "From 5 star",
					Value: 5,
				},
			},
		},
		{
			Name: "Minimum order",
			Key:  "min_order",
			Type: "single_select",
			Options: []models.ProductFilterOption{
				{
					Name:  "From 10 items",
					Value: 10,
				},
				{
					Name:  "From 100 items",
					Value: 100,
				},
				{
					Name:  "From 500 items",
					Value: 500,
				},
			},
		},
	}

	return cc.Success(constants)
}
