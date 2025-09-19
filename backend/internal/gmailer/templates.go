package gmailer

import (
	"bytes"
	"html/template"
	"path/filepath"
)

// TemplateData is dynamic data for placeholders
type TemplateData map[string]interface{}

// LoadTemplate fills an HTML template with dynamic data
func LoadTemplate(templateDir, name string, data TemplateData) (string, error) {
	path := filepath.Join(templateDir, name+".html")
	tmpl, err := template.ParseFiles(path)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}
