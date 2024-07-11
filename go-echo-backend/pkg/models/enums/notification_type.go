package enums

type CmsNotificationType string

var (
	CmsNotificationTypeNewInquiry                 CmsNotificationType = "new_inquiry"
	CmsNotificationTypeNewInquiryQuotation        CmsNotificationType = "new_inquiry_quotation"
	CmsNotificationTypePoDesignNewComment         CmsNotificationType = "po_design_new_comment"
	CmsNotificationTypeInquiryNewNote             CmsNotificationType = "inquiry_new_note"
	CmsNotificationTypeInquirySellerDesignComment CmsNotificationType = "inquiry_seller_design_comment"
)
