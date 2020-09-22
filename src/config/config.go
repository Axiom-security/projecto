package config

import (
	"flag"
	"fmt"
	"io/ioutil"
	"projecto/app"

	"gopkg.in/yaml.v2"
)

const ComponentName = "config"

var (
	config = flag.String("c", "", "path to yaml config filename")
)

type Config struct {
	Addr string `yaml:"addr"`
}

func New() *Config {
	return &Config{}
}

func (c *Config) Setup(a *app.App) (err error) {
	flag.Parse()
	source, err := ioutil.ReadFile(*config)
	if err = yaml.Unmarshal(source, &c); err != nil {
		return fmt.Errorf("can't parse config file '%v': %v", *config, err)
	}
	return nil
}

func (c *Config) Name() (name string) {
	return ComponentName
}

func (c *Config) GetAddr() string {
	return c.Addr
}
