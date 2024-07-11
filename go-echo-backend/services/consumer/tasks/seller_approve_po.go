package tasks

import (
	"context"
	"encoding/json"

	"github.com/hibiken/asynq"
)

type SellerApprovePOTask struct {
	InquirySellerID string `json:"inquiry_seller_id" validate:"required"`
}

func (task SellerApprovePOTask) GetPayload() []byte {
	data, _ := json.Marshal(&task)

	return data
}

// TaskName task name
func (task SellerApprovePOTask) TaskName() string {
	return "seller_approve_po"
}

// Handler handler
func (task SellerApprovePOTask) Handler(ctx context.Context, t *asynq.Task) error {
	var err = workerInstance.BindAndValidate(t.Payload(), &task)
	if err != nil {
		return err
	}

	return err
}

// Dispatch dispatch event
func (task SellerApprovePOTask) Dispatch(ctx context.Context, opts ...asynq.Option) (*asynq.TaskInfo, error) {
	return workerInstance.SendTaskWithContext(ctx, task, opts...)
}
