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
	"strconv"
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
		http.HandleFunc("/addTestUser", app.addTestUser)
		http.HandleFunc("/addTestUserWithLike", app.addTestUserWithLike)

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
	_ = a.usecase.DeleteAll(r.Context())
	a.log.Info("deleting completed")
	w.WriteHeader(http.StatusOK)
}

func (a *application) addTestUser(w http.ResponseWriter, r *http.Request) {
	sexStr := r.URL.Query().Get("sex")
	sex, err := strconv.ParseBool(sexStr)
	if err != nil {
		a.log.Errorf("couldn't parse query parametr sex = %s", sexStr)
	}
	_ = a.usecase.AddTestUser(r.Context(), sex)
	a.log.Info("test user added")
	w.WriteHeader(http.StatusOK)
}

func (a *application) addTestUserWithLike(w http.ResponseWriter, r *http.Request) {
	sexStr := r.URL.Query().Get("sex")
	sex, err := strconv.ParseBool(sexStr)
	if err != nil {
		a.log.Errorf("couldn't parse query parametr sex = %s", sexStr)
	}

	toId := r.URL.Query().Get("toId")

	_ = a.usecase.AddTestUserWithLike(r.Context(), sex, toId)
	a.log.Info("test user with like added")
	w.WriteHeader(http.StatusOK)
}
