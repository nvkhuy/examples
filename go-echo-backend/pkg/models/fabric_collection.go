package models

type FabricCollection struct {
	Model
	ReferenceID string `gorm:"unique" json:"reference_id"`

	Name string `json:"name,omitempty"`
	Slug string `gorm:"unique" json:"slug"`

	VI *FabricCollectionContent `json:"vi,omitempty"`

	Fabrics []Fabric `gorm:"-" json:"fabrics,omitempty"`
}

type FabricInCollection struct {
	FabricID           string     `gorm:"primaryKey" json:"fabric_id,omitempty"`
	FabricCollectionID string     `gorm:"primaryKey" json:"fabric_collection_id,omitempty"`
	CreatedAt          int64      `gorm:"autoCreateTime" json:"created_at,omitempty"`
	UpdatedAt          int64      `gorm:"autoUpdateTime" json:"updated_at,omitempty"`
	DeletedAt          *DeletedAt `sql:"index" json:"deleted_at,omitempty" swaggertype:"primitive,integer"`
}
