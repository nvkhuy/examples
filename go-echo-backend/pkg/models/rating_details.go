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

type RatingDetails struct {
	Stars map[string]struct {
		Count   float64 `json:"count,omitempty"`
		Percent float64 `json:"percent,omitempty"`
	} `json:"stars,omitempty"`
	PercentageOverFit *JsonMetaData `json:"percentage_overfit,omitempty"`
	RatingAverage     float64       `json:"rating_average,omitempty"`
	ReviewsCount      int           `json:"reviews_count,omitempty"`
}

// Value return json value, implement driver.Valuer interface
func (m RatingDetails) Value() (driver.Value, error) {
	ba, err := m.MarshalJSON()
	return string(ba), err
}

// Scan scan value into Jsonb, implements sql.Scanner interface
func (m *RatingDetails) Scan(val interface{}) error {
	if val == nil {
		*m = *new(RatingDetails)
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
	t := RatingDetails{}
	err := json.Unmarshal(ba, &t)
	*m = t
	return err
}

// MarshalJSON to output non base64 encoded []byte
func (m *RatingDetails) MarshalJSON() ([]byte, error) {
	type Alias RatingDetails
	return json.Marshal(&struct {
		*Alias
	}{
		Alias: (*Alias)(m),
	})
}

// UnmarshalJSON to deserialize []byte
func (m *RatingDetails) UnmarshalJSON(b []byte) error {
	type Alias RatingDetails
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(m),
	}
	err := json.Unmarshal(b, &aux)
	return err
}

// GormDataType gorm common data type
func (m RatingDetails) GormDataType() string {
	return "RatingDetails"
}

// GormDBDataType gorm db data type
func (RatingDetails) GormDBDataType(db *gorm.DB, field *schema.Field) string {
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

func (jm RatingDetails) GormValue(ctx context.Context, db *gorm.DB) clause.Expr {
	data, _ := jm.MarshalJSON()
	switch db.Dialector.Name() {
	case "mysql":
		if v, ok := db.Dialector.(*mysql.Dialector); ok && !strings.Contains(v.ServerVersion, "MariaDB") {
			return gorm.Expr("CAST(? AS JSON)", string(data))
		}
	}
	return gorm.Expr("?", string(data))
}
