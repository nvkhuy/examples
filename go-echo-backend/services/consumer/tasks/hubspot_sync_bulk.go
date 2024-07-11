package tasks

import (
	"context"

	"encoding/json"

	"github.com/hibiken/asynq"
)

type HubspotSyncBulkTask struct {
	BulkPurchaseOrderID string `json:"bulk_purchase_order_id" validate:"required"`
	UserID              string `json:"user_id" validate:"required"`
	IsAdmin             bool   `json:"is_admin"`
}

func (task HubspotSyncBulkTask) GetPayload() []byte {
	data, _ := json.Marshal(&task)

	return data
}

// TaskName task name
func (task HubspotSyncBulkTask) TaskName() string {
	return "hubspot_sync_bulk"
}

// Handler handler
func (task HubspotSyncBulkTask) Handler(ctx context.Context, t *asynq.Task) error {
	var err = workerInstance.BindAndValidate(t.Payload(), &task)
	if err != nil {
		return err
	}
	if !workerInstance.App.Config.IsProd() {
		return nil
	}

	// bulk, err := repo.NewBulkPurchaseOrderRepo(workerInstance.App.DB).GetBulkPurchaseOrder(repo.GetBulkPurchaseOrderParams{
	// 	BulkPurchaseOrderID: task.BulkPurchaseOrderID,
	// })
	// if err != nil {
	// 	return err
	// }

	// if bulk.PurchaseOrder != nil && bulk.PurchaseOrder.HubspotDealID != "" {
	// 	deal, err := workerInstance.App.HubspotClient.UpdateDeal(bulk.HubspotDealID, &hubspot.DealPropertiesForm{
	// 		Dealstage: hubspot.DealStageSample,
	// 	})
	// 	if err != nil {
	// 		return eris.Wrap(err, err.Error())
	// 	}

	// 	workerInstance.Logger.Debugf("Update deal %s success", deal.ID)
	// }
	// if err != nil {
	// 	return eris.Wrap(err, err.Error())
	// }
	return err

}

// Dispatch dispatch event
func (task HubspotSyncBulkTask) Dispatch(ctx context.Context, opts ...asynq.Option) (*asynq.TaskInfo, error) {
	return workerInstance.SendTaskWithContext(ctx, task, opts...)
}
