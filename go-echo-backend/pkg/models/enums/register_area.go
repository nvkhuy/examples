package enums

type RegisterArea string

var (
	AreaUS            RegisterArea = "us"
	AreaEU            RegisterArea = "eu"
	AreaJapan         RegisterArea = "japan"
	AreaKorea         RegisterArea = "korea"
	AreaSoutheastAsia RegisterArea = "southeast_asia"
	AreaChina         RegisterArea = "china"
	AreaIndia         RegisterArea = "india"
)

func (register_area RegisterArea) String() string {
	return string(register_area)
}

func (register_area RegisterArea) IconUrl() string {
	return string("https://dev-static.joininflow.io/common/register_area/" + register_area + ".png")
}

func (register_area RegisterArea) DisplayName() string {
	switch register_area {
	case AreaUS:
		return "US"

	case AreaEU:
		return "EU"

	case AreaJapan:
		return "Japan"

	case AreaKorea:
		return "Korea"

	case AreaSoutheastAsia:
		return "Southeast Asia"

	case AreaChina:
		return "China"

	case AreaIndia:
		return "India"
	}

	return string(register_area)
}
