package models

type CartItem struct {
	ID        string    `gorm:"unique" json:"id,omitempty"`
	CreatedAt int64     `gorm:"autoCreateTime" json:"created_at,omitempty"`
	UpdatedAt int64     `gorm:"autoUpdateTime" json:"updated_at,omitempty"`
	DeletedAt DeletedAt `sql:"index" json:"deleted_at,omitempty" swaggertype:"primitive,integer"`

	CartID string `gorm:"primaryKey" json:"cart_id"`

	ProductID string   `gorm:"primaryKey" json:"product_id"`
	Product   *Product `gorm:"-" json:"product"`

	VariantID string   `gorm:"primaryKey" json:"variant_id"`
	Variant   *Variant `gorm:"-" json:"variant"`

	Quantity int `gorm:"default:1" json:"quantity"`
}

type CartItemCreateForm struct {
	ProductID string `gorm:"primaryKey" json:"product_id" validate:"required"`
	VariantID string `gorm:"primaryKey" json:"variant_id" validate:"required"`
	Quantity  int    `gorm:"default:1" json:"quantity" validate:"required"`
}

type CartItemUpdateForm struct {
	Quantity int `gorm:"default:1" json:"quantity" validate:"required"`
}
