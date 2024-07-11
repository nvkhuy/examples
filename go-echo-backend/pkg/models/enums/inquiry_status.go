package enums

type InquiryStatus string

var (
	InquiryStatusNew            InquiryStatus = "new"
	InquiryStatusQuoteInProcess InquiryStatus = "quote_in_process"
	InquiryStatusProduction     InquiryStatus = "production"
	InquiryStatusFinished       InquiryStatus = "finished"
	InquiryStatusClosed         InquiryStatus = "closed"
	InquiryStatusCanceled       InquiryStatus = "canceled"
)

func (p InquiryStatus) String() string {
	return string(p)
}

func (p InquiryStatus) DisplayName() string {
	var name = string(p)

	switch p {
	case InquiryStatusNew:
		name = "New"
	case InquiryStatusQuoteInProcess:
		name = "In Process"
	case InquiryStatusProduction:
		name = "Production"
	case InquiryStatusFinished:
		name = "Sample"
	case InquiryStatusClosed:
		name = "Closed"
	case InquiryStatusCanceled:
		name = "Canceled"
	}

	return name
}
