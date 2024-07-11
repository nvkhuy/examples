package enums

type OBDecorationService string

var (
	OBDecorationServiceWashing OBDecorationService = "washing"
	OBDecorationServiceDrying  OBDecorationService = "drying"
)

func (p OBDecorationService) String() string {
	return string(p)
}

func (p OBDecorationService) DisplayName() string {
	var name = string(p)

	switch p {
	case OBDecorationServiceWashing:
		name = "Washing"
	case OBDecorationServiceDrying:
		name = "Drying"
	}

	return name
}
