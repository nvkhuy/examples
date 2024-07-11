package models

import (
	"context"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/engineeringinflow/inflow-backend/pkg/models/enums"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"
)

// ProductAttributeMetas
type ProductAttributeMetas []*ProductAttributeMeta

type ProductAttributeMeta struct {
	Attribute enums.ProductAttribute `json:"attribute" validate:"required"`
	Value     string                 `json:"value" validate:"required"`
	ColorName string                 `json:"color_name,omitempty"`
}

// Value return json value, implement driver.Valuer interface
func (m ProductAttributeMetas) Value() (driver.Value, error) {
	ba, err := m.MarshalJSON()
	return string(ba), err
}

// Scan scan value into Jsonb, implements sql.Scanner interface
func (m *ProductAttributeMetas) Scan(val interface{}) error {
	if val == nil {
		*m = *new(ProductAttributeMetas)
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
	t := ProductAttributeMetas{}
	err := json.Unmarshal(ba, &t)
	*m = t
	return err
}

// MarshalJSON to output non base64 encoded []byte
func (m *ProductAttributeMetas) MarshalJSON() ([]byte, error) {
	type Alias ProductAttributeMetas
	aux := (*Alias)(m)

	return json.Marshal(&aux)
}

// UnmarshalJSON to deserialize []byte
func (m *ProductAttributeMetas) UnmarshalJSON(b []byte) error {
	if string(b) == "" || string(b) == "null" {
		return nil
	}

	type Alias ProductAttributeMetas
	aux := (*Alias)(m)
	err := json.Unmarshal(b, &aux)
	return err
}

// GormDataType gorm common data type
func (m ProductAttributeMetas) GormDataType() string {
	return "ProductAttributeMetas"
}

// GormDBDataType gorm db data type
func (ProductAttributeMetas) GormDBDataType(db *gorm.DB, field *schema.Field) string {
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

func (jm ProductAttributeMetas) GormValue(ctx context.Context, db *gorm.DB) clause.Expr {
	data, _ := jm.MarshalJSON()

	return gorm.Expr("?", string(data))
}
