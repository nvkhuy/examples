package enums

type OBPaymentTerm string

var (
	OBPaymentTermPartialPayment OBPaymentTerm = "partial_payment"
	OBPaymentTermFullPayment    OBPaymentTerm = "full_payment"
	OBPaymentTermOpenForInvoice OBPaymentTerm = "open_for_invoice"
)

func (p OBPaymentTerm) String() string {
	return string(p)
}

func (p OBPaymentTerm) DisplayName() string {
	var name = string(p)

	switch p {
	case OBPaymentTermPartialPayment:
		name = "30% deposit - 70% within 15 days after receiving the goods"
	case OBPaymentTermFullPayment:
		name = "100% payment within 15 days upon receipt"
	case OBPaymentTermOpenForInvoice:
		name = "Open for Invoice financing"
	}

	return name
}
