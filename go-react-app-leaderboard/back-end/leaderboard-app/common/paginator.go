package common

import "reflect"

type PaginateParams struct {
	Page  int `json:"page" form:"page" query:"page" param:"page"`
	Limit int `json:"limit" form:"limit" query:"limit" param:"limit"`
}

type Paginator interface {
	Data() interface{}
	Page() int
	Total() int
	Limit() int
}

type PaginatorJsonResp struct {
	Data  interface{} `json:"data"`
	Page  int         `json:"page"`
	Limit int         `json:"limit"`
	Total int         `json:"total"`
}

// SlicePaginator struct implements the Paginator interface
type SlicePaginator struct {
	data  interface{}
	total int
	page  int
	limit int
}

// NewSlicePaginator creates a new instance of SlicePaginator
func NewSlicePaginator(data interface{}, page, limit int) *SlicePaginator {
	sp := &SlicePaginator{
		data:  data,
		page:  page,
		limit: limit,
	}
	sp.SetPage(sp.page)
	sp.SetLimit(sp.limit)
	return sp
}

// Data returns the current page of data
func (p *SlicePaginator) Data() interface{} {
	val := reflect.ValueOf(p.data)
	if val.Kind() != reflect.Slice {
		return nil
	}
	p.total = val.Len()

	start := (p.page - 1) * p.limit
	end := start + p.limit

	if start >= val.Len() {
		return nil
	}
	if end > val.Len() {
		end = val.Len()
	}

	return val.Slice(start, end).Interface()
}

// Page returns the current page number
func (p *SlicePaginator) Page() int {
	return p.page
}

// Limit returns the limit of items per page
func (p *SlicePaginator) Limit() int {
	return p.limit
}

// Total returns the total of items
func (p *SlicePaginator) Total() int {
	return p.total
}

// SetPage sets the current page
func (p *SlicePaginator) SetPage(page int) {
	if page < 1 {
		page = 1
	}
	p.page = page
}

// SetLimit sets the current limit
func (p *SlicePaginator) SetLimit(limit int) {
	if limit < 1 {
		limit = 10
	}
	p.limit = limit
}

func (p *SlicePaginator) Json() *PaginatorJsonResp {
	return &PaginatorJsonResp{
		Data:  p.Data(),
		Page:  p.Page(),
		Limit: p.Limit(),
		Total: p.Total(),
	}
}
