package main

import (
	"fmt"

	"github.com/engineeringinflow/inflow-backend/pkg/app"
	"github.com/engineeringinflow/inflow-backend/pkg/config"
	"github.com/engineeringinflow/inflow-backend/pkg/db"
	"github.com/engineeringinflow/inflow-backend/pkg/db/callback"
	"github.com/engineeringinflow/inflow-backend/pkg/logger"
	"github.com/engineeringinflow/inflow-backend/pkg/models"
	"github.com/engineeringinflow/inflow-backend/pkg/models/enums"
)

func FixPoTrackingUserID() {
	var cfg = config.New("../deployment/config/dev/secrets.env")
	logger.Init()

	var app = app.New(cfg).WithDB(db.New(cfg, callback.New(), nil))

	var orders []*models.PurchaseOrder
	app.DB.Select("*").Find(&orders)

	adminActions := []enums.PoTrackingAction{
		enums.PoTrackingActionApproveDesign,
		enums.PoTrackingActionRejectDesign,
		enums.PoTrackingActionUpdateDesign,
		enums.PoTrackingActionUpdateMaterial,
		enums.PoTrackingActionApproveRawMaterial,
		enums.PoTrackingActionMarkMaking,
		enums.PoTrackingActionMarkSubmit,
		enums.PoTrackingActionMarkDelivering,
		enums.PoTrackingActionConfirmDelivered,
	}

	for _, order := range orders {
		var logs []*models.PurchaseOrderTracking
		app.DB.Select("*").Where("purchase_order_id = ? AND action_type IN ?", order.ID, adminActions).Find(&logs)

		for _, log := range logs {
			var trackingUpdates = models.PurchaseOrderTracking{
				UserID:          order.UserID,
				CreatedByUserID: log.UserID,
			}
			var err = app.DB.Model(&models.PurchaseOrderTracking{}).Where("id = ?", log.ID).Updates(trackingUpdates)

			fmt.Printf("Update id=%s err=%+v", order.ID, err)
		}
	}

}
