package models

import (
	"context"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"
)

type CategoryContent struct {
	Name string `json:"name,omitempty"`
	Slug string `json:"slug,omitempty"`
}

// implements sql.Scanner interface
func (d *CategoryContent) Scan(val interface{}) error {
	if val == nil {
		d = nil
		return nil
	}
	var ba []byte
	switch v := val.(type) {
	case []byte:
		ba = v
	case string:
		ba = []byte(v)
	default:
		return errors.New(fmt.Sprint("Failed to unmarshal JSONB value:", val))
	}
	document := CategoryContent{}
	err := json.Unmarshal(ba, &document)
	if err != nil {
		return err
	}
	*d = document
	return nil
}

// implement driver.Valuer interface
func (d CategoryContent) Value() (driver.Value, error) {
	return d.MarshalJSON()
}

func (d CategoryContent) MarshalJSON() ([]byte, error) {
	type Alias CategoryContent
	return json.Marshal(&struct {
		*Alias
	}{
		Alias: (*Alias)(&d),
	})
}

func (d *CategoryContent) UnmarshalJSON(b []byte) error {
	if string(b) == "null" {
		return nil
	}
	type Alias CategoryContent
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(d),
	}
	err := json.Unmarshal(b, &aux)
	return err
}

func (p CategoryContent) GormDataType() string {
	return "CategoryContent"
}

// GormDBDataType gorm db data type
func (CategoryContent) GormDBDataType(db *gorm.DB, field *schema.Field) string {
	switch db.Dialector.Name() {
	case "sqlite":
		return "JSON"
	case "mysql":
		return "JSON"
	case "postgres":
		return "JSONB"
	case "sqlserver":
		return "NVARCHAR(MAX)"
	}
	return ""
}

func (d CategoryContent) GormValue(ctx context.Context, db *gorm.DB) clause.Expr {
	data, _ := d.MarshalJSON()
	switch db.Dialector.Name() {
	case "mysql":
		if v, ok := db.Dialector.(*mysql.Dialector); ok && !strings.Contains(v.ServerVersion, "MariaDB") {
			return gorm.Expr("CAST(? AS JSON)", string(data))
		}
	}
	return gorm.Expr("?", string(data))
}
