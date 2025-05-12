package templates

import (
	"fmt"
	"html/template"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var (
	Core *TemplateManager
)

// TemplateManager 管理多个模板对象
type TemplateManager struct {
	templates map[string]*template.Template
}

type TemplateConfig struct {
	Name string `json:"name" yaml:"name" toml:"name"` // 模板名称
	Path string `json:"path" yaml:"path" toml:"path"` // 模板路径
}

// InitTemplates 创建一个新的 TemplateManager 实例
func InitTemplates(configs []TemplateConfig) *TemplateManager {
	Core = &TemplateManager{
		templates: make(map[string]*template.Template),
	}

	for _, config := range configs {
		// 加载模板, 路径统一改成绝对路径
		absPath, err := filepath.Abs(config.Path)
		if err != nil {
			log.Fatalf("load templates failed, %s", err)
		}
		log.Printf("load templates from %s, name: %s", absPath, config.Name)

		// 检查目录是否存在
		if _, err := os.Stat(absPath); os.IsNotExist(err) {
			log.Panicf("[Warning] templates directory %s does not exist", absPath)
			continue
		}

		Core.loadTemplatesFromDir(config.Name, absPath)
	}
	return Core
}

// LoadTemplatesFromDir 从指定目录加载模板，包括子目录
func (tm *TemplateManager) loadTemplatesFromDir(name, dir string) error {
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

// Render 渲染指定模板, name 板对象，, templateName 模板名称, data 模板数据
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
