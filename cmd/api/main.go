package main

import (
	"github.com/Eretic431/datingTelegramBot/internal/data/postgres"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/xlab/closer"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"log"
)

type application struct {
	config  *config
	log     *zap.SugaredLogger
	users   *postgres.UserRepository
	bot     *tgbotapi.BotAPI
	updates tgbotapi.UpdatesChannel
}

func main() {
	app, cleanup, err := initApp()
	if err != nil {
		log.Fatal("could not init application", err)
	}
	closer.Bind(func() {
		log.Print("stopping server")
		cleanup()
	})

	_, err = zap.NewStdLogAt(app.log.Desugar(), zap.ErrorLevel)
	if err != nil {
		app.log.Fatalw("could not init server logger", "err", err)
	}

	app.handleUpdates()
}

func newLogger(c *config) (*zap.SugaredLogger, func(), error) {
	var logger *zap.Logger
	var err error

	if c.Production {
		logger, err = zap.NewProduction()
	} else {
		conf := zap.NewDevelopmentConfig()
		conf.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		logger, err = conf.Build()
	}

	if err != nil {
		return nil, nil, err
	}

	cleanup := func() {
		_ = logger.Sync()
	}

	return logger.Sugar(), cleanup, nil
}

func newPostgresConfig(c *config, logger *zap.SugaredLogger) *postgres.Config {
	return &postgres.Config{
		PostgresUrl: c.PostgresUrl,
		Logger:      logger,
	}
}
