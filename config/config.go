package config

import (
	"errors"
	"os"
)

type Config struct {
	ApiKey string
}

func LoadConfig() (Config, error) {
	key := os.Getenv("GETBLOCK_API_KEY")
	if key == "" {
		return Config{}, errors.New("missing API key")
	}
	return Config{
		ApiKey: key,
	}, nil
}
