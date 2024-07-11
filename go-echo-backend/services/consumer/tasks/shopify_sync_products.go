package tasks

// import (
// 	"context"
// 	"time"

// 	"encoding/json"

// 	goshopify "github.com/bold-commerce/go-shopify/v3"
// 	"github.com/engineeringinflow/inflow-backend/pkg/models"
// 	"github.com/engineeringinflow/inflow-backend/pkg/repo"
// 	"github.com/engineeringinflow/inflow-backend/pkg/shopify"
// 	"github.com/hibiken/asynq"
// )

// type ShopifySyncProductsTask struct {
// 	CreatedAtMin int64 `json:"created_at_min" validate:"required"`
// 	CreatedAtMax int64 `json:"created_at_max" validate:"required,gt=CreatedAtMin"`
// }

// func (task ShopifySyncProductsTask) GetPayload() []byte {
// 	data, _ := json.Marshal(&task)

// 	return data
// }

// // TaskName task name
// func (task ShopifySyncProductsTask) TaskName() string {
// 	return "shopify_sync_products"
// }

// // Handler handler
// func (task ShopifySyncProductsTask) Handler(ctx context.Context, t *asynq.Task) error {
// 	var err = workerInstance.BindAndValidate(t.Payload(), &task)
// 	if err != nil {
// 		return err
// 	}
// 	var shopChannels []*models.ShopChannel
// 	workerInstance.App.DB.Find(&shopChannels, "coalesce(shop_name,'') <> '' AND coalesce(token,'') <> ''")

// 	for _, shopChannel := range shopChannels {
// 		products, err := shopify.GetInstance().NewClient(shopChannel.ShopName, shopChannel.Token).Product.List(goshopify.ProductListOptions{
// 			ListOptions: goshopify.ListOptions{
// 				CreatedAtMin: time.Unix(task.CreatedAtMin, 0),
// 				CreatedAtMax: time.Unix(task.CreatedAtMax, 0),
// 			},
// 		})
// 		if err != nil {
// 			return err
// 		}

// 		for _, product := range products {
// 			repo.NewShopifyMappingRepo(workerInstance.App.DB).UpdateOrCreateProduct(repo.UpdateOrCreateProductParams{
// 				ShopChannel:    shopChannel,
// 				ShopifyProduct: &product,
// 			})

// 		}
// 		workerInstance.Logger.Debugf("%d shop products are synced", len(products))
// 	}

// 	workerInstance.Logger.Debugf("%d shop channels are synced", len(shopChannels))
// 	return err
// }

// // Dispatch dispatch event
// func (task ShopifySyncProductsTask) Dispatch(ctx context.Context, opts ...asynq.Option) (*asynq.TaskInfo, error) {
// 	return workerInstance.SendTaskWithContext(ctx, task, opts...)
// }
