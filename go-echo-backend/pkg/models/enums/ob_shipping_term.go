package enums

type OBShippingTerm string

var (
	OBShippingTermExWork OBShippingTerm = "ex_work"
	OBShippingTermFOB    OBShippingTerm = "fob"
	OBShippingTermCIF    OBShippingTerm = "cif"
)

func (p OBShippingTerm) String() string {
	return string(p)
}

func (p OBShippingTerm) DisplayName() string {
	var name = string(p)

	switch p {
	case OBShippingTermExWork:
		name = "Ex-works"
	case OBShippingTermFOB:
		name = "FOB"
	case OBShippingTermCIF:
		name = "CIF"
	}

	return name
}

func (p OBShippingTerm) Description() string {
	switch p {
	case OBShippingTermExWork:
		return "Goods Handover at factory's gate"
	case OBShippingTermFOB:
		return "Goods Handover at export place"
	case OBShippingTermCIF:
		return "Goods Handover at import place"
	}

	return string(p)
}
