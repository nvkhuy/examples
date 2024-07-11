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

type SellerBulkQuotationMOQ struct {
	Type       enums.InquirySellerQuotationType `json:"type" validate:"required"`
	UnitPrice  price.Price                      `json:"unit_price"`
	TotalPrice price.Price                      `json:"total_price"`
	UpCharge   price.Price                      `json:"upcharge" validate:"required"`
	Quantity   *int64                           `json:"quantity" validate:"required"`
	LeadTime   *int64                           `json:"lead_time,omitempty" validate:"required"` // unit: days
}

// Value return json value, implement driver.Valuer interface
func (m SellerBulkQuotationMOQ) Value() (driver.Value, error) {
	ba, err := m.MarshalJSON()
	return string(ba), err
}

// Scan scan value into Jsonb, implements sql.Scanner interface
func (m *SellerBulkQuotationMOQ) Scan(val interface{}) error {
	if val == nil {
		*m = *new(SellerBulkQuotationMOQ)
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
	t := SellerBulkQuotationMOQ{}
	err := json.Unmarshal(ba, &t)
	*m = t
	return err
}

// MarshalJSON to output non base64 encoded []byte
func (m *SellerBulkQuotationMOQ) MarshalJSON() ([]byte, error) {
	type Alias SellerBulkQuotationMOQ
	aux := (*Alias)(m)

	return json.Marshal(&aux)
}

// UnmarshalJSON to deserialize []byte
func (m *SellerBulkQuotationMOQ) UnmarshalJSON(b []byte) error {
	if string(b) == "" || string(b) == "null" {
		return nil
	}

	type Alias SellerBulkQuotationMOQ
	aux := (*Alias)(m)
	err := json.Unmarshal(b, &aux)
	return err
}

// GormDataType gorm common data type
func (m SellerBulkQuotationMOQ) GormDataType() string {
	return "SellerBulkQuotationMOQ"
}

// GormDBDataType gorm db data type
func (SellerBulkQuotationMOQ) GormDBDataType(db *gorm.DB, field *schema.Field) string {
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

func (jm SellerBulkQuotationMOQ) GormValue(ctx context.Context, db *gorm.DB) clause.Expr {
	data, _ := jm.MarshalJSON()

	return gorm.Expr("?", string(data))
}
