package config

import (
	"os"

	"github.com/go-playground/validator/v10"
)

type Config struct {
	TelegramToken string `validate:"required"`
	MongoURI      string `validate:"required"`
	OtelURI       string
}

func LoadConfig() (*Config, error) {
	telegramToken := os.Getenv("TELEGRAM_TOKEN")
	mongoURI := os.Getenv("MONGO_URI")
	otelURI := os.Getenv("OTEL_URI")

	config := &Config{
		TelegramToken: telegramToken,
		MongoURI:      mongoURI,
		OtelURI:       otelURI,
	}

	validate := validator.New()
	if err := validate.Struct(config); err != nil {
		return nil, err
	}

	return config, nil
}
