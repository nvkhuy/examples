package tasks

import (
	"context"
	"encoding/json"

	"github.com/engineeringinflow/inflow-backend/pkg/models"
	"github.com/engineeringinflow/inflow-backend/pkg/runner"
	"github.com/hibiken/asynq"
)

type GenerateBlurTask struct {
	FileKeys []string `json:"file_keys"`
}

func (task GenerateBlurTask) GetPayload() []byte {
	data, _ := json.Marshal(&task)

	return data
}

// TaskName task name
func (task GenerateBlurTask) TaskName() string {
	return "generate_blur"
}

// Handler handler
func (task GenerateBlurTask) Handler(ctx context.Context, t *asynq.Task) error {
	var err = workerInstance.BindAndValidate(t.Payload(), &task)
	if err != nil {
		return err
	}
	var products []*models.Product
	workerInstance.App.DB.Select("ID", "Attachments").Find(&products, "attachments is not null")

	var runner = runner.New(5)

	for index, product := range products {
		index := index
		product := product

		runner.Submit(func() {
			if product.Attachments != nil && len(*product.Attachments) > 0 {
				var attachments models.Attachments
				for _, attachment := range *product.Attachments {
					attachment.GetBlurhash()
					attachments = append(attachments, attachment)

				}
				var err = workerInstance.App.DB.Model(&models.Product{}).Where("id = ?", product.ID).Updates(map[string]interface{}{"attachments": attachments}).Error
				workerInstance.App.DB.CustomLogger.Debugf("Update %d/%d product attachments err=%+v", index, len(products), err)
			}
		})

	}

	var blogs []*models.Post
	workerInstance.App.DB.Select("ID", "FeaturedImage").Find(&blogs, "attachments is not null")

	for index, blog := range blogs {
		index := index
		blog := blog

		runner.Submit(func() {
			if blog.FeaturedImage != nil {
				var featuredImage = blog.FeaturedImage.GetBlurhash()
				var err = workerInstance.App.DB.Model(&models.Post{}).Where("id = ?", blog.ID).Updates(map[string]interface{}{"featured_image": featuredImage}).Error
				workerInstance.App.DB.CustomLogger.Debugf("Update %d/%d blog attachments err=%+v", index, len(products), err)
			}
		})

	}

	runner.Wait()

	return err
}

// Dispatch dispatch event
func (task GenerateBlurTask) Dispatch(ctx context.Context, opts ...asynq.Option) (*asynq.TaskInfo, error) {
	return workerInstance.SendTaskWithContext(ctx, task, opts...)
}
