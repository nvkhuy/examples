package tasks

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/engineeringinflow/inflow-backend/pkg/models"
	"github.com/hibiken/asynq"
	"github.com/rotisserie/eris"
)

type UploadProductFileTask struct {
	UploadID   string `json:"upload_id" validate:"required"`
	SiteName   string `json:"site_name" validate:"required"`
	FileKey    string `json:"file_key" validate:"required"`
	ScrapeDate int64  `json:"scrape_date" validate:"required"`
}

func (task UploadProductFileTask) GetPayload() []byte {
	data, _ := json.Marshal(&task)
	return data
}

// TaskName task name
func (task UploadProductFileTask) TaskName() string {
	return "upload_product_file_task"
}

// Handler handler
func (task UploadProductFileTask) Handler(ctx context.Context, t *asynq.Task) error {
	if err := workerInstance.BindAndValidate(t.Payload(), &task); err != nil {
		return err
	}
	var pathArr = strings.Split(task.FileKey, fmt.Sprintf("%s/", task.SiteName))
	if len(pathArr) != 2 {
		return eris.Errorf("invalid file path, site_name:%s, file_key:%s", task.SiteName, task.FileKey)
	}
	var scapeDate = time.Unix(task.ScrapeDate, 0).Format("2006-01-02")
	reqBody, err := json.Marshal(map[string]string{
		"site_name":           task.SiteName,
		"file_name":           pathArr[1],
		"default_scrape_date": scapeDate,
	})
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/v1/web_scraper/manual_insert", workerInstance.App.Config.ServerWebCrawlerBaseURL), bytes.NewBuffer(reqBody))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	var client = http.Client{
		Timeout: 15 * time.Minute,
	}
	resp, err := client.Do(req)
	if err != nil {
		var updateInfo = models.ProductFileUploadInfo{
			Status:     "failed",
			FailReason: err.Error(),
		}
		if err := workerInstance.App.DB.Model(&updateInfo).Where("id = ?", task.UploadID).Updates(&updateInfo).Error; err != nil {
			return err
		}
	} else {
		defer resp.Body.Close()
		var respMetadata models.ProductFileUploadInfoMetadata
		if err := json.NewDecoder(resp.Body).Decode(&respMetadata); err != nil {
			return err
		}
		var updateInfo = models.ProductFileUploadInfo{
			Status:   "finished",
			Metadata: &respMetadata,
		}
		if err := workerInstance.App.DB.Model(&updateInfo).Where("id = ?", task.UploadID).Updates(&updateInfo).Error; err != nil {
			return err
		}
	}
	return nil
}

// Dispatch dispatch event
func (task UploadProductFileTask) Dispatch(ctx context.Context, opts ...asynq.Option) (*asynq.TaskInfo, error) {
	return workerInstance.SendTaskWithContext(ctx, task, opts...)
}
