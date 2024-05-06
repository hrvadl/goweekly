package cfg

import (
	"errors"
	"os"
)

const PortKey = "TRANSLATOR_PORT"

type Config struct {
	Port string
}

func Must(c *Config, err error) *Config {
	if err != nil {
		panic(err)
	}
	return c
}

func New() (*Config, error) {
	port := os.Getenv(PortKey)
	if port == "" {
		return nil, errors.New("port cannot be empty")
	}
	return &Config{
		Port: port,
	}, nil
}
