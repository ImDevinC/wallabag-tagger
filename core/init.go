package core

import (
	_ "github.com/joho/godotenv/autoload"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	WallabagUrl    string `envconfig:"WALLABAG_URL"`
	ClientID       string `envconfig:"CLIENT_ID"`
	ClientSecret   string `envconfig:"CLIENT_SECRET"`
	Username       string `envconfig:"USERNAME"`
	Password       string `envconfig:"PASSWORD"`
	GoogleAIApiKey string `envconfig:"GOOGLE_AI_API_KEY"`
	Ollama         struct {
		Model string `envconfig:"OLLAMA_MODEL"`
		URL   string `envconfig:"OLLAMA_URL"`
	}
}

func GetConfigFromEnv() (Config, error) {
	var config Config
	if err := envconfig.Process("", &config); err != nil {
		return Config{}, err
	}
	return config, nil
}
