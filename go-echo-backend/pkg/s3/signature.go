package s3

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

const policyDocument = `
{ "expiration": "%s",
  "conditions": [
    {"bucket": "%s"},
    ["starts-with", "$key", "%s"],
	{"acl": "%s"},
	["starts-with", "$content-type", "%s"],
	["content-length-range", 1, %d],
    {"x-amz-credential": "%s"},
    {"x-amz-algorithm": "%s"},
    {"x-amz-date": "%s" }
  ]
}
`

const (
	expirationFormat = "2006-01-02T15:04:05.000Z"
	timeFormat       = "20060102T150405Z"
	shortTimeFormat  = "20060102"
	acl              = "public-read"
	algorithm        = "AWS4-HMAC-SHA256"
)

type Signature struct {
	Key         string `json:"key"`
	URL         string `json:"url"`
	Policy      string `json:"policy"`
	Credential  string `json:"x-amz-credential"`
	Algorithm   string `json:"x-amz-algorithm"`
	Signature   string `json:"x-amz-signature"`
	Date        string `json:"x-amz-date"`
	ACL         string `json:"acl"`
	ContentType string `json:"content-type"`
}

type GenerateSignatureParams struct {
	Key           string
	ContentType   string
	ExpiryMinutes int
	ACL           string
	Bucket        string
	MaxFileSize   int
}

func (client *Client) GenerateSignature(params GenerateSignatureParams) Signature {
	var t = time.Now().Add(time.Minute * time.Duration(params.ExpiryMinutes))
	var formattedShortTime = t.UTC().Format(shortTimeFormat)
	var date = t.UTC().Format(timeFormat)
	var cred = fmt.Sprintf("%s/%s/%s/s3/aws4_request", client.cfg.AWSAccessKeyID, formattedShortTime, client.cfg.AWSS3Region)
	var defaultACL = params.ACL
	if defaultACL == "" {
		defaultACL = acl
	}
	var maxFileSize = params.MaxFileSize
	if maxFileSize == 0 {
		maxFileSize = 40971520
	}
	b64Policy := fmt.Sprintf(policyDocument,
		t.UTC().Format(expirationFormat),
		params.Bucket,
		params.Key,
		defaultACL,
		params.ContentType,
		maxFileSize,
		cred,
		algorithm,
		date,
	)

	// Generate policy
	policy := base64.StdEncoding.EncodeToString([]byte(b64Policy))

	// Generate signature
	h1 := makeHmac([]byte("AWS4"+client.cfg.AWSSecretAccessKey), []byte(date[:8]))
	h2 := makeHmac(h1, []byte(client.cfg.AWSS3Region))
	h3 := makeHmac(h2, []byte("s3"))
	h4 := makeHmac(h3, []byte("aws4_request"))
	signature := hex.EncodeToString(makeHmac(h4, []byte(policy)))

	// Base url
	url := fmt.Sprintf("https://%s.s3.amazonaws.com", params.Bucket)

	return Signature{
		Key:         params.Key,
		URL:         url,
		ACL:         defaultACL,
		Algorithm:   algorithm,
		Credential:  cred,
		Date:        date,
		Policy:      policy,
		Signature:   signature,
		ContentType: params.ContentType,
	}

}

func (client *Client) GetS3SignedURL(params GenerateSignatureParams) string {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(client.cfg.AWSS3Region),
	})
	if err != nil {
		return ""
	}

	var obj = &s3.PutObjectInput{
		Bucket: aws.String(params.Bucket),
		Key:    aws.String(params.Key),
	}

	if params.ACL != "" {
		obj.ACL = aws.String(params.ACL)
	}

	if params.ContentType != "" {
		obj.ContentType = aws.String(params.ContentType)
	}

	r, _ := s3.New(sess).PutObjectRequest(obj)

	if params.ExpiryMinutes == 0 {
		params.ExpiryMinutes = 15
	}

	// Create the pre-signed url with an expiry
	url, err := r.Presign(time.Duration(params.ExpiryMinutes) * time.Minute)
	if err != nil {
		client.logger.Debugf("Failed to generate a pre-signed url: ", err)
		return ""
	}

	return url
}

func makeHmac(key []byte, data []byte) []byte {
	hash := hmac.New(sha256.New, key)
	hash.Write(data)
	return hash.Sum(nil)
}
