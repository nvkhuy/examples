package tasks

import (
	"context"

	"encoding/json"

	"github.com/engineeringinflow/inflow-backend/pkg/models"
	"github.com/engineeringinflow/inflow-backend/pkg/models/enums"
	"github.com/engineeringinflow/inflow-backend/pkg/repo"
	"github.com/hibiken/asynq"
)

type CreateUserNotificationTask struct {
	Message          string                           `json:"message" validate:"required"`
	UserID           string                           `json:"user_id" validate:"required"`
	NotificationType enums.UserNotificationType       `json:"notification_type" validate:"required"`
	Metadata         *models.UserNotificationMetadata `json:"metadata"`
}

func (task CreateUserNotificationTask) GetPayload() []byte {
	data, _ := json.Marshal(&task)

	return data
}

// TaskName task name
func (task CreateUserNotificationTask) TaskName() string {
	return "create_user_notification"
}

// Handler handler
func (task CreateUserNotificationTask) Handler(ctx context.Context, t *asynq.Task) error {
	var err = workerInstance.BindAndValidate(t.Payload(), &task)
	if err != nil {
		return err
	}

	// CreateFromPayload notification for admin
	_, err = repo.NewUserNotificationRepo(workerInstance.App.DB).CreateUserNotification(models.UserNotificationForm{
		Message:          task.Message,
		NotificationType: task.NotificationType,
		UserID:           task.UserID,
		Metadata:         task.Metadata,
	})

	return err

}

// Dispatch dispatch event
func (task CreateUserNotificationTask) Dispatch(ctx context.Context, opts ...asynq.Option) (*asynq.TaskInfo, error) {
	return workerInstance.SendTaskWithContext(ctx, task, opts...)
}
