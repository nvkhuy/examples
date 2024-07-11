package models

import (
	"context"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"
)

type BulkPurchaseOrderAdditionalItems []*BulkPurchaseOrderAdditionalItem

func (items BulkPurchaseOrderAdditionalItems) UpdateTotal() {
	for i := range items {
		items[i].TotalPrice = items[i].UnitPrice.MultipleInt(items[i].Qty).ToPtr()
	}
}

// Value return json value, implement driver.Valuer interface
func (m BulkPurchaseOrderAdditionalItems) Value() (driver.Value, error) {
	ba, err := m.MarshalJSON()
	return string(ba), err
}

// Scan scan value into Jsonb, implements sql.Scanner interface
func (m *BulkPurchaseOrderAdditionalItems) Scan(val interface{}) error {
	if val == nil {
		*m = *new(BulkPurchaseOrderAdditionalItems)
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
	t := BulkPurchaseOrderAdditionalItems{}
	err := json.Unmarshal(ba, &t)
	*m = t
	return err
}

// MarshalJSON to output non base64 encoded []byte
func (m *BulkPurchaseOrderAdditionalItems) MarshalJSON() ([]byte, error) {
	type Alias BulkPurchaseOrderAdditionalItems
	aux := (*Alias)(m)

	return json.Marshal(&aux)
}

// UnmarshalJSON to deserialize []byte
func (m *BulkPurchaseOrderAdditionalItems) UnmarshalJSON(b []byte) error {
	if string(b) == "" || string(b) == "null" {
		return nil
	}

	type Alias BulkPurchaseOrderAdditionalItems
	aux := (*Alias)(m)
	err := json.Unmarshal(b, &aux)
	return err
}

// GormDataType gorm common data type
func (m BulkPurchaseOrderAdditionalItems) GormDataType() string {
	return "BulkPurchaseOrderAdditionalItems"
}

// GormDBDataType gorm db data type
func (BulkPurchaseOrderAdditionalItems) GormDBDataType(db *gorm.DB, field *schema.Field) string {
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

func (jm BulkPurchaseOrderAdditionalItems) GormValue(ctx context.Context, db *gorm.DB) clause.Expr {
	data, _ := jm.MarshalJSON()

	return gorm.Expr("?", string(data))
}
