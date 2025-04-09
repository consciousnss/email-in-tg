package config

import (
	"os"

	"github.com/go-playground/validator/v10"
)

type Config struct {
	ImapServer    string `validate:"required"`
	TelegramToken string `validate:"required"`
	MongoURI      string `validate:"required"`
}

func LoadConfig() (*Config, error) {
	imapServer := os.Getenv("IMAP_SERVER")
	telegramToken := os.Getenv("TELEGRAM_TOKEN")
	mongoURI := os.Getenv("MONGO_URI")

	config := &Config{
		ImapServer:    imapServer,
		TelegramToken: telegramToken,
		MongoURI:      mongoURI,
	}

	validate := validator.New()
	if err := validate.Struct(config); err != nil {
		return nil, err
	}

	return config, nil
}
