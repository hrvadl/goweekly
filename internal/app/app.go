package app

import (
	"time"

	"github.com/hrvadl/go-weekly/internal/adapter"
	"github.com/hrvadl/go-weekly/internal/crawler"
	"github.com/hrvadl/go-weekly/internal/tg"
	"github.com/hrvadl/go-weekly/internal/tg/formatter"
	"github.com/hrvadl/go-weekly/internal/translator"
	"github.com/hrvadl/go-weekly/pkg/logger"
)

type Config struct {
	TranslateBatchRequests int
	TranslateRetries       int
	TranslateTimeout       time.Duration
	TranslateInterval      time.Duration

	ArticlesRetries int
	ArticlesTimeout time.Duration

	TgToken  string
	TgChatID string
}

func NewWithDefaults(token, chatID string) *GoWeekly {
	cfg := Config{
		TranslateBatchRequests: 7,
		TranslateRetries:       5,
		TranslateTimeout:       10 * time.Second,
		TranslateInterval:      10 * time.Second,
		ArticlesRetries:        3,
		ArticlesTimeout:        30 * time.Second,
		TgToken:                token,
		TgChatID:               chatID,
	}
	return &GoWeekly{cfg}
}

func New(cfg Config) *GoWeekly {
	return &GoWeekly{cfg}
}

type GoWeekly struct {
	cfg Config
}

func (o GoWeekly) TranslateAndSend() {
	start := time.Now()
	crawler := crawler.New(o.cfg.ArticlesTimeout, o.cfg.ArticlesRetries)
	fmt := formatter.NewMarkdown()
	bot := tg.NewBot(tg.BotConfig{
		ParseMode: formatter.MarkdownType,
		Token:     o.cfg.TgToken,
		ChatID:    o.cfg.TgChatID,
	})
	translator := translator.NewLingvaClient(&translator.Config{
		Timeout:         o.cfg.TranslateTimeout,
		Retries:         o.cfg.TranslateRetries,
		RetriesInterval: o.cfg.TranslateInterval / 2,
		BatchRequests:   o.cfg.TranslateBatchRequests,
		BatchInterval:   o.cfg.TranslateInterval,
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

	adp := adapter.NewArticleSender(bot, fmt)
	adp.SendWeekly(articles)
	logger.Info("Finished sending all the weekly articles")
}
