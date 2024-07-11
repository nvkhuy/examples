package enums

type OBDevelopmentService string

var (
	OBDevelopmentServiceDesignSketches            OBDevelopmentService = "design_sketches"
	OBDevelopmentServiceFabricMaterialCombination OBDevelopmentService = "fabric_material_combination"
)

func (p OBDevelopmentService) String() string {
	return string(p)
}

func (p OBDevelopmentService) DisplayName() string {
	var name = string(p)

	switch p {
	case OBDevelopmentServiceDesignSketches:
		name = "Design Sketches"
	case OBDevelopmentServiceFabricMaterialCombination:
		name = "Fabric & Material Combination"
	}

	return name
}
