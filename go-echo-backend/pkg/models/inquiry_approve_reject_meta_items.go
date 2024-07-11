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

type InquiryApproveRejectMetaItems []*InquiryApproveRejectMeta

// Value return json value, implement driver.Valuer interface
func (m InquiryApproveRejectMetaItems) Value() (driver.Value, error) {
	ba, err := m.MarshalJSON()
	return string(ba), err
}

// Scan scan value into Jsonb, implements sql.Scanner interface
func (m *InquiryApproveRejectMetaItems) Scan(val interface{}) error {
	if val == nil {
		*m = *new(InquiryApproveRejectMetaItems)
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
	t := InquiryApproveRejectMetaItems{}
	err := json.Unmarshal(ba, &t)
	*m = t
	return err
}

// MarshalJSON to output non base64 encoded []byte
func (m *InquiryApproveRejectMetaItems) MarshalJSON() ([]byte, error) {
	type Alias InquiryApproveRejectMetaItems
	aux := (*Alias)(m)

	return json.Marshal(&aux)
}

// UnmarshalJSON to deserialize []byte
func (m *InquiryApproveRejectMetaItems) UnmarshalJSON(b []byte) error {
	if string(b) == "" || string(b) == "null" {
		return nil
	}

	type Alias InquiryApproveRejectMetaItems
	aux := (*Alias)(m)
	err := json.Unmarshal(b, &aux)
	return err
}

// GormDataType gorm common data type
func (m InquiryApproveRejectMetaItems) GormDataType() string {
	return "InquiryApproveRejectMetaItems"
}

// GormDBDataType gorm db data type
func (InquiryApproveRejectMetaItems) GormDBDataType(db *gorm.DB, field *schema.Field) string {
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

func (jm InquiryApproveRejectMetaItems) GormValue(ctx context.Context, db *gorm.DB) clause.Expr {
	data, _ := jm.MarshalJSON()

	return gorm.Expr("?", string(data))
}
