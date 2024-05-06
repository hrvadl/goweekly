package app

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"golang.org/x/sync/errgroup"

	"github.com/hrvadl/goweekly/crawler/internal/cfg"
	"github.com/hrvadl/goweekly/crawler/internal/clients/rabbitmq/article"
	"github.com/hrvadl/goweekly/crawler/internal/crawler"
)

type App struct {
	cfg       cfg.Config
	crawler   crawler.Crawler
	publisher *article.Publisher
	log       *slog.Logger
}

func New(cfg cfg.Config, l *slog.Logger) *App {
	return &App{
		cfg:       cfg,
		crawler:   *crawler.New(time.Second*30, 10),
		log:       l,
		publisher: article.NewPublisher(),
	}
}

func (a *App) MustRun() {
	if err := a.Run(); err != nil {
		panic(err)
	}
}

func (a *App) Run() error {
	if err := a.publisher.Connect(a.cfg.RabbitMQAddr); err != nil {
		return fmt.Errorf("failed to connect to rabbitmq broker: %w", err)
	}
	defer a.publisher.Close()

	articles, err := a.crawler.ParseArticles()
	if err != nil {
		return err
	}

	g := new(errgroup.Group)
	for _, art := range articles {
		g.Go(func() error {
			ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
			defer cancel()
			return a.publisher.Publish(ctx, art)
		})
	}

	return g.Wait()
}
