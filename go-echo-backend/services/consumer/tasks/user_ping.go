package tasks

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/engineeringinflow/inflow-backend/pkg/models"
	"github.com/hibiken/asynq"
)

type UserPingTask struct {
	UserID       string `json:"user_id" validate:"required"`
	IsOffline    *bool  `json:"is_offline"`
	LastOnlineAt *int64 `json:"last_online_at"`
}

func (task UserPingTask) GetPayload() []byte {
	data, _ := json.Marshal(&task)

	return data
}

// TaskName task name
func (task UserPingTask) TaskName() string {
	return "user_ping"
}

// Handler handler
func (task UserPingTask) Handler(ctx context.Context, t *asynq.Task) error {
	var err = workerInstance.BindAndValidate(t.Payload(), &task)
	if err != nil {
		return err
	}

	var key = fmt.Sprintf("user_ping_%s", task.UserID)

	cancel, err := workerInstance.App.DB.Locker.AcquireLock(key, time.Second*10)
	if err != nil {
		return err
	}

	defer cancel()

	var updates = models.User{
		IsOffline:    task.IsOffline,
		LastOnlineAt: task.LastOnlineAt,
	}
	err = workerInstance.App.DB.Model(&models.User{}).Where("id = ?", task.UserID).Updates(&updates).Error

	return err

}

// Dispatch dispatch event
func (task UserPingTask) Dispatch(ctx context.Context, opts ...asynq.Option) (*asynq.TaskInfo, error) {
	return workerInstance.SendTaskWithContext(ctx, task, opts...)
}
