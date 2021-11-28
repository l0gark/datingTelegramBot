package main

import "github.com/caarlos0/env"

type config struct {
	Production bool   `env:"PRODUCTION" envDefault:"false"`
	Port       string `env:"PORT" envDefault:"80"`
}

func getConfig() (*config, error) {
	c := &config{}
	if err := env.Parse(c); err != nil {
		return nil, err
	}

	return c, nil
}
