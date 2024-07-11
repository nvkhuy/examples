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
	"github.com/samber/lo"
)

type CreatePOPaymentInvoiceForMultipleItemsTask struct {
	ReCreate          bool   `json:"re_create"`
	CheckoutSessionID string `json:"checkout_session_id" validate:"required"`
	ApprovedByUserID  string `json:"approved_by_user_id"`
}

func (task CreatePOPaymentInvoiceForMultipleItemsTask) GetPayload() []byte {
	data, _ := json.Marshal(&task)
	return data
}

// TaskName task name
func (task CreatePOPaymentInvoiceForMultipleItemsTask) TaskName() string {
	return "create_po_payment_invoice_for_multiple_items"
}

// Handler handler
func (task CreatePOPaymentInvoiceForMultipleItemsTask) Handler(ctx context.Context, t *asynq.Task) error {
	var err = workerInstance.BindAndValidate(t.Payload(), &task)
	if err != nil {
		return err
	}

	resp, err := repo.NewInvoiceRepo(workerInstance.App.DB).CreateMultiplePurchaseOrderInvoice(repo.CreateMultiplePurchaseOrderInvoiceParams{
		CheckoutSessionID: task.CheckoutSessionID,
		ReCreate:          task.ReCreate,
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
			for _, purchaseOrder := range resp.PurchaseOrders {
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
	}

	if len(resp.PurchaseOrders) > 0 {
		TrackCustomerIOTask{
			UserID: resp.PurchaseOrders[0].UserID,
			Event:  customerio.EventPoBuyerPaymentSucceeded,
			Data: resp.PaymentTransaction.GetCustomerIOMetadata(map[string]interface{}{
				"receipt_url": resp.PurchaseOrders[0].ReceiptURL,
			}),
		}.Dispatch(ctx)
	}

	for _, purchaseOrder := range resp.PurchaseOrders {
		TrackCustomerIOTask{
			UserID: workerInstance.App.Config.InflowMerchandiseGroupEmail,
			Event:  customerio.EventPoCreated,
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

	}

	var eventData = lo.Map(resp.PurchaseOrders, func(item *models.PurchaseOrder, index int) map[string]interface{} {
		return item.GetCustomerIOMetadata(nil)
	})
	t.ResultWriter().Write(helper.ToJson(&eventData))

	return err
}

// Dispatch dispatch event
func (task CreatePOPaymentInvoiceForMultipleItemsTask) Dispatch(ctx context.Context, opts ...asynq.Option) (*asynq.TaskInfo, error) {
	return workerInstance.SendTaskWithContext(ctx, task, opts...)
}
