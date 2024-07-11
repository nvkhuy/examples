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

type Attachments []*Attachment

// Value return json value, implement driver.Valuer interface
func (m Attachments) Value() (driver.Value, error) {
	ba, err := m.MarshalJSON()
	return string(ba), err
}

// Scan scan value into Jsonb, implements sql.Scanner interface
func (m *Attachments) Scan(val interface{}) error {
	if val == nil {
		*m = *new(Attachments)
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
	t := Attachments{}
	err := json.Unmarshal(ba, &t)
	*m = t
	return err
}

// MarshalJSON to output non base64 encoded []byte
func (m *Attachments) MarshalJSON() ([]byte, error) {
	type Alias Attachments
	aux := (*Alias)(m)

	return json.Marshal(&aux)
}

// UnmarshalJSON to deserialize []byte
func (m *Attachments) UnmarshalJSON(b []byte) error {
	if string(b) == "" || string(b) == "null" {
		return nil
	}

	type Alias Attachments
	aux := (*Alias)(m)
	err := json.Unmarshal(b, &aux)
	return err
}

// GormDataType gorm common data type
func (m Attachments) GormDataType() string {
	return "Attachments"
}

// GormDBDataType gorm db data type
func (Attachments) GormDBDataType(db *gorm.DB, field *schema.Field) string {
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

func (m Attachments) GormValue(ctx context.Context, db *gorm.DB) clause.Expr {
	data, _ := m.MarshalJSON()
	switch db.Dialector.Name() {
	case "mysql":
		if v, ok := db.Dialector.(*mysql.Dialector); ok && !strings.Contains(v.ServerVersion, "MariaDB") {
			return gorm.Expr("CAST(? AS JSON)", string(data))
		}
	}
	return gorm.Expr("?", string(data))
}

func (m Attachments) GenerateFileURL() Attachments {
	for index := range m {
		m[index] = m[index].GenerateFileURL()
	}

	return m
}
