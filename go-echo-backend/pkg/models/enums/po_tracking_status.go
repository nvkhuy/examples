package enums

type PoTrackingStatus string

var (
	PoTrackingStatusNew                PoTrackingStatus = "new"
	PoTrackingStatusWaitingForApproved PoTrackingStatus = "waiting_for_approved"
	PoTrackingStatusDesignApproved     PoTrackingStatus = "design_approved"
	PoTrackingStatusDesignRejected     PoTrackingStatus = "design_rejected"
	PoTrackingStatusRawMaterial        PoTrackingStatus = "raw_material"
	PoTrackingStatusMaking             PoTrackingStatus = "making"
	PoTrackingStatusSubmit             PoTrackingStatus = "submit"
	PoTrackingStatusDelivering         PoTrackingStatus = "delivering"
	PoTrackingStatusDeliveryConfirmed  PoTrackingStatus = "delivery_confirmed"
	PoTrackingStatusCanceled           PoTrackingStatus = "canceled"

	PoTrackingStatusPaymentReceived PoTrackingStatus = "payment_received"
)

func (p PoTrackingStatus) String() string {
	return string(p)
}

func (p PoTrackingStatus) DisplayName() string {
	var name = string(p)

	switch p {
	case PoTrackingStatusNew:
		name = "New"
	case PoTrackingStatusWaitingForApproved:
		name = "Waiting For Approval"
	case PoTrackingStatusDesignApproved:
		name = "Design Approved"
	case PoTrackingStatusDesignRejected:
		name = "Design Rejected"
	case PoTrackingStatusRawMaterial:
		name = "Raw Material"
	case PoTrackingStatusMaking:
		name = "Making"
	case PoTrackingStatusSubmit:
		name = "Submit"
	case PoTrackingStatusDelivering:
		name = "Delivering"
	case PoTrackingStatusDeliveryConfirmed:
		name = "Delivery Confirmed"
	case PoTrackingStatusCanceled:
		name = "Canceled"
	case PoTrackingStatusPaymentReceived:
		name = "Payment Received"
	}

	return name
}

type PoCatalogTrackingStatus string

var (
	PoCatalogTrackingStatusNew       PoCatalogTrackingStatus = "new"
	PoCatalogTrackingStatusPayment   PoCatalogTrackingStatus = "payment"
	PoCatalogTrackingStatusDispatch  PoCatalogTrackingStatus = "dispatch"
	PoCatalogTrackingStatusDelivery  PoCatalogTrackingStatus = "delivery"
	PoCatalogTrackingStatusCompleted PoCatalogTrackingStatus = "completed"
)

func (p PoCatalogTrackingStatus) String() string {
	return string(p)
}

func (p PoCatalogTrackingStatus) DisplayName() string {
	var name = string(p)

	switch p {
	case PoCatalogTrackingStatusNew:
		name = "New"
	case PoCatalogTrackingStatusPayment:
		name = "Payment"
	case PoCatalogTrackingStatusDispatch:
		name = "Dispatch"
	case PoCatalogTrackingStatusDelivery:
		name = "Delivery"
	case PoCatalogTrackingStatusCompleted:
		name = "Approval"
	}

	return name
}
