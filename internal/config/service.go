package config

import (
	"github.com/go-playground/validator/v10"
	"os"
)

type Config struct {
	ImapServer string `validate:"required"`
	Username   string `validate:"required"`
	Password   string `validate:"required"`
}

func LoadConfig() (*Config, error) {
	imapServer := os.Getenv("IMAP_SERVER")
	username := os.Getenv("IMAP_USERNAME")
	password := os.Getenv("IMAP_PASSWORD")

	config := &Config{
		ImapServer: imapServer,
		Username:   username,
		Password:   password,
	}

	validate := validator.New()
	if err := validate.Struct(config); err != nil {
		return nil, err
	}

	return config, nil
}
