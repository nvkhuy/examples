package models

import (
	"github.com/engineeringinflow/inflow-backend/pkg/models/enums"
	"github.com/lib/pq"
)

type Document struct {
	Model
	Title         string               `gorm:"not null" json:"title,omitempty"`
	Content       string               `gorm:"not null" json:"content,omitempty"`
	Slug          string               `gorm:"not null;unique" json:"slug"`
	Vi            *DocumentContent     `json:"vi,omitempty"`
	PublishedAt   *int64               `json:"published_at,omitempty"`
	Status        enums.DocumentStatus `gorm:"not null;default:'new'" json:"status,omitempty"`
	FeaturedImage *Attachment          `gorm:"not null" json:"featured_image,omitempty"`
	CategoryID    string               `gorm:"not null" json:"category_id,omitempty"`
	UserID        string               `gorm:"not null" json:"user_id,omitempty"`
	VisibleTo     pq.StringArray       `gorm:"type:varchar(100)[]" json:"visible_to"`

	User     *User             `gorm:"-" json:"user,omitempty"`
	Category *DocumentCategory `gorm:"-" json:"category,omitempty"`
	Tags     []*DocumentTag    `gorm:"-" json:"tags,omitempty"`
}

type CreateDocumentRequest struct {
	JwtClaimsInfo
	Title         string               `json:"title,omitempty" validate:"required"`
	Content       string               `json:"content,omitempty" validate:"required"`
	Vi            *DocumentContent     `json:"vi,omitempty"`
	PublishedAt   *int64               `json:"published_at,omitempty"`
	Status        enums.DocumentStatus `json:"status,omitempty" gorm:"default:'new'"`
	FeaturedImage *Attachment          `json:"featured_image,omitempty" validate:"required"`
	CategoryID    string               `json:"category_id,omitempty" validate:"required"`
	UserID        string               `json:"user_id,omitempty"`
	TagIDs        []string             `json:"tag_ids,omitempty"`
	VisibleTo     []enums.Role         `json:"visible_to" validate:"required"`
}

type UpdateDocumentRequest struct {
	CreateDocumentRequest
	DocumentID string `param:"document_id" validate:"required"`
}

type GetDocumentParams struct {
	JwtClaimsInfo
	DocumentID string `json:"document_id,omitempty"`
	Slug       string `param:"slug" query:"slug" form:"slug" validate:"required"`
}

type GetDocumentListParams struct {
	PaginationParams
	JwtClaimsInfo

	Statuses    []enums.DocumentStatus `json:"statuses" query:"statuses" form:"statuses"`
	CategoryIDs []string               `json:"categories" query:"categories" form:"categories"`
	TagIDs      []string               `json:"tags" query:"tags" form:"tags"`
	Roles       []string               `json:"roles" query:"roles" form:"roles"`

	Language enums.LanguageCode `json:"-"`
}

type DeleteDocumentParams struct {
	JwtClaimsInfo
	DocumentID string `param:"document_id" query:"document_id" path:"document_id" validate:"required"`
}
