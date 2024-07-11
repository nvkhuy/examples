package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"

	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type PurchaseOrderChanges struct {
	Requirement string `json:"requirement"`

	TrimButton string `json:"trim_button"`
	TrimThread string `json:"trim_thread"`
	TrimZipper string `json:"trim_zipper"`
	TrimLabel  string `json:"trim_label"`

	FabricName   string   `json:"fabric_name"`
	FabricWeight *float64 `json:"fabric_weight"`

	Attachments         *Attachments `json:"attachments,omitempty"`
	FabricAttachments   *Attachments `json:"fabric_attachments,omitempty"`
	TechpackAttachments *Attachments `json:"techpack_attachments,omitempty"`
}

// Value return json value, implement driver.Valuer interface
func (p PurchaseOrderChanges) Value() (driver.Value, error) {
	ba, err := p.MarshalJSON()
	return string(ba), err
}

// Scan scan value into Jsonb, implements sql.Scanner interface
func (p *PurchaseOrderChanges) Scan(val interface{}) error {
	if val == nil {
		*p = *new(PurchaseOrderChanges)
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
	t := PurchaseOrderChanges{}
	err := json.Unmarshal(ba, &t)
	*p = t
	return err
}

// MarshalJSON to output non base64 encoded []byte
func (p *PurchaseOrderChanges) MarshalJSON() ([]byte, error) {
	type Alias PurchaseOrderChanges
	return json.Marshal(&struct {
		*Alias
	}{
		Alias: (*Alias)(p),
	})
}

// UnmarshalJSON to deserialize []byte
func (p *PurchaseOrderChanges) UnmarshalJSON(b []byte) error {
	if string(b) == "" || string(b) == "null" {
		return nil
	}

	type Alias PurchaseOrderChanges
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(p),
	}
	err := json.Unmarshal(b, &aux)
	return err
}

// GormDataType gorm common data type
func (p PurchaseOrderChanges) GormDataType() string {
	return "PurchaseOrderChanges"
}

// GormDBDataType gorm db data type
func (PurchaseOrderChanges) GormDBDataType(db *gorm.DB, field *schema.Field) string {
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
