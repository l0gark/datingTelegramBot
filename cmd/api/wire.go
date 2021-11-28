//+build wireinject

package main

import "github.com/google/wire"

func initApp() (*application, func(), error) {
	wire.Build(
		getConfig,
		newLogger,
		wire.Struct(new(application), "*"),
	)

	return nil, nil, nil
}
