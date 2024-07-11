package models

// ShippingMethod ShippingMethod's model
type ShippingMethod struct {
	Model

	Name         string `json:"name,omitempty"`
	Code         string `gorm:"not null;unique" json:"code,omitempty"`
	Description  string `json:"description,omitempty"`
	ShippingType string `gorm:"default:normal" json:"shipping_type,omitempty"`
	SortOrder    int    `gorm:"default:0" json:"sort_order,omitempty"`
	Active       bool   `gorm:"default:false" json:"active,omitempty"`
}

type ShippingMethodUpdateForm struct {
	JwtClaimsInfo

	Name         string `json:"name,omitempty"`
	Code         string `json:"code,omitempty"`
	Description  string `json:"description,omitempty"`
	ShippingType string `gorm:"default:normal" json:"shipping_type,omitempty"`
	SortOrder    int    `gorm:"default:0" json:"sort_order,omitempty"`
	Active       bool   `gorm:"default:false" json:"active,omitempty"`
}
