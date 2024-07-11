package models

import "github.com/engineeringinflow/inflow-backend/pkg/models/price"

type BulkPurchaseOrderItems []*BulkPurchaseOrderItem

type BulkPurchaseOrderItem struct {
	ID        string    `gorm:"primaryKey" json:"id,omitempty"`
	CreatedAt int64     `gorm:"autoCreateTime" json:"created_at,omitempty"`
	UpdatedAt int64     `gorm:"autoUpdateTime" json:"updated_at,omitempty"`
	DeletedAt DeletedAt `sql:"index" json:"deleted_at,omitempty" swaggertype:"primitive,integer"`

	PurchaseOrderID     string       `json:"purchase_order_id"`
	BulkPurchaseOrderID string       `json:"bulk_purchase_order_id"`
	ColorName           string       `json:"color_name"`
	Size                string       `json:"size"`
	Style               string       `json:"style,omitempty"`
	Sku                 string       `json:"sku,omitempty"`
	Qty                 int64        `gorm:"default:1" json:"qty"`
	UnitPrice           *price.Price `gorm:"type:decimal(20,4);default:0.0" json:"unit_price,omitempty"`
	TotalPrice          *price.Price `gorm:"type:decimal(20,4);default:0.0" json:"total_price,omitempty"`
}

type BulkPurchaseOrderItemCreateForm struct {
	ColorName string `json:"color_name"`
	Size      string `json:"size"`
	Qty       int    `gorm:"default:1" json:"qty"`
}

type BulkPurchaseOrderItemsUpdateForm struct {
	InquiryID      string                             `json:"inquiry_id"`
	Items          []*BulkPurchaseOrderItemCreateForm `json:"items"`
	NoteToSupplier string                             `json:"note_to_supplier"`

	JwtClaimsInfo
}

type GetBulkPurchaseOrderItemsParams struct {
	InquiryID string `param:"inquiry_id" validate:"required"`

	JwtClaimsInfo
}

func (items BulkPurchaseOrderItems) GetSizeQty() (v SizeQty) {
	v = SizeQty{}

	for _, item := range items {
		if value, found := v[item.Size]; found {
			v[item.Size] = value + float64(item.Qty)
		} else {
			v[item.Size] = 0
		}
	}
	return
}
