package tasks

// import (
// 	"context"
// 	"strconv"

// 	"encoding/json"

// 	"github.com/engineeringinflow/inflow-backend/pkg/models"
// 	"github.com/engineeringinflow/inflow-backend/pkg/repo"
// 	"github.com/engineeringinflow/inflow-backend/pkg/shopify"
// 	"github.com/hibiken/asynq"
// )

// type ShopifySyncProductTask struct {
// 	ProductID  string              `json:"product_id" validate:"required"`
// 	LocationID string              `json:"location_id" validate:"required"`
// 	ClientInfo *shopify.ClientInfo `json:"client_info" validate:"required"`
// }

// func (task ShopifySyncProductTask) GetPayload() []byte {
// 	data, _ := json.Marshal(&task)

// 	return data
// }

// // TaskName task name
// func (task ShopifySyncProductTask) TaskName() string {
// 	return "shopify_sync_product"
// }

// // Handler handler
// func (task ShopifySyncProductTask) Handler(ctx context.Context, t *asynq.Task) error {
// 	var err = workerInstance.BindAndValidate(t.Payload(), &task)
// 	if err != nil {
// 		return err
// 	}

// 	v, err := strconv.Atoi(task.ProductID)
// 	if err != nil {
// 		return err
// 	}

// 	var shopifyClient = shopify.GetInstance().NewClient(task.ClientInfo.ShopName, task.ClientInfo.Token)

// 	product, err := shopifyClient.Product.Get(int64(v), nil)
// 	if err != nil {
// 		return err
// 	}

// 	var shopChannel models.ShopChannel
// 	err = workerInstance.App.DB.First(&shopChannel, "shop_name = ? AND token = ?", task.ClientInfo.ShopName, task.ClientInfo.Token).Error
// 	if err != nil {
// 		return err
// 	}

// 	_, err = repo.NewShopifyMappingRepo(workerInstance.App.DB).UpdateOrCreateProduct(repo.UpdateOrCreateProductParams{
// 		ShopChannel:    &shopChannel,
// 		ShopifyProduct: product,
// 	})

// 	return err
// }

// // Dispatch dispatch event
// func (task ShopifySyncProductTask) Dispatch(ctx context.Context, opts ...asynq.Option) (*asynq.TaskInfo, error) {
// 	return workerInstance.SendTaskWithContext(ctx, task, opts...)
// }
