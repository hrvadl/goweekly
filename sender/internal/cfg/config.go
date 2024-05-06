package cfg

import (
	"errors"
	"os"
)

const (
	TokenKey     = "SENDER_TG_TOKEN"
	ChatIDKey    = "SENDER_CHAT_ID"
	ParseModeKey = "SENDER_PARSE_MODE"
	PortKey      = "SENDER_PORT"
)

type Config struct {
	Token     string
	ChatID    string
	ParseMode string
	Port      string
}

func Must(cfg *Config, err error) *Config {
	if err != nil {
		panic(err)
	}
	return cfg
}

func New() (*Config, error) {
	token := os.Getenv(TokenKey)
	if token == "" {
		return nil, errors.New("token can not be empty")
	}

	chatID := os.Getenv(ChatIDKey)
	if chatID == "" {
		return nil, errors.New("chatID can not be empty")
	}

	parseMode := os.Getenv(ParseModeKey)
	if parseMode == "" {
		return nil, errors.New("parse mode can not be empty")
	}

	port := os.Getenv(PortKey)
	if port == "" {
		return nil, errors.New("port can not be empty")
	}

	return &Config{
		Token:     token,
		ChatID:    chatID,
		ParseMode: parseMode,
		Port:      port,
	}, nil
}
