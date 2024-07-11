package models

import (
	"strings"

	"gorm.io/gorm"
)

func makeTitle(object *Variant) string {
	var data []string
	if object.Color != "" {
		data = append(data, object.Color)
	}
	if object.Size != "" {
		data = append(data, object.Size)
	}
	if object.Material != "" {
		data = append(data, object.Material)
	}
	result := strings.Join(data, ",")
	return result
}

func (object *Variant) BeforeSave(tx *gorm.DB) (err error) {
	if object.Title != "" {
		return
	}
	tx.Statement.SetColumn("title", makeTitle(object))
	return
}
