package utils

import (
	"strings"
	"text/template"
)

func RenderingTemplate[T any](templateName, templateString string, data T) (string, error) {
	parse, err := template.New(templateName).Parse(templateString)
	if err != nil {
		return "", err
	}
	builder := strings.Builder{}
	err = parse.Execute(&builder, data)
	if err != nil {
		return "", err
	}
	return builder.String(), nil
}
