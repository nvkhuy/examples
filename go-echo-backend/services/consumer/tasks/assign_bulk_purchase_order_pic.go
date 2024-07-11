package tasks

import (
	"context"

	"encoding/json"

	"github.com/engineeringinflow/inflow-backend/pkg/customerio"
	"github.com/engineeringinflow/inflow-backend/pkg/helper"
	"github.com/engineeringinflow/inflow-backend/pkg/models"
	"github.com/hibiken/asynq"
	"github.com/rotisserie/eris"
)

type AssignBulkPurchaseOrderPICTask struct {
	AssignerID          string `json:"assignor_id" validate:"required"`
	AssigneeID          string `json:"assignee_id" validate:"required"`
	BulkPurchaseOrderID string `json:"bulk_purchase_order_id" validate:"required"`
}

func (task AssignBulkPurchaseOrderPICTask) GetPayload() []byte {
	data, _ := json.Marshal(&task)

	return data
}

// TaskName task name
func (task AssignBulkPurchaseOrderPICTask) TaskName() string {
	return "assign_bulk_purchase_order_pic"
}

// Handler handler
func (task AssignBulkPurchaseOrderPICTask) Handler(ctx context.Context, t *asynq.Task) error {
	var err = workerInstance.BindAndValidate(t.Payload(), &task)
	if err != nil {
		return err
	}

	var bulkPurchaseOrder models.BulkPurchaseOrder
	err = workerInstance.App.DB.Select("ID", "ReferenceID").First(&bulkPurchaseOrder, "id = ?", task.BulkPurchaseOrderID).Error
	if err != nil {
		return err
	}

	var assigner models.User
	err = workerInstance.App.DB.Select("ID", "Name", "Email", "Role", "Team").First(&assigner, "id = ?", task.AssignerID).Error
	if err != nil {
		return err
	}

	var eventDatra = bulkPurchaseOrder.GetCustomerIOMetadata(map[string]interface{}{
		"assigner": assigner.GetCustomerIOMetadata(nil),
	})

	_, _ = TrackCustomerIOTask{
		UserID: task.AssigneeID,
		Event:  customerio.EventAdminBulkPurchaseOrderAssignPIC,
		Data:   eventDatra,
	}.Dispatch(ctx)

	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	t.ResultWriter().Write(helper.ToJson(&eventDatra))
	return err

}

// Dispatch dispatch event
func (task AssignBulkPurchaseOrderPICTask) Dispatch(ctx context.Context, opts ...asynq.Option) (*asynq.TaskInfo, error) {
	return workerInstance.SendTaskWithContext(ctx, task, opts...)
}
