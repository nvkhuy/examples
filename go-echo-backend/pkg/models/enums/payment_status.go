package enums

type PaymentStatus string

var (
	PaymentStatusRefunded       PaymentStatus = "refunded"
	PaymentStatusPaid           PaymentStatus = "paid"
	PaymentStatusUnpaid         PaymentStatus = "unpaid"
	PaymentStatusWaitingConfirm PaymentStatus = "waiting_confirm" // Bank transfer

)

func (p PaymentStatus) DisplayName() string {
	var name = string(p)
	switch p {
	case PaymentStatusPaid:
		return "Paid"

	case PaymentStatusUnpaid:
		return "Unpaid"

	case PaymentStatusWaitingConfirm:
		return "Waiting confirm"

	case PaymentStatusRefunded:
		return "Refunded"
	}
	return name
}
