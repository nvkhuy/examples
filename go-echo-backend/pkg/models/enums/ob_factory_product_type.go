package enums

type OBFactoryProductType string

var (
	OBFactoryProductTypeBlouse            OBFactoryProductType = "blouse"
	OBFactoryProductTypeBodysuit          OBFactoryProductType = "bodysuit"
	OBFactoryProductTypeCamisole          OBFactoryProductType = "camisole"
	OBFactoryProductTypeCrochetKnitted    OBFactoryProductType = "crochet_knitted"
	OBFactoryProductTypeDress             OBFactoryProductType = "dress"
	OBFactoryProductTypeHoodie            OBFactoryProductType = "hoodie"
	OBFactoryProductTypeJacketCoatBlaze   OBFactoryProductType = "jacket_coat_blaze"
	OBFactoryProductTypeJumpsuit          OBFactoryProductType = "jumpsuit"
	OBFactoryProductTypeLegging           OBFactoryProductType = "legging"
	OBFactoryProductTypePant              OBFactoryProductType = "pant"
	OBFactoryProductTypePoloShirt         OBFactoryProductType = "polo_shirt"
	OBFactoryProductTypeSweatShirt        OBFactoryProductType = "sweat_shirt"
	OBFactoryProductTypeTShirt            OBFactoryProductType = "t_shirt"
	OBFactoryProductTypeSkirt             OBFactoryProductType = "skirt"
	OBFactoryProductTypeSkort             OBFactoryProductType = "skort"
	OBFactoryProductTypeTankTop           OBFactoryProductType = "tank_top"
	OBFactoryProductTypeUnderwearSwimwear OBFactoryProductType = "underwear_swimwear"
	OBFactoryProductTypeVestSuit          OBFactoryProductType = "vest_suit"
	OBFactoryProductTypeCrochetKnitting   OBFactoryProductType = "crochet_knitting"
)

func (p OBFactoryProductType) String() string {
	return string(p)
}

func (p OBFactoryProductType) DisplayName() string {
	var name = string(p)

	switch p {
	case OBFactoryProductTypeBlouse:
		name = "Blouse"
	case OBFactoryProductTypeBodysuit:
		name = "Bodysuit"
	case OBFactoryProductTypeCamisole:
		name = "Camisole"
	case OBFactoryProductTypeCrochetKnitted:
		name = "Crochet/Knitted"
	case OBFactoryProductTypeDress:
		name = "Dress"
	case OBFactoryProductTypeHoodie:
		name = "Hoodie"
	case OBFactoryProductTypeJacketCoatBlaze:
		name = "Jacket/Coat/Blaze"
	case OBFactoryProductTypeJumpsuit:
		name = "Jumpsuit"
	case OBFactoryProductTypeLegging:
		name = "Legging"
	case OBFactoryProductTypePant:
		name = "Pant"
	case OBFactoryProductTypePoloShirt:
		name = "Polo Shirt"
	case OBFactoryProductTypeSweatShirt:
		name = "Sweat Shirt"
	case OBFactoryProductTypeTShirt:
		name = "T-Shirt"
	case OBFactoryProductTypeSkirt:
		name = "Skirt"
	case OBFactoryProductTypeSkort:
		name = "Skort"
	case OBFactoryProductTypeTankTop:
		name = "Tank top"
	case OBFactoryProductTypeUnderwearSwimwear:
		name = "Underwear/Swimwear"
	case OBFactoryProductTypeVestSuit:
		name = "Vest/Suit"
	case OBFactoryProductTypeCrochetKnitting:
		name = "Crochet knitting"
	}

	return name
}
