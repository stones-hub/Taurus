package internal

import (
	"Taurus/internal/controller"

	"github.com/google/wire"
)

type Injector struct {
	DemoCtrl *controller.DemoCtrl
}

// Injector is the injector for the internal package
var InjectorSet = wire.NewSet(wire.Struct(new(Injector), "*"), controller.DemoCtrlSet)
