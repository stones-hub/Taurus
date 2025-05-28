package internal

import (
	"Taurus/internal/controller"

	"github.com/google/wire"
)

var Core *Injector

type Injector struct {
	ValidateCtrl *controller.ValidateCtrl
	TraceCtrl    *controller.TraceCtrl
	MidCtrl      *controller.MidCtrl
	ConsulCtrl   *controller.ConsulCtrl
}

// Injector is the injector for the internal package
var InjectorSet = wire.NewSet(wire.Struct(new(Injector), "*"),
	controller.ValidateCtrlSet,
	controller.TraceCtrlSet,
	controller.MidCtrlSet,
	controller.ConsulCtrlSet,
)
