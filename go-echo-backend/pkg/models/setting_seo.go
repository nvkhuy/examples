package models

import "gorm.io/datatypes"

type SettingSEOLanguageGroup struct {
	Route string     `json:"route"`
	EN    SettingSEO `json:"en"`
	VI    SettingSEO `json:"vi"`
}

type SettingSEO struct {
	Model

	Route        string             `gorm:"index:idx_setting_seo_route_language_code,unique" json:"route,omitempty"`
	LanguageCode string             `gorm:"index:idx_setting_seo_route_language_code,unique" json:"language_code,omitempty"`
	Title        string             `json:"title,omitempty"`
	Description  string             `json:"description,omitempty"`
	Keywords     string             `json:"keywords,omitempty"`
	Thumbnail    *Attachment        `json:"thumbnail,omitempty"`
	Metadata     *datatypes.JSONMap `json:"metadata,omitempty"`
}

type SettingSEOSlice []SettingSEO

type CreateSettingSEOForm struct {
	JwtClaimsInfo
	Route        string             `json:"route,omitempty" validate:"required"`
	LanguageCode string             `json:"language_code" validate:"omitempty,required,oneof=en vi"`
	Title        string             `json:"title,omitempty"`
	Description  string             `json:"description,omitempty"`
	Keywords     string             `json:"keywords,omitempty"`
	Metadata     *datatypes.JSONMap `json:"metadata,omitempty"`
	Thumbnail    *Attachment        `json:"thumbnail"`
}

type UpdateSettingSEOForm struct {
	JwtClaimsInfo
	ID string `json:"id" param:"id" query:"id" validate:"required"`

	Route        string             `json:"route,omitempty"`
	LanguageCode string             `json:"language_code" validate:"omitempty,oneof=en vi"`
	Title        string             `json:"title,omitempty"`
	Description  string             `json:"description,omitempty"`
	Keywords     string             `json:"keywords,omitempty"`
	Metadata     *datatypes.JSONMap `json:"metadata,omitempty"`
	Thumbnail    *Attachment        `json:"thumbnail"`
}

type DeleteSettingSEOForm struct {
	JwtClaimsInfo
	ID string `json:"id" param:"id" query:"id" validate:"required"`
}
