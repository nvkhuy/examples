package enums

type RegisterBusiness string

var (
	FindProduct      RegisterBusiness = "find_product"
	FindManufacturer RegisterBusiness = "find_manufacturer"
	FindDesigner     RegisterBusiness = "find_designer"
	FindOther        RegisterBusiness = "find_other"
)

func (register_business RegisterBusiness) String() string {
	return string(register_business)
}

func (register_business RegisterBusiness) IconUrl() string {
	return string("https://dev-static.joininflow.io/common/" + register_business + ".png")
}

func (register_business RegisterBusiness) DisplayName() string {
	switch register_business {
	case FindProduct:
		return "Find Product"

	case FindManufacturer:
		return "Find Manufacturers"

	case FindDesigner:
		return "Find Designer"

	case FindOther:
		return "Find Something else"

	}

	return string(register_business)
}
