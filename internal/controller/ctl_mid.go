package controller

import (
	"Taurus/pkg/httpx"
	"Taurus/pkg/logx"
	"net/http"

	"github.com/google/wire"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

// Author: yelei
// Email: 61647649@qq.com
// Date: 2025-06-13

type MidCtrl struct {
}

var MidCtrlSet = wire.NewSet(wire.Struct(new(MidCtrl), "*"))

func (midCtrl *MidCtrl) TestMid(w http.ResponseWriter, r *http.Request) {

	span := trace.SpanFromContext(r.Context())
	span.SetAttributes(attribute.String("user.id", "123"))
	span.SetAttributes(attribute.String("user.name", "TestMid-Name"))

	logx.Core.Info("default", "this is mid test")
	httpx.SendResponse(w, http.StatusOK, "mid test", nil)
}
