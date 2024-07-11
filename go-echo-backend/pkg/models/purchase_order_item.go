package models

import (
	"github.com/engineeringinflow/inflow-backend/pkg/models/price"
)

type PurchaseOrderItem struct {
	Model

	PurchaseOrderID   string       `json:"purchase_order_id,omitempty"`
	CheckoutSessionID string       `json:"checkout_session_id,omitempty"`
	Style             string       `json:"style,omitempty"`
	Sku               string       `json:"sku,omitempty"`
	Size              string       `json:"size,omitempty"`
	NumSamples        int64        `json:"num_samples,omitempty"`
	UnitPrice         *price.Price `json:"unit_price,omitempty"`
	TotalPrice        *price.Price `json:"total_price,omitempty"`

	ProductID string   `gorm:"size:100" json:"product_id,omitempty"`
	Product   *Product `gorm:"-" json:"product,omitempty"`

	UserID string `gorm:"size:100" json:"user_id,omitempty"`

	VariantID string   `gorm:"size:200" json:"variant_id,omitempty"`
	Variant   *Variant `gorm:"-" json:"variant,omitempty"`

	Color string `gorm:"size:200" json:"color,omitempty"`
	Notes string `gorm:"size:2000" json:"notes,omitempty"`

	FabricID   string  `gorm:"size:200" json:"fabric_id,omitempty"`
	FabricName string  `gorm:"size:200" json:"fabric_name,omitempty"`
	Fabric     *Fabric `gorm:"-" json:"fabric,omitempty"`

	Attachments Attachments `json:"attachments,omitempty"`

	Quantity int64 `json:"quantity,omitempty" validate:"required"`
}
