package config

type Core struct {
	Port string `env:"PORT"`
}

var Config Core
