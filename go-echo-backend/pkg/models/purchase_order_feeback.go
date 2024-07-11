package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"

	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type PurchaseOrderFeedback struct {
	Visual            *PurchaseOrderFeedbackDetails `json:"visual"`
	Form              *PurchaseOrderFeedbackDetails `json:"form,omitempty"`
	Measurements      *PurchaseOrderFeedbackDetails `json:"measurements,omitempty"`
	Workmanship       *PurchaseOrderFeedbackDetails `json:"workmanship,omitempty"`
	FabricTrimQuality *PurchaseOrderFeedbackDetails `json:"fabric_trim_quality,omitempty"`
	Printing          *PurchaseOrderFeedbackDetails `json:"printing,omitempty"`
	Embroidery        *PurchaseOrderFeedbackDetails `json:"embroidery,omitempty"`
	WashingDyeing     *PurchaseOrderFeedbackDetails `json:"washing_dyeing,omitempty"`
	OtherFeedback     *PurchaseOrderFeedbackDetails `json:"other_feedback,omitempty"`
}

// Value return json value, implement driver.Valuer interface
func (p PurchaseOrderFeedback) Value() (driver.Value, error) {
	ba, err := p.MarshalJSON()
	return string(ba), err
}

// Scan scan value into Jsonb, implements sql.Scanner interface
func (p *PurchaseOrderFeedback) Scan(val interface{}) error {
	if val == nil {
		*p = *new(PurchaseOrderFeedback)
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
	t := PurchaseOrderFeedback{}
	err := json.Unmarshal(ba, &t)
	*p = t
	return err
}

// MarshalJSON to output non base64 encoded []byte
func (p *PurchaseOrderFeedback) MarshalJSON() ([]byte, error) {
	type Alias PurchaseOrderFeedback
	return json.Marshal(&struct {
		*Alias
	}{
		Alias: (*Alias)(p),
	})
}

// UnmarshalJSON to deserialize []byte
func (p *PurchaseOrderFeedback) UnmarshalJSON(b []byte) error {
	if string(b) == "" || string(b) == "null" {
		return nil
	}

	type Alias PurchaseOrderFeedback
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(p),
	}
	err := json.Unmarshal(b, &aux)
	return err
}

// GormDataType gorm common data type
func (p PurchaseOrderFeedback) GormDataType() string {
	return "PurchaseOrderFeedback"
}

// GormDBDataType gorm db data type
func (PurchaseOrderFeedback) GormDBDataType(db *gorm.DB, field *schema.Field) string {
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

type PurchaseOrderFeedbackDetails struct {
	MeetExpectation  bool   `json:"meet_expectation"`
	NeedsImprovement string `json:"needs_improvement"`
}
