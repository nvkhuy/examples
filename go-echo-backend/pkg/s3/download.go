package s3

import (
	"context"
	"errors"
	"io"
	"sync"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

type DownloadParams struct {
	Bucket string
	Key    string
}

func (client *Client) Download(params DownloadParams) ([]byte, error) {
	session, err := client.NewSession()
	if err != nil {
		client.logger.Errorf("Create session err: %v", err)
		return nil, err
	}

	head, err := s3.New(session).HeadObject(&s3.HeadObjectInput{
		Bucket: aws.String(params.Bucket),
		Key:    aws.String(params.Key),
	})
	if err != nil {
		client.logger.Errorf("Head object err: %v", err)
		return nil, err
	}

	if awsErr, ok := err.(awserr.Error); ok && awsErr.Code() == s3.ErrCodeNoSuchKey {
		return nil, errors.New("not found")
	}
	if err != nil {
		return nil, err
	}

	var buf = aws.NewWriteAtBuffer(make([]byte, 0, int(*head.ContentLength)))

	var downloader = s3manager.NewDownloaderWithClient(s3.New(session), func(d *s3manager.Downloader) {
		d.PartSize = 10 * 1024 * 1024
	})
	_, err = downloader.Download(buf, &s3.GetObjectInput{
		Bucket: aws.String(params.Bucket),
		Key:    aws.String(params.Key),
	})
	if err != nil {
		client.logger.Errorf("Download file from bucket = %s key = %s error: %v", params.Bucket, params.Key, err)
		return nil, err
	}

	return buf.Bytes(), nil
}

type DownloadFileParams struct {
	Bucket   string
	Key      string
	Output   *aws.WriteAtBuffer
	FileName string
}

func (client *Client) DownloadFile(params *DownloadFileParams) error {
	session, err := client.NewSession()
	if err != nil {
		client.logger.Debugf("Create session err: %v\n", err)
		return err
	}

	var downloader = s3manager.NewDownloaderWithClient(s3.New(session))
	_, err = downloader.Download(params.Output, &s3.GetObjectInput{
		Bucket: aws.String(params.Bucket),
		Key:    aws.String(params.Key),
	})
	if err != nil {
		client.logger.Debugf("Download file from bucket = %s key = %s error: %v\n", params.Bucket, params.Key, err)
		return err
	}

	return nil
}

func (client *Client) DownloadFiles(params []*DownloadFileParams) error {
	var wg = sync.WaitGroup{}

	for _, param := range params {
		wg.Add(1)
		go func(param *DownloadFileParams) {
			defer wg.Done()
			client.DownloadFile(param)

		}(param)
	}

	wg.Wait()

	return nil
}

type GetObjectParams struct {
	Bucket string
	Key    string
}

func (client *Client) GetObject(params *GetObjectParams) ([]byte, error) {
	session, err := client.NewSession()
	if err != nil {
		client.logger.Errorf("Create session err: %v", err)
		return nil, err
	}

	var s3Client = s3.New(session)

	var getInput = &s3.GetObjectInput{
		Bucket: aws.String(params.Bucket),
		Key:    aws.String(params.Key),
	}

	resp, err := s3Client.GetObjectWithContext(context.TODO(), getInput)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}
