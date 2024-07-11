package models

type DocumentTag struct {
	Model
	Name string `gorm:"not null;unique;size:200" json:"name"`
}

type CreateDocumentTagRequest struct {
	JwtClaimsInfo
	Name string `json:"name" validate:"required"`
}

type UpdateDocumentTagRequest struct {
	JwtClaimsInfo
	DocumentTagID string `json:"document_tag_id" param:"document_tag_id" validate:"required"`
	Name          string `json:"name" validate:"required"`
}
type DeleteDocumentTagRequest struct {
	JwtClaimsInfo
	DocumentTagID string `json:"document_tag_id" param:"document_tag_id" validate:"required"`
}

type GetDocumentTagListParams struct {
	PaginationParams
	JwtClaimsInfo
}

type DocumentTags []*DocumentTag

func (tags DocumentTags) IDs() []string {
	var IDs []string
	for _, tag := range tags {
		IDs = append(IDs, tag.ID)
	}
	return IDs
}
