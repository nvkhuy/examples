package repo

import (
	"github.com/engineeringinflow/inflow-backend/pkg/db"
	"github.com/engineeringinflow/inflow-backend/pkg/logger"
	"github.com/engineeringinflow/inflow-backend/pkg/models"
	"github.com/engineeringinflow/inflow-backend/pkg/repo/query"
	"github.com/engineeringinflow/inflow-backend/pkg/repo/query/queryfunc"
	"github.com/jinzhu/copier"
	"github.com/rotisserie/eris"
	"google.golang.org/api/sheets/v4"
	"gorm.io/gorm/clause"
)

type ReleaseNoteRepo struct {
	db       *db.DB
	logger   *logger.Logger
	sheetAPI *sheets.Service
}

func NewReleaseNoteRepo(db *db.DB) *ReleaseNoteRepo {
	return &ReleaseNoteRepo{
		db:     db,
		logger: logger.New("repo/ReleaseNote"),
	}
}

type CreateReleaseNotesParams struct {
	models.JwtClaimsInfo
	Title       string `json:"title"`
	Description string `json:"description"`
	ReleaseDate int64  `json:"release_date"`
}

func (r *ReleaseNoteRepo) Create(params CreateReleaseNotesParams) (result *models.ReleaseNote, err error) {
	var note models.ReleaseNote
	err = copier.Copy(&note, &params)
	if err != nil {
		return nil, eris.Wrap(err, "")
	}

	if err = r.db.Model(&models.ReleaseNote{}).Clauses(clause.Returning{}).Create(&note).Error; err != nil {
		return
	}
	result = &note
	return
}

type UpdateReleaseNotesParams struct {
	models.PaginationParams
	models.JwtClaimsInfo
	ID          string `param:"id" query:"id" form:"id" validate:"required"`
	Title       string `json:"title"`
	Description string `json:"description"`
	ReleaseDate int64  `json:"release_date"`
}

func (r *ReleaseNoteRepo) Update(params UpdateReleaseNotesParams) (result *models.ReleaseNote, err error) {
	var note models.ReleaseNote
	err = copier.Copy(&note, &params)
	if err != nil {
		return nil, eris.Wrap(err, "")
	}

	if err = r.db.Model(&note).
		Clauses(clause.Returning{}).
		Where("id = ?", params.ID).
		Updates(&note).Error; err != nil {
		return
	}
	result = &note
	return
}

type PaginateReleaseNotesParams struct {
	models.PaginationParams
	models.JwtClaimsInfo
}

func (r *ReleaseNoteRepo) Paginate(params PaginateReleaseNotesParams) *query.Pagination {
	var builder = queryfunc.NewReleaseNoteBuilder(queryfunc.ReleaseNoteBuilderOptions{
		QueryBuilderOptions: queryfunc.QueryBuilderOptions{
			Role: params.GetRole(),
		},
	})
	var result = query.New(r.db, builder).
		WhereFunc(func(builder *query.Builder) {
		}).
		Page(params.Page).
		Limit(params.Limit).
		PagingFunc()
	return result
}

type DeleteReleaseNotesParams struct {
	models.PaginationParams
	models.JwtClaimsInfo
	ID string `json:"id" param:"id" query:"id" form:"id" validate:"required"`
}

func (r *ReleaseNoteRepo) Delete(params DeleteReleaseNotesParams) (err error) {
	err = r.db.Unscoped().Delete(&models.ReleaseNote{}, "id = ?", params.ID).Error
	return
}
