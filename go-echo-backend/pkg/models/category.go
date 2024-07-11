package models

import (
	"github.com/engineeringinflow/inflow-backend/pkg/models/enums"
	"github.com/lib/pq"
)

// Category model
type Category struct {
	// Model
	ID        string    `gorm:"primaryKey" json:"id,omitempty"`
	CreatedAt int64     `gorm:"autoCreateTime" json:"created_at,omitempty"`
	UpdatedAt int64     `gorm:"autoUpdateTime" json:"updated_at,omitempty"`
	DeletedAt DeletedAt `sql:"index" json:"deleted_at,omitempty" swaggertype:"primitive,integer"`

	Name string           `json:"name,omitempty"`
	Slug string           `gorm:"not null;unique" json:"slug,omitempty"`
	Vi   *CategoryContent `json:"vi,omitempty"`

	ParentCategoryID *string `json:"parent_category_id,omitempty"`

	CategoryType enums.CategoryType `gorm:"default:'product'" json:"category_type,omitempty"`

	Order int `gorm:"default:0" json:"order,omitempty"`

	// Photo
	Icon *Attachment `json:"icon,omitempty"`

	Children      []*Category    `gorm:"foreignKey:parent_category_id;references:id" json:"children,omitempty"`
	TopProductIDs pq.StringArray `gorm:"type:varchar(200)[]" json:"top_product_ids" swaggertype:"array,string"`
	TopProducts   []*Product     `gorm:"-" json:"top_products,omitempty"`
}

type CategorySlice []*Category

type CategoryCreateForm struct {
	JwtClaimsInfo

	Name             string             `json:"name,omitempty" validate:"required"`
	ParentCategoryID string             `json:"parent_category_id,omitempty"`
	CategoryType     enums.CategoryType `json:"category_type"`
	Slug             string             `json:"slug"`
	Icon             *Attachment        `json:"icon,omitempty"`
	Order            int                `json:"order,omitempty"`
	TopProductIDs    pq.StringArray     `json:"top_product_ids,omitempty" swaggertype:"array,string"`
	Vi               *CategoryContent   `json:"vi,omitempty"`
}

type CategoryUpdateForm struct {
	JwtClaimsInfo

	CategoryID string `param:"category_id" validate:"required"`

	Name             string             `json:"name,omitempty"`
	ParentCategoryID *string            `json:"parent_category_id,omitempty"`
	CategoryType     enums.CategoryType `gorm:"default:'product'" json:"category_type"`
	Slug             string             `json:"slug"`
	Icon             *Attachment        `json:"icon,omitempty"`
	Order            int                `gorm:"default:0" json:"order,omitempty"`
	TopProductIDs    pq.StringArray     `json:"top_product_ids,omitempty" swaggertype:"array,string"`
	Vi               *CategoryContent   `json:"vi,omitempty"`
}

type CategoryResponse struct {
	ID               string              `json:"id"`
	Slug             string              `json:"slug"`
	Name             string              `json:"name,omitempty"`
	ParentCategoryID string              `json:"parent_category_id"`
	Children         []*CategoryResponse `json:"children"`
	Icon             *Attachment         `json:"icon,omitempty"`
	TopProductIds    []string            `json:"top_product_ids,omitempty"`
	TopProducts      []*Product          `json:"top_products,omitempty"`
	Vi               *CategoryContent    `json:"vi,omitempty"`
	TotalProduct     int                 `json:"total_product,omitempty"`
}

type CategoryTreeResponse struct {
	Records []*CategoryResponse `json:"records"`
}
