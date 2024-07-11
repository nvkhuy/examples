package repo

import (
	"net/http"

	"github.com/engineeringinflow/inflow-backend/pkg/db"
	"github.com/engineeringinflow/inflow-backend/pkg/errs"
	"github.com/engineeringinflow/inflow-backend/pkg/logger"
	"github.com/engineeringinflow/inflow-backend/pkg/models"
	"github.com/engineeringinflow/inflow-backend/pkg/repo/query"
	"github.com/engineeringinflow/inflow-backend/pkg/repo/query/queryfunc"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	simpleslug "github.com/gosimple/slug"
	"github.com/jinzhu/copier"
	"github.com/rotisserie/eris"
)

type DocumentCategoryRepo struct {
	db     *db.DB
	logger *logger.Logger
}

func NewDocumentCategoryRepo(db *db.DB) *DocumentCategoryRepo {
	return &DocumentCategoryRepo{
		db:     db,
		logger: logger.New("repo/document_category"),
	}
}

func (r *DocumentCategoryRepo) CreateDocumentCategory(req *models.CreateDocumentCategoryRequest) (*models.DocumentCategory, error) {
	category := &models.DocumentCategory{}
	err := copier.Copy(category, req)
	if err != nil {
		return nil, eris.Wrap(err, "copy attribute error")
	}
	category.Slug = simpleslug.Make(category.Name)

	if err = r.db.Create(category).Error; err != nil {
		ok, _ := r.db.IsDuplicateConstraint(err)
		if ok {
			return nil, errs.New(http.StatusBadRequest, "Category name is already existed")
		}
		return nil, eris.Wrap(err, "create document category error")
	}

	return category, nil
}

func (r *DocumentCategoryRepo) UpdateDocumentCategory(req *models.UpdateDocumentCategoryRequest) (*models.DocumentCategory, error) {
	var category = &models.DocumentCategory{}
	err := copier.Copy(category, req)
	if err != nil {
		return nil, eris.Wrap(err, "copy attribute error")
	}
	category.ID = req.DocumentCategoryID
	category.Slug = simpleslug.Make(category.Name)

	sqlResult := r.db.Clauses(clause.Returning{}).
		Model(&category).Where("id = ?", category.ID).Updates(&category)

	if sqlResult.Error != nil {
		ok, _ := r.db.IsDuplicateConstraint(sqlResult.Error)
		if ok {
			return nil, errs.New(http.StatusBadRequest, "Category name is already existed")
		}
		return nil, eris.Wrap(sqlResult.Error, "update document category error")
	}
	if sqlResult.RowsAffected == 0 {
		return nil, errs.ErrRecordNotFound
	}

	return category, nil
}

func (r *DocumentCategoryRepo) GetDocumentCategory(params *models.GetDocumentCategoryParams) (*models.DocumentCategory, error) {
	var builder = queryfunc.NewDocumentCategoryBuilder(queryfunc.DocumentCategoryBuilderOptions{
		QueryBuilderOptions: queryfunc.QueryBuilderOptions{
			Role: params.GetRole(),
		},
	})
	var category models.DocumentCategory
	var err = query.New(r.db, builder).
		WhereFunc(func(builder *query.Builder) {
			builder.Where("cate.id = ?", params.DocumentCategoryID)
		}).
		FirstFunc(&category)

	if err != nil {
		if r.db.IsRecordNotFoundError(err) {
			return nil, errs.ErrRecordNotFound
		}
		return nil, err
	}

	return &category, nil
}

func (r *DocumentCategoryRepo) DeleteDocumentCategory(params *models.GetDocumentCategoryParams) error {
	err := r.db.Transaction(func(tx *gorm.DB) error {
		sqlResult := r.db.Unscoped().Delete(&models.DocumentCategory{}, "id = ?", params.DocumentCategoryID)
		if sqlResult.RowsAffected == 0 {
			return errs.ErrRecordNotFound
		}
		if sqlResult.Error != nil {
			return sqlResult.Error
		}
		return tx.Model(&models.Document{}).Where("category_id = ?", params.DocumentCategoryID).
			UpdateColumn("category_id", "").Error
	})
	return err
}

func (r *DocumentCategoryRepo) GetDocumentCategoryList(params *models.GetDocumentCategoryListParams) []*models.DocumentCategory {
	var builder = queryfunc.NewDocumentCategoryBuilder(queryfunc.DocumentCategoryBuilderOptions{
		QueryBuilderOptions: queryfunc.QueryBuilderOptions{
			Role: params.GetRole(),
		},
	})
	if params.Limit == 0 {
		params.Limit = 100
	}
	if params.Page == 0 {
		params.Page = 1
	}
	var result []*models.DocumentCategory
	query.New(r.db, builder).
		WhereFunc(func(builder *query.Builder) {
			if params.Keyword != "" {
				var q = "%" + params.Keyword + "%"
				builder.Where("cate.name ILIKE ?", q)
			}
		}).
		Page(params.Page).
		Limit(params.Limit).
		FindFunc(&result)

	return result
}
