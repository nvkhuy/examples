package models

import (
	"github.com/engineeringinflow/inflow-backend/pkg/db"
	"github.com/engineeringinflow/inflow-backend/pkg/errs"
)

// Validate validate category's info for create
func (record *Category) Validate(db *db.DB) error {
	if record.IsSlugExists(db) {
		return errs.ErrCategoryTaken
	}
	return nil
}

/******************** Utils ********************/
func (record *Category) IsSlugExists(db *db.DB) bool {
	if record.Slug != "" {
		var err = db.Select("ID").Where("id <> ?", record.ID).Where("slug = ?", record.Slug).Take(&Category{}).Error
		return !db.IsRecordNotFoundError(err)

	}
	return false
}
