package repo

import (
	"errors"

	"github.com/engineeringinflow/inflow-backend/pkg/db"
	"github.com/engineeringinflow/inflow-backend/pkg/models"
	"github.com/engineeringinflow/inflow-backend/pkg/models/enums"
	"gorm.io/gorm/clause"
)

type SettingInquiryRepo struct {
	db *db.DB
}

func NewSettingInquiryRepo(db *db.DB) *SettingInquiryRepo {
	return &SettingInquiryRepo{
		db: db,
	}
}

type SettingInquiryCreateParams struct {
	models.JwtClaimsInfo
	Type        enums.SettingInquiry `json:"type,omitempty" validate:"omitempty,oneof=rfq_edit_timeout"`
	EditTimeout int64                `json:"edit_timeout,omitempty"`
	UpdatedBy   string               `json:"updated_by,omitempty"`
}

func (r *SettingInquiryRepo) Create(params SettingInquiryCreateParams) (setting *models.SettingInquiry, err error) {
	if !params.Type.IsValid() {
		err = errors.New("invalid setting inquiry type")
		return
	}

	var result = models.SettingInquiry{
		Type:        params.Type,
		EditTimeout: params.EditTimeout,
		UpdatedBy:   params.JwtClaimsInfo.GetUserID(),
	}
	if err = r.db.Model(&result).
		Clauses(clause.Returning{}).
		Create(&result).Error; err != nil {
		return
	}
	setting = &result
	return
}

type GetSettingInquiryParams struct {
	models.JwtClaimsInfo
	Type enums.SettingInquiry `json:"type,omitempty" param:"type" query:"type" form:"type" validate:"omitempty,oneof=rfq_edit_timeout"`
}

func (r *SettingInquiryRepo) Get(params GetSettingInquiryParams) (*models.SettingInquiry, error) {
	var result models.SettingInquiry
	var err = r.db.Where("type = ?", params.Type).First(&result).Error
	if err != nil {
		if r.db.IsRecordNotFoundError(err) {
			return &result, nil
		}
		return nil, err
	}
	return &result, nil
}

type SettingInquiryUpdateParams struct {
	models.JwtClaimsInfo
	Type        enums.SettingInquiry `json:"type,omitempty" param:"type" query:"type" form:"type" validate:"omitempty,oneof=rfq_edit_timeout"`
	EditTimeout int64                `json:"edit_timeout,omitempty"`
	UpdatedBy   string               `json:"updated_by,omitempty"`
}

func (r *SettingInquiryRepo) Update(params SettingInquiryUpdateParams) (inquiry *models.SettingInquiry, err error) {
	if !params.Type.IsValid() {
		err = errors.New("invalid setting inquiry type")
		return
	}
	var result = models.SettingInquiry{
		Type:        params.Type,
		EditTimeout: params.EditTimeout,
		UpdatedBy:   params.JwtClaimsInfo.GetUserID(),
	}

	if err = r.db.Model(&result).
		Clauses(clause.Returning{}).
		Where("type = ?", params.Type).
		Updates(&result).Error; err != nil {
		return
	}
	inquiry = &result
	return
}
