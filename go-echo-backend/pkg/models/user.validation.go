package models

import (
	"github.com/engineeringinflow/inflow-backend/pkg/db"
	"github.com/engineeringinflow/inflow-backend/pkg/errs"
)

// Validate validate user's info for create
func (user *User) Validate(db *db.DB) error {
	if user.IsEmailExists(db) {
		return errs.ErrPhoneTaken
	}
	return nil
}

/******************** Utils ********************/
func (user *User) IsEmailExists(db *db.DB) bool {
	if user.Email != "" {
		var err = db.Select("ID", "Email").Where("id <> ?", user.ID).Where("email = ?", user.Email).Take(&User{}).Error
		return !db.IsRecordNotFoundError(err)

	}
	return false
}
