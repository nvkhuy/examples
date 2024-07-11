package tasks

import (
	"context"
	"encoding/json"
	"time"

	"github.com/engineeringinflow/inflow-backend/pkg/customerio"
	"github.com/hibiken/asynq"
)

type TrackCustomerIOTask struct {
	UserID string                 `json:"user_id" validate:"required"`
	Event  customerio.Event       `json:"event" validate:"required"`
	Data   map[string]interface{} `json:"data"`
}

func (task TrackCustomerIOTask) GetPayload() []byte {
	data, _ := json.Marshal(&task)
	return data
}

// TaskName task name
func (task TrackCustomerIOTask) TaskName() string {
	return "track_customerio"
}

// Handler handler
func (task TrackCustomerIOTask) Handler(ctx context.Context, t *asynq.Task) error {
	var err = workerInstance.BindAndValidate(t.Payload(), &task)
	if err != nil {
		return err
	}
	err = workerInstance.App.CustomerIOClient.Track.Track(task.UserID, string(task.Event), task.Data)

	return err
}

// Dispatch dispatch event
func (task TrackCustomerIOTask) Dispatch(ctx context.Context, opts ...asynq.Option) (*asynq.TaskInfo, error) {
	return workerInstance.SendTaskWithContext(ctx, task, opts...)
}

func (task TrackCustomerIOTask) DispatchIn(d time.Duration, opts ...asynq.Option) (*asynq.TaskInfo, error) {
	return workerInstance.SendTaskIn(task, d, opts...)
}

func (task TrackCustomerIOTask) DispatchAt(at time.Time, opts ...asynq.Option) (*asynq.TaskInfo, error) {
	return workerInstance.SendTaskAt(task, at, opts...)
}
