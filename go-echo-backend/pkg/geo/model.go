package geo

import "github.com/engineeringinflow/inflow-backend/pkg/models/enums"

type PlaceIndex string

var (
	PlaceIndexGrab    PlaceIndex = "Inflow-Grab"
	PlaceIndexDefault PlaceIndex = "Inflow"
)

type Address struct {
	Lat *float64 `gorm:"type:decimal(10,8);default:null" json:"lat"`
	Lng *float64 `gorm:"type:decimal(11,8);default:null" json:"lng"`

	FormattedAddress string            `json:"formatted_address"`
	City             string            `json:"city"`
	Country          enums.CountryName `json:"country"`
	CountryCode      enums.CountryCode `gorm:"default:'VN'" json:"country_code"`
	PostalCode       string            `json:"postal_code"`
	County           string            `json:"county"`
	District         string            `json:"district"`
	Street           string            `json:"street"`
	Number           string            `json:"number"`
	Neighborhood     string            `json:"neighborhood"`
	State            string            `json:"state"`
	StateCode        string            `json:"state_code"`
	BuildingName     string            `json:"building_name"`
	PlaceID          string            `json:"place_id"`
}

type LatLng struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}
