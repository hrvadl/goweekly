package app

import (
	"fmt"
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
}

func New(cfg cfg.Config) *App {
	lingva := lingva.NewClient(&lingva.Config{
		BatchRequests:   7,
		Retries:         5,
		Timeout:         10 * time.Second,
		RetriesInterval: 10 * time.Second,
		BatchInterval:   10 * time.Second,
	})
	srv := grpc.NewServer()
	translator.Register(srv, lingva)

	return &App{
		cfg:        cfg,
		gRPCServer: srv,
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

	if err := a.gRPCServer.Serve(l); err != nil {
		return fmt.Errorf("failed to server grpc translator server: %w", err)
	}

	return nil
}

func (a *App) Stop() {}
