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

type InquiryApproveRejectMeta struct {
	Comment     string      `json:"comment,omitempty"`
	Reason      string      `json:"reason"`
	Attachments Attachments `json:"attachments"`
	CreatedAt   int64       `gorm:"autoCreateTime" json:"created_at,omitempty"`
}

// Value return json value, implement driver.Valuer interface
func (m InquiryApproveRejectMeta) Value() (driver.Value, error) {
	ba, err := m.MarshalJSON()
	return string(ba), err
}

// Scan scan value into Jsonb, implements sql.Scanner interface
func (m *InquiryApproveRejectMeta) Scan(val interface{}) error {
	if val == nil {
		*m = *new(InquiryApproveRejectMeta)
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
	t := InquiryApproveRejectMeta{}
	err := json.Unmarshal(ba, &t)
	*m = t
	return err
}

// MarshalJSON to output non base64 encoded []byte
func (m *InquiryApproveRejectMeta) MarshalJSON() ([]byte, error) {
	type Alias InquiryApproveRejectMeta
	return json.Marshal(&struct {
		*Alias
	}{
		Alias: (*Alias)(m),
	})
}

// UnmarshalJSON to deserialize []byte
func (m *InquiryApproveRejectMeta) UnmarshalJSON(b []byte) error {
	if string(b) == "" || string(b) == "null" {
		return nil
	}

	type Alias InquiryApproveRejectMeta
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(m),
	}
	err := json.Unmarshal(b, &aux)
	return err
}

// GormDataType gorm common data type
func (m InquiryApproveRejectMeta) GormDataType() string {
	return "InquiryApproveRejectMeta"
}

// GormDBDataType gorm db data type
func (InquiryApproveRejectMeta) GormDBDataType(db *gorm.DB, field *schema.Field) string {
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

func (jm InquiryApproveRejectMeta) GormValue(ctx context.Context, db *gorm.DB) clause.Expr {
	data, _ := jm.MarshalJSON()
	switch db.Dialector.Name() {
	case "mysql":
		if v, ok := db.Dialector.(*mysql.Dialector); ok && !strings.Contains(v.ServerVersion, "MariaDB") {
			return gorm.Expr("CAST(? AS JSON)", string(data))
		}
	}
	return gorm.Expr("?", string(data))
}
