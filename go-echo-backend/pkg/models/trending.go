package models

import (
	"github.com/engineeringinflow/inflow-backend/pkg/models/enums"
	"github.com/lib/pq"
)

type Trending struct {
	Model
	Name               string               `gorm:"type:varchar(200)"  json:"name"`
	Description        string               `json:"description"`
	ProductTrendingIDs pq.StringArray       `gorm:"type:varchar(200)[]" json:"product_trending_ids,omitempty"`
	Status             enums.TrendingStatus `json:"status,omitempty" gorm:"default:'new'"`
	CoverAttachment    *Attachment          `json:"cover_attachment,omitempty"`
	IsAutoCreate       *bool                `gorm:"default:false" json:"is_auto_create,omitempty"`
	CollectionURL      string               `json:"collection_url,omitempty"`

	ProductTrendings []AnalyticProductTrending `gorm:"-" json:"product_trendings,omitempty"`
	Total            int                       `gorm:"-" json:"total,omitempty"`
}
