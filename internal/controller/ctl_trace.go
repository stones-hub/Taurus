package controller

import (
	"Taurus/pkg/httpx"
	"Taurus/pkg/logx"
	"net/http"

	"github.com/google/wire"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

/*
Note: 测试trace中间件，记录请求的详细信息
*/

type TraceCtrl struct {
}

var TraceCtrlSet = wire.NewSet(wire.Struct(new(TraceCtrl), "*"))

func (traceCtrl *TraceCtrl) TestTraceMiddleware(w http.ResponseWriter, r *http.Request) {
	logx.Core.Info("default", "this is trace middleware test")

	// 获取当前 span
	span := trace.SpanFromContext(r.Context())
	span.SetAttributes(attribute.String("user.id", "123"))

	httpx.SendResponse(w, http.StatusOK, "trace middleware test", nil)

}
