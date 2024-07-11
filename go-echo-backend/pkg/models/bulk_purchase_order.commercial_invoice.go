package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/engineeringinflow/inflow-backend/pkg/models/enums"
	"github.com/engineeringinflow/inflow-backend/pkg/models/price"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type SizeQty map[string]float64

type BulkPurchaseOrderCommercialInvoiceItems []*BulkPurchaseOrderCommercialInvoiceItem
type BulkPurchaseOrderCommercialInvoiceItem struct {
	ID            string      `json:"id,omitempty"`
	Color         string      `json:"color,omitempty"`
	Size          *SizeQty    `json:"size,omitempty"`
	SizeName      string      `json:"size_name"`
	UnitPrice     price.Price `json:"unit_price,omitempty"`
	TotalQuantity int64       `json:"total_quantity,omitempty"`
	TotalAmount   price.Price `json:"total_amount,omitempty"`
	Qty           int64       `json:"qty"`
	ActualQty     int64       `json:"actual_qty"`
	ItemCode      string      `json:"item_code,omitempty"`
	HSCode        string      `json:"hs_code,omitempty"`
}

type BulkPurchaseOrderCommercialInvoice struct {
	InvoiceType   enums.InvoiceType                       `json:"invoice_type,omitempty"`
	InvoiceNumber int                                     `json:"invoice_number,omitempty"`
	DueDate       int64                                   `json:"due_date,omitempty"`
	IssuedDate    int64                                   `json:"issued_date,omitempty"`
	Currency      string                                  `json:"currency,omitempty"`
	CountryCode   string                                  `json:"country_code,omitempty"`
	Items         BulkPurchaseOrderCommercialInvoiceItems `json:"items,omitempty"`
	Vendor        *InvoiceParty                           `json:"vendor,omitempty"`
	Consignee     *InvoiceParty                           `json:"consignee,omitempty"`
	Shipper       *InvoiceParty                           `json:"shipper,omitempty"`
	Status        enums.InvoiceStatus                     `json:"status,omitempty" validate:"omitempty,oneof=paid unpaid"`
	Note          string                                  `json:"note,omitempty"`

	Pricing
}

// Value return json value, implement driver.Valuer interface
func (m BulkPurchaseOrderCommercialInvoice) Value() (driver.Value, error) {
	ba, err := json.Marshal(&m)
	return string(ba), err
}

// Scan scan value into jsonb, implements sql.Scanner interface
func (m *BulkPurchaseOrderCommercialInvoice) Scan(val interface{}) error {
	var ba []byte
	switch v := val.(type) {
	case []byte:
		ba = v
	case string:
		ba = []byte(v)
	default:
		return errors.New(fmt.Sprint("Failed to unmarshal jsonB value:", val))
	}
	var t BulkPurchaseOrderCommercialInvoice
	err := json.Unmarshal(ba, &t)
	*m = BulkPurchaseOrderCommercialInvoice(t)
	return err
}

// MarshalJSON to output non base64 encoded []byte
func (m BulkPurchaseOrderCommercialInvoice) MarshalJSON() ([]byte, error) {
	type Alias BulkPurchaseOrderCommercialInvoice
	var copy = new(BulkPurchaseOrderCommercialInvoice)
	*copy = m

	return json.Marshal(&struct {
		*Alias
	}{
		Alias: (*Alias)(&m),
	})
}

// UnmarshalJSON to deserialize []byte
func (m *BulkPurchaseOrderCommercialInvoice) UnmarshalJSON(b []byte) error {
	if string(b) == "" || string(b) == "null" {
		return nil
	}

	type Alias BulkPurchaseOrderCommercialInvoice
	var copy = &struct {
		*Alias
	}{
		Alias: (*Alias)(m),
	}
	err := json.Unmarshal(b, &copy)
	if err != nil {
		return err
	}

	*m = BulkPurchaseOrderCommercialInvoice(*copy.Alias)
	return err
}

// GormDataType gorm common data type
func (m BulkPurchaseOrderCommercialInvoice) GormDataType() string {
	return "BulkPurchaseOrderCommercialInvoice"
}

// GormDBDataType gorm db data type
func (BulkPurchaseOrderCommercialInvoice) GormDBDataType(db *gorm.DB, field *schema.Field) string {
	switch db.Dialector.Name() {
	case "sqlite":
		return "json"
	case "mysql":
		return "json"
	case "postgres":
		return "jsonB"
	}
	return ""
}

func (s SizeQty) GetSizeDescription() string {
	var parts []string
	for size := range s {
		parts = append(parts, size)
	}

	return strings.Join(parts, ",")
}
