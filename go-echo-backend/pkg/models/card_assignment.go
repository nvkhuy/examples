package models

type CardAssignment struct {
	ID        string    `gorm:"primaryKey" json:"id,omitempty"`
	CreatedAt int64     `gorm:"autoCreateTime" json:"created_at,omitempty"`
	UpdatedAt int64     `gorm:"autoUpdateTime" json:"updated_at,omitempty"`
	DeletedAt DeletedAt `sql:"index" json:"deleted_at,omitempty" swaggertype:"primitive,integer"`

	OrderCardID string `json:"order_card_id" validate:"required"`
	UserID      string `json:"user_id" validate:"required"`

	User *User `gorm:"-" json:"user,omitempty"`
}
