package repo

import (
	"github.com/engineeringinflow/inflow-backend/pkg/db"
	"github.com/engineeringinflow/inflow-backend/pkg/errs"
	"github.com/engineeringinflow/inflow-backend/pkg/helper"
	"github.com/engineeringinflow/inflow-backend/pkg/logger"
	"github.com/engineeringinflow/inflow-backend/pkg/models"
	"github.com/engineeringinflow/inflow-backend/pkg/models/enums"
	"github.com/engineeringinflow/inflow-backend/pkg/repo/query"
	"github.com/engineeringinflow/inflow-backend/pkg/repo/query/queryfunc"

	"github.com/jinzhu/copier"
	"github.com/lib/pq"
	"github.com/rotisserie/eris"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type DocumentRepo struct {
	db     *db.DB
	logger *logger.Logger
}

func NewDocumentRepo(db *db.DB) *DocumentRepo {
	return &DocumentRepo{
		db:     db,
		logger: logger.New("repo/Document"),
	}
}

func (r *DocumentRepo) CreateDocument(req *models.CreateDocumentRequest) (*models.Document, error) {
	document := &models.Document{}
	err := copier.Copy(document, req)
	if err != nil {
		return nil, eris.Wrap(err, "copy attribute error")
	}
	if document.Vi != nil && document.Vi.Status == "" {
		document.Vi.Status = enums.DocumentStatusNew
	}

	if err := r.db.Select("ID").First(&models.DocumentCategory{}, "id = ?", document.CategoryID).Error; err != nil {
		if r.db.IsRecordNotFoundError(err) {
			return nil, errs.ErrDocumentCateNotFound
		}
		return nil, err
	}
	if len(req.TagIDs) > 0 {
		var tags = make(models.DocumentTags, 0, len(req.TagIDs))
		if err := r.db.Find(&tags, "id IN ?", req.TagIDs).Error; err != nil {
			return nil, err
		}
		var dbTagIDs = tags.IDs()
		for _, id := range req.TagIDs {
			if !helper.StringContains(dbTagIDs, id) {
				return nil, eris.Wrapf(errs.ErrDocumentTagNotFound, "document_tag_id:%s", id)
			}
		}
	}

	if err := r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Omit(clause.Associations).Create(document).Error; err != nil {
			return err
		}
		if len(req.TagIDs) > 0 {
			var taggedDocumentsToCreate = make([]*models.TaggedDocument, 0, len(req.TagIDs))
			for _, tagID := range req.TagIDs {
				taggedDocumentsToCreate = append(taggedDocumentsToCreate, &models.TaggedDocument{
					DocumentID:    document.ID,
					DocumentTagID: tagID,
				})
			}
			return tx.Create(taggedDocumentsToCreate).Error
		}
		return nil
	}); err != nil {
		return nil, err
	}

	return document, nil
}

func (r *DocumentRepo) UpdateDocument(req *models.UpdateDocumentRequest) (*models.Document, error) {
	var document = &models.Document{}
	err := copier.Copy(document, req)
	if err != nil {
		return nil, eris.Wrap(err, "copy attribute error")
	}
	document.ID = req.DocumentID

	if document.Vi != nil && document.Vi.Status == "" {
		document.Vi.Status = enums.DocumentStatusNew
	}

	if err := r.db.Select("ID").First(&models.Document{}, "id = ?", document.ID).Error; err != nil {
		if r.db.IsRecordNotFoundError(err) {
			return nil, errs.ErrDocumentNotFound
		}
		return nil, err
	}

	if err := r.db.Select("ID").First(&models.DocumentCategory{}, "id = ?", document.CategoryID).Error; err != nil {
		if r.db.IsRecordNotFoundError(err) {
			return nil, errs.ErrDocumentCateNotFound
		}
		return nil, err
	}

	if len(req.TagIDs) > 0 {
		var tags = make(models.DocumentTags, 0, len(req.TagIDs))
		if err := r.db.Select("ID").Find(&tags, "id IN ?", req.TagIDs).Error; err != nil {
			return nil, err
		}
		var dbTagIDs = tags.IDs()
		for _, id := range req.TagIDs {
			if !helper.StringContains(dbTagIDs, id) {
				return nil, eris.Wrapf(errs.ErrDocumentTagNotFound, "document_tag_id:%s", id)
			}
		}
	}

	if err := r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Omit(clause.Associations).Model(&document).Where("id = ?", document.ID).
			Updates(&document).Error; err != nil {
			return err
		}
		if err := tx.Unscoped().Delete(&models.TaggedDocument{}, "document_id = ?", document.ID).Error; err != nil {
			return err
		}
		if len(req.TagIDs) > 0 {
			var taggedDocumentsToCreate = make([]*models.TaggedDocument, 0, len(req.TagIDs))
			for _, tagID := range req.TagIDs {
				taggedDocumentsToCreate = append(taggedDocumentsToCreate, &models.TaggedDocument{
					DocumentID:    document.ID,
					DocumentTagID: tagID,
				})
			}
			return tx.Create(taggedDocumentsToCreate).Error
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return document, nil
}

func (r *DocumentRepo) GetDocument(params *models.GetDocumentParams) (*models.Document, error) {
	var builder = queryfunc.NewDocumentBuilder(queryfunc.DocumentBuilderOptions{
		QueryBuilderOptions: queryfunc.QueryBuilderOptions{
			Role: params.GetRole(),
		},
	})
	var document models.Document
	var err = query.New(r.db, builder).
		WhereFunc(func(builder *query.Builder) {
			builder.Where("(d.slug = ? or d.vi->>'slug' = ?)", params.Slug, params.Slug)
		}).
		FirstFunc(&document)

	if err != nil {
		if r.db.IsRecordNotFoundError(err) {
			return nil, errs.ErrRecordNotFound
		}
		return nil, err
	}

	return &document, nil
}

func (r *DocumentRepo) DeleteDocument(params *models.DeleteDocumentParams) error {

	var sqlResult = r.db.Unscoped().Delete(&models.Document{}, "id = ?", params.DocumentID)

	if sqlResult.RowsAffected == 0 {
		return errs.ErrRecordNotFound
	}
	if sqlResult.Error != nil {
		return sqlResult.Error
	}

	return nil
}

func (r *DocumentRepo) GetDocumentList(params *models.GetDocumentListParams) *query.Pagination {
	var builder = queryfunc.NewDocumentBuilder(queryfunc.DocumentBuilderOptions{
		QueryBuilderOptions: queryfunc.QueryBuilderOptions{
			Role: params.GetRole(),
		},
	})

	if params.Limit == 0 {
		params.Limit = 20
	}

	var result = query.New(r.db, builder).
		WhereFunc(func(builder *query.Builder) {
			if len(params.CategoryIDs) > 0 {
				builder.Where("d.category_id IN ?", params.CategoryIDs)
			}
			if len(params.Statuses) > 0 {
				switch params.Language {
				case enums.LanguageCodeVietnam:
					builder.Where("vi ->> 'status' IN ?", params.Statuses)
				case enums.LanguageCodeEnglish:
					builder.Where("d.status IN ?", params.Statuses)
				default:
					builder.Where("(d.status IN ? OR vi ->> 'status' IN ?)", params.Statuses, params.Statuses)
				}
			}
			if len(params.TagIDs) > 0 {
				builder.Where("td.document_tag_id IN ?", params.TagIDs)
			}
			if len(params.Keyword) > 0 {
				q := "%" + params.Keyword + "%"
				builder.Where("(d.title ILIKE ? OR unaccent(vi ->> 'title') ILIKE unaccent(?))", q, q)
			}
			if len(params.Roles) > 0 {
				builder.Where("count_elements(d.visible_to,?) >= 1", pq.StringArray(params.Roles))
			}
		}).
		Page(params.Page).
		Limit(params.Limit).
		PagingFunc()

	return result
}
