package models

import (
	"github.com/engineeringinflow/inflow-backend/pkg/models/enums"
	"github.com/lib/pq"
)

// Shop Shop's model
type PageSection struct {
	Model

	Title         string                `json:"title,omitempty"`
	Content       string                `json:"content,omitempty"`
	Vi            *PageSectionContent   `json:"vi,omitempty"`
	SectionType   enums.PageSectionType `json:"section_type" validate:"oneof=category_main catalog_drop catalog_collection catalog_closet,omitempty"`
	PageID        string                `json:"page_id"`
	Order         int                   `gorm:"default:0" json:"order"`
	Metadata      *PageSectionMetadata  `json:"metadata,omitempty"`
	ProductIds    pq.StringArray        `gorm:"type:varchar(100)[]" json:"product_ids"`
	CategoryIds   pq.StringArray        `gorm:"type:varchar(100)[]" json:"category_ids"`
	CollectionIds pq.StringArray        `gorm:"type:varchar(100)[]" json:"collection_ids"`

	Products    []*Product    `gorm:"-" json:"products,omitempty"`
	Categories  []*Category   `gorm:"-" json:"categories,omitempty"`
	Collections []*Collection `gorm:"-" json:"collections,omitempty"`
}

type PageSectionCreateForm struct {
	Title       string                `json:"title,omitempty"`
	Content     string                `json:"content,omitempty"`
	Vi          *PageSectionContent   `json:"vi,omitempty"`
	SectionType enums.PageSectionType `json:"section_type" validate:"oneof=category_main catalog_drop catalog_collection catalog_closet,omitempty"`
	PageID      string                `json:"page_id"`
	Order       int                   `json:"order"`
	Metadata    Attachments           `json:"metadata"`
	ForRole     enums.Role            `json:"-"`
	CatalogIDs  pq.StringArray        `gorm:"type:varchar(100)[]" json:"catalogs_ids"`
}

type PageSectionUpdateForm struct {
	JwtClaimsInfo

	ID            string                `json:"id,omitempty"`
	Title         string                `json:"title,omitempty"`
	Vi            *PageSectionContent   `json:"vi,omitempty"`
	Order         int                   `json:"order"`
	SectionType   enums.PageSectionType `json:"section_type" validate:"oneof=category_main catalog_drop catalog_collection catalog_closet,omitempty"`
	Content       string                `json:"content,omitempty"`
	Metadata      *PageSectionMetadata  `json:"metadata,omitempty"`
	ProductIds    pq.StringArray        `gorm:"type:varchar(100)[]" json:"product_ids"`
	CategoryIds   pq.StringArray        `gorm:"type:varchar(100)[]" json:"category_ids"`
	CollectionIds pq.StringArray        `gorm:"type:varchar(100)[]" json:"collection_ids"`
}
