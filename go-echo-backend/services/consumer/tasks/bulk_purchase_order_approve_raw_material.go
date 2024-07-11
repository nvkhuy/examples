package tasks

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/engineeringinflow/inflow-backend/pkg/models"
	"github.com/engineeringinflow/inflow-backend/pkg/repo"
	"github.com/hibiken/asynq"
	"github.com/jinzhu/copier"
	"math"
	"time"
)

type BulkPurchaseOrderRawMaterialApproveTask struct {
	models.JwtClaimsInfo

	BulkPurchaseOrderID string `json:"purchase_order_id" param:"purchase_order_id" query:"purchase_order_id" validate:"required"`

	ItemIDs []string `json:"item_ids" param:"item_ids" query:"item_ids"`
}

func (task BulkPurchaseOrderRawMaterialApproveTask) GetPayload() []byte {
	data, _ := json.Marshal(&task)
	return data
}

// TaskName task name
func (task BulkPurchaseOrderRawMaterialApproveTask) TaskName() string {
	return "bulk_purchase_order_raw_material_approve"
}

// Handler handler
func (task BulkPurchaseOrderRawMaterialApproveTask) Handler(ctx context.Context, t *asynq.Task) (err error) {
	err = workerInstance.BindAndValidate(t.Payload(), &task)
	if err != nil {
		return err
	}
	now := time.Now().Unix()
	var bpo models.BulkPurchaseOrder
	if err = workerInstance.App.DB.Model(&models.BulkPurchaseOrder{}).
		Select("approve_raw_material_at").Where("id = ?", task.BulkPurchaseOrderID).First(&bpo).Error; err != nil {
		return
	}
	if bpo.ApproveRawMaterialAt != nil {
		if math.Abs(float64(*bpo.ApproveRawMaterialAt-now)) > 5 {
			err = errors.New("bulk PO raw material approval deadline extended")
			return
		}
	}

	var params repo.BulkPurchaseBuyerApproveRawMaterialParams
	err = copier.Copy(&params, &task)
	if err != nil {
		return err
	}
	_, err = repo.NewBulkPurchaseOrderRepo(workerInstance.App.DB).BulkPurchaseOrderBuyerApproveRawMaterial(params)
	return err
}

// Dispatch dispatch event
func (task BulkPurchaseOrderRawMaterialApproveTask) Dispatch(ctx context.Context, opts ...asynq.Option) (*asynq.TaskInfo, error) {
	return workerInstance.SendTaskWithContext(ctx, task, opts...)
}

func (task BulkPurchaseOrderRawMaterialApproveTask) DispatchAt(at time.Time, opts ...asynq.Option) (*asynq.TaskInfo, error) {
	return workerInstance.SendTaskAt(task, at, opts...)
}
