package config

import (
	"os"
)

type Config struct {
	ApiKey string
}

func LoadConfig() Config {
	return Config{
		ApiKey: os.Getenv("GETBLOCK_API_KEY"),
	}
}
