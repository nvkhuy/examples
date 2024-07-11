package enums

type ProductUnit string

var (
	ProductUnitPiece ProductUnit = "piece"
	ProductUnitPair  ProductUnit = "pair"
	ProductUnitBox   ProductUnit = "box"
)

func (p ProductUnit) String() string {
	return string(p)
}

func (p ProductUnit) DisplayName() string {
	var name = string(p)

	switch p {
	case ProductUnitPiece:
		name = "Piece"
	case ProductUnitPair:
		name = "Pair"
	case ProductUnitBox:
		name = "Box"
	}

	return name
}
