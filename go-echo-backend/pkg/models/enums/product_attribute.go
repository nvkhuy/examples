package enums

type ProductAttribute string

var (
	ProductAttributeSize  ProductAttribute = "size"
	ProductAttributeColor ProductAttribute = "color"
)

func (p ProductAttribute) String() string {
	return string(p)
}
