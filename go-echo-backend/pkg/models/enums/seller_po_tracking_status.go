package enums

type SellerPoTrackingStatus string

var (
	SellerPoTrackingStatusNew                    SellerPoTrackingStatus = "new"
	SellerPoTrackingStatusWaitingForPayment      SellerPoTrackingStatus = "waiting_for_payment"
	SellerPoTrackingStatusPaymentConfirmed       SellerPoTrackingStatus = "payment_confirmed"
	SellerPoTrackingStatusRejectPO               SellerPoTrackingStatus = "reject_po"
	SellerPoTrackingStatusDesignApproval         SellerPoTrackingStatus = "design_approval"
	SellerPoTrackingStatusDesignApprovedBySeller SellerPoTrackingStatus = "design_approved_by_seller"
	SellerPoTrackingStatusDesignApprovedByBuyer  SellerPoTrackingStatus = "design_approved_by_buyer"
	SellerPoTrackingStatusDesignApprovedByAdmin  SellerPoTrackingStatus = "design_approved_by_admin"
	SellerPoTrackingStatusRawMaterial            SellerPoTrackingStatus = "raw_material"
	SellerPoTrackingStatusRawMaterialSkipped     SellerPoTrackingStatus = "raw_material_skipped"
	SellerPoTrackingStatusMaking                 SellerPoTrackingStatus = "making"
	SellerPoTrackingStatusSubmit                 SellerPoTrackingStatus = "submit"
	SellerPoTrackingStatusDelivering             SellerPoTrackingStatus = "delivering"
	SellerPoTrackingStatusDeliveryConfirmed      SellerPoTrackingStatus = "delivery_confirmed"
)

func (p SellerPoTrackingStatus) String() string {
	return string(p)
}

func (p SellerPoTrackingStatus) DisplayName() string {
	var name = string(p)

	switch p {
	case SellerPoTrackingStatusNew:
		name = "New"
	case SellerPoTrackingStatusWaitingForPayment:
		name = "Waiting for payment"
	case SellerPoTrackingStatusPaymentConfirmed:
		name = "Payment confirmed"
	}

	return name
}
