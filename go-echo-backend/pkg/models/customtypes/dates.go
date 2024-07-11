package customtypes

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"gorm.io/datatypes"
)

type Dates []*datatypes.Date

const (
	separator = ","
	layout    = "2006-01-02"
)

func (dates Dates) Value() (driver.Value, error) {
	var result []driver.Value
	for _, date := range dates {
		value, _ := date.Value()
		result = append(result, value)
	}
	return result, nil
}

func (dates *Dates) Scan(value interface{}) (err error) {
	fmt.Printf("value: %s", value)
	stringValue := value.(string)
	stringValue = strings.TrimSuffix(stringValue, "}")
	stringValue = strings.TrimPrefix(stringValue, "{")

	datesString := strings.Split(stringValue, separator)
	for _, dateString := range datesString {
		dateTime, _ := time.Parse(layout, dateString)
		date := datatypes.Date(dateTime)
		*dates = append(*dates, &date)
	}
	return nil
}

func (dates *Dates) UnmarshalJSON(b []byte) error {
	if string(b) == "" || string(b) == "null" {
		return nil
	}

	var datesString []string
	_ = json.Unmarshal(b, &datesString)
	for _, dateString := range datesString {
		dateTime, err := time.Parse(layout, dateString)
		if err != nil {
			return errors.New("date is invalid format, correct format is YYYY-MM-DD")
		}
		date := datatypes.Date(dateTime)
		*dates = append(*dates, &date)
	}
	return nil
}

func (dates Dates) MarshalJSON() ([]byte, error) {
	var times []string
	for _, date := range dates {
		dateTime := time.Time(*date)
		times = append(times, dateTime.Format(layout))
	}
	return convertArrayStringToBytes(times), nil
}

func convertArrayStringToBytes(strArray []string) []byte {
	result, _ := json.Marshal(strArray)
	return result
}

func (dates Dates) String() string {
	result := "["
	lenDate := len(dates)
	for index, date := range dates {
		result += time.Time(*date).Format(layout)
		if index != lenDate-1 {
			result += ","
		}
	}
	return result + "]"
}

func (dates Dates) ToTimes() []time.Time {
	result := make([]time.Time, 0)
	for _, date := range dates {
		result = append(result, time.Time(*date))
	}
	return result
}
