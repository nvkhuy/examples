package models

import "github.com/engineeringinflow/inflow-backend/pkg/models/enums"

type InquiryCollection struct {
	Model

	Name   string `gorm:"index:idx_name_user_id,unique" json:"name,omitempty"`
	UserID string `gorm:"index:idx_name_user_id,unique" json:"user_id,omitempty"`
}

type InquiryCollectionUpdateForm struct {
	Name    string     `json:"name,omitempty" validate:"required"`
	UserID  string     `json:"user_id,omitempty" validate:"required"`
	ForRole enums.Role `json:"-"`
}
