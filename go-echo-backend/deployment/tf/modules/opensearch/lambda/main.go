package main

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	opensearch "github.com/opensearch-project/opensearch-go/v2"
	requestsigner "github.com/opensearch-project/opensearch-go/v2/signer/awsv2"
)

type LogItem struct {
	ID        string `json:"@id"`
	Timestamp string `json:"@timestamp"`
	Owner     string `json:"@owner"`
	LogGroup  string `json:"@log_group"`
	LogStream string `json:"@log_stream"`
	Message   string `json:"@message"`
}

type Config struct {
	Endpoint        string `json:"endpoint"`
	Region          string `json:"region"`
	AccessKey       string `json:"access_key"`
	SecretKey       string `json:"secret_key"`
	SessionToken    string `json:"session_token"`
	Service         string `json:"service"`
	IndexNamePrefix string `json:"index_name_prefix"`
	Timezone        string `json:"timezone"`
}

var cfg = &Config{}

func HandleRequest(ctx context.Context, event events.CloudwatchLogsEvent) {
	bulkUpload(ctx, event.AWSLogs.Data)
}

func printJSON(i interface{}) {
	data, _ := json.MarshalIndent(i, "", "   ")
	fmt.Println(string(data))
}

func getCredentialProvider(accessKey, secretAccessKey, token string) aws.CredentialsProviderFunc {
	return func(ctx context.Context) (aws.Credentials, error) {
		c := &aws.Credentials{
			AccessKeyID:     accessKey,
			SecretAccessKey: secretAccessKey,
			SessionToken:    token,
		}
		return *c, nil
	}
}

func bulkUpload(ctx context.Context, rawData string) {
	data, err := base64.StdEncoding.DecodeString(rawData)
	if err != nil {
		log.Fatal("Decode aws logs data error", err)
		return
	}

	zr, err := gzip.NewReader(bytes.NewBuffer(data))
	if err != nil {
		log.Fatal("Gzip reader error", err)
		return
	}
	defer zr.Close()

	var d events.CloudwatchLogsData
	err = json.NewDecoder(zr).Decode(&d)
	if err != nil {
		log.Fatal("Decode aws logs data error", err)
		return
	}

	// awsCfg, err := config.LoadDefaultConfig(ctx)
	awsCfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion(cfg.Region),
		config.WithCredentialsProvider(
			getCredentialProvider(cfg.AccessKey, cfg.SecretKey, cfg.SessionToken),
		),
	)
	if err != nil {
		log.Fatal("Loading aws config error", err) // don't log.fatal in a production-ready app
		return
	}

	// create an AWS request Signer and load AWS configuration using default config folder or env vars.
	signer, err := requestsigner.NewSignerWithService(awsCfg, cfg.Service) // "aoss" for Amazon OpenSearch Serverless
	if err != nil {
		log.Fatal("Request signer error", err) // don't log.fatal in a production-ready app
		return
	}

	// create an opensearch client and use the request-signer
	client, err := opensearch.NewClient(opensearch.Config{
		Addresses: []string{cfg.Endpoint},
		Signer:    signer,
	})
	if err != nil {
		log.Fatal("New opensearch client err", err)
		return
	}

	var payload = transformData(cfg.IndexNamePrefix, d)
	if payload == nil {
		log.Fatal("Transformed data is empty")
		return
	}

	_, err = client.Bulk(payload)
	if err != nil {
		log.Fatal("Perform bulk err", err)
		return
	}

}
func transformData(indexNamePrefix string, d events.CloudwatchLogsData) *bytes.Buffer {
	if d.MessageType == "CONTROL_MESSAGE" {
		return nil
	}

	var bulkRequestBody = bytes.NewBuffer(nil)
	for _, e := range d.LogEvents {
		var date = time.UnixMilli(e.Timestamp).In(getLoc(cfg.Timezone))

		var indexName = fmt.Sprintf("%s-%s", indexNamePrefix, date.Format("2006.01.02"))
		var raw = `{"index":{"_index":"%s","_id":"%s"}}`
		var action = fmt.Sprintf(raw, indexName, e.ID)

		var item = LogItem{
			ID:        e.ID,
			Timestamp: date.Format(time.RFC3339),
			Owner:     d.Owner,
			LogGroup:  d.LogGroup,
			LogStream: d.LogStream,
			Message:   e.Message,
		}

		data, err := json.Marshal(&item)
		if err != nil {
			log.Fatalln("JSON marshal error", item)
			return nil
		}

		bulkRequestBody.WriteString(action)
		bulkRequestBody.WriteString("\n")
		bulkRequestBody.WriteString(string(data))
		bulkRequestBody.WriteString("\n")

	}

	return bulkRequestBody
}

func getLoc(tz string) *time.Location {
	loc, err := time.LoadLocation(tz)
	if err != nil {
		return time.UTC
	}

	return loc
}

func main() {
	var endpoint = os.Getenv("ENDPOINT")
	var region = os.Getenv("AWS_REGION")
	var accessKey = os.Getenv("AWS_ACCESS_KEY_ID")
	var secretKey = os.Getenv("AWS_SECRET_ACCESS_KEY")
	var sessionToken = os.Getenv("AWS_SESSION_TOKEN")
	var service = os.Getenv("SERVICE")

	var indexNamePrefix = os.Getenv("INDEX_NAME_PREFIX")
	if indexNamePrefix == "" {
		indexNamePrefix = "cwl"
	}
	var timezone = os.Getenv("TIMEZONE")
	if timezone == "" {
		timezone = "Asia/Ho_Chi_Minh"
	}

	cfg = &Config{
		Endpoint:        endpoint,
		Region:          region,
		AccessKey:       accessKey,
		SecretKey:       secretKey,
		SessionToken:    sessionToken,
		Service:         service,
		IndexNamePrefix: indexNamePrefix,
		Timezone:        timezone,
	}

	if !strings.HasPrefix("http", cfg.Endpoint) {
		cfg.Endpoint = fmt.Sprintf("https://%s", cfg.Endpoint)
	}

	printJSON(cfg)

	lambda.Start(HandleRequest)
}
