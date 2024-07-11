package enums

type InquirySkuStatus string

var (
	InquirySkuStatusNew                 InquirySkuStatus = "new"
	InquirySkuStatusWaitingForQuotation InquirySkuStatus = "waiting_for_quotation"
	InquirySkuStatusWaitingForApproval  InquirySkuStatus = "waiting_for_approval"
	InquirySkuStatusApproved            InquirySkuStatus = "approved"
	InquirySkuStatusRejected            InquirySkuStatus = "rejected"
)

func (p InquirySkuStatus) String() string {
	var name = string(p)

	return name
}

func (p InquirySkuStatus) DisplayName() string {
	var name = string(p)

	switch p {
	case InquirySkuStatusNew:
		name = "New"
	case InquirySkuStatusWaitingForQuotation:
		name = "Waiting for quotation"
	case InquirySkuStatusWaitingForApproval:
		name = "Reviewing quotation"
	case InquirySkuStatusApproved:
		name = "Approved"
	case InquirySkuStatusRejected:
		name = "Rejected"

	}

	return name
}
