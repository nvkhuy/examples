package tasks

import (
	"context"
	"encoding/json"
	"github.com/engineeringinflow/inflow-backend/pkg/models"
	"github.com/engineeringinflow/inflow-backend/pkg/models/enums"
	"github.com/engineeringinflow/inflow-backend/pkg/zalo"
	"github.com/hibiken/asynq"
)

type RefreshTokenZaloTask struct{}

func (task RefreshTokenZaloTask) GetPayload() []byte {
	data, _ := json.Marshal(&task)
	return data
}

// TaskName task name
func (task RefreshTokenZaloTask) TaskName() string {
	return "refresh_token_zalo"
}

// Handler handler
func (task RefreshTokenZaloTask) Handler(ctx context.Context, t *asynq.Task) (err error) {
	if workerInstance.App.Config.IsProd() {
		return
	}

	var config models.ZaloConfig
	if err = workerInstance.App.DB.Model(&models.ZaloConfig{}).
		Where("id = ?", enums.ZaloConfigKeyRefreshToken).
		First(&config).Error; err != nil {
		return
	}

	_, err = zalo.OARefreshToken(workerInstance.App.DB, zalo.OARefreshTokenParams{})
	return
}

// Dispatch dispatch event
func (task RefreshTokenZaloTask) Dispatch(ctx context.Context, opts ...asynq.Option) (*asynq.TaskInfo, error) {
	return workerInstance.SendTaskWithContext(ctx, task, opts...)
}
