package models

// Category Category's model
type BlogCategory struct {
	// Model
	ID        string    `gorm:"primaryKey" json:"id,omitempty"`
	CreatedAt int64     `gorm:"autoCreateTime" json:"created_at,omitempty"`
	UpdatedAt int64     `gorm:"autoUpdateTime" json:"updated_at,omitempty"`
	DeletedAt DeletedAt `sql:"index" json:"deleted_at,omitempty" swaggertype:"primitive,integer"`

	Name        string               `json:"name,omitempty"`
	Slug        string               `gorm:"not null;unique" json:"slug,omitempty"`
	Description string               `json:"description,omitempty"`
	Subtitle    string               `json:"sub_title,omitempty"`
	Vi          *BlogCategoryContent `json:"vi,omitempty"`

	Order     int `gorm:"default:0" json:"order,omitempty"`
	TotalPost int `gorm:"-" json:"total_post,omitempty"`
}

type BlogCategoryUpdateForm struct {
	BlogCategoryID string `param:"blog_category_id" validate:"required"`

	BlogCategoryCreateForm
}

type BlogCategoryCreateForm struct {
	Name        string               `json:"name,omitempty"`
	Subtitle    string               `json:"sub_title,omitempty"`
	Description string               `json:"description,omitempty"`
	Order       int                  `gorm:"default:0" json:"order,omitempty"`
	Vi          *BlogCategoryContent `json:"vi,omitempty"`

	JwtClaimsInfo
}
