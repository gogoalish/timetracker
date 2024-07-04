package config

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/pkg/errors"
)

type Config struct {
	DBURL  string
	APIURL string
	Host   string
	Port   string
}

func New() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, errors.Wrap(err, "error loading .env:")
	}

	return &Config{
		DBURL:  os.Getenv("DB_URL"),
		APIURL: os.Getenv("API_URL"),
		Host:   os.Getenv("HOST"),
		Port:   os.Getenv("PORT"),
	}, nil
}
