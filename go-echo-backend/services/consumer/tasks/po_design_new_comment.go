package tasks

import (
	"context"
	"fmt"

	"encoding/json"

	"github.com/engineeringinflow/inflow-backend/pkg/models"
	"github.com/engineeringinflow/inflow-backend/pkg/models/enums"
	"github.com/hibiken/asynq"
)

type PoDesignNewCommentTask struct {
	PurchaseOrderID string `json:"purchase_order_id" validate:"required"`
	UserID          string `json:"user_id" validate:"required"`
	CommentID       string `json:"comment_id" validate:"required"`
}

func (task PoDesignNewCommentTask) GetPayload() []byte {
	data, _ := json.Marshal(&task)

	return data
}

// TaskName task name
func (task PoDesignNewCommentTask) TaskName() string {
	return "po_design_new_comment"
}

// Handler handler
func (task PoDesignNewCommentTask) Handler(ctx context.Context, t *asynq.Task) error {
	var err = workerInstance.BindAndValidate(t.Payload(), &task)
	if err != nil {
		return err
	}

	var purchaseOrder models.PurchaseOrder
	err = workerInstance.App.DB.Select("ID", "ReferenceID").First(&purchaseOrder, "id = ?", task.PurchaseOrderID).Error
	if err != nil {
		return err
	}

	CreateUserNotificationTask{
		UserID:           task.UserID,
		Message:          fmt.Sprintf("New comment on sample purchase order %s", purchaseOrder.ReferenceID),
		NotificationType: enums.UserNotificationTypePoDesignNewComment,
		Metadata: &models.UserNotificationMetadata{
			CommentID:                task.CommentID,
			PurchaseOrderID:          purchaseOrder.ID,
			PurchaseOrderReferenceID: purchaseOrder.ReferenceID,
		},
	}.Dispatch(ctx)

	return err

}

// Dispatch dispatch event
func (task PoDesignNewCommentTask) Dispatch(ctx context.Context, opts ...asynq.Option) (*asynq.TaskInfo, error) {
	return workerInstance.SendTaskWithContext(ctx, task, opts...)
}
