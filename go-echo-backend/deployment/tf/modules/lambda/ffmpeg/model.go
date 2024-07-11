package main

type APIRequestType = string

var (
	APIRequestTypeThumbnail  APIRequestType = "thumbnail"
	APIRequestTypeConversion APIRequestType = "conversion"
)

type APIParams struct {
	FileKey string `mapstructure:"file_key" json:"file_key" validate:"required"`

	Type    APIRequestType `mapstructure:"type" json:"type" validate:"oneof=thumbnail conversion"`
	NoCache bool           `mapstructure:"no_cache" json:"no_cache"`
	Token   string         `mapstructure:"token" json:"token" validate:"required"`
}
