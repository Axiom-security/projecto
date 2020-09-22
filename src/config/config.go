package config

import (
	"projecto/app"
)

const ComponentName = "config"

type Config struct {
}

func New() *Config {
	return &Config{}
}

func (c *Config) Setup(a *app.App) (err error) {
	return nil
}

func (c *Config) Name() (name string) {
	return ComponentName
}
