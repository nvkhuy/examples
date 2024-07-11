package enums

type OBLeadTime string

var (
	OBLeadTime10 OBLeadTime = "10"
	OBLeadTime15 OBLeadTime = "15"
	OBLeadTime20 OBLeadTime = "20"
)

func (p OBLeadTime) String() string {
	return string(p)
}

func (p OBLeadTime) DisplayName() string {
	var name = string(p)

	switch p {
	case OBLeadTime10:
		name = "10 days CMPT"
	case OBLeadTime15:
		name = "15 days CMPT"
	case OBLeadTime20:
		name = "20 days CMPT"
	}

	return name
}
