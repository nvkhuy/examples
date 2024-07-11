package models

import (
	"context"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/engineeringinflow/inflow-backend/pkg/models/enums"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"
)

type PoRawMaterialMeta struct {
	MaterialType          string                    `json:"material_type" validate:"required"`
	Name                  string                    `json:"name" validate:"required"`
	ColorName             string                    `json:"color_name,omitempty"`
	Weight                *float64                  `json:"weight,omitempty"`
	Description           string                    `json:"description,omitempty"`
	Status                enums.PoRawMaterialStatus `json:"status,omitempty"`
	Attachments           Attachments               `json:"attachments,omitempty"`
	BuyerApproved         *bool                     `json:"buyer_approved,omitempty"`
	ReferenceID           string                    `json:"reference_id,omitempty"`
	WaitingForSendToBuyer *bool                     `json:"waiting_for_send_to_buyer,omitempty"`
}

// Value return json value, implement driver.Valuer interface
func (m PoRawMaterialMeta) Value() (driver.Value, error) {
	ba, err := m.MarshalJSON()
	return string(ba), err
}

// Scan scan value into Jsonb, implements sql.Scanner interface
func (m *PoRawMaterialMeta) Scan(val interface{}) error {
	if val == nil {
		*m = *new(PoRawMaterialMeta)
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
	t := PoRawMaterialMeta{}
	err := json.Unmarshal(ba, &t)
	*m = t
	return err
}

// MarshalJSON to output non base64 encoded []byte
func (m *PoRawMaterialMeta) MarshalJSON() ([]byte, error) {
	type Alias PoRawMaterialMeta
	return json.Marshal(&struct {
		*Alias
	}{
		Alias: (*Alias)(m),
	})
}

// UnmarshalJSON to deserialize []byte
func (m *PoRawMaterialMeta) UnmarshalJSON(b []byte) error {
	if string(b) == "" || string(b) == "null" {
		return nil
	}

	type Alias PoRawMaterialMeta
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(m),
	}
	err := json.Unmarshal(b, &aux)
	return err
}

// GormDataType gorm common data type
func (m PoRawMaterialMeta) GormDataType() string {
	return "PoRawMaterialMeta"
}

// GormDBDataType gorm db data type
func (PoRawMaterialMeta) GormDBDataType(db *gorm.DB, field *schema.Field) string {
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

func (jm PoRawMaterialMeta) GormValue(ctx context.Context, db *gorm.DB) clause.Expr {
	data, _ := jm.MarshalJSON()
	switch db.Dialector.Name() {
	case "mysql":
		if v, ok := db.Dialector.(*mysql.Dialector); ok && !strings.Contains(v.ServerVersion, "MariaDB") {
			return gorm.Expr("CAST(? AS JSON)", string(data))
		}
	}
	return gorm.Expr("?", string(data))
}
