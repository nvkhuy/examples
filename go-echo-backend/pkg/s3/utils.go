package s3

import (
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/s3"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

var (
	ErrNotFound = errors.New("not found")
)

func (client *Client) CheckFile(bucket, key string) (string, error) {
	session, err := client.NewSession()
	if err != nil {
		client.logger.Debugf("Create session err: %v\n", err)
		return "", err
	}

	_, err = s3.New(session).HeadObject(&s3.HeadObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case "NotFound": // s3.ErrCodeNoSuchKey does not work, aws is missing this error code so we hardwire a string
				return "", ErrNotFound
			default:
				return "", err
			}
		}
		return "", err
	}

	var url = fmt.Sprintf("https://%s.s3.amazonaws.com/%s", bucket, key)

	return url, nil
}

type Blurhash struct {
	Blurhash            string `json:"blurhash"`
	BlurhashData        string `json:"blurhash_data"`
	BlurhashAvg         string `json:"blurhash_avg"`
	BlurhashThumbnail   string `json:"blurhash_thumbnail"`
	BlurhashImageWidth  string `json:"blurhash_image_width"`
	BlurhashImageHeight string `json:"blurhash_image_height"`
}

func (client *Client) GetBlurhashMetadata(bucket, key string) (*Blurhash, error) {
	session, err := client.NewSession()
	if err != nil {
		client.logger.Debugf("Create session err: %v\n", err)
		return nil, err
	}

	headObject, err := s3.New(session).HeadObject(&s3.HeadObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case "NotFound": // s3.ErrCodeNoSuchKey does not work, aws is missing this error code so we hardwire a string
				return nil, ErrNotFound
			default:
				return nil, err
			}
		}
		return nil, err
	}

	if headObject.Metadata != nil {
		var blurhash = Blurhash{
			Blurhash:            getMetadataValue(headObject.Metadata, "blurhash"),
			BlurhashAvg:         getMetadataValue(headObject.Metadata, "blurhash_avg"),
			BlurhashImageHeight: getMetadataValue(headObject.Metadata, "blurhash_image_height"),
			BlurhashImageWidth:  getMetadataValue(headObject.Metadata, "blurhash_image_width"),
			BlurhashThumbnail:   getMetadataValue(headObject.Metadata, "blurhash_thumbnail"),
		}

		return &blurhash, nil
	}

	return nil, errors.New("blurhash not found")
}

func getMetadataValue(m map[string]*string, key string) string {
	caser := cases.Title(language.English)

	if v, ok := m[key]; ok && aws.StringValue(v) != "" {
		return aws.StringValue(v)
	}

	if v, ok := m[caser.String(key)]; ok && aws.StringValue(v) != "" {
		return aws.StringValue(v)
	}

	return ""
}
