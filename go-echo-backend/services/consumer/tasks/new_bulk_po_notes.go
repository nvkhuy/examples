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

type NewBulkPONotesTask struct {
	UserID              string              `json:"user_id" validate:"required"`
	BulkPurchaseOrderID string              `json:"bulk_purchase_order_id" validate:"required"`
	MentionUserIDs      []string            `json:"mention_user_ids" validate:"required"`
	Message             string              `json:"message"`
	Attachments         *models.Attachments `json:"attachments"`
}

func (task NewBulkPONotesTask) GetPayload() []byte {
	data, _ := json.Marshal(&task)

	return data
}

// TaskName task name
func (task NewBulkPONotesTask) TaskName() string {
	return "new_bulk_po_notes"
}

// Handler handler
func (task NewBulkPONotesTask) Handler(ctx context.Context, t *asynq.Task) error {
	var err = workerInstance.BindAndValidate(t.Payload(), &task)
	if err != nil {
		return err
	}

	bulkPurchaseOrder, err := repo.NewBulkPurchaseOrderRepo(workerInstance.App.DB).GetBulkPurchaseOrder(repo.GetBulkPurchaseOrderParams{
		BulkPurchaseOrderID: task.BulkPurchaseOrderID,
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
		"sender":                              sender.GetCustomerIOMetadata(nil),
		"admin_bulk_purchase_order_notes_url": fmt.Sprintf("%s/bulks/%s/notes", workerInstance.App.Config.AdminPortalBaseURL, bulkPurchaseOrder.ID),
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
			Event:  customerio.EventAdminBulkPONewNotes,
			Data:   bulkPurchaseOrder.GetCustomerIOMetadata(extras),
		}.DispatchIn(time.Second * 3)

	}

	return err
}

// Dispatch dispatch event
func (task NewBulkPONotesTask) Dispatch(ctx context.Context, opts ...asynq.Option) (*asynq.TaskInfo, error) {
	return workerInstance.SendTaskWithContext(ctx, task, opts...)
}
