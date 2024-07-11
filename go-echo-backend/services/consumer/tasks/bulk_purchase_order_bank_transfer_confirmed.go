package tasks

import (
	"context"

	"github.com/engineeringinflow/inflow-backend/pkg/db"
	"github.com/engineeringinflow/inflow-backend/pkg/models"

	"encoding/json"

	"github.com/engineeringinflow/inflow-backend/pkg/models/enums"
	"github.com/hibiken/asynq"
)

type BulkPurchaseOrderBankTransferConfirmedTask struct {
	ApprovedByUserID    string                 `json:"aprroved_by_user_id" validate:"required"`
	BulkPurchaseOrderID string                 `json:"bulk_purchase_order_id" validate:"required"`
	Milestone           enums.PaymentMilestone `json:"milestone" validate:"required"`
}

func (task BulkPurchaseOrderBankTransferConfirmedTask) GetPayload() []byte {
	data, _ := json.Marshal(&task)

	return data
}

// TaskName task name
func (task BulkPurchaseOrderBankTransferConfirmedTask) TaskName() string {
	return "bulk_purchase_order_bank_transfer_confirmed"
}

// Handler handler
func (task BulkPurchaseOrderBankTransferConfirmedTask) Handler(ctx context.Context, t *asynq.Task) error {
	var err = workerInstance.BindAndValidate(t.Payload(), &task)
	if err != nil {
		return err
	}

	switch task.Milestone {
	case enums.PaymentMilestoneFirstPayment:
		CreateBulkPoFirstPaymentInvoiceTask{
			ApprovedByUserID:    task.ApprovedByUserID,
			BulkPurchaseOrderID: task.BulkPurchaseOrderID,
		}.Dispatch(ctx)

	case enums.PaymentMilestoneFinalPayment:
		CreateBulkPoFinalPaymentInvoiceTask{
			ApprovedByUserID:    task.ApprovedByUserID,
			BulkPurchaseOrderID: task.BulkPurchaseOrderID,
			ReCreate:            true,
		}.Dispatch(ctx)
	}

	err = task.updatePaymentTransactions(workerInstance.App.DB)
	return err

}

func (task BulkPurchaseOrderBankTransferConfirmedTask) updatePaymentTransactions(db *db.DB) (err error) {
	err = db.Model(&models.PaymentTransaction{}).
		Where("bulk_purchase_order_id = ?", task.BulkPurchaseOrderID).
		Where("milestone = ?", task.Milestone).
		Update("status", enums.PaymentStatusPaid).Error
	return
}

// Dispatch dispatch event
func (task BulkPurchaseOrderBankTransferConfirmedTask) Dispatch(ctx context.Context, opts ...asynq.Option) (*asynq.TaskInfo, error) {
	return workerInstance.SendTaskWithContext(ctx, task, opts...)
}
