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

package telemetry

import (
	"log"

	"go.opentelemetry.io/otel/trace"
)

// 注册配置的调用链

var tracerRegistry = make(map[string]trace.Tracer)

func RegisterTracer(name string, tracer trace.Tracer) {
	if _, exists := tracerRegistry[name]; exists {
		log.Printf("Tracer %s already registered", name)
	}
	tracerRegistry[name] = tracer
}

func GetTracer(name string) trace.Tracer {
	if tracer, exists := tracerRegistry[name]; exists {
		return tracer
	}

	return Provider.Tracer("default")
}
