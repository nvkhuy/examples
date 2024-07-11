package enums

type FactoryTourStatus string

var (
	FactoryTourStatusActive   FactoryTourStatus = "active"
	FactoryTourStatusInactive FactoryTourStatus = "inactive"
)

func (p FactoryTourStatus) String() string {
	return string(p)
}

func (p FactoryTourStatus) DisplayName() string {
	var name = string(p)

	switch p {
	case FactoryTourStatusActive:
		name = "Active"
	case FactoryTourStatusInactive:
		name = "Inactive"
	}

	return name
}
