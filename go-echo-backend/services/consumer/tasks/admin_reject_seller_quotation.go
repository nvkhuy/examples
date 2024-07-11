package tasks

import (
	"context"
	"encoding/json"

	"github.com/hibiken/asynq"
)

type AdminRejectSellerQuotationTask struct {
	AdminID         string `json:"admin_id" validate:"required"`
	InquirySellerID string `json:"inquiry_seller_id" validate:"required"`
}

func (task AdminRejectSellerQuotationTask) GetPayload() []byte {
	data, _ := json.Marshal(&task)

	return data
}

// TaskName task name
func (task AdminRejectSellerQuotationTask) TaskName() string {
	return "admin_reject_seller_quotation"
}

// Handler handler
func (task AdminRejectSellerQuotationTask) Handler(ctx context.Context, t *asynq.Task) error {
	var err = workerInstance.BindAndValidate(t.Payload(), &task)
	if err != nil {
		return err
	}

	return err
}

// Dispatch dispatch event
func (task AdminRejectSellerQuotationTask) Dispatch(ctx context.Context, opts ...asynq.Option) (*asynq.TaskInfo, error) {
	return workerInstance.SendTaskWithContext(ctx, task, opts...)
}
