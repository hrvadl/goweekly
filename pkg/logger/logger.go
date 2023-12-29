package logger

import (
	"fmt"
	"log/slog"
	"os"
)

func Infof(msg string, args ...any) {
	slog.Info(fmt.Sprintf(msg, args...))
}

func Info(msg string) {
	slog.Info(msg)
}

func Debugf(msg string, args ...any) {
	slog.Debug(fmt.Sprintf(msg, args...))
}

func Debug(msg string) {
	slog.Debug(msg)
}

func Errorf(msg string, args ...any) {
	slog.Error(fmt.Sprintf(msg, args...))
}

func Error(msg string) {
	slog.Error(msg)
}

func Fatalf(msg string, args ...any) {
	Fatal(fmt.Sprintf(msg, args...))
}

func Fatal(msg string, args ...any) {
	Errorf(msg, args...)
	os.Exit(1)
}
