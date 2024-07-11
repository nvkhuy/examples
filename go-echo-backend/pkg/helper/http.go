package helper

import (
	"fmt"
	"net/http"
	"sort"
	"strings"
)

var quoteEscaper = strings.NewReplacer("\"", "", "\t", "", "\n", " ")
var quoteJsonEscaper = strings.NewReplacer("\"", `"`)

func EscapeQuotes(s string) string {
	return quoteEscaper.Replace(s)
}

func EscapeJsonQuotes(s string) string {
	return quoteJsonEscaper.Replace(s)
}

func ComposeHeaders(hdrs http.Header) string {
	str := make([]string, 0, len(hdrs))
	for _, k := range SortHeaderKeys(hdrs) {
		var v string
		if k == "Cookie" {
			cv := strings.TrimSpace(strings.Join(hdrs[k], ", "))

			v = strings.TrimSpace(fmt.Sprintf("%25s: %s", k, cv))
		} else {
			v = strings.TrimSpace(fmt.Sprintf("%25s: %s", k, strings.Join(hdrs[k], ", ")))
		}
		if v != "" {
			str = append(str, "\t"+v)
		}
	}
	return strings.Join(str, "\n")
}

func SortHeaderKeys(hdrs http.Header) []string {
	keys := make([]string, 0, len(hdrs))
	for key := range hdrs {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}

func CopyHeaders(hdrs http.Header) http.Header {
	nh := http.Header{}
	for k, v := range hdrs {
		nh[k] = v
	}
	return nh
}

func FormatHeaders(hdrs http.Header) map[string]string {
	var result = map[string]string{}

	for k, v := range hdrs {
		result[k] = strings.Join(v, " ")
	}
	return result
}
