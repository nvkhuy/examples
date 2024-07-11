package enums

type SellerQuotationFilter string

var (
	SellerQuotationFilterNew                SellerQuotationFilter = "new"
	SellerQuotationFilterSent               SellerQuotationFilter = "sent"
	SellerQuotationFilterWaitingForApproval SellerQuotationFilter = "waiting_for_approval"
)

func (p SellerQuotationFilter) String() string {
	return string(p)
}

func (p SellerQuotationFilter) DisplayName() string {
	var name = string(p)

	switch p {
	case SellerQuotationFilterNew:
		name = "New"
	case SellerQuotationFilterSent:
		name = "Sent"
	case SellerQuotationFilterWaitingForApproval:
		name = "Waiting for approval"
	}

	return name
}
