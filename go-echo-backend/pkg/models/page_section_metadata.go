package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"

	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type PageSectionMetadata struct {
	Href            string           `json:"href"`
	CategoriesHref  []CategoryHref   `json:"categories_href"`
	ProductsHref    []ProductHref    `json:"products_href"`
	CollectionsHref []CollectionHref `json:"collections_href"`
}

type CategoryHref struct {
	CategoryID string `json:"category_id"`
	Href       string `json:"href"`
}

type ProductHref struct {
	ProductID string `json:"product_id"`
	Href      string `json:"href"`
}

type CollectionHref struct {
	CollectionID string `json:"collection_id"`
	Href         string `json:"href"`
}

// Value return json value, implement driver.Valuer interface
func (m PageSectionMetadata) Value() (driver.Value, error) {
	ba, err := json.Marshal(&m)
	return string(ba), err
}

// Scan scan value into jsonb, implements sql.Scanner interface
func (m *PageSectionMetadata) Scan(val interface{}) error {
	var ba []byte
	switch v := val.(type) {
	case []byte:
		ba = v
	case string:
		ba = []byte(v)
	default:
		return errors.New(fmt.Sprint("Failed to unmarshal jsonB value:", val))
	}
	var t PageSectionMetadata
	err := json.Unmarshal(ba, &t)
	*m = PageSectionMetadata(t)
	return err
}

// MarshalJSON to output non base64 encoded []byte
func (m PageSectionMetadata) MarshalJSON() ([]byte, error) {
	type Alias PageSectionMetadata
	var copy = new(PageSectionMetadata)
	*copy = m

	return json.Marshal(&struct {
		*Alias
	}{
		Alias: (*Alias)(&m),
	})
}

// UnmarshalJSON to deserialize []byte
func (m *PageSectionMetadata) UnmarshalJSON(b []byte) error {
	if string(b) == "" || string(b) == "null" {
		return nil
	}

	type Alias PageSectionMetadata
	var copy = &struct {
		*Alias
	}{
		Alias: (*Alias)(m),
	}
	err := json.Unmarshal(b, &copy)
	if err != nil {
		return err
	}

	*m = PageSectionMetadata(*copy.Alias)
	return err
}

// GormDataType gorm common data type
func (m PageSectionMetadata) GormDataType() string {
	return "PageSectionMetadata"
}

// GormDBDataType gorm db data type
func (PageSectionMetadata) GormDBDataType(db *gorm.DB, field *schema.Field) string {
	switch db.Dialector.Name() {
	case "sqlite":
		return "json"
	case "mysql":
		return "json"
	case "postgres":
		return "jsonB"
	}
	return ""
}
