//go:build wireinject
// +build wireinject

package main

import (
	"github.com/Eretic431/datingTelegramBot/internal/data/postgres"
	"github.com/Eretic431/datingTelegramBot/internal/usecase"
	"github.com/google/wire"
)

func initApp() (*application, func(), error) {
	wire.Build(
		getConfig,
		newLogger,
		newPostgresConfig,
		postgres.NewPsqlPool,
		postgres.NewUserRepository,
		postgres.NewLikeRepository,
		wire.Struct(new(postgres.UserRepository), "*"),
		wire.Struct(new(postgres.LikeRepository), "*"),
		newTgBot,
		newTgBotUpdatesChan,
		usecase.NewUsecase,
		wire.Struct(new(application), "*"),
	)

	return nil, nil, nil
}
