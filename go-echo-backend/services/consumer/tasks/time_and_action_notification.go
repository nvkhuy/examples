package tasks

import (
	"context"
	"encoding/json"
	"time"

	"github.com/engineeringinflow/inflow-backend/pkg/customerio"
	"github.com/engineeringinflow/inflow-backend/pkg/repo"
	"github.com/hibiken/asynq"
)

type TimeAndActionNotificationTask struct {
	ID string `json:"id" validate:"required"`
}

func (task TimeAndActionNotificationTask) GetPayload() []byte {
	data, _ := json.Marshal(&task)
	return data
}

// TaskName task name
func (task TimeAndActionNotificationTask) TaskName() string {
	return "time_and_action_notification"
}

// Handler handler
func (task TimeAndActionNotificationTask) Handler(ctx context.Context, t *asynq.Task) error {
	var err = workerInstance.BindAndValidate(t.Payload(), &task)
	if err != nil {
		return err
	}

	tna, err := repo.NewTNARepo(workerInstance.App.DB).Get(repo.GetTNAParams{
		ID: task.ID,
	})
	if err != nil {
		return err
	}

	err = workerInstance.Scheduler.Unregister(tna.ID)
	if err != nil {
		workerInstance.Logger.Errorf("Unregister task_id=%s error=%+v", tna.ID, err)
	}

	for _, userID := range tna.AssigneeIDs {
		// send brand link to brand owner
		_, _ = TrackCustomerIOTask{
			UserID: userID,
			Event:  customerio.EventAdminSendTNANotification,
			Data:   tna.GetCustomerIOMetadata(nil),
		}.DispatchAt(time.Unix(tna.DateFrom, 0), asynq.TaskID(tna.ID))
	}

	return err
}

// Dispatch dispatch event
func (task TimeAndActionNotificationTask) Dispatch(ctx context.Context, opts ...asynq.Option) (*asynq.TaskInfo, error) {
	return workerInstance.SendTaskWithContext(ctx, task, opts...)
}

func (task TimeAndActionNotificationTask) DispatchIn(d time.Duration, opts ...asynq.Option) (*asynq.TaskInfo, error) {
	return workerInstance.SendTaskIn(task, d, opts...)
}

func (task TimeAndActionNotificationTask) DispatchAt(at time.Time, opts ...asynq.Option) (*asynq.TaskInfo, error) {
	return workerInstance.SendTaskAt(task, at, opts...)
}

func (task TimeAndActionNotificationTask) DynamicDispatchAt(id string, at time.Time, opts ...asynq.Option) (*asynq.TaskInfo, error) {
	var options = append([]asynq.Option{
		asynq.ProcessAt(at),
	}, opts...)

	return workerInstance.SendDynamicTask(asynq.NewTask(id, task.GetPayload(), options...), task.Handler)
}
