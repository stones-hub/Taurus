// Copyright (c) 2025 Taurus Team. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Author: yelei
// Email: 61647649@qq.com
// Date: 2025-06-13

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

// InitTemplates
func InitTemplates(configs []TemplateConfig) *TemplateManager {
	Core = &TemplateManager{
		templates: make(map[string]*template.Template),
	}

	for _, config := range configs {

		if _, ok := Core.templates[config.Name]; ok {
			log.Printf("[Warning] template %s already exists", config.Name)
			continue
		}

		// add template, use absolute path
		absPath, err := filepath.Abs(config.Path)
		if err != nil {
			log.Fatalf("load templates failed, %s", err)
		}
		log.Printf("load templates from %s, name: %s", absPath, config.Name)

		// 检查目录是否存在
		if _, err := os.Stat(absPath); os.IsNotExist(err) {
			log.Printf("[Warning] templates directory %s does not exist", absPath)
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
