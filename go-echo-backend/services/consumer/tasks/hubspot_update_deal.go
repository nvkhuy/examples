package tasks

import (
	"context"

	"encoding/json"

	"github.com/engineeringinflow/inflow-backend/pkg/helper"
	"github.com/engineeringinflow/inflow-backend/pkg/hubspot"
	"github.com/hibiken/asynq"
	"github.com/rotisserie/eris"
)

type HubspotUpdateDealTask struct {
	DealID string                      `json:"deal_id" validate:"required"`
	Data   *hubspot.DealPropertiesForm `json:"data" validate:"required"`
}

func (task HubspotUpdateDealTask) GetPayload() []byte {
	data, _ := json.Marshal(&task)

	return data
}

// TaskName task name
func (task HubspotUpdateDealTask) TaskName() string {
	return "hubspot_update_deal"
}

// Handler handler
func (task HubspotUpdateDealTask) Handler(ctx context.Context, t *asynq.Task) error {
	var err = workerInstance.BindAndValidate(t.Payload(), &task)
	if err != nil {
		return err
	}
	if !workerInstance.App.Config.IsProd() {
		return nil
	}

	deal, err := workerInstance.App.HubspotClient.UpdateDeal(task.DealID, task.Data)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	t.ResultWriter().Write(helper.ToJson(deal))

	return err

}

// Dispatch dispatch event
func (task HubspotUpdateDealTask) Dispatch(ctx context.Context, opts ...asynq.Option) (*asynq.TaskInfo, error) {
	return workerInstance.SendTaskWithContext(ctx, task, opts...)
}
