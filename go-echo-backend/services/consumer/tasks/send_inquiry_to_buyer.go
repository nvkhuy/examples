package tasks

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/engineeringinflow/inflow-backend/pkg/customerio"
	"github.com/engineeringinflow/inflow-backend/pkg/models"
	"github.com/engineeringinflow/inflow-backend/pkg/models/enums"
	"github.com/engineeringinflow/inflow-backend/pkg/repo"
	"github.com/engineeringinflow/inflow-backend/pkg/repo/query/queryfunc"
	"github.com/hibiken/asynq"
)

type SendInquiryToBuyerTask struct {
	UserID string                 `json:"user_id" validate:"required"`
	Event  customerio.Event       `json:"event" validate:"required"`
	Data   map[string]interface{} `json:"data" validate:"required"`
}

func (task SendInquiryToBuyerTask) GetPayload() []byte {
	data, _ := json.Marshal(&task)

	return data
}

// TaskName task name
func (task SendInquiryToBuyerTask) TaskName() string {
	return "send_inquiry_to_buyer"
}

// Handler handler
func (task SendInquiryToBuyerTask) Handler(ctx context.Context, t *asynq.Task) (err error) {
	err = workerInstance.BindAndValidate(t.Payload(), &task)
	if err != nil {
		return err
	}

	var products []*models.Product
	products, _ = repo.NewProductRepo(workerInstance.App.DB).
		GetProductByPageSectionType(string(enums.PageSectionTypeCatalogDrop), queryfunc.ProductBuilderOptions{})

	if len(products) > 0 {
		imgInc := 1
		for _, p := range products {
			if p.Attachments != nil && len(*p.Attachments) > 0 {
				att := *p.Attachments
				task.Data[fmt.Sprintf("img_%d", imgInc)] = att[0].GenerateFileURL().FileURL
				imgInc += 1
			}
		}
	}

	err = workerInstance.App.CustomerIOClient.Track.Track(task.UserID, string(task.Event), task.Data)

	return
}

// Dispatch dispatch event
func (task SendInquiryToBuyerTask) Dispatch(ctx context.Context, opts ...asynq.Option) (*asynq.TaskInfo, error) {
	return workerInstance.SendTaskWithContext(ctx, task, opts...)
}
