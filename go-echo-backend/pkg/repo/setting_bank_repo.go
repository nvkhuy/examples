package repo

import (
	"strings"

	"github.com/engineeringinflow/inflow-backend/pkg/db"
	"github.com/engineeringinflow/inflow-backend/pkg/models"
	"github.com/engineeringinflow/inflow-backend/pkg/repo/query"
	"github.com/engineeringinflow/inflow-backend/pkg/repo/query/queryfunc"
	"github.com/jinzhu/copier"
	"gorm.io/gorm/clause"
)

type SettingBankRepo struct {
	db *db.DB
}

func NewSettingBankRepo(db *db.DB) *SettingBankRepo {
	return &SettingBankRepo{
		db: db,
	}
}

type GetSettingBankInfosParams struct {
	models.JwtClaimsInfo
	IsDisabled *bool `json:"-"`
}

func (s *SettingBankRepo) GetSettingBankInfos(params GetSettingBankInfosParams) ([]*models.SettingBank, error) {
	var records []*models.SettingBank
	var err = query.New(s.db, queryfunc.NewSettingBankBuilder(queryfunc.SettingBankBuilderOptions{
		QueryBuilderOptions: queryfunc.QueryBuilderOptions{
			Role: params.GetRole(),
		},
	})).WhereFunc(func(builder *query.Builder) {
		if params.IsDisabled != nil {
			builder.Where("st.is_disabled = ?", *params.IsDisabled)
		}
	}).
		FindFunc(&records)
	return records, err
}

type PaginateSettingBankParams struct {
	models.PaginationParams
	models.JwtClaimsInfo
}

func (s *SettingBankRepo) PaginateSettingBank(params PaginateSettingBankParams) (qp *query.Pagination) {
	var result = query.New(s.db, queryfunc.NewSettingBankBuilder(queryfunc.SettingBankBuilderOptions{
		QueryBuilderOptions: queryfunc.QueryBuilderOptions{
			Role: params.GetRole(),
		},
	})).
		Limit(params.Limit).
		Page(params.Page).
		PagingFunc()
	return result
}
func (s *SettingBankRepo) CreateSettingBanks(form models.SettingBanksForm) (bank models.SettingBank, err error) {
	form.CountryCode = strings.TrimSpace(strings.ToUpper(form.CountryCode))
	err = copier.Copy(&bank, &form)
	if err != nil {
		return models.SettingBank{}, err
	}
	err = s.db.Create(&bank).Error
	return
}
func (s *SettingBankRepo) UpdateSettingBanks(form models.SettingBanksForm) (bank models.SettingBank, err error) {
	form.CountryCode = strings.TrimSpace(strings.ToUpper(form.CountryCode))
	err = copier.Copy(&bank, &form)
	if err != nil {
		return models.SettingBank{}, err
	}
	err = s.db.Updates(&bank).Error
	return
}
func (s *SettingBankRepo) DeleteSettingBanks(form models.DeleteSettingBanksForm) (err error) {
	err = s.db.Unscoped().Delete(&models.SettingBank{}, "id = ?", form.ID).Error
	return
}
func (s *SettingBankRepo) DeleteSettingBanksByCountryCode(form models.DeleteSettingBanksByCountryCodeForm) (banks models.SettingBankSlice, err error) {
	deleteCountry := strings.TrimSpace(strings.ToUpper(form.CountryCode))
	err = s.db.Unscoped().Clauses(clause.Returning{}).Where("country_code = ?", deleteCountry).Delete(&banks).Error
	return
}
