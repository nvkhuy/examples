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

type AuditValue struct {
	PoQuantity float64 `gorm:"type:decimal(20,4);default:0.0" json:"po_quantity,omitempty"`
	Actual     float64 `gorm:"type:decimal(20,4);default:0.0" json:"actual,omitempty"`
}

type BulkPurchaseOrderOrderClosingItem struct {
	ID            string      `json:"id,omitempty"`
	Color         string      `json:"color,omitempty"`
	Size          *AuditValue `json:"size,omitempty"`
	UnitPrice     *AuditValue `json:"unit_price,omitempty"`
	TotalQuantity *AuditValue `json:"total_quantity,omitempty"`
	TotalAmount   *AuditValue `json:"total_amount,omitempty"`
}

type BulkPurchaseOrderOrderClosingItems []*BulkPurchaseOrderOrderClosingItem

type BulkPurchaseOrderOrderClosing struct {
	OrderClosingDate    int64                                `json:"order_closing_date,omitempty"`
	PurchaseOrderNumber string                               `json:"purchase_order_number,omitempty"`
	PurchaseOrderDate   int                                  `json:"purchase_order_date,omitempty"`
	OrderItems          []*BulkPurchaseOrderOrderClosingItem `json:"order_items,omitempty"`
}

// Value return json value, implement driver.Valuer interface
func (m BulkPurchaseOrderOrderClosing) Value() (driver.Value, error) {
	ba, err := m.MarshalJSON()
	return string(ba), err
}

// Scan scan value into Jsonb, implements sql.Scanner interface
func (m *BulkPurchaseOrderOrderClosing) Scan(val interface{}) error {
	if val == nil {
		*m = *new(BulkPurchaseOrderOrderClosing)
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
	t := BulkPurchaseOrderOrderClosing{}
	err := json.Unmarshal(ba, &t)
	*m = t
	return err
}

// MarshalJSON to output non base64 encoded []byte
func (m *BulkPurchaseOrderOrderClosing) MarshalJSON() ([]byte, error) {
	type Alias BulkPurchaseOrderOrderClosing
	return json.Marshal(&struct {
		*Alias
	}{
		Alias: (*Alias)(m),
	})
}

// UnmarshalJSON to deserialize []byte
func (m *BulkPurchaseOrderOrderClosing) UnmarshalJSON(b []byte) error {
	if string(b) == "" || string(b) == "null" {
		return nil
	}

	type Alias BulkPurchaseOrderOrderClosing
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(m),
	}
	err := json.Unmarshal(b, &aux)
	return err
}

// GormDataType gorm common data type
func (m BulkPurchaseOrderOrderClosing) GormDataType() string {
	return "BulkPurchaseOrderOrderClosing"
}

// GormDBDataType gorm db data type
func (BulkPurchaseOrderOrderClosing) GormDBDataType(db *gorm.DB, field *schema.Field) string {
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

func (jm BulkPurchaseOrderOrderClosing) GormValue(ctx context.Context, db *gorm.DB) clause.Expr {
	data, _ := jm.MarshalJSON()
	switch db.Dialector.Name() {
	case "mysql":
		if v, ok := db.Dialector.(*mysql.Dialector); ok && !strings.Contains(v.ServerVersion, "MariaDB") {
			return gorm.Expr("CAST(? AS JSON)", string(data))
		}
	}
	return gorm.Expr("?", string(data))
}
