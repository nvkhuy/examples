package enums

type OBMOQType string

var (
	OBMOQTypeLT100 OBMOQType = "lt_100"
	OBMOQTypeLT300 OBMOQType = "lt_300"
	OBMOQTypeGT500 OBMOQType = "gt_500"
)

func (p OBMOQType) String() string {
	return string(p)
}

func (p OBMOQType) DisplayName() string {
	var name = string(p)

	switch p {
	case OBMOQTypeLT100:
		name = "50 - 100 pcs / color"
	case OBMOQTypeLT300:
		name = "200 - 300 pcs / color"
	case OBMOQTypeGT500:
		name = "Over 500 pcs / color"
	}

	return name
}
