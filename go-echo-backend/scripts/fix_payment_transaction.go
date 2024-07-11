package main

import (
	"fmt"

	"github.com/engineeringinflow/inflow-backend/pkg/app"
	"github.com/engineeringinflow/inflow-backend/pkg/config"
	"github.com/engineeringinflow/inflow-backend/pkg/db"
	"github.com/engineeringinflow/inflow-backend/pkg/db/callback"
	"github.com/engineeringinflow/inflow-backend/pkg/logger"
	"github.com/engineeringinflow/inflow-backend/pkg/models"
	"github.com/engineeringinflow/inflow-backend/pkg/repo/query"
	"github.com/engineeringinflow/inflow-backend/pkg/repo/query/queryfunc"
	"github.com/samber/lo"
)

func FixPaymentTransactions() {
	var cfg = config.New("../deployment/config/prod/env.json")
	logger.Init()

	var app = app.New(cfg).WithDB(db.New(cfg, callback.New(), nil))

	app.DB.AutoMigrate(&models.PaymentTransaction{})

	var paymentTransactions []*models.PaymentTransaction
	app.DB.Find(&paymentTransactions)

	for index, trans := range paymentTransactions {
		if trans.BulkPurchaseOrderID != "" && trans.Metadata == nil {
			var bulkPO models.BulkPurchaseOrder
			var err = app.DB.Select("ID", "ReferenceID", "InquiryID").First(&bulkPO, "id = ?", trans.BulkPurchaseOrderID).Error
			if err != nil {
				fmt.Printf("Update %d/%d id=%s err=%+v\n", index, len(paymentTransactions), trans.ID, err)
				continue
			}

			var inquiry models.Inquiry
			err = app.DB.Select("ID", "ReferenceID").First(&inquiry, "id = ?", bulkPO.InquiryID).Error
			if err != nil {
				fmt.Printf("Update %d/%d id=%s err=%+v\n", index, len(paymentTransactions), trans.ID, err)
				continue
			}

			trans.Metadata = &models.PaymentTransactionMetadata{
				InquiryID:                    inquiry.ID,
				InquiryReferenceID:           inquiry.ReferenceID,
				BulkPurchaseOrderID:          bulkPO.ID,
				BulkPurchaseOrderReferenceID: bulkPO.ReferenceID,
			}
			err = app.DB.Model(&models.PaymentTransaction{}).Where("bulk_purchase_order_id = ?", bulkPO.ID).UpdateColumn("Metadata", trans.Metadata).Error
			if err != nil {
				fmt.Printf("Update %d/%d id=%s err=%+v\n", index, len(paymentTransactions), trans.ID, err)
				continue
			}
			fmt.Printf("Update %d/%d id=%s\n", index, len(paymentTransactions), trans.ID)
			continue
		}

		if trans.PurchaseOrderID != "" && trans.Metadata == nil {
			var purchaseOrder models.PurchaseOrder
			var err = app.DB.Select("ID", "ReferenceID", "InquiryID").First(&purchaseOrder, "id = ?", trans.PurchaseOrderID).Error
			if err != nil {
				fmt.Printf("Update %d/%d id=%s err=%+v\n", index, len(paymentTransactions), trans.ID, err)
				continue
			}

			var inquiry models.Inquiry
			err = app.DB.Select("ID", "ReferenceID").First(&inquiry, "id = ?", purchaseOrder.InquiryID).Error
			if err != nil {
				fmt.Printf("Update %d/%d id=%s err=%+v\n", index, len(paymentTransactions), trans.ID, err)
				continue
			}

			trans.Metadata = &models.PaymentTransactionMetadata{
				InquiryID:                inquiry.ID,
				InquiryReferenceID:       inquiry.ReferenceID,
				PurchaseOrderID:          purchaseOrder.ID,
				PurchaseOrderReferenceID: purchaseOrder.ReferenceID,
			}
			err = app.DB.Model(&models.PaymentTransaction{}).Where("purchase_order_id = ?", purchaseOrder.ID).UpdateColumn("Metadata", trans.Metadata).Error
			if err != nil {
				fmt.Printf("Update %d/%d id=%s err=%+v\n", index, len(paymentTransactions), trans.ID, err)
				continue
			}
			fmt.Printf("Update %d/%d id=%s\n", index, len(paymentTransactions), trans.ID)
			continue
		}

		if trans.CheckoutSessionID != "" && trans.Metadata == nil {
			var purchaseOrders []*models.PurchaseOrder
			query.New(app.DB, queryfunc.NewPurchaseOrderBuilder(queryfunc.PurchaseOrderBuilderOptions{})).
				Where("po.checkout_session_id = ?", trans.CheckoutSessionID).
				FindFunc(&purchaseOrders)

			trans.Metadata = &models.PaymentTransactionMetadata{
				InquiryIDs: lo.Map(purchaseOrders, func(item *models.PurchaseOrder, index int) string {
					return item.InquiryID
				}),
				InquiryReferenceIDs: lo.Map(purchaseOrders, func(item *models.PurchaseOrder, index int) string {
					return item.Inquiry.ReferenceID
				}),
				PurchaseOrderIDs: lo.Map(purchaseOrders, func(item *models.PurchaseOrder, index int) string {
					return item.ID
				}),
				PurchaseOrderReferenceIDs: lo.Map(purchaseOrders, func(item *models.PurchaseOrder, index int) string {
					return item.ReferenceID
				}),
			}
			var err = app.DB.Model(&models.PaymentTransaction{}).Where("checkout_session_id = ?", trans.CheckoutSessionID).UpdateColumn("Metadata", trans.Metadata).Error
			if err != nil {
				fmt.Printf("Update %d/%d id=%s err=%+v\n", index, len(paymentTransactions), trans.ID, err)
				continue
			}
			fmt.Printf("Update %d/%d id=%s\n", index, len(paymentTransactions), trans.ID)
			continue
		}

	}

}
