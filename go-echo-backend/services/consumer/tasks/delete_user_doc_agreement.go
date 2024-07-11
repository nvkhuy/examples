package tasks

import (
	"context"
	"encoding/json"
	"github.com/engineeringinflow/inflow-backend/pkg/models/enums"
	"github.com/engineeringinflow/inflow-backend/pkg/repo"
	"github.com/hibiken/asynq"
)

type DeleteUserDocAgreementTask struct {
	SettingDocType enums.SettingDoc `json:"setting_doc_type,omitempty"`
}

func (task DeleteUserDocAgreementTask) GetPayload() []byte {
	data, _ := json.Marshal(&task)

	return data
}

// TaskName task name
func (task DeleteUserDocAgreementTask) TaskName() string {
	return "delete_user_doc_agreement"
}

// Handler handler
func (task DeleteUserDocAgreementTask) Handler(ctx context.Context, t *asynq.Task) error {
	var err = workerInstance.BindAndValidate(t.Payload(), &task)
	if err != nil {
		return err
	}

	err = repo.NewUserDocAgreementRepo(workerInstance.App.DB).Delete(repo.DeleteUserDocAgreementTypeParams{
		SettingDocType: task.SettingDocType,
	})

	return err
}

// Dispatch dispatch event
func (task DeleteUserDocAgreementTask) Dispatch(ctx context.Context, opts ...asynq.Option) (*asynq.TaskInfo, error) {
	return workerInstance.SendTaskWithContext(ctx, task, opts...)
}
