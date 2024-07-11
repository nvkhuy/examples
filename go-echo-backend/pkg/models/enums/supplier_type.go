package enums

type SupplierType string

var (
	SupplierTypeManufacturer    SupplierType = "manufacturer"
	SupplierTypeMill            SupplierType = "mill"
	SupplierTypeAccessory       SupplierType = "accessory"
	SupplierTypeService         SupplierType = "service"
	SupplierTypeProductDesigner SupplierType = "product_designer"
)

func (p SupplierType) String() string {
	return string(p)
}

func (p SupplierType) DisplayName() string {
	var name = string(p)

	switch p {
	case SupplierTypeManufacturer:
		name = "Manufacturer"
	case SupplierTypeMill:
		name = "Mill"
	case SupplierTypeAccessory:
		name = "Accessory"
	case SupplierTypeService:
		name = "Service"
	case SupplierTypeProductDesigner:
		name = "Product Designer"
	}

	return name
}
