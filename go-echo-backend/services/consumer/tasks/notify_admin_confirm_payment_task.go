package tasks

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/engineeringinflow/inflow-backend/pkg/customerio"
	"github.com/engineeringinflow/inflow-backend/pkg/models"
	"github.com/engineeringinflow/inflow-backend/pkg/repo"
	"github.com/hibiken/asynq"
	"github.com/samber/lo"
)

type NotifyAdminConfirmPaymentTask struct {
	PaymentTransactionID string `json:"payment_transaction_id" validate:"required"`
}

func (task NotifyAdminConfirmPaymentTask) GetPayload() []byte {
	data, _ := json.Marshal(&task)
	return data
}

// TaskName task name
func (task NotifyAdminConfirmPaymentTask) TaskName() string {
	return "notify_admin_confirm_payment_task"
}

// Handler handler
func (task NotifyAdminConfirmPaymentTask) Handler(ctx context.Context, t *asynq.Task) error {
	if err := workerInstance.BindAndValidate(t.Payload(), &task); err != nil {
		return err
	}
	payment, err := repo.NewPaymentTransactionRepo(workerInstance.App.DB).
		GetPaymentTransaction(repo.GetPaymentTransactionsParams{PaymentTransactionID: task.PaymentTransactionID, IncludeDetails: true})
	if err != nil {
		return err
	}
	var assigneeIDs []string
	var orderReferenceIDs []string

	if len(payment.PurchaseOrders) > 0 {
		for _, po := range payment.PurchaseOrders {
			assigneeIDs = append(assigneeIDs, po.AssigneeIDs...)
			orderReferenceIDs = append(orderReferenceIDs, po.ReferenceID)
		}
	}
	if len(payment.BulkPurchaseOrders) > 0 {
		for _, bpo := range payment.BulkPurchaseOrders {
			assigneeIDs = append(assigneeIDs, bpo.AssigneeIDs...)
			orderReferenceIDs = append(orderReferenceIDs, bpo.ReferenceID)
		}
	}

	assigneeIDs = lo.Uniq(assigneeIDs)
	var transferAt = time.Unix(payment.CreatedAt, 0).Format(time.RFC1123)
	var attachment *models.Attachment
	if payment.Attachments != nil {
		attachments := *payment.Attachments
		attachment = attachments[0]
	}

	for _, assigneeID := range assigneeIDs {
		TrackCustomerIOTask{
			Event:  customerio.EventBuyerCheckoutThroughBankTransfer,
			UserID: assigneeID,
			Data: map[string]interface{}{
				"buyer_name":           payment.User.Name,
				"payment_reference_id": payment.ReferenceID,
				"payment_link":         fmt.Sprintf("%s/payments/%s/overview", workerInstance.App.Config.AdminPortalBaseURL, payment.ID),
				"transfer_at":          transferAt,
				"order_items":          strings.Join(orderReferenceIDs, ", "),
				"payment_attachment":   attachment.GenerateFileURL(),
			},
		}.Dispatch(ctx)
	}

	return nil
}

// Dispatch dispatch event
func (task NotifyAdminConfirmPaymentTask) Dispatch(ctx context.Context, opts ...asynq.Option) (*asynq.TaskInfo, error) {
	return workerInstance.SendTaskWithContext(ctx, task, opts...)
}
