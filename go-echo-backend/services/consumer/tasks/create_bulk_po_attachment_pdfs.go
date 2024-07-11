package tasks

import (
	"context"
	"encoding/json"

	"github.com/engineeringinflow/inflow-backend/pkg/errs"
	"github.com/engineeringinflow/inflow-backend/pkg/helper"
	"github.com/engineeringinflow/inflow-backend/pkg/repo"
	"github.com/hibiken/asynq"
	"github.com/rotisserie/eris"
)

type CreateBulkPoAttachmentPDFsTask struct {
	BulkPurchaseOrderID string `json:"bulk_purchase_order_id" validate:"required"`
	ReCreate            bool   `json:"re_create"`
}

func (task CreateBulkPoAttachmentPDFsTask) GetPayload() []byte {
	data, _ := json.Marshal(&task)
	return data
}

// TaskName task name
func (task CreateBulkPoAttachmentPDFsTask) TaskName() string {
	return "create_bulk_po_attachment_pdfs"
}

// Handler handler
func (task CreateBulkPoAttachmentPDFsTask) Handler(ctx context.Context, t *asynq.Task) (err error) {
	err = workerInstance.BindAndValidate(t.Payload(), &task)
	if err != nil {
		return
	}

	bulk, err := repo.NewInvoiceRepo(workerInstance.App.DB).CreateBulkDebitNotes(repo.CreateBulkDebitNotesParams{
		BulkPurchaseOrderID: task.BulkPurchaseOrderID,
		ReCreate:            task.ReCreate,
	})
	if err != nil {
		if eris.Is(err, errs.ErrBulkPoInvoiceAlreadyGenerated) {
			return nil
		}
		return err
	}

	var results = bulk.GetCustomerIOMetadata(nil)

	t.ResultWriter().Write(helper.ToJson(&results))

	return err
}

// Dispatch dispatch event
func (task CreateBulkPoAttachmentPDFsTask) Dispatch(ctx context.Context, opts ...asynq.Option) (*asynq.TaskInfo, error) {
	return workerInstance.SendTaskWithContext(ctx, task, opts...)
}
