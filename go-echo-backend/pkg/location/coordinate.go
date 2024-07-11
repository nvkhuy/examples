package location

import (
	"encoding/json"

	"github.com/engineeringinflow/inflow-backend/pkg/models/enums"
)

type Coordinate struct {
	Lat *float64 `gorm:"type:decimal(10,8);default:null" json:"lat,omitempty"`
	Lng *float64 `gorm:"type:decimal(11,8);default:null" json:"lng,omitempty"`

	AddressNumber    string `json:"address_number,omitempty"`
	FormattedAddress string `json:"formatted_address,omitempty"`
	Street           string `json:"street,omitempty"`
	Level1           string `gorm:"column:level_1" json:"level_1,omitempty" `
	Level2           string `gorm:"column:level_2" json:"level_2,omitempty"`
	Level3           string `gorm:"column:level_3" json:"level_3,omitempty"`
	Level4           string `gorm:"column:level_4" json:"level_4,omitempty"`

	PostalCode string `json:"postal_code"`

	CountryCode enums.CountryCode `gorm:"default:'VN'" json:"country_code"`

	TimezoneName   enums.Timezone `json:"timezone_name,omitempty"`
	TimezoneOffset int            `json:"timezone_offset,omitempty"`
	PlaceID        string         `json:"place_id,omitempty"`
}

func (coordinate *Coordinate) ToJsonString() string {
	data, err := json.Marshal(coordinate)
	if err != nil {
		return ""
	}

	return string(data)
}
