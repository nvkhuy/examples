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

type PoProductionMeta struct {
	Name                   string      `json:"name"`
	Description            string      `json:"description,omitempty"`
	Attachments            Attachments `json:"attachments,omitempty"`
	ProductAttachments     Attachments `json:"product_attachments,omitempty"`
	FinalReadyMaterialDate *int64      `json:"final_ready_material_date,omitempty"`
	ActualCutDate          *int64      `json:"actual_cut_date,omitempty"`
	ActualStartSewingDate  *int64      `json:"actual_start_sewing_date,omitempty"`
	ActualFinalOutputDate  *int64      `json:"actual_final_output_date,omitempty"`
	ActualFinalPackingDate *int64      `json:"actual_final_packing_date,omitempty"`

	ActualInlineDate          *int64 `json:"actual_inline_date,omitempty"`
	ActualWashDate            *int64 `json:"actual_wash_date,omitempty"`
	ActualFinalInspectionDate *int64 `json:"actual_final_inspection_date,omitempty"`
}

// Value return json value, implement driver.Valuer interface
func (m PoProductionMeta) Value() (driver.Value, error) {
	ba, err := m.MarshalJSON()
	return string(ba), err
}

// Scan scan value into Jsonb, implements sql.Scanner interface
func (m *PoProductionMeta) Scan(val interface{}) error {
	if val == nil {
		*m = *new(PoProductionMeta)
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
	t := PoProductionMeta{}
	err := json.Unmarshal(ba, &t)
	*m = t
	return err
}

// MarshalJSON to output non base64 encoded []byte
func (m *PoProductionMeta) MarshalJSON() ([]byte, error) {
	type Alias PoProductionMeta
	return json.Marshal(&struct {
		*Alias
	}{
		Alias: (*Alias)(m),
	})
}

// UnmarshalJSON to deserialize []byte
func (m *PoProductionMeta) UnmarshalJSON(b []byte) error {
	if string(b) == "" || string(b) == "null" {
		return nil
	}

	type Alias PoProductionMeta
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(m),
	}
	err := json.Unmarshal(b, &aux)
	return err
}

// GormDataType gorm common data type
func (m PoProductionMeta) GormDataType() string {
	return "PoProductionMeta"
}

// GormDBDataType gorm db data type
func (PoProductionMeta) GormDBDataType(db *gorm.DB, field *schema.Field) string {
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

func (jm PoProductionMeta) GormValue(ctx context.Context, db *gorm.DB) clause.Expr {
	data, _ := jm.MarshalJSON()
	switch db.Dialector.Name() {
	case "mysql":
		if v, ok := db.Dialector.(*mysql.Dialector); ok && !strings.Contains(v.ServerVersion, "MariaDB") {
			return gorm.Expr("CAST(? AS JSON)", string(data))
		}
	}
	return gorm.Expr("?", string(data))
}
func (jm *PoProductionMeta) GenerateFileURL() *PoProductionMeta {
	jm.Attachments = jm.Attachments.GenerateFileURL()
	jm.ProductAttachments = jm.ProductAttachments.GenerateFileURL()

	return jm
}
