package enums

type AdsVideoSection string

var (
	AdsVideoSectionRFQ       AdsVideoSection = "rfq"
	AdsVideoSectionSample    AdsVideoSection = "sample"
	AdsVideoSectionBulk      AdsVideoSection = "bulk"
	AdsVideoSectionCatalogue AdsVideoSection = "catalogue"
)

func (p AdsVideoSection) String() string {
	return string(p)
}

func (p AdsVideoSection) DisplayName() string {
	var name = string(p)

	switch p {
	case AdsVideoSectionRFQ:
		name = "RFQ"
	case AdsVideoSectionSample:
		name = "Sample"
	case AdsVideoSectionBulk:
		name = "Bulk"
	case AdsVideoSectionCatalogue:
		name = "Catalogue"
	}

	return name
}
