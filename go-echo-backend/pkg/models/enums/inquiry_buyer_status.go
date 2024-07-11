package enums

type InquiryBuyerStatus string

var (
	InquiryBuyerStatusNew                InquiryBuyerStatus = "new"
	InquiryBuyerStatusWaitingForApproved InquiryBuyerStatus = "waiting_for_approved"
	InquiryBuyerStatusApproved           InquiryBuyerStatus = "approved"
	InquiryBuyerStatusRejected           InquiryBuyerStatus = "rejected"
)

func (p InquiryBuyerStatus) String() string {
	return string(p)
}

func (p InquiryBuyerStatus) DisplayName() string {
	var name = string(p)

	switch p {
	case InquiryBuyerStatusNew:
		name = "New"
	case InquiryBuyerStatusWaitingForApproved:
		name = "In Process"
	case InquiryBuyerStatusApproved:
		name = "Approved"
	case InquiryBuyerStatusRejected:
		name = "Rejected"
	}

	return name
}
