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

type PageRepo struct {
	db     *db.DB
	logger *logger.Logger
}

func NewPageRepo(db *db.DB) *PageRepo {
	return &PageRepo{
		db:     db,
		logger: logger.New("repo/Page"),
	}
}

type PaginatePageParams struct {
	models.PaginationParams

	Name    string `json:"name" query:"name" form:"name"`
	ForRole enums.Role
}

type SearchPageParams struct {
	ForRole enums.Role
}

// func (r *PageRepo) PaginateCategories(params PaginatePageParams) *query.Pagination {
// 	var builder = queryfunc.NewUserBuilder(queryfunc.PageBuilderOptions{
// 		ForRole: params.ForRole,
// 	})

// 	var result = query.New(r.adb, builder).
// 		WhereFunc(func(builder *query.Builder) { }).
// 		Page(params.Page).
// 		Limit(params.Limit).
// 		PagingFunc()

// 	return result
// }

func (r *PageRepo) ListPage(params SearchPageParams) []*models.Page {
	var builder = queryfunc.NewPageBuilder(queryfunc.PageBuilderOptions{
		QueryBuilderOptions: queryfunc.QueryBuilderOptions{
			Role: params.ForRole,
		},
	})
	var result []*models.Page
	query.New(r.db, builder).
		WhereFunc(func(builder *query.Builder) {
		}).
		FindFunc(&result)

	return result
}

func (r *PageRepo) CreatePage(form models.PageUpdateForm) (*models.Page, error) {
	var Page models.Page
	err := copier.Copy(&Page, &form)
	if err != nil {
		return nil, eris.Wrap(err, "")
	}

	err = r.db.Omit(clause.Associations).Where(models.Page{PageType: Page.PageType}).FirstOrCreate(&Page).Error
	if err != nil {
		if r.db.IsRecordNotFoundError(err) {
			return nil, errs.ErrRecordNotFound
		}
		return nil, eris.Wrap(err, "")
	}

	return &Page, nil
}

type GetPageByIDParams struct {
	models.JwtClaimsInfo

	PageID string
}

func (r *PageRepo) GetPageByID(params GetPageByIDParams) (*models.Page, error) {
	var builder = queryfunc.NewPageBuilder(queryfunc.PageBuilderOptions{
		QueryBuilderOptions: queryfunc.QueryBuilderOptions{
			Role: params.GetRole(),
		},
	})
	var Page models.Page
	var err = query.New(r.db, builder).
		WhereFunc(func(builder *query.Builder) {
			builder.Where("p.id = ?", params.PageID)
		}).
		FirstFunc(&Page)

	if err != nil {
		if r.db.IsRecordNotFoundError(err) {
			return nil, errs.ErrRecordNotFound
		}
		return nil, err
	}

	return &Page, nil
}

type GetPageDetailByIDParams struct {
	models.JwtClaimsInfo

	PageID string `json:"page_id" param:"page_id" validate:"required"`
}

func (r *PageRepo) GetPageDetailByID(params GetPageDetailByIDParams) (*models.PageDetailResponse, error) {
	var page *models.Page
	page, err := r.GetPageByID(GetPageByIDParams{
		JwtClaimsInfo: params.JwtClaimsInfo,
		PageID:        params.PageID,
	})
	if err != nil {
		return nil, err
	}

	var builder = queryfunc.NewPageSectionBuilder(queryfunc.PageSectionBuilderOptions{})
	var sections []*models.PageSection
	err = query.New(r.db, builder).
		WhereFunc(func(builder *query.Builder) {
			builder.Where("p.page_id = ?", params.PageID)
		}).
		FindFunc(&sections)

	if err != nil {
		if r.db.IsRecordNotFoundError(err) {
			return nil, errs.ErrRecordNotFound
		}
		return nil, err
	}

	var detailResponse = &models.PageDetailResponse{
		ID:      page.ID,
		Title:   page.Title,
		Url:     page.Url,
		Content: sections,
	}

	return detailResponse, err
}

func (r *PageRepo) GetPageByType(PageType string, options queryfunc.PageBuilderOptions) (*models.Page, error) {
	var builder = queryfunc.NewPageBuilder(options)
	var Page models.Page
	var err = query.New(r.db, builder).
		WhereFunc(func(builder *query.Builder) {
			builder.Where("p.page_type = ?", PageType)
		}).
		FirstFunc(&Page)

	if err != nil {
		if r.db.IsRecordNotFoundError(err) {
			return nil, errs.ErrRecordNotFound
		}
		return nil, err
	}

	return &Page, nil
}

func (r *PageRepo) UpdatePageByID(pageID string, form models.PageWithSectionUpdateForm) (*models.Page, error) {
	var update models.Page

	var err = copier.Copy(&update, &form)
	if err != nil {
		return nil, eris.Wrap(err, err.Error())
	}

	err = r.db.Omit(clause.Associations).Model(&models.Page{}).Where("id = ?", pageID).Updates(&update).Error
	if err != nil {
		return nil, err
	}

	// Update/create/delete sections
	var activeIDs []string
	for _, st := range form.Content {
		// CreateFromPayload new section
		if st.ID == "" {
			var update models.PageSection
			err = copier.Copy(&update, &st)
			if err != nil {
				return nil, eris.Wrap(err, err.Error())
			}
			update.PageID = pageID
			// update.Order = idx
			err = r.db.Omit(clause.Associations).Where(models.PageSection{}).Create(&update).Error
			if err != nil {
				return nil, eris.Wrap(err, err.Error())
			}
			activeIDs = append(activeIDs, update.ID)
		} else {
			// Update exist section
			activeIDs = append(activeIDs, st.ID)
			var update models.PageSection
			var err = copier.Copy(&update, &st)
			if err != nil {
				return nil, eris.Wrap(err, err.Error())
			}
			// update.Order = idx
			err = r.db.Omit(clause.Associations).Model(&models.PageSection{}).Where("id = ?", st.ID).Updates(&update).Error
			if err != nil {
				return nil, eris.Wrap(err, err.Error())
			}
		}
	}

	err = r.db.Unscoped().Delete(&models.PageSection{}, "page_id = ? and id NOT IN (?)", pageID, activeIDs).Error
	if err != nil {
		return nil, eris.Wrap(err, err.Error())
	}

	return r.GetPageByID(GetPageByIDParams{
		JwtClaimsInfo: form.JwtClaimsInfo,
		PageID:        pageID,
	})
}

func (r *PageRepo) DeletePageByID(PageID string) error {

	var err = r.db.Unscoped().Delete(&models.Page{}, "id = ?", PageID).Error

	return err
}

func (r *PageRepo) PageCatalog(PageID string) ([]*models.PageSection, error) {
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

	return sections, err
}

func (r *PageRepo) PageHome(PageID string) ([]*models.PageSection, error) {
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

	// for _, section := range sections {
	// 	var metadata interface{}

	// 	// Home Top
	// 	if section.SectionType == enums.PageSectionHomeTop {
	// 		var images []*models.Attachment
	// 		jsonIDs, _ := section.Metadata.MarshalJSON()
	// 		if err := json.Unmarshal(jsonIDs, &images); err != nil {
	// 			images = []*models.Attachment{}
	// 		}
	// 		metadata = images
	// 	}

	// 	// Home client
	// 	if section.SectionType == enums.PageSectionHomeClient {
	// 		var images []*models.Attachment
	// 		jsonIDs, _ := section.Metadata.MarshalJSON()
	// 		if err := json.Unmarshal(jsonIDs, &images); err != nil {
	// 			images = []*models.Attachment{}
	// 		}
	// 		metadata = images
	// 	}

	// 	// Home collections
	// 	if section.SectionType == enums.PageSectionHomeCollection {
	// 		var ids []string
	// 		jsonIDs, _ := section.Metadata.MarshalJSON()
	// 		if err := json.Unmarshal(jsonIDs, &ids); err != nil {
	// 			ids = []string{}
	// 		}

	// 		collections, err := NewCollectionRepo(r.db).GetCollectionByIDs(ids, queryfunc.CollectionBuilderOptions{})
	// 		if err != nil {
	// 			return nil, err
	// 		}
	// 		metadata = collections
	// 	}

}
