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

type PoLogisticMeta struct {
	CourierName    string `json:"courier_name,omitempty" validate:"required"`
	TrackingNumber string `json:"tracking_number,omitempty" validate:"required"`
	TrackingUrl    string `json:"tracking_url,omitempty"`

	EstDeliveredAt       *int64 `json:"est_delivered_at,omitempty"`
	EstDeliveredDuration *int64 `json:"est_delivered_duration,omitempty"`

	Attachments *Attachments `json:"attachments,omitempty"`
}

// Value return json value, implement driver.Valuer interface
func (m PoLogisticMeta) Value() (driver.Value, error) {
	ba, err := m.MarshalJSON()
	return string(ba), err
}

// Scan scan value into Jsonb, implements sql.Scanner interface
func (m *PoLogisticMeta) Scan(val interface{}) error {
	if val == nil {
		*m = *new(PoLogisticMeta)
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
	t := PoLogisticMeta{}
	err := json.Unmarshal(ba, &t)
	*m = t
	return err
}

// MarshalJSON to output non base64 encoded []byte
func (m *PoLogisticMeta) MarshalJSON() ([]byte, error) {
	type Alias PoLogisticMeta
	return json.Marshal(&struct {
		*Alias
	}{
		Alias: (*Alias)(m),
	})
}

// UnmarshalJSON to deserialize []byte
func (m *PoLogisticMeta) UnmarshalJSON(b []byte) error {
	if string(b) == "" || string(b) == "null" {
		return nil
	}

	type Alias PoLogisticMeta
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(m),
	}
	err := json.Unmarshal(b, &aux)
	return err
}

// GormDataType gorm common data type
func (m PoLogisticMeta) GormDataType() string {
	return "PoLogisticMeta"
}

// GormDBDataType gorm db data type
func (PoLogisticMeta) GormDBDataType(db *gorm.DB, field *schema.Field) string {
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

func (jm PoLogisticMeta) GormValue(ctx context.Context, db *gorm.DB) clause.Expr {
	data, _ := jm.MarshalJSON()
	switch db.Dialector.Name() {
	case "mysql":
		if v, ok := db.Dialector.(*mysql.Dialector); ok && !strings.Contains(v.ServerVersion, "MariaDB") {
			return gorm.Expr("CAST(? AS JSON)", string(data))
		}
	}
	return gorm.Expr("?", string(data))
}
func (jm *PoLogisticMeta) GenerateFileURL() *PoLogisticMeta {
	var values = jm.Attachments.GenerateFileURL()
	jm.Attachments = &values
	return jm
}
