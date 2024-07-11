package models

import (
	"fmt"
	"gorm.io/gorm"
	"strings"
)

func (s *SettingSEO) BeforeCreate(tx *gorm.DB) (err error) {
	if s != nil {
		return s.generateRoute(tx)
	}
	return
}

func (s *SettingSEO) BeforeUpdate(tx *gorm.DB) (err error) {
	if s != nil {
		return s.generateRoute(tx)
	}
	return
}

func (s *SettingSEO) generateRoute(tx *gorm.DB) (err error) {
	if s == nil {
		return
	}
	path := strings.ToLower(strings.TrimSpace(s.Route))
	origin := path
	isExits := true
	inc := 0
	for isExits {
		var find SettingSEO
		q := tx.Select("id").Where("route = ? and language_code = ?", path, s.LanguageCode)
		if s.ID != "" {
			q = q.Where("id != ?", s.ID)
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
	s.Route = path

	return
}
