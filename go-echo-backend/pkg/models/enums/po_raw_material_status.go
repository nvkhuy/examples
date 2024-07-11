package enums

type PoRawMaterialStatus string

var (
	PoRawMaterialStatusDying             PoRawMaterialStatus = "dying"
	PoRawMaterialStatusProcessing        PoRawMaterialStatus = "processing"
	PoRawMaterialStatusWaitingForApprove PoRawMaterialStatus = "waiting_for_approval"
	PoRawMaterialStatusApproved          PoRawMaterialStatus = "approved"
)

func (p PoRawMaterialStatus) String() string {
	return string(p)
}

func (p PoRawMaterialStatus) DisplayName() string {
	var name = string(p)

	switch p {
	case PoRawMaterialStatusDying:
		name = "Dying"
	case PoRawMaterialStatusProcessing:
		name = "Processing"
	case PoRawMaterialStatusWaitingForApprove:
		name = "Waiting for approval"
	case PoRawMaterialStatusApproved:
		name = "Approved"
	}

	return name
}
