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

func MigrateOrderCartItem() {
	var cfg = config.New("../deployment/config/local/env.json")
	logger.Init()

	var app = app.New(cfg).WithDB(db.New(cfg, callback.New(), nil))

	var inquiryCartItems []*models.InquiryCartItem

	if err := app.DB.Find(&inquiryCartItems).Error; err != nil {
		fmt.Printf("Get inquiry_cart_items err=%+v\n", err)
		return
	}
	var inquiryIDs []string
	for _, item := range inquiryCartItems {
		inquiryIDs = append(inquiryIDs, item.InquiryID)
	}

	var purchaseOrders models.PurchaseOrders
	if err := app.DB.Select("ID", "InquiryID").Find(&purchaseOrders, "inquiry_id IN ?", inquiryIDs).Error; err != nil {
		fmt.Printf("Get purchase_orders err=%+v\n", err)
		return
	}
	var mapInquiryIDToPurchaseOrderID = make(map[string]string)
	for _, po := range purchaseOrders {
		if po.InquiryID != "" {
			mapInquiryIDToPurchaseOrderID[po.InquiryID] = po.ID
		}
	}

	var bulkOrderItems []*models.BulkPurchaseOrderItem
	if err := app.DB.Find(&bulkOrderItems).Error; err != nil {
		fmt.Printf("Get bulk_order_items err=%+v\n", err)
		return
	}
	var orderCartItemsToCreate = make([]*models.OrderCartItem, 0, len((inquiryCartItems)))
	for _, item := range inquiryCartItems {
		poID, ok := mapInquiryIDToPurchaseOrderID[item.InquiryID]
		if !ok {
			fmt.Printf("purchase_order_id not found,inquiry_id:%s\n", item.InquiryID)
			continue
		}
		var orderCartItem = &models.OrderCartItem{
			Model: models.Model{
				ID:        item.ID,
				CreatedAt: item.CreatedAt,
				UpdatedAt: item.UpdatedAt,
				DeletedAt: &item.DeletedAt,
			},
			PurchaseOrderID:    poID,
			Size:               item.Size,
			Sku:                item.Sku,
			ColorName:          item.ColorName,
			Qty:                item.Qty,
			NoteToSupplier:     item.NoteToSupplier,
			UnitPrice:          item.UnitPrice,
			TotalPrice:         item.TotalPrice,
			CheckoutSessionID:  item.CheckoutSessionID,
			WaitingForCheckout: item.WaitingForCheckout,
			Style:              item.Style,
		}
		orderCartItemsToCreate = append(orderCartItemsToCreate, orderCartItem)
	}

	for _, item := range bulkOrderItems {
		var orderCartItem = &models.OrderCartItem{
			Model: models.Model{
				ID:        item.ID,
				CreatedAt: item.CreatedAt,
				UpdatedAt: item.UpdatedAt,
				DeletedAt: &item.DeletedAt,
			},
			BulkPurchaseOrderID: item.PurchaseOrderID,
			Size:                item.Size,
			Sku:                 item.Sku,
			ColorName:           item.ColorName,
			Qty:                 item.Qty,
			UnitPrice:           *item.UnitPrice,
			TotalPrice:          *item.TotalPrice,
			Style:               item.Style,
		}
		orderCartItemsToCreate = append(orderCartItemsToCreate, orderCartItem)
	}

	if err := app.DB.Clauses(clause.OnConflict{Columns: []clause.Column{{Name: "id"}}, DoNothing: true}).
		Create(&orderCartItemsToCreate).Error; err != nil {
		fmt.Printf("Migrate order_cart_items err=%+v\n", err)
		return
	}

}
