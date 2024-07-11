package controllers

import (
	"fmt"
	"github.com/engineeringinflow/inflow-backend/pkg/models"
	"github.com/engineeringinflow/inflow-backend/pkg/models/enums"
	"github.com/engineeringinflow/inflow-backend/pkg/repo"
	"github.com/engineeringinflow/inflow-backend/pkg/repo/query/queryfunc"
	"github.com/engineeringinflow/inflow-backend/services/consumer/tasks"
	"github.com/labstack/echo/v4"

	"github.com/rotisserie/eris"
)

// PaginateProduct
// @Tags Admin-Product
// @Summary PaginateProduct
// @Description PaginateProduct
// @Accept  json
// @Produce  json
// @Param keyword query string false "Keyword"
// @Param category_id query string false "Category ID"
// @Param shop_ids query array false "Shop IDs"
// @Param page query int false "Page number"
// @Success 200 {object} models.ProductResponse
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/products [get]
func PaginateProduct(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.PaginateProductParams

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	err = cc.Bind(&params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	params.JwtClaimsInfo = claims
	var result = repo.NewProductRepo(cc.App.DB).PaginateProducts(params)

	return cc.Success(result)
}

// GenerateProductSlug
// @Tags Admin-Product
// @Summary PaginateProduct
// @Description PaginateProduct
// @Accept  json
// @Produce  json
// @Param keyword query string false "Keyword"
// @Param category_id query string false "Category ID"
// @Param shop_ids query array false "Shop IDs"
// @Param page query int false "Page number"
// @Success 200 {object} models.ProductResponse
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/products/fabric/price [patch]
func GenerateProductSlug(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params models.GenerateProductSlugParams

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}
	params.JwtClaimsInfo = claims
	err = repo.NewProductRepo(cc.App.DB).GenerateSlug()
	if err != nil {
		return eris.Wrap(err, err.Error())
	}
	return cc.Success("generated")
}

// FetchProductTypesPrice
// @Tags Admin-Product
// @Summary PaginateProduct
// @Description PaginateProduct
// @Accept  json
// @Produce  json
// @Param keyword query string false "Keyword"
// @Param category_id query string false "Category ID"
// @Param shop_ids query array false "Shop IDs"
// @Param page query int false "Page number"
// @Success 200 {object} models.ProductResponse
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/products/types/price [patch]
func FetchProductTypesPrice(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params models.FetchProductTypesPriceParams

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	err = cc.Bind(&params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}
	params.JwtClaimsInfo = claims
	result, err := repo.NewProductTypesPriceRepo(cc.App.DB).
		WithSheetAPI(cc.App.SheetAPI).
		FetchProductTypesPrice(&params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}
	return cc.Success(result)
}

// FetchFabricTypesPrice
// @Tags Admin-Product
// @Summary PaginateProduct
// @Description PaginateProduct
// @Accept  json
// @Produce  json
// @Param keyword query string false "Keyword"
// @Param category_id query string false "Category ID"
// @Param shop_ids query array false "Shop IDs"
// @Param page query int false "Page number"
// @Success 200 {object} models.ProductResponse
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/products/fabric/price [patch]
func FetchFabricTypesPrice(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params models.FetchRWDFabricPriceParams

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}
	params.JwtClaimsInfo = claims
	result, err := repo.NewRWDFabricPriceRepo(cc.App.DB).
		WithSheetAPI(cc.App.SheetAPI).
		FetchPrice(&params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}
	return cc.Success(result)
}

// FetchProductTypesPriceImagesURL
// @Tags Admin-Product
// @Summary PaginateProduct
// @Description PaginateProduct
// @Accept  json
// @Produce  json
// @Param keyword query string false "Keyword"
// @Param category_id query string false "Category ID"
// @Param shop_ids query array false "Shop IDs"
// @Param page query int false "Page number"
// @Success 200 {object} models.ProductResponse
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/products/types/price/images [patch]
func FetchProductTypesPriceImagesURL(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params models.PatchSheetImageURLParams

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}
	params.JwtClaimsInfo = claims
	go func(ctx *models.CustomContext) {
		_ = repo.NewProductTypesPriceRepo(ctx.App.DB).
			WithSheetAPI(ctx.App.SheetAPI).PatchSheetImageURL(&params)
	}(cc)

	return cc.Success("request send")
}

// PaginateProductTypesPrice
// @Tags Admin-Product
// @Summary PaginateProduct
// @Description PaginateProduct
// @Accept  json
// @Produce  json
// @Param keyword query string false "Keyword"
// @Param category_id query string false "Category ID"
// @Param shop_ids query array false "Shop IDs"
// @Param page query int false "Page number"
// @Success 200 {object} models.ProductResponse
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/products/types/price [get]
func PaginateProductTypesPrice(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params models.PaginateProductTypesPriceParams

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	err = cc.Bind(&params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	params.JwtClaimsInfo = claims
	result := repo.NewProductTypesPriceRepo(cc.App.DB).PaginateProductTypesPrice(&params)
	return cc.Success(result)
}

// PaginateProductTypesPriceVine
// @Tags Admin-Product
// @Summary PaginateProduct
// @Description PaginateProduct
// @Accept  json
// @Produce  json
// @Param keyword query string false "Keyword"
// @Param category_id query string false "Category ID"
// @Param shop_ids query array false "Shop IDs"
// @Param page query int false "Page number"
// @Success 200 {object} models.ProductResponse
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/products/types/price/vine [get]
func PaginateProductTypesPriceVine(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params models.PaginateProductTypesPriceParams

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	err = cc.Bind(&params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	params.JwtClaimsInfo = claims
	result := repo.NewProductTypesPriceRepo(cc.App.DB).PaginateProductTypesPriceVine(&params)
	return cc.Success(result)
}

// PaginateRWDFabricPriceVine
// @Tags Admin-Product
// @Summary PaginateProduct
// @Description PaginateProduct
// @Accept  json
// @Produce  json
// @Param keyword query string false "Keyword"
// @Param category_id query string false "Category ID"
// @Param shop_ids query array false "Shop IDs"
// @Param page query int false "Page number"
// @Success 200 {object} models.ProductResponse
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/products/fabric/price/vine [get]
func PaginateRWDFabricPriceVine(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params models.PaginateRWDFabricPriceParams

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	params.JwtClaimsInfo = claims
	result := repo.NewRWDFabricPriceRepo(cc.App.DB).Vine(&params)
	return cc.Success(result)
}

// PaginateProductTypesPriceQuote
// @Tags Admin-Product
// @Summary PaginateProduct
// @Description PaginateProduct
// @Accept  json
// @Produce  json
// @Param keyword query string false "Keyword"
// @Param category_id query string false "Category ID"
// @Param shop_ids query array false "Shop IDs"
// @Param page query int false "Page number"
// @Success 200 {object} models.ProductResponse
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/products/types/price/quote [get]
func PaginateProductTypesPriceQuote(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params models.PaginateProductTypesPriceQuoteParams

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	params.JwtClaimsInfo = claims
	result := repo.NewProductTypesPriceRepo(cc.App.DB).PaginateProductTypesPriceQuote(&params)
	return cc.Success(result)
}

// CreateProduct create product
// @Tags Admin-Product
// @Summary create product
// @Description create product
// @Accept  json
// @Produce  json
// @Param data body models.ProductCreateForm true "Form"
// @Success 200 {object} models.Product
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/products/create [post]
func CreateProduct(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var form models.ProductCreateForm

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	err = cc.BindAndValidate(&form)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	form.JwtClaimsInfo = claims
	product, err := repo.NewProductRepo(cc.App.DB).CreateProduct(form)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	if product != nil {
		_, _ = tasks.CreateSysNotificationTask{
			SysNotification: models.SysNotification{
				Name:    fmt.Sprintf("New Product - %s", product.Name),
				Type:    enums.SysNotificationCreateProductType,
				Message: fmt.Sprintf("New Product - %s", product.Name),
			},
		}.Dispatch(c.Request().Context())
	}

	return cc.Success(product)
}

// GetProductQRCode create product
// @Tags Admin-Product
// @Summary create product
// @Description create product
// @Accept  json
// @Produce  json
// @Param data body models.ProductCreateForm true "Form"
// @Success 200 {object} models.Product
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/products/qr_code [get]
func GetProductQRCode(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var form models.GetProductQRCodeParams

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	err = cc.BindAndValidate(&form)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	form.JwtClaimsInfo = claims
	form.Logo = cc.App.Config.QRCodeLogoURL
	form.Bucket = cc.App.Config.AWSS3StorageBucket
	var url string
	url, err = repo.NewProductRepo(cc.App.DB).
		GetQRCode(form)

	if err != nil {
		return eris.Wrap(err, err.Error())
	}
	return cc.Success(url)
}

// ExportProducts
// @Tags Admin-Product
// @Summary create product
// @Description create product
// @Accept  json
// @Produce  json
// @Param data body models.ProductCreateForm true "Form"
// @Success 200 {object} models.Product
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/products/export [get]
func ExportProducts(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.PaginateProductParams

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	params.JwtClaimsInfo = claims
	resp, err := repo.NewProductRepo(cc.App.DB).ExportExcel(params)

	if err != nil {
		return eris.Wrap(err, err.Error())
	}
	return cc.Success(resp)
}

// UpdateProduct Update Product
// @Tags Admin-Product
// @Summary Update Product
// @Description Update Product
// @Accept  json
// @Produce  json
// @Param user_id path string true "ID"
// @Param data body models.ProductCreateForm true "Form"
// @Success 200 {object} models.Product
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/products/{product_id} [put]
func UpdateProduct(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var form models.ProductUpdateForm

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	err = cc.BindAndValidate(&form)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	form.JwtClaimsInfo = claims
	product, err := repo.NewProductRepo(cc.App.DB).UpdateProduct(form)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	return cc.Success(product)
}

// GetProduct
// @Tags Admin-Product
// @Summary Product Detail
// @Description Product Detail
// @Accept  json
// @Produce  json
// @Param product_id query string true "ProductID"
// @Success 200 {object} models.Product
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/products/get [get]
func GetProduct(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.GetProductParams

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, "")
	}

	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, "")
	}

	params.JwtClaimsInfo = claims
	product, err := repo.NewProductRepo(cc.App.DB).GetProduct(params)
	if err != nil {
		return eris.Wrap(err, "")
	}

	var response models.ProductDetailWithVariant
	response.Product = product

	variants, err := repo.NewVariantRepo(cc.App.DB).GetVariantsByProductID(params.ProductID, queryfunc.VariantBuilderOptions{})
	if err != nil {
		return eris.Wrap(err, "")
	}
	response.Variants = variants

	options, err := repo.NewProductAttributeRepo(cc.App.DB).GetProductAttributesByProductID(params.ProductID, queryfunc.ProductAttributeBuilderOptions{})
	if err != nil {
		return eris.Wrap(err, "")
	}
	response.Options = options

	priceTiers, err := repo.NewQuantityPriceTierRepo(cc.App.DB).GetQuantityPriceTierByProductID(params.ProductID, queryfunc.QuantityPriceTierBuilderOptions{})
	if err != nil {
		return eris.Wrap(err, "")
	}
	response.QuantityPriceTiers = priceTiers

	return cc.Success(response)
}

// AdminDeleteProduct
// @Tags Admin-Product
// @Summary Delete Product
// @Description Delete Product
// @Accept  json
// @Produce  json
// @Param Product_id path string true "ID"
// @Success 200 {object} models.M
// @Header 200 {string} Bearer YOUR_TOKEN
// @Security ApiKeyAuth
// @Failure 404 {object} errs.Error
// @Router /api/v1/admin/products/{product_id}/delete [delete]
func AdminDeleteProduct(c echo.Context) error {
	var cc = c.(*models.CustomContext)
	var params repo.DeleteProductParams

	claims, err := cc.GetJwtClaimsInfo()
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	err = cc.BindAndValidate(&params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	params.JwtClaimsInfo = claims
	err = repo.NewProductRepo(cc.App.DB).DeleteProduct(params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	return cc.Success("Deleted")
}

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
// @Router /api/v1/admin/products/search [get]
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
