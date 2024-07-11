package tasks

import (
	"context"
	"encoding/json"

	"github.com/engineeringinflow/inflow-backend/pkg/repo"
	"github.com/hibiken/asynq"
)

type SyncCustomerIOUserTask struct {
	UserID string `json:"user_id" validate:"required"`
}

func (task SyncCustomerIOUserTask) GetPayload() []byte {
	data, _ := json.Marshal(&task)

	return data
}

// TaskName task name
func (task SyncCustomerIOUserTask) TaskName() string {
	return "sync_customerio_user"
}

// Handler handler
func (task SyncCustomerIOUserTask) Handler(ctx context.Context, t *asynq.Task) error {
	var err = workerInstance.BindAndValidate(t.Payload(), &task)
	if err != nil {
		return err
	}

	user, err := repo.NewUserRepo(workerInstance.App.DB).GetCustomerIOUser(task.UserID)
	if err != nil {
		return err
	}

	err = workerInstance.App.CustomerIOClient.Track.Identify(task.UserID, user.GetCustomerIOMetadata(nil))

	return err
}

// Dispatch dispatch event
func (task SyncCustomerIOUserTask) Dispatch(ctx context.Context, opts ...asynq.Option) (*asynq.TaskInfo, error) {
	return workerInstance.SendTaskWithContext(ctx, task, opts...)
}
