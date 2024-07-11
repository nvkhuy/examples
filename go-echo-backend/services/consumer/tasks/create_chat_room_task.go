package tasks

import (
	"context"
	"encoding/json"

	"github.com/engineeringinflow/inflow-backend/pkg/models"
	"github.com/engineeringinflow/inflow-backend/pkg/models/enums"
	"github.com/engineeringinflow/inflow-backend/pkg/repo"
	"github.com/hibiken/asynq"
)

type CreateChatRoomTask struct {
	UserID              string     `json:"user_id"`
	Role                enums.Role `json:"role"`
	InquiryID           string     `json:"inquiry_id"`
	PurchaseOrderID     string     `json:"purchase_order_id"`
	BulkPurchaseOrderID string     `json:"bulk_purchase_order_id"`
	BuyerID             string     `json:"buyer_id"`
	SellerID            string     `json:"seller_id"`
}

func (task CreateChatRoomTask) GetPayload() []byte {
	data, _ := json.Marshal(&task)

	return data
}

// TaskName task name
func (task CreateChatRoomTask) TaskName() string {
	return "create_chat_room"
}

// Handler handler
func (task CreateChatRoomTask) Handler(ctx context.Context, t *asynq.Task) error {
	var err = workerInstance.BindAndValidate(t.Payload(), &task)
	if err != nil {
		return err
	}
	var claims = models.JwtClaimsInfo{}
	claims.SetUserID(task.UserID)
	claims.SetRole(task.Role)

	if _, err := repo.NewChatRoomRepo(workerInstance.App.DB).CreateChatRoom(&models.CreateChatRoomRequest{
		JwtClaimsInfo:       claims,
		InquiryID:           task.InquiryID,
		PurchaseOrderID:     task.PurchaseOrderID,
		BulkPurchaseOrderID: task.BulkPurchaseOrderID,
		BuyerID:             task.BuyerID,
		SellerID:            task.SellerID,
	}); err != nil {
		return err
	}

	return nil

}

// Dispatch dispatch event
func (task CreateChatRoomTask) Dispatch(ctx context.Context, opts ...asynq.Option) (*asynq.TaskInfo, error) {
	return workerInstance.SendTaskWithContext(ctx, task, opts...)
}
