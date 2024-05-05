package cfg

import (
	"errors"
	"os"
)

const (
	TranslatorAddrKey = "TRANSLATOR_ADDR"
	SenderAddrKey     = "SENDER_ADDR"
)

type Config struct {
	TranslatorAddr string
	SenderAddr     string
}

func Must(c *Config, err error) *Config {
	if err != nil {
		panic(err)
	}
	return c
}

func New() (*Config, error) {
	translatorAddr := os.Getenv(TranslatorAddrKey)
	if translatorAddr == "" {
		return nil, errors.New("translator service addr cannot be empty")
	}

	senderAddr := os.Getenv(SenderAddrKey)
	if translatorAddr == "" {
		return nil, errors.New("sender service addr cannot be empty")
	}

	return &Config{
		TranslatorAddr: translatorAddr,
		SenderAddr:     senderAddr,
	}, nil
}
