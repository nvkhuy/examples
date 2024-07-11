package tasks

import (
	"context"
	"encoding/json"
	"github.com/engineeringinflow/inflow-backend/pkg/models"
	"github.com/engineeringinflow/inflow-backend/pkg/repo"
	"github.com/hibiken/asynq"
)

type CreateSysNotificationTask struct {
	models.SysNotification
}

func (task CreateSysNotificationTask) GetPayload() []byte {
	data, _ := json.Marshal(&task)
	return data
}

// TaskName task name
func (task CreateSysNotificationTask) TaskName() string {
	return "create_sys_notification"
}

// Handler handler
func (task CreateSysNotificationTask) Handler(ctx context.Context, t *asynq.Task) (err error) {
	err = workerInstance.BindAndValidate(t.Payload(), &task)
	if err != nil {
		return err
	}
	_, err = repo.NewSysNotificationRepo(workerInstance.App.DB).Create(repo.CreateSysNotificationsParams{
		SysNotification: task.SysNotification,
	})

	return err
}

// Dispatch dispatch event
func (task CreateSysNotificationTask) Dispatch(ctx context.Context, opts ...asynq.Option) (*asynq.TaskInfo, error) {
	return workerInstance.SendTaskWithContext(ctx, task, opts...)
}
