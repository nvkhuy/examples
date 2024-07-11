package errs

import (
	"strings"

	"github.com/rotisserie/eris"
)

var jsonFormat = eris.NewDefaultJSONFormat(eris.FormatOptions{
	InvertOutput: true, // flag that inverts the error output (wrap errors shown first)
	WithTrace:    true, // flag that enables stack trace output
	InvertTrace:  true, // flag that inverts the stack trace output (top of call stack shown first)
	WithExternal: true,
})

var stringFormat = eris.NewDefaultStringFormat(eris.FormatOptions{
	InvertOutput: true, // flag that inverts the error output (wrap errors shown first)
	WithTrace:    true, // flag that enables stack trace output
	InvertTrace:  true, // flag that inverts the stack trace output (top of call stack shown first)
	WithExternal: false,
})

// FormatErisJSON format json
func FormatErisJSON(err error) map[string]interface{} {

	return eris.ToCustomJSON(err, jsonFormat)
}

func ParseErisJSON(err error) (message string, stack []string) {
	var result = eris.ToCustomJSON(err, jsonFormat)
	message = err.Error()

	if v, ok := result["message"].(string); ok {
		message = v
	}

	if m, ok := result["root"].(map[string]interface{}); ok {
		if msg, ok := m["message"].(string); ok {
			message = msg
		}

		if st, ok := m["stack"].([]string); ok {
			for i := 0; i < len(st); i++ {
				if i > 6 {
					return
				}
				stack = append(stack, strings.Replace(st[i], "/app/", "", 1))
			}

		}

	}

	return
}

// FormatErisString format string
func FormatErisString(err error) string {

	return eris.ToCustomString(err, stringFormat)
}
