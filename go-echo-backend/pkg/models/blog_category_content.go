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

type BlogCategoryContent struct {
	Name        string `json:"name,omitempty"`
	Subtitle    string `json:"sub_title,omitempty"`
	Description string `json:"description,omitempty"`
	Slug        string `json:"slug,omitempty"`
}

// Value return json value, implement driver.Valuer interface
func (b BlogCategoryContent) Value() (driver.Value, error) {
	ba, err := b.MarshalJSON()
	return string(ba), err
}

// Scan scan value into Jsonb, implements sql.Scanner interface
func (b *BlogCategoryContent) Scan(val interface{}) error {
	if val == nil {
		*b = *new(BlogCategoryContent)
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
	t := BlogCategoryContent{}
	err := json.Unmarshal(ba, &t)
	*b = t
	return err
}

// MarshalJSON to output non base64 encoded []byte
func (b *BlogCategoryContent) MarshalJSON() ([]byte, error) {
	type Alias BlogCategoryContent
	return json.Marshal(&struct {
		*Alias
	}{
		Alias: (*Alias)(b),
	})
}

// UnmarshalJSON to deserialize []byte
func (b *BlogCategoryContent) UnmarshalJSON(_b []byte) error {
	if string(_b) == "" || string(_b) == "null" {
		return nil
	}

	type Alias BlogCategoryContent
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(b),
	}
	err := json.Unmarshal(_b, &aux)
	return err
}

// GormDataType gorm common data type
func (BlogCategoryContent) GormDataType() string {
	return "BlogCategoryContent"
}

// GormDBDataType gorm db data type
func (BlogCategoryContent) GormDBDataType(db *gorm.DB, field *schema.Field) string {
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

func (b BlogCategoryContent) GormValue(ctx context.Context, db *gorm.DB) clause.Expr {
	data, _ := b.MarshalJSON()
	switch db.Dialector.Name() {
	case "mysql":
		if v, ok := db.Dialector.(*mysql.Dialector); ok && !strings.Contains(v.ServerVersion, "MariaDB") {
			return gorm.Expr("CAST(? AS JSON)", string(data))
		}
	}
	return gorm.Expr("?", string(data))
}
