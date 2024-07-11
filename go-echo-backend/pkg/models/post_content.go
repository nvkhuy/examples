package models

import (
	"context"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/engineeringinflow/inflow-backend/pkg/models/enums"
	"strings"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"
)

type PostContent struct {
	Title            string           `json:"title,omitempty"`
	Status           enums.PostStatus `json:"status,omitempty" gorm:"default:'new'"`
	Content          string           `json:"content,omitempty"`
	ContentURL       string           `json:"content_url,omitempty"`
	ShortDescription string           `json:"short_description,omitempty"`
	Slug             string           `gorm:"not null;unique" json:"slug,omitempty"`
	SettingSEO       *SettingSEO      `gorm:"-" json:"setting_seo,omitempty"`
	SettingSeoID     string           `json:"setting_seo_id"`
}

// Value return json value, implement driver.Valuer interface
func (p PostContent) Value() (driver.Value, error) {
	ba, err := p.MarshalJSON()
	return string(ba), err
}

// Scan scan value into Jsonb, implements sql.Scanner interface
func (p *PostContent) Scan(val interface{}) error {
	if val == nil {
		*p = *new(PostContent)
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
	t := PostContent{}
	err := json.Unmarshal(ba, &t)
	*p = t
	return err
}

// MarshalJSON to output non base64 encoded []byte
func (p *PostContent) MarshalJSON() ([]byte, error) {
	type Alias PostContent
	return json.Marshal(&struct {
		*Alias
	}{
		Alias: (*Alias)(p),
	})
}

// UnmarshalJSON to deserialize []byte
func (p *PostContent) UnmarshalJSON(b []byte) error {
	if string(b) == "" || string(b) == "null" {
		return nil
	}

	type Alias PostContent
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(p),
	}
	err := json.Unmarshal(b, &aux)
	return err
}

// GormDataType gorm common data type
func (p PostContent) GormDataType() string {
	return "PostContent"
}

// GormDBDataType gorm db data type
func (PostContent) GormDBDataType(db *gorm.DB, field *schema.Field) string {
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

func (p PostContent) GormValue(ctx context.Context, db *gorm.DB) clause.Expr {
	data, _ := p.MarshalJSON()
	switch db.Dialector.Name() {
	case "mysql":
		if v, ok := db.Dialector.(*mysql.Dialector); ok && !strings.Contains(v.ServerVersion, "MariaDB") {
			return gorm.Expr("CAST(? AS JSON)", string(data))
		}
	}
	return gorm.Expr("?", string(data))
}
