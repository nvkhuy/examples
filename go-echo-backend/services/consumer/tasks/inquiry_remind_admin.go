package tasks

import (
	"context"
	"fmt"

	"encoding/json"

	"github.com/engineeringinflow/inflow-backend/pkg/customerio"
	"github.com/engineeringinflow/inflow-backend/pkg/helper"
	"github.com/engineeringinflow/inflow-backend/pkg/repo"
	"github.com/hibiken/asynq"
	"github.com/rotisserie/eris"
)

type InquiryRemindAdminTask struct {
}

func (task InquiryRemindAdminTask) GetPayload() []byte {
	data, _ := json.Marshal(&task)

	return data
}

// TaskName task name
func (task InquiryRemindAdminTask) TaskName() string {
	return "inquiry_remind_admin"
}

// Handler handler
func (task InquiryRemindAdminTask) Handler(ctx context.Context, t *asynq.Task) error {
	var err = workerInstance.BindAndValidate(t.Payload(), &task)
	if err != nil {
		return err
	}

	var inquiries = repo.NewInquiryRepo(workerInstance.App.DB).GetInquiryRemindAdmin()

	if len(inquiries) == 0 {
		workerInstance.Logger.Debugf("No inquiries to remind admin")
		return nil
	}

	var eventData = map[string]interface{}{
		"inquiries":            inquiries.GetCustomerIOMetadata(),
		"total_inquiries":      len(inquiries),
		"admin_inquiries_link": fmt.Sprintf("%s/inquiries", workerInstance.App.Config.AdminPortalBaseURL),
	}

	_, _ = TrackCustomerIOTask{
		UserID: workerInstance.Config.App.Config.InflowSaleGroupEmail,
		Event:  customerio.EventNewInquiryRemindAdmin,
		Data:   eventData,
	}.Dispatch(ctx)

	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	t.ResultWriter().Write(helper.ToJson(&eventData))

	return err

}

// Dispatch dispatch event
func (task InquiryRemindAdminTask) Dispatch(ctx context.Context, opts ...asynq.Option) (*asynq.TaskInfo, error) {
	return workerInstance.SendTaskWithContext(ctx, task, opts...)
}
