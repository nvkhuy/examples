package enums

type BulkPurchaseOrderStatus string

var (
	BulkPurchaseOrderStatusNew                BulkPurchaseOrderStatus = "new"
	BulkPurchaseOrderStatusQuoteInProcess     BulkPurchaseOrderStatus = "quote_in_process"
	BulkPurchaseOrderStatusWaitingForApproved BulkPurchaseOrderStatus = "waiting_for_approved"
	BulkPurchaseOrderStatusWaitingForPo       BulkPurchaseOrderStatus = "waiting_for_po"
	BulkPurchaseOrderStatusWaitingForPayment  BulkPurchaseOrderStatus = "waiting_for_payment"
	BulkPurchaseOrderStatusPartiallyPaid      BulkPurchaseOrderStatus = "partially_paid" // Bank transfer
	BulkPurchaseOrderStatusFullyPaid          BulkPurchaseOrderStatus = "fully_paid"     // Bank transfer

)
