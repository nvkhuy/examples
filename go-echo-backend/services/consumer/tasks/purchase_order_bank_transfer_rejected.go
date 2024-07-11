package tasks

import (
	"context"
	"fmt"

	"encoding/json"

	"github.com/engineeringinflow/inflow-backend/pkg/customerio"
	"github.com/engineeringinflow/inflow-backend/pkg/models"
	"github.com/engineeringinflow/inflow-backend/pkg/models/enums"
	"github.com/engineeringinflow/inflow-backend/pkg/repo"
	"github.com/hibiken/asynq"
)

type PurchaseOrderBankTransferRejectedTask struct {
	ApprovedByUserID string `json:"aprroved_by_user_id" validate:"required"`
	PurchaseOrderID  string `json:"purchase_order_id" validate:"required"`
}

func (task PurchaseOrderBankTransferRejectedTask) GetPayload() []byte {
	data, _ := json.Marshal(&task)

	return data
}

// TaskName task name
func (task PurchaseOrderBankTransferRejectedTask) TaskName() string {
	return "purchase_order_bank_transfer_rejected"
}

// Handler handler
func (task PurchaseOrderBankTransferRejectedTask) Handler(ctx context.Context, t *asynq.Task) error {
	var err = workerInstance.BindAndValidate(t.Payload(), &task)
	if err != nil {
		return err
	}

	userAdmin, err := repo.NewUserRepo(workerInstance.App.DB).GetShortUserInfo(task.ApprovedByUserID)
	if err != nil {
		return err
	}

	purchaseOrder, err := repo.NewPurchaseOrderRepo(workerInstance.App.DB).GetPurchaseOrder(repo.GetPurchaseOrderParams{
		PurchaseOrderID: task.PurchaseOrderID,
	})
	if err != nil {
		return err
	}

	CreateInquiryAuditTask{
		Form: models.InquiryAuditCreateForm{
			InquiryID:   purchaseOrder.InquiryID,
			ActionType:  enums.AuditActionTypeInquiryAdminMarkAsUnPaid,
			UserID:      userAdmin.ID,
			Description: fmt.Sprintf("Admin %s has confirmed the payment", userAdmin.Name),
		},
	}.Dispatch(ctx)

	TrackCustomerIOTask{
		UserID: purchaseOrder.InquiryID,
		Event:  customerio.EventPoBuyerPaymentFailed,
		Data:   purchaseOrder.GetCustomerIOMetadata(nil),
	}.Dispatch(ctx)
	return err

}

// Dispatch dispatch event
func (task PurchaseOrderBankTransferRejectedTask) Dispatch(ctx context.Context, opts ...asynq.Option) (*asynq.TaskInfo, error) {
	return workerInstance.SendTaskWithContext(ctx, task, opts...)
}
