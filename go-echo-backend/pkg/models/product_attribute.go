package models

import (
	"github.com/engineeringinflow/inflow-backend/pkg/models/enums"
)

// ProductAttribute's model
type ProductAttribute struct {
	Model

	Name      string                 `json:"name,omitempty"`              // 'Size'
	ProductID string                 `json:"product_id,omitempty"`        // 'VariantID
	Values    ProductValidAttributes `json:"values" validate:"omitempty"` // [red, green]
	Order     int                    `json:"order"`
}

// Product Product's model
type ProductAttributeUpdateForm struct {
	Name    string                 `json:"name,omitempty"`
	Values  ProductValidAttributes `json:"values"`
	ForRole enums.Role             `json:"-"`
}
