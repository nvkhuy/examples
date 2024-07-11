package models

import (
	"fmt"
	"strings"

	"github.com/engineeringinflow/inflow-backend/pkg/geo"
	"github.com/engineeringinflow/inflow-backend/pkg/location"
	"github.com/engineeringinflow/inflow-backend/pkg/models/enums"
)

type Coordinate struct {
	Model
	location.Coordinate `gorm:"embedded"`

	// Address: just use for query in BE
	AddressID   string            `gorm:"-" json:"address_id,omitempty"`
	AddressName string            `gorm:"-" json:"address_name,omitempty"`
	AddressType enums.AddressType `gorm:"-" json:"address_type,omitempty"`
}

func (coordinate *Coordinate) GetLatLng() {
	var input = coordinate.FormattedAddress
	if input == "" {
		input = fmt.Sprintf("%s %s, %s, %s, %s, %s", coordinate.AddressNumber, coordinate.Street, coordinate.Level3, coordinate.Level2, coordinate.Level1, coordinate.CountryCode.GetCountryName())
	}

	result, err := geo.GetInstance().SearchPlaceIndexForText(geo.SearchPlaceIndexForTextParams{
		Address:     input,
		CountryCode: coordinate.CountryCode,
		Language:    enums.LanguageCodeVietnam,
	})
	if err == nil && len(result) > 0 {
		coordinate.Lat = result[0].Lat
		coordinate.Lng = result[0].Lng
	}
}
func (coordinate *Coordinate) Display() string {
	var input = coordinate.FormattedAddress
	if input == "" {
		var parts []string

		if coordinate.AddressNumber != "" && coordinate.Street != "" {
			parts = append(parts, fmt.Sprintf("%s %s", coordinate.AddressNumber, coordinate.Street))
		} else {
			if coordinate.AddressNumber != "" {
				parts = append(parts, coordinate.AddressNumber)
			}

			if coordinate.Street != "" {
				parts = append(parts, coordinate.Street)
			}
		}

		if coordinate.Level3 != "" {
			parts = append(parts, coordinate.Level3)
		}

		if coordinate.Level2 != "" {
			parts = append(parts, coordinate.Level2)
		}
		if coordinate.Level1 != "" {
			parts = append(parts, coordinate.Level1)
		}

		if coordinate.CountryCode.GetCountryName() != "" {
			parts = append(parts, coordinate.CountryCode.GetCountryName().String())
		}

		input = strings.Join(parts, ", ")
	}
	return input
}
