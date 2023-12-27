package main

import (
	"os"
	"time"

	"github.com/hrvadl/go-weekly/internal/crawler"
	"github.com/hrvadl/go-weekly/internal/tg"
	"github.com/hrvadl/go-weekly/internal/tg/formatter"
	"github.com/hrvadl/go-weekly/internal/translator"
	"github.com/hrvadl/go-weekly/pkg/logger"
)

const (
	articlesURL     = "https://golangweekly.com/issues/latest"
	articlesRetries = 3
	articlesTimeout = 30 * time.Second
)

const (
	translateRetries       = 5
	translateBatchRequests = 7
	translateInterval      = 10 * time.Second
	translateTimeout       = 10 * time.Second
)

const (
	tgTokenKey = "TG_TOKEN"
	tgChatID   = "@goweeklych"
)

func main() {
	start := time.Now()
	crawler := crawler.New(articlesURL, articlesTimeout, articlesRetries)
	bot := tg.NewBot(tg.URL, os.Getenv(tgTokenKey), tgChatID)
	formatter := formatter.NewMarkdownV2()
	translator := translator.NewLingvaClient(&translator.Config{
		Timeout:         translateTimeout,
		Retries:         translateRetries,
		RetriesInterval: translateInterval / 2,
		BatchRequests:   translateBatchRequests,
		BatchInterval:   translateInterval,
		URL:             translator.LingvaAPIURL,
	})

	articles, err := crawler.ParseArticles()
	if err != nil {
		logger.Fatalf("Cannot parse articles: %v\n", err)
	}

	logger.Infof(
		"Successfully parsed articles in %v: %v\n",
		time.Since(start).String(),
		articles,
	)

	if err := translator.TranslateArticles(articles); err != nil {
		logger.Fatalf("Failed to translate articles: %v\n", err)
	}

	logger.Infof(
		"Successfully translated articles in %v: %v\n",
		time.Since(start).String(),
		articles,
	)

	formatted := formatter.FormatArticles(articles)
	bot.SendMessagesThroughoutWeek(formatted)
}
