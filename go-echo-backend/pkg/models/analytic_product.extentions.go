package models

import (
	"errors"
	"gorm.io/gorm"
)

func (AnalyticProduct) BeforeCreate(tx *gorm.DB) (err error) {
	return errors.New("write action not allowed")
}

func (AnalyticProduct) BeforeUpdate(tx *gorm.DB) (err error) {
	return errors.New("write action not allowed")
}

func (AnalyticProduct) BeforeDelete(tx *gorm.DB) (err error) {
	return errors.New("write action not allowed")
}

func (AnalyticProductChanges) BeforeCreate(tx *gorm.DB) (err error) {
	return errors.New("write action not allowed")
}

func (AnalyticProductChanges) BeforeUpdate(tx *gorm.DB) (err error) {
	return errors.New("write action not allowed")
}

func (AnalyticProductChanges) BeforeDelete(tx *gorm.DB) (err error) {
	return errors.New("write action not allowed")
}
