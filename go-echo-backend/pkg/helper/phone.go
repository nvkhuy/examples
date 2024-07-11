package helper

import (
	"fmt"
	"strings"

	"github.com/engineeringinflow/inflow-backend/pkg/models/enums"
	"github.com/nyaruka/phonenumbers"
)

func FormatNationalPhoneNumber(phone string, countryCode enums.CountryCode) string {
	if strings.HasPrefix(phone, "+") {
		return phone
	}

	value, err := phonenumbers.Parse(phone, string(countryCode))
	if err != nil {
		return phone
	}

	return fmt.Sprintf("+%d%d", *value.CountryCode, *value.NationalNumber)
}

func FormatLocalPhoneNumber(phone string, countryCode enums.CountryCode) string {
	value, err := phonenumbers.Parse(phone, string(countryCode))
	if err != nil {
		return phone
	}

	return fmt.Sprintf("%d", *value.NationalNumber)
}

func ParsePhoneNumber(phone string, region string) (*phonenumbers.PhoneNumber, error) {
	return phonenumbers.Parse(phone, region)
}
