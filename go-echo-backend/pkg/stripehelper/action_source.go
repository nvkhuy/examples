package stripehelper

type ActionSource string

var (
	ActionSourceCreatePaymentLink ActionSource = "create_payment_link"

	ActionSourceInquiryPayment       ActionSource = "inquiry_payment"
	ActionSourceBulkPODepositPayment ActionSource = "bulk_po_deposit_payment"
	ActionSourceBulkPOFirstPayment   ActionSource = "bulk_po_first_payment"
	ActionSourceBulkPOSecondPayment  ActionSource = "bulk_po_second_payment"
	ActionSourceBulkPOFinalPayment   ActionSource = "bulk_po_final_payment"
	ActionSourceMultiInquiryPayment  ActionSource = "multi_inquiry_payment"
	ActionSourceMultiPOPayment       ActionSource = "multi_po_payment"
	ActionSourceOrderCartPayment     ActionSource = "order_cart_payment"
)

func (ac ActionSource) IsValid() bool {
	switch ac {
	case ActionSourceCreatePaymentLink,
		ActionSourceInquiryPayment,
		ActionSourceBulkPODepositPayment,
		ActionSourceBulkPOFirstPayment,
		ActionSourceBulkPOSecondPayment,
		ActionSourceBulkPOFinalPayment,
		ActionSourceMultiInquiryPayment,
		ActionSourceMultiPOPayment,
		ActionSourceOrderCartPayment:
		return true
	}

	return false
}
