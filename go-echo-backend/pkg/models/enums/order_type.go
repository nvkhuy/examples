package enums

type OrderGroupType string

var (
	OrderGroupTypeRFQ    OrderGroupType = "RFQ"
	OrderGroupTypeSample OrderGroupType = "Sample"
	OrderGroupTypeBulk   OrderGroupType = "Bulk"
)

func (status OrderGroupType) String() string {
	return string(status)
}

func (status OrderGroupType) DisplayName() string {
	var name = string(status)

	switch status {
	case OrderGroupTypeRFQ:
		name = "RFQ"
	case OrderGroupTypeSample:
		name = "Sample"
	case OrderGroupTypeBulk:
		name = "Bulk"
	}

	return name
}
