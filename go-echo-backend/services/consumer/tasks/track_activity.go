package tasks

import (
	"context"

	"encoding/json"

	"github.com/engineeringinflow/inflow-backend/pkg/customerio"
	"github.com/engineeringinflow/inflow-backend/pkg/models"
	"github.com/engineeringinflow/inflow-backend/pkg/models/enums"
	"github.com/engineeringinflow/inflow-backend/pkg/repo"
	"github.com/hibiken/asynq"
)

type TrackActivityTask struct {
	UserID                string                       `json:"user_id" validate:"required"`
	CountryCode           enums.CountryCode            `json:"country_code"`
	UserTrackActivityForm models.UserTrackActivityForm `json:"user_track_activity_form"`
}

func (task TrackActivityTask) GetPayload() []byte {
	data, _ := json.Marshal(&task)

	return data
}

// TaskName task name
func (task TrackActivityTask) TaskName() string {
	return "track_activity"
}

// Handler handler
func (task TrackActivityTask) Handler(ctx context.Context, t *asynq.Task) error {
	var err = workerInstance.BindAndValidate(t.Payload(), &task)
	if err != nil {
		return err
	}

	updates, err := repo.NewUserRepo(workerInstance.App.DB).TrackActivity(task.UserID, task.UserTrackActivityForm)
	if err != nil {
		return err
	}

	SyncCustomerIOUserTask{
		UserID: task.UserID,
	}.Dispatch(ctx)

	TrackCustomerIOTask{
		UserID: task.UserID,
		Event:  customerio.EventTrackActivity,
		Data:   updates,
	}.Dispatch(ctx)

	return nil
}

// Dispatch dispatch event
func (task TrackActivityTask) Dispatch(ctx context.Context, opts ...asynq.Option) (*asynq.TaskInfo, error) {
	return workerInstance.SendTaskWithContext(ctx, task, opts...)
}
