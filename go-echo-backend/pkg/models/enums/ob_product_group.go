package enums

type OBProductGroup string

var (
	OBProductGroupClothing OBProductGroup = "clothing"
	OBProductGroupShoes    OBProductGroup = "shoes"
	OBProductGroupBag      OBProductGroup = "bag"
	OBProductGroupCap      OBProductGroup = "cap"
	OBProductGroupTent     OBProductGroup = "tent"
	OBProductGroupCurtain  OBProductGroup = "curtain"
	OBProductGroupBedding  OBProductGroup = "bedding"
	OBProductGroupBlanket  OBProductGroup = "blanket"
	OBProductGroupPillow   OBProductGroup = "pillow"
)

func (p OBProductGroup) String() string {
	return string(p)
}

func (p OBProductGroup) DisplayName() string {
	var name = string(p)

	switch p {
	case OBProductGroupClothing:
		name = "Clothing"
	case OBProductGroupShoes:
		name = "Shoes"
	case OBProductGroupBag:
		name = "Bag"
	case OBProductGroupCap:
		name = "Cap"
	case OBProductGroupTent:
		name = "Tent"
	case OBProductGroupCurtain:
		name = "Curtain"
	case OBProductGroupBedding:
		name = "Bedding"
	case OBProductGroupBlanket:
		name = "Blanket"
	case OBProductGroupPillow:
		name = "Pillow"
	}

	return name
}
