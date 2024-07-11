package models

import (
	"database/sql/driver"
	"encoding/json"

	"errors"
	"fmt"

	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type PaymentTransactionMetadata struct {
	InquiryID          string `json:"inquiry_id,omitempty"`
	InquiryReferenceID string `json:"inquiry_reference_id,omitempty"`

	PurchaseOrderID          string `json:"purchase_order_id,omitempty"`
	PurchaseOrderReferenceID string `json:"purchase_order_reference_id,omitempty"`

	InquiryIDs          []string `json:"inquiry_ids,omitempty"`
	InquiryReferenceIDs []string `json:"inquiry_reference_ids,omitempty"`

	PurchaseOrderIDs          []string `json:"purchase_order_ids,omitempty"`
	PurchaseOrderReferenceIDs []string `json:"purchase_order_reference_ids,omitempty"`

	BulkPurchaseOrderID          string `json:"bulk_purchase_order_id,omitempty"`
	BulkPurchaseOrderReferenceID string `json:"bulk_purchase_order_reference_id,omitempty"`

	BulkPurchaseOrderIDs          []string `json:"bulk_purchase_order_ids,omitempty"`
	BulkPurchaseOrderReferenceIDs []string `json:"bulk_purchase_order_reference_ids,omitempty"`
}

// Value return json value, implement driver.Valuer interface
func (m PaymentTransactionMetadata) Value() (driver.Value, error) {
	ba, err := json.Marshal(&m)
	return string(ba), err
}

// Scan scan value into jsonb, implements sql.Scanner interface
func (m *PaymentTransactionMetadata) Scan(val interface{}) error {
	var ba []byte
	switch v := val.(type) {
	case []byte:
		ba = v
	case string:
		ba = []byte(v)
	default:
		return errors.New(fmt.Sprint("Failed to unmarshal jsonB value:", val))
	}
	var t PaymentTransactionMetadata
	err := json.Unmarshal(ba, &t)
	*m = PaymentTransactionMetadata(t)
	return err
}

// MarshalJSON to output non base64 encoded []byte
func (m PaymentTransactionMetadata) MarshalJSON() ([]byte, error) {
	type Alias PaymentTransactionMetadata
	var copy = new(PaymentTransactionMetadata)
	*copy = m

	return json.Marshal(&struct {
		*Alias
	}{
		Alias: (*Alias)(&m),
	})
}

// UnmarshalJSON to deserialize []byte
func (m *PaymentTransactionMetadata) UnmarshalJSON(b []byte) error {
	if string(b) == "" || string(b) == "null" {
		return nil
	}

	type Alias PaymentTransactionMetadata
	var copy = &struct {
		*Alias
	}{
		Alias: (*Alias)(m),
	}
	err := json.Unmarshal(b, &copy)
	if err != nil {
		return err
	}

	*m = PaymentTransactionMetadata(*copy.Alias)
	return err
}

// GormDataType gorm common data type
func (m PaymentTransactionMetadata) GormDataType() string {
	return "PaymentTransactionMetadata"
}

// GormDBDataType gorm db data type
func (PaymentTransactionMetadata) GormDBDataType(db *gorm.DB, field *schema.Field) string {
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
