package tasks

import (
	"context"

	"encoding/json"

	"github.com/engineeringinflow/inflow-backend/pkg/models"
	"github.com/engineeringinflow/inflow-backend/pkg/models/enums"
	"github.com/engineeringinflow/inflow-backend/pkg/repo"
	"github.com/hibiken/asynq"
)

type CreateCmsNotificationTask struct {
	Message          string                       `json:"message" validate:"required"`
	NotificationType enums.CmsNotificationType    `json:"notification_type" validate:"required"`
	Metadata         *models.NotificationMetadata `json:"metadata"`
}

func (task CreateCmsNotificationTask) GetPayload() []byte {
	data, _ := json.Marshal(&task)

	return data
}

// TaskName task name
func (task CreateCmsNotificationTask) TaskName() string {
	return "create_cms_notification"
}

// Handler handler
func (task CreateCmsNotificationTask) Handler(ctx context.Context, t *asynq.Task) error {
	var err = workerInstance.BindAndValidate(t.Payload(), &task)
	if err != nil {
		return err
	}

	// CreateFromPayload notification for admin
	_, err = repo.NewCmsNotificationRepo(workerInstance.App.DB).CreateCmsNotification(models.CmsNotificationForm{
		Message:          task.Message,
		NotificationType: task.NotificationType,
		Metadata:         task.Metadata,
	})

	return err

}

// Dispatch dispatch event
func (task CreateCmsNotificationTask) Dispatch(ctx context.Context, opts ...asynq.Option) (*asynq.TaskInfo, error) {
	return workerInstance.SendTaskWithContext(ctx, task, opts...)
}
