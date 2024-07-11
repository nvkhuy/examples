package main

import (
	"encoding/json"
	"fmt"

	"github.com/engineeringinflow/inflow-backend/pkg/app"
	"github.com/engineeringinflow/inflow-backend/pkg/config"
	"github.com/engineeringinflow/inflow-backend/pkg/db"
	"github.com/engineeringinflow/inflow-backend/pkg/db/callback"
	"github.com/engineeringinflow/inflow-backend/pkg/helper"
	"github.com/engineeringinflow/inflow-backend/pkg/logger"
	"github.com/engineeringinflow/inflow-backend/pkg/models"
	"github.com/engineeringinflow/inflow-backend/pkg/models/enums"
	"gorm.io/datatypes"
)

type BulkPurchaseOrderAlias struct {
	ID      string
	PpsInfo datatypes.JSON `gorm:"column:pps_info"`
}
type BulkPurchaseOrderTrackingAlias struct {
	ID              string
	PurchaseOrderID string
	Metadata        datatypes.JSON
}

func MigrateBulkPps() {
	var cfg = config.New("../deployment/config/local/env.json")
	logger.Init()

	var app = app.New(cfg).WithDB(db.New(cfg, callback.New(), nil))

	// Get all bpo need to migrate
	var bpoNeedMigrate []BulkPurchaseOrderAlias
	result := app.DB.Model(&models.BulkPurchaseOrder{}).
		Select("ID", "PpsInfo").
		Where("pps_info is not null and jsonb_typeof(pps_info) ='object'").
		Find(&bpoNeedMigrate)
	if result.Error != nil {
		fmt.Printf("Get bpo err=%+v\n", result.Error)
		return
	}
	fmt.Printf("BPO need to migrate length: %d", result.RowsAffected)

	var bpoIDs []string
	var bpoIdToPpsInfo = make(map[string]*models.PoPpsMetas)
	// Transform pps_info of bpo from object to array
	var migrateBpos []models.BulkPurchaseOrder
	for _, bpo := range bpoNeedMigrate {
		bpoIDs = append(bpoIDs, bpo.ID)
		var ppsInfo models.PoPpsMeta
		if err := json.Unmarshal(bpo.PpsInfo, &ppsInfo); err != nil {
			fmt.Printf("Unmarshal pps info err=%+v\n", err)
			return
		}
		ppsInfo.ID = helper.GenerateXID()
		ppsInfo.Status = enums.PpsStatusApproved

		ppsArr := []*models.PoPpsMeta{&ppsInfo}
		migrateBpos = append(migrateBpos, models.BulkPurchaseOrder{
			Model: models.Model{
				ID: bpo.ID,
			},
			PpsInfo: (*models.PoPpsMetas)(&ppsArr),
		})
		bpoIdToPpsInfo[bpo.ID] = (*models.PoPpsMetas)(&ppsArr)
	}
	// Get bpo tracking need to migrate
	var migrateBpoTrackings []models.BulkPurchaseOrderTracking
	result = app.DB.Model(&models.BulkPurchaseOrderTracking{}).
		Select("ID", "PurchaseOrderID", "Metadata").
		Where("purchase_order_id IN ? AND action_type = ? AND metadata -> 'after' -> 'pps_info' is not null", bpoIDs, "update_pps").
		Find(&migrateBpoTrackings)
	if result.Error != nil {
		fmt.Printf("Get bpo tracking err=%+v\n", result.Error)
		return
	}
	fmt.Printf("BPO need to migrate length: %d", result.RowsAffected)

	// Transform pps_info["after"] of bpo tracking from object to array
	for idx, bpoTracking := range migrateBpoTrackings {
		updatedPps, ok := bpoIdToPpsInfo[bpoTracking.PurchaseOrderID]
		if !ok {
			fmt.Printf("The pps_info is empty bpo_id=%s", bpoTracking.PurchaseOrderID)
			return
		}
		migrateBpoTrackings[idx].Metadata.After = map[string]interface{}{
			"pps_info": updatedPps,
		}
	}
	// Migrate transformed bpo pps info
	for _, bpo := range migrateBpos {
		if err := app.DB.Model(&models.BulkPurchaseOrder{}).Where("id = ?", bpo.ID).UpdateColumn("pps_info", bpo.PpsInfo).Error; err != nil {
			fmt.Printf("Update bpo pps info err=%+v\n", err)
			return
		}
	}
	// Migrate transformed bpo tracking metadata
	for _, bpoTracking := range migrateBpoTrackings {
		if err := app.DB.Model(&models.BulkPurchaseOrderTracking{}).Where("id = ?", bpoTracking.ID).UpdateColumn("metadata", bpoTracking.Metadata).Error; err != nil {
			fmt.Printf("Update bpo tracking metadata err=%+v\n", err)
			return
		}
	}
}
