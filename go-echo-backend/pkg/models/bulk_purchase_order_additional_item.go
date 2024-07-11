package models

import (
	"context"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/engineeringinflow/inflow-backend/pkg/models/price"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"
)

type BulkPurchaseOrderAdditionalItem struct {
	Attachments         Attachments  `json:"attachments,omitempty"`
	BulkPurchaseOrderID string       `json:"bulk_purchase_order_id,omitempty"`
	Sku                 string       `json:"sku,omitempty"`
	Qty                 int64        `gorm:"default:1" json:"qty"`
	UnitPrice           *price.Price `gorm:"type:decimal(20,4);default:0.0" json:"unit_price,omitempty"`
	TotalPrice          *price.Price `gorm:"type:decimal(20,4);default:0.0" json:"total_price,omitempty"`
	Description         string       `json:"description,omitempty"`
}

// Value return json value, implement driver.Valuer interface
func (m BulkPurchaseOrderAdditionalItem) Value() (driver.Value, error) {
	ba, err := m.MarshalJSON()
	return string(ba), err
}

// Scan scan value into Jsonb, implements sql.Scanner interface
func (m *BulkPurchaseOrderAdditionalItem) Scan(val interface{}) error {
	if val == nil {
		*m = *new(BulkPurchaseOrderAdditionalItem)
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
	t := BulkPurchaseOrderAdditionalItem{}
	err := json.Unmarshal(ba, &t)
	*m = t
	return err
}

// MarshalJSON to output non base64 encoded []byte
func (m *BulkPurchaseOrderAdditionalItem) MarshalJSON() ([]byte, error) {
	type Alias BulkPurchaseOrderAdditionalItem
	aux := (*Alias)(m)

	return json.Marshal(&aux)
}

// UnmarshalJSON to deserialize []byte
func (m *BulkPurchaseOrderAdditionalItem) UnmarshalJSON(b []byte) error {
	if string(b) == "" || string(b) == "null" {
		return nil
	}

	type Alias BulkPurchaseOrderAdditionalItem
	aux := (*Alias)(m)
	err := json.Unmarshal(b, &aux)
	return err
}

// GormDataType gorm common data type
func (m BulkPurchaseOrderAdditionalItem) GormDataType() string {
	return "BulkPurchaseOrderAdditionalItem"
}

// GormDBDataType gorm db data type
func (BulkPurchaseOrderAdditionalItem) GormDBDataType(db *gorm.DB, field *schema.Field) string {
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

func (jm BulkPurchaseOrderAdditionalItem) GormValue(ctx context.Context, db *gorm.DB) clause.Expr {
	data, _ := jm.MarshalJSON()

	return gorm.Expr("?", string(data))
}
