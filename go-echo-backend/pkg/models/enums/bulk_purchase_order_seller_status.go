package enums

type BulkPurchaseOrderSellerStatus string

var (
	BulkPurchaseOrderSellerStatusWaitingForQuotation BulkPurchaseOrderSellerStatus = "waiting_for_quotation"
	BulkPurchaseOrderSellerStatusWaitingForApproval  BulkPurchaseOrderSellerStatus = "waiting_for_approval"
	BulkPurchaseOrderSellerStatusApproved            BulkPurchaseOrderSellerStatus = "approved"
	BulkPurchaseOrderSellerStatusRejected            BulkPurchaseOrderSellerStatus = "rejected"
)

func (p BulkPurchaseOrderSellerStatus) String() string {
	return string(p)
}

func (p BulkPurchaseOrderSellerStatus) DisplayName() string {
	var name = string(p)

	switch p {
	case BulkPurchaseOrderSellerStatusWaitingForQuotation:
		name = "Waiting for quotation"
	case BulkPurchaseOrderSellerStatusWaitingForApproval:
		name = "Waiting for approval"
	case BulkPurchaseOrderSellerStatusApproved:
		name = "Approved"
	case BulkPurchaseOrderSellerStatusRejected:
		name = "Rejected"
	}

	return name
}
