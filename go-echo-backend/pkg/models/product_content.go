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

type ProductContent struct {
	Name             string `json:"name,omitempty"`
	Slug             string `json:"slug"`
	ShortDescription string `json:"short_description,omitempty"`
	Description      string `json:"description,omitempty"`
}

// implements sql.Scanner interface
func (p *ProductContent) Scan(val interface{}) error {
	if val == nil {
		p = nil
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
	content := ProductContent{}
	err := json.Unmarshal(ba, &content)
	if err != nil {
		return err
	}
	*p = content
	return nil
}

// implement driver.Valuer interface
func (p ProductContent) Value() (driver.Value, error) {
	return p.MarshalJSON()
}

func (p ProductContent) MarshalJSON() ([]byte, error) {
	type Alias ProductContent
	return json.Marshal(&struct {
		*Alias
	}{
		Alias: (*Alias)(&p),
	})
}

func (p *ProductContent) UnmarshalJSON(b []byte) error {
	if string(b) == "null" {
		return nil
	}
	type Alias ProductContent
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(p),
	}
	err := json.Unmarshal(b, &aux)
	return err
}

func (p ProductContent) GormDataType() string {
	return "ProductContent"
}

// GormDBDataType gorm db data type
func (ProductContent) GormDBDataType(db *gorm.DB, field *schema.Field) string {
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

func (p ProductContent) GormValue(ctx context.Context, db *gorm.DB) clause.Expr {
	data, _ := p.MarshalJSON()
	switch db.Dialector.Name() {
	case "mysql":
		if v, ok := db.Dialector.(*mysql.Dialector); ok && !strings.Contains(v.ServerVersion, "MariaDB") {
			return gorm.Expr("CAST(? AS JSON)", string(data))
		}
	}
	return gorm.Expr("?", string(data))
}
