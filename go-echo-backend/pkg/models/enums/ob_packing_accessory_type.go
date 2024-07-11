package enums

type OBPackingAccessoryType string

var (
	OBPackingAccessoryTypeHanger             OBPackingAccessoryType = "hanger"
	OBPackingAccessoryTypeSafetyPin          OBPackingAccessoryType = "safety_pin"
	OBPackingAccessoryTypeScotchTape         OBPackingAccessoryType = "scotch_tape"
	OBPackingAccessoryTypePolybag            OBPackingAccessoryType = "polybag"
	OBPackingAccessoryTypeCarton             OBPackingAccessoryType = "carton"
	OBPackingAccessoryTypeTags               OBPackingAccessoryType = "tags"
	OBPackingAccessoryTypeTissuePaper        OBPackingAccessoryType = "tissue_paper"
	OBPackingAccessoryTypeButterPaper        OBPackingAccessoryType = "butter_paper"
	OBPackingAccessoryTypePlasticClip        OBPackingAccessoryType = "plastic_clip"
	OBPackingAccessoryTypePaperBoard         OBPackingAccessoryType = "paper_board"
	OBPackingAccessoryTypeButterfly          OBPackingAccessoryType = "butterfly"
	OBPackingAccessoryTypeShirtCollarSupport OBPackingAccessoryType = "shirt_collar_support"
	OBPackingAccessoryTypeShirtBackSupport   OBPackingAccessoryType = "shirt_back_support"
	OBPackingAccessoryTypeTagPin             OBPackingAccessoryType = "tag_pin"
	OBPackingAccessoryTypePriceTag           OBPackingAccessoryType = "price_tag"
	OBPackingAccessoryTypeBallHeadPin        OBPackingAccessoryType = "ball_head_pin"
	OBPackingAccessoryTypeInnerBox           OBPackingAccessoryType = "inner_box"
	OBPackingAccessoryTypeFoam               OBPackingAccessoryType = "foam"
	OBPackingAccessoryTypeTagGun             OBPackingAccessoryType = "tag_gun"
	OBPackingAccessoryTypeClip               OBPackingAccessoryType = "clip"
	OBPackingAccessoryTypePlasticAdjuster    OBPackingAccessoryType = "plastic_adjuster"
	OBPackingAccessoryTypeShirtBox           OBPackingAccessoryType = "shirt_box"
)

func (p OBPackingAccessoryType) String() string {
	return string(p)
}

func (p OBPackingAccessoryType) DisplayName() string {
	var name = string(p)

	switch p {
	case OBPackingAccessoryTypeHanger:
		name = "Hanger"
	case OBPackingAccessoryTypeSafetyPin:
		name = "Safety Pin"
	case OBPackingAccessoryTypeScotchTape:
		name = "Scotch Tape"
	case OBPackingAccessoryTypePolybag:
		name = "Polybag"
	case OBPackingAccessoryTypeCarton:
		name = "Carton"
	case OBPackingAccessoryTypeTags:
		name = "Tags"
	case OBPackingAccessoryTypeTissuePaper:
		name = "Tissue Paper"
	case OBPackingAccessoryTypeButterPaper:
		name = "Butter Paper"
	case OBPackingAccessoryTypePlasticClip:
		name = "Plastic clip"
	case OBPackingAccessoryTypePaperBoard:
		name = "Paper board"
	case OBPackingAccessoryTypeButterfly:
		name = "Butterfly"
	case OBPackingAccessoryTypeShirtCollarSupport:
		name = "Shirt Collar Support"
	case OBPackingAccessoryTypeShirtBackSupport:
		name = "Shirt back support"
	case OBPackingAccessoryTypeTagPin:
		name = "Tag pin"
	case OBPackingAccessoryTypePriceTag:
		name = "Price Tag"
	case OBPackingAccessoryTypeBallHeadPin:
		name = "Ball head pin"
	case OBPackingAccessoryTypeInnerBox:
		name = "Inner box"
	case OBPackingAccessoryTypeFoam:
		name = "Foam"
	case OBPackingAccessoryTypeClip:
		name = "Clip"
	case OBPackingAccessoryTypePlasticAdjuster:
		name = "Plastic adjuster"
	case OBPackingAccessoryTypeShirtBox:
		name = "Shirt box"

	}

	return name
}
