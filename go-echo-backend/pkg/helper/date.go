package helper

import (
	"math"
	"time"

	"github.com/engineeringinflow/inflow-backend/pkg/models/enums"
)

// DefaultTimezone timezone
const DefaultTimezone enums.Timezone = "Asia/Saigon"

// RoundTime round
func RoundTime(input float64) int {
	var result float64

	if input < 0 {
		result = math.Ceil(input - 0.5)
	} else {
		result = math.Floor(input + 0.5)
	}

	// only interested in integer, ignore fractional
	i, _ := math.Modf(result)

	return int(i)
}

// DiffDays days
func DiffDays(d time.Duration) int {
	return RoundTime(d.Seconds() / 86400)
}

// DiffWeeks weeks
func DiffWeeks(d time.Duration) int {
	return RoundTime(d.Seconds() / 604800)
}

// DiffMonths months
func DiffMonths(d time.Duration) int {
	return RoundTime(d.Seconds() / 2600640)
}

// DiffYears years
func DiffYears(d time.Duration) int {
	return RoundTime(d.Seconds() / 31207680)
}

// TimeToTimeString to time string
func TimeToTimeString(t time.Time, format ...string) string {
	var layout = "15:04:05"
	if len(format) > 0 {
		layout = format[0]
	}

	return t.Format(layout)
}

// TimeToDateString to time string
func TimeToDateString(t time.Time, format ...string) string {
	var layout = "2006-01-02 15:04:05"
	if len(format) > 0 {
		layout = format[0]
	}

	return t.Format(layout)
}

// CompareBeforeTime compare t1 > t2
func CompareBeforeTime(t1 string, t2 string) (result bool, err error) {
	var layout = "15:04:05"
	tc1, err := time.Parse(layout, t1)
	if err != nil {
		return false, err
	}
	tc2, err := time.Parse(layout, t2)
	if err != nil {
		return false, err
	}

	var value = tc1.Unix() > tc2.Unix()
	return value, nil
}

func IsSameDate(t1 time.Time, t2 time.Time) bool {
	return t1.Year() == t2.Year() && t1.Month() == t2.Month() && t1.Day() == t2.Day()
}
