package models

import (
	"fmt"
	"github.com/engineeringinflow/inflow-backend/pkg/helper"
	"github.com/gosimple/slug"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (f *FabricCollection) BeforeUpdate(tx *gorm.DB) (err error) {
	if f != nil {
		f.generateSlug(tx)
	}
	return
}

func (f *FabricCollection) BeforeCreate(tx *gorm.DB) (err error) {
	if f != nil {
		f.generateSlug(tx)
	}
	if f.ReferenceID == "" {
		var id = helper.GenerateFabricCollectionReferenceID()
		tx.Statement.SetColumn("ReferenceID", id)
		tx.Statement.AddClauseIfNotExists(clause.OnConflict{
			Columns: []clause.Column{{Name: "reference_id"}},
			DoUpdates: clause.Assignments(func() map[string]interface{} {
				id = helper.GenerateFabricCollectionReferenceID()
				return map[string]interface{}{"reference_id": id}
			}()),
		})
	}
	return
}

func (f *FabricCollection) generateSlug(tx *gorm.DB) {
	if f == nil {
		return
	}
	path := slug.Make(f.Name)
	origin := path
	isExits := true
	inc := 0
	for isExits {
		var find FabricCollection
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

	// VI
	if f.VI != nil && f.VI.Name != "" {
		path = slug.Make(f.VI.Name)
		origin = path
		isExits = true
		inc = 0
		for isExits {
			var find FabricCollection
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
	}
	return
}
