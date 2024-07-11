package tasks

import (
	"context"

	"encoding/json"

	"github.com/hibiken/asynq"
	"github.com/thaitanloi365/go-sendgrid"

	"github.com/engineeringinflow/inflow-backend/pkg/mailer"
)

type SendMailTask struct {
	Email      string                 `json:"email" validate:"email"`
	TemplateID mailer.TemplateID      `json:"template_id" validate:"required"`
	Data       map[string]interface{} `json:"data"`
}

func (task SendMailTask) GetPayload() []byte {
	data, _ := json.Marshal(&task)

	return data
}

// TaskName task name
func (task SendMailTask) TaskName() string {
	return "send_mail"
}

// Handler handler
func (task SendMailTask) Handler(ctx context.Context, t *asynq.Task) error {
	var err = workerInstance.BindAndValidate(t.Payload(), &task)
	if err != nil {
		return err
	}

	err = workerInstance.App.Mailer.Send(sendgrid.SendMailParams{
		Email:      task.Email,
		TemplateID: string(task.TemplateID),
		Data:       task.Data,
	})

	return err
}

// Dispatch dispatch event
func (task SendMailTask) Dispatch(ctx context.Context, opts ...asynq.Option) (*asynq.TaskInfo, error) {
	return workerInstance.SendTaskWithContext(ctx, task, opts...)
}
