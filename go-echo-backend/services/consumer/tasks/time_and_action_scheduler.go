package tasks

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/engineeringinflow/inflow-backend/pkg/models"
	"github.com/hibiken/asynq"
)

type TimeAndActionSchedulerActionType string

var (
	TimeAndActionSchedulerActionTypeCreate TimeAndActionSchedulerActionType = "create"
	TimeAndActionSchedulerActionTypeUpdate TimeAndActionSchedulerActionType = "update"
	TimeAndActionSchedulerActionTypeDelete TimeAndActionSchedulerActionType = "delete"
)

type TimeAndActionSchedulerTask struct {
	ID         string                           `json:"id" validate:"required"`
	ActionType TimeAndActionSchedulerActionType `json:"action_type" validate:"required"`
}

func (task TimeAndActionSchedulerTask) GetPayload() []byte {
	data, _ := json.Marshal(&task)
	return data
}

// TaskName task name
func (task TimeAndActionSchedulerTask) TaskName() string {
	return "time_and_action_scheduler"
}

// Handler handler
func (task TimeAndActionSchedulerTask) Handler(ctx context.Context, t *asynq.Task) error {
	var err = workerInstance.BindAndValidate(t.Payload(), &task)
	if err != nil {
		return err
	}

	var tna models.TNA
	err = workerInstance.App.DB.Select("ID", "DateFrom", "DateTo").First(&tna, "id = ?", task.ID).Error
	if err != nil {
		return err
	}

	var startTaskID = fmt.Sprintf("start_%s", tna.ID)
	var endTaskID = fmt.Sprintf("end_%s", tna.ID)

	err = workerInstance.Scheduler.Unregister(startTaskID)
	if err != nil {
		workerInstance.Logger.Errorf("Unregister task_id=%s error=%+v", startTaskID, err)
	}

	err = workerInstance.Scheduler.Unregister(endTaskID)
	if err != nil {
		workerInstance.Logger.Errorf("Unregister task_id=%s error=%+v", endTaskID, err)
	}

	if task.ActionType == TimeAndActionSchedulerActionTypeDelete {
		return nil
	}
	if tna.DateFrom > time.Now().Unix() {
		_, _ = TimeAndActionNotificationTask{
			ID: tna.ID,
		}.DynamicDispatchAt(startTaskID, time.Unix(tna.DateFrom, 0))
	}

	if tna.DateTo > time.Now().Unix() {
		_, _ = TimeAndActionNotificationTask{
			ID: tna.ID,
		}.DynamicDispatchAt(endTaskID, time.Unix(tna.DateTo, 0))
	}

	return err
}

// Dispatch dispatch event
func (task TimeAndActionSchedulerTask) Dispatch(ctx context.Context, opts ...asynq.Option) (*asynq.TaskInfo, error) {
	return workerInstance.SendTaskWithContext(ctx, task, opts...)
}

func (task TimeAndActionSchedulerTask) DispatchIn(d time.Duration, opts ...asynq.Option) (*asynq.TaskInfo, error) {
	return workerInstance.SendTaskIn(task, d, opts...)
}

func (task TimeAndActionSchedulerTask) DispatchAt(at time.Time, opts ...asynq.Option) (*asynq.TaskInfo, error) {
	return workerInstance.SendTaskAt(task, at, opts...)
}
