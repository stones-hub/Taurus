//go:build wireinject
// +build wireinject

package internal

import (
	"github.com/google/wire"
)

func BuildInjector() (*Injector, func(), error) {
	panic(
		// 构建injector
		wire.Build(InjectorSet),
	)
	return new(Injector), func() {}, nil
}
