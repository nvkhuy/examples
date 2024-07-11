package models

import "github.com/lib/pq"

// AdsVideo's model
type AdsVideo struct {
	// Model
	ID        string    `gorm:"primaryKey" json:"id,omitempty"`
	CreatedAt int64     `gorm:"autoCreateTime" json:"created_at,omitempty"`
	UpdatedAt int64     `gorm:"autoUpdateTime" json:"updated_at,omitempty"`
	DeletedAt DeletedAt `sql:"index" json:"deleted_at,omitempty" swaggertype:"primitive,integer"`

	Description string         `json:"description,omitempty"`
	Thumbnail   *Attachment    `json:"thumbnail,omitempty"`
	URL         string         `json:"url,omitempty"`
	Sections    pq.StringArray `gorm:"type:varchar(200)[]" json:"sections,omitempty"`
}

type AdsVideoUpdateForm struct {
	JwtClaimsInfo
	AdsVideoID string `json:"ads_video_id,omitempty" param:"ads_video_id"  validate:"required"`

	Description string         `json:"description,omitempty"`
	Thumbnail   *Attachment    `json:"thumbnail,omitempty"`
	URL         string         `json:"url,omitempty"`
	Sections    pq.StringArray `gorm:"type:varchar(200)[]" json:"sections,omitempty"`
}

type AdsVideoCreateForm struct {
	JwtClaimsInfo

	Description string         `json:"description,omitempty"`
	Thumbnail   *Attachment    `json:"thumbnail,omitempty"`
	URL         string         `json:"url,omitempty"`
	Sections    pq.StringArray `gorm:"type:varchar(200)[]" json:"sections,omitempty"`
}
