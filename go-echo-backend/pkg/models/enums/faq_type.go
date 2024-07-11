package enums

type FaqType string

var (
	FaqTypeBuyer  FaqType = "buyer"
	FaqTypeSeller FaqType = "seller"
)

func (p FaqType) String() string {
	return string(p)
}

func (p FaqType) DisplayName() string {
	var name = string(p)

	switch p {
	case FaqTypeBuyer:
		name = "Buyer"
	case FaqTypeSeller:
		name = "Seller"
	}

	return name
}
