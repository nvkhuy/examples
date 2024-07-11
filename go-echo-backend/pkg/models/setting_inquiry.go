package models

import "github.com/engineeringinflow/inflow-backend/pkg/models/enums"

type SettingInquiry struct {
	Type        enums.SettingInquiry `gorm:"primaryKey" json:"type,omitempty"`
	EditTimeout int64                `json:"edit_timeout,omitempty"`
	UpdatedBy   string               `json:"updated_by,omitempty"`

	CreatedAt int64      `gorm:"autoCreateTime" json:"created_at,omitempty"`
	UpdatedAt int64      `gorm:"autoUpdateTime" json:"updated_at,omitempty"`
	DeletedAt *DeletedAt `sql:"index" json:"deleted_at,omitempty" swaggertype:"primitive,integer"`
}
