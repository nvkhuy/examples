package tasks

import (
	"context"
	"encoding/json"

	"github.com/hibiken/asynq"
)

type DeleteCustomerIOUserTask struct {
	UserID string `json:"user_id" validate:"required"`
	IsNew  bool   `json:"is_new"`
}

func (task DeleteCustomerIOUserTask) GetPayload() []byte {
	data, _ := json.Marshal(&task)

	return data
}

// TaskName task name
func (task DeleteCustomerIOUserTask) TaskName() string {
	return "delete_customerio_user"
}

// Handler handler
func (task DeleteCustomerIOUserTask) Handler(ctx context.Context, t *asynq.Task) error {
	var err = workerInstance.BindAndValidate(t.Payload(), &task)
	if err != nil {
		return err
	}

	err = workerInstance.App.CustomerIOClient.Track.Delete(task.UserID)

	return err
}

// Dispatch dispatch event
func (task DeleteCustomerIOUserTask) Dispatch(ctx context.Context, opts ...asynq.Option) (*asynq.TaskInfo, error) {
	return workerInstance.SendTaskWithContext(ctx, task, opts...)
}
