package enums

type RegisterQuantity string

var (
	Quantity_Range_50 RegisterQuantity = "0-50"
	Quantity_Range_100 RegisterQuantity = "50-100"
	Quantity_Range_500 RegisterQuantity = "100-500"
	Quantity_Range_1000 RegisterQuantity = "500-1000"
	Quantity_Range_10000 RegisterQuantity = "1000-10000"
)

func (register_quantity RegisterQuantity) String() string {
	return string(register_quantity)
}

func (register_quantity RegisterQuantity) DisplayName() string {
	switch register_quantity {
	case Quantity_Range_50:
		return "Under 50"

	case Quantity_Range_100:
		return "50-100"

	case Quantity_Range_500:
		return "100-500"

	case Quantity_Range_1000:
		return "500-1000"

	case Quantity_Range_10000:
		return "1000-10000"

	}

	return string(register_quantity)
}
