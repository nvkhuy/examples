package tasks

import (
	"context"

	"encoding/json"

	"github.com/engineeringinflow/inflow-backend/pkg/helper"
	"github.com/engineeringinflow/inflow-backend/pkg/hubspot"
	"github.com/hibiken/asynq"
	"github.com/rotisserie/eris"
)

type HubspotCreateDealTask struct {
	Data *hubspot.DealPropertiesForm `json:"data" validate:"required"`
}

func (task HubspotCreateDealTask) GetPayload() []byte {
	data, _ := json.Marshal(&task)

	return data
}

// TaskName task name
func (task HubspotCreateDealTask) TaskName() string {
	return "hubspot_create_deal"
}

// Handler handler
func (task HubspotCreateDealTask) Handler(ctx context.Context, t *asynq.Task) error {
	var err = workerInstance.BindAndValidate(t.Payload(), &task)
	if err != nil {
		return err
	}
	if !workerInstance.App.Config.IsProd() {
		return nil
	}

	deal, err := workerInstance.App.HubspotClient.CreateDeal(task.Data)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	t.ResultWriter().Write(helper.ToJson(deal))

	return err

}

// Dispatch dispatch event
func (task HubspotCreateDealTask) Dispatch(ctx context.Context, opts ...asynq.Option) (*asynq.TaskInfo, error) {
	return workerInstance.SendTaskWithContext(ctx, task, opts...)
}
