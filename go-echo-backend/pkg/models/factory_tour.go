package models

import "github.com/engineeringinflow/inflow-backend/pkg/models/enums"

// Faq Faq's model
type FactoryTour struct {
	Model

	Email   string                  `json:"email,omitempty" validate:"required"`
	Phone   string                  `json:"phone,omitempty"`
	Name    string                  `json:"name,omitempty"`
	Company string                  `json:"company,omitempty"`
	Reason  string                  `json:"reason,omitempty"`
	Status  enums.FactoryTourStatus `gorm:"default:'active'" validate:"oneof=active inactive,omitempty" json:"status,omitempty"`
}

// Faq Faq's model
type FactoryTourUpdateForm struct {
	Email   string     `json:"email,omitempty" validate:"required"`
	Phone   string     `json:"phone,omitempty"`
	Name    string     `json:"name,omitempty"`
	Company string     `json:"company,omitempty"`
	Reason  string     `json:"reason,omitempty"`
	Status  string     `json:"status,omitempty"`
	ForRole enums.Role `json:"-"`
}
