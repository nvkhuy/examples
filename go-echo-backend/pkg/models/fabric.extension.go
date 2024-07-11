package models

import (
	"fmt"

	"github.com/engineeringinflow/inflow-backend/pkg/helper"
	"github.com/gosimple/slug"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (f *Fabric) BeforeUpdate(tx *gorm.DB) (err error) {
	if f != nil {
		f.generateSlug(tx)
	}
	return
}

func (f *Fabric) BeforeCreate(tx *gorm.DB) (err error) {
	if f != nil {
		f.generateSlug(tx)
	}
	if f.ReferenceID == "" {
		var id = helper.GenerateFabricReferenceID()
		tx.Statement.SetColumn("ReferenceID", id)
		tx.Statement.AddClauseIfNotExists(clause.OnConflict{
			Columns: []clause.Column{{Name: "reference_id"}},
			DoUpdates: clause.Assignments(func() map[string]interface{} {
				id = helper.GenerateFabricReferenceID()
				return map[string]interface{}{"reference_id": id}
			}()),
		})
	}
	return
}

func (f *Fabric) generateSlug(tx *gorm.DB) {
	if f == nil {
		return
	}
	path := slug.Make(f.FabricType)
	origin := path
	isExits := true
	inc := 0
	for isExits {
		var find Fabric
		q := tx.Select("id").Where("slug = ?", path)
		if f.ID != "" {
			q = q.Where("id != ?", f.ID)
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
	f.Slug = path

	if f.VI != nil {
		if f.VI.FabricType != "" {
			path = slug.Make(f.VI.FabricType)
			origin = path
			isExits = true
			inc = 0
			for isExits {
				var find Fabric
				var q = tx.Select("ID").Where("vi ->> 'slug' = ?", path)
				if f.ID != "" {
					q = q.Where("id != ?", f.ID)
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
			f.VI.Slug = path
		} else {
			f.Slug = ""
		}
	}
}
