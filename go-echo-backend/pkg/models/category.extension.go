package models

import (
	"fmt"

	"github.com/gosimple/slug"
	"gorm.io/gorm"
)

func (c *Category) BeforeUpdate(tx *gorm.DB) (err error) {
	if c != nil {
		c.generateSlug(tx)
	}
	return
}

func (c *Category) BeforeCreate(tx *gorm.DB) (err error) {
	if c != nil {
		c.generateSlug(tx)
	}
	return
}

func (c *Category) generateSlug(tx *gorm.DB) {
	if c == nil {
		return
	}
	path := slug.Make(c.Name)
	origin := path
	isExits := true
	inc := 0
	for isExits {
		var find Category
		q := tx.Select("id").Where("slug = ?", path)
		if c.ID != "" {
			q = q.Where("id != ?", c.ID)
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
	c.Slug = path

	if c.Vi != nil && c.Vi.Name != "" {
		path = slug.Make(c.Vi.Name)
		origin = path
		isExits = true
		inc = 0
		for isExits {
			var find Category
			q := tx.Select("id").Where("vi ->> 'slug' = ?", path)
			if c.ID != "" {
				q = q.Where("id != ?", c.ID)
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
		c.Vi.Slug = path
	} else {
		c.Vi = &CategoryContent{}
	}
}
