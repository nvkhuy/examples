package tasks

import (
	"context"

	"encoding/json"

	"github.com/engineeringinflow/inflow-backend/pkg/customerio"
	"github.com/engineeringinflow/inflow-backend/pkg/helper"
	"github.com/engineeringinflow/inflow-backend/pkg/models"
	"github.com/engineeringinflow/inflow-backend/pkg/repo"
	"github.com/hibiken/asynq"
	"github.com/rotisserie/eris"
)

type AssignPurchaseOrderPICTask struct {
	AssignerID      string `json:"assignor_id" validate:"required"`
	AssigneeID      string `json:"assignee_id" validate:"required"`
	PurchaseOrderID string `json:"purchase_order_id" validate:"required"`
}

func (task AssignPurchaseOrderPICTask) GetPayload() []byte {
	data, _ := json.Marshal(&task)

	return data
}

// TaskName task name
func (task AssignPurchaseOrderPICTask) TaskName() string {
	return "assign_purchase_order_pic"
}

// Handler handler
func (task AssignPurchaseOrderPICTask) Handler(ctx context.Context, t *asynq.Task) error {
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

	var assigner models.User
	err = workerInstance.App.DB.Select("ID", "Name", "Email", "Role", "Team").First(&assigner, "id = ?", task.AssignerID).Error
	if err != nil {
		return err
	}

	var eventData = purchaseOrder.GetCustomerIOMetadata(map[string]interface{}{
		"assigner": assigner.GetCustomerIOMetadata(nil),
	})

	TrackCustomerIOTask{
		UserID: task.AssigneeID,
		Event:  customerio.EventAdminPurchaseOrderAssignPIC,
		Data:   eventData,
	}.Dispatch(ctx)

	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	t.ResultWriter().Write(helper.ToJson(&eventData))

	return err

}

// Dispatch dispatch event
func (task AssignPurchaseOrderPICTask) Dispatch(ctx context.Context, opts ...asynq.Option) (*asynq.TaskInfo, error) {
	return workerInstance.SendTaskWithContext(ctx, task, opts...)
}
