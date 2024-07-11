package repo

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/engineeringinflow/inflow-backend/pkg/db"
	"github.com/engineeringinflow/inflow-backend/pkg/helper"
	"github.com/engineeringinflow/inflow-backend/pkg/models"
	"github.com/engineeringinflow/inflow-backend/pkg/repo/query"
	"github.com/engineeringinflow/inflow-backend/pkg/repo/query/queryfunc"
	"github.com/engineeringinflow/inflow-backend/pkg/s3"
	"gorm.io/gorm"
)

type ProductFileUploadRepo struct {
	db *db.DB
}

func NewProductFileUploadRepo(db *db.DB) *ProductFileUploadRepo {
	return &ProductFileUploadRepo{
		db: db,
	}
}

func (r *ProductFileUploadRepo) UploadFile(req *models.UploadProductFileRequest) (*models.ProductFileUploadInfo, error) {
	var s3Client = s3.New(r.db.Configuration)
	sess, err := s3Client.NewSession()
	if err != nil {
		return nil, err
	}
	var fileName = req.File.Filename
	var fileKey = fmt.Sprintf("web-scraper/%s/%s_%s.csv", req.SiteName, time.Unix(int64(req.ScrapeDate), 0).Format("2006-01-02"), helper.GenerateXID())
	fileData, err := req.File.Open()
	if err != nil {
		return nil, err
	}

	var uploadInfo = models.ProductFileUploadInfo{
		SiteName: req.SiteName,
		Attachment: models.Attachment{
			FileName:    fileName,
			FileKey:     fileKey,
			ContentType: "text/csv",
		},
		ScrapeDate: int64(req.ScrapeDate),
		Status:     "started",
	}
	if err := r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&uploadInfo).Error; err != nil {
			return err
		}
		var svc = s3manager.NewUploader(sess)
		var uploadParams = &s3manager.UploadInput{
			Body:        fileData,
			Bucket:      aws.String(r.db.Configuration.AWSS3WebScraperBucket),
			Key:         aws.String(fileKey),
			ContentType: aws.String("text/csv"),
		}
		_, err = svc.Upload(uploadParams)
		if err != nil {
			return err
		}
		return nil
	}); err != nil {
		return nil, err
	}

	return &uploadInfo, nil
}

func (r *ProductFileUploadRepo) GetProductFileList(params *models.GetProductFileListRequest) *query.Pagination {
	var builder = queryfunc.NewProductFileUploadInfoBuilder(queryfunc.ProductFileUploadInfoOptions{
		QueryBuilderOptions: queryfunc.QueryBuilderOptions{
			Role: params.GetRole(),
		},
	})
	if params.Page == 0 {
		params.Page = 1
	}
	if params.Limit == 0 {
		params.Limit = 12
	}
	result := query.New(r.db, builder).
		WhereFunc(func(builder *query.Builder) {

		}).
		Page(params.Page).
		Limit(params.Limit).
		PagingFunc()

	return result
}
