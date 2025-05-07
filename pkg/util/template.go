package util

import (
	"html/template"
	"net/http"
	"path/filepath"
)

// TemplateUtil is a utility for handling HTML templates.
type TemplateUtil struct {
	templates *template.Template
}

// NewTemplateUtil creates a new TemplateUtil instance.
func NewTemplateUtil(templateDir string) (*TemplateUtil, error) {
	tmpl, err := template.ParseGlob(filepath.Join(templateDir, "*.html"))
	if err != nil {
		return nil, err
	}
	return &TemplateUtil{templates: tmpl}, nil
}

// Render renders a template with the given data.
func (tu *TemplateUtil) Render(w http.ResponseWriter, templateName string, data interface{}) error {
	return tu.templates.ExecuteTemplate(w, templateName, data)
}
