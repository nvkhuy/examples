package helper

import (
	"github.com/thaitanloi365/go-utils/random"
)

func GenerateUserDummyPassword() string {
	var id = random.String(8, random.Alphanumeric)
	return id
}
