package tasks

import (
	"context"
	"encoding/json"

	"github.com/hibiken/asynq"
)

type AdminApproveSellerQuotationTask struct {
	AdminID         string `json:"admin_id" validate:"required"`
	InquirySellerID string `json:"inquiry_seller_id" validate:"required"`
}

func (task AdminApproveSellerQuotationTask) GetPayload() []byte {
	data, _ := json.Marshal(&task)

	return data
}

// TaskName task name
func (task AdminApproveSellerQuotationTask) TaskName() string {
	return "admin_approve_seller_quotation"
}

// Handler handler
func (task AdminApproveSellerQuotationTask) Handler(ctx context.Context, t *asynq.Task) error {
	var err = workerInstance.BindAndValidate(t.Payload(), &task)
	if err != nil {
		return err
	}

	return err
}

// Dispatch dispatch event
func (task AdminApproveSellerQuotationTask) Dispatch(ctx context.Context, opts ...asynq.Option) (*asynq.TaskInfo, error) {
	return workerInstance.SendTaskWithContext(ctx, task, opts...)
}
