package internal

import (
	"Taurus/internal/controller"

	"github.com/google/wire"
)

var Core *Injector

type Injector struct {
	DemoCtrl *controller.DemoCtrl
}

// Injector is the injector for the internal package
var InjectorSet = wire.NewSet(wire.Struct(new(Injector), "*"), controller.DemoCtrlSet)
