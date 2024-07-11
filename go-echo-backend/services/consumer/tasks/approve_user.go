package tasks

import (
	"context"
	"fmt"

	"encoding/json"

	"github.com/engineeringinflow/inflow-backend/pkg/customerio"
	"github.com/engineeringinflow/inflow-backend/pkg/helper"
	"github.com/engineeringinflow/inflow-backend/pkg/models"
	"github.com/engineeringinflow/inflow-backend/pkg/models/enums"
	"github.com/hibiken/asynq"
	"github.com/rotisserie/eris"
)

type ApproveUserTask struct {
	UserID string `json:"user_id" validate:"required"`
}

func (task ApproveUserTask) GetPayload() []byte {
	data, _ := json.Marshal(&task)

	return data
}

// TaskName task name
func (task ApproveUserTask) TaskName() string {
	return "approve_user"
}

// Handler handler
func (task ApproveUserTask) Handler(ctx context.Context, t *asynq.Task) error {
	var err = workerInstance.BindAndValidate(t.Payload(), &task)
	if err != nil {
		return err
	}

	var user models.User
	err = workerInstance.App.DB.Select("Name", "Email", "ID", "Role").First(&user, "id = ?", task.UserID).Error
	if err != nil {
		return err
	}

	var baseUrl string
	var eventName = customerio.EventNotifyUserApproved

	if user.Role == enums.RoleSeller {
		baseUrl = workerInstance.App.DB.Configuration.SellerPortalBaseURL
		eventName = customerio.EventNotifySellerApproved
	} else {
		baseUrl = workerInstance.App.DB.Configuration.BrandPortalBaseURL
	}

	var data = map[string]interface{}{
		"login_url": fmt.Sprintf("%s/login", baseUrl),
		"email":     user.Email,
		"user_name": user.Name,
	}

	TrackCustomerIOTask{
		UserID: user.ID,
		Event:  eventName,
		Data:   data,
	}.Dispatch(ctx)

	TrackCustomerIOTask{
		UserID: user.ID,
		Event:  customerio.EventWelcomeToBoard,
		Data:   data,
	}.Dispatch(ctx)

	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	t.ResultWriter().Write(helper.ToJson(&data))
	return err

}

// Dispatch dispatch event
func (task ApproveUserTask) Dispatch(ctx context.Context, opts ...asynq.Option) (*asynq.TaskInfo, error) {
	return workerInstance.SendTaskWithContext(ctx, task, opts...)
}
