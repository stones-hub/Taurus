package util

import (
	"fmt"
	"html/template"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var (
	Templates *TemplateManager
)

// TemplateManager 管理多个模板对象
type TemplateManager struct {
	templates map[string]*template.Template
}

// NewTemplateManager 创建一个新的 TemplateManager 实例
func NewTemplateManager() *TemplateManager {
	return &TemplateManager{
		templates: make(map[string]*template.Template),
	}
}

// LoadTemplatesFromDir 从指定目录加载模板，包括子目录
func (tm *TemplateManager) LoadTemplatesFromDir(name, dir string) error {
	tmpl := template.New(name)

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(info.Name(), ".html") {
			_, err := tmpl.ParseFiles(path)
			if err != nil {
				return err
			}
		}
		return nil
	})

	if err != nil {
		return err
	}

	tm.templates[name] = tmpl
	return nil
}

// AddTemplate 动态添加模板 name 模板对象， templateName 模板名称， content 模板内容
func (tm *TemplateManager) AddTemplate(name, templateName, content string) error {
	tmpl, exists := tm.templates[name]
	if !exists {
		tmpl = template.New(name)
		tm.templates[name] = tmpl
	}

	_, err := tmpl.New(templateName).Parse(content)
	return err
}

// Render 渲染指定模板
func (tm *TemplateManager) Render(name, templateName string, data interface{}) (string, error) {
	tmpl, exists := tm.templates[name]
	if !exists {
		return "", fmt.Errorf("template %s does not exist", name)
	}

	var sb strings.Builder
	err := tmpl.ExecuteTemplate(&sb, templateName, data)
	if err != nil {
		return "", err
	}

	return sb.String(), nil
}

func init() {
	Templates = NewTemplateManager()
	err := Templates.LoadTemplatesFromDir("default", "./templates")
	if err != nil {
		log.Printf("Failed to load templates: %v \n", err)
	}
}
