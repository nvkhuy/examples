package models

import (
	"fmt"

	"github.com/gosimple/slug"
	"gorm.io/gorm"
)

func (p *AsFeaturedIn) BeforeUpdate(tx *gorm.DB) (err error) {
	if p != nil {
		p.generateSlug(tx)
	}
	return
}

func (p *AsFeaturedIn) BeforeCreate(tx *gorm.DB) (err error) {
	if p != nil {
		p.generateSlug(tx)
	}
	return
}

func (p *AsFeaturedIn) generateSlug(tx *gorm.DB) {
	if p == nil {
		return
	}
	path := slug.Make(p.Title)
	origin := path
	isExits := true
	inc := 0
	for isExits {
		var find AsFeaturedIn
		q := tx.Select("id").Where("slug = ?", path)
		if p.ID != "" {
			q = q.Where("id != ?", p.ID)
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
	p.Slug = path

	if p.VI != nil && p.VI.Title != "" {
		path = slug.Make(p.VI.Title)
		origin = path
		isExits = true
		inc = 0
		for isExits {
			var find AsFeaturedIn
			var q = tx.Select("ID").Where("vi ->> 'slug' = ?", path)
			if p.ID != "" {
				q = q.Where("id != ?", p.ID)
			}
			q.Select("ID").First(&find)

			if find.ID == "" {
				isExits = false
			}
			if isExits {
				inc += 1
				path = fmt.Sprintf("%s-%d", origin, inc)
			}
		}
		p.VI.Slug = path
	}

}
