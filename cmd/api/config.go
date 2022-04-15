package main

import "github.com/caarlos0/env"

type config struct {
	Production  bool   `env:"PRODUCTION" envDefault:"false"`
	Port        string `env:"PORT" envDefault:"80"`
	PostgresUrl string `env:"POSTGRES_URL" envDefault:"postgres://postgres:tgbot@localhost:5432/postgres?sslmode=disable"`
	TgBotToken  string `env:"BOT_TOKEN" envDefault:"5085296878:AAFm-DFK-_S0Eg7GpXhkq0S6XN1W9BhLVLU"`
}

func getConfig() (*config, error) {
	c := &config{}
	if err := env.Parse(c); err != nil {
		return nil, err
	}

	return c, nil
}
