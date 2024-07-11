package enums

type BulkPoTrackingStatus string

var (
	BulkPoTrackingStatusNew                    BulkPoTrackingStatus = "new"
	BulkPoTrackingStatusWaitingForSubmitOrder  BulkPoTrackingStatus = "waiting_for_submit_order"
	BulkPoTrackingStatusWaitingForQuotation    BulkPoTrackingStatus = "waiting_for_quotation"
	BulkPoTrackingStatusFirstPayment           BulkPoTrackingStatus = "first_payment"
	BulkPoTrackingStatusFirstPaymentConfirm    BulkPoTrackingStatus = "first_payment_confirm"
	BulkPoTrackingStatusFirstPaymentConfirmed  BulkPoTrackingStatus = "first_payment_confirmed"
	BulkPoTrackingStatusSecondPayment          BulkPoTrackingStatus = "second_payment"
	BulkPoTrackingStatusSecondPaymentConfirm   BulkPoTrackingStatus = "second_payment_confirm"
	BulkPoTrackingStatusSecondPaymentConfirmed BulkPoTrackingStatus = "second_payment_confirmed"
	BulkPoTrackingStatusRawMaterial            BulkPoTrackingStatus = "raw_material"
	BulkPoTrackingStatusPps                    BulkPoTrackingStatus = "pps" // pre production step
	BulkPoTrackingStatusProduction             BulkPoTrackingStatus = "production"
	BulkPoTrackingStatusQc                     BulkPoTrackingStatus = "qc"
	BulkPoTrackingStatusSubmit                 BulkPoTrackingStatus = "submit"
	BulkPoTrackingStatusFinalPayment           BulkPoTrackingStatus = "final_payment"
	BulkPoTrackingStatusFinalPaymentConfirm    BulkPoTrackingStatus = "final_payment_confirm"
	BulkPoTrackingStatusFinalPaymentConfirmed  BulkPoTrackingStatus = "final_payment_confirmed"
	BulkPoTrackingStatusDelivering             BulkPoTrackingStatus = "delivering"
	BulkPoTrackingStatusDeliveryConfirmed      BulkPoTrackingStatus = "delivery_confirmed"
	BulkPoTrackingStatusDelivered              BulkPoTrackingStatus = "delivered"
)

func (p BulkPoTrackingStatus) String() string {
	return string(p)
}

func (p BulkPoTrackingStatus) DisplayName() string {
	var name = string(p)

	switch p {
	case BulkPoTrackingStatusNew:
		name = "New"
	case BulkPoTrackingStatusWaitingForSubmitOrder:
		name = "Waiting for submit order"
	case BulkPoTrackingStatusWaitingForQuotation:
		name = "Waiting for quotation"
	case BulkPoTrackingStatusFirstPayment:
		name = "First payment"
	case BulkPoTrackingStatusFirstPaymentConfirm:
		name = "First payment confirm"
	case BulkPoTrackingStatusRawMaterial:
		name = "Raw material"
	case BulkPoTrackingStatusProduction:
		name = "Production"
	case BulkPoTrackingStatusQc:
		name = "QC"
	case BulkPoTrackingStatusSubmit:
		name = "Submit"
	case BulkPoTrackingStatusFinalPayment:
		name = "Final payment"
	case BulkPoTrackingStatusFinalPaymentConfirm:
		name = "Final payment confirm"
	case BulkPoTrackingStatusDelivering:
		name = "Delivering"
	case BulkPoTrackingStatusDeliveryConfirmed:
		name = "Delivery Confirmed"
	case BulkPoTrackingStatusFirstPaymentConfirmed:
		name = "First payment confirmed"
	case BulkPoTrackingStatusSecondPayment:
		name = "Second payment"
	case BulkPoTrackingStatusSecondPaymentConfirm:
		name = "Second payment confirm"
	case BulkPoTrackingStatusSecondPaymentConfirmed:
		name = "Seconde payment confirmed"
	case BulkPoTrackingStatusPps:
		name = "Pps"
	case BulkPoTrackingStatusFinalPaymentConfirmed:
		name = "Final payment confirmed"
	case BulkPoTrackingStatusDelivered:
		name = "Delivered"
	}

	return name
}
