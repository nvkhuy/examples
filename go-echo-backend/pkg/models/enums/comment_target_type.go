package enums

type CommentTargetType string

var (
	CommentTargetTypeOrderCard            CommentTargetType = "order_card"
	CommentTargetTypeInquirySku           CommentTargetType = "inquiry_sku"
	CommentTargetTypeInquiryInternalNotes CommentTargetType = "inquiry_internal_notes"
	CommentTargetTypePOInternalNotes      CommentTargetType = "purchase_order_internal_notes"
	CommentTargetTypeBulkPOInternalNotes  CommentTargetType = "bulk_purchase_order_internal_notes"

	CommentTargetTypePurchaseOrderDesign  CommentTargetType = "purchase_order_design"
	CommentTargetTypeInquirySellerDesign  CommentTargetType = "inquiry_seller_design"
	CommentTargetTypeInquirySellerRequest CommentTargetType = "inquiry_seller_request"
	CommentTargetTypeSellerPoUpload       CommentTargetType = "seller_po_upload"
	CommentTargetTypeSellerPoDesign       CommentTargetType = "seller_po_design"
	CommentTargetTypeSellerPoFinalDesign  CommentTargetType = "seller_po_final_design"
	CommentTargetTypeSellerPoRawMaterial  CommentTargetType = "seller_po_raw_material"

	// CommentTargetTypeInquirySellerFabric   CommentTargetType = "inquiry_seller_fabric"
	// CommentTargetTypeInquirySellerDesign   CommentTargetType = "inquiry_seller_design"
	// CommentTargetTypeInquirySellerTechpack CommentTargetType = "inquiry_seller_techpack"
)
