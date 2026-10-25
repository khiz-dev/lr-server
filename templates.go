package main

import (
	"fmt"
	"os"
)

type template struct {
	path string
}

type templateEnum string

const (
	AdminMetricsTemplate templateEnum = "metrics"
)

func getTemplates() (map[string]*template) {
	templates := map[string]*template{
		string(AdminMetricsTemplate): {path: "./templates/admin/metrics.html"},
	}
	return templates
}

func getTemplate(t templateEnum) (*template, error) {
	templates := getTemplates()
	tmpl, ok := templates[string(t)]
	if !ok {
		return nil, fmt.Errorf("template %s not found", t)
	}
	return tmpl, nil
}

func (t *template) render(data any) ([]byte, error) {
	// Read the template file
	content, err := os.ReadFile(t.path)
	if err != nil {
		return nil, err
	}

	// Format the content with the provided data
	renderedContent := fmt.Sprintf(string(content), data)

	return []byte(renderedContent), nil
}