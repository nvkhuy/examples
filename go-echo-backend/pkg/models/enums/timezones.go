package enums

import "time"

type Timezone string

func (tz Timezone) GetLocation() *time.Location {
	loc, _ := time.LoadLocation(string(tz))

	return loc
}

func (tz Timezone) Now() time.Time {
	return time.Now().In(tz.GetLocation())
}

func (tz Timezone) String() string {
	return string(tz)
}

func (tz Timezone) DefaultIfInvalid() Timezone {
	if tz == "" {
		return Timezone(CountryCodeUS.GetTimezone())
	}

	return tz
}

func (tz Timezone) ZoneOffset() int {
	_, offset := tz.Now().Zone()

	return offset
}
