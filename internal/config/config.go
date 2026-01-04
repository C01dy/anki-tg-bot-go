package config

import (
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	PORT          string `env:"PORT"`
	DB            string `env:"DB"`
	TelegramToken string `env:"TELEGRAM_TOKEN"`
}

func GetConfig(path string) (*Config, error) {
	var cfg Config

	err := cleanenv.ReadConfig(path, &cfg)
	if err != nil {
		return nil, fmt.Errorf("couldn't find env file by %s path: %w", path, err)
	}

	return &cfg, nil
}
