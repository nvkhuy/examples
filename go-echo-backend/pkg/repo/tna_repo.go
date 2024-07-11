package repo

import (
	"errors"
	"github.com/engineeringinflow/inflow-backend/pkg/db"
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

type TNARepo struct {
	db     *db.DB
	logger *logger.Logger
}

func NewTNARepo(db *db.DB) *TNARepo {
	return &TNARepo{
		db:     db,
		logger: logger.New("repo/TNA"),
	}
}

type CreateTNAsParams struct {
	models.JwtClaimsInfo
	ReferenceID  string            `gorm:"not null" json:"reference_id"`
	OrderType    enums.InquiryType `json:"order_type" validate:"omitempty,oneof=bulk sample"`
	Title        string            `json:"title"`
	SubTitle     string            `json:"sub_title"`
	Comment      string            `json:"comment"`
	DateFrom     int64             `json:"date_from"`
	DateTo       int64             `json:"date_to"`
	AssigneeIDs  pq.StringArray    `json:"assignee_ids"`
	Dependencies pq.StringArray    `json:"dependencies"`
}

func (r *TNARepo) Create(params CreateTNAsParams) (result *models.TNA, err error) {
	var tna models.TNA
	err = copier.Copy(&tna, &params)
	if err != nil {
		return nil, eris.Wrap(err, "")
	}

	if tna.ReferenceID == "" {
		err = errors.New("required reference id")
		return
	}

	if err = r.db.Model(&models.User{}).Where("id IN ?", []string(params.AssigneeIDs)).Find(&tna.Assignees).Error; err != nil {
		return
	}

	if err = r.db.Model(&models.TNA{}).
		Clauses(clause.Returning{}).
		Create(&tna).Error; err != nil {
		return
	}

	result, err = r.Get(GetTNAParams{
		JwtClaimsInfo: params.JwtClaimsInfo,
		ID:            tna.ID,
	})
	return
}

type UpdateTNAsParams struct {
	models.PaginationParams
	models.JwtClaimsInfo
	ID           string            `json:"id" param:"id" query:"id" form:"id" validate:"required"`
	OrderType    enums.InquiryType `json:"order_type" validate:"omitempty,oneof=bulk sample"`
	ReferenceID  string            `json:"reference_id"`
	Title        string            `json:"title"`
	SubTitle     string            `json:"sub_title"`
	Comment      string            `json:"comment"`
	DateFrom     int64             `json:"date_from"`
	DateTo       int64             `json:"date_to"`
	AssigneeIDs  pq.StringArray    `json:"assignee_ids"`
	Dependencies pq.StringArray    `json:"dependencies"`
}

func (r *TNARepo) Update(params UpdateTNAsParams) (result *models.TNA, err error) {
	var tna models.TNA
	err = copier.Copy(&tna, &params)
	if err != nil {
		return nil, eris.Wrap(err, "")
	}

	if err = r.db.Model(&models.User{}).Where("id IN ?", []string(params.AssigneeIDs)).
		Find(&tna.Assignees).Error; err != nil {
		return
	}

	if err = r.db.Model(&tna).
		Clauses(clause.Returning{}).
		Where("id = ?", params.ID).
		Updates(&tna).Error; err != nil {
		return
	}

	result, err = r.Get(GetTNAParams{
		JwtClaimsInfo: params.JwtClaimsInfo,
		ID:            tna.ID,
	})
	return
}

type PaginateTNAsParams struct {
	models.PaginationParams
	models.JwtClaimsInfo
	ReferenceID string `json:"reference_id" query:"reference_id" form:"reference_id" param:"reference_id"`
}

func (r *TNARepo) Paginate(params PaginateTNAsParams) *query.Pagination {
	var builder = queryfunc.NewTNABuilder(queryfunc.TNABuilderOptions{
		QueryBuilderOptions: queryfunc.QueryBuilderOptions{
			Role: params.GetRole(),
		},
	})
	if params.Page < 0 {
		params.Page = 1
	}
	if params.Limit == 0 {
		params.Limit = 20
	}
	var result = query.New(r.db, builder).
		WhereFunc(func(builder *query.Builder) {
			if params.ReferenceID != "" {
				builder.Where("reference_id = ?", params.ReferenceID)
			}
		}).
		Page(params.Page).
		Limit(params.Limit).
		PagingFunc()

	return result
}

type DeleteTNAsParams struct {
	models.PaginationParams
	models.JwtClaimsInfo
	ID string `json:"id" param:"id" query:"id" form:"id" validate:"required"`
}

func (r *TNARepo) Delete(params DeleteTNAsParams) (err error) {
	err = r.db.Transaction(func(tx *gorm.DB) (err error) {
		if err = r.db.Model(&models.TNA{}).
			Where("? = ANY (dependencies)", params.ID).
			Update("dependencies", gorm.Expr("array_remove(dependencies, ?)", params.ID)).Error; err != nil {
			return
		}
		err = r.db.Unscoped().Delete(&models.TNA{}, "id = ?", params.ID).Error
		return
	})
	return
}

type GetTNAParams struct {
	models.JwtClaimsInfo
	ID          string `json:"id" param:"id" query:"id" form:"id" validate:"required"`
	ReferenceID string `json:"reference_id" param:"reference_id" query:"reference_id" form:"reference_id"`
}

func (r *TNARepo) Get(params GetTNAParams) (result *models.TNA, err error) {
	builder := queryfunc.NewTNABuilder(queryfunc.TNABuilderOptions{
		QueryBuilderOptions: queryfunc.QueryBuilderOptions{
			Role: params.GetRole(),
		},
	})
	var tna models.TNA
	err = query.New(r.db, builder).
		WhereFunc(func(builder *query.Builder) {
			if params.ID != "" {
				builder.Where("id = ?", params.ID)
			}
			if params.ReferenceID != "" {
				builder.Where("reference_id = ?", params.ReferenceID)
			}
		}).
		FirstFunc(&tna)

	result = &tna
	return
}
