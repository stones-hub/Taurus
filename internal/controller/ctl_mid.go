package controller

import (
	"Taurus/pkg/httpx"
	"log"
	"net/http"

	"github.com/google/wire"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type MidCtrl struct {
}

var MidCtrlSet = wire.NewSet(wire.Struct(new(MidCtrl), "*"))

func (midCtrl *MidCtrl) TestMid(w http.ResponseWriter, r *http.Request) {
	log.Println("-------------------------------- TestMid --------------------------------")

	span := trace.SpanFromContext(r.Context())
	span.SetAttributes(attribute.String("user.id", "123"))
	span.SetAttributes(attribute.String("user.name", "TestMid-Name"))

	httpx.SendResponse(w, http.StatusOK, "mid test", nil)
}
