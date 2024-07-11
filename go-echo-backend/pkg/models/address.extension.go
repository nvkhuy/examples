package models

import (
	"strings"

	"github.com/engineeringinflow/inflow-backend/pkg/db"
	"github.com/thaitanloi365/go-utils"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (address *Address) BeforeCreate(tx *gorm.DB) (err error) {
	address.PhoneNumber = strings.TrimSpace(address.PhoneNumber)
	address.ID = address.GenerateID()
	tx.Statement.AddClause(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		UpdateAll: true,
	})

	return
}

func (address *Address) BeforeSave(tx *gorm.DB) (err error) {
	address.ID = address.GenerateID()
	tx.Statement.AddClause(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		UpdateAll: true,
	})

	return
}

func (userAddress *Address) GenerateID() string {
	userAddress.ID = utils.MD5(userAddress.CoordinateID + userAddress.Email + userAddress.PhoneNumber + userAddress.Name + string(userAddress.AddressType))
	return userAddress.ID
}

func (userAddress *Address) CreateOrUpdate(db *db.DB) error {
	if userAddress.Coordinate != nil {
		if err := userAddress.Coordinate.CreateOrUpdate(db); err == nil {
			userAddress.CoordinateID = userAddress.Coordinate.ID
		}
	}
	userAddress.ID = userAddress.GenerateID()

	return db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		UpdateAll: true,
	}).Create(userAddress).Error

}
