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
	DB   struct {
		Host     string `yaml:"host"`
		Port     int    `yaml:"port"`
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		Database string `yaml:"database"`
	} `yaml:"db"`
}

func New() *Config {
	return &Config{}
}

func (c *Config) Setup(a *app.App) (err error) {
	flag.Parse()
	source, err := ioutil.ReadFile(*config)
	if err != nil {
		return fmt.Errorf("can't read config file '%v': %v", *config, err)
	}
	if err = yaml.Unmarshal(source, &c); err != nil {
		return fmt.Errorf("can't parse config file '%v': %v", *config, err)
	}
	return nil
}

func (c *Config) Name() (name string) {
	return ComponentName
}

func (c *Config) GetDB() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s dbname=%s password=%s sslmode=disable",
		c.DB.Host,
		c.DB.Port,
		c.DB.User,
		c.DB.Database,
		c.DB.Password,
	)
}

func (c *Config) GetAddr() string {
	return c.Addr
}
