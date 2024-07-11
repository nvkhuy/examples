package models

import (
	"context"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/engineeringinflow/inflow-backend/pkg/models/price"
	"github.com/lib/pq"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"
)

type Fabric struct {
	Model
	ReferenceID  string  `gorm:"unique" json:"reference_id"`
	FabricType   string  `gorm:"size:100" json:"fabric_type,omitempty"`
	FabricID     string  `gorm:"size:100" json:"fabric_id,omitempty"`
	FabricWeight float64 `json:"fabric_weight,omitempty"`
	MOQ          float64 `json:"moq"` //minimum order quantity

	Attachments       *Attachments       `json:"attachments"`
	ManufacturerIDs   pq.StringArray     `gorm:"type:varchar(200)[]" json:"manufacturer_ids,omitempty"`
	FabricCostings    *FabricCostings    `json:"fabric_costings,omitempty"`
	Slug              string             `gorm:"unique;size:200" json:"slug,omitempty"`
	FabricCollections []FabricCollection `gorm:"-" json:"fabric_collections,omitempty"`
	Manufacturers     []User             `gorm:"-" json:"manufacturers,omitempty"`

	Title       string         `gorm:"size:2000" json:"title"`
	Description string         `gorm:"size:2000" json:"description,omitempty"`
	Colors      pq.StringArray `gorm:"type:varchar(200)[]" json:"colors"`
	VI          *FabricContent `json:"vi,omitempty"`
}

type FabricCostings []*FabricCosting

type FabricCosting struct {
	From           int64       `json:"from,omitempty"`
	To             int64       `json:"to,omitempty"`
	Price          price.Price `json:"price,omitempty"`
	ProcessingTime string      `json:"processing_time,omitempty"`
}

// Value return json value, implement driver.Valuer interface
func (m FabricCostings) Value() (driver.Value, error) {
	ba, err := m.MarshalJSON()
	return string(ba), err
}

// Scan scan value into Jsonb, implements sql.Scanner interface
func (m *FabricCostings) Scan(val interface{}) error {
	if val == nil {
		*m = *new(FabricCostings)
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
	t := FabricCostings{}
	err := json.Unmarshal(ba, &t)
	*m = t
	return err
}

// MarshalJSON to output non base64 encoded []byte
func (m *FabricCostings) MarshalJSON() ([]byte, error) {
	type Alias FabricCostings
	aux := (*Alias)(m)

	return json.Marshal(&aux)
}

// UnmarshalJSON to deserialize []byte
func (m *FabricCostings) UnmarshalJSON(b []byte) error {
	if string(b) == "" || string(b) == "null" {
		return nil
	}

	type Alias FabricCostings
	aux := (*Alias)(m)
	err := json.Unmarshal(b, &aux)
	return err
}

// GormDataType gorm common data type
func (m FabricCostings) GormDataType() string {
	return "FabricCostings"
}

// GormDBDataType gorm db data type
func (FabricCostings) GormDBDataType(db *gorm.DB, field *schema.Field) string {
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

func (m FabricCostings) GormValue(ctx context.Context, db *gorm.DB) clause.Expr {
	data, _ := m.MarshalJSON()
	switch db.Dialector.Name() {
	case "mysql":
		if v, ok := db.Dialector.(*mysql.Dialector); ok && !strings.Contains(v.ServerVersion, "MariaDB") {
			return gorm.Expr("CAST(? AS JSON)", string(data))
		}
	}
	return gorm.Expr("?", string(data))
}
