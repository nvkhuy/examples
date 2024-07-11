package models

type ProductClass struct {
	ProductID string    `gorm:"primaryKey" json:"product_id,omitempty"`
	Class     string    `gorm:"primaryKey" json:"class,omitempty"`
	Conf      float64   `json:"conf,omitempty"`
	CreatedAt int64     `gorm:"autoCreateTime" json:"created_at,omitempty"`
	UpdatedAt int64     `gorm:"autoUpdateTime" json:"updated_at,omitempty"`
	DeletedAt DeletedAt `sql:"index" json:"deleted_at,omitempty" swaggertype:"primitive,integer"`
}
