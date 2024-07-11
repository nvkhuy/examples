package enums

type InquirySellerStatus string

var (
	InquirySellerStatusNew                 InquirySellerStatus = "new"
	InquirySellerStatusOfferRejected       InquirySellerStatus = "offer_rejected"
	InquirySellerStatusWaitingForQuotation InquirySellerStatus = "waiting_for_quotation"
	InquirySellerStatusWaitingForApproval  InquirySellerStatus = "waiting_for_approval"
	InquirySellerStatusApproved            InquirySellerStatus = "approved"
	InquirySellerStatusRejected            InquirySellerStatus = "rejected"
)

func (p InquirySellerStatus) String() string {
	return string(p)
}

func (p InquirySellerStatus) DisplayName() string {
	var name = string(p)

	switch p {
	case InquirySellerStatusNew:
		name = "New"
	case InquirySellerStatusOfferRejected:
		name = "Offer rejected"
	case InquirySellerStatusWaitingForQuotation:
		name = "Waiting for quotation"
	case InquirySellerStatusWaitingForApproval:
		name = "Waiting for approval"
	case InquirySellerStatusApproved:
		name = "Approved"
	case InquirySellerStatusRejected:
		name = "Rejected"
	}

	return name
}
