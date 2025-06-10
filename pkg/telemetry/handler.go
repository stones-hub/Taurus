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
