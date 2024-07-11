package helper

import (
	"bytes"
	"fmt"
	"html/template"

	textTemplate "text/template"
)

func ParseHTML(templateName string, i interface{}) (string, error) {
	t, err := template.ParseFiles(templateName)
	if err != nil {
		return "", err
	}

	var tpl bytes.Buffer
	if err := t.Execute(&tpl, i); err != nil {
		return "", err
	}

	return tpl.String(), nil
}

func ParseHTMLBytes(templateName string, i interface{}) ([]byte, error) {
	t, err := template.ParseFiles(templateName)
	if err != nil {
		return nil, err
	}

	var tpl bytes.Buffer
	if err := t.Execute(&tpl, i); err != nil {
		return nil, err
	}

	return tpl.Bytes(), nil
}

func TemplateStr(msg string, params interface{}, funcs ...template.FuncMap) string {
	var mergeFuncs = textTemplate.FuncMap{
		"Caller": func() string {
			return fmt.Sprintf("%s:%s:%s", GetFuncName(16), GetFuncName(17), GetFuncName(18))
		},
	}
	for _, funcMap := range funcs {
		for name, f := range funcMap {
			mergeFuncs[name] = f
		}

	}

	t := textTemplate.Must(textTemplate.New("message").Funcs(mergeFuncs).Parse(msg))
	buf := &bytes.Buffer{}
	err := t.Execute(buf, params)
	if err == nil {
		return buf.String()
	}
	return msg
}
