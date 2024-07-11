package s3

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/engineeringinflow/inflow-backend/pkg/config"
	"github.com/engineeringinflow/inflow-backend/pkg/helper"
	"github.com/engineeringinflow/inflow-backend/pkg/logger"
	"github.com/go-resty/resty/v2"
)

type Client struct {
	cfg    *config.Configuration
	logger *logger.Logger
}

func New(cfg *config.Configuration) *Client {
	var client = &Client{
		cfg:    cfg,
		logger: logger.New("s3"),
	}

	return client
}

func (client *Client) NewSession() (*session.Session, error) {
	session, err := session.NewSession(&aws.Config{
		Region:      aws.String(client.cfg.AWSS3Region),
		Credentials: credentials.NewStaticCredentials(client.cfg.AWSAccessKeyID, client.cfg.AWSSecretAccessKey, ""),
	})

	return session, err
}

type GetBlurhashParams struct {
	ThumbnailSize string `json:"thumbnail_size" query:"thumbnail_size" form:"thumbnail_size"`
	FileKey       string `json:"file_key" query:"file_key" form:"file_key" param:"file_key"`
	NoCache       bool   `json:"no_cache" query:"no_cache" form:"no_cache" param:"no_cache"`
}

func (c *Client) GetBlurhash(params GetBlurhashParams) (*Blurhash, error) {
	var key = strings.TrimPrefix(params.FileKey, "/")

	if params.ThumbnailSize != "" && helper.IsImageExt(params.FileKey) {
		if blurhash, err := c.GetBlurhashMetadata(c.cfg.AWSS3StorageBucket, key); err == nil && !params.NoCache && blurhash.Blurhash != "" {
			return blurhash, nil
		}

		var keyWithSize = fmt.Sprintf("blur/%s/%s", params.ThumbnailSize, key)
		token, err := c.cfg.GetMediaToken(keyWithSize)
		if err != nil {
			return nil, err
		}

		var url = fmt.Sprintf("%s?token=%s&size=%s&key=%s", c.cfg.LambdaAPIBlurURL, token, params.ThumbnailSize, params.FileKey)
		if params.NoCache {
			url = fmt.Sprintf("%s?token=%s&size=%s&key=%s&no_cache=%t", c.cfg.LambdaAPIBlurURL, token, params.ThumbnailSize, params.FileKey, params.NoCache)
		}

		var blurhash Blurhash
		resp, err := resty.New().
			SetDebug(true).R().
			SetHeader("Content-Type", "application/json").
			SetHeader("Accept", "application/json").
			Get(url)
		if err != nil {
			return nil, err
		}

		err = json.Unmarshal(resp.Body(), &blurhash)

		return &blurhash, err

	}

	return nil, errors.New("not able to generate blurhash")

}
