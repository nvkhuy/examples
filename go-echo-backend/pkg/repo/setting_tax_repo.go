package repo

import (
	"github.com/biter777/countries"
	"strings"
	"time"

	"github.com/engineeringinflow/inflow-backend/pkg/db"
	"github.com/engineeringinflow/inflow-backend/pkg/models"
	"github.com/engineeringinflow/inflow-backend/pkg/models/enums"
	"github.com/engineeringinflow/inflow-backend/pkg/repo/query"
	"github.com/engineeringinflow/inflow-backend/pkg/repo/query/queryfunc"
	"github.com/jinzhu/copier"
)

type SettingTaxRepo struct {
	db *db.DB
}

func NewSettingTaxRepo(db *db.DB) *SettingTaxRepo {
	return &SettingTaxRepo{
		db: db,
	}
}

type PaginateSettingTaxParams struct {
	models.PaginationParams
	models.JwtClaimsInfo
}

func (r *SettingTaxRepo) PaginateSettingTaxes(params PaginateSettingTaxParams) (qp *query.Pagination) {
	var result = query.New(r.db, queryfunc.NewSettingTaxBuilder(queryfunc.SettingTaxBuilderOptions{
		QueryBuilderOptions: queryfunc.QueryBuilderOptions{
			Role: params.GetRole(),
		},
	})).
		Limit(params.Limit).
		Page(params.Page).
		PagingFunc()
	return result
}

type GetSettingTaxParams struct {
	models.JwtClaimsInfo

	CountryCodes []enums.CountryCode `json:"country_code" query:"country_code" param:"country_code"`
}

func (r *SettingTaxRepo) GetSettingTaxes(params GetSettingTaxParams) ([]*models.SettingTax, error) {
	var records []*models.SettingTax
	var err = query.New(r.db, queryfunc.NewSettingTaxBuilder(queryfunc.SettingTaxBuilderOptions{
		QueryBuilderOptions: queryfunc.QueryBuilderOptions{
			Role: params.GetRole(),
		},
	})).WhereFunc(func(builder *query.Builder) {
		if len(params.CountryCodes) > 0 {
			builder.Where("st.country_code IN ?", params.CountryCodes)

			builder.Where("st.date_affected > ?", time.Now().Unix())
		}
	}).
		Limit(1).
		FindFunc(&records)
	return records, err
}

func (r *SettingTaxRepo) CreateSettingTax(form models.CreateSettingTaxForm) (*models.SettingTax, error) {
	var tax models.SettingTax

	var err = copier.Copy(&tax, &form)
	if err != nil {
		return nil, err
	}

	tax.CurrencyCode = countries.ByName(strings.TrimSpace(strings.ToUpper(form.CountryCode))).Currency().Alpha()
	err = r.db.Create(&tax).Error
	if err != nil {
		return nil, err
	}

	return &tax, err
}

func (r *SettingTaxRepo) UpdateSettingTax(form models.UpdateSettingTaxForm) (*models.SettingTax, error) {
	var tax models.SettingTax
	var err = copier.Copy(&tax, &form)
	if err != nil {
		return nil, err
	}

	tax.CurrencyCode = countries.ByName(strings.TrimSpace(strings.ToUpper(form.CountryCode))).Currency().Alpha()
	err = r.db.Model(&models.SettingTax{}).Where("id = ?", form.TaxID).Updates(&tax).Error
	if err != nil {
		return nil, err
	}

	return &tax, err
}

func (r *SettingTaxRepo) DeleteSettingTax(form models.DeleteSettingTaxForm) (err error) {
	return r.db.Unscoped().Delete(&models.SettingTax{}, "id = ?", form.TaxID).Error
}

func (r *SettingTaxRepo) GetSettingTax(form models.GetSettingTaxForm) (*models.SettingTax, error) {
	var record models.SettingTax
	var err = r.db.First(&record, "id = ?", form.TaxID).Error
	return &record, err
}

func (r *SettingTaxRepo) GetAffectedSettingTax(form models.GetAffectedSettingTaxForm) (*models.SettingTax, error) {
	var record models.SettingTax
	var err error
	if form.CurrencyCode == "" {
		form.CurrencyCode = enums.Currency(countries.ByName(strings.TrimSpace(strings.ToUpper(form.CountryCode))).Currency().Alpha())
	}
	q := r.db.Where("date_affected < EXTRACT(EPOCH FROM now())").Order("date_affected DESC")
	if form.CountryCode != "" {
		form.CountryCode = strings.ToUpper(strings.TrimSpace(form.CountryCode))
		q = q.Where("country_code = ?", form.CountryCode)
	} else if form.CurrencyCode != "" {
		q = q.Where("currency_code = ?", form.CurrencyCode)
	}
	err = q.First(&record).Error
	return &record, err
}
