package enums

type CategoryType string

var (
	CategoryOfProduct CategoryType = "product"
	CategoryOfDesigner  CategoryType = "designer"
)

func (p CategoryType) String() string {
	return string(p)
}

