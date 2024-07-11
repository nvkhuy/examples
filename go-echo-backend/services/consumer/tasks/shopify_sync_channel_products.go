package tasks

// import (
// 	"context"

// 	"encoding/json"

// 	"github.com/engineeringinflow/inflow-backend/pkg/models"
// 	"github.com/engineeringinflow/inflow-backend/pkg/repo"
// 	"github.com/engineeringinflow/inflow-backend/pkg/shopify"
// 	"github.com/hibiken/asynq"
// )

// type ShopifySyncChannelProductsTask struct {
// 	ClientInfo *shopify.ClientInfo `json:"client_info" validate:"required"`
// 	LocationID string              `json:"location_id" validate:"required"`
// }

// func (task ShopifySyncChannelProductsTask) GetPayload() []byte {
// 	data, _ := json.Marshal(&task)

// 	return data
// }

// // TaskName task name
// func (task ShopifySyncChannelProductsTask) TaskName() string {
// 	return "shopify_sync_channel_products"
// }

// // Handler handler
// func (task ShopifySyncChannelProductsTask) Handler(ctx context.Context, t *asynq.Task) error {
// 	var err = workerInstance.BindAndValidate(t.Payload(), &task)
// 	if err != nil {
// 		return err
// 	}
// 	products, err := shopify.GetInstance().NewClient(task.ClientInfo.ShopName, task.ClientInfo.Token).Product.List(nil)
// 	if err != nil {
// 		return err
// 	}

// 	var shopChannel models.ShopChannel
// 	err = workerInstance.App.DB.First(&shopChannel, "shop_name = ? AND token = ?", task.ClientInfo.ShopName, task.ClientInfo.Token).Error
// 	if err != nil {
// 		return err
// 	}

// 	for _, product := range products {
// 		repo.NewShopifyMappingRepo(workerInstance.App.DB).UpdateOrCreateProduct(repo.UpdateOrCreateProductParams{
// 			ShopChannel:    &shopChannel,
// 			ShopifyProduct: &product,
// 		})

// 	}
// 	workerInstance.Logger.Debugf("%d shop products are synced", len(products))

// 	return err
// }

// // Dispatch dispatch event
// func (task ShopifySyncChannelProductsTask) Dispatch(ctx context.Context, opts ...asynq.Option) (*asynq.TaskInfo, error) {
// 	return workerInstance.SendTaskWithContext(ctx, task, opts...)
// }
