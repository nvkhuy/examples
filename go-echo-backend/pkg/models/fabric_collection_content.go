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

type FabricCollectionContent struct {
	Name string `json:"name,omitempty"`
	Slug string `gorm:"unique" json:"slug"`
}

// Value return json value, implement driver.Valuer interface
func (p FabricCollectionContent) Value() (driver.Value, error) {
	ba, err := p.MarshalJSON()
	return string(ba), err
}

// Scan scan value into Jsonb, implements sql.Scanner interface
func (p *FabricCollectionContent) Scan(val interface{}) error {
	if val == nil {
		*p = *new(FabricCollectionContent)
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
	t := FabricCollectionContent{}
	err := json.Unmarshal(ba, &t)
	*p = t
	return err
}

// MarshalJSON to output non base64 encoded []byte
func (p *FabricCollectionContent) MarshalJSON() ([]byte, error) {
	type Alias FabricCollectionContent
	return json.Marshal(&struct {
		*Alias
	}{
		Alias: (*Alias)(p),
	})
}

// UnmarshalJSON to deserialize []byte
func (p *FabricCollectionContent) UnmarshalJSON(b []byte) error {
	if string(b) == "" || string(b) == "null" {
		return nil
	}

	type Alias FabricCollectionContent
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(p),
	}
	err := json.Unmarshal(b, &aux)
	return err
}

// GormDataType gorm common data type
func (p FabricCollectionContent) GormDataType() string {
	return "FabricCollectionContent"
}

// GormDBDataType gorm db data type
func (FabricCollectionContent) GormDBDataType(db *gorm.DB, field *schema.Field) string {
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

func (p FabricCollectionContent) GormValue(ctx context.Context, db *gorm.DB) clause.Expr {
	data, _ := p.MarshalJSON()
	switch db.Dialector.Name() {
	case "mysql":
		if v, ok := db.Dialector.(*mysql.Dialector); ok && !strings.Contains(v.ServerVersion, "MariaDB") {
			return gorm.Expr("CAST(? AS JSON)", string(data))
		}
	}
	return gorm.Expr("?", string(data))
}
