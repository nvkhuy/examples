package tasks

import (
	"context"
	"github.com/engineeringinflow/inflow-backend/pkg/db"
	"github.com/engineeringinflow/inflow-backend/pkg/models"
	"github.com/engineeringinflow/inflow-backend/pkg/models/enums"

	"encoding/json"

	"github.com/hibiken/asynq"
)

type PurchaseOrderBankTransferConfirmedTask struct {
	ApprovedByUserID string `json:"aprroved_by_user_id" validate:"required"`
	PurchaseOrderID  string `json:"purchase_order_id" validate:"required"`
}

func (task PurchaseOrderBankTransferConfirmedTask) GetPayload() []byte {
	data, _ := json.Marshal(&task)

	return data
}

// TaskName task name
func (task PurchaseOrderBankTransferConfirmedTask) TaskName() string {
	return "purchase_order_bank_transfer_confirmed"
}

// Handler handler
func (task PurchaseOrderBankTransferConfirmedTask) Handler(ctx context.Context, t *asynq.Task) error {
	var err = workerInstance.BindAndValidate(t.Payload(), &task)
	if err != nil {
		return err
	}

	_, _ = CreatePOPaymentInvoiceTask{
		PurchaseOrderID:  task.PurchaseOrderID,
		ApprovedByUserID: task.ApprovedByUserID,
	}.Dispatch(ctx)
	err = task.updatePaymentTransactions(workerInstance.App.DB)
	return err
}

func (task PurchaseOrderBankTransferConfirmedTask) updatePaymentTransactions(db *db.DB) (err error) {
	err = db.Model(&models.PaymentTransaction{}).
		Where("purchase_order_id = ?", task.PurchaseOrderID).
		Update("status", enums.PaymentStatusPaid).Error
	return
}

// Dispatch dispatch event
func (task PurchaseOrderBankTransferConfirmedTask) Dispatch(ctx context.Context, opts ...asynq.Option) (*asynq.TaskInfo, error) {
	return workerInstance.SendTaskWithContext(ctx, task, opts...)
}
