//+build wireinject

package main

import (
	"github.com/Eretic431/datingTelegramBot/internal/data/postgres"
	"github.com/google/wire"
)

func initApp() (*application, func(), error) {
	wire.Build(
		getConfig,
		newLogger,
		newPostgresConfig,
		postgres.NewPsqlPool,
		wire.Struct(new(postgres.UserRepository), "*"),
		wire.Struct(new(application), "*"),
	)

	return nil, nil, nil
}
