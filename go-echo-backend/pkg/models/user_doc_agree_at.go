package models

import "github.com/engineeringinflow/inflow-backend/pkg/models/enums"

type UserDocAgreement struct {
	UserID         string           `gorm:"primaryKey" json:"user_id,omitempty"`
	SettingDocType enums.SettingDoc `gorm:"primaryKey" json:"setting_doc_type"`
	AgreeAt        int64            `json:"agree_at,omitempty"`

	CreatedAt int64      `gorm:"autoCreateTime" json:"created_at,omitempty"`
	UpdatedAt int64      `gorm:"autoUpdateTime" json:"updated_at,omitempty"`
	DeletedAt *DeletedAt `sql:"index" json:"deleted_at,omitempty" swaggertype:"primitive,integer"`
}
