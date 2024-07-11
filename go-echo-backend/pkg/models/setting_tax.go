package models

import "github.com/engineeringinflow/inflow-backend/pkg/models/enums"

type SettingTax struct {
	Model

	CountryCode   string  `gorm:"column:country_code;type:varchar(20);not null" json:"country_code,omitempty"`
	CurrencyCode  string  `gorm:"type:varchar(20)" json:"currency_code,omitempty"`
	TaxPercentage float64 `gorm:"type:decimal(20,4);default:0.0" json:"tax_percentage,omitempty"`
	DateAffected  int64   `json:"date_affected,omitempty"`

	Taxes []*SettingTax `gorm:"-" json:"taxes,omitempty"`
}

type CreateSettingTaxForm struct {
	JwtClaimsInfo

	CountryCode   string  `json:"country_code,omitempty" validate:"required"`
	CurrencyCode  string  `json:"currency_code,omitempty"`
	TaxPercentage float64 `json:"tax_percentage,omitempty" validate:"required,min=0,max=100"`
	DateAffected  int64   `json:"date_affected,omitempty" validate:"required"`
}

type UpdateSettingTaxForm struct {
	JwtClaimsInfo

	TaxID string `param:"tax_id" json:"tax_id" validate:"required"`

	CurrencyCode  string  `json:"currency_code,omitempty"`
	CountryCode   string  `json:"country_code,omitempty" validate:"required"`
	TaxPercentage float64 `json:"tax_percentage,omitempty" validate:"required,min=0,max=100"`
	DateAffected  int64   `json:"date_affected,omitempty" validate:"required"`
}

type DeleteSettingTaxForm struct {
	JwtClaimsInfo
	TaxID string `param:"tax_id" query:"tax_id" json:"tax_id"  validate:"required"`
}

type GetSettingTaxForm struct {
	JwtClaimsInfo
	TaxID string `param:"tax_id" query:"tax_id" json:"tax_id"  validate:"required"`
}

type GetAffectedSettingTaxForm struct {
	CountryCode  string         `json:"country_code"`
	CurrencyCode enums.Currency `json:"currency_code"`
}
