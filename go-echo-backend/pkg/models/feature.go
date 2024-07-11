package models

import (
	"database/sql/driver"
	"encoding/json"

	"errors"
	"fmt"

	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type Feature struct {
	Type string `json:"type"`
	Name string `json:"name"`
}

// Value return json value, implement driver.Valuer interface
func (m Feature) Value() (driver.Value, error) {
	if (m == Feature{}) {
		return nil, nil
	}

	ba, err := json.Marshal(&m)
	return string(ba), err
}

// Scan scan value into jsonb, implements sql.Scanner interface
func (m *Feature) Scan(val interface{}) error {
	var ba []byte
	switch v := val.(type) {
	case []byte:
		ba = v
	case string:
		ba = []byte(v)
	default:
		return errors.New(fmt.Sprint("Failed to unmarshal jsonB value:", val))
	}
	var t Feature
	err := json.Unmarshal(ba, &t)
	*m = Feature(t)
	return err
}

// MarshalJSON to output non base64 encoded []byte
func (m Feature) MarshalJSON() ([]byte, error) {
	if (m == Feature{}) {
		return []byte("null"), nil
	}
	type Alias Feature
	var copy = new(Feature)
	*copy = m

	return json.Marshal(&struct {
		*Alias
	}{
		Alias: (*Alias)(&m),
	})
}

// UnmarshalJSON to deserialize []byte
func (m *Feature) UnmarshalJSON(b []byte) error {
	if string(b) == "" || string(b) == "null" {
		return nil
	}

	type Alias Feature
	var copy = &struct {
		*Alias
	}{
		Alias: (*Alias)(m),
	}
	err := json.Unmarshal(b, &copy)
	if err != nil {
		return err
	}

	*m = Feature(*copy.Alias)
	return err
}

// GormDataType gorm common data type
func (m Feature) GormDataType() string {
	return "Feature"
}

// GormDBDataType gorm db data type
func (Feature) GormDBDataType(db *gorm.DB, field *schema.Field) string {
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
