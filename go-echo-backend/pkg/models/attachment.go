package models

import (
	"context"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/engineeringinflow/inflow-backend/pkg/config"
	"github.com/engineeringinflow/inflow-backend/pkg/downloader"
	"github.com/engineeringinflow/inflow-backend/pkg/helper"
	"github.com/engineeringinflow/inflow-backend/pkg/s3"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"
)

type Attachment struct {
	ContentType     string                 `json:"content_type,omitempty"`
	FileURL         string                 `json:"file_url,omitempty" validate:"omitempty,startswith=http"`
	FileKey         string                 `json:"file_key,omitempty" validate:"required"`
	FileName        string                 `json:"file_name,omitempty"`
	FileTitle       string                 `json:"file_title,omitempty"`
	FileDescription string                 `json:"file_description,omitempty"`
	ThumbnailURL    string                 `json:"thumbnail_url,omitempty" validate:"omitempty,startswith=http"`
	Metadata        map[string]interface{} `json:"metadata,omitempty"`
	Blurhash        *s3.Blurhash           `json:"blurhash,omitempty"`
}

// Value return json value, implement driver.Valuer interface
func (m Attachment) Value() (driver.Value, error) {
	ba, err := m.MarshalJSON()
	return string(ba), err
}

// Scan scan value into Jsonb, implements sql.Scanner interface
func (m *Attachment) Scan(val interface{}) error {
	if val == nil {
		*m = *new(Attachment)
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
	t := Attachment{}
	err := json.Unmarshal(ba, &t)
	*m = t
	return err
}

// MarshalJSON to output non base64 encoded []byte
func (m *Attachment) MarshalJSON() ([]byte, error) {
	type Alias Attachment
	return json.Marshal(&struct {
		*Alias
	}{
		Alias: (*Alias)(m),
	})
}

// UnmarshalJSON to deserialize []byte
func (m *Attachment) UnmarshalJSON(b []byte) error {
	if string(b) == "" || string(b) == "null" {
		return nil
	}

	type Alias Attachment
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(m),
	}
	err := json.Unmarshal(b, &aux)
	return err
}

// GormDataType gorm common data type
func (m Attachment) GormDataType() string {
	return "Attachment"
}

// GormDBDataType gorm db data type
func (Attachment) GormDBDataType(db *gorm.DB, field *schema.Field) string {
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

func (jm Attachment) GormValue(ctx context.Context, db *gorm.DB) clause.Expr {
	data, _ := jm.MarshalJSON()
	switch db.Dialector.Name() {
	case "mysql":
		if v, ok := db.Dialector.(*mysql.Dialector); ok && !strings.Contains(v.ServerVersion, "MariaDB") {
			return gorm.Expr("CAST(? AS JSON)", string(data))
		}
	}
	return gorm.Expr("?", string(data))
}

func (jm *Attachment) GenerateFileURL() *Attachment {
	if jm.FileKey != "" {
		jm.FileURL = fmt.Sprintf("%s/api/v1/common/attachments/%s", config.GetInstance().ServerBaseURL, jm.FileKey)
	}
	return jm
}

func (m *Attachment) GetBlurhash() *s3.Blurhash {
	if (m.Blurhash == nil || m.Blurhash.BlurhashData == "") && m.FileKey != "" && helper.IsImageExt(m.FileKey) {
		v, err := s3.New(config.GetInstance()).GetBlurhash(s3.GetBlurhashParams{
			FileKey:       m.FileKey,
			ThumbnailSize: "128w",
		})
		if err != nil {
			log.Println("UnmarshalJSON Attachment get blurhash error", err)
		}
		m.Blurhash = v
	}

	return m.Blurhash
}

func (m *Attachment) GeTThumbnailURL(thumbnailSize string) string {
	var cfg = config.GetInstance()
	var key = strings.TrimPrefix(m.FileKey, "/")
	var imageUrl = fmt.Sprintf("%s/api/v1/common/attachments/%s?thumbnail_size=%s", cfg.ServerBaseURL, key, thumbnailSize)

	return imageUrl

}

func (m *Attachment) DownloadThumbnail() ([]byte, error) {
	return downloader.DownloadFile(m.GeTThumbnailURL("360w"))

}
