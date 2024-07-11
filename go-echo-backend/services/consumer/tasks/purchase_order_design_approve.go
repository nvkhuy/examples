package tasks

import (
	"context"
	"encoding/json"
	"errors"
	"math"
	"time"

	"github.com/engineeringinflow/inflow-backend/pkg/customerio"
	"github.com/engineeringinflow/inflow-backend/pkg/models"
	"github.com/engineeringinflow/inflow-backend/pkg/repo"
	"github.com/hibiken/asynq"
	"github.com/jinzhu/copier"
	"github.com/rotisserie/eris"
)

type PurchaseOrderDesignApproveTask struct {
	models.JwtClaimsInfo
	PurchaseOrderID   string                           `json:"purchase_order_id" param:"purchase_order_id" query:"purchase_order_id" validate:"required"`
	ApproveRejectMeta *models.InquiryApproveRejectMeta `json:"approve_reject_meta" param:"approve_reject_meta" query:"approve_reject_meta"`
}

func (task PurchaseOrderDesignApproveTask) GetPayload() []byte {
	data, _ := json.Marshal(&task)
	return data
}

// TaskName task name
func (task PurchaseOrderDesignApproveTask) TaskName() string {
	return "purchase_order_design_approve"
}

// Handler handler
func (task PurchaseOrderDesignApproveTask) Handler(ctx context.Context, t *asynq.Task) (err error) {
	err = workerInstance.BindAndValidate(t.Payload(), &task)
	if err != nil {
		return err
	}
	now := time.Now().Unix()
	var po models.PurchaseOrder
	if err = workerInstance.App.DB.Model(&models.PurchaseOrder{}).
		Select("approve_design_at").Where("id = ?", task.PurchaseOrderID).First(&po).Error; err != nil {
		return
	}
	if po.ApproveDesignAt == nil {
		err = errors.New("design rejected")
		return
	}
	if po.ApproveDesignAt != nil {
		if math.Abs(float64(*po.ApproveDesignAt-now)) > 5 {
			err = errors.New("design approval deadline extended")
			return
		}
	}

	var params repo.BuyerApproveDesignParams
	err = copier.Copy(&params, &task)
	if err != nil {
		return err
	}
	purchaseOrder, err := repo.NewPurchaseOrderRepo(workerInstance.App.DB).BuyerApproveDesign(params)
	if err != nil {
		return eris.Wrap(err, err.Error())
	}

	TrackCustomerIOTask{
		UserID: purchaseOrder.UserID,
		Event:  customerio.EventPoApproveDesign,
		Data:   purchaseOrder.GetCustomerIOMetadata(nil),
	}.Dispatch(ctx)

	return err
}

// Dispatch dispatch event
func (task PurchaseOrderDesignApproveTask) Dispatch(ctx context.Context, opts ...asynq.Option) (*asynq.TaskInfo, error) {
	return workerInstance.SendTaskWithContext(ctx, task, opts...)
}

func (task PurchaseOrderDesignApproveTask) DispatchAt(at time.Time, opts ...asynq.Option) (*asynq.TaskInfo, error) {
	return workerInstance.SendTaskAt(task, at, opts...)
}
