package enums

type LabelDimension string

var (
	LabelDimension4060MM LabelDimension = "40_60_mm"
	LabelDimension6070MM LabelDimension = "60_70_mm"
	LabelDimension2713MM LabelDimension = "27_13_mm"
)

func (p LabelDimension) String() string {
	return string(p)
}
