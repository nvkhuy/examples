package models

import (
	"fmt"

	simpleslug "github.com/gosimple/slug"
	"gorm.io/gorm"
)

func (record *BlogCategory) BeforeSave(tx *gorm.DB) error {
	if record != nil {
		record.generateSlug(tx)
	}
	return nil
}

func (record *BlogCategory) BeforeUpdate(tx *gorm.DB) error {
	if record != nil {
		record.generateSlug(tx)
	}
	return nil
}

func (record *BlogCategory) BeforeCreate(tx *gorm.DB) (err error) {
	if record != nil {
		record.generateSlug(tx)
	}
	return nil
}

func (record *BlogCategory) generateSlug(tx *gorm.DB) {
	if record == nil {
		return
	}
	path := simpleslug.Make(record.Name)
	origin := path
	isExits := true
	inc := 0
	for isExits {
		var find BlogCategory
		q := tx.Select("id").Where("slug = ?", path)
		if record.ID != "" {
			q = q.Where("id != ?", record.ID)
		}
		q.First(&find)

		if find.ID == "" {
			isExits = false
		}
		if isExits {
			inc += 1
			path = fmt.Sprintf("%s-%d", origin, inc)
		}
	}
	record.Slug = path

	if record.Vi != nil && record.Vi.Name != "" {
		path = simpleslug.Make(record.Vi.Name)
		origin = path
		isExits = true
		inc = 0
		for isExits {
			var find BlogCategory
			q := tx.Select("id").Where("vi ->> 'slug' = ?", path)
			if record.ID != "" {
				q = q.Where("id != ?", record.ID)
			}
			q.First(&find)

			if find.ID == "" {
				isExits = false
			}
			if isExits {
				inc += 1
				path = fmt.Sprintf("%s-%d", origin, inc)
			}
		}
		record.Vi.Slug = path
	} else {
		record.Vi = &BlogCategoryContent{}
	}
}
