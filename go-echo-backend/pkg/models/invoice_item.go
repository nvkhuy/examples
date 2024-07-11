package models

import (
	"context"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/engineeringinflow/inflow-backend/pkg/models/price"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"
)

type InvoiceItem struct {
	ItemCode          string       `json:"item_code,omitempty"`
	Description       string       `json:"description,omitempty"`
	HSCode            string       `json:"hs_code,omitempty"`
	FabricComposition string       `json:"fabric_composition,omitempty"`
	Color             string       `json:"color,omitempty"`
	Size              *SizeQty     `json:"size,omitempty"`
	SizeName          string       `json:"size_name"`
	TotalQuantity     int64        `json:"total_quantity,omitempty"`
	UnitPrice         *price.Price `json:"unit_price,omitempty"`
	TotalAmount       *price.Price `json:"total_amount,omitempty"`
}

// Value return json value, implement driver.Valuer interface
func (m InvoiceItem) Value() (driver.Value, error) {
	ba, err := m.MarshalJSON()
	return string(ba), err
}

// Scan scan value into Jsonb, implements sql.Scanner interface
func (m *InvoiceItem) Scan(val interface{}) error {
	if val == nil {
		*m = *new(InvoiceItem)
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
	t := InvoiceItem{}
	err := json.Unmarshal(ba, &t)
	*m = t
	return err
}

// MarshalJSON to output non base64 encoded []byte
func (m *InvoiceItem) MarshalJSON() ([]byte, error) {
	type Alias InvoiceItem
	return json.Marshal(&struct {
		*Alias
	}{
		Alias: (*Alias)(m),
	})
}

// UnmarshalJSON to deserialize []byte
func (m *InvoiceItem) UnmarshalJSON(b []byte) error {
	if string(b) == "" || string(b) == "null" {
		return nil
	}

	type Alias InvoiceItem
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(m),
	}
	err := json.Unmarshal(b, &aux)
	return err
}

// GormDataType gorm common data type
func (m InvoiceItem) GormDataType() string {
	return "InvoiceItem"
}

// GormDBDataType gorm db data type
func (InvoiceItem) GormDBDataType(db *gorm.DB, field *schema.Field) string {
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

func (jm InvoiceItem) GormValue(ctx context.Context, db *gorm.DB) clause.Expr {
	data, _ := jm.MarshalJSON()
	switch db.Dialector.Name() {
	case "mysql":
		if v, ok := db.Dialector.(*mysql.Dialector); ok && !strings.Contains(v.ServerVersion, "MariaDB") {
			return gorm.Expr("CAST(? AS JSON)", string(data))
		}
	}
	return gorm.Expr("?", string(data))
}
