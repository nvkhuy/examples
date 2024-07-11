package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"

	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type ProductFileUploadInfoMetadata struct {
	TotalSuccessRecords int    `json:"total_success_records"`
	TotalFailedRecords  int    `json:"total_failed_records"`
	TotalRecords        int    `json:"total_records"`
	TimeElapsed         string `json:"time_elapsed"`
}

// Value return json value, implement driver.Valuer interface
func (m ProductFileUploadInfoMetadata) Value() (driver.Value, error) {
	ba, err := json.Marshal(&m)
	return string(ba), err
}

// Scan scan value into jsonb, implements sql.Scanner interface
func (m *ProductFileUploadInfoMetadata) Scan(val interface{}) error {
	var ba []byte
	switch v := val.(type) {
	case []byte:
		ba = v
	case string:
		ba = []byte(v)
	default:
		return errors.New(fmt.Sprint("Failed to unmarshal jsonB value:", val))
	}
	var t ProductFileUploadInfoMetadata
	err := json.Unmarshal(ba, &t)
	*m = ProductFileUploadInfoMetadata(t)
	return err
}

// MarshalJSON to output non base64 encoded []byte
func (m ProductFileUploadInfoMetadata) MarshalJSON() ([]byte, error) {
	type Alias ProductFileUploadInfoMetadata
	var copy = new(ProductFileUploadInfoMetadata)
	*copy = m

	return json.Marshal(&struct {
		*Alias
	}{
		Alias: (*Alias)(&m),
	})
}

// UnmarshalJSON to deserialize []byte
func (m *ProductFileUploadInfoMetadata) UnmarshalJSON(b []byte) error {
	if string(b) == "" || string(b) == "null" {
		return nil
	}

	type Alias ProductFileUploadInfoMetadata
	var copy = &struct {
		*Alias
	}{
		Alias: (*Alias)(m),
	}
	err := json.Unmarshal(b, &copy)
	if err != nil {
		return err
	}

	*m = ProductFileUploadInfoMetadata(*copy.Alias)
	return err
}

// GormDataType gorm common data type
func (m ProductFileUploadInfoMetadata) GormDataType() string {
	return "ProductFileUploadInfoMetadata"
}

// GormDBDataType gorm db data type
func (ProductFileUploadInfoMetadata) GormDBDataType(db *gorm.DB, field *schema.Field) string {
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
