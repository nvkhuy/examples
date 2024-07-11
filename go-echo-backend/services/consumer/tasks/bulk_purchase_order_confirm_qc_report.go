package tasks

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/engineeringinflow/inflow-backend/pkg/models"
	"github.com/engineeringinflow/inflow-backend/pkg/models/enums"
	"github.com/engineeringinflow/inflow-backend/pkg/repo"
	"github.com/hibiken/asynq"
	"github.com/jinzhu/copier"
	"math"
	"time"
)

type BulkPurchaseQCApproveTask struct {
	models.JwtClaimsInfo

	BulkPurchaseOrderID string                     `json:"bulk_purchase_order_id" param:"bulk_purchase_order_id" query:"bulk_purchase_order_id" validate:"required"`
	TrackingStatus      enums.BulkPoTrackingStatus `json:"tracking_status" param:"tracking_status" query:"tracking_status"`
	TrackingAction      enums.BulkPoTrackingAction `json:"tracking_action" param:"tracking_action" query:"tracking_action"`

	UserID string `json:"-"`
}

func (task BulkPurchaseQCApproveTask) GetPayload() []byte {
	data, _ := json.Marshal(&task)
	return data
}

// TaskName task name
func (task BulkPurchaseQCApproveTask) TaskName() string {
	return "bulk_purchase_order_qc_approve"
}

// Handler handler
func (task BulkPurchaseQCApproveTask) Handler(ctx context.Context, t *asynq.Task) (err error) {
	err = workerInstance.BindAndValidate(t.Payload(), &task)
	if err != nil {
		return err
	}
	now := time.Now().Unix()
	var bpo models.BulkPurchaseOrder
	if err = workerInstance.App.DB.Model(&models.BulkPurchaseOrder{}).
		Select("approve_qc_at").Where("id = ?", task.BulkPurchaseOrderID).First(&bpo).Error; err != nil {
		return
	}
	if bpo.ApproveQCAt != nil {
		if math.Abs(float64(*bpo.ApproveQCAt-now)) > 5 {
			err = errors.New("bulk PO QC approval deadline extended")
			return
		}
	}

	var params repo.BulkPurchaseOrderUpdateTrackingStatusParams
	err = copier.Copy(&params, &task)
	if err != nil {
		return err
	}

	params.TrackingAction = enums.BulkPoTrackingActionBuyerApproveQc
	params.TrackingStatus = enums.BulkPoTrackingStatusSubmit
	_, err = repo.NewBulkPurchaseOrderRepo(workerInstance.App.DB).BulkPurchaseOrderUpdateTrackingStatus(params)
	return err
}

// Dispatch dispatch event
func (task BulkPurchaseQCApproveTask) Dispatch(ctx context.Context, opts ...asynq.Option) (*asynq.TaskInfo, error) {
	return workerInstance.SendTaskWithContext(ctx, task, opts...)
}

func (task BulkPurchaseQCApproveTask) DispatchAt(at time.Time, opts ...asynq.Option) (*asynq.TaskInfo, error) {
	return workerInstance.SendTaskAt(task, at, opts...)
}
