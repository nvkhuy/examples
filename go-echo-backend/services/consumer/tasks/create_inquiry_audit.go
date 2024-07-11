package tasks

import (
	"context"

	"encoding/json"

	"github.com/engineeringinflow/inflow-backend/pkg/models"
	"github.com/engineeringinflow/inflow-backend/pkg/repo"
	"github.com/hibiken/asynq"
	"github.com/rotisserie/eris"
)

type CreateInquiryAuditTask struct {
	Form models.InquiryAuditCreateForm `json:"form" validate:"required"`
}

func (task CreateInquiryAuditTask) GetPayload() []byte {
	data, _ := json.Marshal(&task)

	return data
}

// TaskName task name
func (task CreateInquiryAuditTask) TaskName() string {
	return "create_inquiry_audit"
}

// Handler handler
func (task CreateInquiryAuditTask) Handler(ctx context.Context, t *asynq.Task) error {
	var err = workerInstance.BindAndValidate(t.Payload(), &task)
	if err != nil {
		return err
	}

	// InquiryAudit
	err = repo.NewInquiryAuditRepo(workerInstance.App.DB).CreateInquiryAudit(task.Form)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}
	return err

}

// Dispatch dispatch event
func (task CreateInquiryAuditTask) Dispatch(ctx context.Context, opts ...asynq.Option) (*asynq.TaskInfo, error) {
	return workerInstance.SendTaskWithContext(ctx, task, opts...)
}
