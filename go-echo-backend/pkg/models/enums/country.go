package enums

import (
	"strings"
)

type CountryTimezone = map[CountryCode]Timezone

type CountryCurrency = map[CountryCode]Currency

type CountryNameData = map[CountryCode]CountryName

type CountryName string

type CountryCode string

type CountryCodes []CountryCode

func (countries CountryCodes) Contains(c CountryCode) bool {
	for _, v := range countries {
		if v == c {
			return true
		}
	}

	return false
}

const (
	CountryNameUS CountryName = "United State"
	CountryNameSG CountryName = "Singapore"
	CountryNameVN CountryName = "Viet Nam"
)

const (
	CountryCodeUS CountryCode = "US"
	CountryCodeSG CountryCode = "SG"
	CountryCodeVN CountryCode = "VN"
)

var AsiaCountries = CountryCodes{
	CountryCodeSG,
	CountryCodeVN,
}
var CountryTimezones CountryTimezone = map[CountryCode]Timezone{
	CountryCodeUS: "Asia/Singapore",
}

var CountryNames CountryNameData = map[CountryCode]CountryName{
	CountryCodeUS: CountryNameUS,
}

var CountryCurrencies CountryCurrency = map[CountryCode]Currency{
	CountryCodeUS: USD,
	CountryCodeVN: VND,
}

func (c CountryCode) GetTimezone() Timezone {
	return CountryTimezones[c]
}

func (c CountryCode) GetCurrency() Currency {
	return CountryCurrencies[c]
}

func (c CountryCode) GetCountryName() CountryName {
	return CountryNames[c]
}

func (c CountryCode) String() string {
	return string(c)
}

func (c CountryName) String() string {
	return string(c)
}

func (c CountryCode) ToUpper() CountryCode {
	return CountryCode(strings.ToUpper(c.String()))
}

func (c CountryCode) ToLower() string {
	return strings.ToLower(c.String())
}

func (c CountryCode) DefaultIfInvalid() CountryCode {
	if c == "" {
		return CountryCodeUS
	}

	return c
}

func (c CountryName) GetCountryCode() CountryCode {
	if strings.EqualFold(c.String(), CountryNameSG.String()) {
		return CountryCodeSG
	}

	return CountryCode(c.String())
}
