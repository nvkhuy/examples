package repo

import (
	"github.com/engineeringinflow/inflow-backend/pkg/db"
	"github.com/engineeringinflow/inflow-backend/pkg/models"
	"github.com/engineeringinflow/inflow-backend/pkg/models/enums"
	"gorm.io/gorm/clause"
	"time"
)

type UserDocAgreementRepo struct {
	db *db.DB
}

func NewUserDocAgreementRepo(db *db.DB) *UserDocAgreementRepo {
	return &UserDocAgreementRepo{
		db: db,
	}
}

type GetUserDocAgreementParams struct {
	models.JwtClaimsInfo

	UserID         string           `json:"user_id" query:"user_id" param:"user_id"`
	SettingDocType enums.SettingDoc `json:"setting_doc_type" query:"setting_doc_type" param:"setting_doc_type" validate:"omitempty,oneof=nda tnc"`
}

func (r *UserDocAgreementRepo) Get(params GetUserDocAgreementParams) (res *models.UserDocAgreement, err error) {
	var result models.UserDocAgreement
	r.db.Model(&models.UserDocAgreement{}).
		Where("user_id = ?", params.GetUserID()).
		Where("setting_doc_type = ?", params.SettingDocType).First(&result)
	res = &result
	return
}

type CreateUserDocAgreementParams struct {
	models.JwtClaimsInfo

	UserID         string           `json:"user_id" query:"user_id" param:"user_id"`
	SettingDocType enums.SettingDoc `json:"setting_doc_type" query:"setting_doc_type" param:"setting_doc_type" validate:"omitempty,oneof=nda tnc"`
}

func (r *UserDocAgreementRepo) Create(params CreateUserDocAgreementParams) (res *models.UserDocAgreement, err error) {
	var agreement = models.UserDocAgreement{
		UserID:         params.GetUserID(),
		SettingDocType: params.SettingDocType,
		AgreeAt:        time.Now().Unix(),
	}
	if err = r.db.Model(&agreement).Clauses(clause.Returning{}).Create(&agreement).Error; err != nil {
		return
	}
	res = &agreement
	return
}

type DeleteUserDocAgreementTypeParams struct {
	models.JwtClaimsInfo

	SettingDocType enums.SettingDoc `json:"setting_doc_type" query:"setting_doc_type" param:"setting_doc_type" validate:"omitempty,oneof=nda tnc"`
}

func (r *UserDocAgreementRepo) Delete(params DeleteUserDocAgreementTypeParams) (err error) {
	err = r.db.Unscoped().Delete(&models.UserDocAgreement{}, "setting_doc_type = ?", params.SettingDocType).Error
	return
}
