package tasks

import (
	"context"
	"encoding/json"

	"github.com/hibiken/asynq"
)

type AdminRejectSellerBulkPurchaseOrderQuotation struct {
	AdminID           string `json:"admin_id" validate:"required"`
	SellerQuotationID string `json:"seller_quotation_id" validate:"required"`
}

func (task AdminRejectSellerBulkPurchaseOrderQuotation) GetPayload() []byte {
	data, _ := json.Marshal(&task)

	return data
}

// TaskName task name
func (task AdminRejectSellerBulkPurchaseOrderQuotation) TaskName() string {
	return "admin_reject_seller_bulk_purchase_order_quotation"
}

// Handler handler
func (task AdminRejectSellerBulkPurchaseOrderQuotation) Handler(ctx context.Context, t *asynq.Task) error {
	var err = workerInstance.BindAndValidate(t.Payload(), &task)
	if err != nil {
		return err
	}

	return err
}

// Dispatch dispatch event
func (task AdminRejectSellerBulkPurchaseOrderQuotation) Dispatch(ctx context.Context, opts ...asynq.Option) (*asynq.TaskInfo, error) {
	return workerInstance.SendTaskWithContext(ctx, task, opts...)
}
