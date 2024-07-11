package enums

type PurchaseOrderStatus string

var (
	PurchaseOrderStatusPaid           PurchaseOrderStatus = "paid"
	PurchaseOrderStatusUnpaid         PurchaseOrderStatus = "unpaid"
	PurchaseOrderStatusFailed         PurchaseOrderStatus = "failed"
	PurchaseOrderStatusPending        PurchaseOrderStatus = "pending"
	PurchaseOrderStatusWaitingConfirm PurchaseOrderStatus = "waiting_confirm" // Bank transfer
	PurchaseOrderStatusCanceled       PurchaseOrderStatus = "canceled"
)
