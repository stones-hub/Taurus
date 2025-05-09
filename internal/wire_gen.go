// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package internal

import (
	"Taurus/internal/controller"
	"Taurus/internal/controller/ws"
)

// Injectors from wire.go:

func BuildInjector() (*Injector, func(), error) {
	demoCtrl := &controller.DemoCtrl{}
	demoWs := &ws.DemoWs{}
	injector := &Injector{
		DemoCtrl: demoCtrl,
		DemoWs:   demoWs,
	}
	return injector, func() {
	}, nil
}
