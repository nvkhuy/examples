package tasks

import (
	"context"

	"encoding/json"

	"github.com/engineeringinflow/inflow-backend/pkg/helper"
	"github.com/engineeringinflow/inflow-backend/pkg/hubspot"
	"github.com/hibiken/asynq"
	"github.com/rotisserie/eris"
)

type HubspotCreateContactTask struct {
	Data *hubspot.ContactPropertiesForm `json:"data" validate:"required"`
}

func (task HubspotCreateContactTask) GetPayload() []byte {
	data, _ := json.Marshal(&task)

	return data
}

// TaskName task name
func (task HubspotCreateContactTask) TaskName() string {
	return "hubspot_create_contact"
}

// Handler handler
func (task HubspotCreateContactTask) Handler(ctx context.Context, t *asynq.Task) error {
	var err = workerInstance.BindAndValidate(t.Payload(), &task)
	if err != nil {
		return err
	}
	if !workerInstance.App.Config.IsProd() {
		return nil
	}

	contact, err := workerInstance.App.HubspotClient.CreateContact(task.Data)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	t.ResultWriter().Write(helper.ToJson(contact))

	return err

}

// Dispatch dispatch event
func (task HubspotCreateContactTask) Dispatch(ctx context.Context, opts ...asynq.Option) (*asynq.TaskInfo, error) {
	return workerInstance.SendTaskWithContext(ctx, task, opts...)
}
