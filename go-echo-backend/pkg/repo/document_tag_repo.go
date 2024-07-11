package repo

import (
	"github.com/engineeringinflow/inflow-backend/pkg/db"
	"github.com/engineeringinflow/inflow-backend/pkg/errs"
	"github.com/engineeringinflow/inflow-backend/pkg/logger"
	"github.com/engineeringinflow/inflow-backend/pkg/models"
	"github.com/engineeringinflow/inflow-backend/pkg/repo/query"
	"github.com/engineeringinflow/inflow-backend/pkg/repo/query/queryfunc"
	"gorm.io/gorm"
)

type DocumentTagRepo struct {
	db     *db.DB
	logger *logger.Logger
}

func NewDocumentTagRepo(db *db.DB) *DocumentTagRepo {
	return &DocumentTagRepo{
		db:     db,
		logger: logger.New("repo/document_tag"),
	}
}

func (r *DocumentTagRepo) CreateDocumentTag(req *models.CreateDocumentTagRequest) (*models.DocumentTag, error) {
	var tag = &models.DocumentTag{Name: req.Name}

	if err := r.db.Create(tag).Error; err != nil {
		ok, _ := r.db.IsDuplicateConstraint(err)
		if ok {
			return nil, errs.ErrDocumentTagExisted
		}
		return nil, err
	}

	return tag, nil
}

func (r *DocumentTagRepo) GetDocumentTagList(params *models.GetDocumentTagListParams) *query.Pagination {
	var builder = queryfunc.NewDocumentTagBuilder(queryfunc.DocumentTagBuilderOptions{
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
	var result = query.New(r.db, builder).
		WhereFunc(func(builder *query.Builder) {
			if params.Keyword != "" {
				var q = "%" + params.Keyword + "%"
				builder.Where("t.name ILIKE ?", q)
			}
		}).
		Page(params.Page).
		Limit(params.Limit).
		PagingFunc()

	return result
}

func (r *DocumentTagRepo) UpdateDocumentTag(req *models.UpdateDocumentTagRequest) (*models.DocumentTag, error) {
	var tag = &models.DocumentTag{
		Model: models.Model{ID: req.DocumentTagID},
		Name:  req.Name,
	}
	var sqlResult = r.db.
		Model(&models.DocumentTag{}).Where("id = ?", req.DocumentTagID).Updates(tag)

	if sqlResult.Error != nil {
		ok, _ := r.db.IsDuplicateConstraint(sqlResult.Error)
		if ok {
			return nil, errs.ErrDocumentTagExisted
		}
		return nil, sqlResult.Error
	}
	if sqlResult.RowsAffected == 0 {
		return nil, errs.ErrRecordNotFound
	}

	return tag, nil
}

func (r *DocumentTagRepo) DeleteDocumentTag(req *models.DeleteDocumentTagRequest) error {
	var err = r.db.Transaction(func(tx *gorm.DB) error {
		sqlResult := tx.Unscoped().Delete(&models.DocumentTag{}, "id = ?", req.DocumentTagID)
		if sqlResult.RowsAffected == 0 {
			return errs.ErrRecordNotFound
		}
		if sqlResult.Error != nil {
			return sqlResult.Error
		}
		return tx.Unscoped().Delete(&models.TaggedDocument{}, "document_tag_id = ?", req.DocumentTagID).Error
	})
	return err
}
