package models

import (
	"context"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"
)

type OBFabricTypeMeta struct {
	FabricName     string       `json:"fabric_name,omitempty"`
	FabricValue    string       `json:"fabric_value,omitempty"`
	FabricImageURL string       `json:"fabric_image_url,omitempty"`
	Attachments    *Attachments `json:"attachments,omitempty"`
	UploadedByUser *bool        `json:"uploaded_by_user,omitempty"`
}

// Value return json value, implement driver.Valuer interface
func (m OBFabricTypeMeta) Value() (driver.Value, error) {
	ba, err := m.MarshalJSON()
	return string(ba), err
}

// Scan scan value into Jsonb, implements sql.Scanner interface
func (m *OBFabricTypeMeta) Scan(val interface{}) error {
	if val == nil {
		*m = *new(OBFabricTypeMeta)
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
	t := OBFabricTypeMeta{}
	err := json.Unmarshal(ba, &t)
	*m = t
	return err
}

// MarshalJSON to output non base64 encoded []byte
func (m *OBFabricTypeMeta) MarshalJSON() ([]byte, error) {
	type Alias OBFabricTypeMeta
	return json.Marshal(&struct {
		*Alias
	}{
		Alias: (*Alias)(m),
	})
}

// UnmarshalJSON to deserialize []byte
func (m *OBFabricTypeMeta) UnmarshalJSON(b []byte) error {
	if string(b) == "" || string(b) == "null" {
		return nil
	}

	type Alias OBFabricTypeMeta
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(m),
	}
	err := json.Unmarshal(b, &aux)
	return err
}

// GormDataType gorm common data type
func (m OBFabricTypeMeta) GormDataType() string {
	return "OBFabricTypeMeta"
}

// GormDBDataType gorm db data type
func (OBFabricTypeMeta) GormDBDataType(db *gorm.DB, field *schema.Field) string {
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

func (jm OBFabricTypeMeta) GormValue(ctx context.Context, db *gorm.DB) clause.Expr {
	data, _ := jm.MarshalJSON()
	switch db.Dialector.Name() {
	case "mysql":
		if v, ok := db.Dialector.(*mysql.Dialector); ok && !strings.Contains(v.ServerVersion, "MariaDB") {
			return gorm.Expr("CAST(? AS JSON)", string(data))
		}
	}
	return gorm.Expr("?", string(data))
}
