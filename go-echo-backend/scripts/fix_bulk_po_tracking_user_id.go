package main

import (
	"fmt"

	"github.com/engineeringinflow/inflow-backend/pkg/app"
	"github.com/engineeringinflow/inflow-backend/pkg/config"
	"github.com/engineeringinflow/inflow-backend/pkg/db"
	"github.com/engineeringinflow/inflow-backend/pkg/db/callback"
	"github.com/engineeringinflow/inflow-backend/pkg/logger"
	"github.com/engineeringinflow/inflow-backend/pkg/models"
)

func FixBulkPoTrackingUserID() {
	var cfg = config.New("../deployment/config/local/secrets.env")
	logger.Init()

	var app = app.New(cfg).WithDB(db.New(cfg, callback.New(), nil))

	var orders []*models.BulkPurchaseOrder
	app.DB.Select("*").Find(&orders)

	for _, order := range orders {
		var logs []*models.BulkPurchaseOrderTracking
		app.DB.Select("*").Where("purchase_order_id = ?", order.ID).Find(&logs)

		for _, log := range logs {
			var trackingUpdates = models.BulkPurchaseOrderTracking{
				UserID:          order.UserID,
				CreatedByUserID: log.UserID,
			}
			var err = app.DB.Model(&models.BulkPurchaseOrderTracking{}).Where("id = ?", log.ID).Updates(trackingUpdates)

			fmt.Printf("Update id=%s err=%+v", order.ID, err)
		}
	}

}
