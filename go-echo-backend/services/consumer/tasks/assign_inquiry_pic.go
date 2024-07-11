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

type AssignInquiryPICTask struct {
	AssignerID string `json:"assignor_id" validate:"required"`
	AssigneeID string `json:"assignee_id" validate:"required"`
	InquiryID  string `json:"inquiry_id" validate:"required"`
}

func (task AssignInquiryPICTask) GetPayload() []byte {
	data, _ := json.Marshal(&task)

	return data
}

// TaskName task name
func (task AssignInquiryPICTask) TaskName() string {
	return "assign_inquiry_pic"
}

// Handler handler
func (task AssignInquiryPICTask) Handler(ctx context.Context, t *asynq.Task) error {
	var err = workerInstance.BindAndValidate(t.Payload(), &task)
	if err != nil {
		return err
	}

	var inquiry models.Inquiry
	err = workerInstance.App.DB.Select("ID", "ReferenceID").First(&inquiry, "id = ?", task.InquiryID).Error
	if err != nil {
		return err
	}

	var assigner models.User
	err = workerInstance.App.DB.Select("ID", "Name", "Email", "Role", "Team").First(&assigner, "id = ?", task.AssignerID).Error
	if err != nil {
		return err
	}

	var eventData = inquiry.GetCustomerIOMetadata(map[string]interface{}{
		"assigner":    assigner.GetCustomerIOMetadata(nil),
		"assignee_id": task.AssigneeID,
	})
	TrackCustomerIOTask{
		UserID: task.AssigneeID,
		Event:  customerio.EventAdminInquiryAssignPIC,
		Data:   eventData,
	}.Dispatch(ctx)

	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	t.ResultWriter().Write(helper.ToJson(&eventData))

	return err

}

// Dispatch dispatch event
func (task AssignInquiryPICTask) Dispatch(ctx context.Context, opts ...asynq.Option) (*asynq.TaskInfo, error) {
	return workerInstance.SendTaskWithContext(ctx, task, opts...)
}
