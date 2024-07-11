package models

import "github.com/engineeringinflow/inflow-backend/pkg/models/enums"

type AsFeaturedIn struct {
	ID        string `gorm:"primaryKey" json:"id,omitempty"`
	CreatedAt int64  `gorm:"autoCreateTime" json:"created_at,omitempty"`
	UpdatedAt int64  `gorm:"autoUpdateTime" json:"updated_at,omitempty"`

	Slug  string      `json:"slug" gorm:"unique"`
	Title string      `json:"title,omitempty"`
	Image *Attachment `json:"image,omitempty"`
	Link  string      `json:"link,omitempty"`
	Logo  *Attachment `json:"logo,omitempty"`

	PublishedAt *int64           `json:"published_at,omitempty"`
	Status      enums.PostStatus `json:"status,omitempty" gorm:"default:'new'"`

	VI *AsFeaturedInContent `json:"vi,omitempty"`
}

type AsFeaturedInCreateForm struct {
	Model
	JwtClaimsInfo

	Title  string           `json:"title,omitempty" validate:"required"`
	Image  *Attachment      `json:"image,omitempty" validate:"required"`
	Link   string           `json:"link,omitempty" validate:"required,startswith=http"`
	Logo   *Attachment      `json:"logo,omitempty" validate:"required"`
	Status enums.PostStatus `json:"status,omitempty" gorm:"default:'new'"`

	VI *AsFeaturedInContent `json:"vi,omitempty"`
}

type AsFeaturedInUpdateForm struct {
	JwtClaimsInfo

	AsFeaturedInID string `json:"as_featured_in_id" param:"as_featured_in_id" validate:"required"`

	Status enums.PostStatus `json:"status,omitempty" gorm:"default:'new'"`
	Title  string           `json:"title,omitempty"`
	Image  *Attachment      `json:"image,omitempty"`
	Link   string           `json:"link,omitempty"`
	Logo   *Attachment      `json:"logo,omitempty"`

	VI *AsFeaturedInContent `json:"vi,omitempty"`
}
