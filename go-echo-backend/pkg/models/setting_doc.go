package models

import "github.com/engineeringinflow/inflow-backend/pkg/models/enums"

type SettingDoc struct {
	Type      enums.SettingDoc `gorm:"primaryKey" json:"type,omitempty"`
	Document  *Attachment      `json:"document,omitempty"`
	Metadata  *JsonMetaData    `json:"metadata,omitempty"`
	UpdatedBy string           `json:"updated_by,omitempty"`

	CreatedAt int64      `gorm:"autoCreateTime" json:"created_at,omitempty"`
	UpdatedAt int64      `gorm:"autoUpdateTime" json:"updated_at,omitempty"`
	DeletedAt *DeletedAt `sql:"index" json:"deleted_at,omitempty" swaggertype:"primitive,integer"`
}
