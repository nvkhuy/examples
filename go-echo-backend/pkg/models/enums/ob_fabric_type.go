package enums

// Using for manufacturer seller
type OBFabricType string

var (
	OBFabricTypeSequin    OBFabricType = "sequin"
	OBFabricTypeVelvet    OBFabricType = "velvet"
	OBFabricTypeSatin     OBFabricType = "satin"
	OBFabricTypeSilk      OBFabricType = "silk"
	OBFabricTypePuLeather OBFabricType = "pu_leather"
	OBFabricTypeFur       OBFabricType = "fur"
	OBFabricTypeKnit      OBFabricType = "knit"
	OBFabricTypeWoven     OBFabricType = "woven"
	OBFabricTypeDenim     OBFabricType = "denim"
	OBFabricTypeThickness OBFabricType = "thickness"
)

func (p OBFabricType) String() string {
	return string(p)
}

func (p OBFabricType) DisplayName() string {
	var name = string(p)

	switch p {
	case OBFabricTypeSequin:
		name = "Sequin"
	case OBFabricTypeVelvet:
		name = "Velvet"
	case OBFabricTypeSatin:
		name = "Satin"
	case OBFabricTypeSilk:
		name = "Silk"
	case OBFabricTypePuLeather:
		name = "PuLeather"
	case OBFabricTypeFur:
		name = "Fur"
	case OBFabricTypeKnit:
		name = "Knit"
	case OBFabricTypeWoven:
		name = "Woven"
	case OBFabricTypeDenim:
		name = "Denim"
	case OBFabricTypeThickness:
		name = "Thickness"
	}

	return name
}
