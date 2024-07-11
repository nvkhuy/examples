package main

import (
	"errors"
	"fmt"
	"io"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

func UploadS3(key string, data io.Reader, contentType string) (string, error) {
	s, err := session.NewSession()
	if err != nil {
		return "", err
	}

	var svc = s3manager.NewUploader(s)
	var uploadParams = &s3manager.UploadInput{
		Bucket:      aws.String(config.AwsStorageBucket),
		Key:         aws.String(key),
		ACL:         aws.String("private"),
		Body:        data,
		ContentType: aws.String(contentType),
	}

	result, err := svc.Upload(uploadParams)
	if err != nil {
		return "", err
	}

	return result.Location, nil
}

func DownloadS3(key string, output io.WriterAt) error {
	s, err := session.NewSession()
	if err != nil {
		return err
	}

	var downloader = s3manager.NewDownloaderWithClient(s3.New(s))
	_, err = downloader.Download(output, &s3.GetObjectInput{
		Bucket: aws.String(config.AwsStorageBucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return err
	}

	return nil
}

func CheckFile(key string) (string, error) {
	s, err := session.NewSession()
	if err != nil {
		return "", err
	}
	_, err = s3.New(s).HeadObject(&s3.HeadObjectInput{
		Bucket: aws.String(config.AwsStorageBucket),
		Key:    aws.String(key),
	})
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case "NotFound": // s3.ErrCodeNoSuchKey does not work, aws is missing this error code so we hardwire a string
				return "", errors.New("file not found")
			default:
				return "", err
			}
		}
		return "", err
	}

	var url = fmt.Sprintf("https://%s.s3.amazonaws.com/%s", config.AwsStorageBucket, key)

	return url, nil
}
