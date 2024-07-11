package models

import (
	"mime/multipart"
)

type ProductFileUploadInfo struct {
	Model
	SiteName   string                         `json:"site_name"`
	Status     string                         `json:"status"`
	Attachment Attachment                     `json:"attachment"`
	FailReason string                         `json:"fail_reason"`
	ScrapeDate int64                          `json:"scrape_date"`
	Metadata   *ProductFileUploadInfoMetadata `json:"metadata"`
}

type UploadProductFileRequest struct {
	JwtClaimsInfo
	SiteName   string                `form:"site_name" validate:"required"`
	ScrapeDate int                   `form:"scrape_date" validate:"required"`
	File       *multipart.FileHeader `form:"file" validate:"required"`

	// Attachment Attachment `json:"attachment" validate:"required"`
}

type GetProductFileListRequest struct {
	JwtClaimsInfo
	PaginationParams
	// Attachment Attachment `json:"attachment" validate:"required"`
}
