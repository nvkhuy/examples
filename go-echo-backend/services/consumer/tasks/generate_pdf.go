package tasks

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/engineeringinflow/inflow-backend/pkg/pdf"
	"github.com/engineeringinflow/inflow-backend/pkg/repo"
	"github.com/hibiken/asynq"
	"github.com/rotisserie/eris"
)

type GeneratePDFTask struct {
	URL               string `json:"url" validate:"required"`
	Selector          string `json:"selector" validate:"required"`
	FileName          string `json:"file_name" validate:"required"`
	Landscape         bool   `json:"landscape"`
	PrintBackground   bool   `json:"print_background"`
	PreferCssPageSize bool   `json:"prefer_css_page_size"`
}

func (task GeneratePDFTask) GetPayload() []byte {
	data, _ := json.Marshal(&task)

	return data
}

// TaskName task name
func (task GeneratePDFTask) TaskName() string {
	return "generate_pdf"
}

// Handler handler
func (task GeneratePDFTask) Handler(ctx context.Context, t *asynq.Task) error {
	var err = workerInstance.BindAndValidate(t.Payload(), &task)
	if err != nil {
		return err
	}
	var start = time.Now()

	data, err := pdf.New(workerInstance.App.DB.Configuration).GetPDF(pdf.GetPDFParams{
		URL:               task.URL,
		Selector:          task.Selector,
		Landscape:         task.Landscape,
		PrintBackground:   task.PrintBackground,
		PreferCssPageSize: task.PreferCssPageSize,
	})
	if err != nil {
		return eris.Wrapf(err, "Generate pdf %s from url=%s selector=%s failed", task.FileName, task.URL, task.Selector)
	}

	if len(data) == 0 {
		return eris.Errorf("Generate pdf %s from url=%s selector=%s failed", task.FileName, task.URL, task.Selector)
	}

	var fileName = fmt.Sprintf("%s.pdf", task.FileName)
	var filePath = fmt.Sprintf("%s/files/%s", workerInstance.App.Config.EFSPath, fileName)
	fileInfo, err := repo.NewExportRepo(workerInstance.App.DB).WriteFile(filePath, data)
	if err != nil {
		return err
	}

	var elapsedTime = time.Since(start)
	workerInstance.Logger.Debugf("Generate pdf %s from url=%s selector=%s elapsed=%s success", fileInfo.DownloadURL, task.URL, task.Selector, elapsedTime)
	return err
}

// Dispatch dispatch event
func (task GeneratePDFTask) Dispatch(ctx context.Context, opts ...asynq.Option) (*asynq.TaskInfo, error) {
	return workerInstance.SendTaskWithContext(ctx, task, opts...)
}
