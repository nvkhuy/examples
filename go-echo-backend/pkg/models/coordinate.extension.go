package models

import (
	"github.com/engineeringinflow/inflow-backend/pkg/db"
	"github.com/thaitanloi365/go-utils"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (coordinate *Coordinate) BeforeCreate(tx *gorm.DB) (err error) {

	var dataStr = coordinate.Coordinate.ToJsonString()
	coordinate.ID = utils.MD5(dataStr)

	tx.Statement.AddClause(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		UpdateAll: true,
	})

	return
}

func (coordinate *Coordinate) BeforeSave(tx *gorm.DB) (err error) {
	var dataStr = coordinate.Coordinate.ToJsonString()
	coordinate.ID = utils.MD5(dataStr)

	tx.Statement.AddClause(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		UpdateAll: true,
	})

	return
}

func (coordinate *Coordinate) CreateOrUpdate(db *db.DB, needFetchGeo ...bool) (err error) {
	if len(needFetchGeo) > 0 && needFetchGeo[0] {
		coordinate.GetLatLng()
	}

	var dataStr = coordinate.Coordinate.ToJsonString()
	coordinate.ID = utils.MD5(dataStr)

	var result = db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		UpdateAll: true,
	}).Create(coordinate)
	if result.Error != nil {
		return
	}

	return
}

func (coordinate *Coordinate) HasLatLng() bool {
	return coordinate.Lat != nil && coordinate.Lng != nil
}
