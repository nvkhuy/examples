package models

type DocumentCategory struct {
	Model
	Name        string `gorm:"not null;unique" json:"name"`
	Slug        string `gorm:"not null;unique" json:"slug"`
	Description string `json:"description,omitempty"`

	Order int `gorm:"default:0" json:"order,omitempty"`
}

type CreateDocumentCategoryRequest struct {
	JwtClaimsInfo
	Name        string `json:"name" validate:"required"`
	Description string `json:"description,omitempty"`
	Order       int    `json:"order,omitempty"`
}

type UpdateDocumentCategoryRequest struct {
	CreateDocumentCategoryRequest
	DocumentCategoryID string `param:"document_category_id" validate:"required"`
}

type GetDocumentCategoryParams struct {
	JwtClaimsInfo
	DocumentCategoryID string `param:"document_category_id" validate:"required"`
}

type GetDocumentCategoryListParams struct {
	PaginationParams
	JwtClaimsInfo
}
