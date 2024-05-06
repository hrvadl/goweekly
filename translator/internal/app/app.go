package app

import (
	"fmt"
	"log/slog"
	"net"
	"time"

	"google.golang.org/grpc"

	"github.com/hrvadl/goweekly/translator/internal/cfg"
	"github.com/hrvadl/goweekly/translator/internal/grpc/translator"
	"github.com/hrvadl/goweekly/translator/internal/platform/translator/lingva"
)

type App struct {
	cfg        cfg.Config
	gRPCServer *grpc.Server
	log        *slog.Logger
}

func New(cfg cfg.Config, log *slog.Logger) *App {
	lingva := lingva.NewClient(&lingva.Config{
		Retries:         5,
		RetriesInterval: 10 * time.Second,
		Logger:          log,
	})
	srv := grpc.NewServer()
	translator.Register(srv, lingva, log)

	return &App{
		cfg:        cfg,
		gRPCServer: srv,
		log:        log,
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

	a.log.Info("Starting translator service", slog.String("port", a.cfg.Port))
	if err := a.gRPCServer.Serve(l); err != nil {
		return fmt.Errorf("failed to server grpc translator server: %w", err)
	}

	return nil
}

func (a *App) Stop() {}
