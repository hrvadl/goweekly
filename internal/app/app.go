package app

import (
	"time"

	"github.com/hrvadl/go-weekly/internal/crawler"
	"github.com/hrvadl/go-weekly/internal/tg"
	"github.com/hrvadl/go-weekly/internal/tg/formatter"
	"github.com/hrvadl/go-weekly/pkg/logger"
)

type Config struct {
	TgToken  string
	TgChatID string
}

func New(cfg Config) *GoWeekly {
	return &GoWeekly{
		cfg,
		make([]Writer, 0),
		make([]Redactor, 0),
	}
}

type GoWeekly struct {
	cfg       Config
	writers   []Writer
	redactors []Redactor
}

func (o *GoWeekly) AddRedactor(r Redactor) {
	o.redactors = append(o.redactors, r)
}

func (o *GoWeekly) AddWriter(w Writer) {
	o.writers = append(o.writers, w)
}

func (o GoWeekly) TranslateAndSend() {
	start := time.Now()
	bot := tg.NewBot(o.cfg.TgToken, o.cfg.TgChatID)
	formatter := formatter.NewMarkdown()

	articles, err := o.collectArticles()
	if err != nil {
		logger.Fatalf("Cannot parse articles: %v\n", err)
	}

	logger.Infof(
		"Successfully parsed articles in %v: %v\n",
		time.Since(start).String(),
		articles,
	)

	if articles, err = o.processArticles(articles); err != nil {
		logger.Fatalf("Failed to process articles: %v\n", err)
	}

	logger.Infof(
		"Successfully translated articles in %v: %v\n",
		time.Since(start).String(),
		articles,
	)

	bot.SendMessagesThroughoutWeek(formatter.FormatArticles(articles))
	logger.Info("Finished sending all the weekly articles")
}

func (o GoWeekly) collectArticles() ([]crawler.Article, error) {
	result := make([]crawler.Article, 0)

	for i := 0; i < len(o.writers); i++ {
		arts, err := o.writers[i].GetArticles()
		if err != nil {
			return nil, err
		}
		result = append(result, arts...)
	}

	return result, nil
}

func (o GoWeekly) processArticles(articles []crawler.Article) ([]crawler.Article, error) {
	var err error
	for i := 0; i < len(o.redactors); i++ {
		articles, err = o.redactors[i].Review(articles)
		if err != nil {
			return nil, err
		}
	}

	return articles, nil
}
