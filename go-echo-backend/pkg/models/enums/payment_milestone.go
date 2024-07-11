package enums

type PaymentMilestone string

var (
	PaymentMilestoneDeposit       PaymentMilestone = "deposit"
	PaymentMilestoneFirstPayment  PaymentMilestone = "first_payment"
	PaymentMilestoneSecondPayment PaymentMilestone = "second_payment"
	PaymentMilestoneFinalPayment  PaymentMilestone = "final_payment"
)

func (p PaymentMilestone) DisplayName() string {
	switch p {
	case PaymentMilestoneDeposit:
		return "Deposit"
	case PaymentMilestoneFirstPayment:
		return "1st Payment"
	case PaymentMilestoneSecondPayment:
		return "2nd Payment"
	case PaymentMilestoneFinalPayment:
		return "Final Payment"
	}

	return string(p)
}
