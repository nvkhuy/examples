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

type HubspotSyncInquiryTask struct {
	InquiryID string `json:"inquiry_id" validate:"required"`
	UserID    string `json:"user_id" validate:"required"`
	IsAdmin   bool   `json:"is_admin"`
}

func (task HubspotSyncInquiryTask) GetPayload() []byte {
	data, _ := json.Marshal(&task)

	return data
}

// TaskName task name
func (task HubspotSyncInquiryTask) TaskName() string {
	return "hubspot_sync_inquiry"
}

// Handler handler
func (task HubspotSyncInquiryTask) Handler(ctx context.Context, t *asynq.Task) error {
	var err = workerInstance.BindAndValidate(t.Payload(), &task)
	if err != nil {
		return err
	}
	if !workerInstance.App.Config.IsProd() {
		return nil
	}

	inquiry, err := repo.NewInquiryRepo(workerInstance.App.DB).GetInquiryByID(repo.GetInquiryByIDParams{
		InquiryID: task.InquiryID,
	})
	if err != nil {
		return err
	}

	if inquiry.HubspotDealID == "" && len(inquiry.AssigneeIDs) > 0 && inquiry.UserID != "" {
		var hubspotOwnerID string
		workerInstance.App.DB.Model(&models.User{}).Select("HubspotOwnerID").
			First(&hubspotOwnerID, "id IN ? AND COALESCE(hubspot_owner_id,'') <> ?", inquiry.AssigneeIDs, "")

		var hubspotContactID string
		workerInstance.App.DB.Model(&models.User{}).Select("HubspotContactID").
			First(&hubspotContactID, "id = ? AND COALESCE(hubspot_owner_id,'') <> ?", inquiry.UserID, "")

		deal, err := workerInstance.App.HubspotClient.CreateDeal(&hubspot.DealPropertiesForm{
			Amount:           "",
			Dealname:         fmt.Sprintf("%s - %s", inquiry.ReferenceID, inquiry.Title),
			Pipeline:         hubspot.PipelineManufacturing,
			Dealstage:        hubspot.DealStagePending,
			HubspotOwnerID:   hubspotOwnerID,
			HubspotContactID: hubspotContactID,
		})
		if err != nil {
			return eris.Wrap(err, err.Error())
		}
		err = workerInstance.App.DB.Transaction(func(tx *gorm.DB) error {
			err = tx.Model(&models.Inquiry{}).Where("id = ?", inquiry.ID).UpdateColumn("HubspotDealID", deal.ID).Error
			if err != nil {
				return eris.Wrap(err, err.Error())
			}

			if inquiry.PurchaseOrder != nil && inquiry.PurchaseOrder.ID != "" {
				err = tx.Model(&models.PurchaseOrder{}).Where("id = ?", inquiry.PurchaseOrder.ID).UpdateColumn("HubspotDealID", deal.ID).Error
			}
			return err
		})
		if err != nil {
			return err
		}

		t.ResultWriter().Write(helper.ToJson(deal))
		workerInstance.Logger.Debugf("Create deal %s success", deal.ID)

	} else {
		if inquiry.PurchaseOrder != nil && inquiry.PurchaseOrder.Status == enums.PurchaseOrderStatusPaid {
			deal, err := workerInstance.App.HubspotClient.UpdateDeal(inquiry.HubspotDealID, &hubspot.DealPropertiesForm{
				Dealstage: hubspot.DealStageSample,
			})
			if err != nil {
				return eris.Wrap(err, err.Error())
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
func (task HubspotSyncInquiryTask) Dispatch(ctx context.Context, opts ...asynq.Option) (*asynq.TaskInfo, error) {
	return workerInstance.SendTaskWithContext(ctx, task, opts...)
}
