package tasks

import (
	"context"
	"fmt"

	"encoding/json"

	"github.com/engineeringinflow/inflow-backend/pkg/customerio"
	"github.com/engineeringinflow/inflow-backend/pkg/models"
	"github.com/engineeringinflow/inflow-backend/pkg/models/enums"
	"github.com/engineeringinflow/inflow-backend/pkg/repo"
	"github.com/engineeringinflow/inflow-backend/pkg/repo/query/queryfunc"
	"github.com/hibiken/asynq"
)

type BulkPurchaseOrderBankTransferRejectedTask struct {
	ApprovedByUserID string `json:"aprroved_by_user_id" validate:"required"`
	InquiryID        string `json:"inquiry_id"`
}

func (task BulkPurchaseOrderBankTransferRejectedTask) GetPayload() []byte {
	data, _ := json.Marshal(&task)

	return data
}

// TaskName task name
func (task BulkPurchaseOrderBankTransferRejectedTask) TaskName() string {
	return "bulk_purchase_order_bank_transfer_rejected"
}

// Handler handler
func (task BulkPurchaseOrderBankTransferRejectedTask) Handler(ctx context.Context, t *asynq.Task) error {
	var err = workerInstance.BindAndValidate(t.Payload(), &task)
	if err != nil {
		return err
	}

	userAdmin, err := repo.NewUserRepo(workerInstance.App.DB).GetShortUserInfo(task.ApprovedByUserID)
	if err != nil {
		return err
	}

	inquiry, err := repo.NewInquiryRepo(workerInstance.App.DB).GetInquiryByID(repo.GetInquiryByIDParams{
		InquiryID: task.InquiryID,
		InquiryBuilderOptions: queryfunc.InquiryBuilderOptions{
			IncludePurchaseOrder: true,
		},
		JwtClaimsInfo: *models.NewJwtClaimsInfo().SetRole(enums.RoleSuperAdmin).SetUserID(task.ApprovedByUserID),
	})
	if err != nil {
		return err
	}

	CreateInquiryAuditTask{
		Form: models.InquiryAuditCreateForm{
			InquiryID:   task.InquiryID,
			ActionType:  enums.AuditActionTypeInquiryAdminMarkAsPaid,
			UserID:      userAdmin.ID,
			Description: fmt.Sprintf("Admin %s has confirmed the payment", userAdmin.Name),
		},
	}.Dispatch(ctx)

	var purchaseOrder = inquiry.PurchaseOrder
	inquiry.PurchaseOrder = nil
	purchaseOrder.Inquiry = inquiry

	TrackCustomerIOTask{
		UserID: workerInstance.App.Config.InflowMerchandiseGroupEmail,
		Event:  customerio.EventPoCreated,
		Data:   purchaseOrder.GetCustomerIOMetadata(nil),
	}.Dispatch(ctx)

	return err

}

// Dispatch dispatch event
func (task BulkPurchaseOrderBankTransferRejectedTask) Dispatch(ctx context.Context, opts ...asynq.Option) (*asynq.TaskInfo, error) {
	return workerInstance.SendTaskWithContext(ctx, task, opts...)
}
