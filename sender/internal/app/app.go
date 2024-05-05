package app

import (
	"fmt"
	"log/slog"
	"net"

	"google.golang.org/grpc"

	"github.com/hrvadl/goweekly/sender/internal/cfg"
	"github.com/hrvadl/goweekly/sender/internal/grpc/sender"
	"github.com/hrvadl/goweekly/sender/internal/platform/sender/tg"
)

type App struct {
	gRPCServer *grpc.Server
	cfg        cfg.Config
	log        *slog.Logger
}

func New(cfg cfg.Config, l *slog.Logger) *App {
	b := tg.NewBot(cfg.Token, cfg.ChatID, cfg.ParseMode)

	srv := grpc.NewServer()
	sender.Register(srv, b, l)

	return &App{
		gRPCServer: srv,
		cfg:        cfg,
		log:        l,
	}
}

func (a *App) MustRun() {
	if err := a.Run(); err != nil {
		panic(err)
	}
}

func (a *App) Run() error {
	l, err := net.Listen("tcp", net.JoinHostPort("", a.cfg.Port))
	if err != nil {
		return err
	}

	a.log.Info("Starting sender service", slog.String("port", a.cfg.Port))
	if err := a.gRPCServer.Serve(l); err != nil {
		return fmt.Errorf("grpc server failed to listen: %w", err)
	}

	return nil
}

func (a *App) Stop() {}
