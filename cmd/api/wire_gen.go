// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//+build !wireinject

package main

import (
	"github.com/Eretic431/datingTelegramBot/internal/data/postgres"
)

// Injectors from wire.go:

func initApp() (*application, func(), error) {
	mainConfig, err := getConfig()
	if err != nil {
		return nil, nil, err
	}
	sugaredLogger, cleanup, err := newLogger(mainConfig)
	if err != nil {
		return nil, nil, err
	}
	postgresConfig := newPostgresConfig(mainConfig, sugaredLogger)
	pool, cleanup2, err := postgres.NewPsqlPool(postgresConfig)
	if err != nil {
		cleanup()
		return nil, nil, err
	}
	userRepository := &postgres.UserRepository{
		DB: pool,
	}
	likeRepository := &postgres.LikeRepository{
		DB: pool,
	}
	botAPI, err := newTgBot(mainConfig)
	if err != nil {
		cleanup2()
		cleanup()
		return nil, nil, err
	}
	updatesChannel := newTgBotUpdatesChan(botAPI)
	mainApplication := &application{
		config:  mainConfig,
		log:     sugaredLogger,
		users:   userRepository,
		likes:   likeRepository,
		bot:     botAPI,
		updates: updatesChannel,
	}
	return mainApplication, func() {
		cleanup2()
		cleanup()
	}, nil
}
