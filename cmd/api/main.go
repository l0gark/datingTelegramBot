package main

import (
	"github.com/Eretic431/datingTelegramBot/internal"
	"github.com/Eretic431/datingTelegramBot/internal/data/postgres"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/xlab/closer"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"log"
	"net/http"
)

type application struct {
	config  *config
	usecase internal.Usecase
	log     *zap.SugaredLogger
	users   *postgres.UserRepository
	likes   *postgres.LikeRepository
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

	go func() {
		http.HandleFunc("/deleteAll", app.deleteAllHandler)

		if err = http.ListenAndServe(":8090", nil); err != nil {
			log.Fatalf("couldn't start listen and serve with err = %e", err)
			return
		}
	}()

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

func (a *application) deleteAllHandler(w http.ResponseWriter, r *http.Request) {
	a.usecase.DeleteAll(r.Context())
	a.log.Info("deleting completed")
	w.WriteHeader(http.StatusOK)
}
