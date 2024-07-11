package models

import "github.com/engineeringinflow/inflow-backend/pkg/models/enums"

type BuyerDataAnalyticRFQ struct {
	Total             int64 `json:"total,omitempty"`
	Approved          int64 `json:"approved,omitempty"`
	Rejected          int64 `json:"rejected,omitempty"`
	PurchaseOrder     int64 `json:"purchase_order,omitempty"`
	BulkPurchaseOrder int64 `json:"bulk_purchase_order,omitempty"`
}

type BuyerDataAnalyticPendingTask struct {
	ID          string       `json:"id,omitempty"`
	ProductName string       `json:"product_name,omitempty"`
	Quantity    *int64       `json:"quantity,omitempty"`
	Attachments *Attachments `json:"attachments,omitempty"`
	OrderID     string       `json:"order_id,omitempty"`
	Status      string       `json:"status,omitempty"`
}

type BuyerDataAnalyticPendingPayment struct {
	ID          string                 `json:"id,omitempty"`
	Quantity    int64                  `json:"quantity,omitempty"`
	Attachments *Attachments           `json:"attachments,omitempty"`
	ProductName string                 `json:"product_name,omitempty"`
	OrderID     string                 `json:"order_id,omitempty"`
	Amount      float64                `json:"amount,omitempty"`
	Currency    enums.Currency         `json:"currency,omitempty"`
	Milestone   enums.PaymentMilestone `json:"milestone"`
}

type BuyerDataAnalyticTotalStyleProduce struct {
	TotalPurchaseOrder      int64                    `json:"total_purchase_order,omitempty"`
	TotalBulkPurchaserOrder int64                    `json:"total_bulk_purchaser_order,omitempty"`
	TotalShipped            int64                    `json:"total_shipped,omitempty"`
	Charts                  []TotalStyleProduceChart `json:"charts"`
}

type TotalStyleProduceChart struct {
	Date              string `json:"date"`
	Timestamp         int64  `json:"timestamp"`
	PurchaseOrder     int64  `json:"purchase_order,omitempty"`
	BulkPurchaseOrder int64  `json:"bulk_purchase_order,omitempty"`
	Shipped           int64  `json:"shipped,omitempty"`
}
