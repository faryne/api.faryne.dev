package config

type Core struct {
	Port           string `env:"PORT"`
	Environmet     string `env:"ENVIRONMENT"`
	TelegraphToken string `env:"TELEGRAPH_TOKEN"`
	Pixiv          struct {
		Username string `env:"NEKOMAID_PIXIV_USERNAME"`
		Password string `env:"NEKOMAID_PIXIV_PASSWORD"`
	}
	Tinami struct {
		ApiKey string `env:"NEKOMAID_TINAMI_APIKEY"`
	}
	AWS struct {
		Key    string `env:"AWS_KEY"`
		Secret string `env:"AWS_SECRET"`
	}
}

var Config Core
