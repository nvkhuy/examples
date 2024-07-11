package repo

import (
	"bytes"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/engineeringinflow/inflow-backend/pkg/helper"
	"github.com/engineeringinflow/inflow-backend/pkg/models/enums"
	"github.com/engineeringinflow/inflow-backend/pkg/s3"

	"github.com/gosimple/slug"
	"google.golang.org/api/sheets/v4"

	"github.com/engineeringinflow/inflow-backend/pkg/db"
	"github.com/engineeringinflow/inflow-backend/pkg/errs"
	"github.com/engineeringinflow/inflow-backend/pkg/logger"
	"github.com/engineeringinflow/inflow-backend/pkg/models"
	"github.com/engineeringinflow/inflow-backend/pkg/repo/query"
	"github.com/engineeringinflow/inflow-backend/pkg/repo/query/queryfunc"
	"github.com/jinzhu/copier"

	"github.com/rotisserie/eris"
	"github.com/yeqown/go-qrcode/writer/standard"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ProductRepo struct {
	db       *db.DB
	adb      *db.DB
	logger   *logger.Logger
	sheetAPI *sheets.Service
}

func NewProductRepo(db *db.DB) *ProductRepo {
	return &ProductRepo{
		db:     db,
		logger: logger.New("repo/Product"),
	}
}

func (r *ProductRepo) WithSheetAPI(api *sheets.Service) *ProductRepo {
	r.sheetAPI = api
	return r
}

func (r *ProductRepo) WithAnalyticDB(db *db.DB) *ProductRepo {
	r.adb = db
	return r
}

type PaginateProductParams struct {
	models.PaginationParams
	models.JwtClaimsInfo

	ParentCategoryID   string   `json:"parent_category_id" query:"parent_category_id" form:"parent_category_id"`
	Name               string   `json:"name" query:"name" form:"name"`
	CategoryID         string   `json:"category_id" query:"category_id" form:"category_id"`
	CategorySlug       string   `json:"category_slug" query:"category_slug"`
	SubCategorySlug    string   `json:"sub_category_slug" query:"sub_category_slug"`
	ShopIDs            []string `json:"shop_ids" query:"shop_ids" form:"shop_ids"`
	ReadyToShip        bool     `json:"ready_to_ship" query:"ready_to_ship" form:"ready_to_ship"`
	DailyDeal          bool     `json:"daily_deal" query:"daily_deal" form:"daily_deal"`
	RatingStar         float32  `json:"rating_star" query:"rating_star" form:"rating_star"`
	MinOrder           int      `json:"min_order" query:"min_order" form:"min_order"`
	ProductType        string   `json:"product_type" query:"product_type" form:"product_type"`
	ExceptedProductIDs []string `json:"excepted_product_ids" query:"excepted_product_ids" form:"excepted_product_ids"`
	Tags               []string `json:"tags" query:"tags" param:"tags"`
	RecommendProductID string   `json:"recommend_product_id" query:"recommend_product_id" param:"recommend_product_id"`
	ProductClass       string   `json:"product_class" param:"product_class" query:"product_class" form:"product_class"`
}

func (r *ProductRepo) PaginateProducts(params PaginateProductParams) *query.Pagination {
	var builder = queryfunc.NewProductBuilder(queryfunc.ProductBuilderOptions{
		QueryBuilderOptions: queryfunc.QueryBuilderOptions{
			Role: params.GetRole(),
		},
	})

	var result = query.New(r.db, builder).
		WhereFunc(func(builder *query.Builder) {
			if params.ParentCategoryID != "" {
				builder.Where("p.category_id = ?", params.ParentCategoryID)
			}

			if params.SubCategorySlug != "" {
				params.SubCategorySlug = strings.ToLower(strings.TrimSpace(params.SubCategorySlug))
				builder.Where("ct.slug = ? OR ct.vi ->> 'slug' = ?", params.SubCategorySlug, params.SubCategorySlug)
			} else if params.CategorySlug != "" {
				params.CategorySlug = strings.ToLower(strings.TrimSpace(params.CategorySlug))
				builder.Where("pr.slug = ? OR pr.vi ->> 'slug' = ?", params.CategorySlug, params.CategorySlug)
			}

			if params.CategoryID != "" {
				var cateIds = NewCategoryRepo(r.db).GetChildCategoryIDs(params.CategoryID)
				cateIds = append(cateIds, params.CategoryID)
				builder.Where("p.category_id IN (?)", cateIds)
			}

			if params.ExceptedProductIDs != nil && len(params.ExceptedProductIDs) > 0 {
				builder.Where("p.id NOT IN (?)", params.ExceptedProductIDs)
			}

			if len(params.ShopIDs) > 0 {
				builder.Where("p.shop_id IN ?", params.ShopIDs)
			}

			if params.ReadyToShip {
				builder.Where("p.ready_to_ship = ?", true)
			}

			if params.DailyDeal {
				builder.Where("p.daily_deal = ?", true)
			}

			if params.RatingStar > 0 {
				builder.Where("p.rating_star >= ?", params.RatingStar)
			}

			if params.MinOrder > 0 {
				builder.Where("p.min_order >= ?", params.MinOrder)
			}

			if params.ProductType != "" {
				builder.Where("p.product_type = ?", params.ProductType)
			}

			if len(params.Tags) > 0 {
				for _, tag := range params.Tags {
					tag = strings.TrimSpace(strings.ToLower(tag))
					if enums.ProductTag(tag) == enums.ProductTagTrending {
						builder.Where("p.is_trending = ?", true)
					}
				}
			}

			if strings.TrimSpace(params.Keyword) != "" {
				var q = "%" + params.Keyword + "%"
				builder.Where("p.name ILIKE @query", sql.Named("query", q))
			}
		}).
		Page(params.Page).
		Limit(params.Limit).
		OrderBy("p.created_at DESC").
		PagingFunc()

	return result
}

type PaginateRecommendProductParams struct {
	models.PaginationParams
	models.JwtClaimsInfo

	RecommendProductID         string `json:"recommend_product_id" query:"recommend_product_id" param:"recommend_product_id"`
	RecommendPlatformProductID string `json:"recommend_platform_product_id" query:"recommend_platform_product_id" param:"recommend_platform_product_id"`
	ProductClass               string `json:"product_class" param:"product_class" query:"product_class" form:"product_class"`
}

func (r *ProductRepo) PaginateRecommendations(params PaginateRecommendProductParams) *query.Pagination {
	var builder = queryfunc.NewProductRecommendBuilder(queryfunc.ProductRecommendBuilderOptions{
		QueryBuilderOptions: queryfunc.QueryBuilderOptions{
			Role: params.GetRole(),
		},
	})

	var result = query.New(r.db, builder).
		WhereFunc(func(builder *query.Builder) {
			if params.RecommendProductID != "" {
				var class models.ProductClass
				if err := r.db.Model(&models.ProductClass{}).
					Where("product_id = ?", params.RecommendProductID).Order("conf DESC").First(&class).Error; err != nil {
					return
				}
				if class.Class != "" {
					builder.Where("pc.class = ?", class.Class)
				}
			} else if params.RecommendPlatformProductID != "" {
				var class models.AnalyticProductClass
				if err := r.adb.Model(&models.AnalyticProductClass{}).
					Where("product_id = ?", params.RecommendPlatformProductID).Order("conf DESC").First(&class).Error; err != nil {
					return
				}
				if class.Class != "" {
					builder.Where("pc.class = ?", class.Class)
				}
			} else if params.ProductClass != "" {
				builder.Where("pc.class = ?", params.ProductClass)
			}
		}).
		Page(params.Page).
		Limit(params.Limit).
		OrderBy("pc.conf DESC").
		PagingFunc()

	return result
}

func (r *ProductRepo) PaginateJustForYouProducts(params PaginateProductParams) *query.Pagination {
	var builder = queryfunc.NewProductBuilder(queryfunc.ProductBuilderOptions{
		QueryBuilderOptions: queryfunc.QueryBuilderOptions{
			Role: params.GetRole(),
		},
	})

	if params.Limit == 0 {
		params.Limit = 20
	}

	var result = query.New(r.db, builder).
		WhereFunc(func(builder *query.Builder) {
			if params.CategoryID != "" {
				builder.Where("p.category_id IN (?)", params.CategoryID)
			}

			if params.ExceptedProductIDs != nil && len(params.ExceptedProductIDs) > 0 {
				builder.Where("p.id NOT IN (?)", params.ExceptedProductIDs)
			}

			if strings.TrimSpace(params.Keyword) != "" {
				var q = "%" + params.Keyword + "%"
				builder.Where("name ILIKE @query", sql.Named("query", q))
			}
		}).
		Page(params.Page).
		Limit(params.Limit).
		PagingFunc()

	return result
}

func (r *ProductRepo) PaginateBestSellingProducts(params PaginateProductParams) *query.Pagination {
	var builder = queryfunc.NewProductBuilder(queryfunc.ProductBuilderOptions{
		QueryBuilderOptions: queryfunc.QueryBuilderOptions{
			Role: params.GetRole(),
		},
	})

	if params.Limit == 0 {
		params.Limit = 20
	}

	var result = query.New(r.db, builder).
		WhereFunc(func(builder *query.Builder) {
			if params.CategoryID != "" {
				var childCategories = NewCategoryRepo(r.db).GetChildCategoryIDs(params.CategoryID)
				builder.Where("p.category_id IN (?)", childCategories)
			}

			if params.ExceptedProductIDs != nil && len(params.ExceptedProductIDs) > 0 {
				builder.Where("p.id NOT IN (?)", params.ExceptedProductIDs)
			}

			if strings.TrimSpace(params.Keyword) != "" {
				var q = "%" + params.Keyword + "%"
				builder.Where("name ILIKE @query", sql.Named("query", q))
			}
		}).
		Page(params.Page).
		Limit(params.Limit).
		PagingFunc()

	return result
}

func (r *ProductRepo) CreateProduct(form models.ProductCreateForm) (*models.Product, error) {
	var product models.Product
	err := copier.Copy(&product, &form)
	if err != nil {
		return nil, eris.Wrap(err, "")
	}

	err = r.db.Clauses(clause.Returning{}).Omit(clause.Associations).Create(&product).Error
	if err != nil {
		if r.db.IsRecordNotFoundError(err) {
			return nil, errs.ErrProductExisted
		}
		return nil, eris.Wrap(err, "")
	}

	// CreateFromPayload price tiers
	if len(form.QuantityPriceTiers) > 0 {
		var priceTiers []*models.QuantityPriceTier
		err = copier.Copy(&priceTiers, &form.QuantityPriceTiers)
		if err != nil {
			return nil, eris.Wrap(err, "")
		}

		for _, attr := range priceTiers {
			attr.ProductID = product.ID
			attr.Unit = product.TradeUnit
		}

		err = r.db.Omit(clause.Associations).Create(&priceTiers).Error
		if err != nil {
			return nil, eris.Wrap(err, "")
		}
	}

	// CreateFromPayload variants
	if len(form.Variants) > 0 {
		var variants []*models.Variant
		err = copier.Copy(&variants, &form.Variants)
		if err != nil {
			return nil, eris.Wrap(err, "")
		}
		for _, variant := range variants {
			variant.ProductID = product.ID
		}
		err = r.db.Omit(clause.Associations).Create(&variants).Error
		if err != nil {
			return nil, eris.Wrap(err, "")
		}
	}

	if form.Attachments != nil && len(*form.Attachments) > 0 {
		for _, attachment := range *form.Attachments {
			attachment.GetBlurhash()
		}
	}

	return r.GetProduct(GetProductParams{
		JwtClaimsInfo: form.JwtClaimsInfo,
		ProductID:     product.ID,
	})
}

type GetProductParams struct {
	models.JwtClaimsInfo
	ProductID string `param:"product_id" query:"product_id"`
	SlugID    string `param:"slug" query:"slug"`
}

func (r *ProductRepo) GetProduct(params GetProductParams) (*models.Product, error) {
	if params.SlugID == "" && params.ProductID == "" {
		err := errors.New("empty id")
		return nil, err
	}
	var builder = queryfunc.NewProductBuilder(queryfunc.ProductBuilderOptions{
		QueryBuilderOptions: queryfunc.QueryBuilderOptions{
			Role: params.GetRole(),
		},
		GetDetails: true,
	})
	var product models.Product
	var err = query.New(r.db, builder).
		WhereFunc(func(builder *query.Builder) {
			if params.ProductID != "" {
				builder.Where("p.id = ?", params.ProductID)
			} else if params.SlugID != "" {
				builder.Where("p.slug = ?", params.SlugID)
			}
		}).
		FirstFunc(&product)

	if err != nil {
		if r.db.IsRecordNotFoundError(err) {
			return nil, errs.ErrProductNotFound
		}
		return nil, err
	}

	return &product, nil
}

func (r *ProductRepo) UpdateProduct(form models.ProductUpdateForm) (*models.Product, error) {
	var update models.Product
	update.ID = form.ProductID
	var err = copier.Copy(&update, &form)
	if err != nil {
		return nil, eris.Wrap(err, err.Error())
	}

	if form.Attachments != nil && len(*form.Attachments) > 0 {
		var attachments models.Attachments
		for _, attachment := range *form.Attachments {
			attachment.GetBlurhash()
			attachments = append(attachments, attachment)
		}

		update.Attachments = &attachments
	}

	err = r.db.Clauses(clause.Returning{}).
		Omit(clause.Associations).Model(&update).
		Where("id = ?", form.ProductID).Updates(&update).Error
	if err != nil {
		return nil, err
	}

	// CreateFromPayload variants
	if len(form.Variants) > 0 {
		var variants []*models.Variant
		err = copier.Copy(&variants, &form.Variants)
		if err != nil {
			return nil, eris.Wrap(err, "")
		}
		for _, variant := range variants {
			variant.ProductID = update.ID
		}
		err = r.db.Transaction(func(tx *gorm.DB) (e error) {
			if e = r.db.Unscoped().Delete(&models.Variant{}, "product_id = ?", update.ID).Error; e != nil {
				return
			}
			if e = tx.Create(&variants).Error; e != nil {
				return
			}
			return
		})
		if err != nil {
			return nil, eris.Wrap(err, "")
		}
	}

	return r.GetProduct(GetProductParams{
		JwtClaimsInfo: form.JwtClaimsInfo,
		ProductID:     form.ProductID,
	})
}

type DeleteProductParams struct {
	models.JwtClaimsInfo
	ProductID string `param:"product_id" validate:"required"`
}

func (r *ProductRepo) DeleteProduct(params DeleteProductParams) error {
	return r.db.Unscoped().Transaction(func(tx *gorm.DB) error {
		var err = tx.Delete(&models.Product{}, "id = ?", params.ProductID).Error
		if err != nil {
			return err
		}

		return tx.Unscoped().Delete(&models.ProductAttribute{}, "product_id = ?", params.ProductID).Error
	})
}

// Product Variant

func (r *ProductRepo) CreateProductAttribute(form models.ProductAttributeUpdateForm) (*models.ProductAttribute, error) {
	var option models.ProductAttribute
	err := copier.Copy(&option, &form)
	if err != nil {
		return nil, eris.Wrap(err, "")
	}

	err = r.db.Omit(clause.Associations).Create(&option).Error
	if err != nil {
		if r.db.IsRecordNotFoundError(err) {
			return nil, errs.ErrProductExisted
		}
		return nil, eris.Wrap(err, "")
	}

	return &option, nil
}

func (r *ProductRepo) GetProductAttributeByID(optionID string, options queryfunc.ProductBuilderOptions) (*models.ProductAttribute, error) {
	var builder = queryfunc.NewProductBuilder(options)
	var option models.ProductAttribute
	var err = query.New(r.db, builder).
		WhereFunc(func(builder *query.Builder) {
			builder.Where("p.id = ?", optionID)
		}).
		FirstFunc(&option)

	if err != nil {
		if r.db.IsRecordNotFoundError(err) {
			return nil, errs.ErrProductNotFound
		}
		return nil, err
	}

	return &option, nil
}

func (r *ProductRepo) GetProductByIDs(ProductIDs []string, options queryfunc.ProductBuilderOptions) ([]*models.Product, error) {
	var builder = queryfunc.NewProductBuilder(options)
	var products []*models.Product
	var err = query.New(r.db, builder).
		WhereFunc(func(builder *query.Builder) {
			builder.Where("p.id IN (?)", ProductIDs)
		}).
		FindFunc(&products)

	if err != nil {
		if r.db.IsRecordNotFoundError(err) {
			return nil, errs.ErrCategoryNotFound
		}
		return nil, err
	}

	return products, nil
}

func (r *ProductRepo) GetProductByPageSectionType(Type string, options queryfunc.ProductBuilderOptions) (products []*models.Product, err error) {
	var builder = queryfunc.NewPageSectionBuilder(queryfunc.PageSectionBuilderOptions{})
	var PageSection models.PageSection
	err = query.New(r.db, builder).
		WhereFunc(func(builder *query.Builder) {
			builder.Where("p.section_type = ?", Type)
		}).
		FirstFunc(&PageSection)
	if err != nil {
		if r.db.IsRecordNotFoundError(err) {
			err = errs.ErrRecordNotFound
		}
		return nil, eris.Wrap(err, fmt.Sprintf("cannot find page_section_type %s", Type))
	}
	ids := PageSection.ProductIds
	products, err = NewProductRepo(r.db).GetProductByIDs(ids, queryfunc.ProductBuilderOptions{})
	return
}

type PaginateCollectionProductParams struct {
	PaginateProductParams

	CollectionID string `param:"collection_id" validate:"required"`
}

func (r *ProductRepo) PaginateCollectionProduct(params PaginateCollectionProductParams) *query.Pagination {
	var builder = queryfunc.NewProductBuilder(queryfunc.ProductBuilderOptions{
		QueryBuilderOptions: queryfunc.QueryBuilderOptions{
			Role: params.GetRole(),
		},
		GetByCollection: true,
	})

	if params.Limit == 0 {
		params.Limit = 20
	}

	var result = query.New(r.db, builder).
		WhereFunc(func(builder *query.Builder) {
			builder.Where("c.collection_id = ?", params.CollectionID)
		}).
		Page(params.Page).
		Limit(params.Limit).
		PagingFunc()

	return result
}

type GetAllCollectionProductParams struct {
	CollectionID string `param:"collection_id" validate:"required"`

	models.JwtClaimsInfo
}

func (r *ProductRepo) GetAllCollectionProduct(params GetAllCollectionProductParams) []*models.Product {
	var builder = queryfunc.NewProductBuilder(queryfunc.ProductBuilderOptions{
		QueryBuilderOptions: queryfunc.QueryBuilderOptions{
			Role: params.GetRole(),
		},
		GetByCollection: true,
	})

	var result []*models.Product
	query.New(r.db, builder).
		WhereFunc(func(builder *query.Builder) {
			builder.Where("c.collection_id = ?", params.CollectionID)
		}).
		FindFunc(&result)

	return result
}

type PaginateShopProductParams struct {
	models.PaginationParams
	models.JwtClaimsInfo

	ShopID string `param:"shop_id"`
}

func (r *ProductRepo) PaginateShopProduct(params PaginateShopProductParams) *query.Pagination {
	var builder = queryfunc.NewProductBuilder(queryfunc.ProductBuilderOptions{
		QueryBuilderOptions: queryfunc.QueryBuilderOptions{
			Role: params.GetRole(),
		},
	})

	if params.Limit == 0 {
		params.Limit = 20
	}

	var result = query.New(r.db, builder).
		WhereFunc(func(builder *query.Builder) {
			builder.Where("p.shop_id = ?", params.ShopID)
		}).
		Page(params.Page).
		Limit(params.Limit).
		PagingFunc()

	return result
}

func (r *ProductRepo) GenerateSlug() (err error) {
	err = r.db.Transaction(func(tx *gorm.DB) (e error) {
		var products models.Products
		r.db.Find(&products)
		for i, pd := range products {
			e = tx.Model(pd).Where("id = ?", pd.ID).UpdateColumn("slug", fmt.Sprintf("%d", i)).Error
			if e != nil {
				return
			}
		}
		uniquePath := make(map[string]bool)
		for _, pd := range products {
			path := slug.Make(pd.Name)
			origin := path
			isExits := true
			inc := 0
			for isExits {
				isExits = uniquePath[path]
				if isExits {
					inc += 1
					path = fmt.Sprintf("%s-%d", origin, inc)
				}
			}
			e = tx.Model(pd).Where("id = ?", pd.ID).UpdateColumn("slug", path).Error
			if e != nil {
				return
			}
			uniquePath[path] = true
		}
		return
	})
	return
}

func (r *ProductRepo) GetQRCode(params models.GetProductQRCodeParams) (url string, err error) {
	var s3Client = s3.New(r.db.Configuration)

	if params.Logo == "" {
		err = errors.New("empty logo")
		return
	}
	if params.Bucket == "" {
		err = errors.New("empty bucket")
		return
	}
	var find models.Product
	err = r.db.Select("id", "qr_code").First(&find, "id = ?", params.ProductId).Error
	if err != nil {
		return
	}
	if find.ID == "" {
		err = errs.ErrProductNotFound
		return
	}
	if !params.Override && find.QRCode != "" {
		url = find.QRCode
		return
	}

	buf, err := helper.GenerateQRCode(helper.GenerateQRCodeOptions{
		Content: params.URL,
		ImageOptions: []standard.ImageOption{
			standard.WithLogoImageFileJPEG(params.Logo),
			standard.WithLogoSizeMultiplier(1),
		},
	})
	if err != nil {
		return
	}

	var contentType = models.ContentTypeImageJPG
	url = fmt.Sprintf("uploads/media/qr_code/product-%s%s", params.ProductId, contentType.GetExtension())
	_, _ = s3Client.UploadFile(s3.UploadFileParams{
		Data:        bytes.NewReader(buf.Bytes()),
		Bucket:      params.Bucket,
		ContentType: string(contentType),
		ACL:         "private",
		Key:         url,
	})
	err = r.db.Model(&models.Product{}).Where("id = ?", params.ProductId).UpdateColumn("qr_code", url).Error
	return
}

func (r *ProductRepo) ExportExcel(params PaginateProductParams) (*models.Attachment, error) {
	var result = r.PaginateProducts(params)
	if result == nil || result.Records == nil {
		return nil, errors.New("empty response")
	}
	trans, ok := result.Records.([]*models.Product)
	if !ok {
		return nil, errors.New("empty response")
	}

	fileContent, err := models.Products(trans).ToExcel()
	if err != nil {
		return nil, err
	}

	var contentType = models.ContentTypeXLSX
	url := fmt.Sprintf("uploads/products/export/export_products_user_%s%s", params.GetUserID(), contentType.GetExtension())
	_, err = s3.New(r.db.Configuration).UploadFile(s3.UploadFileParams{
		Data:        bytes.NewReader(fileContent),
		Bucket:      r.db.Configuration.AWSS3StorageBucket,
		ContentType: string(contentType),
		ACL:         "private",
		Key:         url,
	})
	if err != nil {
		return nil, err
	}
	var resp = models.Attachment{
		FileKey:     url,
		ContentType: string(contentType),
	}
	return &resp, err
}
