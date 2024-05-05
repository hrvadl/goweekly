package app

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"time"

	"github.com/hrvadl/goweekly/protos/gen/go/v1/sender"
	"github.com/hrvadl/goweekly/protos/gen/go/v1/translator"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/hrvadl/goweekly/crawler/internal/cfg"
	"github.com/hrvadl/goweekly/crawler/internal/crawler"
	"github.com/hrvadl/goweekly/crawler/internal/formatter"
)

type App struct {
	cfg     cfg.Config
	crawler crawler.Crawler
	log     *slog.Logger
}

func New(cfg cfg.Config, l *slog.Logger) *App {
	return &App{
		cfg:     cfg,
		crawler: *crawler.New(time.Second*30, 10),
		log:     l,
	}
}

func (a *App) MustRun() {
	if err := a.Run(); err != nil {
		panic(err)
	}
}

func (a *App) Run() error {
	a.log.Info("Connecting to the translator service...", "addr", a.cfg.TranslatorAddr)
	translatorConn, err := grpc.Dial(
		a.cfg.TranslatorAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return fmt.Errorf("failed to connect to translator GRPC server: %w", err)
	}

	defer translatorConn.Close()
	translatorClient := translator.NewTranslateServiceClient(translatorConn)

	a.log.Info("Connecting to the sender service...", "addr", a.cfg.SenderAddr)
	senderConn, err := grpc.Dial(
		a.cfg.SenderAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return fmt.Errorf("failed to connect to sender GRPC server: %w", err)
	}

	defer senderConn.Close()
	senderClient := sender.NewSenderServiceClient(senderConn)

	articles, err := a.crawler.ParseArticles()
	if err != nil {
		return err
	}

	a.log.Debug("Parsed articles", "articles", articles)

	for i, article := range articles {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()
		res, err := translatorClient.Translate(
			ctx,
			&translator.TranslateRequest{Message: article.Content},
		)
		if err != nil {
			return fmt.Errorf("failed to translate article: %w", err)
		}

		a.log.Debug("Translated article", "content", res.Message)
		articles[i].Content = res.Message
	}

	ctx, sCancel := context.WithTimeout(context.Background(), time.Second*10)
	defer sCancel()
	stream, err := senderClient.Send(ctx)
	if err != nil {
		return fmt.Errorf("failed to create article stream: %w", err)
	}

	fmter := formatter.NewMarkdown()
	for _, a := range articles {
		err := stream.Send(&sender.SendRequest{
			Message: fmter.FormatArticle(a),
		})
		if err != nil {
			return err
		}
	}

	_, err = stream.CloseAndRecv()
	if err != nil && !errors.Is(err, io.EOF) {
		return fmt.Errorf("failed to close stream: %w", err)
	}

	return nil
}
