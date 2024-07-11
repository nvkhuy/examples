package repo

import (
	"github.com/engineeringinflow/inflow-backend/pkg/db"
	"github.com/engineeringinflow/inflow-backend/pkg/errs"
	"github.com/engineeringinflow/inflow-backend/pkg/logger"
	"github.com/engineeringinflow/inflow-backend/pkg/models"
	"github.com/engineeringinflow/inflow-backend/pkg/models/enums"
	"github.com/engineeringinflow/inflow-backend/pkg/repo/query"
	"github.com/engineeringinflow/inflow-backend/pkg/repo/query/queryfunc"
	"github.com/jinzhu/copier"

	"github.com/rotisserie/eris"

	"gorm.io/gorm/clause"
)

type PageSectionRepo struct {
	db     *db.DB
	logger *logger.Logger
}

func NewPageSectionRepo(db *db.DB) *PageSectionRepo {
	return &PageSectionRepo{
		db:     db,
		logger: logger.New("repo/PageSection"),
	}
}

type PaginatePageSectionParams struct {
	models.PaginationParams

	Name    string `json:"name" query:"name" form:"name"`
	ForRole enums.Role
}

// func (r *PageSectionRepo) PaginateCategories(params PaginatePageSectionParams) *query.Pagination {
// 	var builder = queryfunc.NewUserBuilder(queryfunc.PageSectionBuilderOptions{
// 		ForRole: params.ForRole,
// 	})

// 	var result = query.New(r.adb, builder).
// 		WhereFunc(func(builder *query.Builder) { }).
// 		PageSection(params.PageSection).
// 		Limit(params.Limit).
// 		PagingFunc()

// 	return result
// }

func (r *PageSectionRepo) CreatePageSection(form models.PageSectionCreateForm) (*models.PageSection, error) {
	var PageSection models.PageSection
	err := copier.Copy(&PageSection, &form)
	if err != nil {
		return nil, eris.Wrap(err, "")
	}

	err = r.db.Omit(clause.Associations).Create(&PageSection).Error
	if err != nil {
		if r.db.IsRecordNotFoundError(err) {
			return nil, errs.ErrRecordNotFound
		}
		return nil, eris.Wrap(err, "")
	}

	return &PageSection, nil
}

func (r *PageSectionRepo) GetPageSectionByID(PageSectionID string, options queryfunc.PageSectionBuilderOptions) (*models.PageSection, error) {
	var builder = queryfunc.NewPageSectionBuilder(options)
	var PageSection models.PageSection
	var err = query.New(r.db, builder).
		WhereFunc(func(builder *query.Builder) {
			builder.Where("p.id = ?", PageSectionID)
		}).
		FirstFunc(&PageSection)

	if err != nil {
		if r.db.IsRecordNotFoundError(err) {
			return nil, errs.ErrRecordNotFound
		}
		return nil, err
	}

	return &PageSection, nil
}

func (r *PageSectionRepo) UpdatePageSectionByID(PageSectionID string, form models.PageSectionUpdateForm) (*models.PageSection, error) {
	var update models.PageSection

	var err = copier.Copy(&update, &form)
	if err != nil {
		return nil, eris.Wrap(err, err.Error())
	}

	err = r.db.Omit(clause.Associations).Model(&models.PageSection{}).Where("id = ?", PageSectionID).Updates(&update).Error
	if err != nil {
		return nil, err
	}

	return r.GetPageSectionByID(PageSectionID, queryfunc.PageSectionBuilderOptions{
		QueryBuilderOptions: queryfunc.QueryBuilderOptions{
			Role: form.GetRole(),
		},
	})
}

func (r *PageSectionRepo) DeletePageSectionByID(PageSectionID string) error {

	var err = r.db.Unscoped().Delete(&models.PageSection{}, "id = ?", PageSectionID).Error

	return err
}

func (r *PageSectionRepo) GetSectionByPageID(PageID string, options queryfunc.PageSectionBuilderOptions) ([]*models.PageSection, error) {
	var builder = queryfunc.NewPageSectionBuilder(queryfunc.PageSectionBuilderOptions{})
	var sections []*models.PageSection
	var err = query.New(r.db, builder).
		WhereFunc(func(builder *query.Builder) {
			builder.Where("p.page_id = ?", PageID)
		}).
		FindFunc(&sections)

	if err != nil {
		if r.db.IsRecordNotFoundError(err) {
			return nil, errs.ErrRecordNotFound
		}
		return nil, err
	}

	return sections, nil
}
