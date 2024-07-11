package enums

type LabelSize string

var (
	LabelSizeXS LabelSize = "xs"
	LabelSizeS  LabelSize = "s"
	LabelSizeM  LabelSize = "m"
	LabelSizeL  LabelSize = "l"
	LabelSizeXL LabelSize = "xl"
)

func (p LabelSize) String() string {
	return string(p)
}
