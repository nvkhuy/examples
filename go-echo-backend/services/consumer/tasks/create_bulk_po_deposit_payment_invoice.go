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

type CreateBulkPoDepositPaymentInvoiceTask struct {
	BulkPurchaseOrderID string `json:"bulk_purchase_order_id" validate:"required"`
	ApprovedByUserID    string `json:"approved_by_user_id"`
	ReCreate            bool   `json:"re_create"`
}

func (task CreateBulkPoDepositPaymentInvoiceTask) GetPayload() []byte {
	data, _ := json.Marshal(&task)
	return data
}

// TaskName task name
func (task CreateBulkPoDepositPaymentInvoiceTask) TaskName() string {
	return "create_bulk_po_deposit_payment_invoice"
}

// Handler handler
func (task CreateBulkPoDepositPaymentInvoiceTask) Handler(ctx context.Context, t *asynq.Task) error {
	var err = workerInstance.BindAndValidate(t.Payload(), &task)
	if err != nil {
		return err
	}

	bulkPO, err := repo.NewInvoiceRepo(workerInstance.App.DB).CreateBulkDepositInvoice(repo.CreateBulkDepositInvoiceParams{
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
		Event:  customerio.EventBulkPoBuyerDepositSucceeded,
		Data: bulkPO.GetCustomerIOMetadata(map[string]interface{}{
			"receipt_url": bulkPO.DepositReceiptURL,
			"invoice":     bulkPO.DepositInvoice.GetCustomerIOMetadata(),
		}),
	}.Dispatch(ctx)

	for _, assigneeID := range bulkPO.AssigneeIDs {
		TrackCustomerIOTask{
			UserID: assigneeID,
			Event:  customerio.EventBulkPoDepositSucceeded,
			Data: bulkPO.GetCustomerIOMetadata(map[string]interface{}{
				"receipt_url": bulkPO.DepositReceiptURL,
				"invoice":     bulkPO.DepositInvoice.GetCustomerIOMetadata(),
			}),
		}.Dispatch(ctx)
	}

	var eventData = bulkPO.GetCustomerIOMetadata(nil)
	t.ResultWriter().Write(helper.ToJson(&eventData))

	return err
}

// Dispatch dispatch event
func (task CreateBulkPoDepositPaymentInvoiceTask) Dispatch(ctx context.Context, opts ...asynq.Option) (*asynq.TaskInfo, error) {
	return workerInstance.SendTaskWithContext(ctx, task, opts...)
}
