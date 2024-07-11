package enums

type InvoiceType string

var (
	InvoiceTypeInquiry              InvoiceType = "inquiry"
	InvoiceTypeMultipleInquiry      InvoiceType = "multiple_inquiry"
	InvoiceTypeBulkPODepositPayment InvoiceType = "bulk_po_deposit_payment"
	InvoiceTypeBulkPOFirstPayment   InvoiceType = "bulk_po_first_payment"
	InvoiceTypeBulkPOSecondPayment  InvoiceType = "bulk_po_first_payment"
	InvoiceTypeBulkPOFinalPayment   InvoiceType = "bulk_po_final_payment"
)

func (it InvoiceType) DisplayName() string {
	var name = string(it)

	switch it {
	case InvoiceTypeInquiry:
		name = "Sample"

	case InvoiceTypeBulkPODepositPayment:
		name = "Debit Note (Deposit)"

	case InvoiceTypeBulkPOFirstPayment:
		name = "Debit Note (1st Payment)"

	case InvoiceTypeBulkPOSecondPayment:
		name = "Debit Note (2nd Payment)"

	case InvoiceTypeBulkPOFinalPayment:
		name = "Bulk Order"
	}

	return name
}
