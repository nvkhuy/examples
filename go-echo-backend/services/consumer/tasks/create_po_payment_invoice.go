package tasks

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/engineeringinflow/inflow-backend/pkg/customerio"
	"github.com/engineeringinflow/inflow-backend/pkg/errs"
	"github.com/engineeringinflow/inflow-backend/pkg/helper"
	"github.com/engineeringinflow/inflow-backend/pkg/models"
	"github.com/engineeringinflow/inflow-backend/pkg/models/enums"
	"github.com/engineeringinflow/inflow-backend/pkg/repo"
	"github.com/hibiken/asynq"
	"github.com/rotisserie/eris"
)

type CreatePOPaymentInvoiceTask struct {
	PurchaseOrderID string `json:"purchase_order_id" validate:"required"`
	ReCreate        bool   `json:"re_create"`

	ApprovedByUserID string `json:"aprroved_by_user_id" ` // Bank transfer

}

func (task CreatePOPaymentInvoiceTask) GetPayload() []byte {
	data, _ := json.Marshal(&task)
	return data
}

// TaskName task name
func (task CreatePOPaymentInvoiceTask) TaskName() string {
	return "create_po_payment_invoice"
}

// Handler handler
func (task CreatePOPaymentInvoiceTask) Handler(ctx context.Context, t *asynq.Task) error {
	var err = workerInstance.BindAndValidate(t.Payload(), &task)
	if err != nil {
		return err
	}

	purchaseOrder, err := repo.NewInvoiceRepo(workerInstance.App.DB).CreatePurchaseOrderInvoice(repo.CreatePurchaseOrderInvoiceParams{
		PurchaseOrderID: task.PurchaseOrderID,
		ReCreate:        task.ReCreate,
	})
	if err != nil {
		if eris.Is(err, errs.ErrPOInvoiceAlreadyGenerated) {
			return nil
		}
		return err
	}

	if task.ApprovedByUserID != "" {
		userAdmin, err := repo.NewUserRepo(workerInstance.App.DB).GetShortUserInfo(task.ApprovedByUserID)
		if err == nil && userAdmin != nil {
			CreateInquiryAuditTask{
				Form: models.InquiryAuditCreateForm{
					InquiryID:   purchaseOrder.InquiryID,
					ActionType:  enums.AuditActionTypeInquiryAdminMarkAsPaid,
					UserID:      userAdmin.ID,
					Description: fmt.Sprintf("Admin %s has confirmed the payment", userAdmin.Name),
				},
			}.Dispatch(ctx)
		}
	}

	TrackCustomerIOTask{
		UserID: purchaseOrder.UserID,
		Event:  customerio.EventPoBuyerPaymentSucceeded,
		Data: purchaseOrder.GetCustomerIOMetadata(map[string]interface{}{
			"receipt_url": purchaseOrder.ReceiptURL,
		}),
	}.Dispatch(ctx)

	for _, assigneeID := range purchaseOrder.Inquiry.AssigneeIDs {
		TrackCustomerIOTask{
			UserID: assigneeID,
			Event:  customerio.EventPoWorkingOnDesign,
			Data: purchaseOrder.GetCustomerIOMetadata(map[string]interface{}{
				"receipt_url": purchaseOrder.ReceiptURL,
			}),
		}.Dispatch(ctx)
	}

	TrackCustomerIOTask{
		UserID: workerInstance.App.Config.InflowMerchandiseGroupEmail,
		Event:  customerio.EventPoCreated,
		Data: purchaseOrder.GetCustomerIOMetadata(map[string]interface{}{
			"receipt_url": purchaseOrder.ReceiptURL,
		}),
	}.Dispatch(ctx)

	var eventData = purchaseOrder.GetCustomerIOMetadata(nil)
	t.ResultWriter().Write(helper.ToJson(&eventData))

	return err
}

// Dispatch dispatch event
func (task CreatePOPaymentInvoiceTask) Dispatch(ctx context.Context, opts ...asynq.Option) (*asynq.TaskInfo, error) {
	return workerInstance.SendTaskWithContext(ctx, task, opts...)
}
