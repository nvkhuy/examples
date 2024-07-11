package tasks

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/engineeringinflow/inflow-backend/pkg/customerio"
	"github.com/engineeringinflow/inflow-backend/pkg/models"
	"github.com/engineeringinflow/inflow-backend/pkg/repo"
	"github.com/hibiken/asynq"
)

type NewPONotesTask struct {
	UserID          string              `json:"user_id" validate:"required"`
	PurchaseOrderID string              `json:"purchase_order_id" validate:"required"`
	MentionUserIDs  []string            `json:"mention_user_ids" validate:"required"`
	Message         string              `json:"message"`
	Attachments     *models.Attachments `json:"attachments"`
}

func (task NewPONotesTask) GetPayload() []byte {
	data, _ := json.Marshal(&task)

	return data
}

// TaskName task name
func (task NewPONotesTask) TaskName() string {
	return "new_po_notes"
}

// Handler handler
func (task NewPONotesTask) Handler(ctx context.Context, t *asynq.Task) error {
	var err = workerInstance.BindAndValidate(t.Payload(), &task)
	if err != nil {
		return err
	}

	purchaseOrder, err := repo.NewPurchaseOrderRepo(workerInstance.App.DB).GetPurchaseOrder(repo.GetPurchaseOrderParams{
		PurchaseOrderID: task.PurchaseOrderID,
	})
	if err != nil {
		return err
	}

	var sender models.User
	err = workerInstance.App.DB.Select("ID", "Name", "Email").First(&sender, "id = ?", task.UserID).Error
	if err != nil {
		return err
	}

	var extras = map[string]interface{}{
		"sender":                         sender.GetCustomerIOMetadata(nil),
		"admin_purchase_order_notes_url": fmt.Sprintf("%s/samples/%s/customer", workerInstance.App.Config.AdminPortalBaseURL, purchaseOrder.ID),
	}
	if task.Message != "" {
		extras["message"] = task.Message
	}

	if task.Attachments != nil {
		extras["attachments"] = task.Attachments
	}

	for _, userID := range task.MentionUserIDs {
		SyncCustomerIOUserTask{
			UserID: userID,
		}.Dispatch(ctx)

		TrackCustomerIOTask{
			UserID: userID,
			Event:  customerio.EventAdminPONewNotes,
			Data:   purchaseOrder.GetCustomerIOMetadata(extras),
		}.DispatchIn(time.Second * 3)

	}

	return err
}

// Dispatch dispatch event
func (task NewPONotesTask) Dispatch(ctx context.Context, opts ...asynq.Option) (*asynq.TaskInfo, error) {
	return workerInstance.SendTaskWithContext(ctx, task, opts...)
}
