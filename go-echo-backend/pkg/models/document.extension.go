package models

import (
	"fmt"

	"github.com/gosimple/slug"
	"gorm.io/gorm"
)

func (d *Document) BeforeUpdate(tx *gorm.DB) (err error) {
	if d != nil {
		d.generateSlug(tx)
	}
	return
}

func (d *Document) BeforeCreate(tx *gorm.DB) (err error) {
	if d != nil {
		d.generateSlug(tx)
	}
	return
}

func (d *Document) generateSlug(tx *gorm.DB) {
	if d == nil {
		return
	}
	path := slug.Make(d.Title)
	origin := path
	isExits := true
	inc := 0
	for isExits {
		var find Document
		q := tx.Select("id").Where("slug = ?", path)
		if d.ID != "" {
			q = q.Where("id != ?", d.ID)
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
	d.Slug = path

	if d.Vi != nil && d.Vi.Title != "" {
		path = slug.Make(d.Vi.Title)
		origin = path
		isExits = true
		inc = 0
		for isExits {
			var find Document
			q := tx.Select("id").Where("vi ->> 'slug' = ?", path)
			if d.ID != "" {
				q = q.Where("id != ?", d.ID)
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
		d.Vi.Slug = path
	}

	return
}
