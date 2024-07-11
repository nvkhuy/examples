package enums

type OBOutputUnit string

var (
	OBOutputUnitPiece OBOutputUnit = "piece"
	OBOutputUnitMet   OBOutputUnit = "met"
	OBOutputUnitTon   OBOutputUnit = "ton"
)

func (p OBOutputUnit) String() string {
	return string(p)
}

func (p OBOutputUnit) DisplayName() string {
	var name = string(p)

	switch p {
	case OBOutputUnitPiece:
		name = "Piece"
	case OBOutputUnitMet:
		name = "Met"
	case OBOutputUnitTon:
		name = "Ton"
	}

	return name
}
