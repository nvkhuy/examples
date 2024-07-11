package tasks

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/engineeringinflow/inflow-backend/pkg/customerio"
	"github.com/engineeringinflow/inflow-backend/pkg/models"
	"github.com/engineeringinflow/inflow-backend/pkg/models/enums"
	"github.com/engineeringinflow/inflow-backend/pkg/repo"
	"github.com/hibiken/asynq"
)

type CreatePaymentInvoiceTask struct {
	PaymentTransactionID string `json:"payment_transaction_id" validate:"required"`
	ApprovedByUserID     string `json:"approved_by_user_id"`
}

func (task CreatePaymentInvoiceTask) GetPayload() []byte {
	data, _ := json.Marshal(&task)
	return data
}

// TaskName task name
func (task CreatePaymentInvoiceTask) TaskName() string {
	return "create_payment_invoice_task"
}

// Handler handler
func (task CreatePaymentInvoiceTask) Handler(ctx context.Context, t *asynq.Task) error {
	if err := workerInstance.BindAndValidate(t.Payload(), &task); err != nil {
		return err
	}
	invoice, err := repo.NewInvoiceRepo(workerInstance.App.DB).CreatePaymentInvoice(&repo.CreateOrderInvoiceRequest{PaymentTransactionID: task.PaymentTransactionID})
	if err != nil {
		return err
	}
	payment, err := repo.NewPaymentTransactionRepo(workerInstance.App.DB).
		GetPaymentTransaction(repo.GetPaymentTransactionsParams{PaymentTransactionID: task.PaymentTransactionID, IncludeDetails: true})
	if err != nil {
		return err
	}
	var orderReferenceIDs []string

	if len(payment.PurchaseOrders) > 0 {
		// card
		var description = "The payment has been received"
		//bank transfer
		if task.ApprovedByUserID != "" {
			var admin models.User
			if err := workerInstance.App.DB.Select("Name").First(&admin, "id = ?", task.ApprovedByUserID).Error; err != nil {
				return err
			}
			description = fmt.Sprintf("Admin %s has confirmed the payment", admin.Name)
		}
		for _, po := range payment.PurchaseOrders {
			if po.InquiryID != "" {
				CreateInquiryAuditTask{
					Form: models.InquiryAuditCreateForm{
						InquiryID:   po.InquiryID,
						ActionType:  enums.AuditActionTypeInquiryAdminMarkAsPaid,
						UserID:      task.ApprovedByUserID,
						Description: description,
					},
				}.Dispatch(ctx)
			}

			orderReferenceIDs = append(orderReferenceIDs, po.ReferenceID)
		}
	}
	if len(payment.BulkPurchaseOrders) > 0 {
		for _, bpo := range payment.BulkPurchaseOrders {
			orderReferenceIDs = append(orderReferenceIDs, bpo.ReferenceID)

		}
	}

	TrackCustomerIOTask{
		Event:  customerio.EventAdminConfirmPaymentReceived,
		UserID: invoice.UserID,
		Data: map[string]interface{}{
			"order_reference_ids": strings.Join(orderReferenceIDs, ", "),
			"invoice_attachment":  invoice.Document.GenerateFileURL(),
		},
	}.Dispatch(ctx)

	return nil
}

// Dispatch dispatch event
func (task CreatePaymentInvoiceTask) Dispatch(ctx context.Context, opts ...asynq.Option) (*asynq.TaskInfo, error) {
	return workerInstance.SendTaskWithContext(ctx, task, opts...)
}
