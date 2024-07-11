package models

import (
	"github.com/engineeringinflow/inflow-backend/pkg/models/enums"
	"github.com/engineeringinflow/inflow-backend/pkg/models/price"
)

// "github.com/engineeringinflow/inflow-backend/pkg/models/enums"

// Variant's model
type Variant struct {
	Model

	Title       string      `gorm:"column:title" json:"title,omitempty"`
	ProductName string      `json:"product_name,omitempty"`
	ProductID   string      `json:"product_id,omitempty"`
	Sku         string      `json:"sku,omitempty"`
	Price       price.Price `gorm:"type:decimal(20,4);default:0.0" json:"price"`
	Stock       int         `json:"stock,omitempty"`
	MinOrder    int         `json:"min_order,omitempty"`

	Color    string `gorm:"column:color" json:"color,omitempty"`
	Size     string `gorm:"column:size" json:"size,omitempty"`
	Material string `gorm:"column:material" json:"material,omitempty"`
	IsShow   *bool  `gorm:"type:bool;default:false" json:"is_show,omitempty"`

	SourceProductID       string       `json:"source_product_id,omitempty"`
	SourceVariantID       string       `json:"source_variant_id,omitempty"`
	SourceInventoryItemID string       `json:"source_inventory_item_id,omitempty"`
	SourceLocationID      string       `json:"source_location_id,omitempty"`
	Images                *Attachments `json:"images,omitempty"`

	Product *Product `gorm:"-" json:"product,omitempty"`
}

type VariantAttributeItemForm struct {
	Title              string `json:"title,omitempty"`
	ProductAttributeID string `json:"product_attribute_id,omitempty"`
	Value              string `json:"value,omitempty"`
}

type VariantCreateFrom struct {
	Title             string                      `json:"title,omitempty"`
	ProductName       string                      `json:"product_name,omitempty"`
	Price             int                         `json:"price"`
	Stock             int                         `json:"stock"`
	VariantAttributes []*VariantAttributeItemForm `json:"variant_attributes,omitempty"`
	ForRole           enums.Role                  `json:"-"`
}

type VariantAttributeUpdateForm struct {
	Model

	Title       string      `json:"title,omitempty"`
	ProductName string      `json:"product_name,omitempty"`
	ProductID   string      `json:"product_id,omitempty"`
	Sku         string      `json:"sku,omitempty"`
	Price       price.Price `gorm:"type:decimal(20,4);default:0.0" json:"price"`
	Stock       int         `json:"stock,omitempty"`
	MinOrder    int         `json:"min_order,omitempty"`

	Color    string `gorm:"column:color" json:"color,omitempty"`
	Size     string `gorm:"column:size" json:"size,omitempty"`
	Material string `gorm:"column:material" json:"material,omitempty"`
	IsShow   *bool  `gorm:"type:bool;default:false" json:"is_show,omitempty"`

	SourceProductID       string       `json:"source_product_id,omitempty"`
	SourceVariantID       string       `json:"source_variant_id,omitempty"`
	SourceInventoryItemID string       `json:"source_inventory_item_id,omitempty"`
	SourceLocationID      string       `json:"source_location_id,omitempty"`
	Images                *Attachments `json:"images,omitempty"`
}
