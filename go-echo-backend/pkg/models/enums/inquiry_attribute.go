package enums

type InquiryAttribute string

var (
	InquiryAttributeSizeList      InquiryAttribute = "size_list"  // seperate by command
	InquiryAttributeSizeChart     InquiryAttribute = "size_chart" // US, UK
	InquiryAttributeComposition   InquiryAttribute = "composition"
	InquiryAttributeStyleNo       InquiryAttribute = "style_no"
	InquiryAttributeFabricName    InquiryAttribute = "fabric_name"
	InquiryAttributeFabricWeight  InquiryAttribute = "fabric_weight"
	InquiryAttributeColorCount    InquiryAttribute = "color_count"
	InquiryAttributeColorList     InquiryAttribute = "color_list"
	InquiryAttributeProductWeight InquiryAttribute = "product_weight" // Unit: kg
)
