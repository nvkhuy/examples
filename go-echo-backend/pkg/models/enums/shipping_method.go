package enums

type ShippingMethod string

var (
	ShippingMethodFOB ShippingMethod = "fob"
	ShippingMethodCIF ShippingMethod = "cif"
	ShippingMethodEXW ShippingMethod = "exw"
)

func (l ShippingMethod) String() string {
	return string(l)
}

func (l ShippingMethod) DisplayName() string {
	switch l {
	case ShippingMethodFOB:
		return "Collect at loading port"
	case ShippingMethodCIF:
		return "Collect at destination port"
	case ShippingMethodEXW:
		return "Collect at factory"
	}

	return string(l)
}

func (l ShippingMethod) Description() string {
	switch l {
	case ShippingMethodFOB:
		return "Collect at loading port"
	case ShippingMethodCIF:
		return "Collect at destination port"
	case ShippingMethodEXW:
		return "Collect at factory"
	}

	return string(l)
}
