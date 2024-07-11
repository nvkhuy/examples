package tasks

import (
	"context"
	"fmt"
	"time"

	"encoding/json"

	"github.com/engineeringinflow/inflow-backend/pkg/helper"
	"github.com/engineeringinflow/inflow-backend/pkg/hubspot"
	"github.com/engineeringinflow/inflow-backend/pkg/models/enums"
	"github.com/engineeringinflow/inflow-backend/pkg/repo"
	"github.com/hibiken/asynq"
)

type OnboardUserTask struct {
	UserID string `json:"user_id" validate:"required"`
}

func (task OnboardUserTask) GetPayload() []byte {
	data, _ := json.Marshal(&task)

	return data
}

// TaskName task name
func (task OnboardUserTask) TaskName() string {
	return "onboard_user"
}

// Handler handler
func (task OnboardUserTask) Handler(ctx context.Context, t *asynq.Task) error {
	var err = workerInstance.BindAndValidate(t.Payload(), &task)
	if err != nil {
		return err
	}
	user, err := repo.NewUserRepo(workerInstance.App.DB).GetCustomerIOUser(task.UserID)
	if err != nil {
		return err
	}

	var data = user.GetCustomerIOMetadata(map[string]interface{}{
		"login_url": fmt.Sprintf("%s/login", workerInstance.App.Config.BrandPortalBaseURL),
	})
	err = workerInstance.App.CustomerIOClient.Track.Identify(task.UserID, data)
	if err != nil {
		return err
	}

	if user.Role == enums.RoleClient && workerInstance.App.Config.IsProd() {
		HubspotCreateContactTask{
			Data: &hubspot.ContactPropertiesForm{
				Email:          user.Email,
				Firstname:      user.FirstName,
				Lastname:       user.LastName,
				Phone:          user.PhoneNumber,
				Lifecyclestage: "lead",
			},
		}.Dispatch(ctx)

	}

	t.ResultWriter().Write(helper.ToJson(&data))

	return err

}

// Dispatch dispatch event
func (task OnboardUserTask) Dispatch(ctx context.Context, opts ...asynq.Option) (*asynq.TaskInfo, error) {
	return workerInstance.SendTaskWithContext(ctx, task, opts...)
}

func (task OnboardUserTask) DispatchIn(ctx context.Context, duration time.Duration, opts ...asynq.Option) (*asynq.TaskInfo, error) {
	return workerInstance.SendTaskIn(task, duration, opts...)
}
