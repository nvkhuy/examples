package models

import (
	"context"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/engineeringinflow/inflow-backend/pkg/models/enums"
	"github.com/engineeringinflow/inflow-backend/pkg/models/price"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"
)

type InquiryQuotationItem struct {
	Type enums.InquiryType `json:"type,omitempty" validate:"required"`

	Sku      string      `json:"sku,omitempty"`
	Style    string      `json:"style,omitempty"`
	Price    price.Price `json:"price,omitempty" validate:"required"`
	Quantity *int64      `json:"quantity,omitempty" validate:"required"`
	LeadTime *int64      `json:"lead_time,omitempty" validate:"required"` // time when admin send quotation, unit: days
	Accepted *bool       `gorm:"default:false"  json:"accepted,omitempty"`
	Note     string      `json:"note,omitempty"`
}

// Value return json value, implement driver.Valuer interface
func (m InquiryQuotationItem) Value() (driver.Value, error) {
	ba, err := m.MarshalJSON()
	return string(ba), err
}

// Scan scan value into Jsonb, implements sql.Scanner interface
func (m *InquiryQuotationItem) Scan(val interface{}) error {
	if val == nil {
		*m = *new(InquiryQuotationItem)
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
	t := InquiryQuotationItem{}
	err := json.Unmarshal(ba, &t)
	*m = t
	return err
}

// MarshalJSON to output non base64 encoded []byte
func (m *InquiryQuotationItem) MarshalJSON() ([]byte, error) {
	type Alias InquiryQuotationItem
	aux := (*Alias)(m)

	return json.Marshal(&aux)
}

// UnmarshalJSON to deserialize []byte
func (m *InquiryQuotationItem) UnmarshalJSON(b []byte) error {
	if string(b) == "" || string(b) == "null" {
		return nil
	}

	type Alias InquiryQuotationItem
	aux := (*Alias)(m)
	err := json.Unmarshal(b, &aux)
	return err
}

// GormDataType gorm common data type
func (m InquiryQuotationItem) GormDataType() string {
	return "InquiryQuotationItem"
}

// GormDBDataType gorm db data type
func (InquiryQuotationItem) GormDBDataType(db *gorm.DB, field *schema.Field) string {
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

func (jm InquiryQuotationItem) GormValue(ctx context.Context, db *gorm.DB) clause.Expr {
	data, _ := jm.MarshalJSON()

	return gorm.Expr("?", string(data))
}
