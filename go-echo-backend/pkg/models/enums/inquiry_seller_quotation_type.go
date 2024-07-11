package enums

type InquirySellerQuotationType string

var (
	InquirySellerQuotationTypeMCQ InquirySellerQuotationType = "mcq"
	InquirySellerQuotationTypeMOQ InquirySellerQuotationType = "moq"
)

func (p InquirySellerQuotationType) String() string {
	return string(p)
}

func (p InquirySellerQuotationType) DisplayName() string {
	var name = string(p)

	switch p {
	case InquirySellerQuotationTypeMCQ:
		name = "Minimum color quantity"
	case InquirySellerQuotationTypeMOQ:
		name = "Minimum order quantity"
	}

	return name
}
