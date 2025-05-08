package controller

import (
	"Taurus/pkg/httpx"
	"Taurus/pkg/loggerx"
	"net/http"

	"github.com/google/wire"
)

type DemoCtrl struct {
}

var DemoCtrlSet = wire.NewSet(wire.Struct(new(DemoCtrl), "*"))

func (c *DemoCtrl) Get(w http.ResponseWriter, r *http.Request) {
	data, _ := httpx.ParseJson(r)
	loggerx.DefaultLogger.Info("demo get %v", data)
	httpx.SendResponse(w, http.StatusOK, data, nil)
}
