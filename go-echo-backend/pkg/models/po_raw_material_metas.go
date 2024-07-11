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

type PoRawMaterialMetas []*PoRawMaterialMeta

// Value return json value, implement driver.Valuer interface
func (m PoRawMaterialMetas) Value() (driver.Value, error) {
	ba, err := m.MarshalJSON()
	return string(ba), err
}

// Scan scan value into Jsonb, implements sql.Scanner interface
func (m *PoRawMaterialMetas) Scan(val interface{}) error {
	if val == nil {
		*m = *new(PoRawMaterialMetas)
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
	t := PoRawMaterialMetas{}
	err := json.Unmarshal(ba, &t)
	*m = t
	return err
}

// MarshalJSON to output non base64 encoded []byte
func (m *PoRawMaterialMetas) MarshalJSON() ([]byte, error) {
	type Alias PoRawMaterialMetas
	aux := (*Alias)(m)

	return json.Marshal(&aux)
}

// UnmarshalJSON to deserialize []byte
func (m *PoRawMaterialMetas) UnmarshalJSON(b []byte) error {
	if string(b) == "" || string(b) == "null" {
		return nil
	}

	type Alias PoRawMaterialMetas
	aux := (*Alias)(m)
	err := json.Unmarshal(b, &aux)
	return err
}

// GormDataType gorm common data type
func (m PoRawMaterialMetas) GormDataType() string {
	return "PoRawMaterialMetas"
}

// GormDBDataType gorm db data type
func (PoRawMaterialMetas) GormDBDataType(db *gorm.DB, field *schema.Field) string {
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

func (jm PoRawMaterialMetas) GormValue(ctx context.Context, db *gorm.DB) clause.Expr {
	data, _ := jm.MarshalJSON()
	switch db.Dialector.Name() {
	case "mysql":
		if v, ok := db.Dialector.(*mysql.Dialector); ok && !strings.Contains(v.ServerVersion, "MariaDB") {
			return gorm.Expr("CAST(? AS JSON)", string(data))
		}
	}
	return gorm.Expr("?", string(data))
}

func (m PoRawMaterialMetas) GenerateFileURL() PoRawMaterialMetas {
	for i := 0; i < len(m); i++ {
		m[i].Attachments = m[i].Attachments.GenerateFileURL()
	}
	return m
}
