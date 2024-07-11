package models

import (
	"context"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/engineeringinflow/inflow-backend/pkg/models/enums"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"
)

type DocumentContent struct {
	Title   string               `json:"title,omitempty"`
	Status  enums.DocumentStatus `json:"status,omitempty"`
	Content string               `json:"content,omitempty"`
	Slug    string               `json:"slug,omitempty"`
}

// implements sql.Scanner interface
func (d *DocumentContent) Scan(val interface{}) error {
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
	document := DocumentContent{}
	err := json.Unmarshal(ba, &document)
	if err != nil {
		return err
	}
	*d = document
	return nil
}

// implement driver.Valuer interface
func (d DocumentContent) Value() (driver.Value, error) {
	return d.MarshalJSON()
}

func (d DocumentContent) MarshalJSON() ([]byte, error) {
	type Alias DocumentContent
	return json.Marshal(&struct {
		*Alias
	}{
		Alias: (*Alias)(&d),
	})
}

func (d *DocumentContent) UnmarshalJSON(b []byte) error {
	if string(b) == "null" {
		return nil
	}
	type Alias DocumentContent
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(d),
	}
	err := json.Unmarshal(b, &aux)
	return err
}

func (p DocumentContent) GormDataType() string {
	return "DocumentContent"
}

// GormDBDataType gorm db data type
func (DocumentContent) GormDBDataType(db *gorm.DB, field *schema.Field) string {
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

func (d DocumentContent) GormValue(ctx context.Context, db *gorm.DB) clause.Expr {
	data, _ := d.MarshalJSON()
	switch db.Dialector.Name() {
	case "mysql":
		if v, ok := db.Dialector.(*mysql.Dialector); ok && !strings.Contains(v.ServerVersion, "MariaDB") {
			return gorm.Expr("CAST(? AS JSON)", string(data))
		}
	}
	return gorm.Expr("?", string(data))
}
