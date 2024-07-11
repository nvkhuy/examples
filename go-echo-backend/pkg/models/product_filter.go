package models

// "github.com/engineeringinflow/inflow-backend/pkg/models/enums"

// ProductReview's model
type ProductFilter struct {
	Name    string      `json:"name"`
	Key     string      `json:"key"`
	Type    string      `json:"type"`
	Options interface{} `json:"options"`
}

type ProductFilterOption struct {
	Name  string      `json:"name"`
	Value interface{} `json:"value"`
}
