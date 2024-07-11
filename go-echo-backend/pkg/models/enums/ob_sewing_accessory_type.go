package enums

type OBSewingAccessoryType string

var (
	OBSewingAccessoryTypeButton          OBSewingAccessoryType = "button"
	OBSewingAccessoryTypeZipper          OBSewingAccessoryType = "zipper"
	OBSewingAccessoryTypeLining          OBSewingAccessoryType = "lining"
	OBSewingAccessoryTypeInterlining     OBSewingAccessoryType = "interlining"
	OBSewingAccessoryTypeSnapButton      OBSewingAccessoryType = "snap_button"
	OBSewingAccessoryTypeHaspsAndSlider  OBSewingAccessoryType = "hasps_and_slider"
	OBSewingAccessoryTypeEmbroidery      OBSewingAccessoryType = "embroidery"
	OBSewingAccessoryTypeApplique        OBSewingAccessoryType = "applique"
	OBSewingAccessoryTypeBeads           OBSewingAccessoryType = "beads"
	OBSewingAccessoryTypeGlitter         OBSewingAccessoryType = "hlitter"
	OBSewingAccessoryTypeRhinestones     OBSewingAccessoryType = "rhinestones"
	OBSewingAccessoryTypeSequins         OBSewingAccessoryType = "sequins"
	OBSewingAccessoryTypeDrawstring      OBSewingAccessoryType = "drawstring"
	OBSewingAccessoryTypeWaistTies       OBSewingAccessoryType = "waist_ties"
	OBSewingAccessoryTypeBows            OBSewingAccessoryType = "bows"
	OBSewingAccessoryTypeFringe          OBSewingAccessoryType = "fringe"
	OBSewingAccessoryTypePomPom          OBSewingAccessoryType = "pom_pom"
	OBSewingAccessoryTypeTassel          OBSewingAccessoryType = "tassel"
	OBSewingAccessoryTypeLabel           OBSewingAccessoryType = "label"
	OBSewingAccessoryTypeMainLabel       OBSewingAccessoryType = "main_label"
	OBSewingAccessoryTypePULabel         OBSewingAccessoryType = "pu_label"
	OBSewingAccessoryTypePatch           OBSewingAccessoryType = "patch"
	OBSewingAccessoryTypeHookAndLoop     OBSewingAccessoryType = "hook_and_loop"
	OBSewingAccessoryTypeEyeletOrGrommet OBSewingAccessoryType = "eyelet_or_grommet"
	OBSewingAccessoryTypeHookAndEye      OBSewingAccessoryType = "hook_and_eye"
	OBSewingAccessoryTypePadding         OBSewingAccessoryType = "padding"
	OBSewingAccessoryTypeElastic         OBSewingAccessoryType = "elastic"
	OBSewingAccessoryTypeLaceFabric      OBSewingAccessoryType = "lace_fabric"
	OBSewingAccessoryTypeTwillTape       OBSewingAccessoryType = "twill_tape"
	OBSewingAccessoryTypeRib             OBSewingAccessoryType = "rib"
	OBSewingAccessoryTypeBelt            OBSewingAccessoryType = "belt"
	OBSewingAccessoryTypeStrapping       OBSewingAccessoryType = "strapping"
)

func (p OBSewingAccessoryType) String() string {
	return string(p)
}

func (p OBSewingAccessoryType) DisplayName() string {
	var name = string(p)

	switch p {
	case OBSewingAccessoryTypeButton:
		name = "Button"
	case OBSewingAccessoryTypeZipper:
		name = "Zipper"
	case OBSewingAccessoryTypeLining:
		name = "Lining"
	case OBSewingAccessoryTypeInterlining:
		name = "Interlining"
	case OBSewingAccessoryTypeSnapButton:
		name = "Snap Button"
	case OBSewingAccessoryTypeHaspsAndSlider:
		name = "Hasps and Slider"
	case OBSewingAccessoryTypeEmbroidery:
		name = "Embroidery"
	case OBSewingAccessoryTypeApplique:
		name = "Applique"
	case OBSewingAccessoryTypeBeads:
		name = "Beads"
	case OBSewingAccessoryTypeGlitter:
		name = "Glitter"
	case OBSewingAccessoryTypeRhinestones:
		name = "Rhinestones"
	case OBSewingAccessoryTypeSequins:
		name = "Sequins"
	case OBSewingAccessoryTypeDrawstring:
		name = "Drawstring"
	case OBSewingAccessoryTypeWaistTies:
		name = "Waist Ties"
	case OBSewingAccessoryTypeBows:
		name = "Bows"
	case OBSewingAccessoryTypeFringe:
		name = "Fringe"
	case OBSewingAccessoryTypePomPom:
		name = "Pom Pom"
	case OBSewingAccessoryTypeTassel:
		name = "Tassel"
	case OBSewingAccessoryTypeLabel:
		name = "Label"
	case OBSewingAccessoryTypeMainLabel:
		name = "Main Label"
	case OBSewingAccessoryTypePULabel:
		name = "PU Label"
	case OBSewingAccessoryTypeHookAndLoop:
		name = "Hook and Loop"
	case OBSewingAccessoryTypeEyeletOrGrommet:
		name = "Eyelet or Grommet"
	case OBSewingAccessoryTypeHookAndEye:
		name = "Hook and Eye"
	case OBSewingAccessoryTypePadding:
		name = "Padding"
	case OBSewingAccessoryTypeElastic:
		name = "Elastic"
	case OBSewingAccessoryTypeLaceFabric:
		name = "Lace Fabric"
	case OBSewingAccessoryTypeTwillTape:
		name = "Twill Tape"
	case OBSewingAccessoryTypeRib:
		name = "Rib"
	case OBSewingAccessoryTypeBelt:
		name = "Belt"
	case OBSewingAccessoryTypeStrapping:
		name = "Strapping"

	}

	return name
}
