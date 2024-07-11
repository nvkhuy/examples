package tasks

import (
	"context"
	"fmt"

	"encoding/json"

	"github.com/engineeringinflow/inflow-backend/pkg/helper"
	"github.com/engineeringinflow/inflow-backend/pkg/hubspot"
	"github.com/engineeringinflow/inflow-backend/pkg/models"
	"github.com/engineeringinflow/inflow-backend/pkg/models/enums"
	"github.com/engineeringinflow/inflow-backend/pkg/repo"
	"github.com/hibiken/asynq"
	"github.com/rotisserie/eris"
	"gorm.io/gorm"
)

type HubspotSyncPOTask struct {
	PurchaseOrderID string `json:"purchase_order_id" validate:"required"`
	UserID          string `json:"user_id" validate:"required"`
	IsAdmin         bool   `json:"is_admin"`
}

func (task HubspotSyncPOTask) GetPayload() []byte {
	data, _ := json.Marshal(&task)

	return data
}

// TaskName task name
func (task HubspotSyncPOTask) TaskName() string {
	return "hubspot_sync_po"
}

// Handler handler
func (task HubspotSyncPOTask) Handler(ctx context.Context, t *asynq.Task) error {
	var err = workerInstance.BindAndValidate(t.Payload(), &task)
	if err != nil {
		return err
	}
	if !workerInstance.App.Config.IsProd() {
		return nil
	}

	purchaseOrder, err := repo.NewPurchaseOrderRepo(workerInstance.App.DB).GetPurchaseOrder(repo.GetPurchaseOrderParams{
		PurchaseOrderID: task.PurchaseOrderID,
	})
	if err != nil {
		return err
	}

	if purchaseOrder.HubspotDealID == "" {
		var hubspotOwnerID string
		workerInstance.App.DB.Model(&models.User{}).Select("HubspotOwnerID").
			First(&hubspotOwnerID, "id IN ? AND COALESCE(hubspot_owner_id,'') <> ?", purchaseOrder.AssigneeIDs, "")

		var hubspotContactID string
		workerInstance.App.DB.Model(&models.User{}).Select("HubspotOwnerID").
			First(&hubspotOwnerID, "id = ? AND COALESCE(hubspot_owner_id,'') <> ?", purchaseOrder.UserID, "")

		var form = &hubspot.DealPropertiesForm{
			Amount:           "",
			Dealname:         purchaseOrder.ReferenceID,
			Pipeline:         hubspot.PipelineManufacturing,
			Dealstage:        hubspot.DealStagePending,
			HubspotOwnerID:   hubspotOwnerID,
			HubspotContactID: hubspotContactID,
		}

		if purchaseOrder.Inquiry != nil {
			form.Dealname = fmt.Sprintf("%s - %s", purchaseOrder.Inquiry.ReferenceID, purchaseOrder.Inquiry.Title)
		}

		deal, err := workerInstance.App.HubspotClient.CreateDeal(form)
		if err != nil {
			return eris.Wrap(err, err.Error())
		}
		err = workerInstance.App.DB.Transaction(func(tx *gorm.DB) error {
			err = tx.Model(&models.PurchaseOrder{}).Where("id = ?", purchaseOrder.ID).UpdateColumn("HubspotDealID", deal.ID).Error

			if err != nil {
				return eris.Wrap(err, err.Error())
			}

			if purchaseOrder.Inquiry != nil && purchaseOrder.Inquiry.ID != "" {
				err = tx.Model(&models.Inquiry{}).Where("id = ?", purchaseOrder.Inquiry.ID).UpdateColumn("HubspotDealID", deal.ID).Error
			}
			return err
		})
		if err != nil {
			return eris.Wrap(err, err.Error())
		}

		t.ResultWriter().Write(helper.ToJson(deal))
		workerInstance.Logger.Debugf("Create deal %s success", deal.ID)

	} else {
		if purchaseOrder.Status == enums.PurchaseOrderStatusPaid {
			deal, err := workerInstance.App.HubspotClient.UpdateDeal(purchaseOrder.HubspotDealID, &hubspot.DealPropertiesForm{
				Dealstage: hubspot.DealStageSample,
			})
			if err != nil {
				return eris.Wrap(err, err.Error())
			}

			if purchaseOrder.Inquiry != nil && purchaseOrder.Inquiry.ID != "" {
				err = workerInstance.App.DB.Model(&models.Inquiry{}).Where("id = ?", purchaseOrder.Inquiry.ID).UpdateColumn("HubspotDealID", deal.ID).Error
				if err != nil {
					return eris.Wrap(err, err.Error())
				}
			}

			t.ResultWriter().Write(helper.ToJson(deal))
			workerInstance.Logger.Debugf("Update deal %s success", deal.ID)
		}
	}

	if err != nil {
		return eris.Wrap(err, err.Error())
	}
	return err

}

// Dispatch dispatch event
func (task HubspotSyncPOTask) Dispatch(ctx context.Context, opts ...asynq.Option) (*asynq.TaskInfo, error) {
	return workerInstance.SendTaskWithContext(ctx, task, opts...)
}
