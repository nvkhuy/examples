package enums

type PoTrackingAction string

var (
	PoTrackingActionApproveDesign      PoTrackingAction = "approve_design"
	PoTrackingActionRejectDesign       PoTrackingAction = "reject_design"
	PoTrackingActionUpdateDesign       PoTrackingAction = "update_design"
	PoTrackingActionUpdateMaterial     PoTrackingAction = "update_raw_material"
	PoTrackingActionApproveRawMaterial PoTrackingAction = "approve_raw_material"
	PoTrackingActionMarkMaking         PoTrackingAction = "mark_making"
	PoTrackingActionMarkSubmit         PoTrackingAction = "mark_submit"
	PoTrackingActionMarkDelivering     PoTrackingAction = "mark_delivering"
	PoTrackingActionConfirmDelivered   PoTrackingAction = "confirm_delivered"
	PoTrackingActionStageComment       PoTrackingAction = "stage_comment"
	PoTrackingActionPaymentReceived    PoTrackingAction = "payment_received"
	// For seller module
	PoTrackingActionSellerMarkMaking       PoTrackingAction = "seller_mark_making"
	PoTrackingActionSellerMarkSubmit       PoTrackingAction = "seller_mark_submit"
	PoTrackingActionSellerMarkDelivering   PoTrackingAction = "seller_mark_delivering"
	PoTrackingActionSellerConfirmDelivered PoTrackingAction = "seller_confirm_delivered"
	PoTrackingActionSellerApprovedDesign   PoTrackingAction = "seller_approved_design"
	PoTrackingActionSellerApprovedPO       PoTrackingAction = "seller_approved_po"
	PoTrackingActionSellerRejectedPO       PoTrackingAction = "seller_rejected_po"
	PoTrackingActionSellerPaymentReceived  PoTrackingAction = "seller_payment_received"
	PoTrackingActionSellerSkipRawMaterial  PoTrackingAction = "seller_skip_raw_material"

	PoTrackingActionAdminCancel         PoTrackingAction = "admin_cancel"
	PoTrackingActionAdminConfirm        PoTrackingAction = "admin_confirm"
	PoTrackingActionAdminApprovedDesign PoTrackingAction = "admin_approved_design"
	PoTrackingActionAdminUpdatedDesign  PoTrackingAction = "admin_updated_design"
	PoTrackingActionAdminUploadedPO     PoTrackingAction = "admin_uploaded_po"
)
