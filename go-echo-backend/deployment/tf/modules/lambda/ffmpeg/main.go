package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/go-playground/validator/v10"
	"github.com/mitchellh/mapstructure"
)

func respJSON(data interface{}, status int) (*events.APIGatewayProxyResponse, error) {
	jsonData, _ := json.Marshal(&data)
	return &events.APIGatewayProxyResponse{
		StatusCode:      status,
		Headers:         map[string]string{"Content-Type": "application/json"},
		Body:            string(jsonData),
		IsBase64Encoded: false,
	}, nil
}

func redirect(location string) (*events.APIGatewayProxyResponse, error) {
	return &events.APIGatewayProxyResponse{
		StatusCode: 307,
		Headers: map[string]string{
			"Location": location,
		},
	}, nil
}

func handler(ctx context.Context, request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	var params APIParams
	var err = mapstructure.Decode(request.QueryStringParameters, &params)
	if err != nil {
		return nil, err
	}
	fmt.Printf("Processing type=%s file_key=%s token=%s no_cache=%t\n", params.Type, params.FileKey, params.Token, params.NoCache)

	err = validator.New().Struct(&params)
	if err != nil {
		return nil, err
	}

	err = ValidateToken(params.Token, params.FileKey)
	if err != nil {
		return nil, err
	}

	if !params.NoCache {
		var ext = path.Ext(params.FileKey)
		var thumbnailFileKey = strings.ReplaceAll(params.FileKey, ext, ".jpeg")

		if _, e := CheckFile(thumbnailFileKey); e == nil {
			return redirect(fmt.Sprintf("%s/%s", config.AwsStorageUrl, thumbnailFileKey))
		}
	}

	var file *os.File
	var filePath = fmt.Sprintf("/tmp/%s", params.FileKey)
	var dir = path.Dir(filePath)

	if fileInfo, err := os.Stat(filePath); err == nil {
		file, err = os.OpenFile(filePath, os.O_RDWR, 0644)
		if err != nil {
			return nil, err
		}

		if fileInfo.Size() < 1024 {
			err = DownloadS3(params.FileKey, file)
			if err != nil {
				return nil, err
			}
		}

	} else {
		err = os.MkdirAll(dir, 0700)
		if err != nil {
			fmt.Printf("Can not create dir %s err=%+v\n", dir, err)
			return nil, err
		}

		file, err = os.Create(filePath)
		if err != nil {
			fmt.Printf("Can not create file %s err=%+v\n", filePath, err)
			return nil, err
		}

		err = DownloadS3(params.FileKey, file)
		if err != nil {
			return nil, err
		}

	}
	defer func() {
		file.Close()
		os.Remove(filePath)
	}()

	switch params.Type {
	case APIRequestTypeThumbnail:
		thumbnailFileKey, err := ReadFrameAndUpload(filePath, params.FileKey)
		if err != nil {
			return nil, err
		}

		return redirect(fmt.Sprintf("%s/%s", config.AwsStorageUrl, thumbnailFileKey))

		// case APIRequestTypeConversion:
		// 	var response = map[string]interface{}{}
		// 	var wg sync.WaitGroup
		// 	wg.Add(1)
		// 	go func() {
		// 		defer wg.Done()
		// 		thumbnailFileKey, err := ReadFrameAndUpload(filePath, params.FileKey)
		// 		if err != nil {
		// 			return
		// 		}

		// 		response["thumbnail_file_key"] = thumbnailFileKey
		// 		response["thumbnail_url"] = fmt.Sprintf("%s/%s", config.AwsStorageUrl, thumbnailFileKey)

		// 	}()

		// 	wg.Add(1)
		// 	go func() {
		// 		defer wg.Done()
		// 		data, err := MovToMP4(filePath)
		// 		if err != nil {
		// 			return
		// 		}
		// 		var ext = path.Ext(params.FileKey)
		// 		var videoFileKey = strings.ReplaceAll(params.FileKey, ext, ".mp4")

		// 		_, err = UploadS3(videoFileKey, bytes.NewBuffer(data), http.DetectContentType(data))
		// 		if err != nil {
		// 			return
		// 		}

		// 		response["video_file_key"] = videoFileKey
		// 		response["video_url"] = fmt.Sprintf("%s/%s", config.AwsStorageUrl, videoFileKey)

		// 	}()

		// 	wg.Wait()

		// 	return respJSON(response, 200)
	}

	return nil, fmt.Errorf("type is invalid")
}

func main() {
	_, err := NewConfig()
	if err != nil {
		panic(err)
	}

	// resp, err := handler(context.Background(), events.APIGatewayProxyRequest{
	// 	QueryStringParameters: map[string]string{
	// 		"type":     "thumbnail",
	// 		"token":    "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJ1cGxvYWRzL21lZGlhL2NnNWFucjJsbGttNmN0cHZxOGswX2ZhYnJpY19jbXFqMWtjdXFqZTdtN2dpMWhnZy5tcDQ_dGh1bWJuYWlsX3NpemU9MTI4MHciLCJleHAiOjE3MDY2MDc5MjB9.URNBzUSctYXyU1xd1caqWTOHl1gVQ4YeBe4KvWXi0HQ",
	// 		"file_key": "uploads/media/cg5anr2llkm6ctpvq8k0_fabric_cmqj1kcuqje7m7gi1hgg.mp4",
	// 	},
	// })
	// if err != nil {
	// 	panic(err)
	// }
	// PrintJSON(resp)

	lambda.Start(handler)
}
