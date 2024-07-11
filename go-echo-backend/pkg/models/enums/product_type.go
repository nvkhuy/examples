package enums

type ProductType string

var (
	ProductTypeClothing ProductType = "clothing"
	ProductTypeFabric   ProductType = "fabric"
	ProductTypeGraphic  ProductType = "graphic"
)

func (p ProductType) String() string {
	return string(p)
}

func (p ProductType) DisplayName() string {
	var name = string(p)

	switch p {
	case ProductTypeClothing:
		name = "Clothing"
	case ProductTypeFabric:
		name = "Fabric"
	case ProductTypeGraphic:
		name = "Graphic"
	}

	return name
}
