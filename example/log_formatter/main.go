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

package main

import (
	"Taurus/pkg/logx"
	"fmt"
	"time"
)

// if you want to use a custom formatter, you can implement the Formatter interface, and register it to the logx.RegisterFormatter
// then you must set the Formatter name to the config/logger/logger.yaml
type DemoFormatter struct {
}

func (d *DemoFormatter) Format(level logx.LogLevel, file string, line int, message string) string {
	return fmt.Sprintf("[%s] [%s]: %s", time.Now().Format("2006-01-02 15:04:05"), logx.GetLevelSTR(level), message)
}

func init() {
	logx.RegisterFormatter("demo", &DemoFormatter{})
}

func main() {
}
