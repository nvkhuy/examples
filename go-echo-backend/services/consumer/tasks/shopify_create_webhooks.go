package tasks

import (
	"context"

	"encoding/json"

	"github.com/hibiken/asynq"

	"github.com/engineeringinflow/inflow-backend/pkg/shopify"
)

type ShopifyCreateWebhooksTask struct {
	ClientInfo *shopify.ClientInfo `json:"client_info"`
}

func (task ShopifyCreateWebhooksTask) GetPayload() []byte {
	data, _ := json.Marshal(&task)

	return data
}

// TaskName task name
func (task ShopifyCreateWebhooksTask) TaskName() string {
	return "shopify_create_webhooks"
}

// Handler handler
func (task ShopifyCreateWebhooksTask) Handler(ctx context.Context, t *asynq.Task) error {
	var err = workerInstance.BindAndValidate(t.Payload(), &task)
	if err != nil {
		return err
	}

	err = shopify.GetInstance().CreateWebhooks(task.ClientInfo)

	return err
}

// Dispatch dispatch event
func (task ShopifyCreateWebhooksTask) Dispatch(ctx context.Context, opts ...asynq.Option) (*asynq.TaskInfo, error) {
	return workerInstance.SendTaskWithContext(ctx, task, opts...)
}
