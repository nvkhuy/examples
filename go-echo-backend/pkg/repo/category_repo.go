package repo

import (
	"fmt"
	"sort"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/engineeringinflow/inflow-backend/pkg/db"
	"github.com/engineeringinflow/inflow-backend/pkg/errs"
	"github.com/engineeringinflow/inflow-backend/pkg/logger"
	"github.com/engineeringinflow/inflow-backend/pkg/models"
	"github.com/engineeringinflow/inflow-backend/pkg/repo/query"
	"github.com/engineeringinflow/inflow-backend/pkg/repo/query/queryfunc"
	"github.com/gosimple/slug"
	"github.com/jinzhu/copier"
	"github.com/lib/pq"
	"github.com/rotisserie/eris"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type CategoryRepo struct {
	db     *db.DB
	logger *logger.Logger
}

func NewCategoryRepo(db *db.DB) *CategoryRepo {
	return &CategoryRepo{
		db:     db,
		logger: logger.New("repo/category"),
	}
}

type PaginateCategoriesParams struct {
	models.PaginationParams
	models.JwtClaimsInfo

	Name                string   `json:"name" query:"name" form:"name"`
	ParentCategoryID    string   `json:"parent_category_id" query:"parent_category_id" form:"parent_category_id"`
	OrderByTotalProduct bool     `json:"order_by_total_product" query:"order_by_total_product" form:"order_by_total_product"`
	CategoryIDs         []string `json:"category_ids" query:"category_ids" form:"category_ids"`
	ParentOnly          bool
}

func (r *CategoryRepo) PaginateCategoriesParams(params PaginateCategoriesParams) []*models.Category {
	var builder = queryfunc.NewCategoryBuilder(queryfunc.CategoryBuilderOptions{
		QueryBuilderOptions: queryfunc.QueryBuilderOptions{
			Role: params.GetRole(),
		},
	})

	var result []*models.Category
	query.New(r.db, builder).
		WhereFunc(func(builder *query.Builder) {
			if params.ParentCategoryID != "" {
				builder.Where("c.parent_category_id = ?", params.ParentCategoryID)
			}

			if len(params.CategoryIDs) > 0 {
				builder.Where("id IN ?", params.CategoryIDs)
			}
		}).
		Limit(params.Limit).
		FindFunc(&result)

	return result
}

func (r *CategoryRepo) GetParentCategoryIDs() (ids []string) {
	r.db.Raw(`
	SELECT c.parent_category_id
	FROM categories c
	WHERE COALESCE(c.parent_category_id,'') <> ''
	GROUP BY c.parent_category_id
	`).Find(&ids)

	return
}

func (r *CategoryRepo) GetCategoriesV0(params PaginateCategoriesParams) []*models.Category {
	var builder = queryfunc.NewCategoryBuilder(queryfunc.CategoryBuilderOptions{
		QueryBuilderOptions: queryfunc.QueryBuilderOptions{
			Role: params.GetRole(),
		},
	})
	if params.Limit == 0 {
		params.Limit = 1000
	}
	var result []*models.Category
	_ = query.New(r.db, builder).
		WhereFunc(func(builder *query.Builder) {
			if params.ParentOnly {
				builder.Where("cate.parent_category_id = ?", "")
			}
		}).
		Limit(params.Limit).
		FindFunc(&result)

	return result
}

func (r *CategoryRepo) GetCategories(params PaginateCategoriesParams) []*models.Category {
	var parentCategoryIDs = r.GetParentCategoryIDs()
	if len(parentCategoryIDs) == 0 {
		return []*models.Category{}
	}

	var builder = queryfunc.NewCategoryBuilder(queryfunc.CategoryBuilderOptions{
		QueryBuilderOptions: queryfunc.QueryBuilderOptions{
			Role: params.GetRole(),
		},
	})

	var result []*models.Category
	_ = query.New(r.db, builder).
		WhereFunc(func(builder *query.Builder) {
			builder.Where("id IN ?", parentCategoryIDs)
		}).
		Limit(params.Limit).
		FindFunc(&result)

	return result
}

func (r *CategoryRepo) GetCategoryTree(params PaginateCategoriesParams) models.CategoryTreeResponse {
	var result = r.GetCategoriesV0(params)

	var response models.CategoryTreeResponse
	var data []*models.CategoryResponse

	for _, cate := range result {
		if *cate.ParentCategoryID == "" {
			data = append(data, r.BuildCategories(cate, result))
		}
	}

	if params.OrderByTotalProduct {
		sort.Slice(data, func(i, j int) bool {
			return data[i].TotalProduct > data[j].TotalProduct
		})
	}

	response.Records = data

	return response
}

func (r *CategoryRepo) CreateCategory(form models.CategoryCreateForm) (*models.Category, error) {
	var category models.Category
	err := copier.Copy(&category, &form)
	if err != nil {
		return nil, eris.Wrap(err, "")
	}

	if form.Icon == nil {
		category.Icon = &models.Attachment{
			FileKey:      "icon",
			ContentType:  "",
			FileURL:      "https://dev-static.joininflow.io/common/default_category_icon.png",
			ThumbnailURL: "",
		}
	}

	// Commit
	err = r.db.Omit(clause.Associations).Clauses(clause.Returning{}).Create(&category).Error
	if err != nil {
		if r.db.IsRecordNotFoundError(err) {
			return nil, errs.ErrCategoryExisted
		}
		return nil, eris.Wrap(err, "")
	}

	var products []*models.Product
	products, err = r.updateTopProducts(form.TopProductIDs, category.ID)
	if err != nil {
		return nil, err
	}
	category.TopProducts = products

	return &category, nil
}

func (r *CategoryRepo) updateTopProducts(productIds []string, categoryId string) (result []*models.Product, err error) {
	if err = r.db.Model(&models.Category{}).Where("id = ?", categoryId).Update("top_product_ids", pq.StringArray(productIds)).Error; err != nil {
		return
	}
	err = r.db.Where("id IN ?", productIds).Find(&result).Error
	return
}

func (r *CategoryRepo) GetCategoryByID(CategoryID string, options queryfunc.CategoryBuilderOptions) (*models.Category, error) {
	var builder = queryfunc.NewCategoryBuilder(options)
	var Category models.Category
	var err = query.New(r.db, builder).
		WhereFunc(func(builder *query.Builder) {
			builder.Where("cate.id = ?", CategoryID)
		}).
		FirstFunc(&Category)

	if err != nil {
		if r.db.IsRecordNotFoundError(err) {
			return nil, errs.ErrCategoryNotFound
		}
		return nil, err
	}

	return &Category, nil
}

func (r *CategoryRepo) GetCategoryByIDs(CategoryIDs []string, options queryfunc.CategoryBuilderOptions) ([]*models.Category, error) {
	var builder = queryfunc.NewCategoryBuilder(options)
	var categories []*models.Category
	var err = query.New(r.db, builder).
		WhereFunc(func(builder *query.Builder) {
			builder.Where("cate.id IN (?)", CategoryIDs)
		}).
		FindFunc(&categories)

	if err != nil {
		if r.db.IsRecordNotFoundError(err) {
			return nil, errs.ErrCategoryNotFound
		}
		return nil, err
	}

	return categories, nil
}

func (r *CategoryRepo) UpdateCategoryByID(form models.CategoryUpdateForm) (*models.Category, error) {
	category, err := r.GetCategoryByID(form.CategoryID, queryfunc.CategoryBuilderOptions{
		QueryBuilderOptions: queryfunc.QueryBuilderOptions{
			Role: form.GetRole(),
		},
	})
	if err != nil {
		return nil, err
	}

	var update models.Category
	err = copier.Copy(&update, &form)
	if err != nil {
		return nil, eris.Wrap(err, err.Error())
	}
	update.ID = form.CategoryID

	// Validate parent category
	parentId := aws.StringValue(category.ParentCategoryID)
	var parent *models.Category
	if parentId != "" {
		if err = r.db.Select("id", "name", "parent_category_id").First(&parent, "id = ?", parentId).Error; err != nil {
			return nil, errs.ErrParentCategoryNotFound
		}
		grandpaId := aws.StringValue(parent.ParentCategoryID)
		if grandpaId != "" {
			return nil, errs.ErrInvalidParentCategory
		}
	}

	err = r.db.Omit(clause.Associations).Model(&update).Where("id = ?", form.CategoryID).Updates(&update).Error
	if err != nil {
		return nil, err
	}
	err = copier.CopyWithOption(category, &update, copier.Option{IgnoreEmpty: true, DeepCopy: true})
	if err != nil {
		return nil, err
	}

	var products []*models.Product
	products, err = r.updateTopProducts(form.TopProductIDs, form.CategoryID)
	if err != nil {
		return nil, err
	}
	category.TopProducts = products

	return category, err
}

type DeleteCategoryParams struct {
	models.JwtClaimsInfo

	CategoryID string `param:"category_id" validate:"required"`
}

func (r *CategoryRepo) DeleteCategory(params DeleteCategoryParams) error {

	var childCategory models.Category
	var err = r.db.First(&childCategory, "parent_category_id = ?", params.CategoryID).Error
	if childCategory.ID != "" {
		return errs.ErrCategoryChildExisted
	}

	err = r.db.Transaction(func(tx *gorm.DB) error {
		err = tx.Model(&models.Product{}).Where("category_id = ?", params.CategoryID).UpdateColumn("CategoryID", "").Error
		if err != nil {
			return eris.Wrap(err, "")
		}

		err = tx.Unscoped().Delete(&models.Category{}, "id = ?", params.CategoryID).Error

		return err
	})

	return err
}

func (r *CategoryRepo) BuildCategories(parentCate *models.Category, categories []*models.Category) *models.CategoryResponse {
	c := &models.CategoryResponse{
		ID:               parentCate.ID,
		Name:             parentCate.Name,
		Slug:             parentCate.Slug,
		ParentCategoryID: *parentCate.ParentCategoryID,
		Icon:             parentCate.Icon,
		Vi:               parentCate.Vi,
	}

	for _, cate := range categories {
		if *cate.ParentCategoryID == parentCate.ID {
			cateChild := &models.CategoryResponse{
				ID:               cate.ID,
				Name:             cate.Name,
				Slug:             cate.Slug,
				ParentCategoryID: *cate.ParentCategoryID,
				Icon:             cate.Icon,
				Vi:               cate.Vi,
			}
			var childCount int
			r.db.Model(&models.Products{}).Select("count(1)").Where("category_id = ?", cate.ID).Find(&childCount)
			c.TotalProduct += childCount

			var products []*models.Product
			if len(cate.TopProductIDs) > 0 {
				if err := r.db.Where("id IN ?", []string(cate.TopProductIDs)).Find(&products).Error; err == nil {
					cateChild.TopProductIds = cate.TopProductIDs
					cateChild.TopProducts = products
				}
			}
			c.Children = append(c.Children, cateChild)
		}
	}
	return c
}

func (r *CategoryRepo) GetChildCategoryIDs(parentCateID string) []string {
	var c []string

	var children []*models.Category
	var builder = queryfunc.NewCategoryBuilder(queryfunc.CategoryBuilderOptions{})
	query.New(r.db, builder).
		WhereFunc(func(builder *query.Builder) {
			builder.Where("cate.parent_category_id = ?", parentCateID)
		}).
		FindFunc(&children)

	for _, cate := range children {
		c = append(c, cate.ID)
		c = append(c, r.GetChildCategoryIDs(cate.ID)...)
	}
	return c
}

type GenerateCategorySlugParams struct {
	models.JwtClaimsInfo
}

func (r *CategoryRepo) GenerateSlug() (err error) {
	err = r.db.Transaction(func(tx *gorm.DB) (e error) {
		var cates models.CategorySlice
		r.db.Find(&cates)
		for i, cate := range cates {
			e = tx.Model(&models.Category{}).Where("id = ?", cate.ID).UpdateColumn("Slug", fmt.Sprintf("%d", i)).Error
			if e != nil {
				return
			}
		}
		m := make(map[string]*models.Category)
		for _, cate := range cates {
			m[cate.ID] = cate
		}
		uniquePath := make(map[string]bool)
		for _, cate := range cates {
			path := cate.Name
			parentId := aws.StringValue(cate.ParentCategoryID)
			if parentId != "" {
				if pr, ok := m[parentId]; ok {
					path = fmt.Sprintf("%s and %s", pr.Name, path)
				}
			}
			path = slug.Make(path)
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
			e = tx.Model(&models.Category{}).Where("id = ?", cate.ID).UpdateColumn("Slug", path).Error
			if e != nil {
				return
			}
			uniquePath[path] = true
		}
		return
	})
	return
}
