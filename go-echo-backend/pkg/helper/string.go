package helper

import "github.com/gosimple/slug"

func GenerateSlug(text string, lang ...string) string {
	var defaultLang = "en"
	if len(lang) > 0 {
		defaultLang = lang[0]
	}

	return slug.MakeLang(text, defaultLang)
}
