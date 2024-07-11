package html

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/engineeringinflow/inflow-backend/pkg/models"
	"log"
	"regexp"
	"strings"

	"github.com/engineeringinflow/inflow-backend/pkg/s3"
	"golang.org/x/net/html"
)

type CompressHTMLOptions struct {
	uploadHTMLPrefix   string
	uploadImagesPrefix string // example: img/ -> img/image_1.png
	content            string
	s3Client           *s3.Client
	storageBucket      string
	storageURL         string
}

func NewHTMLCompress(uploadHTMLPrefix, uploadPrefix, content, bucket, storageURL string, s3Client *s3.Client) CompressHTMLOptions {
	return CompressHTMLOptions{
		uploadHTMLPrefix:   uploadHTMLPrefix,
		uploadImagesPrefix: uploadPrefix,
		content:            content,
		s3Client:           s3Client,
		storageBucket:      bucket,
		storageURL:         storageURL,
	}
}

// Compress HTML save all image to S3 and return content
func (options CompressHTMLOptions) Compress() (output string, contentURL string, err error) {
	if options.s3Client == nil {
		err = errors.New("s3 client is empty")
		return
	}
	if options.storageBucket == "" {
		err = errors.New("storage bucket is empty")
		return
	}
	if options.storageURL == "" {
		err = errors.New("storage url is empty")
		return
	}
	options.content = options.preExtract(options.content)

	doc, err := html.Parse(strings.NewReader(options.content))
	if err != nil {
		log.Println("Error parsing HTML:", err)
		return
	}

	return options.extractedHTML(doc)
}

func (options CompressHTMLOptions) preExtract(before string) (after string) {
	before = strings.ReplaceAll(before, "nbsp;", "")
	after = strings.ReplaceAll(before, "&amp;", "")
	return
}

func (options CompressHTMLOptions) postExtract(before string) (after string) {
	before = strings.ReplaceAll(before, "nbsp;", "")
	after = strings.ReplaceAll(before, "&amp;", "")
	return
}

func (options CompressHTMLOptions) extractedHTML(n *html.Node) (compressed string, contentURL string, err error) {
	var updatedHTML strings.Builder
	var inc = 1
	var traverse func(*html.Node)

	traverse = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "img" {
			for i, attr := range n.Attr {
				if attr.Key == "src" && options.isBase64Src(attr.Val) {
					var localPath string
					fileURL := fmt.Sprintf("%s_image_%d", options.uploadImagesPrefix, inc)
					localPath, err = options.saveImageFromBase64(attr.Val, fileURL)
					if err != nil {
						return
					}
					inc += 1
					n.Attr[i].Val = localPath
				}
			}
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			traverse(c)
		}
	}

	traverse(n)

	if err = html.Render(&updatedHTML, n); err != nil {
		return
	}

	compressed = options.postExtract(updatedHTML.String())
	contentURL, err = options.saveHTMLToS3([]byte(compressed), options.uploadHTMLPrefix)
	return
}

func (options CompressHTMLOptions) isBase64Src(src string) bool {
	// Use a regular expression to check if the src attribute starts with "data:image/png;base64,"
	re := regexp.MustCompile("^data:image/png;base64,")
	return re.MatchString(src)
}

func (options CompressHTMLOptions) saveImageFromBase64(base64String, fileURL string) (url string, err error) {
	data := strings.SplitN(base64String, ",", 2)[1]

	decoded, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		log.Println("Error decoding base64:", err)
		return
	}

	return options.saveImageToS3(decoded, fileURL)
}

func (options CompressHTMLOptions) saveImageToS3(decoded []byte, fileURL string) (url string, err error) {
	url = fmt.Sprintf("%s%s", fileURL, ".png")
	_, err = options.s3Client.UploadFile(s3.UploadFileParams{
		Data:        bytes.NewReader(decoded),
		Bucket:      options.storageBucket,
		ContentType: "image/png",
		ACL:         "private",
		Key:         url,
	})
	url = fmt.Sprintf("https://%s/%s", options.storageURL, url)
	return
}

func (options CompressHTMLOptions) saveHTMLToS3(decoded []byte, fileURL string) (url string, err error) {
	url = fmt.Sprintf("%s%s", fileURL, models.ContentTypeHTML.GetExtension())
	_, err = options.s3Client.UploadFile(s3.UploadFileParams{
		Data:        bytes.NewReader(decoded),
		Bucket:      options.storageBucket,
		ContentType: string(models.ContentTypeHTML),
		ACL:         "private",
		Key:         url,
	})
	return
}
