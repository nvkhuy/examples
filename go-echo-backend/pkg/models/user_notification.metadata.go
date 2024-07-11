package models

import (
	"database/sql/driver"
	"encoding/json"

	"errors"
	"fmt"

	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type UserNotificationMetadata struct {
	AdminID string `json:"admin_id,omitempty"`

	CommentID string `json:"comment_id,omitempty"`

	InquiryID          string `json:"inquiry_id,omitempty"`
	InquiryReferenceID string `json:"inquiry_reference_id,omitempty"`

	BulkPurchaseOrderID          string `json:"bulk_purchase_order_id"`
	BulkPurchaseOrderReferenceID string `json:"bulk_order_reference_id,omitempty"`

	PurchaseOrderID          string `json:"purchase_order_id"`
	PurchaseOrderReferenceID string `json:"purchase_order_reference_id,omitempty"`

	InquirySellerID string `json:"inquiry_seller_id,omitempty"`
}

// Value return json value, implement driver.Valuer interface
func (m UserNotificationMetadata) Value() (driver.Value, error) {
	if (m == UserNotificationMetadata{}) {
		return nil, nil
	}

	ba, err := json.Marshal(&m)
	return string(ba), err
}

// Scan scan value into jsonb, implements sql.Scanner interface
func (m *UserNotificationMetadata) Scan(val interface{}) error {
	var ba []byte
	switch v := val.(type) {
	case []byte:
		ba = v
	case string:
		ba = []byte(v)
	default:
		return errors.New(fmt.Sprint("Failed to unmarshal jsonB value:", val))
	}
	var t UserNotificationMetadata
	err := json.Unmarshal(ba, &t)
	*m = UserNotificationMetadata(t)
	return err
}

// MarshalJSON to output non base64 encoded []byte
func (m UserNotificationMetadata) MarshalJSON() ([]byte, error) {
	if (m == UserNotificationMetadata{}) {
		return []byte("null"), nil
	}
	type Alias UserNotificationMetadata
	var copy = new(UserNotificationMetadata)
	*copy = m

	return json.Marshal(&struct {
		*Alias
	}{
		Alias: (*Alias)(&m),
	})
}

// UnmarshalJSON to deserialize []byte
func (m *UserNotificationMetadata) UnmarshalJSON(b []byte) error {
	if string(b) == "" || string(b) == "null" {
		return nil
	}

	type Alias UserNotificationMetadata
	var copy = &struct {
		*Alias
	}{
		Alias: (*Alias)(m),
	}
	err := json.Unmarshal(b, &copy)
	if err != nil {
		return err
	}

	*m = UserNotificationMetadata(*copy.Alias)
	return err
}

// GormDataType gorm common data type
func (m UserNotificationMetadata) GormDataType() string {
	return "UserNotificationMetadata"
}

// GormDBDataType gorm db data type
func (UserNotificationMetadata) GormDBDataType(db *gorm.DB, field *schema.Field) string {
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
