package config

import "os"

type Config struct {
	Momo struct {
		ConsumerKey    string
		ConsumerSecret string
		TokenURL       string
		CallbackURL    string
		ApiEndpoint    string
	}
}

func LoadConfig() (*Config, error) {
	var cfg Config

	cfg.Momo.ConsumerKey = os.Getenv("MOMO_CONSUMER_KEY")
	cfg.Momo.ConsumerSecret = os.Getenv("MOMO_CONSUMER_SECRET")
	cfg.Momo.TokenURL = os.Getenv("MOMO_TOKEN_URL")
	cfg.Momo.CallbackURL = os.Getenv("MOMO_CALLBACK_URL")
	cfg.Momo.ApiEndpoint = os.Getenv("MOMO_API_ENDPOINT")

	return &cfg, nil
}
