package internal

import (
	"Taurus/internal/ctrl"
	"Taurus/internal/ctrl/ws"

	"github.com/google/wire"
)

type Injector struct {
	DemoCtrl *ctrl.DemoCtrl
	DemoWs   *ws.DemoWs
}

// Injector is the injector for the internal package
var InjectorSet = wire.NewSet(wire.Struct(new(Injector), "*"), ctrl.DemoCtrlSet, ws.DemoWsSet)
