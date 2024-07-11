package enums

type OBServiceType string

var (
	OBServiceTypeLogisticAndShipping OBServiceType = "logistic_and_shipping"
	OBServiceTypeTesting             OBServiceType = "testing"
	OBServiceTypeDecoration          OBServiceType = "decoration"
)

func (p OBServiceType) String() string {
	return string(p)
}

func (p OBServiceType) DisplayName() string {
	var name = string(p)

	switch p {
	case OBServiceTypeLogisticAndShipping:
		name = "Logistic and Shipping"
	case OBServiceTypeTesting:
		name = "Testing Services"
	case OBServiceTypeDecoration:
		name = "Decoration Services"
	}

	return name
}
