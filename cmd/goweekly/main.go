package main

import (
	"os"
	"time"

	"github.com/hrvadl/go-weekly/internal/app"
	"github.com/hrvadl/go-weekly/internal/crawler"
	"github.com/hrvadl/go-weekly/internal/translator"
	"github.com/hrvadl/go-weekly/pkg/logger"
)

const (
	tgTokenKey = "TG_TOKEN"
	tgChatID   = "@goweeklych"
)

func main() {
	appConfig := app.Config{
		TgToken:  os.Getenv(tgTokenKey),
		TgChatID: tgChatID,
	}
	app := app.New(appConfig)

	// adding applicaiton articles modules
	crawler := crawler.New(crawler.Config{
		Retries: 3,
		Timeout: 30 * time.Second,
	})
	app.AddWriter(crawler)
	translator := translator.NewLingvaClient(translator.Config{
		Timeout:         10 * time.Second,
		Retries:         5,
		RetriesInterval: 10 * time.Second / 2,
		BatchInterval:   10 * time.Second,
		BatchRequests:   7,
	})
	app.AddRedactor(translator)

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
