package tasks

import (
	"context"
	"encoding/json"

	"github.com/engineeringinflow/inflow-backend/pkg/customerio"
	"github.com/engineeringinflow/inflow-backend/pkg/errs"
	"github.com/engineeringinflow/inflow-backend/pkg/helper"
	"github.com/engineeringinflow/inflow-backend/pkg/repo"
	"github.com/hibiken/asynq"
	"github.com/rotisserie/eris"
)

type CreateBulkPoFirstPaymentInvoiceTask struct {
	BulkPurchaseOrderID string `json:"bulk_purchase_order_id" validate:"required"`
	ApprovedByUserID    string `json:"approved_by_user_id"`
	ReCreate            bool   `json:"re_create"`
}

func (task CreateBulkPoFirstPaymentInvoiceTask) GetPayload() []byte {
	data, _ := json.Marshal(&task)
	return data
}

// TaskName task name
func (task CreateBulkPoFirstPaymentInvoiceTask) TaskName() string {
	return "create_bulk_po_first_payment_invoice"
}

// Handler handler
func (task CreateBulkPoFirstPaymentInvoiceTask) Handler(ctx context.Context, t *asynq.Task) error {
	var err = workerInstance.BindAndValidate(t.Payload(), &task)
	if err != nil {
		return err
	}

	bulkPO, err := repo.NewInvoiceRepo(workerInstance.App.DB).CreateBulkFirstPaymentInvoice(repo.CreateBulkFirstPaymentInvoiceParams{
		BulkPurchaseOrderID: task.BulkPurchaseOrderID,
		ReCreate:            task.ReCreate,
	})
	if err != nil {
		if eris.Is(err, errs.ErrBulkPoInvoiceAlreadyGenerated) {
			return nil
		}
		return err
	}

	TrackCustomerIOTask{
		UserID: bulkPO.UserID,
		Event:  customerio.EventBulkPoBuyerFirstPaymentSucceeded,
		Data: bulkPO.GetCustomerIOMetadata(map[string]interface{}{
			"receipt_url": bulkPO.FirstPaymentReceiptURL,
			"invoice":     bulkPO.FirstPaymentInvoice.GetCustomerIOMetadata(),
		}),
	}.Dispatch(ctx)

	for _, assigneeID := range bulkPO.AssigneeIDs {
		TrackCustomerIOTask{
			UserID: assigneeID,
			Event:  customerio.EventBulkPoFirstPaymentSucceeded,
			Data: bulkPO.GetCustomerIOMetadata(map[string]interface{}{
				"receipt_url": bulkPO.FirstPaymentReceiptURL,
				"invoice":     bulkPO.FirstPaymentInvoice.GetCustomerIOMetadata(),
			}),
		}.Dispatch(ctx)
	}

	var eventData = bulkPO.GetCustomerIOMetadata(nil)
	t.ResultWriter().Write(helper.ToJson(&eventData))

	return err
}

// Dispatch dispatch event
func (task CreateBulkPoFirstPaymentInvoiceTask) Dispatch(ctx context.Context, opts ...asynq.Option) (*asynq.TaskInfo, error) {
	return workerInstance.SendTaskWithContext(ctx, task, opts...)
}
