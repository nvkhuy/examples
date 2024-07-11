package main

import (
	"errors"

	"github.com/brianvoe/sjwt"
)

func IncludeStr(slice []string, value string) bool {
	for _, v := range slice {
		if v == value {
			return true
		}
	}

	return false
}
func ValidateToken(token string, fileKey string) error {
	var valid = sjwt.Verify(token, []byte(config.JwtSecret))
	if !valid {
		return errors.New("token is invalid")
	}

	jwt, err := sjwt.Parse(token)
	if err != nil {
		return err
	}

	return jwt.Validate()
}
