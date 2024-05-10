package app

import (
	"context"
	"log/slog"
	"time"

	"github.com/hrvadl/goweekly/core/internal/cfg"
	"github.com/hrvadl/goweekly/core/internal/clients/grpc/sender"
	"github.com/hrvadl/goweekly/core/internal/clients/grpc/translator"
	"github.com/hrvadl/goweekly/core/internal/clients/rabbitmq/article"
	"github.com/hrvadl/goweekly/core/internal/platform/formatter"
	"github.com/hrvadl/goweekly/core/internal/processor"
)

type App struct {
	consumer *article.Consumer
	log      *slog.Logger
	cfg      cfg.Config
}

func Must(a *App, err error) *App {
	if err != nil {
		panic(err)
	}
	return a
}

func New(cfg cfg.Config, log *slog.Logger) (*App, error) {
	fmter := formatter.NewMarkdown()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	sender, err := sender.New(ctx, cfg.SenderAddr, log.With("source", "sender client"))
	if err != nil {
		return nil, err
	}

	ctx, cancel = context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	translator, err := translator.New(
		ctx,
		cfg.TranslatorAddr,
		log.With("source", "translator client"),
	)
	if err != nil {
		return nil, err
	}

	consumer := article.NewConsumer(log, processor.New(fmter, sender, translator))
	consumer.Connect(cfg.RabbitMQAddr)

	return &App{
		consumer: consumer,
		cfg:      cfg,
	}, nil
}

func (a *App) MustRun() {
	if err := a.Run(); err != nil {
		panic(err)
	}
}

func (a *App) Run() error {
	block := make(chan struct{})
	a.consumer.Connect(a.cfg.RabbitMQAddr)
	if err := a.consumer.Consume(); err != nil {
		return err
	}
	defer a.consumer.Close()

	<-block
	return nil
}
