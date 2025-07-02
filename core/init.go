package core

import (
	_ "github.com/joho/godotenv/autoload"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	WallabagUrl    string `envconfig:"WT_WALLABAG_URL"`
	ClientID       string `envconfig:"WT_WALLABAG_CLIENT_ID"`
	ClientSecret   string `envconfig:"WT_WALLABAG_CLIENT_SECRET"`
	Username       string `envconfig:"WT_WALLABAG_USERNAME"`
	Password       string `envconfig:"WT_WALLABAG_PASSWORD"`
	GoogleAIApiKey string `envconfig:"WT_GOOGLE_AI_API_KEY"`
	Ollama         struct {
		Model string `envconfig:"WT_OLLAMA_MODEL"`
		URL   string `envconfig:"WT_OLLAMA_URL"`
	}
}

func GetConfigFromEnv() (Config, error) {
	var config Config
	if err := envconfig.Process("", &config); err != nil {
		return Config{}, err
	}
	return config, nil
}
