package models

type SettingSize struct {
	CreatedAt int64 `gorm:"autoCreateTime" json:"created_at,omitempty"`
	UpdatedAt int64 `gorm:"autoUpdateTime" json:"updated_at,omitempty"`

	ID string `gorm:"unique" json:"id,omitempty"`

	Type string `gorm:"primaryKey" json:"type,omitempty"`
	Name string `gorm:"primaryKey" json:"name,omitempty"`

	Sizes []*SettingSize `gorm:"-" json:"sizes,omitempty"`
}

type SettingSizeCreateForm struct {
	JwtClaimsInfo

	Type string `json:"type" validate:"required"`

	SizeNames []string `json:"size_names" validate:"required"`
}

type SettingSizeIDForm struct {
	JwtClaimsInfo

	SizeID string `param:"size_id" validate:"required"`
}

type SettingSizeDeleteForm struct {
	SettingSizeIDForm
}

type SettingSizeUpdateForm struct {
	JwtClaimsInfo

	SizeID string `param:"size_id" validate:"required"`

	Name string `json:"name" validate:"required"`
}

type SettingSizesUpdateForm struct {
	JwtClaimsInfo

	Type string `json:"type" validate:"required"`

	SizeNames []string `json:"size_names" validate:"required"`
}

type SettingSizeDeleteTypeForm struct {
	JwtClaimsInfo

	Type string `param:"type" validate:"required"`
}

type SettingSizeUpdateTypeForm struct {
	JwtClaimsInfo

	Type string `param:"type" validate:"required"`

	NewType string `json:"new_type,omitempty" validate:"required"`
}

type GetSettingSizeTypeForm struct {
	Type string `json:"type" param:"type" query:"type" validate:"required"`
	JwtClaimsInfo
}
