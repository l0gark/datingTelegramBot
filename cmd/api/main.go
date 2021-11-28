package main

import (
	"github.com/xlab/closer"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"log"
)

type application struct {
	config *config
	log    *zap.SugaredLogger
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
