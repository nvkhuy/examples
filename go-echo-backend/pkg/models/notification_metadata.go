package models

import (
	"database/sql/driver"
	"encoding/json"

	"errors"
	"fmt"

	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type NotificationMetadata struct {
	InquiryID                string `json:"inquiry_id,omitempty"`
	InquiryReferenceID       string `json:"inquiry_reference_id,omitempty"`
	CommentID                string `json:"comment_id,omitempty"`
	PurchaseOrderID          string `json:"purchase_order_id"`
	PurchaseOrderReferenceID string `json:"purchase_order_reference_id,omitempty"`
	InquirySellerID          string `json:"inquiry_seller_id,omitempty"`
}

// Value return json value, implement driver.Valuer interface
func (m NotificationMetadata) Value() (driver.Value, error) {
	if (m == NotificationMetadata{}) {
		return nil, nil
	}

	ba, err := json.Marshal(&m)
	return string(ba), err
}

// Scan scan value into jsonb, implements sql.Scanner interface
func (m *NotificationMetadata) Scan(val interface{}) error {
	var ba []byte
	switch v := val.(type) {
	case []byte:
		ba = v
	case string:
		ba = []byte(v)
	default:
		return errors.New(fmt.Sprint("Failed to unmarshal jsonB value:", val))
	}
	var t NotificationMetadata
	err := json.Unmarshal(ba, &t)
	*m = NotificationMetadata(t)
	return err
}

// MarshalJSON to output non base64 encoded []byte
func (m NotificationMetadata) MarshalJSON() ([]byte, error) {
	if (m == NotificationMetadata{}) {
		return []byte("null"), nil
	}
	type Alias NotificationMetadata
	var copy = new(NotificationMetadata)
	*copy = m

	return json.Marshal(&struct {
		*Alias
	}{
		Alias: (*Alias)(&m),
	})
}

// UnmarshalJSON to deserialize []byte
func (m *NotificationMetadata) UnmarshalJSON(b []byte) error {
	if string(b) == "" || string(b) == "null" {
		return nil
	}

	type Alias NotificationMetadata
	var copy = &struct {
		*Alias
	}{
		Alias: (*Alias)(m),
	}
	err := json.Unmarshal(b, &copy)
	if err != nil {
		return err
	}

	*m = NotificationMetadata(*copy.Alias)
	return err
}

// GormDataType gorm common data type
func (m NotificationMetadata) GormDataType() string {
	return "NotificationMetadata"
}

// GormDBDataType gorm db data type
func (NotificationMetadata) GormDBDataType(db *gorm.DB, field *schema.Field) string {
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
