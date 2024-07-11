package enums

type LabelMaterial string

var (
	LabelMaterialC200   LabelMaterial = "c_200"
	LabelMaterialC250   LabelMaterial = "c_250"
	LabelMaterialC300   LabelMaterial = "c_300"
	LabelMaterialB200   LabelMaterial = "b_200"
	LabelMaterialB250   LabelMaterial = "b_250"
	LabelMaterialPoly   LabelMaterial = "poly"
	LabelMaterialCotton LabelMaterial = "cotton"
)

func (p LabelMaterial) String() string {
	return string(p)
}
