package models

import "github.com/lib/pq"

type Collection struct {
	Model

	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`

	VI *CollectionContent `json:"vi,omitempty"`

	Order       int          `gorm:"default:0" json:"order,omitempty"`
	Attachments *Attachments `json:"attachments,omitempty"`

	ProductGroups []*CollectionProductGroup `json:"product_groups,omitempty"`
	Products      []*Product                `gorm:"-" json:"products,omitempty"`
	ProductIDs    pq.StringArray            `gorm:"type:varchar(200)[]" json:"product_ids,omitempty" swaggertype:"array,string"`
}

type CollectionSlice []*Collection

type CollectionUpdateForm struct {
	JwtClaimsInfo

	CollectionID string `param:"collection_id" validate:"required"`

	Name        string             `json:"name,omitempty"`
	Description string             `json:"description,omitempty"`
	VI          *CollectionContent `json:"vi,omitempty"`

	Attachments *Attachments `json:"attachments,omitempty"`

	ProductIds []string `json:"product_ids,omitempty"`
}

type CollectionCreateForm struct {
	JwtClaimsInfo

	Name        string             `json:"name,omitempty"`
	Description string             `json:"description,omitempty"`
	VI          *CollectionContent `json:"vi,omitempty"`

	Attachments *Attachments `json:"attachments,omitempty"`

	ProductIds pq.StringArray `json:"product_ids,omitempty"`
}
