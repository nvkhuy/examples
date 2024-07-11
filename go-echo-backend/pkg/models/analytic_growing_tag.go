package models

type AnalyticGrowingTags []*AnalyticGrowingTag

type AnalyticGrowingTag struct {
	Model
	Key   string  `gorm:"uniqueIndex:idx_growing_tags"  json:"key"`
	Value string  `json:"value"`
	Rate  float64 `json:"rate"`
}

func (AnalyticGrowingTag) TableName() string {
	return "growing_tags"
}
