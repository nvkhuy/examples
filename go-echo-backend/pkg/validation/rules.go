package validation

import (
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/nyaruka/phonenumbers"
)

func isURL(fl validator.FieldLevel) bool {
	url := fl.Field().String()
	if len(url) > 0 {
		return strings.HasPrefix(url, "http")
	}
	return true
}

func isBirthday(fl validator.FieldLevel) bool {
	layout := "2006-01-02"
	birthday := fl.Field().String()
	_, err := time.Parse(layout, birthday)
	return err == nil
}

func isPhone(fl validator.FieldLevel) bool {
	phone := fl.Field().String()
	num, err := phonenumbers.Parse(phone, "SG")
	if err != nil || num == nil {
		return false
	}
	return phonenumbers.IsValidNumber(num)
}

func isUnix(fl validator.FieldLevel) bool {
	ts := fl.Field().Int()
	if ts == 0 {
		return false
	}
	return true
}

func isTimezone(fl validator.FieldLevel) bool {
	tz := fl.Field().String()
	if len(tz) > 0 {
		_, err := time.LoadLocation(tz)
		return err == nil
	}
	return true
}
