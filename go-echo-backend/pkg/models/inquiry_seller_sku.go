package models

import (
	"github.com/engineeringinflow/inflow-backend/pkg/models/enums"
	"github.com/engineeringinflow/inflow-backend/pkg/models/price"
)

type InquirySellerSku struct {
	ID        string `gorm:"unique" json:"id,omitempty"`
	CreatedAt int64  `gorm:"autoCreateTime" json:"created_at,omitempty"`
	UpdatedAt int64  `gorm:"autoUpdateTime" json:"updated_at,omitempty"`

	InquirySellerID string `gorm:"primaryKey" json:"inquiry_seller_id"`
	InquirySkuID    string `gorm:"primaryKey" json:"inquiry_sku_id"`

	DueDay    *int64         `json:"due_day,omitempty"`
	Price     price.Price    `gorm:"type:decimal(20,4);default:0.0"  json:"price"`
	PriceType string         `json:"price_type"`
	Currency  enums.Currency `gorm:"default:'USD'" json:"currency"`

	Status enums.InquirySkuStatus `json:"status" gorm:"default:'new'"`
}
