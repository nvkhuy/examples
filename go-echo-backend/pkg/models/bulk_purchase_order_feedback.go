package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type BulkPurchaseOrderFeedback struct {
	Visual            *BulkPurchaseOrderFeedbackDetails `json:"visual"`
	Form              *BulkPurchaseOrderFeedbackDetails `json:"form,omitempty"`
	Measurements      *BulkPurchaseOrderFeedbackDetails `json:"measurements,omitempty"`
	Workmanship       *BulkPurchaseOrderFeedbackDetails `json:"workmanship,omitempty"`
	FabricTrimQuality *BulkPurchaseOrderFeedbackDetails `json:"fabric_trim_quality,omitempty"`
	Printing          *BulkPurchaseOrderFeedbackDetails `json:"printing,omitempty"`
	Embroidery        *BulkPurchaseOrderFeedbackDetails `json:"embroidery,omitempty"`
	WashingDyeing     *BulkPurchaseOrderFeedbackDetails `json:"washing_dyeing,omitempty"`
	OtherFeedback     *BulkPurchaseOrderFeedbackDetails `json:"other_feedback,omitempty"`
}

// Value return json value, implement driver.Valuer interface
func (p BulkPurchaseOrderFeedback) Value() (driver.Value, error) {
	ba, err := p.MarshalJSON()
	return string(ba), err
}

// Scan scan value into Jsonb, implements sql.Scanner interface
func (p *BulkPurchaseOrderFeedback) Scan(val interface{}) error {
	if val == nil {
		*p = *new(BulkPurchaseOrderFeedback)
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
	t := BulkPurchaseOrderFeedback{}
	err := json.Unmarshal(ba, &t)
	*p = t
	return err
}

// MarshalJSON to output non base64 encoded []byte
func (p *BulkPurchaseOrderFeedback) MarshalJSON() ([]byte, error) {
	type Alias BulkPurchaseOrderFeedback
	return json.Marshal(&struct {
		*Alias
	}{
		Alias: (*Alias)(p),
	})
}

// UnmarshalJSON to deserialize []byte
func (p *BulkPurchaseOrderFeedback) UnmarshalJSON(b []byte) error {
	if string(b) == "" || string(b) == "null" {
		return nil
	}

	type Alias BulkPurchaseOrderFeedback
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(p),
	}
	err := json.Unmarshal(b, &aux)
	return err
}

// GormDataType gorm common data type
func (p BulkPurchaseOrderFeedback) GormDataType() string {
	return "BulkPurchaseOrderFeedback"
}

// GormDBDataType gorm db data type
func (BulkPurchaseOrderFeedback) GormDBDataType(db *gorm.DB, field *schema.Field) string {
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

type BulkPurchaseOrderFeedbackDetails struct {
	MeetExpectation  bool   `json:"meet_expectation"`
	NeedsImprovement string `json:"needs_improvement"`
}
