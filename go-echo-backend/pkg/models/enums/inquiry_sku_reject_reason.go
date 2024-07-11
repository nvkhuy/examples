package enums

type InquirySkuRejectReason string

var (
	InquirySkuRejectReasonUnreasonablePrice InquirySkuRejectReason = "unreasonable_price"
	InquirySkuRejectReasonChangeMaterial    InquirySkuRejectReason = "change_material"
	InquirySkuRejectReasonChangeQuanity     InquirySkuRejectReason = "change_quantity"
	InquirySkuRejectReasonChangeSize        InquirySkuRejectReason = "change_size"
	InquirySkuRejectReasonOther             InquirySkuRejectReason = "other"
)
