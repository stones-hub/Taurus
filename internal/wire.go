//go:build wireinject
// +build wireinject

package internal

import (
	"github.com/google/wire"
)

func BuildInjector() (*Injector, func(), error) {
	panic(
		wire.Build(InjectorSet),
	)
	return new(Injector), func() {}, nil
}
