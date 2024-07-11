package hubspot

import "time"

type PagingNext struct {
	After string `json:"after,omitempty"`
	Link  string `json:"link,omitempty"`
}

type Paging struct {
	Next *PagingNext `json:"next,omitempty"`
}
type Pagination[T any] struct {
	Results []T     `json:"results"`
	Paging  *Paging `json:"paging"`
}

type DataProperties[T any] struct {
	ID         string     `json:"id,omitempty"`
	Properties T          `json:"properties,omitempty"`
	CreatedAt  *time.Time `json:"createdAt"`
	UpdatedAt  *time.Time `json:"updatedAt"`
	Archived   bool       `json:"archived,omitempty"`
}

type Results[T any] struct {
	Status      string     `json:"status,omitempty"`
	Results     []T        `json:"results,omitempty"`
	StartedAt   *time.Time `json:"startedAt"`
	CompletedAt *time.Time `json:"completedAt"`
}
