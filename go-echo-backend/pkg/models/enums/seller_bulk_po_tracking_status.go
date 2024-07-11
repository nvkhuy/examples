package enums

type SellerBulkPoTrackingStatus string

var (
	// SellerBulkPoTrackingStatusNew                   SellerBulkPoTrackingStatus = "new"
	SellerBulkPoTrackingStatusPO                    SellerBulkPoTrackingStatus = "po"
	SellerBulkPoTrackingStatusPORejected            SellerBulkPoTrackingStatus = "po_rejected"
	SellerBulkPoTrackingStatusWaitingFirstPayment   SellerBulkPoTrackingStatus = "waiting_first_payment"
	SellerBulkPoTrackingStatusWaitingForSubmitOrder SellerBulkPoTrackingStatus = "waiting_for_submit_order"
	SellerBulkPoTrackingStatusWaitingForQuotation   SellerBulkPoTrackingStatus = "waiting_for_quotation"
	SellerBulkPoTrackingStatusFirstPayment          SellerBulkPoTrackingStatus = "first_payment"
	SellerBulkPoTrackingStatusFirstPaymentConfirm   SellerBulkPoTrackingStatus = "first_payment_confirm"

	SellerBulkPoTrackingStatusFirstPaymentConfirmed SellerBulkPoTrackingStatus = "first_payment_confirmed"
	SellerBulkPoTrackingStatusFirstPaymentSkipped   SellerBulkPoTrackingStatus = "first_payment_skipped"
	SellerBulkPoTrackingStatusRawMaterial           SellerBulkPoTrackingStatus = "raw_material"
	SellerBulkPoTrackingStatusPps                   SellerBulkPoTrackingStatus = "pps" // pre production step
	SellerBulkPoTrackingStatusProduction            SellerBulkPoTrackingStatus = "production"
	SellerBulkPoTrackingStatusQc                    SellerBulkPoTrackingStatus = "qc"
	SellerBulkPoTrackingStatusSubmit                SellerBulkPoTrackingStatus = "submit"
	SellerBulkPoTrackingStatusFinalPayment          SellerBulkPoTrackingStatus = "final_payment"
	SellerBulkPoTrackingStatusFinalPaymentConfirm   SellerBulkPoTrackingStatus = "final_payment_confirm"
	SellerBulkPoTrackingStatusFinalPaymentConfirmed SellerBulkPoTrackingStatus = "final_payment_confirmed"
	SellerBulkPoTrackingStatusDelivering            SellerBulkPoTrackingStatus = "delivering"
	SellerBulkPoTrackingStatusDeliveryConfirmed     SellerBulkPoTrackingStatus = "delivery_confirmed"
	SellerBulkPoTrackingStatusDelivered             SellerBulkPoTrackingStatus = "delivered"

	SellerBulkPoTrackingStatusInspection SellerBulkPoTrackingStatus = "inspection"
)

func (p SellerBulkPoTrackingStatus) String() string {
	return string(p)
}

func (p SellerBulkPoTrackingStatus) DisplayName() string {
	var name = string(p)

	switch p {
	// case SellerBulkPoTrackingStatusNew:
	// 	name = "New"
	case SellerBulkPoTrackingStatusPO:
		name = "PO"
	case SellerBulkPoTrackingStatusPORejected:
		name = "PO Rejected"
	case SellerBulkPoTrackingStatusWaitingFirstPayment:
		name = "Waiting first payment"
	case SellerBulkPoTrackingStatusWaitingForSubmitOrder:
		name = "Waiting for submit order"
	case SellerBulkPoTrackingStatusWaitingForQuotation:
		name = "Waiting for quotation"
	case SellerBulkPoTrackingStatusFirstPayment:
		name = "First payment"
	case SellerBulkPoTrackingStatusFirstPaymentConfirm:
		name = "First payment confirm"
	case SellerBulkPoTrackingStatusFirstPaymentSkipped:
		name = "First payment skipped"
	case SellerBulkPoTrackingStatusRawMaterial:
		name = "Raw material"
	case SellerBulkPoTrackingStatusProduction:
		name = "Production"
	case SellerBulkPoTrackingStatusQc:
		name = "QC"
	case SellerBulkPoTrackingStatusSubmit:
		name = "Submit"
	case SellerBulkPoTrackingStatusFinalPayment:
		name = "Final payment"
	case SellerBulkPoTrackingStatusFinalPaymentConfirm:
		name = "Final payment confirm"
	case SellerBulkPoTrackingStatusDelivering:
		name = "Delivering"
	case SellerBulkPoTrackingStatusDeliveryConfirmed:
		name = "Delivery Confirmed"
	case SellerBulkPoTrackingStatusFirstPaymentConfirmed:
		name = "First payment confirmed"
	case SellerBulkPoTrackingStatusPps:
		name = "Pps"
	case SellerBulkPoTrackingStatusFinalPaymentConfirmed:
		name = "Final payment confirmed"
	case SellerBulkPoTrackingStatusDelivered:
		name = "Delivered"
	case SellerBulkPoTrackingStatusInspection:
		name = "Inspection"
	}

	return name
}
