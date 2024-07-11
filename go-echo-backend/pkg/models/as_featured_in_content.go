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

type AsFeaturedInContent struct {
	Slug        string           `json:"slug" gorm:"unique"`
	Title       string           `json:"title,omitempty"`
	Image       *Attachment      `json:"image,omitempty"`
	Link        string           `json:"link,omitempty"`
	Logo        *Attachment      `json:"logo,omitempty"`
	PublishedAt *int64           `json:"published_at,omitempty"`
	Status      enums.PostStatus `json:"status,omitempty" gorm:"default:'new'"`
}

// Value return json value, implement driver.Valuer interface
func (p AsFeaturedInContent) Value() (driver.Value, error) {
	ba, err := p.MarshalJSON()
	return string(ba), err
}

// Scan scan value into Jsonb, implements sql.Scanner interface
func (p *AsFeaturedInContent) Scan(val interface{}) error {
	if val == nil {
		*p = *new(AsFeaturedInContent)
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
	t := AsFeaturedInContent{}
	err := json.Unmarshal(ba, &t)
	*p = t
	return err
}

// MarshalJSON to output non base64 encoded []byte
func (p *AsFeaturedInContent) MarshalJSON() ([]byte, error) {
	type Alias AsFeaturedInContent
	return json.Marshal(&struct {
		*Alias
	}{
		Alias: (*Alias)(p),
	})
}

// UnmarshalJSON to deserialize []byte
func (p *AsFeaturedInContent) UnmarshalJSON(b []byte) error {
	if string(b) == "" || string(b) == "null" {
		return nil
	}

	type Alias AsFeaturedInContent
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(p),
	}
	err := json.Unmarshal(b, &aux)
	return err
}

// GormDataType gorm common data type
func (p AsFeaturedInContent) GormDataType() string {
	return "AsFeaturedInContent"
}

// GormDBDataType gorm db data type
func (AsFeaturedInContent) GormDBDataType(db *gorm.DB, field *schema.Field) string {
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

func (p AsFeaturedInContent) GormValue(ctx context.Context, db *gorm.DB) clause.Expr {
	data, _ := p.MarshalJSON()
	switch db.Dialector.Name() {
	case "mysql":
		if v, ok := db.Dialector.(*mysql.Dialector); ok && !strings.Contains(v.ServerVersion, "MariaDB") {
			return gorm.Expr("CAST(? AS JSON)", string(data))
		}
	}
	return gorm.Expr("?", string(data))
}
