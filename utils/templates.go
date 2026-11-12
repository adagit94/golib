package utils

import (
	"text/template"
)

// Simplified wrapper around template.Must(template.New(name).Parse(tempText)) that returns pointer to Template of passed name with parsed tempText.
func CreateTemplate(name string, tempText string) *template.Template {
	return template.Must(template.New(name).Parse(tempText))
}
