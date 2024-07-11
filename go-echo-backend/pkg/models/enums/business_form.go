package enums

type BusinessForm string

var (
	BusinessFormExport BusinessForm = "export"
	BusinessFormLocal  BusinessForm = "local"
)

func (p BusinessForm) String() string {
	return string(p)
}

func (p BusinessForm) DisplayName() string {
	var name = string(p)

	switch p {
	case BusinessFormExport:
		name = "Export"
	case BusinessFormLocal:
		name = "Local"
	}

	return name
}
