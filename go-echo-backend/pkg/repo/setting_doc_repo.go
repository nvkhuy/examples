package repo

import (
	"errors"
	"github.com/engineeringinflow/inflow-backend/pkg/db"
	"github.com/engineeringinflow/inflow-backend/pkg/models"
	"github.com/engineeringinflow/inflow-backend/pkg/models/enums"
	"gorm.io/gorm/clause"
)

type SettingDocRepo struct {
	db *db.DB
}

func NewSettingDocRepo(db *db.DB) *SettingDocRepo {
	return &SettingDocRepo{
		db: db,
	}
}

type GetSettingDocParams struct {
	models.JwtClaimsInfo
	Type enums.SettingDoc `json:"type,omitempty" param:"type" query:"type" form:"type" validate:"omitempty,oneof=nda tnc"`
}

func (r *SettingDocRepo) Get(params GetSettingDocParams) (doc *models.SettingDoc, err error) {
	var result models.SettingDoc
	err = r.db.Where("type = ?", params.Type).First(&result).Error
	doc = &result
	return
}

type SettingDocCreateParams struct {
	models.JwtClaimsInfo
	Type     enums.SettingDoc     `json:"type,omitempty" validate:"omitempty,oneof=nda tnc"`
	Document *models.Attachment   `json:"document,omitempty"`
	Metadata *models.JsonMetaData `json:"metadata,omitempty"`
}

func (r *SettingDocRepo) Create(params SettingDocCreateParams) (doc *models.SettingDoc, err error) {
	if !params.Type.IsValid() {
		err = errors.New("invalid setting doc type")
		return
	}

	var result = models.SettingDoc{
		Type:      params.Type,
		Document:  params.Document,
		Metadata:  params.Metadata,
		UpdatedBy: params.JwtClaimsInfo.GetUserID(),
	}
	if err = r.db.Model(&result).
		Clauses(clause.Returning{}).
		Create(&result).Error; err != nil {
		return
	}
	doc = &result
	return
}

type SettingDocUpdateParams struct {
	models.JwtClaimsInfo
	Type     enums.SettingDoc     `json:"type,omitempty" param:"type" query:"type" form:"type" validate:"omitempty,oneof=nda tnc"`
	Document *models.Attachment   `json:"document,omitempty"`
	Metadata *models.JsonMetaData `json:"metadata,omitempty"`
}

func (r *SettingDocRepo) Update(params SettingDocUpdateParams) (doc *models.SettingDoc, err error) {
	if !params.Type.IsValid() {
		err = errors.New("invalid setting doc type")
		return
	}
	var result = models.SettingDoc{
		Type:      params.Type,
		Document:  params.Document,
		Metadata:  params.Metadata,
		UpdatedBy: params.JwtClaimsInfo.GetUserID(),
	}
	if err = r.db.Model(&result).
		Clauses(clause.Returning{}).
		Where("type = ?", params.Type).
		Updates(&result).Error; err != nil {
		return
	}
	doc = &result
	return
}
