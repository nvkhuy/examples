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

type InvoiceParty struct {
	ID          string            `json:"id,omitempty"`
	Name        string            `json:"name,omitempty"`
	Email       string            `json:"email,omitempty"`
	Address     string            `json:"address,omitempty"`
	PhoneNumber string            `json:"phone_number,omitempty"`
	UEN         string            `json:"uen,omitempty"`
	CompanyName string            `json:"company_name,omitempty"`
	ContactName string            `json:"contact_name,omitempty"`
	CountryCode enums.CountryCode `json:"country_code,omitempty"`
	TaxID       string            `json:"tax_id,omitempty"`
}

var DefaultVendorForOnlinePayment = &InvoiceParty{
	Name:        "INFLOW GLOBAL PTE. LTD",
	Email:       "khanhle@joininflow.io",
	Address:     "4010 ANG MO KIO AVENUE 10 #07-10 TECHPLACE 1 SINGAPORE (569626)",
	PhoneNumber: "+84 (090) 370 2448",
	ContactName: "Khanh Le",
	CompanyName: "INFLOW GLOBAL PTE. LTD",
	CountryCode: enums.CountryCodeSG,
	UEN:         "202221686D",
}

var DefaultVendorForLocal = &InvoiceParty{
	Name:        "INFLOW COMPANY LIMITED",
	Email:       "khanhle@joininflow.io",
	Address:     "48 HUYNH MAN DAT, WARD 19, BINH THANH DISTRICT, HO CHI MINH CITY, VIETNAM",
	PhoneNumber: "+84 (876) 543 2198",
	ContactName: "Khanh Le",
	CompanyName: "INFLOW COMPANY LIMITED",
	CountryCode: enums.CountryCodeVN,
	TaxID:       "0317398566",
}

// Value return json value, implement driver.Valuer interface
func (m InvoiceParty) Value() (driver.Value, error) {
	ba, err := m.MarshalJSON()
	return string(ba), err
}

// Scan scan value into Jsonb, implements sql.Scanner interface
func (m *InvoiceParty) Scan(val interface{}) error {
	if val == nil {
		*m = *new(InvoiceParty)
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
	t := InvoiceParty{}
	err := json.Unmarshal(ba, &t)
	*m = t
	return err
}

// MarshalJSON to output non base64 encoded []byte
func (m *InvoiceParty) MarshalJSON() ([]byte, error) {
	type Alias InvoiceParty
	return json.Marshal(&struct {
		*Alias
	}{
		Alias: (*Alias)(m),
	})
}

// UnmarshalJSON to deserialize []byte
func (m *InvoiceParty) UnmarshalJSON(b []byte) error {
	if string(b) == "" || string(b) == "null" {
		return nil
	}

	type Alias InvoiceParty
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(m),
	}
	err := json.Unmarshal(b, &aux)
	return err
}

// GormDataType gorm common data type
func (m InvoiceParty) GormDataType() string {
	return "InvoiceParty"
}

// GormDBDataType gorm db data type
func (InvoiceParty) GormDBDataType(db *gorm.DB, field *schema.Field) string {
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

func (jm InvoiceParty) GormValue(ctx context.Context, db *gorm.DB) clause.Expr {
	data, _ := jm.MarshalJSON()
	switch db.Dialector.Name() {
	case "mysql":
		if v, ok := db.Dialector.(*mysql.Dialector); ok && !strings.Contains(v.ServerVersion, "MariaDB") {
			return gorm.Expr("CAST(? AS JSON)", string(data))
		}
	}
	return gorm.Expr("?", string(data))
}
