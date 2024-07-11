package image

// import (
// 	"bytes"
// 	"fmt"

// 	"github.com/aws/aws-sdk-go/aws"
// 	"github.com/engineeringinflow/inflow-backend/pkg/config"
// 	"github.com/engineeringinflow/inflow-backend/pkg/s3"
// 	"github.com/h2non/bimg"
// 	"github.com/rotisserie/eris"
// )

// type Client struct {
// 	cfg *config.Configuration
// }

// func New(cfg *config.Configuration) *Client {
// 	return &Client{
// 		cfg: cfg,
// 	}
// }

// type ResizeImageParams struct {
// 	FileKey       string        `json:"file_key"`
// 	ThumbnailSize ThumbnailSize `json:"thumbnail_size"`
// }

// func (c *Client) ResizeImage(params ResizeImageParams) (string, error) {
// 	var s3Client = s3.New(c.cfg)
// 	var dimension = params.ThumbnailSize.GetDimension()
// 	if !dimension.IsValid() {
// 		return "", eris.New("Invalid size")
// 	}

// 	var resizeKey = fmt.Sprintf("%s/%s", params.ThumbnailSize, params.FileKey)
// 	url, err := s3Client.CheckFile(c.cfg.AWSS3CdnBucket, resizeKey)
// 	if err == nil && url != "" {
// 		return fmt.Sprintf("https://%s/%s", c.cfg.CDNURL, resizeKey), nil
// 	}

// 	output, err := s3Client.Download(s3.DownloadParams{
// 		Bucket: c.cfg.AWSS3StorageBucket,
// 		Key:    params.FileKey,
// 	})
// 	if err != nil {
// 		return "", err
// 	}

// 	var img = bimg.NewImage(output)
// 	data, err := img.Resize(int(dimension.Width), int(dimension.Height))
// 	if err != nil {
// 		return "", err
// 	}

// 	var metadata = map[string]*string{
// 		"size": aws.String(fmt.Sprintf("%d", img.Length())),
// 	}

// 	if meta, err := img.Metadata(); err == nil {
// 		metadata["width"] = aws.String(fmt.Sprintf("%d", meta.Size.Width))
// 		metadata["height"] = aws.String(fmt.Sprintf("%d", meta.Size.Height))
// 		metadata["alpha"] = aws.String(fmt.Sprintf("%t", meta.Alpha))
// 		metadata["type"] = aws.String(meta.Type)
// 		metadata["orientation"] = aws.String(fmt.Sprintf("%d", meta.Orientation))
// 	}

// 	url, err = s3Client.UploadFile(s3.UploadFileParams{
// 		Data:        bytes.NewBuffer(data),
// 		Bucket:      c.cfg.AWSS3CdnBucket,
// 		Key:         resizeKey,
// 		Metadata:    metadata,
// 		ContentType: fmt.Sprintf("image/%s", img.Type()),
// 		ACL:         *aws.String("private"),
// 	})
// 	if err != nil {
// 		return "", err
// 	}

// 	if url == "" {
// 		return "", eris.New("Invalid URL")
// 	}

// 	return fmt.Sprintf("https://%s/%s", c.cfg.CDNURL, resizeKey), nil
// }
