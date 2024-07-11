package models

import (
	"fmt"

	"github.com/engineeringinflow/inflow-backend/pkg/helper"
	"github.com/gosimple/slug"
	"gorm.io/gorm"
)

func (p *Product) BeforeSave(tx *gorm.DB) (err error) {
	if p != nil {
		p.generateSlug(tx)
	}
	return
}

func (p *Product) BeforeUpdate(tx *gorm.DB) (err error) {
	if p != nil {
		p.generateSlug(tx)
	}
	return
}

func (p *Product) BeforeCreate(tx *gorm.DB) (err error) {
	if p != nil {
		p.generateSlug(tx)
	}
	return
}

/*
create -> name -> name
create -> name -> name-1

update -> ten -> name -> name | x -> name-01
*/
func (p *Product) generateSlug(tx *gorm.DB) {
	if p == nil {
		return
	}
	if p.Name == "" {
		return
	}
	path := slug.Make(p.Name)
	origin := path
	isExits := true
	inc := 0
	for isExits {
		var find Product
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

	if p.Vi != nil {
		if p.Vi.Name != "" {
			path = slug.Make(p.Vi.Name)
			origin = path
			isExits = true
			inc = 0
			for isExits {
				var find Product
				q := tx.Select("id").Where("vi ->> 'slug' = ?", path)
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
			p.Vi.Slug = path
		} else {
			p.Vi.Slug = ""
		}
	}
}

func (records Products) ToExcel() ([]byte, error) {

	var data = [][]interface{}{
		{"Slug", "Name", "Category", "Bullet Points", "Gender", "Style", "Country", "Currency"},
	}
	for _, record := range records {
		data = append(data, []interface{}{
			record.Slug,
			record.Name,
			func() string {
				if record.Category != nil {
					return record.Category.Name
				}
				return ""
			}(),
			record.BulletPoints,
			record.Gender,
			record.Style,
			record.CountryCode.GetCountryName(),
			record.Currency,
		})
	}

	return helper.ToExcel(data)

}
