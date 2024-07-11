package models

import "github.com/engineeringinflow/inflow-backend/pkg/models/enums"

// Shop Shop's model
type Page struct {
	Model

	Title    string `json:"title,omitempty"`
	Url      string `json:"url,omitempty"`
	PageType string `json:"page_type,omitempty"`
}

type PageUpdateForm struct {
	JwtClaimsInfo

	Title    string `json:"title,omitempty"`
	Url      string `json:"url,omitempty"`
	PageType string `json:"page_type,omitempty"`
}

type PageDetailResponse struct {
	ID    string `json:"id"`
	Title string `json:"title,omitempty"`
	Url   string `json:"url,omitempty"`

	Content []*PageSection `json:"content,omitempty"`
}

type PageSectionResponse struct {
	ID          string                `json:"id,omitempty"`
	Title       string                `json:"title,omitempty"`
	Content     string                `json:"content,omitempty"`
	SectionType enums.PageSectionType `json:"section_type"`
	Metadata    interface{}           `json:"metadata,omitempty"`
	Products    Products              `json:"products,omitempty"`
	Categories  CategorySlice         `json:"categories,omitempty"`
	Collections CollectionSlice       `json:"collections,omitempty"`
	Order       int                   `json:"order"`
}

type PageWithSectionUpdateForm struct {
	JwtClaimsInfo

	Title   string                   `json:"title,omitempty"`
	Url     string                   `json:"url,omitempty"`
	Content []*PageSectionUpdateForm `json:"content,omitempty"`
}
