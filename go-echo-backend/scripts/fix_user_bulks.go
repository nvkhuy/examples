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

func FixUserBulks() {
	var cfg = config.New("../deployment/config/dev/env.json")
	logger.Init()

	var app = app.New(cfg).WithDB(db.New(cfg, callback.New(), nil))

	var bulks []*models.BulkPurchaseOrder
	app.DB.Select("ID", "InquiryID").Find(&bulks, "coalesce(purchase_order_id,'') = ?", "")

	for index, bulk := range bulks {
		var purchaseOrder models.PurchaseOrder
		var err = app.DB.Select("ID").First(&purchaseOrder, "inquiry_id = ?", bulk.InquiryID).Error
		if err != nil {
			fmt.Printf("Update %d/%d id=%s err=%+v", index, len(bulks), bulk.ID, err)
			continue
		}

		err = app.DB.Model(&models.BulkPurchaseOrder{}).Where("id = ?", bulk.ID).UpdateColumn("purchase_order_id", purchaseOrder.ID).Error
		if err != nil {
			fmt.Printf("Update %d/%d id=%s err=%+v", index, len(bulks), bulk.ID, err)
			continue
		}

		fmt.Printf("Update %d/%d id=%s err=%+v", index, len(bulks), bulk.ID, err)

	}

}
