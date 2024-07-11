package models

import (
	"context"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"
	"strings"
)

type CollectionContent struct {
	Model

	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
}

// Value return json value, implement driver.Valuer interface
func (p CollectionContent) Value() (driver.Value, error) {
	ba, err := p.MarshalJSON()
	return string(ba), err
}

// Scan scan value into Jsonb, implements sql.Scanner interface
func (p *CollectionContent) Scan(val interface{}) error {
	if val == nil {
		*p = *new(CollectionContent)
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
	t := CollectionContent{}
	err := json.Unmarshal(ba, &t)
	*p = t
	return err
}

// MarshalJSON to output non base64 encoded []byte
func (p *CollectionContent) MarshalJSON() ([]byte, error) {
	type Alias CollectionContent
	return json.Marshal(&struct {
		*Alias
	}{
		Alias: (*Alias)(p),
	})
}

// UnmarshalJSON to deserialize []byte
func (p *CollectionContent) UnmarshalJSON(b []byte) error {
	if string(b) == "" || string(b) == "null" {
		return nil
	}

	type Alias CollectionContent
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(p),
	}
	err := json.Unmarshal(b, &aux)
	return err
}

// GormDataType gorm common data type
func (p CollectionContent) GormDataType() string {
	return "CollectionContent"
}

// GormDBDataType gorm db data type
func (CollectionContent) GormDBDataType(db *gorm.DB, field *schema.Field) string {
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

func (p CollectionContent) GormValue(ctx context.Context, db *gorm.DB) clause.Expr {
	data, _ := p.MarshalJSON()
	switch db.Dialector.Name() {
	case "mysql":
		if v, ok := db.Dialector.(*mysql.Dialector); ok && !strings.Contains(v.ServerVersion, "MariaDB") {
			return gorm.Expr("CAST(? AS JSON)", string(data))
		}
	}
	return gorm.Expr("?", string(data))
}
