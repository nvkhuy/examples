package tasks

import (
	"context"
	"encoding/json"
	"github.com/engineeringinflow/inflow-backend/pkg/repo"
	"github.com/hibiken/asynq"
)

type UpdateUserProductClassesTask struct {
	UserID              string `json:"user_id" validate:"required"`
	InquiryID           string `json:"inquiry_id"`
	PurchaseOrderID     string `json:"purchase_order_id"`
	BulkPurchaseOrderID string `json:"bulk_purchase_order_id"`
}

func (task UpdateUserProductClassesTask) GetPayload() []byte {
	data, _ := json.Marshal(&task)
	return data
}

// TaskName task name
func (task UpdateUserProductClassesTask) TaskName() string {
	return "update_user_product_classes"
}

// Handler handler
func (task UpdateUserProductClassesTask) Handler(ctx context.Context, t *asynq.Task) error {
	var err = workerInstance.BindAndValidate(t.Payload(), &task)
	if err != nil {
		return err
	}
	err = repo.NewUserRepo(workerInstance.App.DB).UpdateProductClasses(repo.UpdateUserProductClassesParams{
		UserID:              task.UserID,
		InquiryID:           task.InquiryID,
		PurchaseOrderID:     task.PurchaseOrderID,
		BulkPurchaseOrderID: task.BulkPurchaseOrderID,
	})
	return err
}

// Dispatch dispatch event
func (task UpdateUserProductClassesTask) Dispatch(ctx context.Context, opts ...asynq.Option) (*asynq.TaskInfo, error) {
	return workerInstance.SendTaskWithContext(ctx, task, opts...)
}
