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

type UserBankRepo struct {
	db *db.DB
}

func NewUserBankRepo(db *db.DB) *UserBankRepo {
	return &UserBankRepo{
		db: db,
	}
}

type GetUserBankInfosParams struct {
	models.JwtClaimsInfo

	UserID     string `json:"user_id" query:"user_id" param:"user_id"`
	IsDisabled *bool  `json:"-"`
}

func (s *UserBankRepo) GetUserBankInfos(params GetUserBankInfosParams) ([]*models.UserBank, error) {
	var records []*models.UserBank
	var err = query.New(s.db, queryfunc.NewUserBankBuilder(queryfunc.UserBankBuilderOptions{
		QueryBuilderOptions: queryfunc.QueryBuilderOptions{
			Role: params.GetRole(),
		},
	})).WhereFunc(func(builder *query.Builder) {
		if params.IsDisabled != nil {
			builder.Where("ub.is_disabled = ?", *params.IsDisabled)
		}

		builder.Where("ub.user_id = ?", params.UserID)
	}).
		FindFunc(&records)
	return records, err
}

type PaginateUserBankParams struct {
	models.PaginationParams
	models.JwtClaimsInfo

	UserID string `json:"user_id" query:"user_id" param:"user_id"`
}

func (s *UserBankRepo) PaginateUserBank(params PaginateUserBankParams) (qp *query.Pagination) {
	var result = query.New(s.db, queryfunc.NewUserBankBuilder(queryfunc.UserBankBuilderOptions{
		QueryBuilderOptions: queryfunc.QueryBuilderOptions{
			Role: params.GetRole(),
		},
	})).
		WhereFunc(func(builder *query.Builder) {
			builder.Where("ub.user_id = ?", params.UserID)
		}).
		Limit(params.Limit).
		Page(params.Page).
		PagingFunc()
	return result
}

func (s *UserBankRepo) CreateUserBanks(form models.UserBanksForm) (bank models.UserBank, err error) {
	form.CountryCode = strings.TrimSpace(strings.ToUpper(form.CountryCode))
	err = copier.Copy(&bank, &form)
	if err != nil {
		return models.UserBank{}, err
	}

	bank.UserID = form.GetUserID()
	err = s.db.Create(&bank).Error
	return
}
func (s *UserBankRepo) UpdateUserBanks(form models.UserBanksForm) (bank models.UserBank, err error) {
	form.CountryCode = strings.TrimSpace(strings.ToUpper(form.CountryCode))
	err = copier.Copy(&bank, &form)
	if err != nil {
		return models.UserBank{}, err
	}

	err = s.db.Updates(&bank).Error
	return
}
func (s *UserBankRepo) DeleteUserBanks(form models.DeleteUserBanksForm) (err error) {
	err = s.db.Unscoped().Delete(&models.UserBank{}, "id = ? AND user_id = ?", form.ID, form.GetUserID()).Error
	return
}
func (s *UserBankRepo) DeleteUserBanksByCountryCode(form models.DeleteUserBanksByCountryCodeForm) (banks models.UserBankSlice, err error) {
	deleteCountry := strings.TrimSpace(strings.ToUpper(form.CountryCode))
	err = s.db.Unscoped().Clauses(clause.Returning{}).Where("country_code = ? AND user_id = ?", deleteCountry, form.GetUserID()).Delete(&banks).Error
	return
}
