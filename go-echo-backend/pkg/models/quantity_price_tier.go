package models

import (
	"github.com/engineeringinflow/inflow-backend/pkg/models/enums"
	"github.com/engineeringinflow/inflow-backend/pkg/models/price"
)

// Faq Faq's model
type QuantityPriceTier struct {
	Model

	ProductID   string            `json:"product_id,omitempty" validate:"required"`
	MinQuantity int               `json:"min_quantity,omitempty" validate:"required"`
	MaxQuantity int               `json:"max_quantity,omitempty"`
	Price       price.Price       `gorm:"type:decimal(20,4);default:0.0" json:"price,omitempty" validate:"required"`
	Unit        enums.ProductUnit `json:"unit,omitempty" validate:"required"`
}

type QuantityPriceTierUpdateForm struct {
	ProductID   string            `json:"product_id,omitempty" validate:"required"`
	MinQuantity int               `json:"min_quantity,omitempty" validate:"required"`
	MaxQuantity int               `json:"max_quantity,omitempty"`
	Price       price.Price       `json:"price,omitempty" validate:"required"`
	Unit        enums.ProductUnit `json:"unit,omitempty" validate:"required"`
}
