package models

type TaggedDocument struct {
	DocumentID    string `gorm:"uniqueIndex:idx_tagged_document" json:"document_id"`
	DocumentTagID string `gorm:"uniqueIndex:idx_tagged_document" json:"document_tag_id"`

	CreatedAt int64      `gorm:"autoCreateTime" json:"created_at,omitempty"`
	UpdatedAt int64      `gorm:"autoUpdateTime" json:"updated_at,omitempty"`
	DeletedAt *DeletedAt `sql:"index" json:"deleted_at,omitempty" swaggertype:"primitive,integer"`
}

type TaggedDocuments []*TaggedDocument

func (tds TaggedDocuments) DocumentIDs() []string {
	var documentIDs []string
	for _, tag := range tds {
		documentIDs = append(documentIDs, tag.DocumentID)
	}
	return documentIDs
}

func (tds TaggedDocuments) DocumentTagIDs() []string {
	var tagIDs []string
	for _, tag := range tds {
		tagIDs = append(tagIDs, tag.DocumentTagID)
	}
	return tagIDs
}
