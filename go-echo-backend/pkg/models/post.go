package models

import (
	"github.com/engineeringinflow/inflow-backend/pkg/models/enums"
)

// Post Post's model
type Post struct {
	Model

	Title            string      `json:"title,omitempty"`
	Content          string      `json:"content,omitempty"`
	ContentURL       string      `json:"content_url,omitempty"`
	ShortDescription string      `json:"short_description,omitempty"`
	Slug             string      `gorm:"not null;unique" json:"slug"`
	SettingSEO       *SettingSEO `gorm:"-" json:"setting_seo,omitempty"`
	SettingSeoID     string      `json:"setting_seo_id"`

	VI *PostContent `json:"vi,omitempty"`

	PublishedAt   *int64           `json:"published_at,omitempty"`
	Status        enums.PostStatus `json:"status,omitempty" gorm:"default:'new'"`
	FeaturedImage *Attachment      `json:"featured_image,omitempty"`
	CategoryID    string           `json:"category_id,omitempty"`
	UserID        string           `json:"user_id,omitempty"`

	User     *User         `gorm:"-" json:"user,omitempty"`
	Category *BlogCategory `gorm:"-" json:"category,omitempty"`
}

type PostSlice []*Post

type PostUpdateForm struct {
	PostID string `param:"post_id" validate:"required"`

	PostCreateForm
}

type PostCreateForm struct {
	JwtClaimsInfo
	Title            string           `json:"title,omitempty"`
	Content          string           `json:"content,omitempty"`
	ContentURL       string           `json:"content_url,omitempty"`
	VI               *PostContent     `json:"vi,omitempty"`
	ShortDescription string           `json:"short_description,omitempty"`
	PublishedAt      *int64           `json:"published_at,omitempty"`
	Status           enums.PostStatus `json:"status,omitempty" gorm:"default:'new'"`
	FeaturedImage    *Attachment      `json:"featured_image,omitempty"`
	CategoryID       string           `json:"category_id,omitempty"`
	UserID           string           `json:"user_id,omitempty"`
	SettingSEO       *SettingSEO      `json:"setting_seo,omitempty"`
}

type PostStats struct {
	Total int64 `json:"total,omitempty"`
}
