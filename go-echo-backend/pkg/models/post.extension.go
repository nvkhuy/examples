package models

import (
	"fmt"
	"strings"

	"github.com/gosimple/slug"
	"gorm.io/gorm"
)

func (p *Post) BeforeUpdate(tx *gorm.DB) (err error) {
	if p != nil {
		p.generateSlug(tx)
		p.transformContent()
	}
	return
}

func (p *Post) BeforeCreate(tx *gorm.DB) (err error) {
	if p != nil {
		p.generateSlug(tx)
		p.transformContent()
	}
	return
}

func (p *Post) generateSlug(tx *gorm.DB) {
	if p == nil {
		return
	}
	path := slug.Make(p.Title)
	origin := path
	isExits := true
	inc := 0
	for isExits {
		var find Post
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
			var find Post
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

func (p *Post) transformContent() {
	replaces := make(map[string]string)
	replaces[".style="] = "style="
	for k, v := range replaces {
		p.Content = strings.ReplaceAll(p.Content, k, v)
		if p.VI != nil {
			p.VI.Content = strings.ReplaceAll(p.VI.Content, k, v)
		}
	}
}
