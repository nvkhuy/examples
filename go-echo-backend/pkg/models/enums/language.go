package enums

import "strings"

type LanguageCode string

var (
	LanguageCodeEnglish   LanguageCode = "en"
	LanguageCodeIndonesia LanguageCode = "id"
	LanguageCodeVietnam   LanguageCode = "vi"
)

func (l LanguageCode) String() string {
	return string(l)
}

func (l LanguageCode) ToLower() string {
	return strings.ToLower(string(l))
}

func (l LanguageCode) DefaultIfInvalid() LanguageCode {
	if l == "" {
		return LanguageCodeEnglish
	}

	return l
}
