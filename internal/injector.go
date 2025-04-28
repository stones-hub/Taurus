package internal

import (
	"Taurus/internal/controller"
	"Taurus/internal/controller/ws"

	"github.com/google/wire"
)

type Injector struct {
	DemoCtrl *controller.DemoCtrl
	DemoWs   *ws.DemoWs
}

// Injector is the injector for the internal package
var InjectorSet = wire.NewSet(wire.Struct(new(Injector), "*"), controller.DemoCtrlSet, ws.DemoWsSet)
