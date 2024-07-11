package enums

type PaymentType string

var (
	PaymentTypeCard         PaymentType = "card"
	PaymentTypeBankTransfer PaymentType = "bank_transfer"
)

func (v PaymentType) DisplayName() string {
	switch v {
	case PaymentTypeCard:
		return "Online Payment"

	case PaymentTypeBankTransfer:
		return "Bank Trasnfer"
	}
	return string(v)
}
