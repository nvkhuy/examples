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

// DeliveryAttributes
type DeliveryAttributes []*DeliveryAttribute

type DeliveryAttribute struct {
	Attribute enums.CardAttribute `json:"attribute,omitempty"`
	Value     string              `json:"value,omitempty"`
}

// Value return json value, implement driver.Valuer interface
func (m DeliveryAttributes) Value() (driver.Value, error) {
	ba, err := m.MarshalJSON()
	return string(ba), err
}

// Scan scan value into Jsonb, implements sql.Scanner interface
func (m *DeliveryAttributes) Scan(val interface{}) error {
	if val == nil {
		*m = *new(DeliveryAttributes)
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
	t := DeliveryAttributes{}
	err := json.Unmarshal(ba, &t)
	*m = t
	return err
}

// MarshalJSON to output non base64 encoded []byte
func (m *DeliveryAttributes) MarshalJSON() ([]byte, error) {
	type Alias DeliveryAttributes
	aux := (*Alias)(m)

	return json.Marshal(&aux)
}

// UnmarshalJSON to deserialize []byte
func (m *DeliveryAttributes) UnmarshalJSON(b []byte) error {
	if string(b) == "" || string(b) == "null" {
		return nil
	}

	type Alias DeliveryAttributes
	aux := (*Alias)(m)
	err := json.Unmarshal(b, &aux)
	return err
}

// GormDataType gorm common data type
func (m DeliveryAttributes) GormDataType() string {
	return "DeliveryAttributes"
}

// GormDBDataType gorm db data type
func (DeliveryAttributes) GormDBDataType(db *gorm.DB, field *schema.Field) string {
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

func (jm DeliveryAttributes) GormValue(ctx context.Context, db *gorm.DB) clause.Expr {
	data, _ := jm.MarshalJSON()

	return gorm.Expr("?", string(data))
}
