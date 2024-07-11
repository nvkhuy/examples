package models

import (
	"database/sql/driver"
	"encoding/json"

	"errors"
	"fmt"

	"github.com/engineeringinflow/inflow-backend/pkg/models/enums"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type InvoiceMetadata struct {
	InvoiceType        enums.InvoiceType `json:"invoice_type,omitempty"`
	InquiryID          string            `json:"inquiry_id,omitempty"`
	InquiryReferenceID string            `json:"inquiry_reference_id,omitempty"`

	PurchaseOrderID          string `json:"purchase_order_id,omitempty"`
	PurchaseOrderReferenceID string `json:"purchase_order_reference_id,omitempty"`

	InquiryIDs          []string `json:"inquiry_ids,omitempty"`
	InquiryReferenceIDs []string `json:"inquiry_reference_ids,omitempty"`

	PurchaseOrderIDs          []string `json:"purchase_order_ids,omitempty"`
	PurchaseOrderReferenceIDs []string `json:"purchase_order_reference_ids,omitempty"`

	BulkPurchaseOrderID          string   `json:"bulk_purchase_order_id,omitempty"`
	BulkPurchaseOrderIDs         []string `json:"bulk_purchase_order_ids,omitempty"`
	BulkPurchaseOrderReferenceID string   `json:"bulk_purchase_order_reference_id,omitempty"`

	CheckoutSessionID string `json:"checkout_session_id,omitempty"`

	BulkPurchaseOrderOrderClosingDocAttachment   *Attachment `json:"bulk_purchase_order_order_closing_doc_attachment,omitempty"`
	BulkPurchaseOrderCommercialInvoiceAttachment *Attachment `json:"bulk_purchase_order_commercial_invoice_attachment,omitempty"`

	BulkEstimatedProductionLeadTime *int64 `json:"bulk_estimated_production_lead_time,omitempty"`
}

// Value return json value, implement driver.Valuer interface
func (m InvoiceMetadata) Value() (driver.Value, error) {
	ba, err := json.Marshal(&m)
	return string(ba), err
}

// Scan scan value into jsonb, implements sql.Scanner interface
func (m *InvoiceMetadata) Scan(val interface{}) error {
	var ba []byte
	switch v := val.(type) {
	case []byte:
		ba = v
	case string:
		ba = []byte(v)
	default:
		return errors.New(fmt.Sprint("Failed to unmarshal jsonB value:", val))
	}
	var t InvoiceMetadata
	err := json.Unmarshal(ba, &t)
	*m = InvoiceMetadata(t)
	return err
}

// MarshalJSON to output non base64 encoded []byte
func (m InvoiceMetadata) MarshalJSON() ([]byte, error) {
	type Alias InvoiceMetadata
	var copy = new(InvoiceMetadata)
	*copy = m

	return json.Marshal(&struct {
		*Alias
	}{
		Alias: (*Alias)(&m),
	})
}

// UnmarshalJSON to deserialize []byte
func (m *InvoiceMetadata) UnmarshalJSON(b []byte) error {
	if string(b) == "" || string(b) == "null" {
		return nil
	}

	type Alias InvoiceMetadata
	var copy = &struct {
		*Alias
	}{
		Alias: (*Alias)(m),
	}
	err := json.Unmarshal(b, &copy)
	if err != nil {
		return err
	}

	*m = InvoiceMetadata(*copy.Alias)
	return err
}

// GormDataType gorm common data type
func (m InvoiceMetadata) GormDataType() string {
	return "InvoiceMetadata"
}

// GormDBDataType gorm db data type
func (InvoiceMetadata) GormDBDataType(db *gorm.DB, field *schema.Field) string {
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
