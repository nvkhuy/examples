package enums

type OBOrderType string

var (
	OBOrderTypeFOB OBOrderType = "fob"
	OBOrderTypeCMT OBOrderType = "cmt"
	OBOrderTypeCM  OBOrderType = "cm"
)

func (p OBOrderType) String() string {
	return string(p)
}

func (p OBOrderType) DisplayName() string {
	var name = string(p)

	switch p {
	case OBOrderTypeFOB:
		name = "Full package/FOB"
	case OBOrderTypeCMT:
		name = "CMT"
	case OBOrderTypeCM:
		name = "CM"
	}

	return name
}

func (p OBOrderType) Description() string {
	switch p {
	case OBOrderTypeFOB:
		return "Fabric + CMT"
	case OBOrderTypeCMT:
		return "Cut + Make + Finishing + Trim/Accessory"
	case OBOrderTypeCM:
		return "Sewing thread + Cut + Sew + Finishing"
	}

	return string(p)
}
