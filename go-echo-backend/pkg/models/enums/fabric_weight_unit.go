package enums

type FabricWeightUnit string

var (
	FabricWeightUnitGSM FabricWeightUnit = "gsm"
	FabricWeightUnitOZ  FabricWeightUnit = "oz"
)

func (p FabricWeightUnit) String() string {
	return string(p)
}

func (p FabricWeightUnit) DisplayName() string {
	var name = string(p)

	switch p {
	case FabricWeightUnitGSM:
		name = "GSM"
	case FabricWeightUnitOZ:
		name = "OZ"
	}

	return name
}
