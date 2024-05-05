package main

import (
	"log/slog"
	"os"

	"github.com/hrvadl/goweekly/crawler/internal/app"
	"github.com/hrvadl/goweekly/crawler/internal/cfg"
)

func main() {
	l := setupLogger().With(slog.Int("pid", os.Getpid()), slog.String("source", "sender"))
	l.Info("Starting the application...")

	cfg := *cfg.Must(cfg.New())
	l.Info("Successfuly parsed the configuration")

	app := app.New(cfg, l)
	app.MustRun()
}

func setupLogger() *slog.Logger {
	return slog.New(
		slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		}),
	)
}
