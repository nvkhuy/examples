package tasks

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/engineeringinflow/inflow-backend/pkg/customerio"
	"github.com/engineeringinflow/inflow-backend/pkg/models"
	"github.com/engineeringinflow/inflow-backend/pkg/repo"
	"github.com/hibiken/asynq"
	"github.com/jinzhu/copier"
	"github.com/rotisserie/eris"
	"github.com/samber/lo"
	"math"
	"time"
)

type PurchaseOrderRawMaterialApproveTask struct {
	models.JwtClaimsInfo

	PurchaseOrderID string `json:"purchase_order_id" param:"purchase_order_id" query:"purchase_order_id" validate:"required"`

	ItemIDs []string `json:"item_ids" param:"item_ids" query:"item_ids"`
}

func (task PurchaseOrderRawMaterialApproveTask) GetPayload() []byte {
	data, _ := json.Marshal(&task)
	return data
}

// TaskName task name
func (task PurchaseOrderRawMaterialApproveTask) TaskName() string {
	return "purchase_order_raw_material_approve"
}

// Handler handler
func (task PurchaseOrderRawMaterialApproveTask) Handler(ctx context.Context, t *asynq.Task) (err error) {
	err = workerInstance.BindAndValidate(t.Payload(), &task)
	if err != nil {
		return err
	}
	now := time.Now().Unix()
	var po models.PurchaseOrder
	if err = workerInstance.App.DB.Model(&models.PurchaseOrder{}).
		Select("approve_raw_material_at").Where("id = ?", task.PurchaseOrderID).First(&po).Error; err != nil {
		return
	}
	if po.ApproveRawMaterialAt != nil {
		if math.Abs(float64(*po.ApproveRawMaterialAt-now)) > 5 {
			err = errors.New("raw material approval deadline extended")
			return
		}
	}

	var params repo.PurchaseApproveRawMaterialParams
	err = copier.Copy(&params, &task)
	if err != nil {
		return err
	}
	order, err := repo.NewPurchaseOrderRepo(workerInstance.App.DB).PurchaseOrderApproveRawMaterial(params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	for _, assigneeID := range order.AssigneeIDs {
		var data = order.GetCustomerIOMetadata(nil)
		if order.PoRawMaterials != nil && len(*order.PoRawMaterials) > 0 {
			var list models.PoRawMaterialMetas = lo.Filter(*order.PoRawMaterials, func(item *models.PoRawMaterialMeta, index int) bool {
				return lo.Contains(params.ItemIDs, item.ReferenceID)
			})

			data["approved_raw_materials"] = list.GenerateFileURL()
		}
		_, _ = TrackCustomerIOTask{
			UserID: assigneeID,
			Event:  customerio.EventPoApproveRawMaterial,
			Data:   data,
		}.Dispatch(ctx)
	}

	return err
}

// Dispatch dispatch event
func (task PurchaseOrderRawMaterialApproveTask) Dispatch(ctx context.Context, opts ...asynq.Option) (*asynq.TaskInfo, error) {
	return workerInstance.SendTaskWithContext(ctx, task, opts...)
}

func (task PurchaseOrderRawMaterialApproveTask) DispatchAt(at time.Time, opts ...asynq.Option) (*asynq.TaskInfo, error) {
	return workerInstance.SendTaskAt(task, at, opts...)
}
