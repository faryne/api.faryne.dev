package config

type Core struct {
	Port           string `env:"PORT"`
	TelegraphToken string `env:"TELEGRAPH_TOKEN"`
}

var Config Core
