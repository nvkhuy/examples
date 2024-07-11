package enums

type InquiryPriceType string

var (
	InquiryPriceTypeFOB InquiryPriceType = "fob"
	InquiryPriceTypeCIF InquiryPriceType = "cif"
	InquiryPriceTypeEXW InquiryPriceType = "exw"
)

func (p InquiryPriceType) String() string {
	return string(p)
}

func (p InquiryPriceType) DisplayName() string {
	var name = string(p)

	switch p {
	case InquiryPriceTypeFOB:
		name = "FOB"
	case InquiryPriceTypeCIF:
		name = "CIF"
	case InquiryPriceTypeEXW:
		name = "EXW"
	}

	return name
}
