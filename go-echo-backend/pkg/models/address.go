package models

import (
	"github.com/engineeringinflow/inflow-backend/pkg/location"
	"github.com/engineeringinflow/inflow-backend/pkg/models/enums"
)

// Address Address's model
type Address struct {
	Model

	UserID string `gorm:"not null" json:"user_id"`

	Name        string `json:"name,omitempty"`
	PhoneNumber string `json:"phone_number,omitempty" validate:"isPhone"`
	Email       string `gorm:"type:citext" json:"email,omitempty"`

	AddressType enums.AddressType `gorm:"default:'primary'" json:"address_type,omitempty"`

	CoordinateID string      `json:"coordinate_id,omitempty"`
	Coordinate   *Coordinate `json:"coordinate,omitempty"`
}

type AddressForm struct {
	Name        string              `json:"name,omitempty"`
	PhoneNumber string              `json:"phone_number,omitempty"`
	Email       string              `gorm:"type:citext" json:"email,omitempty"`
	Coordinate  location.Coordinate `json:"coordinate,omitempty"`
}
