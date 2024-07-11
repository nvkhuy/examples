package models

import "github.com/lib/pq"

// Faq Faq's model
type CollectionProductGroup struct {
	Model

	Name          string      `json:"name,omitempty"`
	Description   string      `json:"description,omitempty"`
	FeaturedImage *Attachment `json:"featured_image,omitempty"`
	CollectionID  string      `json:"collection_id,omitempty"`

	ProductIDs pq.StringArray `gorm:"type:varchar(200)[]" json:"product_ids" swaggertype:"array,string"`
	Products   []*Product     `gorm:"-" json:"products,omitempty"`
}

type CollectionProductGroupUpdateForm struct {
	ID            string         `json:"id"`
	Name          string         `json:"name,omitempty"`
	Description   string         `json:"description,omitempty"`
	FeaturedImage *Attachment    `json:"featured_image,omitempty"`
	ProductIDs    pq.StringArray `gorm:"type:varchar(200)[]" json:"product_ids" swaggertype:"array,string"`
}
