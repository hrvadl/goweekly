package main

import (
	"log/slog"
	"os"

	"github.com/hrvadl/goweekly/translator/internal/app"
	"github.com/hrvadl/goweekly/translator/internal/cfg"
)

func main() {
	l := setupLogger().With(slog.Int("pid", os.Getpid()), slog.String("source", "translator"))
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
