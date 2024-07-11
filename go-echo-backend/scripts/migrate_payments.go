package main

import (
	"fmt"

	"github.com/engineeringinflow/inflow-backend/pkg/app"
	"github.com/engineeringinflow/inflow-backend/pkg/config"
	"github.com/engineeringinflow/inflow-backend/pkg/db"
	"github.com/engineeringinflow/inflow-backend/pkg/db/callback"
	"github.com/engineeringinflow/inflow-backend/pkg/logger"
	"github.com/engineeringinflow/inflow-backend/pkg/models"
	"gorm.io/gorm/clause"
)

func MigratePayments() {
	var cfg = config.New("../deployment/config/local/env.json")
	logger.Init()

	var app = app.New(cfg).WithDB(db.New(cfg, callback.New(), nil))

	var payments models.PaymentTransactions
	if err := app.DB.Find(&payments).Error; err != nil {
		fmt.Printf("Get payment err=%+v\n", err)
		return
	}

	for _, payment := range payments {
		if len(payment.PurchaseOrderIDs) == 0 && payment.PurchaseOrderID != "" {
			payment.PurchaseOrderIDs = append(payment.PurchaseOrderIDs, payment.PurchaseOrderID)
			if payment.Metadata != nil {
				payment.Metadata.PurchaseOrderIDs = payment.PurchaseOrderIDs
			} else {
				payment.Metadata = &models.PaymentTransactionMetadata{
					PurchaseOrderIDs: payment.PurchaseOrderIDs,
				}
			}
		}
		if len(payment.BulkPurchaseOrderIDs) == 0 && payment.BulkPurchaseOrderID != "" {
			payment.BulkPurchaseOrderIDs = append(payment.BulkPurchaseOrderIDs, payment.BulkPurchaseOrderID)
			if payment.Metadata != nil {
				payment.Metadata.BulkPurchaseOrderIDs = payment.BulkPurchaseOrderIDs
			} else {
				payment.Metadata = &models.PaymentTransactionMetadata{
					BulkPurchaseOrderIDs: payment.BulkPurchaseOrderIDs,
				}
			}
		}
	}
	if err := app.DB.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		DoUpdates: clause.AssignmentColumns([]string{"purchase_order_ids", "bulk_purchase_order_ids", "metadata"})},
	).
		Create(&payments).Error; err != nil {
		fmt.Printf("Update payment err=%+v\n", err)
		return
	}
}
