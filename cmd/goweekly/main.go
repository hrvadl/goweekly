package main

import (
	"os"
	"time"

	"github.com/hrvadl/go-weekly/internal/app"
	"github.com/hrvadl/go-weekly/pkg/logger"
)

const (
	tgTokenKey = "TG_TOKEN"
	tgChatID   = "@goweeklych"
)

func main() {
	app := app.NewWithDefaults(os.Getenv(tgTokenKey), tgChatID)
	app.TranslateAndSend()

	location, err := time.LoadLocation("UTC")
	if err != nil {
		logger.Fatalf("Error loading time zone: %v", err)
	}

	ticker := time.NewTicker(24 * time.Hour)
	for range ticker.C {
		now := time.Now().In(location)
		if now.Weekday() == time.Wednesday {
			app.TranslateAndSend()
		}
	}
}
