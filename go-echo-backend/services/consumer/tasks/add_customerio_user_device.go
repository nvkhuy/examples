package tasks

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/engineeringinflow/inflow-backend/pkg/models"
	"github.com/hibiken/asynq"
)

type AddCustomerIOUserDeviceTask struct {
	UserID string            `json:"user_id" validate:"required"`
	Device *models.PushToken `json:"device" validate:"required"`
}

func (task AddCustomerIOUserDeviceTask) GetPayload() []byte {
	data, _ := json.Marshal(&task)

	return data
}

// TaskName task name
func (task AddCustomerIOUserDeviceTask) TaskName() string {
	return "add_customerio_user_device"
}

// Handler handler
func (task AddCustomerIOUserDeviceTask) Handler(ctx context.Context, t *asynq.Task) error {
	var err = workerInstance.BindAndValidate(t.Payload(), &task)
	if err != nil {
		return err
	}

	var platform = strings.ToLower(task.Device.Platform)

	if platform == "web" {
		platform = "android"
	}

	err = workerInstance.App.CustomerIOClient.Track.AddDevice(task.UserID, task.Device.Token, platform, map[string]interface{}{
		"last_used": task.Device.LastUsed,
		"user_id":   task.Device.UserID,
		"is_web":    task.Device.Platform == "web",
	})

	return err
}

// Dispatch dispatch event
func (task AddCustomerIOUserDeviceTask) Dispatch(ctx context.Context, opts ...asynq.Option) (*asynq.TaskInfo, error) {
	return workerInstance.SendTaskWithContext(ctx, task, opts...)
}
