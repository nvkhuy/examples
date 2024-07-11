package enums

type BrandType string

var (
	BrandTypeIndividual BrandType = "individual"
	BrandTypeBrand      BrandType = "brand"
)

func (p BrandType) String() string {
	return string(p)
}

func (p BrandType) DisplayName() string {
	var name = string(p)

	switch p {
	case BrandTypeBrand:
		name = "Brand"
	case BrandTypeIndividual:
		name = "Individual"
	}

	return name
}
