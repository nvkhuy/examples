package tasks

import (
	"context"
	"fmt"
	"time"

	"encoding/json"

	"github.com/engineeringinflow/inflow-backend/pkg/customerio"
	"github.com/engineeringinflow/inflow-backend/pkg/helper"
	"github.com/engineeringinflow/inflow-backend/pkg/repo"
	"github.com/hibiken/asynq"
)

type OnboardSellerTask struct {
	UserID string `json:"user_id" validate:"required"`
}

func (task OnboardSellerTask) GetPayload() []byte {
	data, _ := json.Marshal(&task)

	return data
}

// TaskName task name
func (task OnboardSellerTask) TaskName() string {
	return "onboard_seller"
}

// Handler handler
func (task OnboardSellerTask) Handler(ctx context.Context, t *asynq.Task) error {
	var err = workerInstance.BindAndValidate(t.Payload(), &task)
	if err != nil {
		return err
	}

	user, err := repo.NewUserRepo(workerInstance.App.DB).GetCustomerIOUser(task.UserID)
	if err != nil {
		return err
	}

	var data = user.GetCustomerIOMetadata(map[string]interface{}{
		"login_url": fmt.Sprintf("%s/login", workerInstance.App.Config.SellerPortalBaseURL),
	})
	err = workerInstance.App.CustomerIOClient.Track.Identify(task.UserID, data)
	if err != nil {
		return err
	}

	TrackCustomerIOTask{
		UserID: task.UserID,
		Event:  customerio.EventOnboardSeller,
		Data:   data,
	}.Dispatch(ctx)

	t.ResultWriter().Write(helper.ToJson(&data))

	return err

}

// Dispatch dispatch event
func (task OnboardSellerTask) Dispatch(ctx context.Context, opts ...asynq.Option) (*asynq.TaskInfo, error) {
	return workerInstance.SendTaskWithContext(ctx, task, opts...)
}

func (task OnboardSellerTask) DispatchIn(ctx context.Context, duration time.Duration, opts ...asynq.Option) (*asynq.TaskInfo, error) {
	return workerInstance.SendTaskIn(task, duration, opts...)
}
