package models

import "github.com/engineeringinflow/inflow-backend/pkg/models/price"

type InquiryCartItems []*InquiryCartItem

type InquiryCartItem struct {
	ID        string    `gorm:"primaryKey" json:"id,omitempty"`
	CreatedAt int64     `gorm:"autoCreateTime" json:"created_at,omitempty"`
	UpdatedAt int64     `gorm:"autoUpdateTime" json:"updated_at,omitempty"`
	DeletedAt DeletedAt `sql:"index" json:"deleted_at,omitempty" swaggertype:"primitive,integer"`

	InquiryID           string `gorm:"size:200" json:"inquiry_id"`
	PurchaseOrderID     string `gorm:"size:200" json:"purchase_order_id"`
	BulkPurchaseOrderID string `gorm:"size:200" json:"bulk_purchase_order_id"`
	CheckoutSessionID   string `gorm:"size:200" json:"checkout_session_id"`

	Sku       string `json:"sku,omitempty"`
	Style     string `json:"style,omitempty"`
	Color     string `json:"color,omitempty"`
	ColorName string `json:"color_name,omitempty"`
	Size      string `json:"size,omitempty"`

	Qty        int64       `gorm:"default:1" json:"qty"`
	UnitPrice  price.Price `gorm:"type:decimal(20,4);default:0.0" json:"unit_price"`
	TotalPrice price.Price `gorm:"type:decimal(20,4);default:0.0" json:"total_price"`

	WaitingForCheckout *bool `gorm:"default:false" json:"waiting_for_checkout"`

	NoteToSupplier string `gorm:"size:2000" json:"note_to_supplier"`
}

type InquiryCartItemCreateForm struct {
	Color string `json:"color" validate:"required"`
	Size  string `json:"size" validate:"required"`
	Qty   int    `json:"qty" validate:"required"`

	UnitPrice           price.Price `gorm:"type:decimal(20,4);default:0.0" json:"unit_price" validate:"required"`
	ColorName           string      `json:"color_name"`
	NoteToSupplier      string      `json:"note_to_supplier"`
	InquiryID           string      `param:"inquiry_id" json:"inquiry_id" validate:"required"`
	BulkPurchaseOrderID string      `param:"bulk_purchase_order_id" json:"bulk_purchase_order_id"`
}

type InquiryCartItemsUpdateForm struct {
	InquiryID string           `param:"inquiry_id" json:"inquiry_id" validate:"required"`
	Items     []*OrderCartItem `json:"items"`

	JwtClaimsInfo
}

type GetInquiryCartItemsParams struct {
	InquiryID string `param:"inquiry_id" validate:"required"`

	JwtClaimsInfo
}

type InquiryUserCartUpdateItemsForm struct {
	Items []*InquiryCartItemCreateForm `json:"items"`

	JwtClaimsInfo
}
