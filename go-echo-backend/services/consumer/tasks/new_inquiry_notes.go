package tasks

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/engineeringinflow/inflow-backend/pkg/customerio"
	"github.com/engineeringinflow/inflow-backend/pkg/models"
	"github.com/engineeringinflow/inflow-backend/pkg/models/enums"
	"github.com/engineeringinflow/inflow-backend/pkg/repo"
	"github.com/hibiken/asynq"
)

type NewInquiryNotesTask struct {
	UserID         string              `json:"user_id" validate:"required"`
	InquiryID      string              `json:"inquiry_id" validate:"required"`
	MentionUserIDs []string            `json:"mention_user_ids" validate:"required"`
	Message        string              `json:"message"`
	Attachments    *models.Attachments `json:"attachments"`
}

func (task NewInquiryNotesTask) GetPayload() []byte {
	data, _ := json.Marshal(&task)

	return data
}

// TaskName task name
func (task NewInquiryNotesTask) TaskName() string {
	return "new_inquiry_notes"
}

// Handler handler
func (task NewInquiryNotesTask) Handler(ctx context.Context, t *asynq.Task) error {
	var err = workerInstance.BindAndValidate(t.Payload(), &task)
	if err != nil {
		return err
	}

	inquiry, err := repo.NewInquiryRepo(workerInstance.App.DB).GetInquiryByID(repo.GetInquiryByIDParams{
		InquiryID:     task.InquiryID,
		JwtClaimsInfo: *models.NewJwtClaimsInfo().SetRole(enums.RoleSuperAdmin).SetUserID(task.UserID),
	})
	if err != nil {
		return err
	}

	var sender models.User
	err = workerInstance.App.DB.Select("ID", "Name", "Email").First(&sender, "id = ?", task.UserID).Error
	if err != nil {
		return err
	}

	var extras = map[string]interface{}{
		"sender":                  sender.GetCustomerIOMetadata(nil),
		"admin_inquiry_notes_url": fmt.Sprintf("%s/inquiries/%s/customer/notes", workerInstance.App.Config.AdminPortalBaseURL, inquiry.ID),
	}
	if task.Message != "" {
		extras["message"] = task.Message
	}

	if task.Attachments != nil {
		extras["attachments"] = task.Attachments
	}

	for _, userID := range task.MentionUserIDs {
		SyncCustomerIOUserTask{
			UserID: userID,
		}.Dispatch(ctx)

		TrackCustomerIOTask{
			UserID: userID,
			Event:  customerio.EventAdminInquiryNewNotes,
			Data:   inquiry.GetCustomerIOMetadata(extras),
		}.DispatchIn(time.Second * 3)

	}

	return err
}

// Dispatch dispatch event
func (task NewInquiryNotesTask) Dispatch(ctx context.Context, opts ...asynq.Option) (*asynq.TaskInfo, error) {
	return workerInstance.SendTaskWithContext(ctx, task, opts...)
}
